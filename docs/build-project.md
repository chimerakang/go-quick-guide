# 如何建置 Go 專案
by [@chimerakang](https://github.com/chimerakang)

---

## 介紹

您已經使用該命令go run自動編譯原始程式碼並執行生成的可執行檔。儘管此命令對於在命令列上測試程式碼很有用，但分發或部署應用程式需要您將程式碼建置為可共享的二進位可執行文件，或包含可以運行應用程式的機器位元組程式碼的單一檔案。為此，您可以使用 Go 工具鏈來建立和安裝您的程式。

在 Go 中，將原始程式碼轉換為二進位可執行檔的過程稱為建置。建立此可執行檔後，它將不僅包含您的應用程序，還包含在目標平台上執行二進位檔案所需的所有支援程式碼。這意味著 Go 二進位檔案不需要係統依賴項（例如 Go 工具）即可在新系統上運作。將這些可執行檔放入您自己系統上的可執行檔路徑中將允許您從系統上的任何位置執行該程式。這與將程式安裝到系統上是一樣的。

在本教程中，您將使用 Go 工具鏈來運行、建置和安裝範例 greeting，使您能夠有效地使用、分發和部署未來的應用程式。

---
## 第 1 步 建立 Go 模組
Go 程式和函式庫是圍繞模組的核心概念建構的。模組包含有關程式使用的package以及要使用的這些package的版本的資訊。

為了告訴 Go 這是一個 Go 模組，您需要使用以下命令建立一個 Go 模組go mod：
```
$ cd demos/greeting/
$ go mod init greeting
go: creating new go.mod: module greeting
go: to add module requirements and sums:
        go mod tidy

$ go mod tidy
```
如果將來模組的需求發生變化，Go 將提示您運行go mod tidy以更新該模組的需求。現在運行它不會產生額外的效果。

## 第 2 步 — 建置 Go 執行檔
使用go build，您可以為我們的範例 Go 應用程式產生可執行檔，從而允許您將程式分發和部署到您想要的位置。

嘗試這個與main.go.在您的greeting目錄中，執行以下命令：
```
chime@DESKTOP-Q02LKF0 MINGW64 /g/GoProjects/src/go-quick-guide/demos/greeting (master)
$ go build

chime@DESKTOP-Q02LKF0 MINGW64 /g/GoProjects/src/go-quick-guide/demos/greeting (master)
$ ls -l
total 2026
-rw-r--r-- 1 chime 197609      27 Aug  4 13:59 go.mod       
-rwxr-xr-x 1 chime 197609 2071552 Aug  4 14:02 greeting.exe*
-rw-r--r-- 1 chime 197609     219 Aug  4 13:05 main.go 
```

預設情況下go build將為目前平台和體系結構產生可執行檔。例如，如果在某個linux/386系統上構建，則可執行檔將與任何其他系統相容linux/386，即使未安裝 Go。 Go 支援針對其他平台和架構進行構建，您可以在我們的為不同作業系統和架構構建 Go 應用程式文章中閱讀更多相關資訊。

現在，您已經建立了可執行文件，運行它以確保二進位檔案已正確建置。在 macOS 或 Linux 上，執行以下命令：
```
$ ./greeting
```
在 Windows 上，運行：
```
$ greeting.exe
```

在下一節中，本教學將解釋執行檔的命名方式以及如何變更它，以便您可以更好地控製程式的建置過程。

## 第 3 步 — 更改執行檔名稱
現在您已經知道如何產生可執行文件，下一步是確定 Go 如何為二進位檔案選擇名稱並為您的專案自訂該名稱。

當您執行時go build，預設情況下 Go 會自動決定產生的可執行檔的名稱。它透過使用您之前建立的模組來完成此操作。運行 go mod init greeting 時，它創建了一個名為“greeting”的模組，這就是為什麼產生的執行檔案名稱是greeting的原因。

讓我們仔細看看 module 方法。如果您go.mod的專案中有一個帶有module如下聲明的文件：

```
$ cat go.mod 
module greeting

go 1.21.5
```

修改 go.mod
```
module shark
go 1.21.5
```
那麼產生的可執行檔的預設名稱將變成shark.


在需要特定命名約定的更複雜的程式中，這些預設值並是命名檔案的最佳選擇。在這些情況下，最好使用 -o 參數自訂輸出。

若要對此進行測試，請將上一節中建立的可執行檔的名稱變更為hello並將其放置在名為的子資料夾中bin。您不必建立此資料夾； Go 將在建置過程中自行完成此操作。

執行以下：
```
$ go build -o ../bin/hello
$ ls -l ../bin/
total 4048
-rwxr-xr-x 1 chime 197609 2071552 Aug  4 14:13 hello*
```

在本例中，在指定的目錄當中建立新的名為 hello 的新執行檔

若要測試新的可執行文件，請變更到新目錄並執行檔案：

---
## Next

[入門-1](./introduce-1.md)

---

## Prev
[First go](./first-go.md)