// Package susemanager - SUSE Manager api call and support functions
package susemanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "mlmtool/pkg/util/logger"

	"go.uber.org/zap"

	_sumamodels "mlmtool/pkg/models/susemanager"
)

// ActivationKeyListActivationKeys - list all current available activation keys
//
// param: auth
// return:
func (p *Proxy) ActivationKeyListActivationKeys(auth AuthParams) ([]_sumamodels.ActivationkeyGetDetails, error) {
	log.Debug("ActivationKeyListActivationKeys function call started")
	path := "activationkey/listActivationKeys"
	response, err := p.suse.SuseManagerCall(nil, http.MethodGet, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return nil, fmt.Errorf("error while calling list activationKeys manager err: %s", err)
	}
	var resultSuc []_sumamodels.ActivationkeyGetDetails
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("error while calling list software channels manager", zap.Any("error", err))
			return nil, fmt.Errorf("error while calling list software channels manager err: %s", err)
		}
	} else {
		log.Error("Error Bad Request", zap.Any("StatusCode", response.StatusCode))
		return nil, fmt.Errorf("calling list software channels manager Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "ActivationkeyListActivationKeys"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// ActivationKeyGetDetails - get details from activation key
//
// param: auth
// param: keyName
// return:
func (p *Proxy) ActivationKeyGetDetails(auth AuthParams, keyName string) (_sumamodels.ActivationkeyGetDetails, error) {
	log.Debug("ActivationkeyGetDetails function call started")
	var resultSuc _sumamodels.ActivationkeyGetDetails
	body, _ := json.Marshal(map[string]interface{}{"key": keyName})
	path := "activationkey/getDetails"
	response, err := p.suse.SuseManagerCall(body, http.MethodGet, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return resultSuc, fmt.Errorf("error while calling get details activationKeys manager err: %s", err)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return resultSuc, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("error while getting activation key details", zap.Any("error", err))
			return resultSuc, fmt.Errorf("error while getting activation key details manager err: %s", err)
		}
	} else {
		log.Error("fetching activation key details call Failed", zap.Any("StatusCode", response.StatusCode))
		return resultSuc, fmt.Errorf("fetching activation key details call Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "ActivationkeyGetDetails"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// ActivationKeyRemovePackages - remove packages to be installed on registration from activation key
//
// param: auth
// param: keyName
// param: pckgs
// return:
func (p *Proxy) ActivationKeyRemovePackages(auth AuthParams, keyName string, pckgs []_sumamodels.ActivationkeyPackages) (int, error) {
	log.Debug("ActivationKeyRemovePackages call started")
	body, _ := json.Marshal(map[string]interface{}{"key": keyName, "packages": pckgs})
	path := "activationkey/removePackages"
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return 0, fmt.Errorf("error while calling get details activationKeys manager err: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		log.Error("calling list software channels manager Failed", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("calling list software channels manager Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "ActivationkeyRemovePackages"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// ActivationKeyCreate - create activation key
//
// param: auth
// param: keyName
// param: baseChannel
// return:
func (p *Proxy) ActivationKeyCreate(auth AuthParams, keyName string, baseChannel string, entitlement []string) (string, error) {
	log.Debug("ActivationkeyCreate call started")
	// ent := make([]string, 0)
	body, _ := json.Marshal(map[string]interface{}{"baseChannelLabel": baseChannel, "key": keyName, "description": keyName,
		"entitlements": entitlement, "universalDefault": false})
	path := "activationkey/create"
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Creating activationkey failed", zap.Any("error", err))
		return "", fmt.Errorf("error while calling create activationKey. err: %s", err)
	}
	var resultSuc string
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return "", fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			return "", fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		log.Error("Creating activationkey failed", zap.Any("StatusCode", response.StatusCode))
		return "", fmt.Errorf("creating activationkey failed - http statuscode: %s - http response: %s ", fmt.Sprint(response.StatusCode), fmt.Sprint(resultSuc))
	}
	log.Debug("Response from api", zap.Any("api", "ActivationkeyCreate"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// ActivationKeyAddChildChannels - add softwarechannels to activation key
//
// param: auth
// param: keyName
// param: childChannels
// return:
func (p *Proxy) ActivationKeyAddChildChannels(auth AuthParams, keyName string, childChannels []string) (int, error) {
	log.Debug("ActivationkeyKeyAddChildChannels call started")
	body, _ := json.Marshal(map[string]interface{}{"childChannelLabels": childChannels, "key": keyName})
	path := "activationkey/addChildChannels"
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return 0, fmt.Errorf("error while calling get details activationKeys manager err: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		log.Error("calling activationKey/childChannels Failed", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("calling activationKey/childChannels Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "ActivationkeyKeyAddChildChannels"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// ActivationKeyAddServerGroups - add groups to automatically join
//
// param: auth
// param: keyName
// param: groups
// return:
func (p *Proxy) ActivationKeyAddServerGroups(auth AuthParams, keyName string, groups []int) (int, error) {
	log.Debug("ActivationKeyAddServerGroups call started")
	body, _ := json.Marshal(map[string]interface{}{"serverGroupIds": groups, "key": keyName})
	path := "activationkey/addServerGroups"
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return 0, fmt.Errorf("error while calling get details activationKeys manager err: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		log.Error("calling ActivationKey AddServerGroups Failed", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("calling ActivationKey AddServerGroups Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "ActivationKeyAddServerGroups"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

// ActivationKeyDelete - delete activation key
//
// param: auth
// param: keyName
// return:
func (p *Proxy) ActivationKeyDelete(auth AuthParams, keyName string) (int, error) {
	log.Debug("ActivationKeyDelete call started")
	body, _ := json.Marshal(map[string]interface{}{"key": keyName})
	path := "activationkey/delete"
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return 0, fmt.Errorf("error while calling get details activationKeys manager err: %s", err)
	}
	var resultSuc int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return 0, fmt.Errorf("error in handling suse manager response. err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			return 0, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		log.Error("calling ActivationKey Delete Failed", zap.Any("StatusCode", response.StatusCode))
		return 0, fmt.Errorf("calling ActivationKey Delete Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "ActivationKeyDelete"), zap.Any("response", resultSuc))
	return resultSuc, nil
}
