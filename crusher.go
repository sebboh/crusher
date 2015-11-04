package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
			fmt.Println(os.Getenv("MATERIALIZED"))
		} else {
			fmt.Println(os.Getenv("NOT_MATERIALIZED"))
		}
	}

	app.Run(os.Args)
}
