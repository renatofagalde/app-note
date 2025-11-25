package db

import "gorm.io/gorm"

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Database interface {
	Gorm() *gorm.DB
}
