package dbase

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type SMD struct {
	// Our example struct, you can use "-" to ignore a field
	Code        string `csv:"code"`              // code
	A           string `csv:"-"`                 // A
	Device      string `csv:"Type"`              // Type
	Function    string `csv:"Function"`          // Function
	Description string `csv:"Short description"` // Short description
	Case        string `csv:"Case"`              // Case
	Mnf         string `csv:"Mnf"`               // Mnf
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS items(
		code TEXT,
		device TEXT,
		function TEXT,
		description TEXT,
		insertedDatetime DATETIME
	);
	`

	if _, err := db.Exec(sql_table); err != nil {
		panic(err)
	}
}

func StoreItem(db *sql.DB, items []*SMD) {
	sql_additem := `
	INSERT OR REPLACE INTO items(
		code,
		device,
		function,
		description,
		InsertedDatetime
	) values(?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err2 := stmt.Exec(item.Code, item.Device, item.Function, item.Description); err2 != nil {
			panic(err2)
		}
	}
}

func ReadItemsAll(db *sql.DB) []SMD {
	sql := `
	SELECT code, device, function, description FROM items
	ORDER BY datetime(insertedDatetime) DESC
	`
	return readDB(db, sql)
}

func ReadItemsByCode(db *sql.DB, code string) []SMD {
	sql := `
	SELECT code, device, function, description FROM items
	WHERE code = "` + code + `"
	ORDER BY datetime(insertedDatetime) DESC
	`
	return readDB(db, sql)
}

func readDB(db *sql.DB, sql string) []SMD {

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []SMD
	for rows.Next() {
		item := SMD{}
		if err2 := rows.Scan(&item.Code, &item.Device, &item.Function, &item.Description); err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}
