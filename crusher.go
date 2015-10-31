package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "crusher"
	app.Usage = "manage production database views without having to go through engineering"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "materialized, m",
			Usage: "apply if the view should be materialized",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("m") {
			fmt.Println("Materialized")
		} else {
			fmt.Println("Not materialized")
		}
	}

	app.Run(os.Args)
}
