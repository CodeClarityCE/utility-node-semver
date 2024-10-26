# Library - node semver (go)

<br>

<div align="center">
    <img src="https://user-images.githubusercontent.com/124595411/235138790-d86cc2b8-e3ef-43eb-846c-38055748c9db.svg" width="400px" />
</div>

<br>

## Purpose

The node semver go lang library provides support for parsing semver 2.0 spec versions and node semver spec constraints.
The parsed objects can then be used for comparing versions against each other and to evaluate whether a versions satisfies a given constraint.

<br>

## Usage

Recommended usage as a git submodule


1. Add the submodule `git submodule add git@github.com:CodeClarityCE/utility-node-semver.git`
2. Add as dependency `go work use ./node-semver`
4. Import:
    ```go
    import (
        semver "codeclarity.io/node-semver"
    )
    ```
5. Usage:
    ```go

    v1, err := semver.ParseSemver("5.0.0-beta.2")
    if err != nil {
        // ... handle error here
    }

    v2, err := semver.ParseSemver("5.0.0")
    if err != nil {
        // ... handle error here
    }

    v3, err := semver.ParseSemver("4.0.0-beta.2")
    if err != nil {
        // ... handle error here
    }

    c1, err := semver.ParseConstraint("<= 5.0.0")
    if err != nil {
        // ... handle error here
    }

    v1.LT(v2, false) // true
    v1.EQ(v2, false) // false

    v3.LT(v2, false) // true

    semver.Satisfies(v1, c1, false) // true
    semver.Satisfies(v3, c1 false) // false *
    semver.Satisfies(v3, c1 true) // true *

    // * As per node semver spec includePreReleases can be used to allow prerelease versions    
    // with different [major,minor,patch] tuple to satify the given constraint
    //
    // This means:
    //    5.0.0-beta.2 does satisfy the constraint <= 5.0.0
    //    because the [major,minor,patch] is the same
    //
    // but
    //    4.0.0-beta.2 does not satisfy the constraint <= 5.0.0
    //    unless includePreReleases is set to true
    //
    // Read more here: https://github.com/npm/node-semver#prerelease-tags

    ```

<br>

## APIs

- `semver.ParseConstraint(constraintString string) (*Constraint, error)` parses a node semver spec constraint and returns the parsed constraint or an error
- `semver.ParseSemver(versionLiteral string) (*Semver, error)` parses a semver 2.0 spec version and returns the parsed version or an error
- `(v1 *Semver) GE(v2 *Semver, ignorePreRelease bool) bool` evaluates whether v1 is greater or equal to v2. If ignorePreRelease is set to true then the prerelease tag (if any) is ignored, in which case 4.0.0 would equal 4.0.0-beta.2
- `(v1 *Semver) GT(v2 *Semver, ignorePreRelease bool) bool` evaluates whether v1 is strictly greater than v2. If ignorePreRelease is set to true then the prerelease tag (if any) is ignored, in which case 4.0.0 would equal 4.0.0-beta.2
- `(v1 *Semver) LT(v2 *Semver, ignorePreRelease bool) bool` evaluates whether v1 is smaller or equal to v2. If ignorePreRelease is set to true then the prerelease tag (if any) is ignored, in which case 4.0.0 would equal 4.0.0-beta.2
- `(v1 *Semver) LE(v2 *Semver, ignorePreRelease bool) bool` evaluates whether v1 is strictly smaller than to v2. If ignorePreRelease is set to true then the prerelease tag (if any) is ignored, in which case 4.0.0 would equal 4.0.0-beta.2
- `(v1 *Semver) EQ(v2 *Semver, ignorePreRelease bool) bool` evaluates whether v1 is equal to v2. If ignorePreRelease is set to true then the prerelease tag (if any) is ignored, in which case 4.0.0 would equal 4.0.0-beta.2
- `(v1 *Semver) NEQ(v2 *Semver, ignorePreRelease bool) bool` evaluates whether v1 is unequal to v2. If ignorePreRelease is set to true then the prerelease tag (if any) is ignored, in which case 4.0.0 would equal 4.0.0-beta.2
- `(v1 *Semver) Compare(v2 *Semver, ignorePreRelease bool) int` compares v1 to v2 and returns 0 if both are equal, 1 if v1 is strictly greater than v2 and returns -1 if v1 is strictly smaller than v2

- `semver.Satisfies(v *versionTypes.Semver, c *constraints.Constraint, includePreReleases bool) bool`evaluates whether the version `v` satisfies the constraint `c`. As per node semver spec `includePreReleases` can be used to allow prerelease versions with different `[major,minor,patch]` tuple to satify the given constraint
    
   - This means:
        `5.0.0-beta.2` does satisfy the constraint `<= 5.0.0`
        because the `[major,minor,patch]` is the same
    
   - but
        `4.0.0-beta.2` does not satisfy the constraint `<= 5.0.0`
        unless `includePreReleases` is set to true

- `semver.MaxSatisfying(versions []*versionTypes.Semver, c *constraints.Constraint, includePreReleases bool) *versionTypes.Semver` Evaluates the given constraints for each provided version and returns the hightest version that satisfies this constraint  (if any).

- `semver.MaxSatisfyingStrings(versions []string, c *constraints.Constraint, includePreReleases bool) (*versionTypes.Semver, error)` Evaluates the given constraints for each provided version and returns the hightest version that satisfies this constraint  (if any). Equivalent to `semver.MaxSatisfying`, but this function allows users to pass in versions as strings.

- `semver.Sort(sortOrder int, versions []*versions.Semver) []*versions.Semver` sortes a given array of versions. 
    - descending sort 	if sort order == -1 
    - ascending sort 	otherwise

- `semver.SortStrings(sortOrder int, versionStrings []string) ([]string, error)` sortes a given array of versions. Equivalent to `semver.Sort`, but this function allows users to pass in versions as strings
    - descending sort 	if sort order == -1 
    - ascending sort 	otherwise