package dbase

import (
	"testing"
	"os"
	"github.com/gocarina/gocsv"
)

func TestAll(t *testing.T) {
	const dbpath = "foo.db"

	db := InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	items := []*SMD{
		&SMD{Code: "1", Device: "213", Description: "decfor213"},
		&SMD{Code: "2", Device: "214", Description: "decfor214"},
	}
	StoreItem(db, items)

	readItems := ReadItemsAll(db)
	t.Log(readItems)

	items2 := []*SMD{
		&SMD{Code: "1", Device: "215", Description: "-"},
		&SMD{Code: "3", Device: "216", Description: "decfor214"},
	}
	StoreItem(db, items2)

	readItems2 := ReadItemsAll(db)
	t.Log(readItems2)
}

func TestFilltmpBase(t *testing.T) {
	const dbpath = "smd.db"

	db := InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	smddataFile, err := os.OpenFile("test.csv", os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer smddataFile.Close()

	smdlist := []*SMD{}
	if err := gocsv.UnmarshalFile(smddataFile, &smdlist); err != nil { // Load clients from file
		panic(err)
	}

	StoreItem(db, smdlist)

	readItems := ReadItemsAll(db)
	t.Log(readItems)
}

func TestFillBase(t *testing.T) {
	const dbpath = "smd.db"

	db := InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	smddataFile, err := os.OpenFile("databook.csv", os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer smddataFile.Close()

	smdlist := []*SMD{}
	if err := gocsv.UnmarshalFile(smddataFile, &smdlist); err != nil { // Load clients from file
		panic(err)
	}

	StoreItem(db, smdlist)

	readItems := ReadItemsAll(db)
	t.Log(len(readItems))
}
