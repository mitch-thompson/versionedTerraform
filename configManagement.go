package versionedTerraform

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type configStruct struct {
	StableOnly        bool
	LastUpdate        int64
	AvailableVersions []string
}

//ConfigRequiresStable returns bool, error only false if StableOnly: false is set in configuration file
func ConfigRequiresStable(File os.File) (bool, error) {
	fileHandle, err := os.Open(File.Name())
	if err != nil {
		return true, err
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		_line := fileScanner.Text()
		if strings.Contains(_line, "StableOnly: ") {
			isStable := strings.SplitAfter(_line, "StableOnly: ")[1]
			if strings.EqualFold(isStable, "false") {
				return false, nil
			}
		}
	}
	return true, nil
}

//NeedToUpdateAvailableVersions returns bool, error checks if last update was older than 1 day ago
// this prevents us from spamming the list of available terraform versions page
func NeedToUpdateAvailableVersions(fileSystem fs.FS, availableVersions string) (bool, error) {
	//todo this is used a lot abstract it?
	fileHandle, err := fileSystem.Open(availableVersions)
	oneDayAgo := time.Now().AddDate(0, 0, -1).Unix()
	if err != nil {
		return false, err
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		_line := fileScanner.Text()
		if strings.Contains(_line, "LastUpdate: ") {
			lastUpdateTimeString := strings.SplitAfter(_line, "LastUpdate: ")[1]
			lastUpdateTimeString = strings.TrimSpace(lastUpdateTimeString)
			lastUpdateTime, err := strconv.ParseInt(lastUpdateTimeString, 10, 64)
			if err != nil {
				return false, err
			}
			if lastUpdateTime <= oneDayAgo {
				return true, nil
			}
		}
	}
	return false, nil
}

//LoadVersionsFromConfig returns slice of SemVersions and an error from AvailableVersions in configuration file
//This is stored from GetVersionList()
func LoadVersionsFromConfig(fileSystem fs.FS, configFile string) ([]SemVersion, error) {
	fileHandle, err := fileSystem.Open(configFile)
	removeOpenBracket := regexp.MustCompile("\\[")
	removeCloseBracket := regexp.MustCompile("]")
	if err != nil {
		return nil, err
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		_line := fileScanner.Text()
		if strings.Contains(_line, "AvailableVersions: ") {
			var versionList []SemVersion
			_line = strings.SplitAfter(_line, "AvailableVersions: ")[1]
			_line = removeOpenBracket.ReplaceAllString(_line, "")
			_line = removeCloseBracket.ReplaceAllString(_line, "")
			versions := strings.Split(_line, " ")
			for _, version := range versions {
				versionList = append(versionList, *NewSemVersion(version))
			}
			return versionList, nil
		}
	}
	return nil, nil
}

//LoadInstalledVersions returns list of SemVersions and an error from the directory listing of the .versionedTerraform
func LoadInstalledVersions(fileSystem fs.FS) ([]SemVersion, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	var installedTerraformVersions []SemVersion
	terraformRegex := regexp.MustCompile(terraformPrefix)
	if err != nil {
		return nil, err
	}

	for _, f := range dir {
		terraformFileName := f.Name()
		if strings.Contains(terraformFileName, terraformPrefix) {
			terraformVersionString := terraformRegex.ReplaceAllString(terraformFileName, "")
			installedTerraformVersions = append(installedTerraformVersions, *NewSemVersion(terraformVersionString))
		}
	}
	return installedTerraformVersions, nil
}

//UpdateConfig returns an error, and updates configuration file
// adding:
// a new date to the last updated field
// the available versions listed on terraforms website
// the status of if the user wants only stable releases
func UpdateConfig(File os.File, timeNow ...time.Time) error {
	configValues := new(configStruct)

	configValues.AvailableVersions, _ = GetVersionList()
	configValues.StableOnly, _ = ConfigRequiresStable(File)

	var t time.Time
	if len(timeNow) > 0 {
		t = timeNow[0]
	} else {
		t = time.Now()
	}
	configValues.LastUpdate = t.Unix()

	File.Truncate(0)
	File.Seek(0, 0)

	lineToByte := []byte(fmt.Sprintf("StableOnly: %+v\n", configValues.StableOnly))
	File.Write(lineToByte)

	lineToByte = []byte(fmt.Sprintf("LastUpdate: %d\n", configValues.LastUpdate))
	File.Write(lineToByte)
	lineToByte = []byte(fmt.Sprintf("AvailableVersions: %+v\n", configValues.AvailableVersions))
	File.Write(lineToByte)
	return nil
}

//CreateConfig returns error, creates a new configuration file
func CreateConfig(directory string, configFile string) error {
	configFileName := directory + "/" + configFile
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	fileHandler, err := os.Create(configFileName)
	defer fileHandler.Close()

	lineToByte := []byte(fmt.Sprintf("StableOnly: true\n"))
	fileHandler.Write(lineToByte)
	err = UpdateConfig(*fileHandler)
	return err
}
