package versionedTerraform

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"testing/fstest"
	"time"
)

func TestUpdateAvailableVersions(t *testing.T) {
	timeNow := time.Now()
	currentTime := timeNow.Unix()
	twoDaysAgoTime := timeNow.AddDate(0, 0, -2).Unix()

	successUpdate := fmt.Sprintf("LastUpdate: %d", currentTime)
	needsUpdate := fmt.Sprintf("LastUpdate: %d", twoDaysAgoTime)

	fs := fstest.MapFS{
		"successConfig.conf": {Data: []byte(successUpdate)},
		"failConfig.conf":    {Data: []byte(needsUpdate)},
	}

	t.Run("Test success last update time", func(t *testing.T) {
		want := false
		got, err := NeedToUpdateAvailableVersions(fs, "successConfig.conf")
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Errorf("updateAvailableVersions had incorrect output expected %v got %v", want, got)
		}
	})

	t.Run("Test failed last update time", func(t *testing.T) {
		want := true
		got, err := NeedToUpdateAvailableVersions(fs, "failConfig.conf")
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Errorf("updateAvailableVersions had incorrect output expected %v got %v", want, got)
		}
	})
}

func TestAvailableVersions(t *testing.T) {
	availableVersionList := fmt.Sprintf("AvailableVersions: %+v", testVersionList())
	var want []SemVersion
	for _, version := range testVersionList() {
		want = append(want, *NewSemVersion(version))
	}

	fs := fstest.MapFS{
		"successConfig.conf": {Data: []byte(availableVersionList)},
	}

	t.Run("Test success last update time", func(t *testing.T) {
		got, err := LoadVersionsFromConfig(fs, "successConfig.conf")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("LoadInstalledVersions had incorrect output expected %+v\n got %+v", want, got)
		}
	})
}

func TestInstalledVersions(t *testing.T) {
	var want []SemVersion
	testVersionList := testVersionList()
	sort.Strings(testVersionList)
	for _, version := range testVersionList {
		want = append(want, *NewSemVersion(version))
	}

	fs := fstest.MapFS{
		"terraform_0.12.31": {Data: []byte("")},
		"terraform_0.12.30": {Data: []byte("")},
		"terraform_0.11.10": {Data: []byte("")},
		"terraform_0.11.15": {Data: []byte("")},
		"terraform_1.0.1":   {Data: []byte("")},
		"terraform_1.0.12":  {Data: []byte("")},
		"terraform_1.1.1":   {Data: []byte("")},
		"terraform_1.1.2":   {Data: []byte("")},
		"terraform_1.1.3":   {Data: []byte("")},
		"terraform_1.1.4":   {Data: []byte("")},
		"terraform_1.1.5":   {Data: []byte("")},
		"terraform_1.1.6":   {Data: []byte("")},
		"terraform_1.1.7":   {Data: []byte("")},
		"terraform_1.1.8":   {Data: []byte("")},
		"terraform_1.1.9":   {Data: []byte("")},
		"terraform_1.1.10":  {Data: []byte("")},
		"terraform_1.1.11":  {Data: []byte("")},
		"terraform_0.14.0":  {Data: []byte("")},
		"terraform_0.13.1":  {Data: []byte("")},
		"terraform_0.13.0":  {Data: []byte("")},
	}

	t.Run("Test installed versions", func(t *testing.T) {
		got, err := LoadInstalledVersions(fs)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("LoadInstalledVersions had incorrect output expected\n     %+v\n got %+v", want, got)
		}
	})

}
