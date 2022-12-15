package applet

import (
	_ "embed"
	"runtime"

	"github.com/getlantern/systray"
)

//go:embed icons/applet.png
var appletIcon []byte

func RunApplet() {
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
			"Enable battery conservation mode (caps battery to 60% charge)",
			false,
		)
	// TODO: mBatteryConservationMode.SetIcon(icon)
	go func() {
		for {
			<-mBatteryConservationMode.ClickedCh
			if mBatteryConservationMode.Checked() {
				// TODO: Disable battery conservation mode.
				println("disabling battery conservation mode")
				mBatteryConservationMode.Uncheck()
			} else {
				// TODO: Enable battery conservation mode.
				println("enabling battery conservation mode")
				mBatteryConservationMode.Check()
			}
		}
	}()

	systray.AddSeparator()

	// Add quit button.
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	// TODO: mQuit.SetIcon(icon.Data)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}
