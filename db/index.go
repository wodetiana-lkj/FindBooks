package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string // 可选，针对需要时区设置的情况
}

var config = &Config{
	Host:     os.Getenv("POSTGRESQL_IP"),
	Port:     5432,
	User:     os.Getenv("POSTGRESQL_NAME"),
	Password: os.Getenv("POSTGRESQL_PASSWORD"),
	DBName:   "scrapy",
	SSLMode:  "disable",
	TimeZone: "Asia/Shanghai", // 可根据需要调整
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.TimeZone)
}

func Connect() *gorm.DB {
	dsn := config.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect database", err)
		return nil
	}
	return db
}
