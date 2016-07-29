package main

import (
	"github.com/codegangsta/cli"
	"github.com/mozilla-services/product-delivery-tools/post_upload/postupload"
)

func init() {
	cli.AppHelpTemplate = `USAGE:
   {{.Usage}}
VERSION:
   {{.Version}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}{{end}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

// Flags for post_upload
var Flags = []cli.Flag{
	cli.StringFlag{Name: "product, p", Usage: "Set product name to build paths properly."},
	cli.StringFlag{Name: "version, v", Usage: "Set version number to build paths properly."},
	cli.StringFlag{
		Name:  "bucket-prefix",
		Value: "net-mozaws-prod-delivery",
		Usage: "Sets S3 bucket prefix"},
	cli.StringFlag{
		Name:  "url-prefix",
		Value: "https://archive.mozilla.org/",
		Usage: "Sets URL prefix. (Only affects output)"},
	cli.StringFlag{
		Name: "nightly-dir", Value: "nightly",
		Usage: "Set the base directory for nightlies (ie $product/$nightly_dir/}, and the parent directory for release candidates (default 'nightly'}."},
	cli.StringFlag{Name: "branch, b", Usage: "Set branch name to build paths properly."},
	cli.GenericFlag{Name: "buildid, i", Value: new(postupload.BuildID), Usage: "Set buildid to build paths properly."},
	cli.StringFlag{Name: "build-number, n", Usage: "Set buildid to build paths properly."},
	cli.StringFlag{Name: "revision, r"},
	cli.StringFlag{Name: "who, w"},
	cli.BoolFlag{Name: "no-shortdir, S", Usage: "Don't symlink the short dated directories."}, //bool
	cli.StringFlag{Name: "builddir", Usage: "Subdir to arrange packaged unittest build paths properly."},
	cli.StringFlag{Name: "subdir", Usage: "Subdir to arrange packaged unittest build paths properly."},
	cli.StringFlag{Name: "tinderbox-builds-dir", Usage: "Set tinderbox builds dir to build paths properly."},
	cli.BoolFlag{Name: "release-to-latest, l", Usage: "Copy files to $product/$nightly_dir/latest-$branch"},
	cli.BoolFlag{Name: "release-to-dated, d", Usage: "Copy files to $product/$nightly_dir/$datedir-$branch"},
	cli.BoolFlag{Name: "release-to-candidates-dir, c", Usage: "Copy files to $product/$nightly_dir/$version-candidates/build$build_number"},
	cli.BoolFlag{Name: "release-to-mobile-candidates-dir", Usage: "Copy mobile files to $product/$nightly_dir/$version-candidates/build$build_number/$platform"},
	cli.BoolFlag{Name: "release-to-tinderbox-builds, t", Usage: "Copy files to $product/tinderbox-builds/$tinderbox_builds_dir"},
	cli.BoolFlag{Name: "release-to-latest-tinderbox-builds", Usage: "Softlink tinderbox_builds_dir to latest"},
	cli.BoolFlag{Name: "release-to-tinderbox-dated-builds", Usage: "Copy files to $product/tinderbox-builds/$tinderbox_builds_dir/$timestamp"},
	cli.BoolFlag{Name: "release-to-try-builds", Usage: "Copy files to try-builds/$who-$revision"},
	cli.BoolFlag{Name: "signed", Usage: "Don't use unsigned directory for uploaded files"},
	cli.BoolFlag{Name: "dry-run", Usage: "Print the operations which would happen."},
}
