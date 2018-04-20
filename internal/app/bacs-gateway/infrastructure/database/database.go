package database

import (
	"fmt"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/configuration"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"strings"
	_ "github.com/lib/pq"
)

type DatabaseConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (c DatabaseConfiguration) toConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

func (c DatabaseConfiguration) address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

var db *pg.DB

func Connection() *pg.DB {
	return db
}

func StartDatabase() {
	databaseConfiguration := DatabaseConfiguration{
		Port:     configuration.GetInt("db.port"),
		User:     configuration.Get("db.user"),
		Database: configuration.Get("db.database"),
		Host:     configuration.Get("db.host"),
		Password: configuration.Get("db.password"),
	}
	connectionString := databaseConfiguration.toConnectionString()
	fmt.Println("connectionString:", connectionString)

	migrateDatabase(connectionString)

	orm.SetTableNameInflector(func(s string) string {
		return strings.Title(s)
	})

	db = pg.Connect(&pg.Options{
		User:     databaseConfiguration.User,
		Password: databaseConfiguration.Password,
		Database: databaseConfiguration.Database,
		Addr:     databaseConfiguration.address(),
	})

	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
}
