// TODO: This package, package util/suma and usecases/susemanager should be merge composed
package checksumaserver

import (
	"errors"
	"os"
)

const (
	fileSpacecmd      = "/root/.spacecmd/config"
	fileSumaPrimary   = "/root/suma_primary"
	fileSumaSecondary = "/root/suma_secondary"
)

func Primary() bool {
	if _, err := os.Stat(fileSumaPrimary); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func Secondary() bool {
	if _, err := os.Stat(fileSumaSecondary); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func SumaServer() bool {
	if _, err := os.Stat(fileSpacecmd); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
