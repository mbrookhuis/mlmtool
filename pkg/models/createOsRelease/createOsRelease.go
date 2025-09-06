// Package createosrelease structs needed to create osRelease
package createosrelease

// OsReleaseRecord struct holding channel info
type OsReleaseRecord struct {
	Label                string
	ParentChannel        string
	TreePath             string
	ChildChannelsDefault []string
	ChildChannelsExtra   []string
}

// OsReleaseFilterCriteria struct holding information to create a filter
type OsReleaseFilterCriteria struct {
	Field   string
	Matcher string
	Value   string
}
