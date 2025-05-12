package constraints

import (
	"fmt"
	"testing"

	"github.com/CodeClarityCE/utility-node-semver/versions"

	"slices"
)

type RangeToTest struct {
	ConstraintString string
	Constraint       Constraint
}

func compareVersion(v1 versions.Semver, v2 versions.Semver) bool {
	return v1.Major == v2.Major && v1.Minor == v2.Minor && v1.Patch == v2.Patch && v1.PreReleaseTag == v2.PreReleaseTag && v1.MetaData == v2.MetaData
}

func compareRange(r1 Range, r2 Range) bool {
	return r1.StartOp == r2.StartOp && compareVersion(r1.StartVersion, r2.StartVersion) && r1.EndOp == r2.EndOp && compareVersion(r1.EndVersion, r2.EndVersion)
}

func compareConstraint(c1 Constraint, c2 Constraint) bool {

	same := true
	same = slices.Equal(c1.Join, c2.Join)

	if !same {
		return false
	}

	same = len(c1.Ranges) == len(c2.Ranges)

	if !same {
		return false
	}

	for idx, parsedRange := range c1.Ranges {
		if !compareRange(parsedRange, c2.Ranges[idx]) {
			return false
		}
	}

	return true

}

func TestTildeRangeParsing(t *testing.T) {

	fmt.Printf("\n%s Testing tilde (~) range parsing %s\n", "----------------", "----------------")

	// ^*      	-->  	(any)
	// ^1.2.3  	-->  	>=1.2.3 <2.0.0
	// ^1.2    	-->  	>=1.2.0 <2.0.0
	// ^1.2.x   -->  	>=1.2.0 <2.0.0
	// ^1      	-->  	>=1.0.0 <2.0.0
	// ^1.x     -->  	>=1.0.0 <2.0.0
	// ^0.2.3  	-->  	>=0.2.3 <0.3.0
	// ^0.2    	-->  	>=0.2.0 <0.3.0
	// ^0.2.x   -->  	>=0.2.0 <0.3.0
	// ^0.0.3  	-->  	>=0.0.3 <0.0.4
	// ^0.0    	-->  	>=0.0.0 <0.1.0
	// ^0.0.x   -->  	>=0.0.0 <0.1.0
	// ^0      	-->  	>=0.0.0 <1.0.0
	// ^0.x     -->  	>=0.0.0 <1.0.0

	rangesToTest := []RangeToTest{
		{
			ConstraintString: "~*",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~*.*",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~*.*.*",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~ANY",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~1.2.3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~1.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~1.2.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~1",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 2, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~1.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 2, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0.2.3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 2, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 2, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0.2.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 2, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0.0.3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 1, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 1, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0.0.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 1, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~0.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "~1.2.3-beta.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3, PreReleaseTag: "beta.2"},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
	}

	testConstraints(t, rangesToTest)

	fmt.Printf("\n")

}

func TestHyphenatedRangeParsing(t *testing.T) {

	fmt.Printf("\n%s Testing hyphen (-) range parsing %s\n", "----------------", "----------------")

	rangesToTest := []RangeToTest{
		{
			ConstraintString: "1.2.3 - 2.3.4",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3},
						EndOp: LE, EndVersion: versions.Semver{Major: 2, Minor: 3, Patch: 4},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "1.2.3-beta.2 - 2.3.4",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3, PreReleaseTag: "beta.2"},
						EndOp: LE, EndVersion: versions.Semver{Major: 2, Minor: 3, Patch: 4},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "1.2.3 - 2.3.4-beta.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3},
						EndOp: LE, EndVersion: versions.Semver{Major: 2, Minor: 3, Patch: 4, PreReleaseTag: "beta.2"},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "1.2 - 2.3.4",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 0},
						EndOp: LE, EndVersion: versions.Semver{Major: 2, Minor: 3, Patch: 4},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "1.2.3 - 2.3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 2, Minor: 4, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "1.2.3 - 2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 3, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
	}

	testConstraints(t, rangesToTest)

	fmt.Printf("\n")
}

func TestCaretRangeParsing(t *testing.T) {

	fmt.Printf("\n%s Testing caret (^) range parsing %s\n", "----------------", "----------------")

	rangesToTest := []RangeToTest{
		{
			ConstraintString: "^1.2.3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 2, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^0.2.3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 2, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^0.0.3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 3},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 0, Patch: 4},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^1.2.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 2, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^0.0.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 1, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 1, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^1.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 2, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^0.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^0.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^1.2.3-beta.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 3, PreReleaseTag: "beta.2"},
						EndOp: LT, EndVersion: versions.Semver{Major: 2, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "^0.0.3-beta",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 3, PreReleaseTag: "beta"},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 0, Patch: 4},
					},
				},
				Join: []JoinOp{},
			},
		},
	}

	testConstraints(t, rangesToTest)

	fmt.Printf("\n")

}

func TestXRangeParsing(t *testing.T) {

	fmt.Printf("\n%s Testing X range parsing %s\n", "----------------", "----------------")

	rangesToTest := []RangeToTest{
		{
			ConstraintString: "ANY",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "*",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "*.*",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "*.*.*",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "3.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 3, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 4, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "1.2.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "3",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 3, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 4, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "1.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 1, Minor: 2, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 1, Minor: 3, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
	}

	testConstraints(t, rangesToTest)

	fmt.Printf("\n")

}

func TestRangeParsing(t *testing.T) {

	fmt.Printf("\n%s Testing range (>= x.x.x =< x.x.x) parsing %s\n", "----------------", "----------------")

	rangesToTest := []RangeToTest{
		{
			ConstraintString: ">= 5.0.0 <= 7.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LE, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 5.0.0 <= 7.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LE, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: ">= 5.0.0-beta.2 <= 7.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
						EndOp: LE, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 5.0.0 <= 7.0.0-beta.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LE, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 5.0.0 < 7.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: ">= 5.x <= 7.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LE, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: ">= 5.x <= 7.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LE, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: ">= 5 <= 7",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GE, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LE, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 5.x < 7.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 5.x < 7.x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 5 < 7",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> x < 7",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 7, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 0.0.0 < x",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
						EndOp: LT, EndVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "> 0.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: GT, StartVersion: versions.Semver{Major: 0, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		{
			ConstraintString: "<= 5.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: LE, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
	}

	testConstraints(t, rangesToTest)

	fmt.Printf("\n")

}

func TestStaticParsing(t *testing.T) {

	fmt.Printf("\n%s Testing static equality operator (= 5.0.0, !7.0.0) parsing %s\n", "----------------", "----------------")

	rangesToTest := []RangeToTest{
		// {
		// 	ConstraintString: "!5.0.0",
		// 	Constraint: Constraint{
		// 		Ranges: []Range{
		// 			{
		// 				StartOp: NOT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
		// 			},
		// 		},
		// 		Join: []JoinOp{},
		// 	},
		// },
		{
			ConstraintString: "=5.0.0",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: EQ, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0},
					},
				},
				Join: []JoinOp{},
			},
		},
		// {
		// 	ConstraintString: "!5.0.0-beta.2",
		// 	Constraint: Constraint{
		// 		Ranges: []Range{
		// 			{
		// 				StartOp: NOT, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
		// 			},
		// 		},
		// 		Join: []JoinOp{},
		// 	},
		// },
		{
			ConstraintString: "=5.0.0-beta.2",
			Constraint: Constraint{
				Ranges: []Range{
					{
						StartOp: EQ, StartVersion: versions.Semver{Major: 5, Minor: 0, Patch: 0, PreReleaseTag: "beta.2"},
					},
				},
				Join: []JoinOp{},
			},
		},
	}

	testConstraints(t, rangesToTest)

	fmt.Printf("\n")
}

func testConstraints(t *testing.T, rangesToTest []RangeToTest) {
	for _, constraintToTest := range rangesToTest {
		fmt.Printf("\nTesting range parsing: '%s'\n", constraintToTest.ConstraintString)
		parsedConstraint, err := ParseConstraint(constraintToTest.ConstraintString)
		if err != nil {
			fmt.Printf("✗ range parsing test failed for range: '%s'. %s\n", constraintToTest.ConstraintString, err)
			t.Errorf("✗ range parsing test failed for range: '%s'. %s\n", constraintToTest.ConstraintString, err)
		} else {
			if !compareConstraint(parsedConstraint, constraintToTest.Constraint) {
				fmt.Printf("✗ range parsing test failed for range: '%s'\n", constraintToTest.ConstraintString)
				fmt.Printf("Expected %s, but got: %s", constraintToTest.Constraint.Ranges[0], parsedConstraint.Ranges[0])
				t.Errorf("✗ range parsing test failed for range: '%s'\n", constraintToTest.ConstraintString)
			} else {
				fmt.Println("✓ Success")
			}
		}
	}
}
