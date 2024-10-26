package evaluator

import (
	"github.com/CodeClarityCE/utility-node-semver/utils"

	versionTypes "github.com/CodeClarityCE/utility-node-semver/versions"

	constraints "github.com/CodeClarityCE/utility-node-semver/constraints"
)

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
func Satisfies(v versionTypes.Semver, c constraints.Constraint, includePreReleases bool) bool {

	conjunctedConditions := []bool{}
	joinIndex := -1

	for idx, cRange := range c.Ranges {

		endversionEmpty := cRange.EndVersion == versionTypes.Semver{}
		startversionEmpty := cRange.StartVersion == versionTypes.Semver{}
		if !endversionEmpty && !startversionEmpty {

			startSatisified := false
			endSatisified := false

			if cRange.StartOp == constraints.GT {
				startSatisified = v.GT(cRange.StartVersion, false)
			} else if cRange.StartOp == constraints.GE {
				startSatisified = v.GE(cRange.StartVersion, false)
			}

			if !startSatisified {
				conjunctedConditions = append(conjunctedConditions, false)
			} else {
				if cRange.EndOp == constraints.LT {
					endSatisified = v.LT(cRange.EndVersion, false)
				} else if cRange.EndOp == constraints.LE {
					endSatisified = v.LE(cRange.EndVersion, false)
				}

				// According to the nodesemver spec:
				//   "If a version has a prerelease tag (for example, 1.2.3-alpha.3) then it
				//   will only be allowed to satisfy comparator sets if at least one comparator with the
				//   same [major, minor, patch] tuple also has a prerelease tag."
				//
				// This behavior can be suppressed (treating all prerelease versions as if they were normal
				// versions, for the purpose of range matching) by setting the includePreReleases
				if !includePreReleases && (v.PreReleaseTag != "" || cRange.StartVersion.PreReleaseTag != "" || cRange.EndVersion.PreReleaseTag != "") {

					if cRange.StartVersion.PreReleaseTag == "" && cRange.EndVersion.PreReleaseTag == "" {
						conjunctedConditions = append(conjunctedConditions, false)
					} else {
						cmpStart := v.EQ(cRange.StartVersion, true)
						cmpEnd := v.EQ(cRange.EndVersion, true)

						if !cmpStart && !cmpEnd {
							conjunctedConditions = append(conjunctedConditions, false)
						} else {
							conjunctedConditions = append(conjunctedConditions, startSatisified && endSatisified)
						}
					}

				} else {
					conjunctedConditions = append(conjunctedConditions, startSatisified && endSatisified)
				}

			}

		} else {

			res := false

			if cRange.StartOp == constraints.GT {
				res = v.GT(cRange.StartVersion, false)
			} else if cRange.StartOp == constraints.GE {
				res = v.GE(cRange.StartVersion, false)
			} else if cRange.StartOp == constraints.EQ {
				res = v.EQ(cRange.StartVersion, false)
			} else if cRange.StartOp == constraints.LT {
				res = v.LT(cRange.StartVersion, false)
			} else if cRange.StartOp == constraints.LE {
				res = v.LE(cRange.StartVersion, false)
			}

			// According to the nodesemver spec:
			//   "If a version has a prerelease tag (for example, 1.2.3-alpha.3) then it
			//   will only be allowed to satisfy comparator sets if at least one comparator with the
			//   same [major, minor, patch] tuple also has a prerelease tag."
			//
			// This behavior can be suppressed (treating all prerelease versions as if they were normal
			// versions, for the purpose of range matching) by setting the includePreReleases
			if !includePreReleases && (v.PreReleaseTag != "" || cRange.StartVersion.PreReleaseTag != "") {

				if cRange.StartVersion.PreReleaseTag == "" {
					res = false
				} else {
					cmpStart := v.EQ(cRange.StartVersion, true)

					if !cmpStart {
						res = false
					}
				}

			}

			conjunctedConditions = append(conjunctedConditions, res)

		}

		joinIndex++

		if idx <= len(c.Join)-1 && c.Join[joinIndex] == constraints.DISJUNCTION {
			if utils.ContainsOnly(conjunctedConditions, true) {
				return true
			}
			conjunctedConditions = []bool{}
		}

	}

	return utils.ContainsOnly(conjunctedConditions, true)
}

// Evaluates the given constraints for each provided version and returns the hightest version that satisfies this constraint (if any)
func MaxSatisfying(versions []versionTypes.Semver, c constraints.Constraint, includePreReleases bool) versionTypes.Semver {
	if len(versions) == 0 {
		return versionTypes.Semver{}
	}
	max := versions[0]
	for _, version := range versions {
		if Satisfies(version, c, includePreReleases) {
			if version.GT(max, false) {
				max = version
			}
		}
	}
	return max
}

// Evaluates the given constraints for each provided version and returns the hightest version that satisfies this constraint (if any)
// Equivalent to MaxSatisfying, but this function allows users to pass in versions as strings
func MaxSatisfyingStrings(versions []string, c constraints.Constraint, includePreReleases bool) (versionTypes.Semver, error) {
	parsedVersions := []versionTypes.Semver{}
	for _, versionString := range versions {
		parsedVersion, err := versionTypes.ParseSemver(versionString)
		if err != nil {
			return versionTypes.Semver{}, err
		}
		parsedVersions = append(parsedVersions, parsedVersion)
	}
	return MaxSatisfying(parsedVersions, c, includePreReleases), nil
}
