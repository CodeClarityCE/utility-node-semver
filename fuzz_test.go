package semver

import "testing"

// Range strings taken from the parser and evaluator tests.
var rangeSeeds = []string{
	"", "*", "*.*", "*.*.*", "ANY", "3", "3.x", "1.2", "1.2.x", "5.x.x", "5.x.x.",
	"!5.0.0", "!5.0.0-beta.2", "=5.0.0", "=5.0.0-beta.2", ">=1.0.0",
	"< 5.0.0", "< 5.0.0-beta.2", "<= 5.0.0", "> 0.0.0", "> 5.0.0", "> 5.0.0-beta.2", ">= 5.0.0",
	"1.2 - 2.3.4", "1.2.3 - 2", "1.2.3 - 2.3", "1.2.3 - 2.3.4", "1.2.3 - 2.3.4-beta.2", "1.2.3-beta.2 - 2.3.4",
	"> 0.0.0 < x", "> 5 < 7", "> 5.0.0 < 7.0.0", "> 5.0.0 <= 7.0.0", "> 5.0.0 <= 7.0.0-beta.2",
	"> 5.x < 7.0.0", "> 5.x < 7.x", "> x < 7", ">= 5 <= 7", ">= 5.0.0 <= 6.0.0", ">= 5.0.0 <= 6.0.0-beta.2",
	">= 5.0.0 <= 7.0.0", ">= 5.0.0-beta.2 <= 6.0.0", ">= 5.0.0-beta.2 <= 7.0.0", ">= 5.x <= 7.0.0", ">= 5.x <= 7.x",
	"^0", "^0.0", "^0.0.0", "^0.0.3", "^0.0.3-beta", "^0.0.x", "^0.2.3", "^0.x", "^0.x.x",
	"^1.2.3", "^1.2.3-beta.2", "^1.2.x", "^1.x", "^5.6.7",
	"~*", "~*.*", "~*.*.*", "~ANY", "~0", "~0.0", "~0.0.0", "~0.0.3", "~0.0.x", "~0.2", "~0.2.3", "~0.2.x",
	"~0.x", "~1", "~1.2", "~1.2.3", "~1.2.3-beta.2", "~1.2.x", "~1.x", "~5.6.7",
	"4.x && < 7.0.0", "5.2.x && 6.0.0", "5.x.x && < 5.2.5", "~4.5 && ^4.0.0",
	"5.2.6 || 5.2.7 || 5.2.8", ">= 5.0.0 < 5.2.5 || > 5.2.5 <= 5.5.0",
	"4.0.0 || 5.x && < 6.0.0 || 7.0.0", "4.0.0 || 5.x && < 6.0.0 || >= 7.0.0",
	"2.11.2 || 3.6.4 || 4.0.0-beta2 || 4.5.7 || >= 7.1.2",
	"< 2.11.2 || >= 3.0.0 < 3.6.4 || >= 4.0.0 < 4.5.7",
}

// Version strings taken from the ecosystem tests.
var versionSeeds = []string{"1.2.3", "v2.1.0", "5.0.0-beta.2", "dev-master", "dev-feature-branch", "1.0.x-dev"}

// FuzzParseConstraint checks that parsing any range, then evaluating it, never panics.
func FuzzParseConstraint(f *testing.F) {
	for i, r := range rangeSeeds {
		f.Add(r, versionSeeds[i%len(versionSeeds)])
	}
	f.Fuzz(func(t *testing.T, rangeString string, versionString string) {
		for _, ecosystem := range []EcosystemType{NodeJS, Composer} {
			c, err := ParseConstraintWithEcosystem(rangeString, ecosystem)
			if err != nil {
				continue
			}
			_ = c.String()
			v, err := ParseSemverWithEcosystem(versionString, ecosystem)
			if err != nil {
				continue
			}
			Satisfies(v, c, false)
			Satisfies(v, c, true)
		}
	})
}
