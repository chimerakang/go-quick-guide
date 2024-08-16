package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	USERNAME         = "demo"
	PASSWORD         = "demo123"
	NETWORK          = "tcp"
	SERVER           = "127.0.0.1"
	PORT             = 3306
	DATABASE         = "demo"
	MaxLifetime  int = 10
	MaxOpenConns int = 10
	MaxIdleConns int = 10
)

type Userinfo struct {
	ID         int64     `gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;" json:"id,omitempty"`
	Username   string    `gorm:"type:varchar(20) NOT NULL;" json:"username,omitempty"`
	Password   string    `gorm:"type:varchar(100) NOT NULL;" json:"password,omitempty"`
	Department string    `gorm:"type:varchar(100) NOT NULL;" json:"department,omitempty"`
	Status     int32     `gorm:"type:int(5);" json:"status,omitempty"`
	CreatedAt  time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt  time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

func main() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	checkErr(err)
	// defer conn.Close()
	db, err1 := conn.DB()
	if err1 != nil {
		fmt.Println("get db failed:", err)
		return
	}
	db.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	fmt.Println("Connect to MySQL success")
	err = db.Ping()
	checkErr(err)

	// create table by auto migrate
	conn.Debug().AutoMigrate(&Userinfo{})
	// check table
	migrator := conn.Migrator()
	has := migrator.HasTable(&Userinfo{})

	//has := migrator.HasTable("GG")
	if !has {
		fmt.Println("table not exist")
	}

	//插入資料
	stmt, err := db.Prepare("INSERT userinfo SET username=?,department=?,created_at=?")
	checkErr(err)

	res, err := stmt.Exec("chimera", "Develop", "2024-08-01")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	//Update userinfo
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("chimera-update", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//Query userinfo
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	//Delete userinfo by uid
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
