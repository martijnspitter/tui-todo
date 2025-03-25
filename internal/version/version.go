package version

// Version is set during build using ldflags
var Version = "dev"

// GetVersion returns the current version of the application
func GetVersion() string {
	return Version
}
