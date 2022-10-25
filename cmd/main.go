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

var needsStable = true

func main() {
	homeDir, _ := os.UserHomeDir()
	configDirString := homeDir + shortConfigDirString

	// Create configuration directory if it does not exist
	_, err := os.Stat(configDirString)
	if os.IsNotExist(err) {
		err = versionedTerraform.CreateConfig(configDirString, configFileLocation)
	}

	// Create configuration file if it does not exist
	_, err = os.Stat(configDirString + "/" + configFileLocation)
	if os.IsNotExist(err) {
		err = versionedTerraform.CreateConfig(configDirString, configFileLocation)
	}

	// Let the user know if we couldn't create the config directory or file
	if err != nil {
		fmt.Printf("Unable to create config directory: %v", err)
	}

	configDir := os.DirFS(configDirString)
	workingDir := os.DirFS(pwd)
	var versionsFromConfig []versionedTerraform.SemVersion

	flag.Parse()
	args := flag.Args()

	//Load available versions from configuration file
	versionsFromConfig, err = versionedTerraform.LoadVersionsFromConfig(configDir, configFileLocation)
	if err != nil {
		fmt.Printf("Unable to read config: %v\n", err)
		os.Exit(1)
	}

	//Check if we need to update available versions with terraform's website
	//Then update configuration if we do
	//todo move this above loading the config
	needsUpdate, err := versionedTerraform.NeedToUpdateAvailableVersions(configDir, configFileLocation)
	if os.ErrNotExist == err {
		fmt.Printf("Unable to update version: %v\n", err)
	}

	if needsUpdate {
		fileHandle, _ := os.OpenFile(configDirString+"/"+configFileLocation, os.O_RDWR, 0666)
		defer fileHandle.Close()
		versionedTerraform.UpdateConfig(*fileHandle)
	}

	// Load a slice of versions which have already been installed
	installedVersions, err := versionedTerraform.LoadInstalledVersions(configDir)
	if err != nil {
		fmt.Printf("Unable to verify installed verisons: %v", err)
		os.Exit(1)
	}

	var vSlice []string
	for _, v := range versionsFromConfig {
		vSlice = append(vSlice, v.ToString())
	}

	// Check if stable version of terraform is required
	needsStable, err = versionedTerraform.ConfigRequiresStable(configDir, configFileLocation)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open config file, defaulting to stable versions of terraform only")
	}

	// Load version required from terraform directory
	ver, err := versionedTerraform.GetVersionFromFile(workingDir, vSlice, needsStable)
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

	// Execute terraform
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
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
	}
}
