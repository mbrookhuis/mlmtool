// Package susemanager - SUSE Manager api call and support functions
package susemanager

import (
	"encoding/json"
	"fmt"

	sumamodels "mlmtool/pkg/models/susemanager"
	log "mlmtool/pkg/util/logger"
	returnCodes "mlmtool/pkg/util/returnCodes"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ChannelListSoftwareChannels - list all software channels
//
// param: auth
// return:
func (p *Proxy) ChannelListSoftwareChannels(auth AuthParams) ([]sumamodels.ChannelListSoftwareChannels, error) {
	log.Debug("ChannelListSoftwareChannels function call started")
	path := "channel/listSoftwareChannels"
	response, err := p.suse.SuseManagerCall(nil, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return nil, fmt.Errorf("error while calling list software channels manager err: %s", err)
	}
	var resultSuc []sumamodels.ChannelListSoftwareChannels
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("unmarshling error", zap.Any("error", err))
			return nil, fmt.Errorf("error while calling list software channels manager err: %s", err)
		}
	} else {
		log.Error("list software channels call failed", zap.Any("status code", response.StatusCode))
		return nil, fmt.Errorf("calling list software channels manager Failed. Http StatusCode: %s Http Body: %s", fmt.Sprint(response.StatusCode), fmt.Sprint(response.Body))
	}
	log.Debug("Completed ChannelListSoftwareChannels function")
	return resultSuc, nil
}

// ChannelSoftwareListChildren - list software channels from given parent
//
// param: auth
// param: label
// return:
func (p *Proxy) ChannelSoftwareListChildren(auth AuthParams, label string) ([]sumamodels.ChannelSoftwareListChildren, error) {
	log.Debug("Inside ChannelSoftwareListChildren function")
	body, _ := json.Marshal(map[string]interface{}{"channelLabel": label})
	path := "channel/software/listChildren"
	response, err := p.suse.SuseManagerCall(body, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return nil, fmt.Errorf("error while calling list software channels err: %s", err)
	}
	var resultSuc []sumamodels.ChannelSoftwareListChildren
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error("unmarshling error", zap.Any("error", err))
			return nil, fmt.Errorf("error while calling list child software channels , err: %s", err)
		}
	} else {
		return nil, fmt.Errorf("calling list software channels Failed. Http StatusCode: %s Http Body: %s", fmt.Sprint(response.StatusCode), fmt.Sprint(response.Body))
	}
	log.Debug("Completed ChannelSoftwareListChildren function")
	return resultSuc, nil
}

func (p *Proxy) ChannelSoftwareCreateRepo(auth AuthParams, label string, typeRepo string, url string) (sumamodels.ChannelSoftwareCreateRepo, error) {
	log.Debug("Inside ChannelSoftwareCreateRepo function")
	body, err := json.Marshal(map[string]interface{}{"label": label, "type": typeRepo, "url": url})
	var resultSuc sumamodels.ChannelSoftwareCreateRepo
	if err != nil {
		log.Warn(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return resultSuc, errors.New(returnCodes.ErrFailedMarshalling)
	}
	path := "channel/software/createRepo"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Warn(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("error", err))
		return resultSuc, fmt.Errorf(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Warn(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("response", resp), zap.Any("error", err))
			return resultSuc, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Warn(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return resultSuc, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Warn(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return resultSuc, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	log.Debug("Completed ChannelSoftwareCreateRepo function")
	return resultSuc, nil
}

func (p *Proxy) ChannelSoftwareCreate(auth AuthParams, label string, name string, summary string, archLabel string, parentLabel string) (int, error) {
	log.Debug("Inside ChannelSoftwareCreate function")
	var resultSuc int
	body, err := json.Marshal(map[string]interface{}{"label": label, "summary": summary, "archLabel": archLabel, "parentLabel": parentLabel, "name": name})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return 0, errors.New(returnCodes.ErrFailedMarshalling)
	}
	path := "channel/software/create"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("error", err))
		return 0, fmt.Errorf(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("response", resp), zap.Any("error", err))
			return 0, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return 0, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return 0, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	log.Debug("Completed ChannelSoftwareCreate function")
	return resultSuc, nil
}

func (p *Proxy) ChannelSoftwareAssociateRepo(auth AuthParams, channelLabel string, repoLabel string) (sumamodels.ChannelSoftwareListChildren, error) {
	log.Debug("Inside ChannelSoftwareAssociateRepo function")
	var resultSuc sumamodels.ChannelSoftwareListChildren
	body, err := json.Marshal(map[string]interface{}{"channelLabel": channelLabel, "repoLabel": repoLabel})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return resultSuc, errors.New(returnCodes.ErrFailedMarshalling)
	}
	path := "channel/software/associateRepo"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("error", err))
		return resultSuc, fmt.Errorf(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("response", resp), zap.Any("error", err))
			return resultSuc, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return resultSuc, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return resultSuc, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	log.Debug("Completed ChannelSoftwareAssociateRepo function")
	return resultSuc, nil
}

func (p *Proxy) ChannelSoftwareSyncRepo(auth AuthParams, channelLabel string) (int, error) {
	log.Debug("Inside ChannelSoftwareSyncRepo function")
	var resultSuc int
	body, err := json.Marshal(map[string]interface{}{"channelLabel": channelLabel})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return resultSuc, errors.New(returnCodes.ErrFailedMarshalling)
	}
	path := "channel/software/syncRepo"
	response, err := p.suse.SuseManagerCall(body, "POST", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("error", err))
		return resultSuc, fmt.Errorf(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("response", resp), zap.Any("error", err))
			return resultSuc, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return resultSuc, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return resultSuc, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	log.Debug("Completed ChannelSoftwareSyncRepo function")
	return resultSuc, nil
}

func (p *Proxy) ChannelSoftwareIsExisting(auth AuthParams, label string) (bool, error) {
	log.Debug("Inside ChannelSoftwareIsExisting function", zap.Any("Label", label))
	var resultSuc bool
	body, err := json.Marshal(map[string]interface{}{"channelLabel": label})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return false, errors.New(returnCodes.ErrFailedMarshalling)
	}
	path := "channel/software/isExisting"
	response, err := p.suse.SuseManagerCall(body, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("error", err))
		return false, fmt.Errorf(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("response", resp), zap.Any("error", err))
			return false, errors.New(returnCodes.ErrFailedUnMarshalling)
		}

		byteArray, _ := json.Marshal(resp)
		// log.Debug(byteArray)
		err = json.Unmarshal(byteArray, &resultSuc)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return false, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return false, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	log.Debug("Completed ChannelSoftwareIsExisting function")
	return resultSuc, nil
}
