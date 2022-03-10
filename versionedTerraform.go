package versionedTerraform

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func (*Version) latestMajorVersion() {
}

func (*Version) latestMinorVersion() {
}

func (*Version) latestPatchVersion() {
}

type Version struct {
	Version           SemVersion
	availableVersions []SemVersion
	installedVersions []SemVersion
}

const (
	hashicorpUrl             = "https://releases.hashicorp.com/terraform/"
	terraformPrefix          = "terraform_"
	versionedTerraformFolder = "/.versionedTerraform"
)

//getLatestMajorRelease() returns the latest major release from Version
func (v *Version) getLatestMajorRelease() {
	//todo clean up
	for _, release := range v.availableVersions {
		if release.majorVersion == v.Version.majorVersion &&
			release.minorVersion == v.Version.minorVersion &&
			release.patchVersion >= v.Version.patchVersion {
			v.Version = release
		}
	}
}

//getLatestRelease returns the latest release from Version
func (v *Version) getLatestRelease() {
	//todo clean up
	for _, release := range v.availableVersions {
		if release.majorVersion > v.Version.majorVersion {
			v.Version = release
		}
		if release.majorVersion >= v.Version.majorVersion &&
			release.minorVersion > v.Version.minorVersion {
			v.Version = release
		}
		if release.majorVersion >= v.Version.majorVersion &&
			release.minorVersion >= v.Version.minorVersion &&
			release.patchVersion >= v.Version.patchVersion {
			v.Version = release
		}
	}
}

//InstallTerraformVersion installs the defined terraform Version in the application
//configuration directory
func (v *Version) InstallTerraformVersion() error {
	homeDir, _ := os.UserHomeDir()
	resp, err := http.Get(hashicorpUrl +
		v.Version.ToString() +
		"/" + terraformPrefix +
		v.Version.ToString() +
		fileSuffix)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	versionedFileName := homeDir + versionedTerraformFolder + "/" + terraformPrefix + v.Version.ToString()
	versionedFile, err := os.OpenFile(versionedFileName, os.O_WRONLY, 0755)
	if os.IsNotExist(err) {
		//_, err = os.Create(versionedFileName)
		//if err != nil {
		//	return err
		//}
		versionedFile, err = os.OpenFile(versionedFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	defer versionedFile.Close()

	for _, zipFIle := range zipReader.File {
		zr, err := zipFIle.Open()
		if err != nil {
			return err
		}
		unzippedFileBytes, _ := ioutil.ReadAll(zr)

		_, err = versionedFile.Write(unzippedFileBytes)
		if err != nil {
			return err
		}
		zr.Close()
	}
	return nil
}

//NewVersion creates a new Version using sem versioning for determining the
//latest release
func NewVersion(_version string, _vList []string) *Version {
	v := new(Version)
	v.Version = *NewSemVersion(_version)

	for _, release := range _vList {
		v.availableVersions = append(v.availableVersions, *NewSemVersion(release))
	}

	switch {
	case strings.Contains(v.Version.ToString(), latestRelease):
		release := strings.Split(v.Version.ToString(), latestRelease)[1]
		v.Version = *NewSemVersion(release)
		v.getLatestRelease()
	case strings.Contains(v.Version.ToString(), latestPatch):
		release := strings.Split(v.Version.ToString(), latestPatch)[1]
		v.Version = *NewSemVersion(release)
		v.getLatestMajorRelease()
	}

	return v
}

//GetVersionList returns a list of available versions from hashicorp's release page
func GetVersionList() ([]string, error) {
	var versionList []string
	resp, err := http.Get(hashicorpUrl)
	if err != nil {
		return versionList, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return versionList, errors.New("invalid response code")
	}

	body, err := io.ReadAll(resp.Body)

	//todo maybe change this like GetVersionFromFile and consolidate
	bodyText := string(body)
	scanner := bufio.NewScanner(strings.NewReader(bodyText))
	for scanner.Scan() {
		_line := scanner.Text()
		if strings.Contains(_line, "a href") {
			_lineSplice := strings.Split(_line, terraformPrefix)
			if len(_lineSplice) == 2 {
				_line = _lineSplice[1]
				_line = strings.Split(_line, "</a")[0]
				versionList = append(versionList, _line)
			}
		}
	}
	return versionList, nil
}

//removeSpacesVersion removes spaces from Version string for parsing
func removeSpacesVersion(v string) string {
	splitV := strings.Split(v, " ")
	var returnString string
	for _, version := range splitV {
		version := strings.TrimSpace(version)
		returnString += version
	}
	return strings.TrimSpace(returnString)
}

//VersionToString returns string of a Version
func (v *Version) VersionToString() string {
	return v.Version.ToString()
}
