package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestPing(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	defer func () {
		_ = db.Close()
	}()

	mock.ExpectPing()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	gormDatabase := &GormDatabase{db: gormDB}

	mock.ExpectPing()
	if err := gormDatabase.Ping(); err != nil {
		t.Errorf("expected no error on ping, got %v", err)
	}
}

func TestPing_NilDatabase(t *testing.T) {
	var gormDatabase *GormDatabase = nil
	err := gormDatabase.Ping()
	if err == nil {
		t.Error("expected error when calling Ping on nil database, got nil")
	}
}

func TestClose(t *testing.T) {
	db, mock, _ := sqlmock.New()
	
	mock.ExpectPing()

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	gormDatabase := &GormDatabase{db: gormDB}

	mock.ExpectClose()
	gormDatabase.Close()
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
