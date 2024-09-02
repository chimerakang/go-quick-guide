# 基於CA 的TLS 憑證認證

by [@chimerakang](https://github.com/chimerakang)

---
## 介紹
在上一章中，我們提出了一個問題。就是如何確保證書的可靠性和有效性？你要如何確定你Server、Client 的憑證是對的呢？

## CA
為了確保證書的可靠性和有效性，這裡可引入CA 頒發的根證書的概念。其遵守X.509 標準

### 根證書
根憑證（root certificate）是屬於根憑證授權單位（CA）的公鑰憑證。我們可以透過驗證CA 的簽章從而信任CA ，任何人都可以得到CA 的憑證（含公鑰），用以驗證它所簽發的憑證（客戶端、服務端）

它所包含的文件如下：

* 公鑰

* 金鑰

### 產生Key
```
openssl genrsa -out ca.key 2048
```
### 產生金鑰
```
openssl req -new -x509 -days 7200 -key ca.key -out ca.pem
```
### 填寫資訊
```
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:go-grpc-example
Email Address []:
```
### Server
#### 產生CSR
```
openssl req -new -key server.key -out server.csr
```
#### 填寫資訊
```
Country Name (2 letter code) [AU]:
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (eg, fully qualified host name) []:go-grpc-example
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
```
CSR 是Cerificate Signing Request 的英文縮寫，為憑證要求檔案。主要作用是CA 會利用CSR 檔案進行簽章使得攻擊者無法偽裝或竄改原有憑證

#### 基於CA 簽發
```
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.pem
```

### Client
#### 產生Key
```
openssl ecparam -genkey -name secp384r1 -out client.key
```
#### 產生CSR
```
openssl req -new -key client.key -out client.csr
```
#### 基於CA 簽發
```
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem
```
### 整理目錄
至此我們產生了一堆文件，請依照以下目錄結構存放：

```
conf
├── ca.key
├── ca.pem
├── ca.srl
├── client
│   ├── client.csr
│   ├── client.key
│   └── client.pem
└── server
    ├── server.csr
    ├── server.key
    └── server.pem
```    
另外有一些文件是不應該出現在倉庫內，應保密或刪除的。但為了真實示範所以保留著（敲黑板）

## gRPC
接下來將正式開始針對gRPC 進行編碼，改造上一章的程式碼。目標是基於CA 進行TLS 認證🤫

### Server
```go
package main

import (
    "context"
    "log"
    "net"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    pb "github.com/chimerakang/go-grpc-example/proto"
)

...

const PORT = "9001"

func main() {
    cert, err := tls.LoadX509KeyPair("../../conf/server/server.pem", "../../conf/server/server.key")
    if err != nil {
        log.Fatalf("tls.LoadX509KeyPair err: %v", err)
    }

    certPool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("../../conf/ca.pem")
    if err != nil {
        log.Fatalf("ioutil.ReadFile err: %v", err)
    }

    if ok := certPool.AppendCertsFromPEM(ca); !ok {
        log.Fatalf("certPool.AppendCertsFromPEM err")
    }

    c := credentials.NewTLS(&tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientAuth:   tls.RequireAndVerifyClientCert,
        ClientCAs:    certPool,
    })

    server := grpc.NewServer(grpc.Creds(c))
    pb.RegisterSearchServiceServer(server, &SearchService{})

    lis, err := net.Listen("tcp", ":"+PORT)
    if err != nil {
        log.Fatalf("net.Listen err: %v", err)
    }

    server.Serve(lis)
}
```
* tls.LoadX509KeyPair()：從憑證相關檔案讀取和解析訊息，得到憑證公鑰、金鑰對

```go
func LoadX509KeyPair(certFile, keyFile string) (Certificate, error) {
    certPEMBlock, err := ioutil.ReadFile(certFile)
    if err != nil {
        return Certificate{}, err
    }
    keyPEMBlock, err := ioutil.ReadFile(keyFile)
    if err != nil {
        return Certificate{}, err
    }
    return X509KeyPair(certPEMBlock, keyPEMBlock)
}
```
* x509.NewCertPool()：建立一個新的、空的CertPool

* certPool.AppendCertsFromPEM()：嘗試解析所傳入的PEM 編碼的憑證。如果解析成功會將其加到CertPool 中，以便於後面的使用

* credentials.NewTLS：建構基於TLS 的TransportCredentials 選項

* tls.Config：Config 結構用於配置TLS 用戶端或伺服器

在Server，共使用了三個Config 設定項：

（1）Certificates：設定憑證鏈，允許包含一個或多個

（2）ClientAuth：要求必須校驗客戶端的憑證。可依實際情況選用以下參數：

```go
const (
    NoClientCert ClientAuthType = iota
    RequestClientCert
    RequireAnyClientCert
    VerifyClientCertIfGiven
    RequireAndVerifyClientCert
)
```
（3）ClientCAs：設定根憑證的集合，校驗方式使用ClientAuth 中設定的模式

### Client
```go
package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    pb "github.com/chimerakang/go-grpc-example/proto"
)

const PORT = "9001"

func main() {
    cert, err := tls.LoadX509KeyPair("../../conf/client/client.pem", "../../conf/client/client.key")
    if err != nil {
        log.Fatalf("tls.LoadX509KeyPair err: %v", err)
    }

    certPool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("../../conf/ca.pem")
    if err != nil {
        log.Fatalf("ioutil.ReadFile err: %v", err)
    }

    if ok := certPool.AppendCertsFromPEM(ca); !ok {
        log.Fatalf("certPool.AppendCertsFromPEM err")
    }

    c := credentials.NewTLS(&tls.Config{
        Certificates: []tls.Certificate{cert},
        ServerName:   "go-grpc-example",
        RootCAs:      certPool,
    })

    conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
    if err != nil {
        log.Fatalf("grpc.Dial err: %v", err)
    }
    defer conn.Close()

    client := pb.NewSearchServiceClient(conn)
    resp, err := client.Search(context.Background(), &pb.SearchRequest{
        Request: "gRPC",
    })
    if err != nil {
        log.Fatalf("client.Search err: %v", err)
    }

    log.Printf("resp: %s", resp.GetResponse())
}
```
在Client 中絕大部分與Server 一致，不同點的地方是，在Client 請求Server 端時，Client 端會使用根憑證和ServerName 去對Server 端進行校驗

簡單流程大致如下：

1. Client 透過請求得到Server 端的憑證

2. 使用CA 認證的根憑證對Server 端的憑證進行可靠性、有效性等校驗

3. 校驗ServerName 是否可用、有效

當然了，在設定了`tls.RequireAndVerifyClientCert`模式的情況下，Server 也會使用CA 認證的根憑證對Client 端的憑證進行可靠性、有效性等校驗。也就是兩邊都會進行校驗，極大的保證了安全性👍

## 驗證
重新啟動server.go 和執行client.go，查看回應結果是否正常

## 總結
在本章節，我們使用CA 頒發的根憑證對客戶端、服務端的憑證進行了簽發。進一步的提高了兩者的通訊安全

這回是真的大功告成了！

---
## Next: [讓你的服務同時提供HTTP 接口](grpc4.md)