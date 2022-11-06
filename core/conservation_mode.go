package core

import (
	"errors"
	"os"
	"path"
	"strings"
)

// ErrConservationModeNotAvailable indicates that battery
// conservation mode is not available on this system.
//
// Conservation mode limits the battery charge to 60%.
var ErrConservationModeNotAvailable = errors.New(
	"conservation mode is not available on this system")

// IsConservationModeEnabled checks if battery conservation mode is enabled.
// Conservation mode limits the battery charge to 60%.
//
// If conservation mode is not available on this system, this function
// returns ErrConservationModeNotAvailable.
func IsConservationModeEnabled() (bool, error) {
	conservationModeSysFs, err := getConservationModeSysFs()
	if err != nil {
		return false, err
	}
	data, err := os.ReadFile(conservationModeSysFs)
	if os.IsNotExist(err) {
		return false, err
	} else if err != nil {
		return false, err
	}
	return string(data) == "1\n", nil
}

// SetConservationModeStatus enables/disables battery conservation mode.
// Conservation mode limits the battery charge to 60%.
//
// If conservation mode is not available on this system, this function
// returns ErrConservationModeNotAvailable.
func SetConservationModeStatus(mode bool) error {
	conservationModeSysFs, err := getConservationModeSysFs()
	if err != nil {
		return err
	}
	data := []byte("0")
	if mode {
		data = []byte("1")
	}
	err = writeToFile(conservationModeSysFs, data)
	return err
}

func getConservationModeSysFs() (string, error) {
	if !IsIdeapadLaptopKmodLoaded() {
		return "", ErrIdeapadLaptopKmodNotLoaded
	}

	folders, err := os.ReadDir(ideapadAcpiSysFsDir)
	if err != nil {
		return "", err
	}
	for _, folder := range folders {
		if strings.HasPrefix(folder.Name(), "VPC") {
			conservationModeSysFsPath := path.Join(ideapadAcpiSysFsDir, folder.Name(), "conservation_mode")
			if stat, err := os.Lstat(conservationModeSysFsPath); err != nil || stat.IsDir() {
				continue
			}
			return conservationModeSysFsPath, nil
		}
	}
	return "", ErrConservationModeNotAvailable
}
