package postupload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type FileTest struct {
	Src   string
	Dests []string
}

func NewTestRelease() *Release {
	return &Release{
		SourceDir:          "/tmp/src",
		BuildDir:           "build-dir",
		Product:            "product",
		NightlyDir:         "nightly",
		TinderboxBuildsDir: "tbox-win32",
	}
}

func TestReleaseToLatest(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()

	_, err := rel.ToLatest("/tmp/src/nobranch")
	assert.NotNil(err, "no Branch should trigger error.")

	rel.Branch = "l10n"

	_, err = rel.ToLatest("/etc/passwd")
	assert.NotNil(err, "Out of src file should trigger error")

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"product/nightly/latest-l10n/build-dir"}},
		FileTest{"/tmp/src/path1/path2/test.xpi", []string{"product/nightly/latest-l10n/build-dir/path1/path2"}},
		FileTest{"/tmp/src/mar.exe", []string{"product/nightly/latest-l10n/build-dir/mar-tools/win32"}},
	}
	for _, file := range files {
		dests, err := rel.ToLatest(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}

func TestReleaseToDated(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildID = "20150101223305"

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"product/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir"}},
		FileTest{"/tmp/src/path1/path2/test.xpi", []string{"product/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir/path1/path2"}},
	}

	for _, file := range files {
		dests, err := rel.ToDated(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}

func TestReleaseToCandidates(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildNumber = "23"
	rel.Version = "32"
	rel.Signed = false

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"product/nightly/32-candidates/build23/build-dir/subdir"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"product/nightly/32-candidates/build23/unsigned/build-dir/subdir"}},
		FileTest{"/tmp/src/mar.exe", []string{"product/nightly/32-candidates/build23/mar-tools/win32"}},
	}
	for _, file := range files {
		dests, err := rel.ToCandidates(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}

func TestReleaseToMobileCandidates(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildNumber = "23"
	rel.Version = "32"
	rel.Signed = false

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"product/nightly/32-candidates/build23/build-dir/subdir"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"product/nightly/32-candidates/build23/build-dir/subdir"}},
		FileTest{"/tmp/src/mar.exe", []string{"product/nightly/32-candidates/build23/build-dir"}},
	}
	for _, file := range files {
		dests, err := rel.ToMobileCandidates(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}

func TestReleaseToTryBuilds(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildNumber = "23"
	rel.Version = "32"
	rel.Signed = false
	rel.Who = "testuser"
	rel.Revision = "r33"

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"product/try-builds/testuser-r33/build-dir"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"product/try-builds/testuser-r33/build-dir"}},
		FileTest{"/tmp/src/mar.exe", []string{"product/try-builds/testuser-r33/build-dir"}},
	}
	for _, file := range files {
		dests, err := rel.ToTryBuilds(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}
