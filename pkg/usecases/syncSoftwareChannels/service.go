// Package syncsoftwarechannels - update the software channels on SUSE Manager
package syncsoftwarechannels

// ISyncSoftwareChannels - call syncSoftwareChannels
type ISyncSoftwareChannels interface {
	SyncSoftwareChannels() error
}
