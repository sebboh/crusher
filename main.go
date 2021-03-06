package main

import (
	"database/sql"
	"github.com/codegangsta/cli"
	"github.com/devonestes/crusher/crusher"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var blacklist string
var dbURL string

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
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	app.Action = func(c *cli.Context) {
		command := c.Args()[0]
		path := c.Args()[1]
		materialized := c.Bool("m")
		crusher.Run(command, path, materialized, blacklist, db)
	}
	app.Run(os.Args)
}
