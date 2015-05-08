package postupload

// ReleaseToLatest copies files to nightly path
func ReleaseToLatest(branch, tinderboxBuildsDir, uploadDir string, files []string) {

}

// ReleaseToDated copies files to dated
func ReleaseToDated(branch, buildID, product, nightlyDir string, shortDir bool, uploadDir string, files []string) {

}

// ReleaseToCandidates copies files to candidates
func ReleaseToCandidates(buildDir, buildNumber, product, tinderboxBuildsDir,
	version string, signed bool, uploadDir string, files []string) {

}

// ReleaseToMobileCandidates copies files to mobile candidates
func ReleaseToMobileCandidates(version, buildNumber, nightlyDir, product, uploadDir string, files []string) {

}

// ReleaseToTinderboxBuilds copies files to tinderbox builds
func ReleaseToTinderboxBuilds(product, buildID, buildDir, tinderboxBuildsDir, uploadDir string, files []string) {

}

// ReleaseToDatedTinderboxBuilds copies files to dated tinderbox builds
func ReleaseToDatedTinderboxBuilds(product, buildID, buildDir, tinderboxBuildsDir, uploadDir string, files []string) {

}

// ReleaseToTryBuilds copies files to try builds
func ReleaseToTryBuilds(product, who, revision, buildDir, uploadDir string, files []string) {

}
