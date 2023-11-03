package db

import (
	"database/sql"
	"fmt"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var dbUrl = "file:./db.sqlite3"
var Db *sql.DB

type Dog struct {
	Id           int
	Name         string
	OwnerName    string
	Address      string
	WalksPerWeek int
	PricePerWalk int
}

func Init() error {
	err := connect()
	if err != nil {
		return err
	}

	err = createTables()
	if err != nil {
		return err
	}

	return nil
}

func connect() error {
	var err error
	Db, err = sql.Open("libsql", dbUrl)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	return nil
}

func createTables() error {
	_, err := Db.Exec(`
        CREATE TABLE IF NOT EXISTS dogs (
            id INTEGER PRIMARY KEY,
            name TEXT,
            ownerName TEXT,
            address TEXT,
            walksPerWeek INTEGER,
            pricePerWalk INTEGER
        );
    `)
	if err != nil {
		return fmt.Errorf("error creating dogs table: %v", err)
	}

	return nil
}

func GetDogs() ([]Dog, error) {
	rows, err := Db.Query("SELECT * FROM dogs")
	if err != nil {
		return nil, fmt.Errorf("error getting dogs: %v", err)
	}

	var dogs []Dog
	for rows.Next() {
		var dog Dog
		err := rows.Scan(
			&dog.Id,
			&dog.Name,
			&dog.OwnerName,
			&dog.Address,
			&dog.WalksPerWeek,
			&dog.PricePerWalk,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning dog: %v", err)
		}
		dogs = append(dogs, dog)
	}

	return dogs, nil
}

func AddDog(dog Dog) error {
	_, err := Db.Exec(`
        INSERT INTO dogs (
            name,
            ownerName,
            address,
            walksPerWeek,
            pricePerWalk
        ) VALUES (?, ?, ?, ?, ?)
    `, dog.Name, dog.OwnerName, dog.Address, dog.WalksPerWeek, dog.PricePerWalk)
	if err != nil {
		return fmt.Errorf("error adding dog: %v", err)
	}

	return nil
}