// Package susemanager - SUSE Manager api call and support functions
package susemanager

import (
	"encoding/json"
	"fmt"
	log "mlmtool/pkg/util/logger"

	"go.uber.org/zap"

	sumamodels "mlmtool/pkg/models/susemanager"
)

// SystemGroupCreate - create the systemgroup
//
// param: auth
// param: groupName
// param: description
// return:
func (p *Proxy) SystemGroupCreate(auth AuthParams, groupName string, description string) (*sumamodels.SystemGroupGetDetails, error) {
	body, _ := json.Marshal(map[string]any{"name": groupName, "description": description})
	path := "systemgroup/create"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		return nil, err
	}
	var resultSuc sumamodels.SystemGroupGetDetails
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error("error while handling suse manager response for creating system group", zap.Any("error", err))
			return nil, fmt.Errorf("error while handling suse manager response for creating system group, err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("unmarshalling error api response", zap.Any("error", err))
			return nil, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		return nil, fmt.Errorf("error while getting suse manager response for creating system group, Http StatusCode: %s Http Body: %s", fmt.Sprint(response.StatusCode), fmt.Sprint(response.Body))
	}
	return &resultSuc, nil
}

// SystemGroupGetDetails - get the details from the given systemgroup
//
// param: auth
// param: groupName
// return:
func (p *Proxy) SystemGroupGetDetails(auth AuthParams, groupName string) (*sumamodels.SystemGroupGetDetails, error) {
	body, _ := json.Marshal(map[string]any{"systemGroupName": groupName})
	path := "systemgroup/getDetails"
	response, err := p.suse.SuseManagerCall(body, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		return nil, err
	}
	var resultSuc sumamodels.SystemGroupGetDetails
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error("error while handling suse manager response for fetching system group details", zap.Any("error", err))
			return nil, fmt.Errorf("error while handling suse manager response for fetching system group details, err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("unmarshalling error suse-m response", zap.Any("error", err))
			return nil, fmt.Errorf("unmarshalling error: %s", err)
		}
	} else {
		return nil, fmt.Errorf("error while getting suse manager response for fetching system group details, Http StatusCode: %s Http Body: %s", fmt.Sprint(response.StatusCode), fmt.Sprint(response.Body))
	}
	return &resultSuc, nil
}

// SystemGroupListSystemsMinimal - list assigned systems from the given system group
//
// param: auth
// param: groupName
// return:
func (p *Proxy) SystemGroupListSystemsMinimal(auth AuthParams, groupName string) ([]sumamodels.SystemGroupListSystemsMinimal, error) {
	body, _ := json.Marshal(map[string]interface{}{"systemGroupName": groupName})
	path := "systemgroup/listSystemsMinimal"
	response, err := p.suse.SuseManagerCall(body, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("error while fetching system group list systemsMinimal", zap.Any("error", err))
		return nil, fmt.Errorf("error while fetching system group list systemsMinimal err: %s", err)
	}
	var resultSuc []sumamodels.SystemGroupListSystemsMinimal
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("unmarshling error", zap.Any("error", err))
			return nil, fmt.Errorf("unable to process the received data. err: %s", err)
		}
	} else {
		return nil, fmt.Errorf("calling system group list systemsMinimal Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "SystemGroupListSystemsMinimal"), zap.Any("response", resultSuc))
	return resultSuc, nil
}

func (p *Proxy) SystemGroupListActiveSystemsInGroup(auth AuthParams, groupName string) ([]int, error) {
	body, _ := json.Marshal(map[string]interface{}{"systemGroupName": groupName})
	path := "systemgroup/listActiveSystemsInGroup"
	response, err := p.suse.SuseManagerCall(body, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("error while fetching system group list of active systems", zap.Any("error", err))
		return nil, fmt.Errorf("error while fetching system group list active systems err: %s", err)
	}
	var resultSuc []int
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("unmarshling error", zap.Any("error", err))
			return nil, fmt.Errorf("unable to process the received data. err: %s", err)
		}
	} else {
		return nil, fmt.Errorf("calling system group list systemsMinimal Failed. Http StatusCode: %s", fmt.Sprint(response.StatusCode))
	}
	log.Debug("Response from api", zap.Any("api", "SystemGroupListSystemsMinimal"), zap.Any("response", resultSuc))
	return resultSuc, nil
}
