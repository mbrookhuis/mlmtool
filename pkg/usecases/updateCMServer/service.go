// Package updatecmserver update the given server to the given osrelease
package updatecmserver

// IUpdateCMServer update the give server to the given osrelease
type IUpdateCMServer interface {
	UpdateCMServer() error
}
