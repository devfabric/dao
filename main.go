package main

import (
	"fmt"
	"log"

	"github.com/devfabric/dao/dao"
)

type Users struct {
	ID   int64
	Name string
	F1   string
	F2   string
}

func main() {
	//init db
	db, err := dao.NewDB("./")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.CloseDB()
	fmt.Println("connect db success")

	dbsqlx, err := db.GetSqlxIns()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//tx
	{
		tx, err := dbsqlx.Begin()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result, err := tx.Exec("drop table if exists users;")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		affected, err := result.RowsAffected()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("affected:", affected)

		if _, err = dbsqlx.Exec(
			`CREATE TABLE users (
				id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
				name varchar(50) NOT NULL,
				f1 varchar(50) NOT NULL,
				f2 varchar(50) NOT NULL,
				PRIMARY KEY (id)
				) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8`); err != nil {
			fmt.Println(err.Error())
			return
		}
		tx.Commit()
	}

	{

		result, err := dbsqlx.NamedExec(`INSERT INTO users (name, f1, f2) VALUES (:name, :f1, :f2)`, map[string]interface{}{
			"name": "Bin",
			"f1":   "Smuth",
			"f2":   "bensmith@allblacks.nz",
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		affected, err := result.RowsAffected()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("MustExec affected:", affected)
	}

	{
		rows, err := dbsqlx.NamedQuery(`SELECT * FROM users WHERE name=:fn`, map[string]interface{}{"fn": "Bin"})
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var user Users
		for rows.Next() {
			// fmt.Printf("%#v\n", rows)

			err := rows.StructScan(&user)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%#v\n", user)
		}
	}

	{

		fmt.Println("###########################")
		tx := dbsqlx.MustBegin()
		for i := 0; i < 100; i++ {
			name := fmt.Sprintf("Jason-%d", i+1)
			f1 := fmt.Sprintf("Moiron-%d", i+1)
			f2 := fmt.Sprintf("jmoiron-%d", i+1)
			tx.MustExec("INSERT INTO users (name, f1, f2) VALUES (?, ?, ?)", name, f1, f2)
		}
		tx.Commit()

		rows, err := dbsqlx.Queryx(`SELECT * FROM users`)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var user Users
		for rows.Next() {
			// fmt.Printf("%#v\n", rows)

			err := rows.StructScan(&user)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%#v\n", user)
		}
	}
}
