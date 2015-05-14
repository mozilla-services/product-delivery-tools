package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mozilla-services/product-delivery-tools/post_upload/postupload"
)

func main() {
	app := cli.NewApp()
	app.Name = "post_upload"
	app.HideVersion = true
	app.Version = Version
	app.Usage = ""
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jeremy Orem",
			Email: "oremj@mozilla.com",
		},
	}
	app.Action = doMain
	app.Flags = Flags

	app.Run(os.Args)
}

func contextToOptions(c *cli.Context, r *postupload.Release) error {
	r.Branch = c.String("branch")
	r.BuildDir = c.String("builddir")

	if c.String("buildid") == "" {
		r.BuildID = nil
	} else {
		buildID, err := postupload.NewBuildID(c.String("buildid"))
		if err != nil {
			return fmt.Errorf("error in BuildID, %s", err)
		}
		r.BuildID = buildID
	}
	r.BuildNumber = c.String("build-number")
	r.NightlyDir = c.String("nightly-dir")
	r.Revision = c.String("revision")
	r.ShortDir = !c.Bool("no-shortdir")
	r.Signed = c.Bool("signed")
	r.SubDir = c.String("subdir")
	r.TinderboxBuildsDir = c.String("tinderbox-builds-dir")
	r.Version = c.String("version")
	r.Who = c.String("who")

	return nil
}

func eachFile(files []string, f func(string) ([]string, error)) {
	for _, file := range files {
		_, err := f(file)
		if err != nil {
			log.Println(err)
		}
	}
}

func doMain(c *cli.Context) {
	errs := []error{}
	requireArgs := func(args ...string) (hasErrors bool) {
		for _, arg := range args {
			if c.String(arg) == "" {
				hasErrors = true
				errs = append(errs, fmt.Errorf("--%s must be set", arg))
			}
		}
		return
	}

	boolRequireArgs := func(boolArg string, args ...string) bool {
		if c.Bool(boolArg) {
			return requireArgs(args...)
		}
		return false
	}

	if len(c.Args()) < 2 {
		errs = append(errs, errors.New("you must specify a directory and at least one file"))
	}

	requireArgs("product", "bucket")
	boolRequireArgs("release-to-latest", "branch")
	boolRequireArgs("release-to-dated", "branch", "buildid", "nightly-dir")
	boolRequireArgs("release-to-candidates-dir", "version", "build_number")
	boolRequireArgs("release-to-mobile-candidates-dir", "version", "build-number", "builddir")
	boolRequireArgs("release-to-tinderbox-builds", "tinderbox-builds-dir")
	boolRequireArgs("release-to-dated-tinderbox-builds", "tinderbox-builds-dir", "buildid")
	boolRequireArgs("release-to-try-builds", "who", "revision", "builddir")

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println("Error:", err)
		}
		os.Exit(1)
	}

	uploadDir := c.Args()[0]
	files := c.Args()[1:]

	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			log.Fatalf("Error: %s does not exist.\n", f)
		}
	}

	release := postupload.NewRelease(uploadDir, c.String("product"))
	err := contextToOptions(c, release)
	if err != nil {
		log.Fatal("Error parsing options:", err)
	}

	if c.Bool("release-to-latest") {
		eachFile(files, release.ToLatest)
	}
	if c.Bool("release-to-dated") {
		eachFile(files, release.ToDated)
	}

	if c.Bool("release-to-candidates-dir") {
		eachFile(files, release.ToCandidates)
	}

	if c.Bool("release-to-mobile-candidates-dir") {
		eachFile(files, release.ToMobileCandidates)
	}

	if c.Bool("releaset-to-tinderbox-builds") {
		eachFile(files, release.ToTinderboxBuilds)
	}

	if c.Bool("release-to-dated-tinderbox-builds") {
		eachFile(files, release.ToDatedTinderboxBuilds)
	}

	if c.Bool("release-to-try-builds") {
		eachFile(files, release.ToTryBuilds)
	}
}
