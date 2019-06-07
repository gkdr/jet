package tests

import (
	"fmt"
	. "github.com/go-jet/jet/sqlbuilder"
	"github.com/go-jet/jet/tests/.test_files/dvd_rental/test_sample/model"
	"github.com/go-jet/jet/tests/.test_files/dvd_rental/test_sample/table"
	"gotest.tools/assert"
	"testing"
)

func TestUpdateValues(t *testing.T) {
	_, err := table.Link.INSERT(table.Link.URL, table.Link.Name, table.Link.Rel).
		VALUES("http://www.postgresqltutorial.com", "PostgreSQL Tutorial", DEFAULT).
		VALUES("http://www.yahoo.com", "Yahoo", DEFAULT).
		VALUES("http://www.bing.com", "Bing", DEFAULT).
		RETURNING(table.Link.ID).Execute(db)

	assert.NilError(t, err)

	query := table.Link.
		UPDATE(table.Link.Name, table.Link.URL).
		SET("Bong", "http://bong.com").
		WHERE(table.Link.Name.EQ(String("Bing")))

	queryStr, args, err := query.Sql()

	assert.NilError(t, err)
	assert.Equal(t, len(args), 3)
	fmt.Println(queryStr)

	res, err := query.Execute(db)

	assert.NilError(t, err)

	fmt.Println(res)

	links := []model.Link{}

	err = table.Link.SELECT(table.Link.AllColumns).
		WHERE(table.Link.Name.EQ(String("Bong"))).
		Query(db, &links)

	assert.NilError(t, err)

	//spew.Dump(links)
}

func TestUpdateAndReturning(t *testing.T) {
	_, err := table.Link.INSERT(table.Link.URL, table.Link.Name, table.Link.Rel).
		VALUES("http://www.postgresqltutorial.com", "PostgreSQL Tutorial", DEFAULT).
		VALUES("http://www.ask.com", "Ask", DEFAULT).
		VALUES("http://www.ask.com", "Ask", DEFAULT).
		VALUES("http://www.yahoo.com", "Yahoo", DEFAULT).
		VALUES("http://www.bing.com", "Bing", DEFAULT).
		RETURNING(table.Link.ID).Execute(db)

	assert.NilError(t, err)

	stmt := table.Link.
		UPDATE(table.Link.Name, table.Link.URL).
		SET("DuckDuckGo", "http://www.duckduckgo.com").
		WHERE(table.Link.Name.EQ(String("Ask"))).
		RETURNING(table.Link.AllColumns)

	stmtStr, args, err := stmt.Sql()

	assert.NilError(t, err)
	assert.Equal(t, len(args), 3)
	fmt.Println(stmtStr)

	links := []model.Link{}

	err = stmt.Query(db, &links)

	assert.NilError(t, err)

	assert.Equal(t, len(links), 2)

	assert.Equal(t, links[0].Name, "DuckDuckGo")

	assert.Equal(t, links[1].Name, "DuckDuckGo")
}
