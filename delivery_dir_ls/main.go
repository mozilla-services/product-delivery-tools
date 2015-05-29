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
	lister := func(suffix, prefix string) http.Handler {
		return bucketlister.New(
			c.String("bucket-prefix")+"-"+suffix, prefix, deliverytools.AWSConfig)
	}

	http.Handle("/firefox/", lister("firefox", "/firefox/"))

	err := http.ListenAndServe(c.String("addr"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
