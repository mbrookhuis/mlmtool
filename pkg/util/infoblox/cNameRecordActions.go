package infoblox

// IBCreateCNameRecord
//
// param: paramsDns
func (i *Infoblox) IBCreateCNameRecord(paramsDns DNSRecord, name string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Create cName Record")
	err := i.Client.adCreateRecordCName(name, paramsDns.DNSRecord, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return err
	}
	return nil
}

// IBCheckCNameRecord
//
// param: paramsDns
// return:
func (i *Infoblox) IBCheckCNameRecord(name string) ([]ReplyRecordCName, error) {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check cName Record")
	result, err := i.Client.adGetRecordCName(name, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return nil, err
	}
	return result, nil
}

// IBCheckAndCreateCNameRecord
//
// param: paramsDns
func (i *Infoblox) IBCheckAndCreateCNameRecord(paramsDns DNSRecord, name string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check and create cName Record")
	result, err := i.IBCheckCNameRecord(name)
	if err != nil {
		return err
	}
	if len(result) < 1 {
		err = i.IBCreateCNameRecord(paramsDns, name)
		if err != nil {
			return err
		}
	}
	slogger.Infof("cName Record %v for %v created or already present", name, paramsDns.DNSRecord)
	return nil
}

// IBCheckCNameRecordCanonical
//
// param: canonical
// return:
func (i *Infoblox) IBCheckCNameRecordCanonical(canonical string) ([]ReplyRecordCName, error) {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check cName Record")
	result, err := i.Client.adGetRecordCNameCanonical(canonical, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return nil, err
	}
	return result, nil
}

// IBDeleteCNameRecord
//
// param: linkRef
func (i *Infoblox) IBDeleteCNameRecordLink(linkRef string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check cName Record")
	err := i.Client.adDeleteRecordCName(linkRef)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return err
	}
	return nil
}

// IBDeleteCNameRecordName
//
// param: name
func (i *Infoblox) IBDeleteCNameRecordName(name string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check and delete cName Record")
	result, err := i.IBCheckCNameRecord(name)
	if err != nil {
		return err
	}
	for _, record := range result {
		err = i.IBDeleteCNameRecordLink(record.Ref)
		if err != nil {
			return err
		}
	}
	slogger.Infof("cName Record %v deleted", name)
	return nil
}

// IBDeleteCNameRecordCanonical
//
// param: canonical
func (i *Infoblox) IBDeleteCNameRecordCanonical(canonical string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check and delete cName Record")
	result, err := i.IBCheckCNameRecordCanonical(canonical)
	if err != nil {
		return err
	}
	for _, record := range result {
		err = i.IBDeleteCNameRecordLink(record.Ref)
		if err != nil {
			return err
		}
	}
	slogger.Infof("cName Records for canonical %v deleted", canonical)
	return nil
}
