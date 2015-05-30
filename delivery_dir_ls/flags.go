package main

import "github.com/codegangsta/cli"

var Flags = []cli.Flag{
	cli.StringFlag{Name: "addr", Usage: "Set the address on which to listen", Value: ":8888"},
	cli.StringFlag{
		Name:  "bucket-prefix",
		Value: "net-mozaws-prod-delivery",
		Usage: "Sets S3 bucket prefix"},
}
