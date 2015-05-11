package postupload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCopier struct {
	Src    string
	Dest   string
	Copied bool
}

func (c *TestCopier) Copy(src, dest string) error {
	c.Src = src
	c.Dest = dest
	c.Copied = true
	return nil
}

func (c *TestCopier) Reset() {
	c.Src = ""
	c.Dest = ""
	c.Copied = false
}

func NewTestRelease() (*Release, *TestCopier) {
	copier := new(TestCopier)
	return &Release{
		FtpPrefix:  "prefix/ftp",
		PvtPrefix:  "prefix/pvt",
		FtpCopier:  copier,
		PvtCopier:  copier,
		SourceDir:  "/tmp/src",
		BuildDir:   "build-dir",
		Product:    "product",
		NightlyDir: "nightly",
	}, copier
}

func TestReleaseToLatest(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()
	rel.TinderboxBuildsDir = "tbox-win32"

	err := rel.ToLatest("/tmp/src/nobranch")
	assert.NotNil(err, "no Branch should trigger error.")

	rel.Branch = "l10n"

	err = rel.ToLatest("/etc/passwd")
	assert.NotNil(err, "Out of src file should trigger error")

	files := [][]string{
		[]string{"/tmp/src/subdir/file", "prefix/ftp/product/nightly/latest-l10n/build-dir"},
		[]string{"/tmp/src/path1/path2/test.xpi", "prefix/ftp/product/nightly/latest-l10n/build-dir/path1/path2"},
		[]string{"/tmp/src/mar.exe", "prefix/ftp/product/nightly/latest-l10n/build-dir/mar-tools/win32"},
	}
	for _, file := range files {
		copier.Reset()

		err = rel.ToLatest(file[0])
		assert.Nil(err)
		assert.Equal(file[1], copier.Dest)
	}
}

func TestReleaseToDated(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()

	rel.Branch = "l10n"
	rel.BuildID = "20150101223305"

	files := [][]string{
		[]string{"/tmp/src/subdir/file", "prefix/ftp/product/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir"},
		[]string{"/tmp/src/path1/path2/test.xpi", "prefix/ftp/product/nightly/2015/01/2015-01-01-22-33-05-l10n/build-dir/path1/path2"},
	}
	for _, file := range files {
		copier.Reset()

		err := rel.ToDated(file[0])
		assert.Nil(err)
		assert.Equal(file[1], copier.Dest)
	}
}
