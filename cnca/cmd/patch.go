// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use: "patch",
	Short: "Patch an active LTE CUPS userplane or NGC AF TI subscription " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

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

			// patch subscription
			err = AFPatchSubscription(args[0], sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])

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

			// patch userplane
			err = LtePatchUserplane(args[0], up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])
		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

// pfdPatchCmd represents the patch command
var pfdPatchCmd = &cobra.Command{
	Use: "patch",
	Short: "Patch an active NGC AF PFD Transaction or NGC AF PFD Application " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println(errors.New("Missing input(s)"))
			return
		}

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

		if args[0] == "transaction" && args[1] != "" {
			var pfdReportData []byte
			if len(args) > 2 {
				if args[2] == "application" && len(args) > 3 {
					// patch PFD Application

					var s AFPfdData
					if err = yaml.Unmarshal(data, &s); err != nil {
						fmt.Println(err)
						return
					}

					pfdAppData := getPfdAppData(s)

					var app []byte
					app, err = json.Marshal(pfdAppData)
					if err != nil {
						fmt.Println(err)
						return
					}

					pfdReportData, err = AFPatchPfdApplication(args[1], args[3], app)
					if err != nil {
						klog.Info(err)
						if err.Error() == "HTTP failure: 500" && pfdReportData != nil {
							printPfdReport(pfdReportData)
						}
						return
					}
					fmt.Printf("PFD Application %s patched\n", args[3])
					return
				}
			} else {
				// patch PFD Transaction

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
				pfdReportData, err = AFPatchPfdTransaction(args[1], trans)
				if err != nil {
					klog.Info(err)
					if err.Error() == "HTTP failure: 500" && pfdReportData != nil {
						printPfdReports(pfdReportData)
					}
					return
				}
				fmt.Printf("PFD Transaction %s patched\n", args[1])
				return
			}
		}
		fmt.Println(errors.New("Invalid input(s)"))
	},
}

//paPatchCmd represents patch command
var paPatchCmd = &cobra.Command{
	Use: "patch",
	Short: "Patch an active NGC AF PFD Transaction or NGC AF PFD Application " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input(s)"))
			return
		}

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

		var s AFAscUpdateData
		if err = yaml.Unmarshal(data, &s); err != nil {
			fmt.Println(err)
			return
		}

		paAscUpdateData := getPaAscUpdateData(s)

		appSession, err := json.Marshal(paAscUpdateData)
		if err != nil {
			fmt.Println(err)
			return
		}

		if args[0] != "" {
			//Patch app-session data
			appSessionContext, err := AFPatchPaAppSession(args[0], appSession)
			if err != nil {
				klog.Info(err)
				return
			}
			appSessionContext, err = y2j.JSONToYAML(appSessionContext)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("App-Session %s was updated : \n %s", args[0], string(appSessionContext))
			return
		}
		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help = `Patch an active LTE CUPS userplane or NGC AF TI subscription 
using YAML configuration file

Usage:
  cnca patch { <userplane-id> | <subscription-id> } -f <config.yml>

Example:
  cnca patch <userplane-id> -f <config.yml>
  cnca patch <subscription-id> -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	const pfdHelp = `Patch an active NGC AF PFD Transaction or NGC AF PFD 
Application using YAML Configuration file

Usage:
  cnca pfd patch { transaction <transaction-id> |
                   transaction <transaction-id> application <application-id>}
                   -f <config.yml>

Example:
  cnca pfd patch transaction <transaction-id> -f <config.yml>
  cnca pfd patch transaction <transaction-id> 
    application <application-id> -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`

	const paHelp = `Patch an active NGC AF Application Session
using YAML Configuration file

Usage:
  cnca policy-authorization patch <appSessionID> -f <config_patch.yml>

Example:
  cnca policy-authorization patch <appSessionID> -f <config_patch.yml>

Flags: 
  -h, --help       help  
  -f, --filename   YAML configuration file 
`

	// add `patch` command
	cncaCmd.AddCommand(patchCmd)
	patchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = patchCmd.MarkFlagRequired("filename")
	patchCmd.SetHelpTemplate(help)

	// add pfd `patch` command
	pfdCmd.AddCommand(pfdPatchCmd)
	pfdPatchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = pfdPatchCmd.MarkFlagRequired("filename")
	pfdPatchCmd.SetHelpTemplate(pfdHelp)

	//add policy-authorization `patch` command
	paCmd.AddCommand(paPatchCmd)
	paPatchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = paPatchCmd.MarkFlagRequired("filename")
	paPatchCmd.SetHelpTemplate(paHelp)
}

func getPfdAppData(inputPfdAppData AFPfdData) PfdData {

	var pfdAppData PfdData

	pfdAppData.ExternalAppID = inputPfdAppData.Policy.ExternalAppID
	pfdAppData.Self = Link(inputPfdAppData.Policy.Self)

	if inputPfdAppData.Policy.AllowedDelay != nil {
		allowedDelay := DurationSecRm(*inputPfdAppData.Policy.AllowedDelay)
		pfdAppData.AllowedDelay = &allowedDelay
	}
	if inputPfdAppData.Policy.CachingTime != nil {
		cachingTime := DurationSecRo(*inputPfdAppData.Policy.CachingTime)
		pfdAppData.CachingTime = &cachingTime
	}
	if inputPfdAppData.Policy.Pfds != nil {
		pfdAppData.Pfds = make(map[string]Pfd)
	}

	for _, inputPfdData := range inputPfdAppData.Policy.Pfds {
		pfdAppData.Pfds[inputPfdData.PfdID] = Pfd(inputPfdData)
	}

	return pfdAppData
}

func printPfdReports(pfdReportData []byte) {
	//Convert the json PFD Transaction Report into struct
	pfdReports := []PfdReport{}
	err := json.Unmarshal(pfdReportData, &pfdReports)
	if err != nil {
		klog.Info(err)
		return
	}
	fmt.Println("PFD Transaction ID:")
	fmt.Println("    Application IDs:")
	appStatus := make(map[string]string)
	for _, v := range pfdReports {
		for _, str := range v.ExternalAppIds {
			appStatus[str] = string(v.FailureCode)
		}
	}
	for k, v := range appStatus {
		fmt.Printf("      - %s : Failed (Reason: %s)\n", k, v)
	}
}

func printPfdReport(pfdReportData []byte) {
	//Convert the json PFD Transaction Report into struct
	pfdReport := PfdReport{}
	err := json.Unmarshal(pfdReportData, &pfdReport)
	if err != nil {
		klog.Info(err)
		return
	}
	fmt.Println("PFD Transaction ID:")
	fmt.Println("    Application IDs:")

	var appIDList []string
	if pfdReport.ExternalAppIds != nil {
		copy(appIDList, pfdReport.ExternalAppIds)
	}

	for _, v := range appIDList {
		fmt.Printf("      - %s : Failed (Reason: %s)\n", v, string(pfdReport.FailureCode))
	}
}

func getPaAscUpdateData(inputPaAscUpdateData AFAscUpdateData) AppSessionContextUpdateData {

	var paAscUpdateData AppSessionContextUpdateData

	paAscUpdateData.AfAppID = inputPaAscUpdateData.Policy.AfAppID
	paAscUpdateData.AspID = inputPaAscUpdateData.Policy.AspID
	paAscUpdateData.BdtRefID = inputPaAscUpdateData.Policy.BdtRefID
	paAscUpdateData.MpsID = inputPaAscUpdateData.Policy.MpsID
	paAscUpdateData.SponID = inputPaAscUpdateData.Policy.SponID
	paAscUpdateData.SponStatus = SponsoringStatus(inputPaAscUpdateData.Policy.SponStatus)

	//AfRoutReq
	if inputPaAscUpdateData.Policy.AfRoutReq != nil {
		inputAfRoutReq := inputPaAscUpdateData.Policy.AfRoutReq

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
		paAscUpdateData.AfRoutReq = &afRoutReq
	}

	//EvSubsc
	if inputPaAscUpdateData.Policy.EvSubsc != nil {
		inputEvSubsc := inputPaAscUpdateData.Policy.EvSubsc

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
		paAscUpdateData.EvSubsc = &evSubsc
	}

	//MedComponents
	if inputPaAscUpdateData.Policy.MedComponents != nil {
		paAscUpdateData.MedComponents = make(map[string]MediaComponent)
	}

	for _, inputMedComponent := range inputPaAscUpdateData.Policy.MedComponents {
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
		paAscUpdateData.MedComponents[strconv.Itoa(int(medComponent.MedCompN))] = medComponent
	}
	return paAscUpdateData
}
