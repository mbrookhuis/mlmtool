package infoblox

import "fmt"

// adDeleteRecordCName
//
// param: linkRef
func (i *Client) adDeleteRecordCName(linkRef string) error {
	err := i.delete(fmt.Sprintf("%v%v", linkRef, cnameRecordAPI))
	return err
}
