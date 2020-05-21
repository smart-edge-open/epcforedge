// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

// UserplaneFunction :
//  - NONE: No function
//  - SGWU: 4G serving gateway userplane (SGW-U)
//  - PGWU: 4G packet data network (PDN) gateway userplane (PGW-U)
//  - SAEGWU: 4G combination SGW-U and PGW-U
type UserplaneFunction string

// List of UserplaneFunction
const (
	// NONE: No function
	NONE UserplaneFunction = "NONE"
	// SGWU: 4G serving gateway userplane (SGW-U)
	SGWU UserplaneFunction = "SGWU"
	// PGWU: 4G packet data network (PDN) gateway userplane (PGW-U)
	PGWU UserplaneFunction = "PGWU"
	// SAEGWU: 4G combination SGW-U and PGW-U
	SAEGWU UserplaneFunction = "SAEGWU"
)

// CupsUserplaneID identifies CUPS userplane ID
type CupsUserplaneID struct {
	ID string `json:"id,omitempty"`
}

// CupsUserplane CupsUserplane
type CupsUserplane struct {
	ID        string              `json:"id,omitempty"`
	UUID      string              `json:"uuid,omitempty"`
	Function  *UserplaneFunction  `json:"function,omitempty"`
	Config    *UserplaneConfig    `json:"config,omitempty"`
	Selectors []UserplaneSelector `json:"selectors,omitempty"`
	// The UEs that should be entitled to access privileged networks via this
	// userplane.  Note: UEs not in this list will still be able to get a bearer
	// via the userplane. The UEs in this list are just for entitlement
	// purposes. (optional)
	Entitlements []UserplaneEntitlement `json:"entitlements,omitempty"`
}

// UserplaneConfig UserplaneConfig
type UserplaneConfig struct {
	Sxa      *ConfigInfoCpup `json:"sxa,omitempty"`
	Sxb      *ConfigInfoCpup `json:"sxb,omitempty"`
	S1u      *ConfigInfoUp   `json:"s1u,omitempty"`
	S5uSGW   *ConfigInfoUp   `json:"s5u_sgw,omitempty"`
	S5uPGW   *ConfigInfoUp   `json:"s5u_pgw,omitempty"`
	SGi      *ConfigInfoUp   `json:"sgi,omitempty"`
	Breakout []ConfigInfoUp  `json:"breakout,omitempty"`
	DNS      []ConfigInfoUp  `json:"dns,omitempty"`
}

// ConfigInfoCpup Information that the userplane should configure, which relates
// to the control plane (CP) side and the user plane (UP) side.
type ConfigInfoCpup struct {
	CpIPAddress string `json:"cp_ip_address,omitempty"`
	UpIPAddress string `json:"up_ip_address,omitempty"`
}

// ConfigInfoUp Information that the userplane should configure, which relates
// to the user plane (UP) side only.
type ConfigInfoUp struct {
	UpIPAddress string `json:"up_ip_address,omitempty"`
}

// UserplaneSelector UserplaneSelector
type UserplaneSelector struct {
	ID      string           `json:"id,omitempty"`
	Network *SelectorNetwork `json:"network,omitempty"`
	ULI     *SelectorUli     `json:"uli,omitempty"`
	PDN     *SelectorPdn     `json:"pdn,omitempty"`
}

// UserplaneEntitlement UserplaneEntitlement
type UserplaneEntitlement struct {
	ID    string                `json:"id,omitempty"`
	APNs  []string              `json:"apns,omitempty"`
	IMSIs []EntitlementImsiList `json:"imsis,omitempty"`
}

// SelectorNetwork SelectorNetwork
type SelectorNetwork struct {
	MCC string `json:"mcc,omitempty"`
	MNC string `json:"mnc,omitempty"`
}

// SelectorUli SelectorUli
type SelectorUli struct {
	TAI  *Ulitai  `json:"tai,omitempty"`
	ECGI *Uliecgi `json:"ecgi,omitempty"`
}

// SelectorPdn SelectorPdn
type SelectorPdn struct {
	APNs []string `json:"apns,omitempty"`
}

// Ulitai Ulitai
type Ulitai struct {
	// Tracking area code (TAC), which is typically an unsigned integer
	// from 1 to 2^16, inclusive.
	TAC int64 `json:"tac,omitempty"`
}

// Uliecgi Uliecgi
type Uliecgi struct {
	// E-UTRAN cell identifier (ECI), which is typically an unsigned integer
	// from 1 to 2^32, inclusive.
	ECI int64 `json:"eci,omitempty"`
}

// EntitlementImsiList EntitlementImsiList
type EntitlementImsiList struct {
	Begin string `json:"begin,omitempty"`
	End   string `json:"end,omitempty"`
}
