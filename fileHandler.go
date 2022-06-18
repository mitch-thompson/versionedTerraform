package versionedTerraform

import (
	"bufio"
	"io/fs"
	"regexp"
	"strings"
)

var needsStable bool

//GetVersionFromFile returns Version pointer and error
//Iterates through files in current directory and sending to parseVersionFromFile
//todo this should be (Version) GetVers...
func GetVersionFromFile(fileSystem fs.FS, versionList []string, needsStableValue bool) (*Version, error) {
	needsStable = needsStableValue
	var versionFinal Version
	versionFinal = *NewVersion(">= 0.0.0", versionList)
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return &versionFinal, err
	}

	for _, f := range dir {
		version, isFinished, err := parseVersionFromFile(fileSystem, f.Name(), versionList)
		if err != nil {
			return &versionFinal, err
		}
		if isFinished {
			versionFinal = *version
			break
		}
	}
	return &versionFinal, nil
}

//todo same here
//parseVersionFromFile returns Version pointer, bool, and error
//bool returns true if required_version is found
//if required_version is found sends to NewVersion to determine required version for the application
func parseVersionFromFile(f fs.FS, fileName string, versionList []string) (*Version, bool, error) {
	fileHandle, err := f.Open(fileName)
	regex := regexp.MustCompile("required_version\\s+?=")
	isComment := "^\\s?#"
	removeQuotes := regexp.MustCompile("\"")
	if err != nil {
		return &Version{}, false, err
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		_line := fileScanner.Text()
		isComment, _ := regexp.MatchString(isComment, _line)
		if strings.Contains(_line, "required_version") && !isComment {
			_line = regex.ReplaceAllString(_line, "")
			_line = removeQuotes.ReplaceAllString(_line, "")
			return NewVersion(_line, versionList), true, nil
		}
	}

	return &Version{}, false, nil
}
