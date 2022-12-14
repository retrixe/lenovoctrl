package applet

import (
	_ "embed"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus/v5"
)

//go:embed icon.png
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
	systray.SetIcon(appletIcon)
	systray.SetTooltip("Control Lenovo laptop settings.")
	if runtime.GOOS != "linux" {
		systray.SetTitle("lenovoctrl")
	}

	// Add conservation mode button.
	mBatteryConservationMode :=
		systray.AddMenuItemCheckbox(
			"Battery Conservation Mode",
			"Toggle battery conservation mode (caps battery to 60% charge).",
			false,
		)
	// Fetch status every second.
	go func() {
		for {
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
			<-time.After(1 * time.Second)
		}
	}()
	// Add click handler.
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
	mQuit := systray.AddMenuItem("Quit Lenovoctrl", "Close the Lenovoctrl applet.")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}
