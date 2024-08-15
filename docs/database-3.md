# 資料庫操作 gorm
by [@chimerakang](https://github.com/chimerakang)

---
## ORM 介紹
ORM，英文叫 Object Relational Mapping，翻譯成中文為物件關聯對映。
ORM 在網站開發結構中，是在『資料庫』和『 Model 資料容器』兩者之間，
簡單來說，它是一個幫助使用者更簡便、安全的去從資料庫讀取資料，因為 ORM 的一個特性為: 透過程式語言，去操作資料庫語言( SQL )，而這也是實作了物件導向的概念，產生的一種工具模式，

### ORM的優缺點

#### ORM 的優點:

* 安全性
首先使用 ORM 最大的優點就是在此，有一種常見的網路攻擊叫做 SQL 注入( SQL injection )，就是駭客在傳輸到網站伺服器裡的資料直接寫 SQL ，而我們網站某段 SQL 直接讀取該駭客傳來的資料並執行，如果傳來的是正常的資料如 Email ，就會沒事，但如果傳來的是 SQL 語句，且包含『 Delete 』這種會害死大家的詞，那網站的資料就可能會被惡意移除。
但如果是透過 ORM 的方式，就變成是操作程式語言，像文章開頭的例子，where 的 age 可以改成吃變數，如 query_age，但程式就會自動判斷，query_age 內如果是奇怪的值，像是 “Update…” 等等之類的，就會自動擋掉。

* 簡化性
簡單解釋一下開頭的程式語言，如果用 SQL 語句完整寫完會是:
```Select * From users Where age = 30```
是不是比用 ORM 版本的囉唆多了呀？ORM 版本會簡化，幾乎都會用到的 From 和 Select 的部分，可能會想說這又什麼好簡化的，久了不就習慣了。但當程式越寫越多時，會對一些重複性極高的寫法感到疲乏，而人在疲乏時就有可能犯一些愚蠢的錯誤，像是打錯字, 打翻杯子(?) 之類的...ORM 就能較為避免此狀況發生。

* 通用性
因為 ORM 是在程式語言和資料庫之間，因此就算我們的網站未來有資料庫轉移的問題，也比較不會遇到需要改寫程式的狀況。因為不同的資料庫間，SQL 語句的語法也會稍有差異，我們來比較下面 MySQL 和 MsSQL 做兩件同樣事
```
// MySQL
SELECT * FROM TestTable WHERE id=12 LIMIT 10
// MsSQL
SELECT TOP 10 * FROM TestTable WHERE id=12
```
這兩個程式碼都要做同一件事，可是寫法會有所不同，因此在替換資料庫時就會需要全面檢查網站的 SQL 程式碼，或是一個 PHP 用習慣 MySQL 的工程師，換到一家使用 PHP 配 MsSQL 的公司，那SQL 部分就需要好好調整一番。

但是使用 ORM 就比較不會有這個狀況，ORM 是跟開發者對程式語言的熟悉度與泛用 SQL 語法概念較有關係

#### 缺點
講完優點，那要來提提缺點的部分:

* 效能
這是所有 ORM 和程式工具的通病，當達成了方便性，通常都會犧牲到效能的問題，因為等於要多了『把程式語言轉譯成 SQL語言』這項工作，不過各大程式語言的 ORM 都有持續改善，
緩慢的效能大多都是 SQL 語句設計不良，而較少是因為使用 ORM 的關係

* 學習曲線高
對於完全初學者來說，ORM 必須融合 SQL 語言和程式語言兩個不同的概念和語法，因此學習曲線比單寫 SQL 高了一些，不過這算是針對初學者的問題。

* 複雜查詢維護性低
對於複雜的查詢，ORM 的使用上就較為力不從心。我們目前遇到的都算是簡單的查詢，
如果是跨好幾個表格，且要針對部分欄位做 Sum, count 等工作，那就變成除了 ORM 的寫法外，還要額外寫入原生的 SQL 語法。

## GORM
用raw sql語法操作db跟用orm其實沒有什麼誰好誰不好，如果sql語法寫起來跟吃飯喝水一樣easy的話，不一定要用orm。如果sql語法沒寫的那麼精準，也不想常常被DBA手刀打下去的話，是可以考慮使用orm來滿足工作需求。

Golang首推的ORM套件還是以 GORM為首選，
特點如下

* 全功能ORM
* 關聯(Has One，Has Many，Belongs To，Many To Many，多態，單表繼承)
* Create，Save，Update，Delete，Find 中鉤子方法
* 支持Preload、Joins的預加載
* 事務，嵌套事務，Save Point，Rollback To Saved Point
* Context，預編譯模式，DryRun 模式
* 批量插入，FindInBatches，Find/Create with Map，使用SQL 表達式、Context * Valuer 進行CRUD
* SQL 構建器，Upsert，數據庫鎖，Optimizer/Index/Comment Hint，命名參數，子查詢
* 複合主鍵，索引，約束
* 自動遷移
* 自定義Logger
* 靈活的可擴展插件API：Database Resolver（多數據庫，讀寫分離）、Prometheus…-
* 每個特性都經過了測試的重重考驗
* 開發者友好

更多的細節可以參考: [GORM指南](https://gorm.io/zh_CN/docs/index.html)

### DB連線
```go
package main

import (
	"fmt"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	UserName     string = "root"
	Password     string = "password"
	Addr         string = "127.0.0.1"
	Port         int    = 3306
	Database     string = "test"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
)

func main() { 
	//組合sql連線字串
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)    
	//連接MySQL
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}  
    //設定ConnMaxLifetime/MaxIdleConns/MaxOpenConns
	db, err1 := conn.DB()
	if err1 != nil {
		fmt.Println("get db failed:", err)
		return
	}
	db.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)	
}
```

### model的定義
對gorm來說，model struct定義好，gorm才能進行mapping，mapping表請參考[models](https://gorm.io/zh_CN/docs/models.html)
```go
type User struct {
	ID        int64     `gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;" json:"id,omitempty"`
	Username  string    `gorm:"type:varchar(20) NOT NULL;" json:"username,omitempty"`
	Password  string    `gorm:"type:varchar(100) NOT NULL;" json:"password,omitempty"`
	Status    int32     `gorm:"type:int(5);" json:"status,omitempty"`
	CreatedAt time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}
```

### table 操作
```go
package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	UserName     string = "root"
	Password     string = "password"
	Addr         string = "127.0.0.1"
	Port         int    = 3306
	Database     string = "test"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
)

type User struct {
	ID        int64     `gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;" json:"id,omitempty"`
	Username  string    `gorm:"type:varchar(20) NOT NULL;" json:"username,omitempty"`
	Password  string    `gorm:"type:varchar(100) NOT NULL;" json:"password,omitempty"`
	Status    int32     `gorm:"type:int(5);" json:"status,omitempty"`
	CreatedAt time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

func main() {

	//組合sql連線字串
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)
	//連接MySQL
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}

	db, err1 := conn.DB()
	if err1 != nil {
		fmt.Println("get db failed:", err)
		return
	}

	db.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
    
    //產生table
	conn.Debug().AutoMigrate(&User{})
    //判斷有沒有table存在
    migrator := conn.Migrator()
    has := migrator.HasTable(&User{})
	//has := migrator.HasTable("GG")
	if !has {
		fmt.Println("table not exist")
	}
}

```

### CREATE
create可透過model或是語法進行新增資料
```go
//新增model資料
user := User{UserName: "tester", Password: "12333", Status: 1}
//
result := conn.Debug().Create(&user)
if result.Error != nil {
    fmt.Println("Create failt")
}
if result.RowsAffected != 1 {
    fmt.Println("RowsAffected Number failt")
}
//執行結果：[rows:1] INSERT INTO `users` (`user_name`,`password`,`status`,`created_at`,`updated_at`) VALUES ('tester','12333',1,'2020-09-30 13:52:49.913','2020-09-30 13:52:49.913')

//只insert特定欄位值
conn.Debug().Select("UserName", "Password").Create(&user)
//執行結果：INSERT INTO `users` (`user_name`,`password`) VALUES ('tester','12333')

//不insert特定欄位值，但是這個語法好像遇到auto_increment的pk會Duplicate pk
conn.Debug().Omit("status").Create(&user)
//執行結果：INSERT INTO `users` (`user_name`,`password`,`created_at`,`updated_at`,`id`) VALUES ('tester','12333','2020-09-30 14:00:46.006','2020-09-30 14:00:46.006',10)，ID自己產生於語法中

//BATCH INSERT，v2版新增的特點
    users := []User{{UserName: "tester", Password: "12333", Status: 1}, {UserName: "gger", Password: "132333", Status: 1}, {UserName: "ininder", Password: "12333", Status: 1}}

result := conn.Debug().Create(&users)
if result.Error != nil {
    fmt.Println("Create failt")
}
fmt.Println("result.RowsAffected:", result.RowsAffected)
//執行結果：INSERT INTO `users` (`user_name`,`password`,`status`,`created_at`,`updated_at`) VALUES ('tester','12333',1,'2020-09-30 14:07:00.133','2020-09-30 14:07:00.133'),('gger','132333',1,'2020-09-30 14:07:00.133','2020-09-30 14:07:00.133'),('ininder','12333',1,'2020-09-30 14:07:00.133','2020-09-30 14:07:00.133')

//除了用slice struct外，也可以用map進行新增資料，map跟struct有蠻大的差異，struct會有default值，map是沒指定的就不會出現在insert語法上
conn.Debug().Model(&User{}).Create(map[string]interface{}{
    "UserName": "gg", "Password": "18",
})
//執行結果：INSERT INTO `users` (`password`,`user_name`) VALUES ('18','gg')

//也可以執行raw語法
conn.Exec("INSERT INTO `users` (`password`,`user_name`) VALUES (?,?)", "G123", "999")
```
### SELECT
* Find：取得全部筆數
* First：根據pk升冪取得第一筆
* Last：根據pk降冪取得第一筆
* Take：取得一筆資料

對gorm的內部定義來說，查無資料視為錯誤，所以要判斷查無資料還是是其他錯誤可以用下面語法判斷
```go

var user User
var users []User

res := conn.Debug().Find(&users)
fmt.Println(res.RowsAffected)
//SELECT * FROM `users`
conn.Debug().First(&user)
//SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1
conn.Debug().Take(&user)
//SELECT * FROM `users` WHERE `users`.`id` = 1 LIMIT 1
conn.Debug().Last(&user)
//SELECT * FROM `users` WHERE `users`.`id` = 1 ORDER BY `users`.`id` DESC LIMIT 1
```
### WHERE Clause
where裡面的寫法類似raw sql語法，使用?做佔位符號
```go
conn.Debug().Where("user_name = ?", "tester4").Find(&user)
//SELECT * FROM `users` WHERE user_name = 'tester4'
conn.Debug().Where("user_name IN ?", []string{"tester4", "tester3"}).Find(&users)
//SELECT * FROM `users` WHERE user_name IN ('tester4','tester3')，IN語法可以用slice
conn.Debug().Where("user_name LIKE ?", "%tester%").Find(&users)
//SELECT * FROM `users` WHERE user_name LIKE '%tester%'
```
有個很好用的語法Pluck，可以把select出來的值變成slice的
```go
conn.Debug().Table("users").Pluck("user_name", &us)
//us: [tester1 tester2 tester3 tester4 tester tester tester tester tester tester tester gger ininder jinzhu gg gg 999]
```
還有更多的語法可以參考官方文件 [GORM QUERY](https://gorm.io/docs/query.html)

### UPDATE與DELETE
* Save：更新整個struct
* Update：更新單個欄位
* Updates：更新多個欄位，使用strcut或是map
```go
//使用Save
conn.First(&user)
user.Password = "GG"
conn.Debug().Save(&user)
//UPDATE `users` SET `user_name`='tester1',`password`='GG',`status`=1,`created_at`='2020-09-30 05:46:22',`updated_at`='2020-09-30 15:46:07.403' WHERE `id` = 1

//使用Update
conn.Debug().Model(&User{}).Where("id = ?", 1).Update("password", "helloGG")
//UPDATE `users` SET `password`='helloGG',`updated_at`='2020-09-30 15:48:45.676' WHERE id = 1

//使用Updates
conn.Debug().Model(&User{}).Where("id = ?", 1).Updates(User{UserName: "hello", Password: "GG"})
//UPDATE `users` SET `user_name`='hello',`password`='GG',`updated_at`='2020-09-30 15:51:27.78' WHERE id = 1
conn.Debug().Model(&User{}).Where("id = ?", 1).Updates(map[string]interface{}{"UserName": "hello", "Password": "GG"})
//UPDATE `users` SET `password`='GG',`user_name`='hello',`updated_at`='2020-09-30 15:51:27.784' WHERE id = 1

//DELETE
conn.Debug().Where("id = ?", "1").Delete(&User{})
//DELETE FROM `users` WHERE id = '1'
```

---
## NEXT: [GORM進階使用](./database-3a.md)
## [Gorm番外篇](./database-3plus.md)
