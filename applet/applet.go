package applet

import (
	_ "embed"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus/v5"
	"github.com/ncruces/zenity"
)

//go:embed icon.png
var appletIcon []byte

var conn *dbus.Conn

func RunApplet() {
	var err error
	conn, err = dbus.ConnectSystemBus()
	if err != nil {
		zenity.Error("Failed to connect to system bus: " + err.Error())
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

	mBatteryConservationMode :=
		systray.AddMenuItemCheckbox(
			"Battery Conservation Mode",
			"Toggle battery conservation mode (caps battery to 60% charge).",
			false,
		)
	go func() {
		for {
			conservationMode, err := GetConservationModeStatus()
			if err != nil {
				zenity.Error("Failed to get conservation mode status: " + err.Error())
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
	go func() {
		for {
			<-mBatteryConservationMode.ClickedCh
			if mBatteryConservationMode.Checked() {
				if err := SetConservationMode(false); err != nil {
					zenity.Error("Failed to set conservation mode status: " + err.Error())
				}
				mBatteryConservationMode.Uncheck()
			} else {
				if err := SetConservationMode(true); err != nil {
					zenity.Error("Failed to set conservation mode status: " + err.Error())
				}
				mBatteryConservationMode.Check()
			}
		}
	}()

	mKeyboardFnLock :=
		systray.AddMenuItemCheckbox(
			"Keyboard Fn Lock",
			"Toggle keyboard Fn Lock.",
			false,
		)
	go func() {
		for {
			keyboardFnLock, err := GetKeyboardFnLockStatus()
			if err != nil {
				zenity.Error("Failed to get keyboard Fn Lock status: " + err.Error())
			} else if keyboardFnLock == -1 {
				mKeyboardFnLock.Hide()
			} else if keyboardFnLock == 0 {
				mKeyboardFnLock.Uncheck()
			} else {
				mKeyboardFnLock.Check()
			}
			<-time.After(1 * time.Second)
		}
	}()
	go func() {
		for {
			<-mKeyboardFnLock.ClickedCh
			if mKeyboardFnLock.Checked() {
				if err := SetKeyboardFnLock(false); err != nil {
					zenity.Error("Failed to set keyboard Fn Lock status: " + err.Error())
				}
				mKeyboardFnLock.Uncheck()
			} else {
				if err := SetKeyboardFnLock(true); err != nil {
					zenity.Error("Failed to set keyboard Fn Lock status: " + err.Error())
				}
				mKeyboardFnLock.Check()
			}
		}
	}()

	systray.AddSeparator()

	// Add autostart button.
	mAutostart :=
		systray.AddMenuItemCheckbox(
			"Autostart on Login",
			"Toggle the applet autostarting on login.",
			false,
		)
	go func() {
		for {
			autostartEnabled, err := IsAutostartEnabled()
			if err != nil {
				zenity.Error("Failed to get autostart status: " + err.Error())
			} else if autostartEnabled {
				mAutostart.Check()
			} else {
				mAutostart.Uncheck()
			}
			<-time.After(1 * time.Second)
		}
	}()
	go func() {
		for {
			<-mAutostart.ClickedCh
			if mAutostart.Checked() {
				if err := SetAutostartEnabled(false); err != nil {
					zenity.Error("Failed to set autostart status: " + err.Error())
				}
				mAutostart.Uncheck()
			} else {
				if err := SetAutostartEnabled(true); err != nil {
					zenity.Error("Failed to set autostart status: " + err.Error())
				}
				mAutostart.Check()
			}
		}
	}()

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
