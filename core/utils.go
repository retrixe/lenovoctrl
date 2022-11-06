package core

import (
	"errors"
	"os"
	"strings"
)

const ideapadAcpiSysFsDir = "/sys/bus/platform/drivers/ideapad_acpi/"

// ErrIdeapadLaptopKmodNotLoaded indicates that the ideapad_laptop kernel module is not loaded.
//
// This package requires ideapad_laptop for most of its functions to work.
var ErrIdeapadLaptopKmodNotLoaded = errors.New("ideapad_laptop kernel module is not loaded")

// IsIdeapadLaptopKmodLoaded checks whether or not the ideapad_laptop kernel module is loaded.
//
// This package requires ideapad_laptop for most of its functions to work.
func IsIdeapadLaptopKmodLoaded() bool {
	modulesInfo, err := os.ReadFile("/proc/modules")
	if err != nil {
		return false
	}
	modules := strings.Split(string(modulesInfo), "\n")
	for _, module := range modules {
		if strings.HasPrefix(module, "ideapad_laptop ") {
			return true
		}
	}
	return false
}

// writeToFile writes data to the named file.
// It will NOT create the file or attempt to truncate it.
func writeToFile(name string, data []byte) error {
	f, err := os.OpenFile(name, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
