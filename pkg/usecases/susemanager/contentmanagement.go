// Package susemanager - SUSE Manager api call and support functions
package susemanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	sumamodels "mlmtool/pkg/models/susemanager"
	log "mlmtool/pkg/util/logger"
	returnCodes "mlmtool/pkg/util/returnCodes"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ContentManagementListProjects - list content projects
//
// param: auth
// return:
func (p *Proxy) ContentManagementListProjects(auth AuthParams) ([]sumamodels.ContentManagementListProjects, error) {
	log.Debug("contentManagement.listProjects function called")
	path := "contentmanagement/listProjects"
	response, err := p.suse.SuseManagerCall(nil, http.MethodGet, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("error", err))
		return nil, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	var result []sumamodels.ContentManagementListProjects
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
			return nil, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return nil, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

func (p *Proxy) ContentManagementLookupProject(auth AuthParams, project string) (sumamodels.ContentManagementListProjects, error) {
	log.Debug("contentManagement.lookupProject function called")
	var result sumamodels.ContentManagementListProjects

	path := "contentmanagement/lookupProject"
	body, err := json.Marshal(map[string]any{"projectLabel": project})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, err
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("action", "suseManagerResponse"), zap.Any("error", err))
		return result, err
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrHandlingSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
			return result, err
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, err
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

// ContentManagementCreate - create
//
// param: auth
// param: projectLabel
// param: name
// param: description
// return:
func (p *Proxy) ContentManagementCreate(auth AuthParams, projectLabel string, name string, description string) (sumamodels.ContentManagementListProjects, error) {
	log.Debug("contentManagement.create")
	var result sumamodels.ContentManagementListProjects
	path := "contentmanagement/createProject"
	body, err := json.Marshal(map[string]any{"name": name, "description": description, "projectLabel": projectLabel})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrFailedMarshalling)
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrFailedMarshalling, zap.Any("action", "HandleSuseManagerResponse"), zap.Any("error", err))
			return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		fmt.Println(returnCodes.ErrHTTPSuseManagerResponse)
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

// ContentManagementAttachSource - list attached channels
//
// param: auth
// param: projectLabel
// param: sourceType
// param: sourceLabel
// return:
func (p *Proxy) ContentManagementAttachSource(auth AuthParams, projectLabel string, sourceType string, sourceLabel string) (sumamodels.ContentManagementSource, error) {
	log.Debug("contentManagement.source")
	var result sumamodels.ContentManagementSource
	path := "contentmanagement/attachSource"
	body, err := json.Marshal(map[string]any{"projectLabel": projectLabel, "sourceType": sourceType, "sourceLabel": sourceLabel})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrFailedMarshalling)
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

// ContentManagementDetachSource - detach a source from a project
//
// param: auth
// param: projectLabel
// param: sourceType
// param: sourceLabel
// return:
func (p *Proxy) ContentManagementDetachSource(auth AuthParams, projectLabel string, sourceType string, sourceLabel string) error {
	log.Debug("contentManagement.source")
	var result int
	path := "contentmanagement/detachSource"
	body, err := json.Marshal(map[string]any{"projectLabel": projectLabel, "sourceType": sourceType, "sourceLabel": sourceLabel})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return errors.New(returnCodes.ErrFailedMarshalling)
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
		return errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
			return errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return nil
}

// ContentManagementListFilters - list available filters for content management
//
// param: auth
// return:
func (p *Proxy) ContentManagementListFilters(auth AuthParams) ([]sumamodels.ContentManagementFilter, error) {
	log.Debug("contentManagement.listFilters function called")
	path := "contentmanagement/listFilters"
	response, err := p.suse.SuseManagerCall(nil, http.MethodGet, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
		return nil, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	var result []sumamodels.ContentManagementFilter
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

// ContentManagementCreateFilter - create filter for content management
//
// param: auth
// param: name
// param: rule
// param: entityType
// param: criteria
// return:
func (p *Proxy) ContentManagementCreateFilter(auth AuthParams, name string, rule string, entityType string, criteria sumamodels.FilterCriteria) (sumamodels.ContentManagementFilter, error) {
	log.Debug("contentManagement.createFilter function called")
	var result sumamodels.ContentManagementFilter
	path := "contentmanagement/createFilter"
	body, err := json.Marshal(map[string]any{"name": name, "rule": rule, "entityType": entityType, "criteria": criteria})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrFailedMarshalling)
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

// ContentManagementAttachFilter - attach a filter to a specific project
//
// param: auth
// param: projectLabel
// param: filterID
// return:
func (p *Proxy) ContentManagementAttachFilter(auth AuthParams, projectLabel string, filterID int) (sumamodels.ContentManagementFilter, error) {
	log.Debug("contentManagement.attachFilter function called")
	var result sumamodels.ContentManagementFilter
	path := "contentmanagement/attachFilter"
	body, err := json.Marshal(map[string]any{"projectLabel": projectLabel, "filterId": filterID})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrFailedMarshalling)
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

// ContentManagementCreateEnvironment - create environment for a project
//
// param: auth
// param: projectLabel
// param: predecessorLabel
// param: envlabel
// param: name
// param: description
// return:
func (p *Proxy) ContentManagementCreateEnvironment(auth AuthParams, projectLabel string, predecessorLabel string, envlabel string, name string, description string) (sumamodels.ContentManagementEnvironment, error) {
	log.Debug("contentManagement.createEnvironment function called")
	var result sumamodels.ContentManagementEnvironment
	path := "contentmanagement/createEnvironment"
	body, err := json.Marshal(map[string]any{"projectLabel": projectLabel, "predecessorLabel": predecessorLabel, "envLabel": envlabel, "name": name, "description": description})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrFailedMarshalling)
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}

// ContentManagementBuildProject - build a project
//
// param: auth
// param: projectLabel
// return:
func (p *Proxy) ContentManagementBuildProject(auth AuthParams, projectLabel string) (int, error) {
	log.Debug("contentManagement.buildProject function called")
	var result int
	path := "contentmanagement/buildProject"
	body, err := json.Marshal(map[string]any{"projectLabel": projectLabel})
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrFailedMarshalling)
	}
	response, err := p.suse.SuseManagerCall(body, http.MethodPost, auth.Host, path, auth.SessionKey)
	if err != nil {
		log.Error(returnCodes.ErrFailedMarshalling, zap.Any("error", err))
		return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
	}
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrHandlingSuseManagerResponse)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			log.Error(returnCodes.ErrFailedUnMarshalling, zap.Any("error", err))
			return result, errors.New(returnCodes.ErrFailedUnMarshalling)
		}
	} else {
		log.Error(returnCodes.ErrHTTPSuseManagerResponse, zap.Any("HTTP Statuscode", response.StatusCode))
		return result, errors.New(returnCodes.ErrHTTPSuseManagerResponse)
	}
	return result, nil
}
