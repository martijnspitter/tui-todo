package osoperations

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFilePath(fileName string, isDev bool) string {
	if isDev {
		return fmt.Sprintf("./%s", fileName)
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(homeDir, ".tui-todo")
		os.MkdirAll(configDir, 0755)
		return filepath.Join(configDir, fileName)
	}

	return fmt.Sprintf("./%s", fileName)
}
