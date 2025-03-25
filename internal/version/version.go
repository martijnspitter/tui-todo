package version

// Version is set during build using ldflags
var version = "dev"

// GetVersion returns the current version of the application
func GetVersion() string {
	return version
}
