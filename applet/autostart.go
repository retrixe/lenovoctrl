package applet

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
)

const desktopFile = `[Desktop Entry]
Type=Application
Version=1.0
Name=Lenovoctrl Applet
GenericName=Lenovoctrl Applet
Comment=Linux applet to control aspects of Lenovo IdeaPad/Legion devices.
Exec=$exec
Icon=lenovoctrl
Terminal=false
Categories=Settings;Application;
`

func IsAutostartEnabled() (bool, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(filepath.Join(homedir, ".config", "autostart", "lenovoctrl.desktop"))
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func SetAutostartEnabled(enabled bool) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if enabled {
		err = os.MkdirAll(filepath.Join(homedir, ".config", "autostart"), 0755)
		if err != nil {
			return err
		}
		// If /usr/share/applications/lenovoctrl.desktop or ~/.local/share/lenovoctrl.desktop exist,
		// use those, else make a new one.
		desktopFileToWrite := []byte(strings.Replace(desktopFile, "$exec", os.Args[0], 1))
		globalDesktopFilePath := "/usr/share/applications/lenovoctrl.desktop"
		localDesktopFilePath := filepath.Join(homedir, ".local", "share", "applications", "lenovoctrl.desktop")
		if localDesktopFile, err := os.ReadFile(localDesktopFilePath); err == nil {
			desktopFileToWrite = localDesktopFile
		} else if globalDesktopFile, err := os.ReadFile(globalDesktopFilePath); err == nil {
			desktopFileToWrite = globalDesktopFile
		}
		err = os.WriteFile(filepath.Join(homedir, ".config", "autostart", "lenovoctrl.desktop"), desktopFileToWrite, 0644)
		if err != nil {
			return err
		}
	} else {
		err = os.Remove(filepath.Join(homedir, ".config", "autostart", "lenovoctrl.desktop"))
		if err != nil {
			return err
		}
	}
	return nil
}
