package crusher

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// This default value of blacklist is for testing purposes. It will be overwritten by the value given in Run on line 19 for the actual execution of the program
var blacklist = ",blacklists,districts,"
var db *sql.DB

// Run kicks off the logic of the program and is the main abstraction point from the other third party libraries. It is given a command to execute, a path to a SQL file, a boolean that represents whether the resulting view should be materalized or not, a string representing blacklisted relation names, and a pointer to an active database connection pool.
func Run(command string, path string, materialized bool, list string, d *sql.DB) {
	blacklist = list
	db = d
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
		q = fmt.Sprintf("CREATE VIEW %s AS %s;", name, file)
	}

	executeSQL(q)
	fmt.Println(name, "created successfully!")
}

func update(path string, materialized bool) {
	name := strings.Split(path, ".")[0]
	file := openSQLFile(path, name)
	q := ""

	if materialized {
		q = fmt.Sprintf("DROP MATERIALIZED VIEW IF EXISTS %s; CREATE MATERIALIZED VIEW %s AS %s;", name, name, file)
	} else {
		q = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s;", name, file)
	}

	executeSQL(q)
	fmt.Println(name, "updated successfully!")
}

func openSQLFile(path string, name string) string {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		msg := fmt.Sprintf("Something went wrong reading the given SQL file:\n%v", err)
		log.Fatal(msg)
	}
	file, err := validateFile(string(raw), name)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func validateFile(file string, name string) (string, error) {
	// Ensure the file name isn't in the blacklisted names
	r, err := regexp.Compile(fmt.Sprintf(",%s,", name))
	if err != nil {
		return file, fmt.Errorf("Couldn't compile RexExp for checking against blacklist!")
	}
	if r.MatchString(blacklist) == true {
		return file, fmt.Errorf("Your view name is on the blacklist - please choose another!")
	}

	// Ensure the first word in the query is 'select'
	r, err = regexp.Compile(`\Aselect\s+`)
	if err != nil {
		return file, fmt.Errorf("Couldn't compile RexExp for checking for select!")
	}
	if r.MatchString(file) != true {
		return file, fmt.Errorf("Your query needs to be a `select` statement!")
	}

	// Ensure the file has no semi-colons
	r, err = regexp.Compile(`;`)
	if err != nil {
		return file, fmt.Errorf("Couldn't compile RexExp for checking for final semi-colon!")
	}
	if r.MatchString(file) == true {
		return file, fmt.Errorf("Your query cannot contain any semi-colons!")
	}

	// Ensure the file has 0 instances of the words 'create', 'delete', 'refresh', 'update', 'insert', 'drop'
	r, err = regexp.Compile(`\s*create\s+|\s*delete\s+|\s*refresh\s+|\s*update\s+|\s*insert\s+|\s*drop\s+|\s*truncate\s+`)
	if err != nil {
		return file, fmt.Errorf("Couldn't compile RexExp checking for command words!")
	}
	if r.MatchString(file) == true {
		return file, fmt.Errorf("Your query cannot contain any of the following words:\n create - delete - refresh - update - insert - drop - truncate")
	}

	return file, nil
}

func refresh(name string) {
	q := fmt.Sprintf("REFRESH MATERIALIZED VIEW %s;", name)
	executeSQL(q)
	fmt.Println(name, "refreshed successfully!")
}

func executeSQL(query string) {
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
