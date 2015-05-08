package post_upload

func ReleaseToLatest(branch, tinderboxBuildsDir, uploadDir string, files []string) {

}

func ReleaseToDated(branch, buildId, product, nightlyDir string, shortDir bool, uploadDir string, files []string) {

}

func ReleaseToCandidates(buildDir, buildNumber, product, tinderboxBuildsDir,
	version string, signed bool, uploadDir string, files []string) {

}
func ReleaseToMobileCandidates(version, buildNumber, nightlyDir, product, uploadDir string, files []string) {

}

func ReleaseToTinderboxBuilds(product, buildId, buildDir, tinderboxBuildsDir, uploadDir string, files []string) {

}

func ReleaseToDatedTinderboxBuilds(product, buildId, buildDir, tinderboxBuildsDir, uploadDir string, files []string) {

}

func ReleaseToTryBuilds(product, who, revision, buildDir, uploadDir string, files []string) {

}
