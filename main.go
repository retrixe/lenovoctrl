package main

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/retrixe/lenovoctrl/applet"
	"github.com/retrixe/lenovoctrl/daemon"
)

const version = "1.0.0-alpha.0"

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		println("lenovoctrl version " + version)
		return
	} else if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Fprintln(os.Stderr, "./lenovoctrl [-v or --version]")
		fmt.Fprintln(os.Stderr, "./lenovoctrl [-d or --daemon]")
		os.Exit(1)
	} else if len(os.Args) == 2 && (os.Args[1] == "-d" || os.Args[1] == "--daemon") {
		daemon.StartDaemon()
	} else if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "Incorrect usage. Run lenovoctrl --help for more information.")
		os.Exit(1)
	}

	applet.RunApplet()
}
