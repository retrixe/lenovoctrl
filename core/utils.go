package core

import "os"

// WriteFile writes data to the named file.
// It does NOT create the file if it does not exist.
func WriteFile(name string, data []byte) error {
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
