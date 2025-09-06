package infoblox

import (
	"fmt"
	"strings"
)

// IBCreatePTRRecord
//
// param: paramsDns
func (i *Infoblox) IBCreatePTRRecord(paramsDns DNSRecord) error {
	slogger := i.logger.Sugar()
	slogger.Info("Create PTR Record")
	var ipArpaRecord string
	ipArpa := strings.Split(paramsDns.IP, ".")
	ipArpaRecord = fmt.Sprintf("%v.%v.%v.%v.in-addr.arpa", ipArpa[3], ipArpa[2], ipArpa[1], ipArpa[0])
	err := i.Client.adCreateRecordPTR(ipArpaRecord, paramsDns.DNSRecord, paramsDns.IP, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return err
	}
	return nil
}

/*
	Name     string `json:"name"`
	ptrDName string `json:"ptrdname"`
	IPv4Addr string `json:"ipv4addr"`
	View     string `json:"view"`

*/

// IBCheckPTRRecord
//
// param: paramsDns
// return:
func (i *Infoblox) IBCheckPTRRecord(paramsDns DNSRecord) ([]ReplyRecordPTR, error) {
	slogger := i.logger.Sugar()
	slogger.Info("Check PTR Record")
	result, err := i.Client.adGetRecordPTR(paramsDns.DNSRecord, i.View)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return nil, err
	}
	return result, nil
}

// IBCheckAndCreatePTRRecord
//
// param: paramsDns
func (i *Infoblox) IBCheckAndCreatePTRRecord(paramsDns DNSRecord) error {
	slogger := i.logger.Sugar()
	slogger.Info("Check and create PTR-Record")
	result, err := i.IBCheckPTRRecord(paramsDns)
	if err != nil {
		return err
	}
	if len(result) < 1 {
		err = i.IBCreatePTRRecord(paramsDns)
		if err != nil {
			return err
		}
	}
	for _, entry := range result {
		if paramsDns.DNSRecord != entry.PtrDName {
			return fmt.Errorf("record for %v already exists with incorrect ip should: %v", paramsDns.DNSRecord, paramsDns.IP)
		}
	}
	slogger.Infof("Record for %v created or already present", paramsDns.DNSRecord)
	return nil
}

func (i *Infoblox) IBDeletePTRRecordLink(linkRef string) error {
	slogger := i.logger.Sugar()
	slogger.Info("delete PTR Record link")
	err := i.Client.adDeletePTRRecord(linkRef)
	if err != nil {
		slogger.Errorf("Command failed. %v", err)
		return err
	}
	return nil
}

// IBDeletePTRRecord
//
// param: paramsDns
func (i *Infoblox) IBDeletePTRRecord(paramsDns DNSRecord) error {
	slogger := i.logger.Sugar()
	slogger.Info("Delete PTR-Record")
	result, err := i.IBCheckPTRRecord(paramsDns)
	if err != nil {
		return err
	}
	for _, record := range result {
		err = i.IBDeletePTRRecordLink(record.Ref)
		if err != nil {
			return err
		}
	}
	slogger.Infof("Record for %v deleted", paramsDns.DNSRecord)
	return nil
}
