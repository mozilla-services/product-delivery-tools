package main

import (
	"log"
	"net/http"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mozilla-services/product-delivery-tools"
	"github.com/mozilla-services/product-delivery-tools/bucketlister/services"
	"github.com/mozilla-services/product-delivery-tools/metrics"
	"github.com/mozilla-services/product-delivery-tools/mozlog"
)

func main() {
	app := cli.NewApp()
	app.Name = "bucketlister"
	app.HideVersion = true
	app.Version = deliverytools.Version
	app.Usage = "bucketlister [options]"
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
	mozlog.UseMozLogger(c.String("logger"))
	if c.String("dogstatsd-ip") != "" {
		metrics.Metric = &metrics.GodSpeed{
			NameSpace: c.String("dogstatsd-namespace"),
			IP:        c.String("dogstatsd-ip"),
			Port:      c.Int("dogstatsd-port"),
		}
	}
	rootLister := services.NewBucketLister(
		c.String("bucket-prefix")+"-"+deliverytools.ProdBucketMap.Default,
		"",
		deliverytools.AWSConfig,
	)

	lister := func(suffix, prefix string) http.Handler {
		bl := services.NewBucketLister(
			c.String("bucket-prefix")+"-"+suffix, prefix, deliverytools.AWSConfig)

		rootLister.AddBucketLister("/", bl)
		return bl
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rootLister.ServeHTTP(w, r)
		duration := time.Now().Sub(startTime)
		go metrics.Metric.Set("pageload", float64(duration/time.Nanosecond), []string{})
	}))

	for _, mount := range deliverytools.ProdBucketMap.Mounts {
		http.Handle("/"+mount.Prefix, lister(mount.Bucket, "/"+mount.Prefix))
	}

	err := http.ListenAndServe(c.String("addr"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
