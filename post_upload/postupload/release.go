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

// Release contains options for deploying files to S3
type Release struct {
	SourceDir string

	FtpCopier Copier
	PvtCopier Copier

	FtpPrefix string
	PvtPrefix string

	Branch             string
	BuildDir           string
	BuildID            BuildID
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

// NewS3Release returns a new release with s3 copiers.
func NewS3Release(ftpBucket, pvtBucket string) *Release {
	return &Release{
		FtpCopier: &S3Copier{Bucket: ftpBucket},
		PvtCopier: &S3Copier{Bucket: pvtBucket},
	}

}

func (r *Release) nightlyPath() string {
	return filepath.Join(r.FtpPrefix, r.Product, r.NightlyDir)
}

func (r *Release) tinderboxBuildsPath() string {
	return filepath.Join(r.FtpPrefix, r.Product, "tinderbox-builds", r.TinderboxBuildsDir)
}

func (r *Release) candidatesPath() string {
	return filepath.Join(r.FtpPrefix, r.Product, "candidates")
}

func (r *Release) pvtBuildsPath() string {
	return filepath.Join(r.PvtPrefix, r.Product, r.TinderboxBuildsDir)
}

func (r *Release) tryBuildsPath() string {
	return filepath.Join(r.FtpPrefix, r.Product, "try-builds", r.Who+"-"+r.Revision, r.BuildDir)
}

func (r *Release) platform() string {
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

func (r *Release) copyFile(src, dest string, preserveDir bool, copier Copier) error {
	if !strings.HasPrefix(src, r.SourceDir) {
		return fmt.Errorf("%s not in %s", src, r.SourceDir)
	}
	if preserveDir {
		relPath, err := filepath.Rel(r.SourceDir, filepath.Dir(src))
		if err != nil {
			return err
		}
		dest = filepath.Join(dest, relPath)
	}

	return copier.Copy(src, dest)
}

// ToLatest copies files to nightly path
func (r *Release) ToLatest(file string) error {
	if r.Branch == "" {
		return fmt.Errorf("ToLatest: Branch cannot be empty")
	}
	latestPath := filepath.Join(r.nightlyPath(), "latest-"+r.Branch)
	if r.BuildDir != "" {
		latestPath = filepath.Join(latestPath, r.BuildDir)
	}
	marToolsPath := filepath.Join(latestPath, "mar-tools")

	if strings.HasSuffix(file, "crashreporter-symbols.zip") {
		return nil
	}

	if partialMarRe.MatchString(file) {
		return nil
	}

	if strings.HasSuffix(r.Branch, "l10n") && strings.HasSuffix(file, ".xpi") {
		return r.copyFile(file, latestPath, true, r.FtpCopier)
	}

	if isMarTool(file) {
		if platform := r.platform(); platform != "" {
			return r.copyFile(file, filepath.Join(marToolsPath, platform), false, r.FtpCopier)
		}
		return nil
	}

	return r.copyFile(file, latestPath, false, r.FtpCopier)
}

// ToDated copies files to dated
func (r *Release) ToDated(file string) error {
	bID := BuildID(r.BuildID)
	if !bID.Validate() {
		return errors.New("buildID is not valid")
	}

	longDate := fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s",
		bID.Year(), bID.Month(), bID.Day(), bID.Hour(), bID.Minute(), bID.Second(), r.Branch)
	longDatedPath := filepath.Join(r.nightlyPath(), bID.Year(), bID.Month(), longDate)

	if r.BuildDir != "" {
		longDatedPath = filepath.Join(longDatedPath, r.BuildDir)
	}

	if strings.HasSuffix(r.Branch, "l10n") && strings.HasSuffix(file, ".xpi") {
		return r.copyFile(file, longDatedPath, true, r.FtpCopier)
	}

	return r.copyFile(file, longDatedPath, false, r.FtpCopier)
}

// ToCandidates copies files to candidates
func (r *Release) ToCandidates(file string) error {
	path := filepath.Join(r.nightlyPath(), r.Version+"-candidates", "build"+r.BuildNumber)
	marToolsPath := filepath.Join(path, "mar-tools")

	if !r.Signed && strings.Contains(file, "win32") && !strings.Contains(file, "/logs/") {
		path = filepath.Join(path, "unsigned")
	}

	path = filepath.Join(path, r.BuildDir)

	if isMarTool(file) {
		if platform := r.platform(); platform != "" {
			return r.copyFile(file, filepath.Join(marToolsPath, platform), true, r.FtpCopier)
		}
	}

	return r.copyFile(file, path, true, r.FtpCopier)
}

// ToMobileCandidates copies files to mobile candidates
func (r *Release) ToMobileCandidates(file string) error {
	path := filepath.Join(r.nightlyPath(), r.Version+"-candidates", "build"+r.BuildNumber, r.BuildDir)
	return r.copyFile(file, path, true, r.FtpCopier)
}

// ToTinderboxBuilds copies files to tinderbox builds
func (r *Release) ToTinderboxBuilds(file string) error {
	return nil

}

// ToDatedTinderboxBuilds copies files to dated tinderbox builds
func (r *Release) ToDatedTinderboxBuilds(file string) error {
	return nil

}

// ToTryBuilds copies files to try builds
func (r *Release) ToTryBuilds(file string) error {
	if r.Who == "" {
		return errors.New("Who cannot be empty")
	}
	if r.Revision == "" {
		return errors.New("Revision cannot be empty")
	}
	if r.Product == "" {
		return errors.New("Product cannot be empty")
	}
	return r.copyFile(file, r.tryBuildsPath(), false, r.FtpCopier)
}
