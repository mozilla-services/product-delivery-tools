package postupload

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var platforms = []string{"win32", "macosx64", "linux", "linux64", "win64"}
var partialMarRe = regexp.MustCompile(`\.partial\..*\.mar(\.asc)?$`)
var enUsFilesToCopyToL10nRe = regexp.MustCompile(`\.en-US\.(win(32|64)\.zip|.*\.(checksums|complete\.mar|tar.bz2|dmg|exe)(.asc)?)$`)

// Release contains options for deploying files to S3
type Release struct {
	RootDir   string
	SourceDir string

	Branch             string
	BuildDir           string
	BuildID            *BuildID
	BuildNumber        string
	NightlyDir         string
	Product            string
	Revision           string
	ShortDir           bool
	Signed             bool
	SubDir             string
	TinderboxBuildsDir string
	Version            string
	Who                string
}

// NewRelease returns a new release
func NewRelease(sourceDir, product string) *Release {
	return &Release{
		Product:   product,
		RootDir:   "pub",
		SourceDir: sourceDir,
	}
}

func (r *Release) nightlyPath() string {
	return filepath.Join(r.RootDir, r.Product, r.NightlyDir)
}

func (r *Release) tinderboxBuildsPath() string {
	return filepath.Join(r.RootDir, r.Product, "tinderbox-builds", r.TinderboxBuildsDir)
}

func (r *Release) candidatesPath() string {
	return filepath.Join(r.RootDir, r.Product, "candidates")
}

func (r *Release) tryBuildsPath() string {
	return filepath.Join(r.RootDir, r.Product, "try-builds", r.Who+"-"+r.Revision, r.BuildDir)
}

func (r *Release) platform() string {
	if r.TinderboxBuildsDir == "" {
		return ""
	}

	for _, p := range platforms {
		if strings.HasSuffix(r.TinderboxBuildsDir, "-"+p) {
			return p
		}

	}
	return ""
}

func isMarTool(path string) bool {
	name := filepath.Base(path)
	switch name {
	case "mar", "mar.exe", "mbsdiff", "mbsdiff.exe":
		return true
	}
	return false
}

func (r *Release) resolvePath(src, dest string, preserveDir bool) ([]string, error) {
	if !strings.HasPrefix(src, r.SourceDir) {
		return nil, fmt.Errorf("%s not in %s", src, r.SourceDir)
	}
	if preserveDir {
		relPath, err := filepath.Rel(r.SourceDir, filepath.Dir(src))
		if err != nil {
			return nil, err
		}
		dest = filepath.Join(dest, relPath)
	}

	dest = filepath.Join(dest, filepath.Base(src))

	return []string{dest}, nil
}

// ToLatest returns destinations for nightly path
func (r *Release) ToLatest(file string) ([]string, error) {
	if r.Branch == "" {
		return nil, fmt.Errorf("ToLatest: Branch cannot be empty")
	}
	latestPath := r.generateLatestPath()
	marToolsPath := filepath.Join(latestPath, "mar-tools")

	if strings.HasSuffix(file, "crashreporter-symbols.zip") {
		return nil, nil
	}

	if partialMarRe.MatchString(file) {
		return nil, nil
	}

	if strings.HasSuffix(r.Branch, "l10n") && strings.HasSuffix(file, ".xpi") {
		return r.resolvePath(file, latestPath, true)
	}

	if isMarTool(file) {
		if platform := r.platform(); platform != "" {
			return r.resolvePath(file, filepath.Join(marToolsPath, platform), false)
		}
		return nil, nil
	}

	regularDests, err := r.resolvePath(file, latestPath, false)

	if err == nil &&
		!strings.HasSuffix(r.Branch, "l10n") &&
		enUsFilesToCopyToL10nRe.MatchString(file) {

		l10nPath := r.generateLatestPathWithSuffix("-l10n")
		l10nDests, l10nErr := r.resolvePath(file, l10nPath, false)

		finalDests := append(regularDests, l10nDests...)

		return finalDests, l10nErr
	}

	return regularDests, err
}

func (r *Release) generateLatestPath() string {
	return r.generateLatestPathWithSuffix("")
}

func (r *Release) generateLatestPathWithSuffix(branchSuffix string) string {
	latestPath := filepath.Join(r.nightlyPath(), "latest-"+r.Branch+branchSuffix)
	if r.BuildDir != "" {
		latestPath = filepath.Join(latestPath, r.BuildDir)
	}
	return latestPath
}

// ToDated returns destinations for dated
func (r *Release) ToDated(file string) ([]string, error) {
	if r.BuildID == nil {
		return nil, errors.New("BuildID cannot be empty")
	}
	bID := r.BuildID

	longDate := fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s",
		bID.Year(), bID.Month(), bID.Day(), bID.Hour(), bID.Minute(), bID.Second(), r.Branch)
	longDatedPath := filepath.Join(r.nightlyPath(), bID.Year(), bID.Month(), longDate)

	if r.BuildDir != "" {
		longDatedPath = filepath.Join(longDatedPath, r.BuildDir)
	}

	if strings.HasSuffix(r.Branch, "l10n") && strings.HasSuffix(file, ".xpi") {
		return r.resolvePath(file, longDatedPath, true)
	}

	return r.resolvePath(file, longDatedPath, false)
}

// ToCandidates returns destinations for candidates
func (r *Release) ToCandidates(file string) ([]string, error) {
	path := filepath.Join(r.candidatesPath(), r.Version+"-candidates", "build"+r.BuildNumber)
	marToolsPath := filepath.Join(path, "mar-tools")

	if !r.Signed && strings.Contains(file, "win32") && !strings.Contains(file, "/logs/") {
		path = filepath.Join(path, "unsigned")
	}

	path = filepath.Join(path, r.BuildDir)

	if isMarTool(file) {
		if platform := r.platform(); platform != "" {
			return r.resolvePath(file, filepath.Join(marToolsPath, platform), false)
		}
	}

	return r.resolvePath(file, path, true)
}

// ToMobileCandidates returns destinations for mobile candidates
func (r *Release) ToMobileCandidates(file string) ([]string, error) {
	path := filepath.Join(r.candidatesPath(), r.Version+"-candidates", "build"+r.BuildNumber, r.BuildDir)
	return r.resolvePath(file, path, true)
}

// ToTinderboxBuilds returns destinations for tinderbox builds
func (r *Release) ToTinderboxBuilds(file string) ([]string, error) {
	path := filepath.Join(r.tinderboxBuildsPath(), r.BuildDir)
	if strings.HasSuffix(file, ".mar") {
		return nil, nil
	}

	if strings.HasSuffix(r.TinderboxBuildsDir, "l10n") && strings.HasSuffix(file, ".xpi") {
		return r.resolvePath(file, path, true)
	}

	return r.resolvePath(file, path, false)
}

// ToDatedTinderboxBuilds returns destinations for dated tinderbox builds
func (r *Release) ToDatedTinderboxBuilds(file string) ([]string, error) {
	if r.BuildID == nil {
		return nil, errors.New("BuildID cannot be empty")
	}
	path := filepath.Join(r.tinderboxBuildsPath(), fmt.Sprintf("%d", r.BuildID.Time().Unix()), r.BuildDir)

	if strings.HasSuffix(file, ".mar") {
		return nil, nil
	}

	if strings.HasSuffix(r.TinderboxBuildsDir, "l10n") && strings.HasSuffix(file, ".xpi") {
		return r.resolvePath(file, path, true)
	}

	return r.resolvePath(file, path, false)

}

// ToTryBuilds returns destinations for try builds
func (r *Release) ToTryBuilds(file string) ([]string, error) {
	if r.Who == "" {
		return nil, errors.New("Who cannot be empty")
	}
	if r.Revision == "" {
		return nil, errors.New("Revision cannot be empty")
	}
	if r.Product == "" {
		return nil, errors.New("Product cannot be empty")
	}
	return r.resolvePath(file, r.tryBuildsPath(), false)
}
