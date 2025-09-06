package infoblox

// IBCreateZoneAuth
//
// param: domain
func (i *Infoblox) IBCreateZoneAuth(domain string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Create ZoneAuth")
	err := i.Client.adCreateZoneAuth(domain, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return err
	}
	return nil
}

// IBCheckZoneAuth
//
// param: domain
// return:
func (i *Infoblox) IBCheckZoneAuth(domain string) ([]InfobloxReplyZoneAuth, error) {
	slogger := i.logger.Sugar()
	slogger.Info("Check ZoneAuth")
	result, err := i.Client.adGetZoneAuth(domain, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return nil, err
	}
	return result, nil
}

// IBCheckAndCreateZoneAuth
//
// param: domain
func (i *Infoblox) IBCheckAndCreateZoneAuth(domain string) error {
	slogger := i.logger.Sugar()
	slogger.Info("Check and create ZoneAuth")
	result, err := i.IBCheckZoneAuth(domain)
	if err != nil {
		return err
	}
	if len(result) < 1 {
		err = i.IBCreateZoneAuth(domain)
		if err != nil {
			return err
		}
		linkGrid, err := i.IBGridInfoblox()
		if err != nil {
			return err
		}
		err = i.IBRestartInfoblox(linkGrid)
		if err != nil {
			return err
		}
	}
	slogger.Infof("Record for %v created or already present", domain)
	return nil
}
