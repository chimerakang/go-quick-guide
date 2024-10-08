# database/sql 介面
by [@chimerakang](https://github.com/chimerakang)

---
## 介紹
Go為開發資料庫驅動定義了一些標準介面，開發者可以根據定義的介面來開發相應的資料庫驅動，這樣做有一個好處，只要是按照標準介面開發的程式碼， 以後需要遷移資料庫時，不需要任何修改。那麼 Go 都定義了哪些標準介面呢？讓我們來詳細的分析一下

## sql.Register
這個存在於 database/sql 的函式是用來註冊資料庫驅動的，當第三方開發者開發資料庫驅動時，都會實現 init 函式，在 init 裡面會呼叫這個Register(name string, driver driver.Driver)完成本驅動的註冊。

我們來看一下 mymysql、sqlite3 的驅動裡面都是怎麼呼叫的：

```go
//https://github.com/mattn/go-sqlite3 驅動
func init() {
    sql.Register("sqlite3", &SQLiteDriver{})
}

//https://github.com/mikespook/mymysql 驅動
// Driver automatically registered in database/sql
var d = Driver{proto: "tcp", raddr: "127.0.0.1:3306"}
func init() {
    Register("SET NAMES utf8")
    sql.Register("mymysql", &d)
}
```
我們看到第三方資料庫驅動都是透過呼叫這個函式來註冊自己的資料庫驅動名稱以及相應的 driver 實現。在 database/sql 內部透過一個 map 來儲存使用者定義的相應驅動。

```go
var drivers = make(map[string]driver.Driver)

drivers[name] = driver
```
因此透過 database/sql 的註冊函式可以同時註冊多個數據函式庫驅動，只要不重複。

在我們使用 database/sql 介面和第三方函式庫的時候經常看到如下:

```go
   import (
       "database/sql"
        _ "github.com/mattn/go-sqlite3"
   )
```
新手都會被這個 _ 所迷惑，其實這個就是 Go 設計的巧妙之處，我們在變數賦值的時候經常看到這個符號，它是用來忽略變數賦值的佔位符，那麼套件引入用到這個符號也是相似的作用，這兒使用 _ 的意思是引入後面的套件名而不直接使用這個套件中定義的函式，變數等資源。

套件在引入的時候會自動呼叫套件的 init 函式以完成對套件的初始化。因此，我們引入上面的資料庫驅動套件之後會自動去呼叫 init 函式，然後在 init 函式裡面註冊這個資料庫驅動，這樣我們就可以在接下來的程式碼中直接使用這個資料庫驅動了。

## driver.Driver
Driver 是一個數據函式庫驅動的介面，他定義了一個 method： Open(name string)，這個方法回傳一個數據函式庫的 Conn 介面。

```go
type Driver interface {
    Open(name string) (Conn, error)
}
```
回傳的 Conn 只能用來進行一次 goroutine 的操作，也就是說不能把這個 Conn 應用於 Go 的多個 goroutine 裡面。如下程式碼會出現錯誤

```go
go goroutineA (Conn)  //執行查詢操作
go goroutineB (Conn)  //執行插入操作
```

上面這樣的程式碼可能會使 Go 不知道某個操作究竟是由哪個 goroutine 發起的，從而導致資料混亂，比如可能會把 goroutineA 裡面執行的查詢操作的結果回傳給 goroutineB 從而使 B 錯誤地把此結果當成自己執行的插入資料。

第三方驅動都會定義這個函式，它會解析 name 參數來取得相關資料庫的連線資訊，解析完成後，它將使用此資訊來初始化一個 Conn 並回傳它。

## driver.Conn
Conn 是一個數據函式庫連線的介面定義，他定義了一系列方法，這個 Conn 只能應用在一個 goroutine 裡面，不能使用在多個 goroutine 裡面，詳情請參考上面的說明。

```go
type Conn interface {
    Prepare(query string) (Stmt, error)
    Close() error
    Begin() (Tx, error)
}
```
Prepare 函式回傳與當前連線相關的執行 Sql 語句的準備狀態，可以進行查詢、刪除等操作。

Close 函式關閉當前的連線，執行釋放連線擁有的資源等清理工作。因為驅動實現了 database/sql 裡面建議的 conn pool，所以你不用再去實現快取 conn 之類別的，這樣會容易引起問題。

Begin 函式回傳一個代表交易處理的 Tx，透過它你可以進行查詢，更新等操作，或者對交易進行回復 (Rollback)、提交。

### driver.Stmt
Stmt 是一種準備好的狀態，和 Conn 相關聯，而且只能應用於一個 goroutine 中，不能應用於多個 goroutine。

```go
type Stmt interface {
    Close() error
    NumInput() int
    Exec(args []Value) (Result, error)
    Query(args []Value) (Rows, error)
}
```
Close 函式關閉當前的連結狀態，但是如果當前正在執行 query，query 還是有效回傳 rows 資料。

NumInput 函式回傳當前預留參數的個數，當回傳 >=0 時資料庫驅動就會智慧檢查呼叫者的參數。當資料庫驅動套件不知道預留參數的時候，回傳-1。

Exec 函式執行 Prepare 準備好的 sql，傳入參數執行 update/insert 等操作，回傳 Result 資料

Query 函式執行 Prepare 準備好的 sql，傳入需要的參數執行 select 操作，回傳 Rows 結果集

## driver.Tx
交易處理一般就兩個過程，提交或者回復 (Rollback)。資料庫驅動裡面也只需要實現這兩個函式就可以

```go
type Tx interface {
    Commit() error
    Rollback() error
}
```
這兩個函式一個用來提交一個交易，一個用來回復 (Rollback)交易。

## driver.Execer
這是一個 Conn 可選擇實現的介面

```go
type Execer interface {
    Exec(query string, args []Value) (Result, error)
}
```
如果這個介面沒有定義，那麼在呼叫 DB.Exec，就會首先呼叫 Prepare 回傳 Stmt，然後執行 Stmt 的 Exec，然後關閉 Stmt。

## driver.Result
這個是執行 Update/Insert 等操作回傳的結果介面定義

```go
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}
```
LastInsertId 函式回傳由資料庫執行插入操作得到的自增 ID 號。

RowsAffected 函式回傳 query 操作影響的資料條目數。

## driver.Rows
Rows 是執行查詢回傳的結果集介面定義

```go
type Rows interface {
    Columns() []string
    Close() error
    Next(dest []Value) error
}
```
Columns 函式回傳查詢資料庫表的欄位資訊，這個回傳的 slice 和 sql 查詢的欄位一一對應，而不是回傳整個表的所有欄位。

Close 函式用來關閉 Rows 迭代器。

Next 函式用來回傳下一條資料，把資料賦值給 dest。dest 裡面的元素必須是 driver.Value 的值除了 string，回傳的資料裡面所有的 string 都必須要轉換成[]byte。如果最後沒資料了，Next 函式最後回傳 io.EOF。

## driver.RowsAffected
RowsAffected 其實就是一個 int64 的別名，但是他實現了 Result 介面，用來底層實現 Result 的表示方式

```go
type RowsAffected int64

func (RowsAffected) LastInsertId() (int64, error)

func (v RowsAffected) RowsAffected() (int64, error)
```

## driver.Value
Value 其實就是一個空介面，他可以容納任何的資料

```go
type Value interface{}
```
drive 的 Value 是驅動必須能夠操作的 Value，Value 要麼是 nil，要麼是下面的任意一種
```go
int64
float64
bool
[]byte
string   [*]除了 Rows.Next 回傳的不能是 string.
time.Time
```
## driver.ValueConverter
ValueConverter 介面定義了如何把一個普通的值轉化成 driver.Value 的介面

```go
type ValueConverter interface {
    ConvertValue(v interface{}) (Value, error)
}
```
在開發的資料庫驅動套件裡面實現這個介面的函式在很多地方會使用到，這個 ValueConverter 有很多好處：

* 轉化 driver.value 到資料庫表相應的欄位，例如 int64 的資料如何轉化成資料庫表 uint16 欄位

* 把資料庫查詢結果轉化成 driver.Value 值

* 在 scan 函式裡面如何把 driver.Value 值轉化成使用者定義的值

## driver.Valuer
Valuer 介面定義了回傳一個 driver.Value 的方式
```go

type Valuer interface {
    Value() (Value, error)
}
```
很多型別都實現了這個 Value 方法，用來自身與 driver.Value 的轉化。

透過上面的講解，你應該對於驅動的開發有了一個基本的了解，一個驅動只要實現了這些介面就能完成增刪查改等基本操作了，剩下的就是與相應的資料庫進行資料互動等細節問題了，在此不再贅述。

## database/sql
database/sql 在 database/sql/driver 提供的介面基礎上定義了一些更高階的方法，用以簡化資料庫操作，同時內部還建議性地實現一個 conn pool。

```go
type DB struct {
    driver      driver.Driver
    dsn         string
    mu       sync.Mutex // protects freeConn and closed
    freeConn []driver.Conn
    closed   bool
}
```

我們可以看到 Open 函式回傳的是 DB 物件，裡面有一個 freeConn，它就是那個簡易的連線池。它的實現相當簡單或者說簡陋，就是當執行db.prepare -> db.prepareDC的時候會defer dc.releaseConn，然後呼叫db.putConn，也就是把這個連線放入連線池，每次呼叫db.conn的時候會先判斷 freeConn 的長度是否大於 0，大於 0 說明有可以複用的 conn，直接拿出來用就是了，如果不大於 0，則建立一個 conn，然後再回傳之。

---
## Next: [DB介面](./database-2.md)