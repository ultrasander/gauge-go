package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/getgauge-contrib/gauge-go/constants"
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge/common"
)

var pluginDir = ""
var projectRoot = ""
var start = flag.Bool("start", false, "Start go runner")
var initialize = flag.Bool("init", false, "Initialize Go project structure")

func main() {
	flag.Parse()

	setPluginAndProjectRoots()
	if *start {
		startGo()
	} else if *initialize {
		initGo()
	} else {
		printUsage()
	}
}

func startGo() {
	os.Chdir(projectRoot)
	err := gauge.LoadGaugeImpls()
	if err != nil {
		fmt.Printf("Failed to build project: %s\nKilling go runner. \n", err.Error())
		os.Exit(1)
	}
}

func initGo() {
	stepImplDir := filepath.Join(projectRoot, constants.DefaultStepImplDir)
	createDirectory(stepImplDir)
	stepImplFile := filepath.Join(stepImplDir, constants.DefaultStepImplFileName)
	showMessage("create", stepImplFile)
	common.CopyFile(filepath.Join(constants.SkelDir, constants.DefaultStepImplFileName), stepImplFile)
}

func printUsage() {
	flag.PrintDefaults()
}

func showMessage(action, filename string) {
	fmt.Printf(" %s  %s\n", action, filename)
}

func setPluginAndProjectRoots() {
	var err error
	pluginDir, err = os.Getwd()
	if err != nil {
		fmt.Printf("Failed to find current working directory: %s \n", err)
		os.Exit(1)
	}
	projectRoot = os.Getenv(common.GaugeProjectRootEnv)
	if projectRoot == "" {
		fmt.Printf("Could not find %s env. Go Runner exiting...", common.GaugeProjectRootEnv)
		os.Exit(1)
	}

	if !checkIfInSrcPath(projectRoot) {
		fmt.Printf("Project folder must be a subfolder in GOPATH/src folder\n")
		os.Exit(1)
	}
}

func createDirectory(dirPath string) {
	showMessage("create", dirPath)
	if !common.DirExists(dirPath) {
		err := os.MkdirAll(dirPath, common.NewDirectoryPermissions)
		if err != nil {
			fmt.Printf("Failed to make directory. %s\n", err.Error())
		}
	} else {
		fmt.Println("skip ", dirPath)
	}
}

func getGoPaths() []string {
	var paths []string
	if runtime.GOOS == "windows" {
		paths = strings.Split(build.Default.GOPATH, ";")
	} else {
		paths = strings.Split(build.Default.GOPATH, ":")
	}
	return paths
}

func getGoSrcPaths() []string {
	var paths = getGoPaths()
	for i, p := range paths {
		paths[i] = filepath.Join(p, "src")
	}
	return paths
}

func checkIfInSrcPath(dirPath string) bool {
	for _, p := range getGoSrcPaths() {
		if filepath.HasPrefix(dirPath, p) {
			return true
		}
	}
	return false
}
