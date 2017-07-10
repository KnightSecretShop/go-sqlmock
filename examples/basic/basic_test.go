package main

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// a successful case
func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = recordStats(db, 2, 3); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

// a failing test case
func TestShouldRollbackStatUpdatesOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").
		WithArgs(2, 3).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	// now we execute our method
	if err = recordStats(db, 2, 3); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetProviderByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Could not mock the database connection: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "searchable"}).
		AddRow(1, "Foobar", true)

	mock.ExpectQuery("SELECT (.+) FROM user WHERE name LIKE \\$1$").WithArgs("Foobar").WillReturnRows(rows)

	var (
		id         int64
		name       string
		searchable bool
	)

	err = db.QueryRow(`
		SELECT id, name, searchable
		FROM user
		WHERE
			name LIKE $1
	`, "Foobar").Scan(&id, &name, &searchable)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 1 || name != "Foobar" || !searchable {
		t.Fatalf("invalid row scanned")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There are unfulfilled expectation: %s", err)
	}
}
