package main

import (
	"os"

	"github.com/codegangsta/cli"
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

func doMain(c *cli.Context) {

}
