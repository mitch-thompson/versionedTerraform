package versionedTerraform

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
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
		"terraform_0.12.31":      {Data: []byte("")},
		"terraform_0.12.30":      {Data: []byte("")},
		"terraform_0.11.10":      {Data: []byte("")},
		"terraform_0.11.15":      {Data: []byte("")},
		"terraform_1.0.1":        {Data: []byte("")},
		"terraform_1.0.12":       {Data: []byte("")},
		"terraform_1.2.23-alpha": {Data: []byte("")},
		"terraform_1.1.1":        {Data: []byte("")},
		"terraform_1.1.2":        {Data: []byte("")},
		"terraform_1.1.3":        {Data: []byte("")},
		"terraform_1.1.4":        {Data: []byte("")},
		"terraform_1.1.5":        {Data: []byte("")},
		"terraform_1.1.6":        {Data: []byte("")},
		"terraform_1.1.7":        {Data: []byte("")},
		"terraform_1.1.8":        {Data: []byte("")},
		"terraform_1.1.9":        {Data: []byte("")},
		"terraform_1.1.10":       {Data: []byte("")},
		"terraform_1.1.11":       {Data: []byte("")},
		"terraform_0.14.0":       {Data: []byte("")},
		"terraform_0.13.1":       {Data: []byte("")},
		"terraform_0.13.0":       {Data: []byte("")},
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

func TestConfigRequiresStable(t *testing.T) {
	availableVersions, _ := GetVersionList()
	versions := strings.Join(availableVersions, " ")
	cases := []struct {
		name, content, want string
		timeNow             time.Time
	}{
		{"StableOnly True", "StableOnly: true\n" +
			"LastUpdate: 1674481203\n" +
			"AvailableVersions: [1.3.7]",
			"StableOnly: true\n" +
				"LastUpdate: 1286705410\n" +
				"AvailableVersions: [" +
				versions + "]\n",
			time.Date(2010, 10, 10, 10, 10, 10, 10, time.UTC)},
		{"StableOnly False", "StableOnly: false\n" +
			"LastUpdate: 1674481203\n" +
			"AvailableVersions: [1.3.7]",
			"StableOnly: false\n" +
				"LastUpdate: 1286705410\n" +
				"AvailableVersions: [" +
				versions + "]\n",
			time.Date(2010, 10, 10, 10, 10, 10, 10, time.UTC)},
		{"StableOnly not found", "LastUpdate: 1674481203\n" +
			"AvailableVersions: [1.3.7]",
			"StableOnly: true\n" +
				"LastUpdate: 1286705410\n" +
				"AvailableVersions: [" +
				versions + "]\n",
			time.Date(2010, 10, 10, 10, 10, 10, 10, time.UTC)},
	}

	for _, c := range cases {
		t.Run("Test: "+c.name, func(t *testing.T) {
			t.Parallel()

			tempDir := os.TempDir()
			tempFile, err := os.Create(tempDir + "/config")
			defer tempFile.Close()

			if err != nil {
				t.Errorf("Unable to execute test : %v", err)
			}

			UpdateConfig(*tempFile, c.timeNow)

			tempFile.Seek(0, 0)
			data := make([]byte, 1024)
			var got string
			for {
				n, err := tempFile.Read(data)
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Errorf("File reading error : %v", err)
					return
				}
				got += string(data[:n])
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("%v test failed to meet conditions", c.name)
				fmt.Fprintf(os.Stdout, "%v\n", c.want)
				fmt.Fprintf(os.Stdout, "%v\n", got)
			}
		})
	}
}
