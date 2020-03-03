/* SPDX-License-Identifier: Apache-2.0
* Copyright (c) 2020 Intel Corporation
 */

package ngcnef

// PfdContent represents the content of a PFD for an application identifier.
type PfdContent struct {
	// Identifies a PDF of an application identifier.
	PfdID string `json:"pfdID,omitempty"`
	// Represents a 3-tuple with protocol, server ip and server port for
	// UL/DL application traffic.
	FlowDescriptions []string `json:"flowDescriptions,omitempty"`
	// Indicates a URL or a regular expression which is used to match the
	// significant parts of the URL.
	Urls []string `json:"urls,omitempty"`
	// Indicates an FQDN or a regular expression as a domain name matching
	// criteria.
	DomainNames []string `json:"domainNames,omitempty"`
}

// PfdDataForApp represents the PFDs for an application identifier
type PfdDataForApp struct {
	// Identifier of an application.
	AppID ApplicationID `json:"appID"`
	// PFDs for the application identifier.
	Pfds []PfdContent `json:"pfds"`
	// Caching time for an application identifier.
	CachingTime *DateTime `json:"cachingTime,omitempty"`
}
