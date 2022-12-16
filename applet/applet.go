package applet

import (
	_ "embed"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus/v5"
)

//go:embed icons/applet.png
var appletIcon []byte

var conn *dbus.Conn

func RunApplet() {
	var err error
	conn, err = dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	systray.Run(onReady, onExit)
}

func onReady() {
	// Setup systray.
	// TODO: These functions need separation of concerns.
	systray.SetIcon(appletIcon)
	systray.SetTooltip("Control Lenovo laptop settings.")
	if runtime.GOOS != "linux" {
		systray.SetTitle("lenovoctrl")
	}

	// Add conservation mode button.
	mBatteryConservationMode :=
		systray.AddMenuItemCheckbox(
			"Battery Conservation Mode",
			"Enable battery conservation mode (caps battery to 60% charge)",
			false,
		)
	conservationMode, err := GetConservationModeStatus()
	if err != nil {
		panic(err)
	} else if conservationMode == -1 {
		mBatteryConservationMode.Hide()
	} else if conservationMode == 0 {
		mBatteryConservationMode.Uncheck()
	} else {
		mBatteryConservationMode.Check()
	}
	go func() {
		for {
			<-mBatteryConservationMode.ClickedCh
			if mBatteryConservationMode.Checked() {
				if err := SetConservationMode(false); err != nil {
					panic(err)
				}
				mBatteryConservationMode.Uncheck()
			} else {
				if err := SetConservationMode(true); err != nil {
					panic(err)
				}
				mBatteryConservationMode.Check()
			}
		}
	}()

	systray.AddSeparator()

	// Add quit button.
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}
