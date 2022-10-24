package versionedTerraform

import "testing"

func TestSemVersion_VersionInSlice_success(t *testing.T) {
	want := SemVersion{
		version:      "1.1.9",
		isStable:     true,
		majorVersion: 1,
		minorVersion: 1,
		patchVersion: 9,
	}
	ver1 := SemVersion{
		version:      "1.1.10",
		isStable:     true,
		majorVersion: 1,
		minorVersion: 1,
		patchVersion: 10,
	}

	ver2 := SemVersion{
		version:      "0.1.9",
		isStable:     true,
		majorVersion: 0,
		minorVersion: 1,
		patchVersion: 9,
	}

	ver3 := SemVersion{
		version:      "1.0.9",
		isStable:     true,
		majorVersion: 1,
		minorVersion: 0,
		patchVersion: 9,
	}

	var semArray []SemVersion
	semArray = append(semArray, ver1)
	semArray = append(semArray, ver2)
	semArray = append(semArray, ver3)
	semArray = append(semArray, want)

	if !want.VersionInSlice(semArray) {
		t.Errorf("Expected Sem Version to be found in semArray")
	}
}

func TestSemVersion_VersionInSlice_fail(t *testing.T) {
	want := SemVersion{
		version:      "1.1.9",
		isStable:     true,
		majorVersion: 1,
		minorVersion: 1,
		patchVersion: 9,
	}
	ver1 := SemVersion{
		version:      "1.1.10",
		isStable:     true,
		majorVersion: 1,
		minorVersion: 1,
		patchVersion: 10,
	}

	ver2 := SemVersion{
		version:      "0.1.9",
		isStable:     true,
		majorVersion: 0,
		minorVersion: 1,
		patchVersion: 9,
	}

	ver3 := SemVersion{
		version:      "1.0.9",
		isStable:     true,
		majorVersion: 1,
		minorVersion: 0,
		patchVersion: 9,
	}

	var semArray []SemVersion
	semArray = append(semArray, ver1)
	semArray = append(semArray, ver2)
	semArray = append(semArray, ver3)

	if want.VersionInSlice(semArray) {
		t.Errorf("Expected Sem Version to not be found in semArray")
	}
}
