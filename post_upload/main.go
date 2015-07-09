package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mozilla-services/product-delivery-tools"
	"github.com/mozilla-services/product-delivery-tools/post_upload/postupload"
)

func main() {
	app := cli.NewApp()
	app.Name = "post_upload"
	app.HideVersion = true
	app.Version = deliverytools.Version
	app.Usage = "post_upload [options] <directory> <file> [file]..."
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jeremy Orem",
			Email: "oremj@mozilla.com",
		},
	}
	app.Action = doMain
	app.Flags = Flags

	app.RunAndExitOnError()
}

func contextToOptions(c *cli.Context, r *postupload.Release) {
	r.Branch = c.String("branch")
	r.BuildDir = c.String("builddir")
	r.BuildID = c.Generic("buildid").(*postupload.BuildID)
	r.BuildNumber = c.String("build-number")
	r.NightlyDir = c.String("nightly-dir")
	r.Revision = c.String("revision")
	r.ShortDir = !c.Bool("no-shortdir")
	r.Signed = c.Bool("signed")
	r.SubDir = c.String("subdir")
	r.TinderboxBuildsDir = c.String("tinderbox-builds-dir")
	r.Version = c.String("version")
	r.Who = c.String("who")
}

func eachFile(files []string, f func(string) ([]string, error)) {
	for _, file := range files {
		_, err := f(file)
		if err != nil {
			log.Println(err)
		}
	}
}

type pathFunc func(string) ([]string, error)

func doMain(c *cli.Context) {
	errs := []error{}
	requireArgs := func(args ...string) (hasErrors bool) {
		for _, arg := range args {
			if c.String(arg) == "" {
				hasErrors = true
				errs = append(errs, fmt.Errorf("--%s must be set", arg))
			}
		}
		return !hasErrors
	}

	pathActions := []pathFunc{}
	boolRequireArgs := func(f pathFunc, boolArg string, args ...string) bool {
		if c.Bool(boolArg) && requireArgs(args...) {
			pathActions = append(pathActions, f)
			return true
		}
		return false
	}

	if len(c.Args()) < 2 {
		log.Println("you must specify a directory and at least one file")
		os.Exit(1)
	}

	uploadDir := c.Args()[0]
	files := c.Args()[1:]

	requireArgs("product", "bucket-prefix")

	release := postupload.NewRelease(uploadDir, c.String("product"))

	boolRequireArgs(release.ToLatest, "release-to-latest", "branch")
	boolRequireArgs(release.ToDated, "release-to-dated", "branch", "buildid", "nightly-dir")
	boolRequireArgs(release.ToCandidates, "release-to-candidates-dir", "version", "build-number")
	boolRequireArgs(release.ToMobileCandidates, "release-to-mobile-candidates-dir", "version", "build-number", "builddir")
	boolRequireArgs(release.ToTinderboxBuilds, "release-to-tinderbox-builds", "tinderbox-builds-dir")
	boolRequireArgs(release.ToDatedTinderboxBuilds, "release-to-dated-tinderbox-builds", "tinderbox-builds-dir", "buildid")
	boolRequireArgs(release.ToTryBuilds, "release-to-try-builds", "who", "revision", "builddir")

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println("Error:", err)
		}
		os.Exit(1)
	}

	contextToOptions(c, release)

	bucketPrefix := c.String("bucket-prefix")
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			log.Fatalf("Error: %s does not exist.\n", f)
		}
	}

	for _, file := range files {
		dests := []string{}
		for _, action := range pathActions {
			actionDests, err := action(file)
			if err != nil {
				log.Printf("file: %s, err: %s", file, err)
				continue
			}
			dests = append(dests, actionDests...)
		}

		for _, dest := range dests {
			bucket := bucketPrefix + "-" + destToBucket(dest)
			url := c.String("url-prefix") + dest
			if c.Bool("dry-run") {
				fmt.Printf("%s -> %s:%s\n", file, bucket, dest)
				fmt.Fprintln(os.Stderr, url)
				continue
			}
			if err := s3CopyFile(file, bucket, dest); err != nil {
				log.Println(err)
			} else {
				fmt.Fprintln(os.Stderr, url)
			}
		}
	}
}

func destToBucket(dest string) string {
	for _, pathMount := range deliverytools.ProdBucketMap.Mounts {
		if strings.HasPrefix(dest, pathMount.Prefix) {
			return pathMount.Bucket
		}
	}

	return deliverytools.ProdBucketMap.Default
}
