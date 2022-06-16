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

func ConfigRequiresStable(fileSystem fs.FS, configFile string) (bool, error) {
	fileHandle, err := fileSystem.Open(configFile)
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

func UpdateConfig(File os.File) error {
	configValues := new(configStruct)

	configValues.AvailableVersions, _ = GetVersionList()

	timeNow := time.Now()
	configValues.LastUpdate = timeNow.Unix()

	File.Truncate(0)
	File.Seek(0, 0)

	lineToByte := []byte(fmt.Sprintf("LastUpdate: %d\n", configValues.LastUpdate))
	File.Write(lineToByte)
	lineToByte = []byte(fmt.Sprintf("AvailableVersions: %+v\n", configValues.AvailableVersions))
	File.Write(lineToByte)
	return nil
}

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
