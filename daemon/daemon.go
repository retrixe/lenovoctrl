package daemon

import (
	"log"
	"os"

	_ "embed"

	"github.com/godbus/dbus/v5"
)

const version = "1.0.0-alpha.0"

var conn *dbus.Conn
var stdLog = log.New(os.Stdout, "", log.LstdFlags)
var errLog = log.New(os.Stderr, "", log.LstdFlags)

func InitialiseDBusConnection() {
	connection, err := dbus.ConnectSystemBus()
	if err != nil {
		log.Fatalln("Failed to connect to D-Bus system bus!", err)
	} else {
		stdLog.Println("Successfully connected to D-Bus system bus!")
		conn = connection
	}
}

func main() {
	log.SetPrefix("[lenovoctrl] ")
	log.SetOutput(os.Stderr)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		stdLog.Println("lenovoctrl version " + version)
		return
	} else if len(os.Args) >= 2 {
		stdLog.Println("Correct usage: ./lenovoctrl [-v or --version]")
		return
	}

	err := LoadConfig()
	if os.IsNotExist(err) {
		stdLog.Println("Creating new config file at " + GetConfigPath() + ".")
		SaveConfig()
	} else if err != nil {
		log.Fatalln("Failed to read config file for unknown reasons!", err)
	}

	err = ApplyConfig()
	if err != nil {
		log.Println("Failed to apply config! Some settings may not have been applied.", err)
	}

	InitialiseDBusConnection()
	defer conn.Close()
	StartDBusDaemon()
}
