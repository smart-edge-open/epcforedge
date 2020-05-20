// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use: "apply",
	Short: "Apply LTE CUPS userplane or NGC AF TI subscription using YAML " +
		"configuration file",
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// Read file from the filename provided in command
		data, err := readInputData(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		switch c.Kind {
		case "ngc":

			var s AFTrafficInfluSub
			if err = yaml.Unmarshal(data, &s); err != nil {
				fmt.Println(err)
				return
			}

			var sub []byte
			sub, err = yaml.Marshal(s.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			sub, err = y2j.YAMLToJSON(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			// create new subscription
			var subLoc string
			subLoc, err = AFCreateSubscription(sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Subscription URI:", subLoc)

		case "lte":

			var u LTEUserplane
			if err = yaml.Unmarshal(data, &u); err != nil {
				fmt.Println(err)
				return
			}

			up, err := yaml.Marshal(u.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			up, err = y2j.YAMLToJSON(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			// create new LTE userplane
			upID, err := LteCreateUserplane(up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Userplane:", upID)

		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

// pfdApplyCmd represents the apply command
var pfdApplyCmd = &cobra.Command{
	Use: "apply",
	Short: "Apply NGC AF PFD Transaction or NGC AF PFD Application " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// Read file from the filename provided in command
		data, err := readInputData(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		if c.Kind != "ngc_pfd" {
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
			return
		}

		var s AFPfdManagement
		if err = yaml.Unmarshal(data, &s); err != nil {
			fmt.Println(err)
			return
		}

		pfdTransData := getPfdTransData(s)

		var trans []byte
		trans, err = json.Marshal(pfdTransData)
		if err != nil {
			fmt.Println(err)
			return
		}

		// create new AF PFD Transaction
		pfdData, self, err := AFCreatePfdTransaction(trans)
		if err != nil {
			klog.Info(err)
			if err.Error() == "HTTP failure: 500" && pfdData != nil {
				printPfdReports(pfdData)
			}
			return
		}

		if pfdData != nil {
			printPdfTransStatus(pfdData, self)
		} else {
			fmt.Printf("PFD Transaction URI: %s\n", self)
			fmt.Printf("PFD Transaction ID: %s\n",
				getTransIDFromURL(self))
		}
	},
}

//paApplyCmd represents apply command
var paApplyCmd = &cobra.Command{
	Use: "apply",
	Short: "Create NGC AF PCF application session context" +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// Read file from the filename provided in command
		data, err := readInputData(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		if c.Kind != "ngc_policy_authorization" {
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
			return
		}

		var s AFAscReqData
		if err = yaml.Unmarshal(data, &s); err != nil {
			fmt.Println(err)
			return
		}

		paAscReqData := getPaAscReqData(s)
		asc := AppSessionContext{}
		asc.AscReqData = &paAscReqData

		appSession, err := json.Marshal(asc)
		if err != nil {
			fmt.Println(err)
			return
		}

		// create new app-session
		appSessionRespData, appLoc, err := AFCreatePaAppSession(appSession)
		if err != nil {
			klog.Info(err)
			return
		}

		if appSessionRespData != nil {
			//print success response structure in yaml format
			printAscData(appSessionRespData)
		}

		appSessionID := getappSessionIDFromURL(appLoc)
		fmt.Printf("\nappSessionID: %s\n", appSessionID)

	},
}

func init() {

	const help = `Apply LTE CUPS userplane or NGC AF TI subscription
using YAML configuration file

Usage:
  cnca apply -f <config.yml>

Example:
  cnca apply -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	const pfdHelp = `Apply NGC AF PFD Transaction or NGC AF PFD Application
using YAML configuration file

Usage:
  cnca pfd apply -f <config.yml>

Example:
  cnca pfd apply -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	const paHelp = `Create NGC AF PCF application session context
using YAML configuration file

Usage:
  cnca policy-authorization apply -f <config.yml>

Example:
  cnca policy-authorization apply -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	// add `apply` command
	cncaCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = applyCmd.MarkFlagRequired("filename")
	applyCmd.SetHelpTemplate(help)

	// add pfd `apply` command
	pfdCmd.AddCommand(pfdApplyCmd)
	pfdApplyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = pfdApplyCmd.MarkFlagRequired("filename")
	pfdApplyCmd.SetHelpTemplate(pfdHelp)

	//add policy-authorization  (pa) `apply` command
	paCmd.AddCommand(paApplyCmd)
	paApplyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = paApplyCmd.MarkFlagRequired("filename")
	paApplyCmd.SetHelpTemplate(paHelp)
}

func readInputData(cmd *cobra.Command) ([]byte, error) {
	ymlFile, _ := cmd.Flags().GetString("filename")
	if ymlFile == "" {
		return nil, errors.New("YAML file missing")
	}

	data, err := ioutil.ReadFile(ymlFile)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getTransIDFromURL(url string) string {
	urlElements := strings.Split(url, "/")
	for index, str := range urlElements {
		if str == "transactions" {
			return urlElements[index+1]
		}
	}
	return ""
}

func getPfdTransData(inputPfdTransData AFPfdManagement) PfdManagement {

	var pfdTransData PfdManagement

	if inputPfdTransData.Policy.SuppFeat != nil {
		*pfdTransData.SuppFeat =
			SupportedFeatures(*inputPfdTransData.Policy.SuppFeat)
	}

	if inputPfdTransData.Policy.PfdDatas != nil {
		pfdTransData.PfdDatas = make(map[string]PfdData)
	}

	for _, inputPfdAppData := range inputPfdTransData.Policy.PfdDatas {
		var pfdAppData PfdData

		pfdAppData.ExternalAppID = inputPfdAppData.ExternalAppID
		pfdAppData.Self = Link(inputPfdAppData.Self)

		if inputPfdAppData.AllowedDelay != nil {
			allowedDelay := DurationSecRm(*inputPfdAppData.AllowedDelay)
			pfdAppData.AllowedDelay = &allowedDelay
		}
		if inputPfdAppData.CachingTime != nil {
			cachingTime := DurationSecRo(*inputPfdAppData.CachingTime)
			pfdAppData.CachingTime = &cachingTime
		}
		if inputPfdAppData.Pfds != nil {
			pfdAppData.Pfds = make(map[string]Pfd)
		}

		for _, inputPfdData := range inputPfdAppData.Pfds {
			pfdAppData.Pfds[inputPfdData.PfdID] = Pfd(inputPfdData)
		}
		pfdTransData.PfdDatas[pfdAppData.ExternalAppID] = pfdAppData
	}
	return pfdTransData
}

func printPdfTransStatus(pfdTransData []byte, transURI string) {

	//Convert the json PFD Transaction data into struct
	pfdTrans := PfdManagement{}
	err := json.Unmarshal(pfdTransData, &pfdTrans)
	if err != nil {
		klog.Info(err)
		return
	}
	if transURI == "" {
		transURI = string(pfdTrans.Self)
	}

	fmt.Printf("PFD Transaction URI: %s\n", transURI)
	fmt.Printf("PFD Transaction ID: %s\n",
		getTransIDFromURL(transURI))
	fmt.Println("    Application IDs:")

	appStatus := make(map[string]string)
	for k := range pfdTrans.PfdDatas {
		appStatus[k] = "Created"
	}
	for _, v := range pfdTrans.PfdReports {
		for _, str := range v.ExternalAppIds {
			appStatus[str] = string(v.FailureCode)
		}
	}
	for k, v := range appStatus {
		if v != "Created" {
			fmt.Printf("      - %s : Failed (Reason: %s)\n", k, v)
		} else {
			fmt.Printf("      - %s : %s\n", k, v)
		}
	}
}

func getappSessionIDFromURL(appLoc string) string {
	appSessionLoc := strings.Split(appLoc, "/")
	for index, str := range appSessionLoc {
		if str == "app-sessions" {
			return appSessionLoc[index+1]
		}
	}
	return ""
}

func getPaAscReqData(inputPaAscReqData AFAscReqData) AppSessionContextReqData {
	var paAscReqData AppSessionContextReqData

	paAscReqData.AfAppID = inputPaAscReqData.Policy.AfAppID
	paAscReqData.AspID = inputPaAscReqData.Policy.AspID
	paAscReqData.BdtRefID = inputPaAscReqData.Policy.BdtRefID
	paAscReqData.Dnn = inputPaAscReqData.Policy.Dnn
	paAscReqData.IPDomain = inputPaAscReqData.Policy.IPDomain
	paAscReqData.MpsID = inputPaAscReqData.Policy.MpsID
	paAscReqData.NotifURI = inputPaAscReqData.Policy.NotifURI
	paAscReqData.SponID = inputPaAscReqData.Policy.SponID
	paAscReqData.Supi = inputPaAscReqData.Policy.Supi
	paAscReqData.Gpsi = inputPaAscReqData.Policy.Gpsi
	paAscReqData.SuppFeat = inputPaAscReqData.Policy.SuppFeat
	paAscReqData.UeIpv4 = IPv4Addr(inputPaAscReqData.Policy.UeIpv4)
	paAscReqData.UeIpv6 = IPv6Addr(inputPaAscReqData.Policy.UeIpv6)
	paAscReqData.UeMac = MacAddr(inputPaAscReqData.Policy.UeMac)
	paAscReqData.SponStatus = SponsoringStatus(inputPaAscReqData.Policy.SponStatus)

	//SliceInfo
	if inputPaAscReqData.Policy.SliceInfo != nil {
		sliceInfo := SNSSAI(*inputPaAscReqData.Policy.SliceInfo)
		paAscReqData.SliceInfo = &sliceInfo
	}

	//AfRoutReq
	if inputPaAscReqData.Policy.AfRoutReq != nil {
		inputAfRoutReq := inputPaAscReqData.Policy.AfRoutReq

		var afRoutReq RoutingRequirement

		afRoutReq.AppReloc = inputAfRoutReq.AppReloc

		//RouteToLocs
		afRoutReq.RouteToLocs = make([]RouteToLocation, len(inputAfRoutReq.RouteToLocs))
		for i, inputRouteToLocs := range inputAfRoutReq.RouteToLocs {
			var routeToLocs RouteToLocation

			routeToLocs.DNAI = DNAI(inputRouteToLocs.DNAI)
			routeToLocs.RouteProfID = inputRouteToLocs.RouteProfID

			//RouteInfo
			inputrouteInfo := inputRouteToLocs.RouteInfo
			var routeInfo RouteInformation
			routeInfo.IPv4Addr = IPv4Addr(inputrouteInfo.IPv4Addr)
			routeInfo.IPv6Addr = IPv6Addr(inputrouteInfo.IPv6Addr)
			routeInfo.PortNumber = inputrouteInfo.PortNumber

			routeToLocs.RouteInfo = routeInfo
			afRoutReq.RouteToLocs[i] = routeToLocs

		}

		//SpVal
		if inputAfRoutReq.SpVal != nil {
			inputSpVal := inputAfRoutReq.SpVal
			var spVal SpatialValidity

			spVal.PresenceInfoList = make(map[string]PresenceInfo)

			for _, inputPresenceInfoList := range inputSpVal.PresenceInfoList {
				var presenceInfoList PresenceInfo

				presenceInfoList.PraID = inputPresenceInfoList.PraID
				presenceInfoList.PresenceState = PresenceState(inputPresenceInfoList.PresenceState)

				//TrackingAreaList
				presenceInfoList.TrackingAreaList = make([]Tai, len(inputPresenceInfoList.TrackingAreaList))
				for i, inputTai := range inputPresenceInfoList.TrackingAreaList {
					var tai Tai

					tai.PlmnID = PlmnID(inputTai.PlmnID)
					tai.Tac = inputTai.Tac

					presenceInfoList.TrackingAreaList[i] = tai
				}

				//EcgiList
				presenceInfoList.EcgiList = make([]Ecgi, len(inputPresenceInfoList.EcgiList))
				for i, inputEcgi := range inputPresenceInfoList.EcgiList {
					var ecgi Ecgi

					ecgi.PlmnID = PlmnID(inputEcgi.PlmnID)
					ecgi.EutraCellID = inputEcgi.EutraCellID

					presenceInfoList.EcgiList[i] = ecgi
				}

				//NcgiList
				presenceInfoList.NcgiList = make([]Ncgi, len(inputPresenceInfoList.NcgiList))
				for i, inputNcgi := range inputPresenceInfoList.NcgiList {
					var ncgi Ncgi

					ncgi.PlmnID = PlmnID(inputNcgi.PlmnID)
					ncgi.NrCellID = inputNcgi.NrCellID

					presenceInfoList.NcgiList[i] = ncgi
				}

				//GlobalRanNodeIDList
				presenceInfoList.GlobalRanNodeIDList = make([]GlobalRanNodeID,
					len(inputPresenceInfoList.GlobalRanNodeIDList))
				for i, inputGlobalRanNodeID := range inputPresenceInfoList.GlobalRanNodeIDList {
					var globalRanNodeID GlobalRanNodeID

					globalRanNodeID.N3IwfID = inputGlobalRanNodeID.N3IwfID
					globalRanNodeID.NgeNbID = inputGlobalRanNodeID.NgeNbID

					//PlmnID
					plmnID := PlmnID(*inputGlobalRanNodeID.PlmnID)
					globalRanNodeID.PlmnID = &plmnID

					//GnbID
					gnbID := GnbID(*inputGlobalRanNodeID.GnbID)
					globalRanNodeID.GnbID = &gnbID

					presenceInfoList.GlobalRanNodeIDList[i] = globalRanNodeID
				}

				spVal.PresenceInfoList[string(presenceInfoList.PresenceState)] = presenceInfoList
			}
			afRoutReq.SpVal = &spVal
		}

		//TempVals
		afRoutReq.TempVals = make([]TemporalValidity, len(inputAfRoutReq.TempVals))
		for i, inputTempVals := range inputAfRoutReq.TempVals {
			var tempVals TemporalValidity

			tempVals.StartTime = inputTempVals.StartTime
			tempVals.StopTime = inputTempVals.StopTime

			afRoutReq.TempVals[i] = tempVals
		}

		//UpPathChgSub
		if inputAfRoutReq.UpPathChgSub != nil {
			inputUpPathChgSub := inputAfRoutReq.UpPathChgSub
			var upPathChgSub UpPathChgEvent

			upPathChgSub.NotificationURI = inputUpPathChgSub.NotificationURI
			upPathChgSub.NotifCorreID = inputUpPathChgSub.NotifCorreID
			upPathChgSub.DnaiChgType = DNAIChangeType(inputUpPathChgSub.DnaiChgType)

			afRoutReq.UpPathChgSub = &upPathChgSub
		}
		paAscReqData.AfRoutReq = &afRoutReq
	}

	//EvSubsc
	if inputPaAscReqData.Policy.EvSubsc != nil {
		inputEvSubsc := inputPaAscReqData.Policy.EvSubsc

		var evSubsc EventsSubscReqData

		//Events
		evSubsc.Events = make([]EventSubscription, len(inputEvSubsc.Events))
		for i, inputEvents := range inputEvSubsc.Events {
			var events EventSubscription

			events.Event = Event(inputEvents.Event)
			events.NotifMethod = NotifMethod(inputEvents.NotifMethod)

			evSubsc.Events[i] = events
		}

		//NotifURI
		evSubsc.NotifURI = inputEvSubsc.NotifURI

		//UsgThres
		if inputEvSubsc.UsgThres != nil {
			usgThres := UsageThreshold(*inputEvSubsc.UsgThres)
			evSubsc.UsgThres = &usgThres
		}
		paAscReqData.EvSubsc = &evSubsc
	}

	//MedComponents
	if inputPaAscReqData.Policy.MedComponents != nil {
		paAscReqData.MedComponents = make(map[string]MediaComponent)
	}

	for _, inputMedComponent := range inputPaAscReqData.Policy.MedComponents {
		var medComponent MediaComponent

		medComponent.ContVer = inputMedComponent.ContVer
		//MedCompN
		if inputMedComponent.MedCompN != 0 {
			medComponent.MedCompN = inputMedComponent.MedCompN
		}
		medComponent.AfAppID = inputMedComponent.AfAppID
		medComponent.MarBwDl = inputMedComponent.MarBwDl
		medComponent.MarBwUl = inputMedComponent.MarBwUl
		medComponent.MirBwDl = inputMedComponent.MirBwDl
		medComponent.MirBwUl = inputMedComponent.MirBwUl
		medComponent.Codecs = inputMedComponent.Codecs

		//AfRoutReq
		if inputMedComponent.AfRoutReq != nil {
			inputAfRoutReq := inputMedComponent.AfRoutReq

			var afRoutReq RoutingRequirement

			afRoutReq.AppReloc = inputAfRoutReq.AppReloc

			//RouteToLocs
			afRoutReq.RouteToLocs = make([]RouteToLocation, len(inputAfRoutReq.RouteToLocs))
			for i, inputRouteToLocs := range inputAfRoutReq.RouteToLocs {
				var routeToLocs RouteToLocation

				routeToLocs.DNAI = DNAI(inputRouteToLocs.DNAI)
				routeToLocs.RouteProfID = inputRouteToLocs.RouteProfID

				//RouteInfo
				inputrouteInfo := inputRouteToLocs.RouteInfo
				var routeInfo RouteInformation
				routeInfo.IPv4Addr = IPv4Addr(inputrouteInfo.IPv4Addr)
				routeInfo.IPv6Addr = IPv6Addr(inputrouteInfo.IPv6Addr)
				routeInfo.PortNumber = inputrouteInfo.PortNumber

				routeToLocs.RouteInfo = routeInfo
				afRoutReq.RouteToLocs[i] = routeToLocs

			}

			//SpVal
			if inputAfRoutReq.SpVal != nil {
				inputSpVal := inputAfRoutReq.SpVal
				var spVal SpatialValidity

				spVal.PresenceInfoList = make(map[string]PresenceInfo)

				for _, inputPresenceInfoList := range inputSpVal.PresenceInfoList {
					var presenceInfoList PresenceInfo

					presenceInfoList.PraID = inputPresenceInfoList.PraID
					presenceInfoList.PresenceState = PresenceState(inputPresenceInfoList.PresenceState)

					//TrackingAreaList
					presenceInfoList.TrackingAreaList = make([]Tai, len(inputPresenceInfoList.TrackingAreaList))
					for i, inputTai := range inputPresenceInfoList.TrackingAreaList {
						var tai Tai

						tai.PlmnID = PlmnID(inputTai.PlmnID)
						tai.Tac = inputTai.Tac

						presenceInfoList.TrackingAreaList[i] = tai
					}

					//EcgiList
					presenceInfoList.EcgiList = make([]Ecgi, len(inputPresenceInfoList.EcgiList))
					for i, inputEcgi := range inputPresenceInfoList.EcgiList {
						var ecgi Ecgi

						ecgi.PlmnID = PlmnID(inputEcgi.PlmnID)
						ecgi.EutraCellID = inputEcgi.EutraCellID

						presenceInfoList.EcgiList[i] = ecgi
					}

					//NcgiList
					presenceInfoList.NcgiList = make([]Ncgi, len(inputPresenceInfoList.NcgiList))
					for i, inputNcgi := range inputPresenceInfoList.NcgiList {
						var ncgi Ncgi

						ncgi.PlmnID = PlmnID(inputNcgi.PlmnID)
						ncgi.NrCellID = inputNcgi.NrCellID

						presenceInfoList.NcgiList[i] = ncgi
					}

					//GlobalRanNodeIDList
					presenceInfoList.GlobalRanNodeIDList = make([]GlobalRanNodeID,
						len(inputPresenceInfoList.GlobalRanNodeIDList))
					for i, inputGlobalRanNodeID := range inputPresenceInfoList.GlobalRanNodeIDList {
						var globalRanNodeID GlobalRanNodeID

						globalRanNodeID.N3IwfID = inputGlobalRanNodeID.N3IwfID
						globalRanNodeID.NgeNbID = inputGlobalRanNodeID.NgeNbID

						//PlmnID
						plmnID := PlmnID(*inputGlobalRanNodeID.PlmnID)
						globalRanNodeID.PlmnID = &plmnID

						//GnbID
						gnbID := GnbID(*inputGlobalRanNodeID.GnbID)
						globalRanNodeID.GnbID = &gnbID

						presenceInfoList.GlobalRanNodeIDList[i] = globalRanNodeID
					}

					spVal.PresenceInfoList[string(presenceInfoList.PresenceState)] = presenceInfoList
				}
				afRoutReq.SpVal = &spVal
			}

			//TempVals
			afRoutReq.TempVals = make([]TemporalValidity, len(inputAfRoutReq.TempVals))
			for i, inputTempVals := range inputAfRoutReq.TempVals {
				var tempVals TemporalValidity

				tempVals.StartTime = inputTempVals.StartTime
				tempVals.StopTime = inputTempVals.StopTime

				afRoutReq.TempVals[i] = tempVals
			}

			//UpPathChgSub
			if inputAfRoutReq.UpPathChgSub != nil {
				inputUpPathChgSub := inputAfRoutReq.UpPathChgSub
				var upPathChgSub UpPathChgEvent

				upPathChgSub.NotificationURI = inputUpPathChgSub.NotificationURI
				upPathChgSub.NotifCorreID = inputUpPathChgSub.NotifCorreID
				upPathChgSub.DnaiChgType = DNAIChangeType(inputUpPathChgSub.DnaiChgType)

				afRoutReq.UpPathChgSub = &upPathChgSub
			}
			medComponent.AfRoutReq = &afRoutReq
		}

		//FStatus
		medComponent.FStatus = FlowStatus(inputMedComponent.FStatus)

		//ResPrio
		medComponent.ResPrio = ReservPriority(inputMedComponent.ResPrio)

		//MedType
		medComponent.MedType = MediaType(inputMedComponent.MedType)

		//MedSubComps
		if inputMedComponent.MedSubComps != nil {
			medComponent.MedSubComps = make(map[string]MediaSubComponent)
		}

		for _, inputMedSubComponent := range inputMedComponent.MedSubComps {
			var medSubComponent MediaSubComponent

			//FNum
			if inputMedSubComponent.FNum != 0 {
				medSubComponent.FNum = inputMedSubComponent.FNum
			}
			medSubComponent.FDescs = inputMedSubComponent.FDescs
			medSubComponent.FStatus = FlowStatus(inputMedSubComponent.FStatus)
			medSubComponent.MarBwDl = inputMedSubComponent.MarBwDl
			medSubComponent.MarBwUl = inputMedSubComponent.MarBwUl
			medSubComponent.TosTrCl = inputMedSubComponent.TosTrCl
			medSubComponent.FlowUsage = FlowUsage(inputMedSubComponent.FlowUsage)

			//EthfDescs
			medSubComponent.EthfDescs = make([]EthFlowDescription, len(inputMedSubComponent.EthfDescs))
			for i, inputEthfDescs := range inputMedSubComponent.EthfDescs {
				var ethfDescs EthFlowDescription

				ethfDescs.DestMacAddr = inputEthfDescs.DestMacAddr
				ethfDescs.EthType = inputEthfDescs.EthType
				ethfDescs.FDesc = inputEthfDescs.FDesc
				ethfDescs.FDir = inputEthfDescs.FDir
				ethfDescs.SourceMacAddr = inputEthfDescs.SourceMacAddr
				ethfDescs.VLANTags = inputEthfDescs.VLANTags

				medSubComponent.EthfDescs[i] = ethfDescs
			}

			medComponent.MedSubComps[strconv.Itoa(int(medSubComponent.FNum))] = medSubComponent
		}
		paAscReqData.MedComponents[strconv.Itoa(int(medComponent.MedCompN))] = medComponent
	}
	return paAscReqData
}

func printAscData(paAscData []byte) {
	paAppSession := AppSessionContext{}
	if err := json.Unmarshal(paAscData, &paAppSession); err != nil {
		fmt.Println(err)
		return
	}

	asc, err := yaml.Marshal(paAppSession)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(asc))

}
