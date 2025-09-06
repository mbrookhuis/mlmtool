package infoblox

import (
	"fmt"
)

// IBCreateARecord
//
// param: paramsDns
func (i *Infoblox) IBCreateARecord(paramsDns DNSRecord) error {
	i.logger.Info("Create A Record")
	err := i.Client.adCreateRecordA(paramsDns.DNSRecord, paramsDns.IP, i.View)
	if err != nil {
		i.logger.Error(fmt.Sprintf("Command failed. %v", err))
		return err
	}
	return nil
}

// IBCheckARecord
//
// param: paramsDns
// return:
func (i *Infoblox) IBCheckARecord(paramsDns DNSRecord) ([]ReplyRecordA, error) {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check A Record")
	result, err := i.Client.adGetRecordA(paramsDns.DNSRecord, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return nil, err
	}
	return result, nil
}

// IBCheckAndCreateARecord
//
// param: paramsDns
func (i *Infoblox) IBCheckAndCreateARecord(paramsDns DNSRecord) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Check and create A-Record")
	result, err := i.IBCheckARecord(paramsDns)
	if err != nil {
		return err
	}
	if len(result) < 1 {
		err = i.IBCreateARecord(paramsDns)
		if err != nil {
			return err
		}
	}
	for _, entry := range result {
		if paramsDns.IP != entry.Ipv4Addr {
			return fmt.Errorf("record for %v already exists with incorrect ip should: %v has: %v", paramsDns.DNSRecord, paramsDns.IP, entry.Ipv4Addr)
		}
	}
	slogger.Infof("Record for %v created or already present", paramsDns.DNSRecord)
	return nil
}

func (i *Infoblox) IBDeleteARecordLink(linkRef string) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("delete A Record link")
	err := i.Client.adDeleteARecord(linkRef)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return err
	}
	return nil
}

// IBDeleteARecord
//
// param: paramsDns
func (i *Infoblox) IBDeleteARecord(paramsDns DNSRecord) error {
	logger, slogger := i.logger, i.logger.Sugar()
	logger.Info("Delete A-Record")
	result, err := i.IBCheckARecord(paramsDns)
	if err != nil {
		return err
	}
	for _, record := range result {
		err = i.IBDeleteARecordLink(record.Ref)
		if err != nil {
			return err
		}
	}
	slogger.Infof("Record for %v deleted", paramsDns.DNSRecord)
	return nil
}
