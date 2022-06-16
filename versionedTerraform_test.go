package versionedTerraform

import (
	"testing"
)

func testVersionList() []string {
	return []string{
		"1.2.23-alpha",
		"1.1.11",
		"1.1.10",
		"1.1.9",
		"1.1.8",
		"1.1.7",
		"1.1.6",
		"1.1.5",
		"1.1.4",
		"1.1.3",
		"1.1.2",
		"1.1.1",
		"1.0.12",
		"1.0.1",
		"0.14.0",
		"0.13.1",
		"0.13.0",
		"0.12.31",
		"0.12.30",
		"0.11.15",
		"0.11.10",
	}
}

func TestGetVersion(t *testing.T) {
	cases := []struct {
		available         []string
		version, expected string
	}{
		{testVersionList(), "0.12.31", "0.12.31"},
		{testVersionList(), "0.12.30", "0.12.30"},
		{testVersionList(), "~> 0.12.30", "0.12.31"},
		{testVersionList(), "~>0.12.30", "0.12.31"},
		{testVersionList(), "~>0.12.4", "0.12.31"},
		{testVersionList(), ">= 0.11.15", "1.1.11"},
		{testVersionList(), ">= 0.12.0", "1.1.11"},
		{testVersionList(), "~> 0.12", "0.12.31"},
		{testVersionList(), "< 0.12", "0.11.15"},
		{testVersionList(), "<= 0.12.31", "0.12.31"},
		{testVersionList(), "~> 0.12.0, < 0.13", "0.12.31"},
		{testVersionList(), "~> 0.12.0, < 0.14", "0.13.1"},
		{testVersionList(), "~> 0.12.0, <= 0.14.0", "0.14.0"},
	}

	for _, c := range cases {
		t.Run("test Version check with various conditions: "+c.version, func(t *testing.T) {
			//t.Parallel()
			got := NewVersion(c.version, c.available)
			if got.Version.version != c.expected {
				t.Errorf("got %q, want %q", got.Version.version, c.expected)
			}
		})
	}
}

func TestRemoveSpacesVersion(t *testing.T) {
	cases := []struct {
		tesValue, want string
	}{
		{"test", "test"},
		{"test ", "test"},
		{" test", "test"},
		{" test ", "test"},
		{" test test ", "testtest"},
	}

	for _, c := range cases {
		t.Run("test remove space in various conditions: "+c.tesValue, func(t *testing.T) {
			t.Parallel()
			got := removeSpacesVersion(c.tesValue)
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestIsVersionGreater(t *testing.T) {
	cases := []struct {
		testName                   string
		testValueOne, testValueTwo Version
		want                       bool
	}{
		{"equal versions",
			*NewVersion("0.12.10", testVersionList()),
			*NewVersion("0.12.10", testVersionList()),
			false,
		},
		{"major greater versions",
			*NewVersion("1.12.10", testVersionList()),
			*NewVersion("0.12.10", testVersionList()),
			true,
		},
		{"major less versions",
			*NewVersion("0.12.10", testVersionList()),
			*NewVersion("1.12.10", testVersionList()),
			false,
		},
		{"minor greater versions",
			*NewVersion("0.13.10", testVersionList()),
			*NewVersion("0.12.10", testVersionList()),
			true,
		},
		{"minor less  versions",
			*NewVersion("0.12.10", testVersionList()),
			*NewVersion("0.13.10", testVersionList()),
			false,
		},
		{"patch greater versions",
			*NewVersion("0.12.11", testVersionList()),
			*NewVersion("0.12.10", testVersionList()),
			true,
		},
		{"patch less versions",
			*NewVersion("0.12.10", testVersionList()),
			*NewVersion("0.12.11", testVersionList()),
			false,
		},
	}

	for _, c := range cases {
		t.Run("test isVersionGreater: "+c.testName, func(t *testing.T) {
			got := isVersionGreater(c.testValueOne, c.testValueTwo)
			if got != c.want {
				t.Errorf("Got %t, want %t", got, c.want)
			}
		})
	}
}

func TestGetVersionList(t *testing.T) {
	//todo write test for this
	//response, _ := getVersionList()
	//for _, Version := range response {
	//	t.Errorf("%v", Version)
	//}
	//t.Errorf("%v", response)
}

func TestInstallTerraformVersion(t *testing.T) {
	//todo write test for this
	//Version := NewVersion("0.12.31", testVersionList())
	//response := Version.InstallTerraformVersion()
	//t.Errorf("%v", response)
}
