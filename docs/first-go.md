# 第一個Golang程式
by [@chimerakang](https://github.com/chimerakang)

---
## 關鍵特性

* 開發快
* 運行快
* 高並發
* 易學習

---

## 開發快

* 強大的函示庫(lib)支持:函式庫很多，可以直接使用 GitHub 上面的函式庫
* 編譯速度很快（模組機制）:帶動整個開發的流程更快速
* 易於分析、易調整（靜態類型）:很多程式的錯誤在編譯期就會挑出來，相對易於除錯
* 易佈署（靜態鏈接，最小化依賴）
* 多傳回值:你函式的回傳值可以是多個
* 跨平台的:只要程式中不碰到 C 函式庫，在 Windows (或 Mac) 寫好的 Golang 網頁程式，可以不經修改就直接發布在 GNU/Linux 伺服器上


## 運行快

* 直接編譯為本機程式碼，類似C/C++
	* （ELF on Unix/PE on Windows）

## 高併發

* 語法級的協程(goroutine)/管道(channel) 支持
* goroutine 比起傳統的執行緒 (thread) 來說輕量得多，在高負載時所需開銷更少

---

## 易學習

* 自帶GC（無記憶體洩漏）
* 沒有class
	* 沒有繼承、多型(OOP, Object-Orentied Programming)
	* 有struct
		* 支援抽象(abstraction)
		* 支援封裝(encapsulation)
	* 有method (Object-Based)
* 沒有template，沒有泛型程式設計（GP, Generic Programming）
	* 有interface

---

## Go 指令
在終端機上執行 go 指令就會看到一系列的指令介紹：

```
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report   
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        work        workspace maintenance
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildconstraint build constraints
        buildmode       build modes
        c               calling between Go and C
        cache           build and test caching
        environment     environment variables
        filetype        file types
        go.mod          the go.mod file
        gopath          GOPATH environment variable
        gopath-get      legacy GOPATH go get
        goproxy         module proxy protocol
        importpath      import path syntax
        modules         modules, module versions, and more
        module-get      module-aware go get
        module-auth     module authentication using go.sum
        packages        package lists and patterns
        private         configuration for downloading non-public code
        testflag        testing flags
        testfunc        testing functions
        vcs             controlling version control with GOVCS

Use "go help <topic>" for more information about that topic.
```

如果對於某個指令特別需要幫助可以用 go help [topic] 指令。

這邊特別介紹的有四個指令：go run、go build、go install、go clean。

### go run
直接執行 go code。
```
$ go run demos/helloworld/main.go 
Hello World: 2024
```

### go build
```
go build
```

這個命令主要用於編譯程式碼。在套件的編譯過程中，若有必要，會同時編譯與之相關聯的套件。

如果是 main 套件，當你執行go build之後，它就會在當前目錄下產生一個可執行檔案，或者使用go build -o 路徑/a.exe

如果某個專案資料夾下有多個檔案，而你只想編譯某個檔案，就可在go build之後加上檔名，例如go build a.go；go build命令預設會編譯當前目錄下的所有 go 檔案。


參數的介紹

-o 指定輸出的檔名，可以帶上路徑，例如 go build -o a/b/c

-i 安裝相應的套件，編譯+go install

-a 更新全部已經是最新的套件的，但是對標準套件不適用

-n 把需要執行的編譯命令顯示出來，但是不執行，這樣就可以很容易的知道底層是如何執行的

-p n 指定可以並行可執行的編譯數目，預設是 CPU 數目

-race 開啟編譯的時候自動檢測資料競爭的情況，目前只支援 64 位的機器

-v 顯示出來我們正在編譯的套件名

-work 顯示出來編譯時候的臨時資料夾名稱，並且如果已經存在的話就不要刪除

-x 顯示出來執行的命令，其實就是和-n的結果類似，只是這個會執行

-ccflags 'arg list' 傳遞參數給 5c, 6c, 8c 呼叫

-compiler name 指定相應的編譯器，gccgo 還是 gc

-gccgoflags 'arg list' 傳遞參數給 gccgo 編譯連線呼叫

-gcflags 'arg list' 傳遞參數給 5g, 6g, 8g 呼叫

-installsuffix suffix 為了和預設的安裝套件區別開來，採用這個字首來重新安裝那些依賴的套件，-race的時候預設已經是-installsuffix race，大家可以透過-n命令來驗證

-ldflags 'flag list' 傳遞參數給 5l, 6l, 8l 呼叫

如果沒有錯誤就產生執行檔於當前目錄。
build 後產生的檔案即是執行檔。

### go install
如果沒有錯誤則產生執行檔於 $GOPATH/bin。
```
go install

# 查看執行檔
ls $GOPATH/bin

> helloWorld
```

### go clean
執行後會將 build 產生的檔案都刪除。(install的不會刪)
```
go clean
```

## 介紹
“你好，世界！”程式是電腦程式設計中經典且歷史悠久的傳統。對於初學者來說，這是一個簡單而完整的第一個程序，並且是確保正確配置環境的好方法。

本教學將引導您完成在 Go 中建立此程式的過程。但是，為了使程式更有趣，您將修改傳統的“Hello, World!”，以便它詢問用戶的姓名。然後您將在問候語中使用該名字。完成本教學後，運行時您將得到一個如下所示的程式：
```
Please enter your name.
Sammy
Hello, Sammy! I'm Go!
```

greeting.go
```go
package main

import (
	"fmt"
	"strings"    
)

func main() {
	fmt.Println("Please enter your name.")
	var name string
	fmt.Scanln(&name)
	name = strings.TrimSpace(name)
	fmt.Printf("Hi, %s! I'm Go!\n", name)
}
```

## 讓我們分解程式碼的不同組成部分。

*package*是一個關鍵字，定義該檔案屬於哪個程式package。每個資料夾只能有一個package，並且每個.go文件必須在其文件頂部聲明相同的package名稱。在此範例中，程式碼屬於main套件。

---

*import*是一個關鍵字，告訴 Go 編譯器您想在此檔案中使用哪些其他套件。這裡導入fmt和strings標準庫套件。fmt提供了開發時有用的格式化和列印功能，strings則提供了字串相關的功能

---

```
name = strings.TrimSpace(name)
```
要從外部輸入獲取字串使用 strings 標準函示庫的 TrimSpaceGo 函數。可以將字串的開頭和結尾刪除所有空格字符，包括換行符。在這種情況下，它會刪除按enter時建立的字串末尾的換行符

---

fmt.Printf是一個 Go 函數，可以在fmt package中找到，它告訴電腦將一些文字列印到螢幕上。

您可以fmt.Printf在函數後面加上一系列字符，例如"Hello, World!"用引號括起來的 。引號內的任何字元都稱為字串。當程式運行時，該fmt.Println函數會將這個字串列印到螢幕上。

```
$ go run ./demos/greeting/main.go
Please enter your name.
chimera
Hi, chimera! I'm Go!
```

您現在已經撰寫玩第一個 Go 程式，它接受用戶的輸入並將其列印回螢幕。

---



## 結論
在本教程中，您編寫了“Hello, World!”程式接受使用者輸入、處理結果並顯示輸出。這個作業將涵蓋基本的輸入/輸出操作、字串操作，以及簡單的控制流程。以下是作業的描述和一些指導:
### 作業題目: 
創建一個簡單的 Golang Madlib 程式
#### 目標:
創建一個程式，要求用戶輸入幾個詞語，然後將這些詞語插入到預定義的故事模板中，最後輸出完整的故事
#### 要求:

使用 fmt 包來處理輸入和輸出。
創建一個包含至少 5 個空缺的短故事模板。
提示用戶輸入不同類型的詞語(如名詞、動詞、形容詞等)。
將用戶輸入的詞語插入到故事模板中。
輸出最終的故事。

提示:

使用 fmt.Print() 或 fmt.Println() 來輸出提示和最終故事。
使用 fmt.Scan() 來獲取用戶輸入。
可以使用字符串連接或 fmt.Sprintf() 來構建最終的故事字符串。

示例結構:
```go
package main

import "fmt"

func main() {
    // 定義變量來存儲用戶輸入
    var noun1, verb1, adjective1 string
    
    // 提示用戶並獲取輸入
    fmt.Print("請輸入一個名詞: ")
    fmt.Scan(&noun1)
    
    // ... 獲取更多輸入 ...
    
    // 構建故事
    story := fmt.Sprintf("從前有一個%s...", noun1)
    
    // 輸出最終故事
    fmt.Println("\n這是您的Madlib故事:")
    fmt.Println(story)
}
```

#### 擴展挑戰:

* 添加錯誤處理,確保用戶輸入有效。
* 允許用戶選擇多個預定義的故事模板。

這個作業將幫助熟悉 Go 的基本語法和概念,同時創建一個有趣的互動程式。

[Build Project Envirement](./build-project.md)