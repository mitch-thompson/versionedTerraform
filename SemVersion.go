package versionedTerraform

import (
	"strconv"
	"strings"
)

const (
	//todo include others if needed
	//todo add comparison i.e. >= 0.11.10, < 0.12.0
	latestRelease      = ">="
	latestPatch        = "~>"
	versionLessOrEqual = "<="
	versionLessThan    = "<"
)

type SemVersion struct {
	version      string
	isStable     bool
	majorVersion int
	minorVersion int
	patchVersion int
}

type SemVersionInterface interface {
	setMajorVersion()
	setMinorVersion()
	setPatchVersion()
}

func NewSemVersion(v string) *SemVersion {
	s := new(SemVersion)
	s.isStable = true
	s.version = removeSpacesVersion(v)

	s.setMajorVersion()
	s.setMinorVersion()
	s.setPatchVersion()
	return s
}

//setMajorVersion setter for SemVersion.majorVersion
func (s *SemVersion) setMajorVersion() {
	version := s.version
	majorVersionString := strings.Split(version, ".")[0]
	s.majorVersion, _ = strconv.Atoi(majorVersionString)
}

//setMinorVersion setter for SemVersion.minorVersion
func (s *SemVersion) setMinorVersion() {
	version := s.version
	minorVersionString := strings.Split(version, ".")[1]
	s.minorVersion, _ = strconv.Atoi(minorVersionString)

}

//setPatchVersion setter for SemVersion.patchVersion
func (s *SemVersion) setPatchVersion() {
	version := s.version
	var err error
	patchStringSlice := strings.Split(version, ".")
	if len(patchStringSlice) < 3 {
		s.patchVersion = 0
		return
	}
	s.patchVersion, err = strconv.Atoi(patchStringSlice[2])
	if err != nil {
		s.isStable = false
		patchStringSlice = strings.Split(patchStringSlice[2], "-")
		s.patchVersion, _ = strconv.Atoi(patchStringSlice[0])
	}
}

//ToString returns string of SemVersion
func (s *SemVersion) ToString() string {
	return s.version
}

//VersionInSlice iterates through slices of SemVersion to check if version is in slice
//Used by main.go to determine if terraform version is currently installed
func (s *SemVersion) VersionInSlice(sSem []SemVersion) bool {
	for _, ver := range sSem {
		if ver.ToString() == s.ToString() {
			return true
		}
	}
	return false
}
