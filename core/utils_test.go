package core

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

func mockOsReadFile(t *testing.T, fileData []byte, err error) func(name string) ([]byte, error) {
	return func(filename string) ([]byte, error) {
		if filename != "/proc/modules" {
			t.Error("Should read /proc/modules")
		}
		return fileData, err
	}
}

func mockOsReadDir(t *testing.T, dirEntries []fs.DirEntry, err error) func(dirname string) ([]fs.DirEntry, error) {
	return func(dirname string) ([]fs.DirEntry, error) {
		if dirname != "/sys/bus/platform/drivers/ideapad_acpi/" {
			t.Error("Should read /sys/bus/platform/drivers/ideapad_acpi/")
		}
		return dirEntries, err
	}
}

type mockDirEntry struct {
	name  string
	isDir bool
}

func (m *mockDirEntry) Name() string {
	return m.name
}

func (m *mockDirEntry) IsDir() bool {
	return m.isDir
}

func (m *mockDirEntry) Type() fs.FileMode {
	return 0
}

func (m *mockDirEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}

const ideapadModuleLoaded = "ideapad_laptop 16384 0 - Live 0xffffffffa0a00000\nzram 36864 2 - Live 0x0000000000000000\n"
const ideapadModuleUnloaded = "zram 36864 2 - Live 0x0000000000000000\n"

func TestIsIdeapadLaptopKmodLoaded(t *testing.T) {
	t.Run("Returns true when module is loaded", func(t *testing.T) {
		osReadFile = mockOsReadFile(t, []byte(ideapadModuleLoaded), nil)
		if !IsIdeapadLaptopKmodLoaded() {
			t.Error("Should return true when ideapad_laptop module is loaded")
		}
	})

	t.Run("Returns false when module isn't loaded", func(t *testing.T) {
		osReadFile = mockOsReadFile(t, []byte(ideapadModuleUnloaded), nil)
		if IsIdeapadLaptopKmodLoaded() {
			t.Error("Should return false when ideapad_laptop module isn't loaded")
		}
	})

	t.Run("Returns false when error is encountered", func(t *testing.T) {
		osReadFile = mockOsReadFile(t, nil, errors.New("error"))
		if IsIdeapadLaptopKmodLoaded() {
			t.Error("Should return false when encountering error reading /proc/modules")
		}
	})
}

func TestGetIdeapadAcpiVpcSysFsDir(t *testing.T) {
	t.Run("Returns error when ideapad_laptop module isn't loaded", func(t *testing.T) {
		osReadFile = mockOsReadFile(t, []byte(ideapadModuleUnloaded), nil)
		_, err := GetIdeapadAcpiVpcSysFsDir()
		if err != ErrIdeapadLaptopKmodNotLoaded {
			t.Error("Should return ErrIdeapadLaptopKmodNotLoaded when ideapad_laptop module isn't loaded")
		}
	})

	t.Run("Returns error when error is encountered", func(t *testing.T) {
		osReadFile = mockOsReadFile(t, []byte(ideapadModuleLoaded), nil)
		osReadDir = mockOsReadDir(t, nil, errors.New("error"))
		_, err := GetIdeapadAcpiVpcSysFsDir()
		if err == nil {
			t.Error("Should return error when error is encountered")
		}
	})

	t.Run("Returns error when no ideapad_acpi VPC dir is found", func(t *testing.T) {
		osReadFile = mockOsReadFile(t, []byte(ideapadModuleLoaded), nil)
		osReadDir = mockOsReadDir(t, []fs.DirEntry{}, nil)
		_, err := GetIdeapadAcpiVpcSysFsDir()
		if err != ErrIdeapadLaptopKmodNotLoaded {
			t.Error("Should return ErrIdeapadLaptopKmodNotLoaded when no ideapad_acpi VPC dir is found")
		}
	})

	t.Run("Returns path to ideapad_acpi VPC dir when found", func(t *testing.T) {
		osReadFile = mockOsReadFile(t, []byte(ideapadModuleLoaded), nil)
		osReadDir = mockOsReadDir(t, []fs.DirEntry{&mockDirEntry{"VPC2004:00", true}}, nil)
		dir, err := GetIdeapadAcpiVpcSysFsDir()
		if err != nil {
			t.Error("Should not return error when ideapad_acpi VPC dir is found")
		}
		const expectedPath = "/sys/bus/platform/drivers/ideapad_acpi/VPC2004:00"
		if dir != expectedPath {
			t.Errorf("Should return path to ideapad_acpi VPC dir when only one is found, expected: %s, found: %s", expectedPath, dir)
		}
	})
}

func TestWriteToFile(t *testing.T) {
	t.Run("Returns error when error is encountered", func(t *testing.T) {
		osOpenFile = func(name string, flag int, perm fs.FileMode) (*os.File, error) {
			return nil, errors.New("error")
		}
		err := writeToFile("/test", []byte{})
		if err == nil {
			t.Error("Should return error when error is encountered")
		}
	})
}
