# Gorm 番外篇
by [@chimerakang](https://github.com/chimerakang)

---
## 為什麼有這一個章節
ORM 的好處是讓程式中的物件與資料庫之間的操作更直覺，資料庫就像程式的關係是很緊密但是就是有些不一樣的地方，代價就是得在程式碼中重現資料庫的設計，也就是說，得把 `table schema` 在程式碼中抄一份，而表之間的關聯也得照著定義一份。

## code-first / schema-first
這裡牽涉到 `ORM package` 的一個分類： `code-first` / `schema-first`。

`Code-first` 以程式的物件為 `data model`，以此為本，產生出 `database schema`，常用於資料庫設計或實作仍未知的情況。

`Schema-first` 先定義 `database schema` ，然後才設計程式的物件去搭配，常用於已經有資料庫的情況。

`Gorm` 屬於 `Code-first`，雖然常常被`schema`問題困擾很久，踩過下面的坑之後，我覺得 `Gorm` 的設計還是很吸引人的。

## Model 不是 Table Schema
帶著 `schema-first` 的概念去理解 `Gorm`會很快，也會很快踩到坑。

Gorm 的 Model ，意義上是物件的樣版，而不是 `Table Schema`，因為用 Model 宣告出來的，並不是 `table`，而是物件，階級於 `table row` 相當：
```go
// 這是個 Model，不是 table schema
type User struct {
  ID           uint
  Name         string
  Email        *string
  Age          uint8
  Birthday     *time.Time
  MemberNumber sql.NullString
  ActivatedAt  sql.NullTime
  CreatedAt    time.Time
  UpdatedAt    time.Time
}

// 這是個 Object，不是 table，比較接近 row
user := User{ Name: "Gary" } // user 比較像是資料表中的一個 row，或是 record
```
`Gorm` 從頭到尾都沒有去定義過 `Table` 。

## Model & field 名稱不會對應到 table 與 column 名稱
前面提到 `Gorm` 從頭到尾都沒有去定義過 `Table` ，那麼查詢的時候 `Table name` 哪來？

Gorm貼心的為你產生的 `Table name`，官方文件的`Model` 舉例叫做 `User`，產生的`SQL`會幫你自動產生 `table name` 就叫做 `users`，你以為它只幫你做了複數轉換？不，它還幫你轉大小寫，而且連 `column name` 都這麼做，例如 :
```go
// 這是個Model
type GoodUser struct {
  LaAge int
}

// 宣告一個物件
gu := GoodUser{18}

// 執行查詢
db.Select("LaAge").Take(&gu)

// 實際執行SQL
// SELECT "la_age" FROM "good_users" LIMIT 1

// Table name:  GoodUser -> good_users ,with 's'
// Column name: LaAge -> la_age
```
注意無論是 `table name` 或是 `column name`，都被大寫改小寫，[駝峰命名法](https://zh.wikipedia.org/zh-tw/%E9%A7%9D%E5%B3%B0%E5%BC%8F%E5%A4%A7%E5%B0%8F%E5%AF%AB)改成[蛇型命名法](https://zh.wikipedia.org/zh-tw/%E8%9B%87%E5%BD%A2%E5%91%BD%E5%90%8D%E6%B3%95)。

### 改用自定義 Table name

根據[官方文件](https://Gorm.io/docs/conventions.html#Pluralized-Table-Name)，對 `Model` 寫個 `function` 吐 table name 就可以了：

```go
// 這是個Model
type GoodUser struct {
  LaAge int
}

// 自定義 table name
// 需實作介面 : type Tabler interface { TableName() string }
func (GoodUser) TableName() string {
  return "GoodUser"
}

// 宣告一個物件
gu := GoodUser{18}

// 執行查詢
db.Select("LaAge").Take(&gu)

// 實際執行SQL
// SELECT "la_age" FROM "GoodUser" LIMIT 1

// Table name:  GoodUser -> GoodUser 被定義了
// Column name: LaAge -> la_age
```

也是有硬幹的做法，但這樣一點都不 ORM:
```
// 透過 db.Table("寫死Table名稱")
db.Table("GoodUser").Take(&gu)
```

### 改用自定義 Column name

根據[官方文件](https://Gorm.io/docs/conventions.html#Column-Name)，想要自定義 column name，可以在 field 後面加上 tag :

```go
// 這是個Model
type GoodUser struct {
  // 注意 tag
  LaAge int `Gorm:"column:LaAge"`
}

// 宣告一個物件
gu := GoodUser{18}

// 執行查詢
db.Select("LaAge").Take(&gu)

// 實際執行SQL
// SELECT "LaAge" FROM "good_users" LIMIT 1

// Table name:  GoodUser -> good_users ,with 's'
// Column name: LaAge -> LaAge 被定義了 
```

### 通用解法
如果現行資料庫有一套命名規則，只是跟 `Gorm` 的規則相左，那麼透過 `Gorm Config` 直接改 [NamingStrategy](https://Gorm.io/docs/Gorm_config.html#NamingStrategy) 應該是最通用的辦法了。

---
## Field name 必須開頭大寫
有些資料庫的欄位名稱偏好小寫，那我想要 `Model` 裡面的 `Field` 也一模一樣行嗎？很抱歉會出錯，因為小寫開頭的 `Field` 是不會被別的 `package` 看到的。

---
## 以為 Gorm 的 first() 就是SQL 的 limit 1 或是 top 1
不知為何官方文件一直用 first 做例子，因為其實這動作有點耗運算資源。

  *The First and Last methods will find the first and last record (respectively) as ordered by primary key.*

也就是說，我今天如果只是想要拿一筆資料來看看，照著官方的舉例一直用 `first` ，其實每次都會做排序，也就是 `order by`，如果有 `primary key (PK)`它就用PK排序，沒有的話它就用第一個欄位。

如果想要 `SQL` 的 `limit 1` 或是 `top 1` 的效果，也就是不論排序給我一個就對了，那得用 `take` 。
```go
type GoodUser struct {
  LaDeSai string
  LaAge   int
}
gu := GoodUser{"What's up", 18}

db.Select("LaAge").First(&gu)
// SELECT "la_age" FROM "good_users" ORDER BY "good_users"."la_de_sai" LIMIT 1
// 注意那個 ORDER BY

db.Select("LaAge").Take(&gu)
// SELECT "la_age" FROM "good_users" LIMIT 1
// 沒有 ORDER BY
```
建議任何查詢動作，除非有十足把握，使用 Gorm 時都先用 [DryRun Mode](https://Gorm.io/docs/sql_builder.html#DryRun-Mode) 看看到底產生了麼SQL code，避免 Gorm 太貼心，導致 DBA 衝進來殺人。


