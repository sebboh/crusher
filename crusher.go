package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/cli"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
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
		parseCmd(c)
	}

	app.Run(os.Args)
}

func parseCmd(c *cli.Context) {
	command := c.Args()[0]
	path := c.Args()[1]
	materialized := c.Bool("m")

	switch command {
	case "create":
		create(path, materialized)
	case "update":
		update(path, materialized)
	case "refresh":
		refresh(path)
	}
}

func create(path string, materialized bool) {
	name := strings.Split(path, ".")[0]
	file := openSQLFile(path, name)
	q := ""

	if materialized {
		q = fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS %s;", name, file)
	} else {
		q = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s;", name, file)
	}

	executeSQL(q)
	fmt.Println(name, "created successfully!")
}

func update(path string, materialized bool) {
	name := strings.Split(path, ".")[0]
	file := openSQLFile(path, name)
	q := ""

	if materialized {
		q = fmt.Sprintf("DROP MATERIALIZED VIEW IF EXISTS %s CASCADE; CREATE MATERIALIZED VIEW %s AS %s;", name, name, file)
	} else {
		q = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s;", name, file)
	}

	executeSQL(q)
	fmt.Println(name, "updated successfully!")
}

func openSQLFile(path string, name string) string {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Something went wrong reading the given SQL file")
	}

	file := string(raw)
	if validateFile(file, name) != nil {
		log.Fatal("Your SQL file does not meet the requirements to use this tool.")
	}
	return file
}

func validateFile(file string, name string) error {
	// Ensure the file name isn't in the blacklisted names
	blacklist := os.Getenv("BLACKLISTED_NAMES")
	r, err := regexp.Compile(fmt.Sprintf("`,%s,`", name))
	if err != nil {
		log.Fatal("Couldn't compile RexExp for checking against blacklist!")
	}
	if r.MatchString(blacklist) != false {
		log.Fatal("Your view name is on the blacklist - please choose another!")
	}

	// Ensure the first word in the query is 'select'
	r, err = regexp.Compile(`\Aselect\s+`)
	if err != nil {
		log.Fatal("Couldn't compile RexExp for checking for select!")
	}
	if r.MatchString(file) != true {
		log.Fatal("Your query needs to be a `select` statement!")
	}

	// Ensure the file has no semi-colons
	r, err = regexp.Compile(`;`)
	if err != nil {
		log.Fatal("Couldn't compile RexExp for checking for final semi-colon!")
	}
	if r.MatchString(file) != false {
		log.Fatal("Your query cannot contain any semi-colons!")
	}

	// Ensure the file has 0 instances of the words 'create', 'delete', 'refresh', 'update', 'insert', 'drop'
	r, err = regexp.Compile(`\s*create\s+|\s*delete\s+|\s*refresh\s+|\s*update\s+|\s*insert\s+|\s*drop\s+`)
	if err != nil {
		log.Fatal("Couldn't compile RexExp checking for command words!")
	}
	if r.MatchString(file) != false {
		log.Fatal("Your query cannot contain any of the following words:\n create - delete - refresh - update - insert - drop")
	}

	return nil
}

func refresh(name string) {
	q := fmt.Sprintf("REFRESH MATERIALIZED VIEW CONCURRENTLY %s;", name)
	executeSQL(q)
	fmt.Println(name, "refreshed successfully!")
}

func executeSQL(query string) {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(query)
}
