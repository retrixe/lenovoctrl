package core

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var ErrLenovoConservationModeNotAvailable = errors.New(
	"lenovo conservation mode is not available on this system")

const IdeapadAcpiSysFs = "/sys/bus/platform/drivers/ideapad_acpi/"

func GetConservationModeSysFs() (string, error) {
	folders, err := os.ReadDir(IdeapadAcpiSysFs)
	if err != nil {
		return "", err
	}
	for _, folder := range folders {
		if strings.HasPrefix(folder.Name(), "VPC") {
			conservationModeSysFsPath := path.Join(IdeapadAcpiSysFs, folder.Name(), "conservation_mode")
			if stat, err := os.Lstat(conservationModeSysFsPath); err != nil || stat.IsDir() {
				continue
			}
			return conservationModeSysFsPath, nil
		}
	}
	return "", ErrLenovoConservationModeNotAvailable
}

func IsConservationModeAvailable() bool {
	modulesInfo, err := exec.Command("lsmod").Output()
	if err != nil {
		log.Println("Failed to run lsmod!")
		return false
	}
	modules := strings.Split(string(modulesInfo), "\n")
	for _, module := range modules {
		if strings.Fields(module)[0] == "ideapad_laptop" {
			conservationModeSysFs, err := GetConservationModeSysFs()
			if err == ErrLenovoConservationModeNotAvailable {
				return false
			} else if err != nil {
				log.Println("An unknown error occurred when checking for Lenovo conservation mode", err)
				return false
			}
			_, err = os.ReadFile(conservationModeSysFs)
			if os.IsNotExist(err) {
				return false
			} else if err != nil {
				log.Println("An unknown error occurred when checking for Lenovo conservation mode", err)
				return false
			}
			return true
		}
	}
	return false
}

func IsConservationModeEnabled() bool {
	conservationModeSysFs, err := GetConservationModeSysFs()
	if err == ErrLenovoConservationModeNotAvailable {
		return false
	} else if err != nil {
		log.Println("An unknown error occurred when checking for Lenovo conservation mode", err)
		return false
	}
	data, err := os.ReadFile(conservationModeSysFs)
	if os.IsNotExist(err) {
		log.Println("Lenovo conservation mode status was checked despite no support for it", err)
		return false
	} else if err != nil {
		log.Println("An unknown error occurred when checking for Lenovo conservation mode", err)
		return false
	}
	return string(data) == "1\n"
}

func SetConservationModeStatus(mode bool) error {
	if !IsConservationModeAvailable() { // Don't accidentally write to the file.
		return ErrLenovoConservationModeNotAvailable
	}
	conservationModeSysFs, err := GetConservationModeSysFs()
	if err != nil {
		log.Println("Failed to set Lenovo conservation mode", err)
		return err
	}
	data := []byte("0")
	if mode {
		data = []byte("1")
	}
	err = WriteFile(conservationModeSysFs, data)
	if err != nil {
		log.Println("Failed to set Lenovo conservation mode", err)
	}
	return err
}
