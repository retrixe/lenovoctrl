package core

import (
	"errors"
	"os"
	"path"
)

// ErrKeyboardFnLockNotAvailable indicates that keyboard Fn Lock is not available on this system.
var ErrKeyboardFnLockNotAvailable = errors.New(
	"keyboard fn lock is not available on this system")

// IsKeyboardFnLockEnabled checks if FnLock is enabled.
//
// If keyboard Fn Lock is not available on this system,
// this function returns ErrKeyboardFnLockNotAvailable.
func IsKeyboardFnLockEnabled() (bool, error) {
	fnLockSysFs, err := getFnLockSysFs()
	if err != nil {
		return false, err
	}
	data, err := os.ReadFile(fnLockSysFs)
	if os.IsNotExist(err) {
		return false, err
	} else if err != nil {
		return false, err
	}
	return string(data) == "1\n", nil
}

// SetFnLockStatus enables/disables keyboard Fn Lock.
//
// If keyboard Fn Lock is not available on this system,
// this function returns ErrKeyboardFnLockNotAvailable.
func SetFnLockStatus(mode bool) error {
	fnLockSysFs, err := getFnLockSysFs()
	if err != nil {
		return err
	}
	data := []byte("0")
	if mode {
		data = []byte("1")
	}
	err = writeToFile(fnLockSysFs, data)
	return err
}

func getFnLockSysFs() (string, error) {
	ideapadAcpiVpcSysFsDir, err := GetIdeapadAcpiVpcSysFsDir()
	if err != nil {
		return "", err
	}
	fnLockSysFsPath := path.Join(ideapadAcpiVpcSysFsDir, "fn_lock")
	if stat, err := os.Lstat(fnLockSysFsPath); err != nil || stat.IsDir() {
		return "", ErrKeyboardFnLockNotAvailable
	}
	return fnLockSysFsPath, nil
}
