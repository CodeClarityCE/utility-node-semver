package versions

import (
	"fmt"
	"testing"
)

type VersionParsingTest struct {
	VersionString  string
	ExpectedResult Semver
}

type ComparisonToTest[T comparable] struct {
	v1             Semver
	v2             Semver
	ExpectedResult T
}

type PreReleaseComparisonToTest[T comparable] struct {
	v1                 Semver
	v2                 Semver
	ExpectedResult     T
	IncludePreReleases bool
}

func TestLEComperator(t *testing.T) {

	fmt.Printf("\n%s Testing LE (<=) comperator %s\n", "----------------", "----------------")

	comparisonsToTest := []ComparisonToTest[bool]{
		// 3.0.0 <= 3.0.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.0 <= 3.1.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.1 <= 3.1.1    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: true,
		},
		// 3.1.1 <= 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 3.1.0 <= 3.1.1    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: true,
		},
		// 3.1.0 <= 3.0.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.0.0 <= 3.1.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: true,
		},
		// 4.0.0 <= 3.0.0    -> false
		{
			v1:             Semver{Major: 4, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.0.0 <= 4.0.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
	}

	testComparisonsBool(t, comparisonsToTest, "<=")

}

func TestLEComperatorPreRelease(t *testing.T) {

	fmt.Printf("\n%s Testing LE (<=) comperator with prerelease versions %s\n", "----------------", "----------------")

	tests := []PreReleaseComparisonToTest[bool]{
		// 3.0.0 <= 3.0.0-beta.2    -> false
		// In a typical release cycle the full version '3.0.0' always comes after the beta version
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.2 <= 3.0.0-beta.2    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.1 <= 3.0.0-beta.2    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.1"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.3 <= 3.0.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.3 <= 3.1.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 1.0.0-alpha < 1.0.0-alpha.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			ExpectedResult: true,
		},
		// 1.0.0-alpha.1 < 1.0.0-alpha.beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			ExpectedResult: true,
		},
		// 1.0.0-alpha.beta < 1.0.0-beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			ExpectedResult: true,
		},
		// 1.0.0-beta < 1.0.0-beta.2
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.2 < 1.0.0-beta.11
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.11 < 1.0.0-rc.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			ExpectedResult: true,
		},
		// 1.0.0-rc.1 < 1.0.0
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: ""},
			ExpectedResult: true,
		},
	}

	testPreReleaseComparisons(t, tests, "<=")

}

func TestLTComperator(t *testing.T) {

	fmt.Printf("\n%s Testing LT (<) comperator %s\n", "----------------", "----------------")

	comparisonsToTest := []ComparisonToTest[bool]{
		// 3.0.0 < 3.0.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.1.0 < 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 3.1.1 < 3.1.1    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: false,
		},
		// 3.1.1 < 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 3.1.0 < 3.1.1    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: true,
		},
		// 3.1.0 < 3.0.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.0.0 < 3.1.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: true,
		},
		// 4.0.0 < 3.0.0    -> false
		{
			v1:             Semver{Major: 4, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.0.0 < 4.0.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
	}

	testComparisonsBool(t, comparisonsToTest, "<")
}

func TestLTComperatorPreRelease(t *testing.T) {

	fmt.Printf("\n%s Testing LT (<) comperator with prerelease versions %s\n", "----------------", "----------------")

	tests := []PreReleaseComparisonToTest[bool]{
		// 3.0.0 < 3.0.0-beta.2    -> true
		// In a typical release cycle the full version '3.0.0' always comes after the beta version
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.2 < 3.0.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.1 < 3.0.0-beta.2    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.1"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.3 < 3.0.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.3 < 3.1.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},

		// 1.0.0-alpha < 1.0.0-alpha.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			ExpectedResult: true,
		},
		// 1.0.0-alpha.1 < 1.0.0-alpha.beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			ExpectedResult: true,
		},
		// 1.0.0-alpha.beta < 1.0.0-beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			ExpectedResult: true,
		},
		// 1.0.0-beta < 1.0.0-beta.2
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.2 < 1.0.0-beta.11
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.11 < 1.0.0-rc.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			ExpectedResult: true,
		},
		// 1.0.0-rc.1 < 1.0.0
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: ""},
			ExpectedResult: true,
		},
	}

	testPreReleaseComparisons(t, tests, "<")

}

func TestGTComperator(t *testing.T) {

	fmt.Printf("\n%s Testing GT (>) comperator %s\n", "----------------", "----------------")

	comparisonsToTest := []ComparisonToTest[bool]{
		// 3.0.0 > 3.0.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.1.0 > 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 3.1.1 > 3.1.1    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: false,
		},
		// 3.1.1 > 3.1.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.0 > 3.1.1    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: false,
		},
		// 3.1.0 > 3.0.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
		// 3.0.0 > 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 4.0.0 > 3.0.0    -> true
		{
			v1:             Semver{Major: 4, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
		// 3.0.0 > 4.0.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
	}

	testComparisonsBool(t, comparisonsToTest, ">")
}

func TestGTComperatorPreRelease(t *testing.T) {

	fmt.Printf("\n%s Testing GT (>) comperator with prerelease versions %s\n", "----------------", "----------------")

	tests := []PreReleaseComparisonToTest[bool]{
		// 3.0.0 > 3.0.0-beta.2    -> true
		// In a typical release cycle the full version '3.0.0' always comes after the beta version
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.2 > 3.0.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.1 > 3.0.0-beta.2    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.1"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.3 > 3.0.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.3 > 3.1.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},

		// 1.0.0-alpha.1 > 1.0.0-alpha
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha"},
			ExpectedResult: true,
		},
		//  1.0.0-alpha.beta > 1.0.0-alpha.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			ExpectedResult: true,
		},
		// 1.0.0-beta > 1.0.0-alpha.beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.2 > 1.0.0-beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.11 > 1.0.0-beta.2
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 1.0.0-rc.1 > 1.0.0-beta.11
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			ExpectedResult: true,
		},
		// 1.0.0 > 1.0.0-rc.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: ""},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			ExpectedResult: true,
		},
	}

	testPreReleaseComparisons(t, tests, ">")

}

func TestGEComperator(t *testing.T) {

	fmt.Printf("\n%s Testing GE (>=) comperator %s\n", "----------------", "----------------")

	comparisonsToTest := []ComparisonToTest[bool]{
		// 3.0.0 >= 3.0.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.0 >= 3.1.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.1 >= 3.1.1    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: true,
		},
		// 3.1.1 >= 3.1.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.0 >= 3.1.1    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: false,
		},
		// 3.1.0 >= 3.0.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
		// 3.0.0 >= 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 4.0.0 >= 3.0.0    -> true
		{
			v1:             Semver{Major: 4, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
		// 3.0.0 >= 4.0.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
	}

	testComparisonsBool(t, comparisonsToTest, ">=")
}

func TestGEComperatorPreRelease(t *testing.T) {

	fmt.Printf("\n%s Testing GE (>=) comperator with prerelease versions %s\n", "----------------", "----------------")

	tests := []PreReleaseComparisonToTest[bool]{
		// 3.0.0 >= 3.0.0-beta.2    -> true
		// In a typical release cycle the full version '3.0.0' always comes after the beta version
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.2 >= 3.0.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.1 >= 3.0.0-beta.2    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.1"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		// 3.0.0-beta.3 >= 3.0.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 3.0.0-beta.3 >= 3.1.0-beta.2    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.3"},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},

		// 1.0.0-alpha.1 >= 1.0.0-alpha
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha"},
			ExpectedResult: true,
		},
		//  1.0.0-alpha.beta >= 1.0.0-alpha.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.1"},
			ExpectedResult: true,
		},
		// 1.0.0-beta >= 1.0.0-alpha.beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "alpha.beta"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.2 >= 1.0.0-beta
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			ExpectedResult: true,
		},
		// 1.0.0-beta.11 >= 1.0.0-beta.2
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		// 1.0.0-rc.1 >= 1.0.0-beta.11
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "beta.11"},
			ExpectedResult: true,
		},
		// 1.0.0 >= 1.0.0-rc.1
		{
			v1:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: ""},
			v2:             Semver{Major: 1, Minor: 0, Patch: 0, PreReleaseTag: "rc.1"},
			ExpectedResult: true,
		},
	}

	testPreReleaseComparisons(t, tests, ">=")

}

func TestEQComperator(t *testing.T) {

	fmt.Printf("\n%s Testing EQ (=) comperator %s\n", "----------------", "----------------")

	comparisonsToTest := []ComparisonToTest[bool]{
		// 3.0.0 = 3.0.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.0 = 3.1.0    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: true,
		},
		// 3.1.1 = 3.1.1    -> true
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: true,
		},
		// 3.1.1 = 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 3.1.0 = 3.1.1    -> false
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: false,
		},
		// 3.1.0 = 3.0.0    -> falsse
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.0.0 = 3.1.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: false,
		},
		// 4.0.0 = 3.0.0    -> true
		{
			v1:             Semver{Major: 4, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
		// 3.0.0 = 4.0.0    -> false
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult: false,
		},
	}

	testComparisonsBool(t, comparisonsToTest, "=")
}

func TestEQComperatorPreRelease(t *testing.T) {

	fmt.Printf("\n%s Testing EQ (=) comperator with prerelease versions %s\n", "----------------", "----------------")

	tests := []PreReleaseComparisonToTest[bool]{
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: true,
		},
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta"},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult: false,
		},
	}

	testPreReleaseComparisons(t, tests, "=")

}

func TestCompareOperator(t *testing.T) {

	fmt.Printf("\n%s Testing compare operator %s\n", "----------------", "----------------")

	comparisonsToTest := []ComparisonToTest[int]{
		// 3.0.0 = 3.0.0    -> 0
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: 0,
		},
		// 3.1.0 = 3.1.0    -> 0
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: 0,
		},
		// 3.1.1 = 3.1.1    -> 0
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: 0,
		},
		// 3.1.1 = 3.1.0    -> 1
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 1},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: 1,
		},
		// 3.1.0 = 3.1.1    -> -1
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult: -1,
		},
		// 3.1.0 = 3.0.0    -> 1
		{
			v1:             Semver{Major: 3, Minor: 1, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: 1,
		},
		// 3.0.0 = 3.1.0    -> -1
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 1, Patch: 0},
			ExpectedResult: -1,
		},
		// 4.0.0 = 3.0.0    -> 1
		{
			v1:             Semver{Major: 4, Minor: 0, Patch: 0},
			v2:             Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult: 1,
		},
		// 3.0.0 = 4.0.0    -> -1
		{
			v1:             Semver{Major: 3, Minor: 0, Patch: 0},
			v2:             Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult: -1,
		},
	}

	testComparisonsInt(t, comparisonsToTest, "=")
}

func TestVersionParsing(t *testing.T) {
	parseTests := []VersionParsingTest{

		// Normal versions
		{
			VersionString:  "0.0.0",
			ExpectedResult: Semver{Major: 0, Minor: 0, Patch: 0},
		},
		{
			VersionString:  "5.66.77",
			ExpectedResult: Semver{Major: 5, Minor: 66, Patch: 77},
		},
		{
			VersionString:  "5.66",
			ExpectedResult: Semver{}, // missing patch
		},

		// Invalid [major, minor, patch] version parts
		{
			VersionString:  "-5.0.0",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.-0.0",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.0.-0",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "+5.0.0",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.+0.0",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.0.+0",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "",
			ExpectedResult: Semver{},
		},
		{
			VersionString:  "ANY",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "*",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "*.*",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "*.*.*",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.66.x",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.*.*",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.*",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.x.x",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.x",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "x",
			ExpectedResult: Semver{}, // invalid character in version
		},
		{
			VersionString:  "5.0.0.0",
			ExpectedResult: Semver{}, // invalid version
		},

		// Build meta data
		{
			VersionString:  "1.0.0+20130313144700",
			ExpectedResult: Semver{Major: 1, Minor: 0, Patch: 0, MetaData: "20130313144700"},
		},
		{
			VersionString:  "1.0.0+exp.sha.5114f85",
			ExpectedResult: Semver{Major: 1, Minor: 0, Patch: 0, MetaData: "exp.sha.5114f85"},
		},
		{
			VersionString:  "1.0.0+exp.sha.5114f85()",
			ExpectedResult: Semver{}, // invalid character in pre release tag
		},
		{
			VersionString:  "1.0.0+21AF26D3+++117B344092BD",
			ExpectedResult: Semver{}, // invalid character in meta data section
		},

		// Pre release tag
		{
			VersionString:  "1.2.3-alpha.1",
			ExpectedResult: Semver{Major: 1, Minor: 2, Patch: 3, PreReleaseTag: "alpha.1"},
		},
		{
			VersionString:  "1.2.3-alpha-1",
			ExpectedResult: Semver{Major: 1, Minor: 2, Patch: 3, PreReleaseTag: "alpha-1"},
		},
		{
			VersionString:  "1.2.3-alpha-1!",
			ExpectedResult: Semver{}, // invalid character in pre release tag
		},
		{
			VersionString:  "1.2.3-alpha--------------2",
			ExpectedResult: Semver{Major: 1, Minor: 2, Patch: 3, PreReleaseTag: "alpha--------------2"},
		},
		{
			VersionString:  "4.18.0-next.1599210529.065f6b95a8f50e8a384176f775a1853a2cd341cf",
			ExpectedResult: Semver{Major: 4, Minor: 18, Patch: 0, PreReleaseTag: "next.1599210529.065f6b95a8f50e8a384176f775a1853a2cd341cf"},
		},
		{
			VersionString:  "4.18.0-next.1599210529.00065",
			ExpectedResult: Semver{}, // the spec disallows numeric identifiers to have leading zeros, i.e. the 00065 identifier
		},

		// Both build meta data and pre release version
		{
			VersionString:  "1.0.0-alpha+001",
			ExpectedResult: Semver{Major: 1, Minor: 0, Patch: 0, MetaData: "001", PreReleaseTag: "alpha"},
		},
		{
			VersionString:  "1.0.0+21AF26D3----117B344092BD",
			ExpectedResult: Semver{Major: 1, Minor: 0, Patch: 0, MetaData: "exp.sha.5114f85", PreReleaseTag: "---117B344092BD"},
		},
	}

	testVersionParsing(t, parseTests)

}

func testPreReleaseComparisons(t *testing.T, comparisons []PreReleaseComparisonToTest[bool], comperatorOperator string) {

	for _, comparisonToTest := range comparisons {
		fmt.Printf("\nTesting comparison: %s %s %s\n", comparisonToTest.v1.String(), comperatorOperator, comparisonToTest.v2.String())
		res := false
		if comperatorOperator == "<=" {
			res = comparisonToTest.v1.LE(comparisonToTest.v2, comparisonToTest.IncludePreReleases)
		} else if comperatorOperator == "<" {
			res = comparisonToTest.v1.LT(comparisonToTest.v2, comparisonToTest.IncludePreReleases)
		} else if comperatorOperator == ">" {
			res = comparisonToTest.v1.GT(comparisonToTest.v2, comparisonToTest.IncludePreReleases)
		} else if comperatorOperator == ">=" {
			res = comparisonToTest.v1.GE(comparisonToTest.v2, comparisonToTest.IncludePreReleases)
		} else if comperatorOperator == "=" {
			res = comparisonToTest.v1.EQ(comparisonToTest.v2, comparisonToTest.IncludePreReleases)
		}

		if res != comparisonToTest.ExpectedResult {
			fmt.Printf("✗ Failed. Expected: %t, but got: %t\n", comparisonToTest.ExpectedResult, res)
			t.Errorf("✗ Failed. Expected: %t, but got: %t\n", comparisonToTest.ExpectedResult, res)
		} else {
			fmt.Println("✓ Success")
		}
	}

}

func testComparisonsBool(t *testing.T, comparisons []ComparisonToTest[bool], comperatorOperator string) {
	for _, comparisonToTest := range comparisons {
		fmt.Printf("\nTesting comparison: %s %s %s\n", comparisonToTest.v1.String(), comperatorOperator, comparisonToTest.v2.String())
		res := false
		if comperatorOperator == "<=" {
			res = comparisonToTest.v1.LE(comparisonToTest.v2, false)
		} else if comperatorOperator == "<" {
			res = comparisonToTest.v1.LT(comparisonToTest.v2, false)
		} else if comperatorOperator == ">" {
			res = comparisonToTest.v1.GT(comparisonToTest.v2, false)
		} else if comperatorOperator == ">=" {
			res = comparisonToTest.v1.GE(comparisonToTest.v2, false)
		} else if comperatorOperator == "=" {
			res = comparisonToTest.v1.EQ(comparisonToTest.v2, false)
		}

		if res != comparisonToTest.ExpectedResult {
			fmt.Printf("✗ Failed. Expected: %t, but got: %t\n", comparisonToTest.ExpectedResult, res)
			t.Errorf("✗ Failed. Expected: %t, but got: %t\n", comparisonToTest.ExpectedResult, res)
		} else {
			fmt.Println("✓ Success")
		}
	}
}

func testComparisonsInt(t *testing.T, comparisons []ComparisonToTest[int], comperatorOperator string) {
	for _, comparisonToTest := range comparisons {
		fmt.Printf("\nTesting comparison: %s %s %s\n", comparisonToTest.v1.String(), comperatorOperator, comparisonToTest.v2.String())
		res := -2
		if comperatorOperator == "=" {
			res = comparisonToTest.v1.Compare(comparisonToTest.v2, false)
		}

		if res != comparisonToTest.ExpectedResult {
			fmt.Printf("✗ Failed. Expected: %d, but got: %d\n", comparisonToTest.ExpectedResult, res)
			t.Errorf("✗ Failed. Expected: %d, but got: %d\n", comparisonToTest.ExpectedResult, res)
		} else {
			fmt.Println("✓ Success")
		}
	}
}

func testVersionParsing(t *testing.T, tests []VersionParsingTest) {
	for _, test := range tests {
		fmt.Printf("\nTesting parsing of %s\n", test.VersionString)
		version, _ := ParseSemver(test.VersionString)

		if !version.EQ(test.ExpectedResult, true) {
			fmt.Printf("✗ Failed. Expected: %s, but got: %s\n", test.ExpectedResult.String(), version.String())
			t.Errorf("✗ Failed. Expected: %s, but got: %s\n", test.ExpectedResult.String(), version.String())
		} else {
			fmt.Println("✓ Success")
		}

	}
}
