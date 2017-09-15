package db

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/sqlite3"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/postgres"
	"gopkg.in/doug-martin/goqu.v4"
	"io/ioutil"
	"bitbucket.org/codefreak/hsmpp/smpp/stringutils"
	"strings"
	"github.com/pkg/errors"
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	database *goqu.Database
	driver string
)

const (
	validationQuery = "select MIN(id) from Item"
	SQLFilesPath           = "./sqls"
)

//CheckAndCreateDB Checks if database exists, if not, creates one with basic tables, admin user and indexes
func CheckAndCreateDB() (*goqu.Database, error) {
	var err error
	if !exists(database) {
		err = create(database)
		if err != nil {
			err = errors.Wrap(err, "Couldn't create database.")
		}
	}
	return database, err
}

func Connect(dbdriver, dsn string) (*goqu.Database, error) {
	driver = dbdriver
	con, err := sql.Open(driver, dsn)
	if err != nil {
		return database, err
	}
	err = con.Ping()
	if err != nil {
		return database, err
	}
	database = goqu.New(driver, con)
	database.Logger(log.StandardLogger())
	return database, nil
}

func Get() *goqu.Database {
	return database
}

// ConnectMock makes a mock database connection for testing purposes
func ConnectMock(driver string, t *testing.T) (*goqu.Database, sqlmock.Sqlmock, error) {
	driver = driver
	con, mock, err := sqlmock.New()
	if err != nil {
		err = errors.Wrap(err, "Couldn't create mock driver.")
	}
	database = goqu.New(driver, con)
	return database, mock, err
}

// exists checks if a database exists
func exists(db *goqu.Database) bool {
	_, err := db.Exec(validationQuery)
	if err != nil {
		return false
	}
	return true
}

// create creates a fresh database, tables, indexes and populates primary data
func create(db *goqu.Database) error {
	fileName := SQLFilesPath + "/fresh-" + driver + ".sql"
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrapf(err, "Couldn't load file %s", fileName)
	}
	query := stringutils.ByteToString(b)
	replacer := strings.NewReplacer("\n", "", "\r", "")
	_, err = db.Exec(replacer.Replace(query))
	if err != nil {
		err = errors.Wrapf(err,"Couldn't load file %s", fileName)
		return err
	}
	return err
}

