package crusher

import (
	"errors"
	"testing"
)

const wordErr = "Your query cannot contain any of the following words:\n create - delete - refresh - update - insert - drop - truncate"

var tests = []struct {
	query string
	name  string
	err   error
}{
	// Test blacklist exclusion
	{"select * from districts", "blacklists", errors.New("Your view name is on the blacklist - please choose another!")},
	// Test that the word `create` is nowhere in the query
	{"select * from (create table districts as select * from districts_view)", "a_districts", errors.New(wordErr)},
	// Test that the word `delete` is nowhere in the query
	{"select * from (delete table districts as select * from districts_view)", "a_districts", errors.New(wordErr)},
	// Test that the word `refresh` is nowhere in the query
	{"select * from (refresh table districts as select * from districts_view)", "a_districts", errors.New(wordErr)},
	// Test that the word `update` is nowhere in the query
	{"select * from (update table districts as select * from districts_view)", "a_districts", errors.New(wordErr)},
	// Test that the word `insert` is nowhere in the query
	{"select * from (insert table districts as select * from districts_view)", "a_districts", errors.New(wordErr)},
	// Test that the word `drop` is nowhere in the query
	{"select * from (drop table districts)", "a_districts", errors.New(wordErr)},
	// Test that the word `truncate` is nowhere in the query
	{"select * from (truncate districts)", "a_districts", errors.New(wordErr)},
}

func TestValidateFile(t *testing.T) {
	query := "select * from districts"
	name := "a_districts"
	_, err := validateFile(query, name)
	if err != nil {
		t.Errorf("\nquery = `%v`\nname = `%v`\ngot error `%v`\nwanted `<nil>`", query, name, err)
	}

	for _, test := range tests {
		_, err := validateFile(test.query, test.name)
		if err == nil {
			t.Errorf("Got `<nil>` error\nWanted `%v`", test.err)
		} else if err.Error() != test.err.Error() {
			t.Errorf("\nquery = `%v`\nname = `%v`\ngot error `%v`\nwant error `%v`", test.query, test.name, err, test.err)
		}
	}
}
