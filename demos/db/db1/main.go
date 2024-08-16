package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "demo"
	PASSWORD = "demo123"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "demo"
)

type User struct {
	ID       string
	Username string
	Password string
}

func main() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Connect to MySQL error:%v\n", err)
		return
	}
	fmt.Println("Connect to MySQL success")
	if err := db.Ping(); err != nil {
		fmt.Printf("Connect to MySQL error:%v\n", err)
		return
	}
	defer db.Close()
	// CreateTable(db)
	InsertUser(db, "test1", "test234")
	QueryUser(db, "test1")
	fmt.Println("The End")
}

// create new table
func CreateTable(db *sql.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS users(
	id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
        username VARCHAR(64),
        password VARCHAR(64)
	); `
	if _, err := db.Exec(sql); err != nil {
		fmt.Printf("Create Table error:%v\n", err)
		return err
	}
	fmt.Println("Create Table success")
	return nil
}

// insert new user into db
func InsertUser(DB *sql.DB, username, password string) error {
	_, err := DB.Exec("insert INTO users(username,password) values(?,?)", username, password)
	if err != nil {
		fmt.Printf("Creat User error：%v\n", err)
		return err
	}
	fmt.Println("Create User success")
	return nil
}

// query user by user name
func QueryUser(db *sql.DB, username string) {
	user := new(User)
	row := db.QueryRow("select * from users where username=?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		fmt.Printf("Reflect user error：%v\n", err)
		return
	}
	fmt.Println("Query user success:", *user)
}
