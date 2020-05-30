package dao

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/devfabric/dao/config"
	"github.com/jmoiron/sqlx"
)

const (
	DB_TYPE_UNDEFINED = "undefined"
	DB_TYPE_MYSQL     = "mysql"
	DB_TYPE_POSTGRES  = "postgres"
)

var (
	dbURLRegex = regexp.MustCompile("(Datasource:\\s*)?(\\S+):(\\S+)@|(Datasource:.*\\s)?(user=\\S+).*\\s(password=\\S+)|(Datasource:.*\\s)?(password=\\S+).*\\s(user=\\S+)")
)

type SqlxDB struct {
	DB              *sqlx.DB
	IsDBInitialized bool
	DatabaseType    string
}

func NewDB(dir string) (*SqlxDB, error) {
	dbConfig, err := config.LoadDBConfig(dir)
	if err != nil {
		return nil, err
	}

	ds := dbConfig.DataSource

	var dbIns *SqlxDB
	switch dbConfig.Type {
	case DB_TYPE_MYSQL:
		dbIns, err = NewUserRegistryMySQL(ds, dbConfig)
		if err != nil {
			return nil, fmt.Errorf("Failed to create MySQL Object %s", err.Error())
		}
		dbIns.DatabaseType = DB_TYPE_MYSQL
		dbIns.IsDBInitialized = true
	case DB_TYPE_POSTGRES:
		//todo create pg ins
		dbIns = &SqlxDB{}
		dbIns.DatabaseType = DB_TYPE_POSTGRES
		dbIns.IsDBInitialized = true
	default:
		return nil, fmt.Errorf("Invalid db.type in config file: '%s'; must be 'mysql'", dbConfig.Type)
	}

	return dbIns, nil
}

func (sdb *SqlxDB) CloseDB() error {
	if sdb.DB != nil && sdb.IsDBInitialized {
		err := sdb.DB.Close()
		if err != nil {
			return err
		}
		sdb.DB = nil
		return nil
	}
	return errors.New("sqlx handle not initialized")
}

func (sdb *SqlxDB) GetSqlxIns() (*sqlx.DB, error) {
	if sdb.DB != nil && sdb.IsDBInitialized {
		return sdb.DB, nil
	}
	return nil, errors.New("sqlx handle not initialized")
}

func GetX509CertificateFromPEM(cert []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(cert)
	if block == nil {
		return nil, errors.New("Failed to PEM decode certificate")
	}
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Error parsing certificate %s", err.Error())
	}
	return x509Cert, nil
}

func checkCertDates(certFile string) error {
	certPEM, err := ioutil.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("Failed to read file '%s' %s", certFile, err.Error())
	}

	cert, err := GetX509CertificateFromPEM(certPEM)
	if err != nil {
		return err
	}

	notAfter := cert.NotAfter
	currentTime := time.Now().UTC()

	if currentTime.After(notAfter) {
		return errors.New("Certificate provided has expired")
	}

	notBefore := cert.NotBefore
	if currentTime.Before(notBefore) {
		return errors.New("Certificate provided not valid until later date")
	}

	return nil
}

func GetClientTLSConfig(cacert string, key string, cert string) (*tls.Config, error) {
	var certs []tls.Certificate
	if cert != "" {
		err := checkCertDates(cert)
		if err != nil {
			return nil, err
		}

		clientCert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}

		certs = append(certs, clientCert)
	}
	rootCAPool := x509.NewCertPool()
	if len(cacert) == 0 {
		return nil, errors.New("No trusted root certificates for TLS were provided")
	}

	caCert, err := ioutil.ReadFile(cacert)
	if err != nil {
		return nil, fmt.Errorf("Failed to read '%s' %s", cacert, err.Error())
	}
	ok := rootCAPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("Failed to process certificate from file %s %s", cacert, err.Error())
	}

	config := &tls.Config{
		Certificates: certs,
		RootCAs:      rootCAPool,
	}

	return config, nil
}

func createMySQLDatabase(dbName string, db *sqlx.DB) error {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return fmt.Errorf("Failed to execute create database query %s", err.Error())
	}

	return nil
}

func createMySQLTables(dbName string, db *sqlx.DB) error {
	// var err error
	// logger.Debug("Creating test101 table if it doesn't exist")
	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS test101 (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,openid VARCHAR(80),pwd VARCHAR(80),phone VARCHAR(20),UNIQUE INDEX indexphone (phone)) DEFAULT CHARSET=utf8 COLLATE utf8_bin"); err != nil {
	// 	return errors.New(fmt.Sprintf("Error creating test101 table %s", err.Error()))
	// }

	// var err error
	// logger.Debug("Creating users table if it doesn't exist")
	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,openid VARCHAR(80),pwd VARCHAR(80),phone VARCHAR(20),UNIQUE INDEX indexphone (phone)) DEFAULT CHARSET=utf8 COLLATE utf8_bin"); err != nil {
	// 	return errors.New(fmt.Sprintf("Error creating users table %s", err.Error()))
	// }

	// logger.Debug("Creating upoint table if it doesn't exist")
	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS upoint (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,phone VARCHAR(20),daytime BIGINT UNSIGNED,tag VARCHAR(50),point VARCHAR(30),INDEX indexphone (phone), INDEX indextime (daytime)) DEFAULT CHARSET=utf8 COLLATE utf8_bin"); err != nil {
	// 	return errors.New(fmt.Sprintf("Error creating upoint table %s", err.Error()))
	// }

	// logger.Debug("Creating word table if it doesn't exist")
	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS word (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,rand VARCHAR(20),second BIGINT UNSIGNED,INDEX indexrand (rand)) DEFAULT CHARSET=utf8 COLLATE utf8_bin"); err != nil {
	// 	return errors.New(fmt.Sprintf("Error creating word table %s", err.Error()))
	// }

	// 	if _, err = db.Exec(
	// 		`CREATE TABLE users(
	// 			id BIGINT UNSIGNED NOT NULL,
	// 			createdat DATETIME(3) NOT NULL,
	// 			updatedat DATETIME(3) NOT NULL,
	// 			deletedat DATETIME(3) ,
	// 			state INT DEFAULT 0,
	// 			access JSON,
	// 			lastseen DATETIME,
	// 			public JSON,
	// -- 			tags JSON,
	// 			PRIMARY KEY(id)
	// 		)`); err != nil {
	// 		return errors.New("Error creating user table" + err.Error())
	// 	}

	return nil
}

func MaskDBCred(str string) string {
	matches := dbURLRegex.FindStringSubmatch(str)
	// If there is a match, there should be three entries: 1 for
	// the match and 9 for submatches (see dbURLRegex regular expression)
	if len(matches) == 10 {
		matchIdxs := dbURLRegex.FindStringSubmatchIndex(str)
		substr := str[matchIdxs[0]:matchIdxs[1]]
		for idx := 1; idx < len(matches); idx++ {
			if matches[idx] != "" {
				if strings.Index(matches[idx], "user=") == 0 {
					substr = strings.Replace(substr, matches[idx], "user=****", 1)
				} else if strings.Index(matches[idx], "password=") == 0 {
					substr = strings.Replace(substr, matches[idx], "password=****", 1)
				} else {
					substr = strings.Replace(substr, matches[idx], "****", 1)
				}
			}
		}
		str = str[:matchIdxs[0]] + substr + str[matchIdxs[1]:len(str)]
	}
	return str
}
