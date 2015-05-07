package main

import "github.com/codegangsta/cli"

var Flags = []cli.Flag{
	cli.StringFlag{Name: "product, p", Usage: "Set product name to build paths properly."},
	cli.StringFlag{Name: "version, v", Usage: "Set version number to build paths properly."},
	cli.StringFlag{
		Name: "nightly-dir", Value: "nightly",
		Usage: "Set the base directory for nightlies (ie $product/$nightly_dir/}, and the parent directory for release candidates (default 'nightly'}."},
	cli.StringFlag{Name: "branch, b", Usage: "Set branch name to build paths properly."},
	cli.StringFlag{Name: "buildid, i", Usage: "Set buildid to build paths properly."},
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
}
