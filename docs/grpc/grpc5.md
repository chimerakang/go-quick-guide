# 對RPC 方法做自訂認證
by [@chimerakang](https://github.com/chimerakang)

---
## 介紹
在前面的章節中，我們介紹了兩種（憑證算一種）可全域認證的方法：

* TLS 憑證認證

* 基於CA 的TLS 憑證認證

* Unary and Stream interceptor

而在實際需求中，常常會對某些模組的RPC 方法做特殊認證或校驗。今天將會講解、實現這塊的功能點

## 課前知識
```go
type PerRPCCredentials interface {
    GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
    RequireTransportSecurity() bool
}
```
在gRPC 中預設定義了PerRPCCredentials，它就是本章節的主角，是gRPC 預設提供用於自訂認證的接口，它的作用是將所需的安全認證資訊添加到每個RPC 方法的上下文中。其包含2 個方法：

* GetRequestMetadata：取得目前請求認證所需的元資料（metadata）

* RequireTransportSecurity：是否需要基於TLS 認證進行安全傳輸

## 目錄結構
新建simple_token_server/server.go 和simple_token_client/client.go，目錄架構如下：

```
go-grpc-example
├── client
│   ├── simple_client
│   ├── simple_http_client
│   ├── simple_token_client
│   └── stream_client
├── conf
├── pkg
├── proto
├── server
│   ├── simple_http_server
│   ├── simple_server
│   ├── simple_token_server
│   └── stream_server
└── vendor
```
## gRPC
### Client
```go
package main

import (
    "context"
    "log"

    "google.golang.org/grpc"

    "github.com/chimerakang/go-grpc-example/pkg/gtls"
    pb "github.com/chimerakang/go-grpc-example/proto"
)

const PORT = "9004"

type Auth struct {
    AppKey    string
    AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
    return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
    return true
}

func main() {
    tlsClient := gtls.Client{
        ServerName: "go-grpc-example",
        CertFile:   "../../conf/server/server.pem",
    }
    c, err := tlsClient.GetTLSCredentials()
    if err != nil {
        log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
    }

    auth := Auth{
        AppKey:    "chimerakang",
        AppSecret: "20221005",
    }
    conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(&auth))
    ...
}
```
在Client 端，重點實現type PerRPCCredentials interface所需的方法，專注於兩點即可：

* struct Auth：GetRequestMetadata、RequireTransportSecurity

* grpc.WithPerRPCCredentials

### Server
```go
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"

    "github.com/chimerakang/go-grpc-example/pkg/gtls"
    pb "github.com/chimerakang/go-grpc-example/proto"
)

type SearchService struct {
    auth *Auth
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
    if err := s.auth.Check(ctx); err != nil {
        return nil, err
    }
    return &pb.SearchResponse{Response: r.GetRequest() + " Token Server"}, nil
}

const PORT = "9004"

func main() {
    ...
}

type Auth struct {
    appKey    string
    appSecret string
}

func (a *Auth) Check(ctx context.Context) error {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return status.Errorf(codes.Unauthenticated, "Token failed")
    }

    var (
        appKey    string
        appSecret string
    )
    if value, ok := md["app_key"]; ok {
        appKey = value[0]
    }
    if value, ok := md["app_secret"]; ok {
        appSecret = value[0]
    }

    if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
        return status.Errorf(codes.Unauthenticated, "Token failed")
    }

    return nil
}

func (a *Auth) GetAppKey() string {
    return "chimerakang"
}

func (a *Auth) GetAppSecret() string {
    return "20221005"
}
```
在Server 端就更簡單了，實際上就是呼叫`metadata.FromIncomingContext`從上下文中取得metadata，再在不同的RPC 方法中進行認證檢查

## 驗證
重新啟動server.go 和client.go，得到以下結果：

```
$ go run client.go
$ resp: gRPC Token Server
```
修改client.go 的值，製造兩者不一致，得到無效結果：

```
$ go run client.go
2024/08/30 18:36:03 client.Search err: rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: tls: failed to verify certificate: x509: certificate relies on legacy Common Name field, use SANs instead"
exit status 1
```
## 一個個加太麻煩
我相信你一定會問一個個加，也太麻煩了吧？有這個想法的你，應當把type PerRPCCredentials interface做成一個攔截器（interceptor）

## 總結
本章節比較簡單，主要是針對RPC 方法的自訂認證進行了介紹，如果是想做全局的，建議是舉一反三從攔截器下手哦。
