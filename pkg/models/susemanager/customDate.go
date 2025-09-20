// Package sumamodels - structs needed for SUSE Manager API Calls
package sumamodels

import (
	"encoding/json"
	returnCodes "mlmtool/pkg/util/returnCodes"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// CustomDate - date format used in api calls
type CustomDate time.Time

// UnmarshalJSON - unmarshal JSON
//
// param: b
// return: error
func (j *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	timeCorrect := false
	t, err := time.Parse("2006-01-02T15:04:05-0700", s)
	if err == nil {
		timeCorrect = true
	}
	if !timeCorrect {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err == nil {
			timeCorrect = true
		}

	}
	if !timeCorrect {
		t, err = time.Parse("Jan 2, 2006, 15:04:05 PM", s)
		if err == nil {
			timeCorrect = true
		}

	}
	if !timeCorrect {
		return errors.New(returnCodes.ErrConversionTime)
	}
	*j = CustomDate(t)
	return nil
}

// MarshalJSON - marshal JSON
//
// return: time
// return: error
func (j CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}
