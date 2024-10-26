package semver

import (
	"sort"

	constraints "github.com/CodeClarityCE/utility-node-semver/constraints"
	evaluator "github.com/CodeClarityCE/utility-node-semver/evaluator"
	versions "github.com/CodeClarityCE/utility-node-semver/versions"
)

// Parses a given node semver constraint string into a constraint object
func ParseConstraint(constraintString string) (constraints.Constraint, error) {
	return constraints.ParseConstraint(constraintString)
}

// Parses a semver string into a semver object
func ParseSemver(versionLiteral string) (versions.Semver, error) {
	return versions.ParseSemver(versionLiteral)
}

// Takes a version and semver constraint
// Returns true if the version satisfies the constraint and false otherwise
//
//	ex: constraint '>= 5.0.0' and version '5.0.0' would return true
//
// As per node semver spec includePreReleases can be used to allow prerelease versions
// with different [major,minor,patch] tuple to satify the given constraint
//
//	ex: includePreReleases 'false' constraint '<= 5.0.0' and version '5.0.0-beta.2' would return true
//	ex: includePreReleases 'false' constraint '<= 5.0.0' and version '4.0.0-beta.2' would return false
//	ex: includePreReleases 'true' constraint '<= 5.0.0' and version '4.0.0-beta.2' would return true
func Satisfies(v versions.Semver, c constraints.Constraint, includePreReleases bool) bool {
	return evaluator.Satisfies(v, c, includePreReleases)
}

// Evaluates the given constraints for each provided version and returns the hightest version that satisfies this constraint (if any)
func MaxSatisfying(versions []versions.Semver, c constraints.Constraint, includePreReleases bool) versions.Semver {
	return evaluator.MaxSatisfying(versions, c, includePreReleases)
}

// Evaluates the given constraints for each provided version and returns the hightest version that satisfies this constraint (if any)
// Equivalent to MaxSatisfying, but this function allows users to pass in versions as strings
func MaxSatisfyingStrings(versions []string, c constraints.Constraint, includePreReleases bool) (versions.Semver, error) {
	return evaluator.MaxSatisfyingStrings(versions, c, includePreReleases)
}

// Sort sortes a given array of versions
//
// descending sort 	if sort order == -1
//
// ascending sort 	otherwise
func Sort(sortOrder int, versions []versions.Semver) []versions.Semver {
	sort.Slice(versions, func(i, j int) bool {
		return sortOrder == -1 && versions[i].GT(versions[j], false) ||
			sortOrder != -1 && versions[i].LT(versions[j], false)
	})
	return versions
}

// SortString sortes a given array of versions
// Equivalent to Sort, but this function allows users to pass in versions as strings
//
// descending sort 	if sort order == -1
//
// ascending sort 	otherwise
func SortStrings(sortOrder int, versionStrings []string) ([]string, error) {
	type VersionVariant struct {
		parsed versions.Semver
		raw    string
	}

	parsedVersions := make([]VersionVariant, len(versionStrings))
	for i, versionString := range versionStrings {
		parsedVersion, err := ParseSemver(versionString)
		if err != nil {
			return nil, err
		}
		parsedVersions[i] = VersionVariant{parsed: parsedVersion, raw: versionString}
	}

	sort.Slice(parsedVersions, func(i, j int) bool {
		if sortOrder == -1 {
			return parsedVersions[i].parsed.GT(parsedVersions[j].parsed, false)
		} else {
			return parsedVersions[i].parsed.LT(parsedVersions[j].parsed, false)
		}
	})

	sortedStrings := make([]string, len(parsedVersions))
	for i, versionVariant := range parsedVersions {
		sortedStrings[i] = versionVariant.raw
	}

	return sortedStrings, nil
}
