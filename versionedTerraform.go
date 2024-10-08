package versionedTerraform

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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

// getLatestMajorRelease() returns the latest major release from Version
func (v *Version) getLatestMajorRelease() {
	for _, release := range v.availableVersions {
		if release.majorVersion == v.Version.majorVersion &&
			release.minorVersion == v.Version.minorVersion &&
			release.patchVersion >= v.Version.patchVersion &&
			(v.Version.isStable || !needsStable) {
			v.Version = release
		}
	}
}

// getGreatestRelease() returns less than release
func (v *Version) getOneLessRelease() {
	var vSlice []Version

	for _, release := range v.availableVersions {
		_v := Version{
			Version:           release,
			availableVersions: v.availableVersions,
			installedVersions: v.installedVersions,
		}

		if isVersionGreater(*v, _v) {
			vSlice = append(vSlice, _v)
		}
	}

	for i, ver := range vSlice {
		if isVersionGreater(ver, *v) || i == 0 {
			*v = ver
		}
	}
}

// getLatestRelease returns the latest release from Version
func (v *Version) getLatestRelease() {
	for _, release := range v.availableVersions {
		if release.majorVersion > v.Version.majorVersion &&
			(release.isStable || !needsStable) {
			v.Version = release
		}
		if release.majorVersion >= v.Version.majorVersion &&
			release.minorVersion > v.Version.minorVersion &&
			(release.isStable || !needsStable) {
			v.Version = release
		}
		if release.majorVersion >= v.Version.majorVersion &&
			release.minorVersion >= v.Version.minorVersion &&
			release.patchVersion >= v.Version.patchVersion &&
			(release.isStable || !needsStable) {
			v.Version = release
		}
	}
}

// InstallTerraformVersion installs the defined terraform Version in the application
// configuration directory
func (v *Version) InstallTerraformVersion() error {
	homeDir, _ := os.UserHomeDir()
	suffix := fileSuffix
	minV := NewSemVersion(minVersion)
	if v.Version.IsLessThan(*minV) {
		suffix = alternateSuffix
	}
	url := hashicorpUrl +
		v.Version.ToString() +
		"/" + terraformPrefix +
		v.Version.ToString() +
		suffix

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download Terraform: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %v", err)
	}

	versionedFileName := homeDir + versionedTerraformFolder + "/" + terraformPrefix + v.Version.ToString()
	versionedFile, err := os.OpenFile(versionedFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer versionedFile.Close()

	for _, zipFile := range zipReader.File {
		if zipFile.Name != "terraform" {
			continue
		}

		zr, err := zipFile.Open()
		if err != nil {
			return fmt.Errorf("failed to open zip file: %v", err)
		}
		defer zr.Close()

		_, err = io.Copy(versionedFile, zr)
		if err != nil {
			return fmt.Errorf("failed to write terraform binary: %v", err)
		}

		break
	}

	return nil
}

// NewVersion creates a new Version using sem versioning for determining the
// latest release
func NewVersion(_version string, _vList []string) *Version {
	v := new(Version)
	v.Version = *NewSemVersion(_version)

	for _, release := range _vList {
		v.availableVersions = append(v.availableVersions, *NewSemVersion(release))
	}

	hasMultiVersions, _ := regexp.MatchString("\\d+,", v.Version.ToString())

	if hasMultiVersions {
		releases := strings.Split(v.Version.ToString(), ",")
		for iteration, _release := range releases {
			_v := new(Version)
			_v.availableVersions = v.availableVersions
			switch {
			case strings.Contains(_release, latestRelease):
				release := strings.Split(_release, latestRelease)[1]
				_v.Version = *NewSemVersion(release)
				_v.getLatestRelease()
			case strings.Contains(_release, latestPatch):
				release := strings.Split(_release, latestPatch)[1]
				_v.Version = *NewSemVersion(release)
				_v.getLatestMajorRelease()
			case strings.Contains(_release, versionLessOrEqual):
				release := strings.Split(_release, versionLessOrEqual)[1]
				_v.Version = *NewSemVersion(release)
			case strings.Contains(_release, versionLessThan):
				release := strings.Split(_release, versionLessThan)[1]
				_v.Version = *NewSemVersion(release)
				_v.getOneLessRelease()
			}

			if isVersionGreater(*_v, *v) || iteration == 0 {
				v = _v
			}
		}

		return v
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
	case strings.Contains(v.Version.ToString(), versionLessOrEqual):
		release := strings.Split(v.Version.ToString(), versionLessOrEqual)[1]
		v.Version = *NewSemVersion(release)
	case strings.Contains(v.Version.ToString(), versionLessThan):
		release := strings.Split(v.Version.ToString(), versionLessThan)[1]
		v.Version = *NewSemVersion(release)
		v.getOneLessRelease()
	}

	return v
}

// GetVersionList returns a list of available versions from hashicorp's release page
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

// removeSpacesVersion removes spaces from Version string for parsing
func removeSpacesVersion(v string) string {
	splitV := strings.Split(v, " ")
	var returnString string
	for _, version := range splitV {
		version := strings.TrimSpace(version)
		returnString += version
	}
	return strings.TrimSpace(returnString)
}

// VersionToString returns string of a Version
func (v *Version) VersionToString() string {
	return v.Version.ToString()
}

// isVersionGreater returns true if v1 is greater than v2
func isVersionGreater(v1 Version, v2 Version) bool {
	if v1.Version.majorVersion != v2.Version.majorVersion {
		if v1.Version.majorVersion > v2.Version.majorVersion {
			return true
		}
		return false
	}

	if v1.Version.minorVersion != v2.Version.minorVersion {
		if v1.Version.minorVersion > v2.Version.minorVersion {
			return true
		}
		return false
	}

	if v1.Version.patchVersion != v2.Version.patchVersion {
		if v1.Version.patchVersion > v2.Version.patchVersion {
			return true
		}
		return false
	}

	return false
}
