package dao

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/devfabric/dao/config"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewUserRegistryMySQL opens a connection to a postgres database
func NewUserRegistryMySQL(datasource string, dbConfig *config.DataBaseConfig) (*SqlxDB, error) {

	dbName := getDBName(datasource)

	re := regexp.MustCompile(`\/([0-9,a-z,A-Z$_$-]+)`)
	connStr := re.ReplaceAllString(datasource, "/")

	if dbConfig.TLS.Enabled {
		tlsConfig, err := GetClientTLSConfig(dbConfig.TLS.CaCert, dbConfig.TLS.KeyFile, dbConfig.TLS.CertFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to get client TLS for MySQL %s", err.Error())
		}
		mysql.RegisterTLSConfig("custom", tlsConfig)
	}

	db, err := sqlx.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to open MySQL database %s, %s", connStr, err.Error())
	}
	count := 0
	for {
		err = db.Ping()
		if err != nil {
			if count > 10 {
				return nil, fmt.Errorf("Failed to connect to MySQL database %s", err.Error())
			}
			count++
			time.Sleep(time.Second * time.Duration(count))
		} else {
			break
		}
	}

	// err = dropMySQLDatabase(dbName, db)
	err = createMySQLDatabase(dbName, db)
	if err != nil {
		return nil, err
	}
	//if err != nil {
	//	log.Logger.Debugf("Failed to create MySQL database %s", err.Error())
	//	return nil, fmt.Errorf("Failed to create MySQL database %s", err.Error())
	//} else {
	//	log.Logger.Debug("create mysql database success")
	//}
	db.Close()

	//log.Logger.Debugf("Connecting to database '%s', using connection string: '%s'", dbName, MaskDBCred(dbConfig.DataSource))
	db, err = sqlx.Open("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database (%s) in MySQL server %s", dbName, err.Error())
	}
	db.SetMaxIdleConns(dbConfig.MaxIdle)
	db.SetMaxOpenConns(dbConfig.MaxOpen)

	err = createMySQLTables(dbName, db)
	if err != nil {
		return nil, fmt.Errorf("Failed to create MySQL tables %s", err.Error())
	}

	return &SqlxDB{DB: db, IsDBInitialized: true}, nil
}

// func createMySQLDatabase(dbName string, db *sqlx.DB) error {
// 	//log.Logger.Debugf("Creating MySQL Database (%s) if it does not exist...", dbName)

// 	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
// 	if err != nil {
// 		return errors.Wrap(err, "Failed to execute create database")
// 	}

// 	return nil
// }

// func dropMySQLDatabase(dbName string, db *sqlx.DB) error {
// 	//log.Logger.Debugf("DROP MySQL Database (%s) if it does not exist...", dbName)

// 	_, err := db.Exec("DROP DATABASE IF EXISTS" + dbName)
// 	if err != nil {
// 		return fmt.Errorf("Failed to execute create database query %s", err.Error())
// 	}

// 	return nil
// }

// func createMySQLTables(dbName string, db *sqlx.DB) error {
// 	// log.Logger.Debug("Creating users table if it doesn't exist")
// 	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id VARCHAR(255) NOT NULL, token blob, type VARCHAR(256), affiliation VARCHAR(1024), attributes TEXT, state INTEGER, max_enrollments INTEGER, level INTEGER DEFAULT 0, PRIMARY KEY (id)) DEFAULT CHARSET=utf8 COLLATE utf8_bin"); err != nil {
// 	// 	return errors.Wrap(err, "Error creating users table")
// 	// }

// 	// log.Logger.Debug("Creating affiliations table if it doesn't exist")
// 	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS affiliations (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(1024) NOT NULL, prekey VARCHAR(1024), level INTEGER DEFAULT 0, PRIMARY KEY (id))"); err != nil {
// 	// 	return errors.Wrap(err, "Error creating affiliations table")
// 	// }
// 	// log.Logger.Debug("Creating index on 'name' in the affiliations table")
// 	// if _, err := db.Exec("CREATE INDEX name_index on affiliations (name)"); err != nil {
// 	// 	if !strings.Contains(err.Error(), "Error 1061") { // Error 1061: Duplicate key name, index already exists
// 	// 		return errors.Wrap(err, "Error creating index on affiliations table")
// 	// 	}
// 	// }
// 	// log.Logger.Debug("Creating certificates table if it doesn't exist")
// 	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS certificates (id VARCHAR(255), serial_number varbinary(128) NOT NULL, authority_key_identifier varbinary(128) NOT NULL, ca_label varbinary(128), status varbinary(128) NOT NULL, reason int, expiry timestamp DEFAULT 0, revoked_at timestamp DEFAULT 0, pem varbinary(4096) NOT NULL, level INTEGER DEFAULT 0, PRIMARY KEY(serial_number, authority_key_identifier)) DEFAULT CHARSET=utf8 COLLATE utf8_bin"); err != nil {
// 	// 	return errors.Wrap(err, "Error creating certificates table")
// 	// }
// 	// log.Debug("Creating credentials table if it doesn't exist")
// 	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS credentials (id VARCHAR(255), revocation_handle varbinary(128) NOT NULL, cred varbinary(4096) NOT NULL, ca_label varbinary(128), status varbinary(128) NOT NULL, reason int, expiry timestamp DEFAULT 0, revoked_at timestamp DEFAULT 0, level INTEGER DEFAULT 0, PRIMARY KEY(revocation_handle)) DEFAULT CHARSET=utf8 COLLATE utf8_bin"); err != nil {
// 	// 	return errors.Wrap(err, "Error creating certificates table")
// 	// }
// 	// log.Debug("Creating revocation_authority_info table if it does not exist")
// 	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS revocation_authority_info (epoch INTEGER, next_handle INTEGER, lasthandle_in_pool INTEGER, level INTEGER DEFAULT 0, PRIMARY KEY (epoch))"); err != nil {
// 	// 	return errors.Wrap(err, "Error creating revocation_authority_info table")
// 	// }
// 	// log.Debug("Creating nonces table if it does not exist")
// 	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS nonces (val VARCHAR(255) NOT NULL, expiry timestamp, level INTEGER DEFAULT 0, PRIMARY KEY (val))"); err != nil {
// 	// 	return errors.Wrap(err, "Error creating nonces table")
// 	// }
// 	// log.Debug("Creating properties table if it does not exist")
// 	// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS properties (property VARCHAR(255), value VARCHAR(256), PRIMARY KEY(property))"); err != nil {
// 	// 	return errors.Wrap(err, "Error creating properties table")
// 	// }
// 	// _, err := db.Exec(db.Rebind("INSERT INTO properties (property, value) VALUES ('identity.level', '0'), ('affiliation.level', '0'), ('certificate.level', '0'), ('credential.level', '0'), ('rcinfo.level', '0'), ('nonce.level', '0')"))
// 	// if err != nil {
// 	// 	if !strings.Contains(err.Error(), "1062") { // MySQL error code for duplicate entry
// 	// 		return err
// 	// 	}
// 	// }
// 	return nil
// }

// getDBName gets database name from connection string
func getDBName(datasource string) string {
	var dbName string
	datasource = strings.ToLower(datasource)

	re := regexp.MustCompile(`(?:\/([^\/?]+))|(?:dbname=([^\s]+))`)
	getName := re.FindStringSubmatch(datasource)
	if getName != nil {
		dbName = getName[1]
		if dbName == "" {
			dbName = getName[2]
		}
	}

	return dbName
}

// GetConnStr gets connection string without database
func getConnStr(datasource string, dbname string) string {
	re := regexp.MustCompile(`(dbname=)([^\s]+)`)
	connStr := re.ReplaceAllString(datasource, fmt.Sprintf("dbname=%s", dbname))
	return connStr
}
