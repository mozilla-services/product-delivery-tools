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
	release := NewRelease("/tmp/src", "firefox")
	release.BuildDir = "build-dir"
	release.NightlyDir = "nightly"
	release.TinderboxBuildsDir = "alder-win32"
	return release
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
		FileTest{"/tmp/src/host/bin/mar.exe", []string{"pub/firefox/nightly/latest-l10n/build-dir/mar-tools/win32/mar.exe"}},
		FileTest{"/tmp/src/mar.exe", []string{"pub/firefox/nightly/latest-l10n/build-dir/mar-tools/win32/mar.exe"}},
	}
	for _, file := range files {
		dests, err := rel.ToLatest(file.Src)
		assert.Nil(err)
		assert.Equal(file.Dests, dests)
	}
}

func TestEnUsToLatest(t *testing.T) {
	rel := NewTestRelease()
	rel.Branch = "mozilla-central"

	filesThatMustBeDuplicated := []string{
		"firefox-55.0a1.en-US.linux-i686.checksums",
		"firefox-55.0a1.en-US.linux-i686.checksums.asc",
		"firefox-55.0a1.en-US.linux-i686.complete.mar",
		"firefox-55.0a1.en-US.linux-x86_64.tar.bz2",
		"firefox-55.0a1.en-US.linux-x86_64.tar.bz2.asc",
		"firefox-55.0a1.en-US.mac.dmg",
		"firefox-55.0a1.en-US.win32.installer-stub.exe",
		"firefox-55.0a1.en-US.win32.installer.exe",
		"firefox-55.0a1.en-US.win32.zip",
	}

	// Not every file is under that list, but it's a representative (enough) subset
	filesThatMustNotReachTheL10nFolder := []string{
		"firefox-55.0a1.en-US.linux-x86_64.common.tests.zip",
		"firefox-55.0a1.en-US.linux-i686.awsy.tests.zip",
		"firefox-55.0a1.en-US.linux-i686.common.tests.zip",
		"firefox-55.0a1.en-US.linux-i686.cppunittest.tests.zip",
		"firefox-55.0a1.en-US.linux-i686.json",
		"firefox-55.0a1.en-US.linux-i686.mochitest.tests.zip",
		"firefox-55.0a1.en-US.linux-i686.mozinfo.json",
		"firefox-55.0a1.en-US.linux-i686.reftest.tests.zip",
		"firefox-55.0a1.en-US.linux-i686.talos.tests.zip",
		"firefox-55.0a1.en-US.linux-i686.test_packages.json",
		"firefox-55.0a1.en-US.linux-i686.txt",
		"firefox-55.0a1.en-US.linux-i686.web-platform.tests.zip",
		"firefox-55.0a1.en-US.win64.xpcshell.tests.zip",
		"firefox-55.0a1.en-US.win32_info.txt",
		"firefox-55.0a1.en-US.mac.web-platform.tests.tar.gz",
		"jsshell-linux-i686.zip",
		"mozharness.zip",
	}

	for _, fileName := range filesThatMustBeDuplicated {
		filePath := "/tmp/src/" + fileName

		expectedDests := []string{
			"pub/firefox/nightly/latest-mozilla-central/build-dir/" + fileName,
			"pub/firefox/nightly/latest-mozilla-central-l10n/build-dir/" + fileName,
		}

		dests, err := rel.ToLatest(filePath)
		assert.Nil(t, err)
		assert.Equal(t, expectedDests, dests)
	}

	for _, fileName := range filesThatMustNotReachTheL10nFolder {
		filePath := "/tmp/src/" + fileName

		expectedDests := []string{
			"pub/firefox/nightly/latest-mozilla-central/build-dir/" + fileName,
		}

		dests, err := rel.ToLatest(filePath)
		assert.Nil(t, err)
		assert.Equal(t, expectedDests, dests)
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

	runTests := func(files []FileTest) {
		for _, file := range files {
			dests, err := rel.ToCandidates(file.Src)
			assert.Nil(err)
			assert.Equal(file.Dests, dests)
		}
	}

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/candidates/32-candidates/build23/build-dir/subdir/file"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"pub/firefox/candidates/32-candidates/build23/unsigned/build-dir/subdir/win32-file"}},
		FileTest{"/tmp/src/host/bin/mar.exe", []string{"pub/firefox/candidates/32-candidates/build23/mar-tools/win32/mar.exe"}},
		FileTest{"/tmp/src/mar.exe", []string{"pub/firefox/candidates/32-candidates/build23/mar-tools/win32/mar.exe"}},
	}

	runTests(files)

	rel.Signed = true

	files = []FileTest{
		FileTest{"/tmp/src/subdir/win32-file", []string{"pub/firefox/candidates/32-candidates/build23/build-dir/subdir/win32-file"}},
	}

	runTests(files)
}

func TestReleaseToMobileCandidates(t *testing.T) {
	assert := assert.New(t)
	rel := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildNumber = "23"
	rel.Version = "32"
	rel.Signed = false

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"pub/firefox/candidates/32-candidates/build23/build-dir/subdir/file"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"pub/firefox/candidates/32-candidates/build23/build-dir/subdir/win32-file"}},
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
