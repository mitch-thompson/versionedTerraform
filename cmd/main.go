package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"versionedTerraform"
)

const (
	configFileLocation   = "config"
	shortConfigDirString = "/.versionedTerraform"
	pwd                  = "."
	terraformPrefix      = "/terraform_"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	configDirString := homeDir + shortConfigDirString
	configDir := os.DirFS(configDirString)
	workingDir := os.DirFS(pwd)
	var versionsFromConfig []versionedTerraform.SemVersion

	flag.Parse()
	args := flag.Args()

	needsUpdate, err := versionedTerraform.NeedToUpdateAvailableVersions(configDir, configFileLocation)
	if needsUpdate {
		fileHandle, _ := os.OpenFile(configDirString+"/"+configFileLocation, os.O_RDWR, 0666)
		defer fileHandle.Close()
		versionedTerraform.UpdateConfig(*fileHandle)
	}
	versionsFromConfig, err = versionedTerraform.LoadVersionsFromConfig(configDir, configFileLocation)
	if err != nil {
		fmt.Printf("Unable to read config: %v", err)
		os.Exit(1)
	}

	installedVersions, err := versionedTerraform.LoadInstalledVersions(configDir)
	if err != nil {
		fmt.Printf("Unable to verify installed verisons: %v", err)
		os.Exit(1)
	}

	var vSlice []string
	for _, v := range versionsFromConfig {
		vSlice = append(vSlice, v.ToString())
	}

	ver, err := versionedTerraform.GetVersionFromFile(workingDir, vSlice)
	if err != nil {
		fmt.Printf("Unable to retrieve terraform version from files: %v", err)
	}

	if !ver.Version.VersionInSlice(installedVersions) {
		fmt.Printf("Installing terraform version %s\n\n", ver.Version.ToString())
		err = ver.InstallTerraformVersion()
		if err != nil {
			fmt.Printf("Unable to install terraform version: %v", err)
		}
	}

	terraformFile := configDirString + terraformPrefix + ver.VersionToString()
	argsForTerraform := append([]string{""}, args...)
	cmd := exec.Cmd{
		Path:   terraformFile,
		Args:   argsForTerraform,
		Env:    os.Environ(),
		Dir:    pwd,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	cmd.Run()
}
