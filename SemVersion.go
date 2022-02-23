package versionedTerraform

import (
	"strconv"
	"strings"
)

const (
	//todo include others if needed
	//todo add comparison i.e. >= 0.11.10, < 0.12.0
	latestRelease = ">="
	latestPatch   = "~>"
)

type SemVersion struct {
	version      string
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
	s.version = removeSpacesVersion(v)

	s.setMajorVersion()
	s.setMinorVersion()
	s.setPatchVersion()
	return s
}

func (s *SemVersion) setMajorVersion() {
	version := s.version
	majorVersionString := strings.Split(version, ".")[0]
	s.majorVersion, _ = strconv.Atoi(majorVersionString)
}

func (s *SemVersion) setMinorVersion() {
	version := s.version
	minorVersionString := strings.Split(version, ".")[1]
	s.minorVersion, _ = strconv.Atoi(minorVersionString)

}

func (s *SemVersion) setPatchVersion() {
	version := s.version
	patchStringSlice := strings.Split(version, ".")
	if len(patchStringSlice) < 3 {
		s.patchVersion = 0
		return
	}
	s.patchVersion, _ = strconv.Atoi(patchStringSlice[2])

}

func (s *SemVersion) ToString() string {
	return s.version
}

func (s *SemVersion) VersionInSlice(sSem []SemVersion) bool {
	for _, ver := range sSem {
		if ver.ToString() == s.ToString() {
			return true
		}
	}
	return false
}
