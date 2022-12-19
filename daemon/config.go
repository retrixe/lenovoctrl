package daemon

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct{}

var config Config = Config{}

func GetConfigPath() string {
	snapdir := os.Getenv("SNAP_DATA")
	path := "/etc/lenovoctrl/config.json"
	if snapdir != "" {
		path = filepath.Join(snapdir, "config.json")
	}
	return path
}

func SaveConfig() {
	err := os.Mkdir(filepath.Dir(GetConfigPath()), 0755)
	if err != nil && !os.IsExist(err) {
		log.Println("Failed to write persistent config file! "+
			"Any changes made during this session may be lost on reboot!", err)
	}

	file, err := json.Marshal(&config)
	if err != nil {
		log.Println("Failed to write persistent config file! "+
			"Any changes made during this session may be lost on reboot!", err)
	}

	err = os.WriteFile(GetConfigPath(), file, os.ModePerm)
	if err != nil {
		log.Println("Failed to write persistent config file! "+
			"Any changes made during this session may be lost on reboot!", err)
	}
}

func LoadConfig() error {
	file, err := os.ReadFile(GetConfigPath())
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}
	return nil
}

func ApplyConfig() error {
	// This is disabled because this setting is stored in firmware.
	// Lenovo Conservation Mode.
	// if lenovo.IsConservationModeAvailable() {
	// 	err := lenovo.SetConservationModeStatus(config.LenovoConservationModeEnabled)
	// 	if err != nil {return err}}

	return nil
}
