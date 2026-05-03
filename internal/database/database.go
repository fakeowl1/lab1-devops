package database

import (
	"fmt"
	"log"
	"notes-service/internal/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDatabase struct {
	db *gorm.DB
}

func New(host, user, port, password, dbName string) (*GormDatabase, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
			user, password, host, port, dbName, 
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
  	logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
  	},
	)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Note{})
	if (err != nil) {
		log.Fatalf("migrate error: %v", err)
	}
	return &GormDatabase{db: db}, nil
}

func (d *GormDatabase) Close() {
	sqlDB, _ := d.db.DB()
	err := sqlDB.Close()
	if (err != nil) {
		log.Fatalf("error while closing database: %v", err)
	}
}

func (d* GormDatabase) Ping() error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database connection is not initialized")
  }
	sqlDB, err := d.db.DB()
	if (err != nil) {
		return err
	}

	return sqlDB.Ping()
}
