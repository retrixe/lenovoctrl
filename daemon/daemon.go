package daemon

import (
	"log"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/retrixe/lenovoctrl/core"
)

// https://dbus.freedesktop.org/doc/dbus-specification.html
// https://dbus.freedesktop.org/doc/dbus-tutorial.html

var stdLog = log.New(os.Stdout, "", log.LstdFlags)

func StartDaemon() {
	log.SetOutput(os.Stderr)

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

	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		log.Fatalln("Failed to connect to D-Bus system bus!", err)
	} else {
		stdLog.Println("Successfully connected to D-Bus system bus!")
	}
	defer conn.Close()
	ListenToDBus(conn)
}

const intro = introspect.IntrospectDeclarationString + `
<node>
	<interface name="com.retrixe.LenovoCtrl.v0">
		<method name="GetConservationModeStatus">
			<arg direction="out" type="n"/>
		</method>
		<method name="SetConservationMode">
		  <arg direction="in" type="b"/>
		</method>
		<method name="GetKeyboardFnLockStatus">
			<arg direction="out" type="n"/>
		</method>
		<method name="SetKeyboardFnLock">
		  <arg direction="in" type="b"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`

func ListenToDBus(conn *dbus.Conn) {
	f := DBusAPI("lenovoctrl v0 API")
	conn.Export(f, "/com/retrixe/LenovoCtrl/v0", "com.retrixe.LenovoCtrl.v0")
	conn.Export(introspect.Introspectable(intro), "/com/retrixe/LenovoCtrl/v0",
		"org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName("com.retrixe.LenovoCtrl.v0", dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatalln("Failed to request D-Bus name com.retrixe.LenovoCtrl.v0", err)
	} else if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalln("D-Bus name com.retrixe.LenovoCtrl.v0 already taken")
	}

	stdLog.Println("Listening on D-Bus name com.retrixe.LenovoCtrl.v0.")
	select {}
}

type DBusAPI string

func (f DBusAPI) GetConservationModeStatus() (int16, *dbus.Error) {
	status, err := core.IsConservationModeEnabled()
	if err == core.ErrConservationModeNotAvailable {
		return -1, nil
	} else if err != nil {
		return -1, dbus.MakeFailedError(err)
	} else if status {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (f DBusAPI) SetConservationMode(status bool) *dbus.Error {
	if err := core.SetConservationModeStatus(status); err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (f DBusAPI) GetKeyboardFnLockStatus() (int16, *dbus.Error) {
	status, err := core.IsKeyboardFnLockEnabled()
	if err == core.ErrKeyboardFnLockNotAvailable {
		return -1, nil
	} else if err != nil {
		return -1, dbus.MakeFailedError(err)
	} else if status {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (f DBusAPI) SetKeyboardFnLock(status bool) *dbus.Error {
	if err := core.SetFnLockStatus(status); err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}
