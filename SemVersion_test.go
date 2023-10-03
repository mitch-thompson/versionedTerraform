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

func TestSemVersion_IsLessThan(t *testing.T) {
	cases := []struct {
		name   string
		want   bool
		s1, s2 SemVersion
	}{
		{"IsEqualTo",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionLess",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MinorVersionLess",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
		},
		{"PatchVersionLess",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionGreater",
			false,
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
		},
		{"MinorVersionGreater",
			false,
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"PatchVersionGreater",
			false,
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
	}
	for _, c := range cases {
		t.Run("Test: "+c.name, func(t *testing.T) {
			t.Parallel()
			got := c.s1.IsLessThan(c.s2)
			if c.want != got {
				t.Errorf("Expected %+v got %+v", c.want, got)
			}
		})
	}
}

func TestSemVersion_IsGreaterThan(t *testing.T) {
	cases := []struct {
		name   string
		want   bool
		s1, s2 SemVersion
	}{
		{"IsEqualTo",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MinorVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
		},
		{"PatchVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionGreater",
			true,
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
		},
		{"MinorVersionGreater",
			true,
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"PatchVersionGreater",
			true,
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
	}
	for _, c := range cases {
		t.Run("Test: "+c.name, func(t *testing.T) {
			t.Parallel()
			got := c.s1.IsGreaterThan(c.s2)
			if c.want != got {
				t.Errorf("Expected %+v got %+v", c.want, got)
			}
		})
	}
}

func TestSemVersion_IsEqualTo(t *testing.T) {
	cases := []struct {
		name   string
		want   bool
		s1, s2 SemVersion
	}{
		{"IsEqualTo",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MinorVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
		},
		{"PatchVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionGreater",
			false,
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
		},
		{"MinorVersionGreater",
			false,
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"PatchVersionGreater",
			false,
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
	}
	for _, c := range cases {
		t.Run("Test: "+c.name, func(t *testing.T) {
			t.Parallel()
			got := c.s1.IsEqualTo(c.s2)
			if c.want != got {
				t.Errorf("Expected %+v got %+v", c.want, got)
			}
		})
	}
}

func TestSemVersion_IsLessOrEqual(t *testing.T) {
	cases := []struct {
		name   string
		want   bool
		s1, s2 SemVersion
	}{
		{"IsEqualTo",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionLess",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MinorVersionLess",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
		},
		{"PatchVersionLess",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionGreater",
			false,
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
		},
		{"MinorVersionGreater",
			false,
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"PatchVersionGreater",
			false,
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
	}
	for _, c := range cases {
		t.Run("Test: "+c.name, func(t *testing.T) {
			t.Parallel()
			got := c.s1.IsLessOrEqual(c.s2)
			if c.want != got {
				t.Errorf("Expected %+v got %+v", c.want, got)
			}
		})
	}
}

func TestSemVersion_IsGreaterOrEqual(t *testing.T) {
	cases := []struct {
		name   string
		want   bool
		s1, s2 SemVersion
	}{
		{"IsEqualTo",
			true,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MinorVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
		},
		{"PatchVersionLess",
			false,
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"MajorVersionGreater",
			true,
			SemVersion{
				version:      "2.2.3",
				isStable:     false,
				majorVersion: 2,
				minorVersion: 2,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
		},
		{"MinorVersionGreater",
			true,
			SemVersion{
				version:      "1.3.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 3,
				patchVersion: 3,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
		{"PatchVersionGreater",
			true,
			SemVersion{
				version:      "1.2.4",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 4,
			},
			SemVersion{
				version:      "1.2.3",
				isStable:     false,
				majorVersion: 1,
				minorVersion: 2,
				patchVersion: 3,
			},
		},
	}
	for _, c := range cases {
		t.Run("Test: "+c.name, func(t *testing.T) {
			t.Parallel()
			got := c.s1.IsGreaterOrEqual(c.s2)
			if c.want != got {
				t.Errorf("Expected %+v got %+v", c.want, got)
			}
		})
	}
}
