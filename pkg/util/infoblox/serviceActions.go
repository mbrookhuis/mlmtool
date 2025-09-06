package infoblox

// IBRestartInfoblox
//
// param: linkGrid
func (i *Infoblox) IBRestartInfoblox(linkGrid string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Restart Infoblox")
	err := i.Client.adRestartInfoblox(linkGrid)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return err
	}
	return nil
}

// IBGridInfoblox
//
// return:
func (i *Infoblox) IBGridInfoblox() (string, error) {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("grid Infoblox")
	result, err := i.Client.adGridInfoblox()
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return "", err
	}
	return result[0].Ref, nil
}
