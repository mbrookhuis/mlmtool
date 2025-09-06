package infoblox

import "fmt"

// adDeletePTRRecord
//
// param: linkRef
func (i *Client) adDeletePTRRecord(linkRef string) error {
	err := i.delete(fmt.Sprintf("%v%v", linkRef, ptrRecordAPI))
	return err
}
