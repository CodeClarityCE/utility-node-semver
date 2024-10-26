package evaluator

import (
	"fmt"
	"testing"

	"github.com/CodeClarityCE/utility-node-semver/versions"

	constraintTypes "github.com/CodeClarityCE/utility-node-semver/constraints"
)

type MaxSatisfyíngConstraintToTest struct {
	ConstraintString             string
	Versions                     []string
	ExpectedMaxSatisfyingVersion versions.Semver
}

type ConstraintToTest struct {
	ConstraintString string
	Version          versions.Semver
	ExpectedResult   bool
}

type ConstraintToTestPreReleases struct {
	ConstraintString   string
	Version            versions.Semver
	ExpectedResult     bool
	IncludePreReleases bool
}

func validateConstraint(c ConstraintToTest) (bool, error) {

	parsedConstraint, err := constraintTypes.ParseConstraint(c.ConstraintString)

	if err != nil {
		return false, err
	}

	satisfies := Satisfies(c.Version, parsedConstraint, false)

	if satisfies != c.ExpectedResult {
		return false, nil
	}

	return true, nil

}

func validateConstraintPreReleases(c ConstraintToTestPreReleases) (bool, error) {

	parsedConstraint, err := constraintTypes.ParseConstraint(c.ConstraintString)

	if err != nil {
		return false, err
	}

	satisfies := Satisfies(c.Version, parsedConstraint, c.IncludePreReleases)

	if satisfies != c.ExpectedResult {
		return false, nil
	}

	return true, nil

}

func TestConjunctionConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing conjuncted constraint expressions (... && ...) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{

		// 5.x.x && < 5.2.5
		// This constraint makes sure that we are using a version with major 5, but below 5.2.5, since 5.2.5 is a new version that was release
		// which is vulnerable
		{
			ConstraintString: "5.x.x && < 5.2.5",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "5.x.x && < 5.2.5",
			Version:          versions.Semver{Major: 5, Minor: 1, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "5.x.x && < 5.2.5",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 5},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "5.x.x && < 5.2.5",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 5},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "5.x.x && < 5.2.5",
			Version:          versions.Semver{Major: 4, Minor: 2, Patch: 5},
			ExpectedResult:   false,
		},

		// The following constraints make no "real-world sense" but they just test if the seperate conditions are evaluated correctly
		// and then evaluated correctly as a whole; as they are joined by a &&

		// 5.2.x && 6.0.0
		{
			ConstraintString: "5.2.x && 6.0.0", // this condition is impossible to satisfy
			Version:          versions.Semver{Major: 5, Minor: 7, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "5.2.x && 6.0.0", // this condition is impossible to satisfy
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},

		// 4.x && < 7.0.0
		{
			ConstraintString: "4.x && < 7.0.0",
			Version:          versions.Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "4.x && < 7.0.0",
			Version:          versions.Semver{Major: 6, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "4.x && < 7.0.0",
			Version:          versions.Semver{Major: 3, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "4.x && < 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},

		// ~4.5 && ^4.0.0
		{
			ConstraintString: "~4.5 && ^4.0.0",
			Version:          versions.Semver{Major: 4, Minor: 5, Patch: 5},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~4.5 && ^4.0.0",
			Version:          versions.Semver{Major: 4, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestDisjunctionConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing disjuncted constraint expressions (... || ...) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{

		// >= 5.0.0 < 5.2.5 || > 5.2.5 <= 5.5.0
		// example this wants to make sure that 5.2.5 is not used since it is vulnerable and to keep compatability with the api we need to use a version betweeen 5.0.0 and 5.5.0
		// in other words we want a version between 5.0.0 and 5.2.5 (excluding 5.2.5) or something higher, while being less than 5.5.0
		{
			ConstraintString: ">= 5.0.0 < 5.2.5 || > 5.2.5 <= 5.5.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0 < 5.2.5 || > 5.2.5 <= 5.5.0",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 6},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0 < 5.2.5 || > 5.2.5 <= 5.5.0",
			Version:          versions.Semver{Major: 5, Minor: 5, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0 < 5.2.5 || > 5.2.5 <= 5.5.0",
			Version:          versions.Semver{Major: 6, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: ">= 5.0.0 < 5.2.5 || > 5.2.5 <= 5.5.0",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 5},
			ExpectedResult:   false,
		},

		// A developer of vulnerable library released patches for each major release
		// This is a real world example from the library called "pg" on npm for vuln CVE-2017-16082
		//
		// We use this constraint to check if our version is vulnerable
		//
		// Affected versions of this library are:
		//  < 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 0, Minor: 4, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 1, Minor: 4, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 2, Minor: 11, Patch: 2},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 3, Minor: 6, Patch: 4},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 3, Minor: 7, Patch: 4},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 4, Minor: 5, Patch: 7},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 4, Minor: 5, Patch: 8},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 1},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 2},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 0, Patch: 5},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 0, Patch: 6},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 1, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 1, Patch: 6},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 1, Patch: 7},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 2, Patch: 5},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 2, Patch: 6},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 3, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 3, Patch: 3},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 3, Patch: 4},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 4, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 4, Patch: 2},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 6, Minor: 4, Patch: 3},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 2},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 3},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 1, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 1, Patch: 2},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7 || >= 5.0.0 < 5.2.1 || >= 6.0.0 < 6.0.5 || >= 6.1.0 < 6.1.6 || >= 6.2.0 < 6.2.5 || >= 6.3.0 < 6.3.3 || >= 6.4.0 < 6.4.2 || >= 7.0.0 < 7.0.2 || >= 7.1.0 < 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 1, Patch: 3},
			ExpectedResult:   false,
		},

		// A developer of vulnerable library released patches for each major release
		// This is a real world example from the library called "pg" on npm for vuln CVE-2017-16082
		//
		// We use this constraint to check if our version is patched
		//
		// Patched versions of this library are:
		//  2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 2, Minor: 11, Patch: 1},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 2, Minor: 11, Patch: 2},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 2, Minor: 11, Patch: 3},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 3, Minor: 6, Patch: 3},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 3, Minor: 6, Patch: 4},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 3, Minor: 6, Patch: 5},
			ExpectedResult:   false,
		},

		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 1, Patch: 1},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 7, Minor: 1, Patch: 2},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || 5.2.1 || 6.0.5 || 6.1.6 || 6.2.5 || 6.3.3 || 6.4.2 || 7.0.2 || 7.0.3 || >= 7.1.2",
			Version:          versions.Semver{Major: 8, Minor: 6, Patch: 5},
			ExpectedResult:   true,
		},

		// The following constraints make no "real-world sense" but they just test if the seperate conditions are evaluated correctly

		// A static list of versions joined by a disjunction
		{
			ConstraintString: "5.2.6 || 5.2.7 || 5.2.8",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 6},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "5.2.6 || 5.2.7 || 5.2.8",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 7},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "5.2.6 || 5.2.7 || 5.2.8",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 8},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "5.2.6 || 5.2.7 || 5.2.8",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 5},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "5.2.6 || 5.2.7 || 5.2.8",
			Version:          versions.Semver{Major: 5, Minor: 2, Patch: 9},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestDisjunctionAndConjunctionConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing combined disjuncted and conjuncted constraint expressions (... || ... && ...) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{
		{
			ConstraintString: "4.0.0 || 5.x && < 6.0.0 || 7.0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "4.0.0 || 5.x && < 6.0.0 || 7.0.0",
			Version:          versions.Semver{Major: 3, Minor: 6, Patch: 5},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "4.0.0 || 5.x && < 6.0.0 || >= 7.0.0",
			Version:          versions.Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "4.0.0 || 5.x && < 6.0.0 || >= 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "4.0.0 || 5.x && < 6.0.0 || >= 7.0.0",
			Version:          versions.Semver{Major: 9, Minor: 5, Patch: 2},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "4.0.0 || 5.x && < 6.0.0 || >= 7.0.0",
			Version:          versions.Semver{Major: 5, Minor: 5, Patch: 2},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "4.0.0 || 5.x && < 6.0.0 || >= 7.0.0",
			Version:          versions.Semver{Major: 6, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestXRangeConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing X range (5.x) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{
		{
			ConstraintString: "ANY",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "ANY",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "*",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "*",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "*.*",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "*.*",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "*.*.*",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "*.*.*",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "3.x",
			Version:          versions.Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "3.x",
			Version:          versions.Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "3.x",
			Version:          versions.Semver{Major: 2, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 1, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "3",
			Version:          versions.Semver{Major: 3, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "3",
			Version:          versions.Semver{Major: 3, Minor: 1, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "3",
			Version:          versions.Semver{Major: 2, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2",
			Version:          versions.Semver{Major: 1, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2",
			Version:          versions.Semver{Major: 1, Minor: 1, Patch: 9},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestRangeConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing range constraint (>= x.x.x =< x.x.x) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{

		// <= 5.0.0
		{
			ConstraintString: "<= 5.0.0",
			Version:          versions.Semver{Major: 1, Minor: 1, Patch: 9},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "<= 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "<= 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "<= 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 1, Patch: 0},
			ExpectedResult:   false,
		},

		// >= 5.0.0
		{
			ConstraintString: ">= 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0",
			Version:          versions.Semver{Major: 4, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},

		// < 5.0.0
		{
			ConstraintString: "< 5.0.0",
			Version:          versions.Semver{Major: 1, Minor: 1, Patch: 9},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "< 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "< 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 1, Patch: 0},
			ExpectedResult:   false,
		},

		// > 5.0.0
		{
			ConstraintString: "> 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "> 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "> 5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 1, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "> 5.0.0",
			Version:          versions.Semver{Major: 4, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},

		// >= 5.0.0 <= 7.0.0
		{
			ConstraintString: ">= 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 4, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: ">= 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: ">= 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 1},
			ExpectedResult:   false,
		},

		// > 5.0.0 <= 7.0.0
		{
			ConstraintString: "> 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "> 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "> 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 4, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "> 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "> 5.0.0 <= 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 1},
			ExpectedResult:   false,
		},

		// > 5.0.0 < 7.0.0
		{
			ConstraintString: "> 5.0.0 < 7.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "> 5.0.0 < 7.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "> 5.0.0 < 7.0.0",
			Version:          versions.Semver{Major: 4, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "> 5.0.0 < 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "> 5.0.0 < 7.0.0",
			Version:          versions.Semver{Major: 7, Minor: 0, Patch: 1},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestTildeConstraint(t *testing.T) {
	fmt.Printf("\n%s Testing tilde constraint (~5.0.0) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{

		// ~1.2.3
		{
			ConstraintString: "~1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 3},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.2.3",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 3},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 2},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 9, Patch: 9},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},

		// ~1.2
		{
			ConstraintString: "~1.2",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.2",
			Version:          versions.Semver{Major: 1, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.2",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.2",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 3},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.2",
			Version:          versions.Semver{Major: 2, Minor: 5, Patch: 55},
			ExpectedResult:   false,
		},

		// ~1.2.x
		{
			ConstraintString: "~1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.2.x",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 3},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.2.x",
			Version:          versions.Semver{Major: 2, Minor: 5, Patch: 55},
			ExpectedResult:   false,
		},

		// ~1
		{
			ConstraintString: "~1",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1",
			Version:          versions.Semver{Major: 1, Minor: 9, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1",
			Version:          versions.Semver{Major: 0, Minor: 9, Patch: 99},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1",
			Version:          versions.Semver{Major: 3, Minor: 9, Patch: 99},
			ExpectedResult:   false,
		},

		// ~1.x
		{
			ConstraintString: "~1.x",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.x",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.x",
			Version:          versions.Semver{Major: 1, Minor: 9, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~1.x",
			Version:          versions.Semver{Major: 0, Minor: 9, Patch: 99},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~1.x",
			Version:          versions.Semver{Major: 3, Minor: 9, Patch: 99},
			ExpectedResult:   false,
		},

		// ~0.2.3
		{
			ConstraintString: "~0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 3},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 1, Patch: 99},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.2.3",
			Version:          versions.Semver{Major: 5, Minor: 1, Patch: 99},
			ExpectedResult:   false,
		},

		// ~0.2
		{
			ConstraintString: "~0.2",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.2",
			Version:          versions.Semver{Major: 0, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.2",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 55},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.2",
			Version:          versions.Semver{Major: 0, Minor: 1, Patch: 55},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.2",
			Version:          versions.Semver{Major: 5, Minor: 5, Patch: 55},
			ExpectedResult:   false,
		},

		// ~0.2.x
		{
			ConstraintString: "~0.2.x",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.2.x",
			Version:          versions.Semver{Major: 0, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.2.x",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 55},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.2.x",
			Version:          versions.Semver{Major: 0, Minor: 1, Patch: 55},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.2.x",
			Version:          versions.Semver{Major: 5, Minor: 5, Patch: 55},
			ExpectedResult:   false,
		},

		// ~0
		{
			ConstraintString: "~0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0",
			Version:          versions.Semver{Major: 0, Minor: 99, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0",
			Version:          versions.Semver{Major: 3, Minor: 99, Patch: 99},
			ExpectedResult:   false,
		},

		// ~0.x
		{
			ConstraintString: "~0.x",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.x",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.x",
			Version:          versions.Semver{Major: 0, Minor: 99, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.x",
			Version:          versions.Semver{Major: 3, Minor: 99, Patch: 99},
			ExpectedResult:   false,
		},

		// ~0.0.0
		{
			ConstraintString: "~0.0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.0.0",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.0.0",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 22},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.0.0",
			Version:          versions.Semver{Major: 2, Minor: 55, Patch: 55},
			ExpectedResult:   false,
		},

		// ~0.0
		{
			ConstraintString: "~0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "~0.0",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.0",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 22},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "~0.0",
			Version:          versions.Semver{Major: 2, Minor: 55, Patch: 55},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestCaretConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing caret constraint (^5.0.0) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{

		// ^1.2.3
		{
			ConstraintString: "^1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 3},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.2.3",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 3},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 2},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^1.2.3",
			Version:          versions.Semver{Major: 1, Minor: 9, Patch: 9},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.2.3",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0.2.3
		{
			ConstraintString: "^0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 3},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 2, Patch: 2},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^0.2.3",
			Version:          versions.Semver{Major: 0, Minor: 3, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0.0.3
		{
			ConstraintString: "^0.0.3",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 3},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0.3",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 4},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^0.0.3",
			Version:          versions.Semver{Major: 0, Minor: 1, Patch: 1},
			ExpectedResult:   false,
		},

		// ^1.2.x
		{
			ConstraintString: "^1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.2.x",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 5},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 6, Patch: 5},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.2.x",
			Version:          versions.Semver{Major: 1, Minor: 1, Patch: 9},
			ExpectedResult:   false,
		},

		// ^0.0.x
		{
			ConstraintString: "^0.0.x",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0.x",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0.x",
			Version:          versions.Semver{Major: 0, Minor: 1, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0.0
		{
			ConstraintString: "^0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0",
			Version:          versions.Semver{Major: 0, Minor: 1, Patch: 0},
			ExpectedResult:   false,
		},

		// ^1.x
		{
			ConstraintString: "^1.x",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.x",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 99},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^1.x",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.x",
			Version:          versions.Semver{Major: 1, Minor: 5, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^1.x",
			Version:          versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0.x
		{
			ConstraintString: "^0.x",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.x",
			Version:          versions.Semver{Major: 0, Minor: 99, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.x",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0.x.x
		{
			ConstraintString: "^0.x.x",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.x.x",
			Version:          versions.Semver{Major: 0, Minor: 99, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.x.x",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0
		{
			ConstraintString: "^0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0",
			Version:          versions.Semver{Major: 0, Minor: 99, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0.0
		{
			ConstraintString: "^0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0",
			Version:          versions.Semver{Major: 0, Minor: 99, Patch: 99},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0",
			Version:          versions.Semver{Major: 0, Minor: 1, Patch: 0},
			ExpectedResult:   false,
		},

		// ^0.0.0
		{
			ConstraintString: "^0.0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "^0.0.0",
			Version:          versions.Semver{Major: 0, Minor: 0, Patch: 1},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "^0.0.0",
			Version:          versions.Semver{Major: 1, Minor: 0, Patch: 0},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestHyphenConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing hypenated constraint (5.0.0 - 7.0.0) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{

		// 1.2.3 - 2.3.4
		{
			ConstraintString: "1.2.3 - 2.3.4",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 3},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2.3 - 2.3.4",
			Version:          versions.Semver{Major: 2, Minor: 3, Patch: 4},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2.3 - 2.3.4",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 2},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2.3 - 2.3.4",
			Version:          versions.Semver{Major: 2, Minor: 3, Patch: 5},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2.3 - 2.3.4",
			Version:          versions.Semver{Major: 1, Minor: 7, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2.3 - 2.3.4",
			Version:          versions.Semver{Major: 2, Minor: 2, Patch: 99},
			ExpectedResult:   true,
		},

		// 1.2 - 2.3.4
		{
			ConstraintString: "1.2 - 2.3.4",
			Version:          versions.Semver{Major: 1, Minor: 2, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2 - 2.3.4",
			Version:          versions.Semver{Major: 2, Minor: 3, Patch: 4},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2 - 2.3.4",
			Version:          versions.Semver{Major: 2, Minor: 3, Patch: 5},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2 - 2.3.4",
			Version:          versions.Semver{Major: 1, Minor: 7, Patch: 99},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "1.2 - 2.3.4",
			Version:          versions.Semver{Major: 1, Minor: 1, Patch: 99},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "1.2 - 2.3.4",
			Version:          versions.Semver{Major: 2, Minor: 2, Patch: 99},
			ExpectedResult:   true,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestEqualityConstraint(t *testing.T) {

	fmt.Printf("\n%s Testing equality constraint (=5.0.0) evaluation %s\n", "----------------", "----------------")

	constraintsToTest := []ConstraintToTest{
		{
			ConstraintString: "=5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:   true,
		},
		{
			ConstraintString: "=5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 1, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "=5.0.0",
			Version:          versions.Semver{Major: 5, Minor: 0, Patch: 1},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "=5.0.0",
			Version:          versions.Semver{Major: 4, Minor: 9, Patch: 0},
			ExpectedResult:   false,
		},
		{
			ConstraintString: "=5.0.0",
			Version:          versions.Semver{Major: 4, Minor: 0, Patch: 9},
			ExpectedResult:   false,
		},
	}

	testConstraints(t, constraintsToTest)

}

func TestMaxSatisfyingEvaluation(t *testing.T) {

	fmt.Printf("\n%s Testing max satifying evaluation %s\n", "----------------", "----------------")

	// Since we tested that each range is correctly evaluated, we dont need to test exhaustively each type of range here
	// We just add 1 test case for each type
	constraintsToTest := []MaxSatisfyíngConstraintToTest{

		// Range
		{
			ConstraintString:             ">= 5.0.0",
			Versions:                     []string{"0.0.0", "2.5.0", "4.9.9"},
			ExpectedMaxSatisfyingVersion: versions.Semver{},
		},
		{
			ConstraintString:             ">= 5.0.0",
			Versions:                     []string{"0.0.0", "2.5.0", "5.0.0"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
		},
		{
			ConstraintString:             ">= 5.0.0",
			Versions:                     []string{"0.0.0", "2.5.0", "5.0.0", "10.11.12"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 10, Minor: 11, Patch: 12},
		},

		// X Range
		{
			ConstraintString:             "5.x.x",
			Versions:                     []string{"0.0.0", "2.5.0", "4.9.9"},
			ExpectedMaxSatisfyingVersion: versions.Semver{},
		},
		{
			ConstraintString:             "5.x.x.",
			Versions:                     []string{"0.0.0", "2.5.0", "5.0.0"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
		},
		{
			ConstraintString:             "5.x.x",
			Versions:                     []string{"0.0.0", "2.5.0", "5.0.0", "5.10.12"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 5, Minor: 10, Patch: 12},
		},

		// Tilde (~) Range
		{
			ConstraintString:             "~5.6.7",
			Versions:                     []string{"0.0.0", "2.5.0", "4.9.9"},
			ExpectedMaxSatisfyingVersion: versions.Semver{},
		},
		{
			ConstraintString:             "~5.6.7",
			Versions:                     []string{"0.0.0", "2.5.0", "5.6.7"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 5, Minor: 6, Patch: 7},
		},
		{
			ConstraintString:             "~5.6.7",
			Versions:                     []string{"0.0.0", "2.5.0", "5.6.99", "5.7.19"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 5, Minor: 6, Patch: 99},
		},

		// Tilde (^) Range
		{
			ConstraintString:             "^5.6.7",
			Versions:                     []string{"0.0.0", "2.5.0", "4.9.9"},
			ExpectedMaxSatisfyingVersion: versions.Semver{},
		},
		{
			ConstraintString:             "^5.6.7",
			Versions:                     []string{"0.0.0", "2.5.0", "5.6.7"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 5, Minor: 6, Patch: 7},
		},
		{
			ConstraintString:             "^5.6.7",
			Versions:                     []string{"0.0.0", "2.5.0", "5.6.99", "5.7.19", "6.0.0"},
			ExpectedMaxSatisfyingVersion: versions.Semver{Major: 5, Minor: 7, Patch: 19},
		},
	}

	testMaxSatisfying(t, constraintsToTest)

}

func TestIncludePreReleases(t *testing.T) {

	// Since we tested that each range is correctly evaluated, we dont need to test exhaustively each type of range here
	// Any kind of operator lile ~,^,x,... desugares to a "normal" range anyhow
	// Additionally we also tested that && and || joined expressions are also evaluated correctly
	//
	// The goal here is just to validate that the includePreReleases works correctly when evaluated Constraints
	constraintsToTest := []ConstraintToTestPreReleases{

		// While not including pre releases
		// 5.0.0-beta.2 satisfies > 5.0.0-beta.2?
		//  no, because the prerelease beta.2 is not strictly greater than beta.2
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 5.0.0-beta.1 satisfies > 5.0.0-beta.2?
		//  no, because the prerelease beta.1 is not smaller than beta.2
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.1"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 5.0.0-beta.5 satisfies > 5.0.0-beta.2?
		//  yes, because the prerelease beta.5 is stictly greater than beta.2
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: false,
		},
		// 5.0.0 satisfies > 5.0.0-beta.2?
		//  yes, because in a typical release cycle 5.0.0-beta.2 would always come before the full release 5.0.0
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:     true,
			IncludePreReleases: false,
		},
		// 2.0.0 satisfies > 5.0.0-beta.2?
		//  no, because 2.0.0 is strictly smaller than 5.0.0
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 2, Minor: 0, Patch: 0},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},

		// While including pre releases
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 9, Patch: 9, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   "> 5.0.0-beta.2",
			Version:            versions.Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult:     false,
			IncludePreReleases: true,
		},

		// While not including pre releases
		// 5.0.0-beta.2 satisfies < 5.0.0-beta.2?
		//  no, because the prerelease beta.2 is not strictly less than beta.2
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 5.0.0-beta.1 satisfies < 5.0.0-beta.2?
		//  yes, because the prerelease beta.1 is less than beta.2
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.1"},
			ExpectedResult:     true,
			IncludePreReleases: false,
		},
		// 5.0.0-beta.5 satisfies < 5.0.0-beta.2?
		//  no, because the prerelease beta.5 is stictly greater than beta.2
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 5.0.0 satisfies < 5.0.0-beta.2?
		//  no, because in a typical release cycle 5.0.0-beta.2 would always come before the full release 5.0.0
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 4.0.0 satisfies < 5.0.0-beta.2?
		//  yes, because 4.0.0 is strictly smaller thanb 5.0.0
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 4, Minor: 0, Patch: 0},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},

		// While including pre releases
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 5, Minor: 9, Patch: 9, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 2, Minor: 2, Patch: 5, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   "< 5.0.0-beta.2",
			Version:            versions.Semver{Major: 8, Minor: 0, Patch: 0},
			ExpectedResult:     false,
			IncludePreReleases: true,
		},

		// 5.0.0-beta.2 satisfies >= 5.0.0-beta.2 <= 6.0.0?
		//  yes, because the prerelease beta.2 is equal to beta.2 (of the start operator)
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult:     true,
			IncludePreReleases: false,
		},
		// 5.0.0-beta.5 satisfies >= 5.0.0-beta.2 <= 6.0.0?
		//  yes, because the prerelease beta.5 is greater than beta.2 (of the start operator)
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: false,
		},
		// 5.0.0 satisfies >= 5.0.0-beta.2 <= 6.0.0?
		//  yes, because in a typical release cycle 5.0.0-beta.2 would always come before the full release 5.0.0
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:     true,
			IncludePreReleases: false,
		},
		// 5.0.5-beta.5 satisfies >= 5.0.0-beta.2 <= 6.0.0?
		//  no, as per semver spec pre releases can only match if the [major, minor, patch] is equal to atleast one operator
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 5, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 5.2.5-beta.5 satisfies >= 5.0.0-beta.2 <= 6.0.0?
		//  no, as per semver spec pre releases can only match if the [major, minor, patch] is equal to atleast one operator
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 2, Patch: 5, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 5.0.0-beta.5 satisfies >= 5.0.0 <= 6.0.0?
		//  no, as per semver spec pre releases can only match if the [major, minor, patch] is equal to atleast one operator
		{
			ConstraintString:   ">= 5.0.0 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},
		// 6.0.0-beta.2 satisfies >= 5.0.0 <= 6.0.0-beta.2?
		//  yes, beta.2 is equal to beta.2
		{
			ConstraintString:   ">= 5.0.0 <= 6.0.0-beta.2",
			Version:            versions.Semver{Major: 6, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult:     true,
			IncludePreReleases: false,
		},
		// 6.0.0-beta.5 satisfies >= 5.0.0 <= 6.0.0-beta.2?
		//  yes, beta.5 is greater than beta.2
		{
			ConstraintString:   ">= 5.0.0 <= 6.0.0-beta.2",
			Version:            versions.Semver{Major: 6, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: false,
		},

		// While including pre releases
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 5, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   ">= 5.0.0-beta.2 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 2, Patch: 5, PreReleaseTag: "beta.5"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   ">= 5.0.0 <= 6.0.0",
			Version:            versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   ">= 5.0.0 <= 6.0.0-beta.2",
			Version:            versions.Semver{Major: 6, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
			ExpectedResult:     true,
			IncludePreReleases: true,
		},
		{
			ConstraintString:   ">= 5.0.0 <= 6.0.0-beta.2",
			Version:            versions.Semver{Major: 6, Minor: 0, Patch: 0, PreReleaseTag: "beta.5"},
			ExpectedResult:     false,
			IncludePreReleases: true,
		},
	}

	testConstraintsPreReleases(t, constraintsToTest)

}

func testConstraintsPreReleases(t *testing.T, constraintsToTest []ConstraintToTestPreReleases) {
	for _, constraintToTest := range constraintsToTest {
		fmt.Printf("\nTesting constraint evaluation. Does '%s' satisfy: '%s'\n", constraintToTest.Version.String(), constraintToTest.ConstraintString)
		res, err := validateConstraintPreReleases(constraintToTest)
		if err != nil {
			fmt.Printf("✗ failed parsing of constraint: '%s'. %s\n", constraintToTest.ConstraintString, err)
			t.Errorf("✗ failed parsing of constraint: '%s'. %s\n", constraintToTest.ConstraintString, err)
		} else {
			if res == false {
				fmt.Printf("✗ evaluation of constraint failed: '%s'\n", constraintToTest.ConstraintString)
				fmt.Printf("Expected %t, but got: %t", constraintToTest.ExpectedResult, !constraintToTest.ExpectedResult)
				t.Errorf("✗ evaluation of constraint failed: '%s'\n", constraintToTest.ConstraintString)
			} else {
				fmt.Println("✓ Success")
			}
		}
	}
}

func testConstraints(t *testing.T, constraintsToTest []ConstraintToTest) {
	for _, constraintToTest := range constraintsToTest {
		fmt.Printf("\nTesting constraint evaluation. Does '%s' satisfy: '%s'\n", constraintToTest.Version.String(), constraintToTest.ConstraintString)
		res, err := validateConstraint(constraintToTest)
		if err != nil {
			fmt.Printf("✗ failed parsing of constraint: '%s'. %s\n", constraintToTest.ConstraintString, err)
			t.Errorf("✗ failed parsing of constraint: '%s'. %s\n", constraintToTest.ConstraintString, err)
		} else {
			if res == false {
				fmt.Printf("✗ evaluation of constraint failed: '%s'\n", constraintToTest.ConstraintString)
				fmt.Printf("Expected %t, but got: %t", constraintToTest.ExpectedResult, !constraintToTest.ExpectedResult)
				t.Errorf("✗ evaluation of constraint failed: '%s'\n", constraintToTest.ConstraintString)
			} else {
				fmt.Println("✓ Success")
			}
		}
	}
}

func testMaxSatisfying(t *testing.T, maxSatisfyíngConstraintsToTest []MaxSatisfyíngConstraintToTest) {
	for _, maxSatisfyíngConstraintToTest := range maxSatisfyíngConstraintsToTest {
		fmt.Printf("\nTesting max satisfying evaluation for '%s' and versions '%v'\n", maxSatisfyíngConstraintToTest.ConstraintString, maxSatisfyíngConstraintToTest.Versions)
		parsedConstraint, err := constraintTypes.ParseConstraint(maxSatisfyíngConstraintToTest.ConstraintString)
		if err != nil {
			fmt.Printf("✗ failed parsing of constraint: '%s'. %s\n", maxSatisfyíngConstraintToTest.ConstraintString, err)
			t.Errorf("✗ failed parsing of constraint: '%s'. %s\n", maxSatisfyíngConstraintToTest.ConstraintString, err)
		} else {

			max, err := MaxSatisfyingStrings(maxSatisfyíngConstraintToTest.Versions, parsedConstraint, true)

			if err != nil {
				fmt.Printf("✗ failed version parsing of included versions: '%s'. %s\n", maxSatisfyíngConstraintToTest.ConstraintString, err)
				t.Errorf("✗ failed version parsing of included versions: '%s'. %s\n", maxSatisfyíngConstraintToTest.ConstraintString, err)
			}

			correct := true
			isNil := false

			// if maxSatisfyíngConstraintToTest.ExpectedMaxSatisfyingVersion == nil && max != nil {
			// 	correct = false
			// 	isNil = true
			// }
			// if maxSatisfyíngConstraintToTest.ExpectedMaxSatisfyingVersion != nil && max == nil {
			// 	correct = false
			// 	isNil = true
			// }
			// if maxSatisfyíngConstraintToTest.ExpectedMaxSatisfyingVersion == nil && max == nil {
			// 	correct = true
			// 	isNil = true
			// }

			if !isNil && !maxSatisfyíngConstraintToTest.ExpectedMaxSatisfyingVersion.EQ(max, true) {
				fmt.Printf("✗ evaluation of max satisfying failed: '%s'\n", maxSatisfyíngConstraintToTest.ConstraintString)
				fmt.Printf("Expected %s, but got: %s", maxSatisfyíngConstraintToTest.ExpectedMaxSatisfyingVersion.String(), max.String())
				t.Errorf("✗ evaluation of max satisfying failed: '%s'\n", maxSatisfyíngConstraintToTest.ConstraintString)
				correct = false
			}

			if correct {
				fmt.Println("✓ Success")
			}

		}
	}
}
