package main

import "testing"

func TestAll(t *testing.T) {
	const dbpath = "foo.db"

	db := InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	items := []SMD{
		SMD{Code: "1", Device: "213", Description: "decfor213"},
		SMD{Code: "2", Device: "214", Description: "decfor214"},
	}
	StoreItem(db, items)

	readItems := ReadItem(db)
	t.Log(readItems)

	items2 := []SMD{
		SMD{Code: "1", Device: "215", Description: "-"},
		SMD{Code: "3", Device: "216", Description: "decfor214"},
	}
	StoreItem(db, items2)

	readItems2 := ReadItem(db)
	t.Log(readItems2)
}
