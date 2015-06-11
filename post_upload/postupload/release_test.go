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
		RootDir:            "pub",
		BuildDir:           "build-dir",
		Product:            "firefox",
		NightlyDir:         "nightly",
		TinderboxBuildsDir: "alder-win32",
	}
}

func mustBuildID(id string) *BuildID {
	b, err := NewBuildID(id)
	if err != nil {
		panic(err)
	}
	return b
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
		FileTest{"/tmp/src/crashreporter-symbols.zip", nil},
		FileTest{"/tmp/src/file.partial.foo.mar", nil},
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/nightly/latest-l10n/build-dir/file"}},
		FileTest{"/tmp/src/path1/path2/test.xpi", []string{"pub/firefox/nightly/latest-l10n/build-dir/path1/path2/test.xpi"}},
		FileTest{"/tmp/src/mar.exe", []string{"pub/firefox/nightly/latest-l10n/build-dir/mar-tools/win32/mar.exe"}},
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
	assert.NotPanics(func() {
		_, err := rel.ToDated("/tmp/src/file")
		assert.NotNil(err)
	})
	rel.BuildID = mustBuildID("20150101223305")

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir/file"}},
		FileTest{"/tmp/src/path1/path2/test.xpi", []string{"pub/firefox/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir/path1/path2/test.xpi"}},
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
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/nightly/32-candidates/build23/build-dir/subdir/file"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"pub/firefox/nightly/32-candidates/build23/unsigned/build-dir/subdir/win32-file"}},
		FileTest{"/tmp/src/mar.exe", []string{"pub/firefox/nightly/32-candidates/build23/mar-tools/win32/mar.exe"}},
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
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/nightly/32-candidates/build23/build-dir/subdir/file"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"pub/firefox/nightly/32-candidates/build23/build-dir/subdir/win32-file"}},
		FileTest{"/tmp/src/mar.exe", []string{"pub/firefox/nightly/32-candidates/build23/build-dir/mar.exe"}},
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
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/try-builds/testuser-r33/build-dir/file"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"pub/firefox/try-builds/testuser-r33/build-dir/win32-file"}},
		FileTest{"/tmp/src/mar.exe", []string{"pub/firefox/try-builds/testuser-r33/build-dir/mar.exe"}},
	}
	for _, file := range files {
		dests, err := rel.ToTryBuilds(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}

func TestReleaseToTinderboxBuilds(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()
	rel.TinderboxBuildsDir = "mozilla-aurora-l10n"

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/tinderbox-builds/mozilla-aurora-l10n/build-dir/file"}},
		FileTest{"/tmp/src/subdir/file.xpi", []string{"pub/firefox/tinderbox-builds/mozilla-aurora-l10n/build-dir/subdir/file.xpi"}},
		FileTest{"/tmp/src/subdir/file.mar", nil},
	}
	for _, file := range files {
		dests, err := rel.ToTinderboxBuilds(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}

func TestReleaseToDatedTinderboxBuilds(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()
	rel.TinderboxBuildsDir = "mozilla-aurora-l10n"
	rel.BuildID = mustBuildID("20150513000000")

	println(rel.BuildID.Time().String())

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/tinderbox-builds/mozilla-aurora-l10n/1431500400/build-dir/file"}},
		FileTest{"/tmp/src/subdir/file.xpi", []string{"pub/firefox/tinderbox-builds/mozilla-aurora-l10n/1431500400/build-dir/subdir/file.xpi"}},
		FileTest{"/tmp/src/subdir/file.mar", nil},
	}
	for _, file := range files {
		dests, err := rel.ToDatedTinderboxBuilds(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}
