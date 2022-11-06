package daemon

import (
	"log"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/retrixe/lenovoctrl/core"
)

// https://dbus.freedesktop.org/doc/dbus-specification.html
// https://dbus.freedesktop.org/doc/dbus-tutorial.html

// TODO: Restrict API to wheel/sudo users.

const intro = introspect.IntrospectDeclarationString + `
<node>
	<interface name="com.retrixe.LenovoCtrl.v0">
		<method name="LenovoGetConservationModeStatus">
			<arg direction="out" type="n"/>
		</method>
		<method name="LenovoSetConservationMode">
		  <arg direction="in" type="b"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`

func StartDBusDaemon() {
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

func (f DBusAPI) LenovoGetConservationModeStatus() (int16, *dbus.Error) {
	if core.IsConservationModeAvailable() {
		if core.IsConservationModeEnabled() {
			return 1, nil
		} else {
			return 0, nil
		}
	} else {
		return -1, nil
	}
}

func (f DBusAPI) LenovoSetConservationMode(status bool) *dbus.Error {
	if err := core.SetConservationModeStatus(status); err != nil {
		return &dbus.Error{
			Name: "A daemon error occurred",
			Body: []interface{}{err.Error()},
		}
	}
	return nil
}
