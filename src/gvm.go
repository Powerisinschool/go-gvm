package main

import (
	"fmt"
	"gvm/web"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const GVM_VERSION = "0.1.0"

func main() {
	args := os.Args
	osArch := strings.ToLower(os.Getenv("PROCESSOR_ARCHITECTURE"))
	err := "no err"
	showErr := false
	// detail := ""

	if osArch == "x86" {
		osArch = "386"
	}

	fmt.Println(strings.Join(args, ","))

	if len(args) < 2 {
		err = "no arguments"
		help(err)
		return
	}

	switch args[1] {
	case "version":
		fmt.Printf("\nGVM Version: " + GVM_VERSION + "\n\n")
		return
	case "help":
		help(err)
		return
	case "?":
		help(err)
		return
	case "arch":
		fmt.Println("Build Architecture: " + osArch)
		return
	case "list":
		list()
		return
	case "install":
		showErr, err = checkSecondArg(args, install, osArch)
	case "use":
		showErr, err = checkSecondArg(args, use, osArch)
	case "uninstall":
		showErr, err = checkSecondArg(args, uninstall, osArch)
	case "goroot":
		showErr, err = checkSecondArg(args, goRoot, osArch)
	}

	if showErr {
		help(err)
	}
}

func list() {
	panic("unimplemented")
}

func use(version string, arch string) bool {
	panic("unimplemented")
}

func install(version string, arch string) bool {
	fmt.Println("")
	if os.Getenv("GOROOT") == "" {
		fmt.Println("No GOROOT set. Set a GOROOT for Go installations with gvm goroot <path>.")
		return false
	}
	if version == "" {
		fmt.Println("Version not specified.")
		return false
	}
	gorootroot := filepath.Clean(os.Getenv("GOROOT") + "\\..")
	return web.Download(version, "windows-"+arch, gorootroot)
	// panic("unimplemented")
}

func uninstall(version string, arch string) bool {
	panic("unimplemented")
}

func goRoot(path string, arch string) bool {
	fmt.Println("")
	if path == "" {
		if os.Getenv("GOROOT") == "" {
			fmt.Println("No GOROOT set.")
		} else {
			fmt.Println("GOROOT: ", os.Getenv("GOROOT"))
			fmt.Println("Other Go versions installed at: ", filepath.Clean(os.Getenv("GOROOT")+"\\.."))
		}
		return false
	}
	newpath := filepath.FromSlash(path)
	//permanently set env var for user and local machine
	//The path should be the same for all windows OSes.
	machineEnvPath := "SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment"
	userEnvPath := "Environment"
	setEnvVar("GOROOT", newpath, machineEnvPath, true)
	setEnvVar("GOROOT", newpath, userEnvPath, false)
	//Also update path for user and local machine
	updatePathVar("PATH", filepath.FromSlash(os.Getenv("GOROOT")), newpath, machineEnvPath, true)
	updatePathVar("PATH", filepath.FromSlash(os.Getenv("GOROOT")), newpath, userEnvPath, false)
	fmt.Println("Set the GOROOT to " + newpath + ". Also updated PATH.")
	fmt.Println("Note: You'll have to start another prompt to see the changes.")
	return false
}

func setEnvVar(envVar string, newVal string, envPath string, machine bool) {
	//this sets the environment variable (GOROOT in this case) for either LOCAL_MACHINE or CURRENT_USER.
	//They are set in the registry. both must be set since the GOROOT could be used from either location.
	regplace := registry.CURRENT_USER
	if machine {
		regplace = registry.LOCAL_MACHINE
	}
	key, err := registry.OpenKey(regplace, envPath, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer key.Close()

	err = key.SetStringValue(envVar, newVal)
	if err != nil {
		fmt.Println("error", err)
	}
}

func updatePathVar(envVar string, oldVal string, newVal string, envPath string, machine bool) {
	//this sets the environment variable for either LOCAL_MACHINE or CURRENT_USER.
	//They are set in the registry. both must be set since the GOROOT could be used from either location.
	regplace := registry.CURRENT_USER
	if machine {
		regplace = registry.LOCAL_MACHINE
	}
	key, err := registry.OpenKey(regplace, envPath, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer key.Close()

	val, _, kerr := key.GetStringValue(envVar)
	if kerr != nil {
		fmt.Println("error", err)
		return
	}
	pvars := strings.Split(val, ";")
	for i, pvar := range pvars {
		if pvar == newVal+"\\bin" {
			//the requested new value already exists in PATH, do nothing
			return
		}
		if pvar == oldVal+"\\bin" {
			pvars = append(pvars[:i], pvars[i+1:]...)
		}
	}
	val = strings.Join(pvars, ";")
	val = val + ";" + newVal + "\\bin"
	err = key.SetStringValue("PATH", val)
	if err != nil {
		fmt.Println("error", err)
	}
}


func help(msg string) {
	fmt.Printf("\nGVM Version: " + GVM_VERSION + "\n\n")
	if msg == "no arguments" {
		fmt.Printf("Expected at least 1 (one) argument to operate with")
	}

	fmt.Println(msg)

	fmt.Println("\nHelp Menu\nAccess this menu by using gvm help or ? at the end of any command")
	fmt.Println("\nUsage:")
	fmt.Printf(" \n")
	fmt.Println("  gvm arch                     : Show architecture of OS.")
	fmt.Println("  gvm install <version>        : The version must be a version of Go.")
	fmt.Println("  gvm goroot [path]            : Sets/appends GOROOT/PATH. Without the extra arg just shows current GOROOT.")
	fmt.Println("  gvm list                     : List the Go installations at or adjacent to GOROOT. Aliased as ls.")
	fmt.Println("  gvm uninstall <version>      : Uninstall specified version of Go. If it was your GOROOT/PATH, make sure to set a new one after.")
	fmt.Println("  gvm use <version>            : Switch to use the specified version. This will set your GOROOT and PATH.")
	fmt.Println("  gvm version                  : Displays the current running version of gvm for Windows. Aliased as v.")
}

type callbackFn func(version string, arch string) bool

func checkSecondArg(args []string, callback callbackFn, arch string) (bool, string) {
	if len(args) < 3 {
		return true, "insufficient arguments"
	}
	callback(args[2], arch)
	return false, "no err"
}
