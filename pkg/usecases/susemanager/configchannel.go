// Package susemanager - SUSE Manager api call and support functions
package susemanager

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	sumamodels "mlmtool/pkg/models/susemanager"
)

// ConfigChannelListGlobals - list configchannels
//
// param: auth
// return:
func (p *Proxy) ConfigChannelListGlobals(auth AuthParams) ([]sumamodels.ConfigChannelListGlobals, error) {
	p.logger.Info("ConfigChannelListGlobals function call started")
	path := "configchannel/listGlobals"
	response, err := p.suse.SuseManagerCall(nil, "GET", auth.Host, path, auth.SessionKey)
	if err != nil {
		p.logger.Error("Error message recieved from suse-manger", zap.Any("error", err))
		return nil, fmt.Errorf("error while getting list of configuration channels. Error: %s", err)
	}
	var result []sumamodels.ConfigChannelListGlobals
	if response.StatusCode == 200 {
		resp, err := HandleSuseManagerResponse(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error while handling suse manager response err: %s", err)
		}
		byteArray, _ := json.Marshal(resp)
		err = json.Unmarshal(byteArray, &result)
		if err != nil {
			p.logger.Error("unmarshling error", zap.Any("error", err))
			return nil, fmt.Errorf("unable to process the received data. err: %s", err)
		}
	} else {
		return nil, fmt.Errorf("fetching configuration channel list failed. Http StatusCode: %s Http Response body: %s", fmt.Sprint(response.StatusCode), fmt.Sprint(string(response.Body)))
	}
	return result, nil
}
