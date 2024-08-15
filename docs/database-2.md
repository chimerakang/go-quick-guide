# database/sql 資料庫介面
by [@chimerakang](https://github.com/chimerakang)

---

## DB
資料庫的物件實例，同時內部也有自己實現的連線池，其池是一個包含多個open和idle連線的池子。使用時被選到的連線會被標記成open, 完成後會被標記成idle，他需要有driver打開或關閉資料庫，管理連線池。同時它也是線程安全的，可以不必重複創立，只需要一個就可以傳遞給多個goroutine使用

一個泳池的泳道數量設定好之後就是固定的, 只要有人使用, 該泳道就被open, 等到有人離開泳道, 該泳道就被視為idel.
maxlifetime, 就視為該泳道的開放使用時間吧, 也許設置成1小時換水清理一次.
但有人還在用的話, 當然就要等它被釋放出來, 才能開始清潔. (有點牽強的例子)



### SetMaxIdleConns
設定空閒連線池的最大連線數

```go
func (db *DB) SetMaxIdleConns(n int) {...}
```

### SetMaxOpenConns
設定最大打開的連線數, 0表示不限制.

```go
func（db * DB）SetMaxOpenConns（n int） {...}
```

### SetConnMaxLifetime
設定連線可重複利用的最大時間長度, 0是預設值表示沒有max life, 總是可重複使用.

```go
func (db *DB) SetConnMaxLifetime(d time.Duration) {...}
```

```go
db.SetConnMaxLifetime(time.Hour)
// 設定每一個連線最大生命週期1hr
```

### Open
```go
func Open(driverName, dataSourceName string) (*DB, error) {...}
```
初始化一個sql.DB對象, 但還沒真正建立連線.
會啟動一個connectionOpener的goroutine.
也初始化一個openerCh channel.
需要連線時, 就對該channel發送數據就好.

```go
db := &DB {
    driver:   driveri,
    dsn:      dataSourceName,
    openerCh: make(chan struct{}, connectionRequestQueueSize),
    lastPut:  make(map[*driverConn]string),
}
```
在這兩種情形下會去建立連線 :

會在第一次呼叫`ping()`, 真正建立連線
呼叫`db.Exec()`或者`db.Query()`, 如果空閒連線池有連線, 就直接取用; 如果沒有就會產生一個新的連線

### Config
`go-sql-driver/mysql`所提供的結構體, 有提供幾個針對DSN的方法

```go
// config的建構式
func NewConfig() *Config
// 把dsn字串剖析轉成config物件
func ParseDSN(dsn string) (cfg *Config, err error)
// 複製一個config物件的副本
func (cfg *Config) Clone() *Config
// 把config結構體格式化成dsn格式的連線字串
func (cfg *Config) FormatDSN() string
```

### Row
呼叫QueryRow()之後返回的單行結果
掃描取得的結果到dest上; 但如果有多個結果, scan做完第一行後, 就會丟棄其他行的資料了.
如果row是沒有資料的, 則是會返回ErrNoRows這錯誤.

```go
func (r *Row) Scan(dest ...interface{}) error {...}
```

### Rows
查詢的結果, 會有個指標從開始直到結束.
呼叫Next(), 就去取得下一行row的資料.

```go
// 主要是用來讓連線釋放回連線池
func (rs *Rows) Close() error {...}
// 對資料做走訪,正常結束的話內部會自動呼叫rows.Close()
func (rs *Rows) Next() bool
// 切換到下一個結果集, 一次查詢是可以返回多個結果集的
func (rs *Rows) NextResultSet() bool
// 
func (rs *Rows) Scan(dest ...interface{}) error
```

### Scanner
```go
type Scanner interface {
    // Scan assigns a value from a database driver.
    Scan(src interface{}) error
}
```

### Stmt
查詢語句, DDL、DML等的prepared sql語句.

### Tx
一個進行中的資料庫事務.
呼叫db.Begin()之後, 會取得一個tx物件, 需要呼叫`Commit()`或`Rollback()`才會結束事務.並且歸還連線回連線池.

### Result
主要是針對insert、update、delete的操作所返回的結果.
Result是個接口, 定義了`LastInsertId()`和`RowsAffected()`
各驅動會實作這接口的兩個方法.
```go
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}
```

### MySqlDriver的實作:
```go
package mysql

type mysqlResult struct {
	affectedRows int64
	insertId     int64
}

func (res *mysqlResult) LastInsertId() (int64, error) {
	return res.insertId, nil
}

func (res *mysqlResult) RowsAffected() (int64, error) {
	return res.affectedRows, nil
}
```

### Nullable Property
各種允許null的基礎型別. 且都有實現Scanner的Scan().

### IsolationLevel
事務用的資料隔離等級.

```go
const（
    LevelDefault  IsolationLevel = iota 
    LevelReadUncommitted 
    LevelReadCommitted 
    LevelWriteCommitted 
    LevelRepeatableRead 
    LevelSnapshot 
    LevelSerializable 
    LevelLinearizable 
）
```

---
## NEXT: [資料庫操作: gorm](./database-3.md)
