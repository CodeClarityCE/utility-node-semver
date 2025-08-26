package semver

import (
	"testing"
)

func TestComposerVersionParsing(t *testing.T) {
	tests := []struct {
		version     string
		ecosystem   EcosystemType
		expectDev   bool
		expectMajor int
		expectMinor int
		expectPatch int
	}{
		{"1.2.3", NodeJS, false, 1, 2, 3},
		{"1.2.3", Composer, false, 1, 2, 3},
		{"dev-master", Composer, true, 999, 999, 999},
		{"dev-feature-branch", Composer, true, 999, 999, 999},
		{"1.0.x-dev", Composer, true, 1, 0, 999},
		{"v2.1.0", NodeJS, false, 2, 1, 0},
		{"v2.1.0", Composer, false, 2, 1, 0},
	}

	for _, test := range tests {
		t.Run(test.version+" with "+string(test.ecosystem), func(t *testing.T) {
			parsed, err := ParseSemverWithEcosystem(test.version, test.ecosystem)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test.version, err)
			}

			if parsed.IsDev != test.expectDev {
				t.Errorf("Expected IsDev=%t, got %t", test.expectDev, parsed.IsDev)
			}
			if parsed.Major != test.expectMajor {
				t.Errorf("Expected Major=%d, got %d", test.expectMajor, parsed.Major)
			}
			if parsed.Minor != test.expectMinor {
				t.Errorf("Expected Minor=%d, got %d", test.expectMinor, parsed.Minor)
			}
			if parsed.Patch != test.expectPatch {
				t.Errorf("Expected Patch=%d, got %d", test.expectPatch, parsed.Patch)
			}
		})
	}
}

func TestComposerDevVersionComparison(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool // true if v1 > v2
	}{
		{"dev-master", "1.0.0", true},       // Dev version > stable
		{"1.0.0", "dev-master", false},      // Stable < dev version
		{"dev-master", "dev-master", false}, // Same dev version
		{"2.0.0", "1.9.9", true},            // Normal comparison
	}

	for _, test := range tests {
		t.Run(test.v1+" GT "+test.v2, func(t *testing.T) {
			v1, err := ParseSemverWithEcosystem(test.v1, Composer)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test.v1, err)
			}

			v2, err := ParseSemverWithEcosystem(test.v2, Composer)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test.v2, err)
			}

			result := v1.GT(v2, false)
			if result != test.expected {
				t.Errorf("Expected %s GT %s = %t, got %t", test.v1, test.v2, test.expected, result)
			}
		})
	}
}

func TestBackwardCompatibility(t *testing.T) {
	// Test that old functions still work
	version, err := ParseSemver("1.2.3")
	if err != nil {
		t.Fatalf("ParseSemver failed: %v", err)
	}
	if version.Major != 1 || version.Minor != 2 || version.Patch != 3 {
		t.Errorf("Expected 1.2.3, got %d.%d.%d", version.Major, version.Minor, version.Patch)
	}

	constraint, err := ParseConstraint(">=1.0.0")
	if err != nil {
		t.Fatalf("ParseConstraint failed: %v", err)
	}
	if constraint.Original != ">=1.0.0" {
		t.Errorf("Expected constraint original to be '>=1.0.0', got '%s'", constraint.Original)
	}
}
