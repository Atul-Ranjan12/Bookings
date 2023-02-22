package main

import "testing"

func TestRun(t *testing.T) {
	db, err := run()

	if err != nil {
		t.Error("failed test...")
	}
	if err := db.SQL.Ping(); err != nil {
		t.Error("Unsuccesful connection to the database. Failed test..")
	}
}
