package main

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		device TEXT,
		description TEXT,
		insertedDatetime DATETIME
	);
	`

	if _, err := db.Exec(sql_table); err != nil {
		panic(err)
	}
}

func StoreItem(db *sql.DB, items []SMD) {
	sql_additem := `
	INSERT OR REPLACE INTO items(
		code,
		device,
		description,
		InsertedDatetime
	) values(?, ?, ?, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err2 := stmt.Exec(item.Code, item.Device, item.Description); err2 != nil {
			panic(err2)
		}
	}
}

func ReadItem(db *sql.DB) []SMD {
	sql_readall := `
	SELECT code, device, description FROM items
	ORDER BY datetime(insertedDatetime) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []SMD
	for rows.Next() {
		item := SMD{}
		if err2 := rows.Scan(&item.Code, &item.Device, &item.Description); err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}
