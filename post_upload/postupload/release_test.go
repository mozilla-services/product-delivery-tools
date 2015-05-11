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
		FtpPrefix: "prefix/ftp",
		PvtPrefix: "prefix/pvt",
		FtpCopier: copier,
		PvtCopier: copier,
	}, copier
}

func TestReleaseToLatest(t *testing.T) {
	assert := assert.New(t)
	rel, copier := NewTestRelease()
	rel.SourceDir = "/tmp/src"
	rel.BuildDir = "build-dir"
	rel.TinderboxBuildsDir = "tbox-win32"

	err := rel.ToLatest("/tmp/src/nobranch")
	assert.NotNil(err)

	rel.Branch = "l10n"

	err = rel.ToLatest("/tmp/src/subdir/file")
	assert.Nil(err)
	assert.Equal("prefix/ftp/latest-l10n/build-dir", copier.Dest)

	err = rel.ToLatest("/etc/passwd")
	assert.NotNil(err)

	copier.Reset()
	err = rel.ToLatest("/tmp/src/path1/path2/test.xpi")
	assert.Nil(err)
	assert.Equal("prefix/ftp/latest-l10n/build-dir/path1/path2", copier.Dest)

	copier.Reset()
	err = rel.ToLatest("/tmp/src/mar.exe")
	assert.Nil(err)
	assert.Equal("prefix/ftp/latest-l10n/build-dir/mar-tools/win32", copier.Dest)
}
