package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/mozilla-services/product-delivery-tools"
	"github.com/mozilla-services/product-delivery-tools/service/bucketlister"
)

func main() {
	app := cli.NewApp()
	app.Name = "delivery_dir_ls"
	app.HideVersion = true
	app.Version = deliverytools.Version
	app.Usage = "delivery_dir_ls [options]"
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

func doMain(c *cli.Context) {
	rootLister := &bucketlister.RootLister{}
	lister := func(suffix, prefix string) http.Handler {
		bl := bucketlister.New(
			c.String("bucket-prefix")+"-"+suffix, prefix, deliverytools.AWSConfig)

		rootLister.AddBucketLister(bl)
		return bl
	}

	http.Handle("/", rootLister)
	http.Handle("/firefox/", lister("firefox", "/firefox/"))
	http.Handle("/firefox/try-builds/", lister("firefox-try", "/firefox/try-builds/"))
	http.Handle("/mobile/", lister("firefox-android", "/mobile/"))
	http.Handle("/mobile/try-builds/", lister("firefox-android-try", "/mobile/try-builds/"))
	http.Handle("/opus/", lister("opus", "/opus/"))
	http.Handle("/thunderbird/", lister("thunderbird", "/thunderbird/"))
	http.Handle("/thunderbird/try-builds/", lister("thunderbird-try", "/thunderbird/try-builds/"))
	http.Handle("/xulrunner/", lister("xulrunner", "/xulrunner/"))
	http.Handle("/xulrunner/try-builds/", lister("xulrunner-try", "/xulrunner/try-builds/"))

	err := http.ListenAndServe(c.String("addr"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
