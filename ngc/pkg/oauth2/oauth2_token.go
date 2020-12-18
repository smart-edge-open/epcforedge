// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package oauth2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go/v4"

	logger "github.com/open-ness/common/log"
)

var log = logger.DefaultLogger.WithField("oauth2", nil)

// Path for OAuth2 Configuration file
const cfgPath string = "configs/oauth2.json"

//TokenVerificationResult Result of the token verification
type TokenVerificationResult string

// Error results of token verification
const (
	StatusBadRequest   = "StatusBadRequest"
	StatusInvalidToken = "StatusInvalidToken"
	StatusSuccess      = "StatusSuccess"
	StatusConfigErr    = "StatusConfigErr"
)

//Config OAuth2 config struct
type Config struct {
	SigningKey string `json:"signingkey"`
	Expiration int64  `json:"expiration"`
}

//PlmnID PLMN ID struct
type PlmnID struct {
	Mcc string `json:"mcc" yaml:"mcc" bson:"mcc" mapstructure:"Mcc"`
	Mnc string `json:"mnc" yaml:"mnc" bson:"mnc" mapstructure:"Mnc"`
}

//NfType Network Function type
type NfType string

//AccessTokenReq NRF access token request
type AccessTokenReq struct {
	GrantType          string  `json:"grant_type"`
	NfInstanceID       string  `json:"nfInstanceId"`
	NfType             NfType  `json:"nfType,omitempty"`
	TargetNfType       NfType  `json:"targetNfType,omitempty"`
	Scope              string  `json:"scope"`
	TargetNfInstanceID string  `json:"targetNfInstanceId,omitempty"`
	RequesterPlmn      *PlmnID `json:"requesterPlmn,omitempty"`
	TargetPlmn         *PlmnID `json:"targetPlmn,omitempty"`
}

//AccessTokenClaims struct
type AccessTokenClaims struct {
	Issuer     string      `json:"issuer"`
	Subject    string      `json:"subject"`
	Audience   interface{} `json:"audience"`
	Scope      string      `json:"scope"`
	Expiration int64       `json:"expiration"`
	jwt.StandardClaims
}

// LoadJSONConfig reads a file located at configPath and unmarshals it to
// config structure
func loadJSONConfig(configPath string, config interface{}) error {
	cfgData, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		log.Infoln(err)
		return err
	}
	return json.Unmarshal(cfgData, config)
}

// GetNEFAccessTokenFromNRF Generates the token. This is the functionality of
// NRF component of 5GC
func GetNEFAccessTokenFromNRF(accessTokenReq AccessTokenReq) (
	NefAccessToken string, err error) {

	var oAuth2Cfg = Config{}

	//Read Json config
	err = loadJSONConfig(cfgPath, &oAuth2Cfg)
	if err != nil {
		log.Errln("Failed to load OAuth2 configuration")
		return NefAccessToken, err
	}

	expiration := time.Now().Add(
		time.Second * time.Duration(oAuth2Cfg.Expiration)).Unix()

	jwtexp,e:= jwt.ParseTime(expiration)
	if e!=nil{
		log.Errln("Failed to parse time",e.Error())
	}

	log.Infoln("Token expires in : ", oAuth2Cfg.Expiration, " seconds")

	var mySigningKey = []byte(oAuth2Cfg.SigningKey)

	//log.Infoln("Expiration Set to ", expiration)
	// Create AccessToken
	var accessTokenClaims = AccessTokenClaims{
		"OpenNESS",             //Issuer:
		"NEF Validation token", //Subject:
		"AF-NEF",               //Audience:
		accessTokenReq.Scope,   //Scope:
		oAuth2Cfg.Expiration,   //Expiration:
		jwt.StandardClaims{ExpiresAt: jwtexp},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, SignedStringErr := token.SignedString(mySigningKey)

	if SignedStringErr != nil {
		log.Info(SignedStringErr)
	}

	return accessToken, nil
}

func fetchNEFAccessTokenFromNRF() (token string, err error) {

	var accessTokenReq AccessTokenReq

	//In case we construct and send it to NRF
	accessTokenReq.GrantType = "client_credentials"

	//Dont have it right now added static UUID
	accessTokenReq.NfInstanceID = "0"
	accessTokenReq.NfType = "AF"
	accessTokenReq.TargetNfType = "NEF"
	accessTokenReq.Scope = "nnrf-nfm"
	accessTokenReq.TargetNfInstanceID = "0" //Instance of NEF

	//POST AccessTokenRequest to NRF /oauth2/token
	return GetNEFAccessTokenFromNRF(accessTokenReq)
}

//GetAccessToken Get the access token to access NEF NF component. This API can
//			     be ported to operator provided token access mechanism
// i/p			: None
// o/p  token	: The access token
//		err		: error code in case of failure or nil in success
func GetAccessToken() (token string, err error) {

	token, err = fetchNEFAccessTokenFromNRF()

	if err != nil {
		log.Info("Failed to get NEF access token ")
	}
	return token, err
}

//ValidateAccessToken Validate the access token
// i/p reqToken : token to be validated
// o/p status : Success/Failure result of the operation
//     err    : error info of the token validation process.
func ValidateAccessToken(reqToken string) (status TokenVerificationResult,
	err error) {

	var oAuth2Cfg = Config{}

	//Read Json config
	err = loadJSONConfig(cfgPath, &oAuth2Cfg)
	if err != nil {
		log.Errln("Failed to load OAuth2 configuration")
		return StatusConfigErr, err
	}
	var mySigningKey = []byte(oAuth2Cfg.SigningKey)
	claims := &AccessTokenClaims{}

	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (
		interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		log.Info(err)

		if err == jwt.ErrSignatureInvalid {
			log.Info("Token is invalid, ErrSignatureInvalid")
			return StatusInvalidToken, err
		}
		//Check for Validation error
		validationErr, ok := err.(*jwt.MalformedTokenError)

		if !ok {
			log.Info(validationErr.Message)
			return StatusInvalidToken, err
		}

		return StatusBadRequest, err
	}
	if !tkn.Valid {
		log.Info("Token is invalid")
		return StatusInvalidToken, errors.New("Token is Invalid")
	}
	log.Info("OAuth2 Token Validation successful")
	return StatusSuccess, nil
}
