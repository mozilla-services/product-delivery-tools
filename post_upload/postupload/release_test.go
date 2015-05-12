package postupload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCopier struct {
	Src    string
	Dest   []string
	Copied bool
}

func (c *TestCopier) Copy(src, dest string) error {
	if c.Src != "" && c.Src != src {
		panic("copier src changed without reset")
	}
	c.Src = src
	c.Dest = append(c.Dest, dest)
	return nil
}

func (c *TestCopier) Reset() {
	c.Src = ""
	c.Dest = []string{}
}

type FileTest struct {
	Src   string
	Dests []string
}

func NewTestRelease() (*Release, *TestCopier) {
	copier := new(TestCopier)
	return &Release{
		FtpPrefix:          "prefix/ftp",
		PvtPrefix:          "prefix/pvt",
		FtpCopier:          copier,
		PvtCopier:          copier,
		SourceDir:          "/tmp/src",
		BuildDir:           "build-dir",
		Product:            "product",
		NightlyDir:         "nightly",
		TinderboxBuildsDir: "tbox-win32",
	}, copier
}

func TestReleaseToLatest(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()

	err := rel.ToLatest("/tmp/src/nobranch")
	assert.NotNil(err, "no Branch should trigger error.")

	rel.Branch = "l10n"

	err = rel.ToLatest("/etc/passwd")
	assert.NotNil(err, "Out of src file should trigger error")

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"prefix/ftp/product/nightly/latest-l10n/build-dir"}},
		FileTest{"/tmp/src/path1/path2/test.xpi", []string{"prefix/ftp/product/nightly/latest-l10n/build-dir/path1/path2"}},
		FileTest{"/tmp/src/mar.exe", []string{"prefix/ftp/product/nightly/latest-l10n/build-dir/mar-tools/win32"}},
	}
	for _, file := range files {
		copier.Reset()

		assert.Nil(rel.ToLatest(file.Src))
		assert.Equal(file.Dests, copier.Dest)
	}
}

func TestReleaseToDated(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildID = "20150101223305"

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"prefix/ftp/product/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir"}},
		FileTest{"/tmp/src/path1/path2/test.xpi", []string{"prefix/ftp/product/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir/path1/path2"}},
	}
	for _, file := range files {
		copier.Reset()

		assert.Nil(rel.ToDated(file.Src))
		assert.Equal(file.Dests, copier.Dest)
	}
}

func TestReleaseToCandidates(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildNumber = "23"
	rel.Version = "32"
	rel.Signed = false

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"prefix/ftp/product/nightly/32-candidates/build23/build-dir/subdir"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"prefix/ftp/product/nightly/32-candidates/build23/unsigned/build-dir/subdir"}},
		FileTest{"/tmp/src/mar.exe", []string{"prefix/ftp/product/nightly/32-candidates/build23/mar-tools/win32"}},
	}
	for _, file := range files {
		copier.Reset()

		assert.Nil(rel.ToCandidates(file.Src))
		assert.Equal(file.Dests, copier.Dest)
	}
}

func TestReleaseToMobileCandidates(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildNumber = "23"
	rel.Version = "32"
	rel.Signed = false

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"prefix/ftp/product/nightly/32-candidates/build23/build-dir/subdir"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"prefix/ftp/product/nightly/32-candidates/build23/build-dir/subdir"}},
		FileTest{"/tmp/src/mar.exe", []string{"prefix/ftp/product/nightly/32-candidates/build23/build-dir"}},
	}
	for _, file := range files {
		copier.Reset()

		assert.Nil(rel.ToMobileCandidates(file.Src))
		assert.Equal(file.Dests, copier.Dest, "src: %s", file.Src)
	}
}

func TestReleaseToTryBuilds(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildNumber = "23"
	rel.Version = "32"
	rel.Signed = false
	rel.Who = "testuser"
	rel.Revision = "r33"

	files := []FileTest{
		FileTest{"/tmp/src/subdir/file", []string{"prefix/ftp/product/try-builds/testuser-r33/build-dir"}},
		FileTest{"/tmp/src/subdir/win32-file", []string{"prefix/ftp/product/try-builds/testuser-r33/build-dir"}},
		FileTest{"/tmp/src/mar.exe", []string{"prefix/ftp/product/try-builds/testuser-r33/build-dir"}},
	}
	for _, file := range files {
		copier.Reset()

		assert.Nil(rel.ToTryBuilds(file.Src))
		assert.Equal(file.Dests, copier.Dest, "src: %s", file.Src)
	}
}
