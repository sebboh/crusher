package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	// Load environment variables, including DB connection info
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = db.QueryRow("select name from schools limit 1").Scan(&name)

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
		fmt.Println(name)
	}

	app.Run(os.Args)
}

func executeSQL(db *sql.DB, query *string) {
	db.Exec(*query)
}
