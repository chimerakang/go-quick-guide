# 撰寫第一個 docker-compose
by [@chimerakang](https://github.com/chimerakang)

---
## 為何需要 Docker-Compose ?
使用多個容器可讓您將容器專用於特製化工作。 每個容器都應該執行一項操作並妥善處理。

以下是您可能需要使用多容器應用程式的一些原因：

* 將容器分開，以不同於資料庫的方式管理 API 和前端。
* 容器可讓您以隔離的方式設定版本及更新版本。
* 雖然您可能會在本機使用資料庫的容器，但您可能需要針對生產中的資料庫使用受控服務。
* 執行多個處理序需要處理序管理員，這會增加容器啟動/關機的複雜性。

## 建立 Docker Compose 檔案
Docker Compose 可協助定義及共用多容器應用程式。 使用 Docker Compose，您可以建立檔案來定義服務。 您可以使用單一命令來啟動所有項目，或將其全部拆毀。

您可以在檔案中定義應用程式堆疊，並在版本控制下，將該檔案保留在專案存放庫的根目錄。 此方法可讓其他人參與您的專案。 他們只需要複製您的存放庫。

1. 在應用程式專案的根目錄中，建立名為 docker-compose.yml 的檔案。

2. 在撰寫檔案中，從定義結構描述版本開始。
    ```YAML
    version: "3.7"
    ```
    在大部分情況下，最好使用最新的支援版本。

3. 定義您要做為應用程式一部分執行的服務或容器。
    ```YAML
    version: "3.7"

    services:
    ```
4. 定義容器的服務項目和映像。    
    ```YAML
    version: "3.7"

    services:
    app:
        image: node:20-alpine
    ```    
    您可以為服務挑選任何名稱。 名稱會自動變成網路別名，這在定義 MySQL 服務時很有用。

5. 新增命令。
    ```YAML
    version: "3.7"

    services:
    app:
        image: node:20-alpine
        command: sh -c "yarn install && yarn run dev"
    ```

6. 指定服務的連接埠，此連接埠會對應至上述命令中的 -p 3000:3000。
    ```YAML
    version: "3.7"

    services:
    app:
        image: node:20-alpine
        command: sh -c "yarn install && yarn run dev"
        ports:
        - 3000:3000
    ```      

7. 指定工作目錄和磁碟區對應
    ```
    YAML
    version: "3.7"

    services:
    app:
        image: node:20-alpine
        command: sh -c "yarn install && yarn run dev"
        ports:
        - 3000:3000
        working_dir: /app
        volumes:
        - ./:/app
    ```      
    在 Docker Compose 磁碟區定義中，您可以使用來自目前目錄的相對路徑。

8. 指定環境變數定義。
    ```YAML
    version: "3.7"

    services:
    app:
        image: node:20-alpine
        command: sh -c "yarn install && yarn run dev"
        ports:
        - 3000:3000
        working_dir: /app
        volumes:
        - ./:/app
        environment:
        MYSQL_HOST: mysql
        MYSQL_USER: root
        MYSQL_PASSWORD: <your-password>
        MYSQL_DB: todos
    ```

9. 新增 MySQL 服務的定義。 以下是您在上方使用的命令：
    ```YAML
    version: "3.7"

    services:
    app:
        # The app service definition
    mysql:
        image: mysql:5.7
    ```    
    服務會自動取得網路別名。 指定要使用的映像。

10. 定義磁碟區對應。

    使用與 services: 相同層級的 volumes: 區段指定磁碟區。 指定映像下的磁碟區對應。

    ```YAML

    version: "3.7"

    services:
    app:
        # The app service definition
    mysql:
        image: mysql:5.7
        volumes:
        - todo-mysql-data:/var/lib/mysql

    volumes:
    todo-mysql-data:
    ```

11. 指定環境變數。
    ```YAML
    version: "3.7"

    services:
    app:
        # The app service definition
    mysql:
        image: mysql:5.7
        volumes:
        - todo-mysql-data:/var/lib/mysql
        environment: 
        MYSQL_ROOT_PASSWORD: <your-password>
        MYSQL_DATABASE: todos

    volumes:
    todo-mysql-data:
    ```    
此時，完整的 docker-compose.yml 看起來會像這樣：
```YAML
version: "3.7"

services:
  web:
    build: ./nginx
    ports:
      - "80:80"
    depends_on:
      - api
  api:
    build:
      context: .
      dockerfile: Dockerfile.stage
    ports:
      - "9999:9999"
```  
## 執行 docker-compose
現在您已擁有 docker-compose.yml 檔案
要為此 Golang API 建立 Docker image，我們在應用程式的目錄中執行以下命令
```
docker-compose build api
```
我們將對運行容器執行相同的操作。我只想提一下，您可以使用此命令來運行容器
```
docker-compose up api
```

### 設定 Nginx 伺服器（當然在 Docker 中）
讓我們將 Nginx 檔案放入nginx應用程式目錄中命名的新目錄中。在nginx目錄中，建立一個新文件，nginx.conf並使用以下配置命名：
```yaml
#nginx/nginx.conf
events {
    worker_connections 1024;
}
http {
  server_tokens off;
  server {
    listen 80;
    root  /var/www;

    location / {
      index index.html;
    }

    # Media: images, icons, video, audio, HTC
    location ~* \.(?:jpg|jpeg|gif|png|ico|cur|gz|svg|svgz|mp4|ogg|ogv|webm|htc)$ {
      expires 1d;
      access_log off;
      add_header Cache-Control "public";
    }

    # CSS and Javascript
    location ~* \.(?:css|js)$ {
      expires 1d;
      access_log off;
      add_header Cache-Control "public";
    }

    location /api/ {
      proxy_set_header X-Forwarded-For $remote_addr;
      proxy_set_header Host            $http_host;
      proxy_pass http://api:9999/;
    }
  }
}
```

此組態設定一個 Nginx 伺服器來偵聽連接埠 80 並將所有`/api`請求轉送到在連接埠 9999 上執行的 Golang API http://api:9999。它是 docker-compose.yml 檔案中定義的 Golang API 服務的名稱。

若要在 Docker 容器中建置 Nginx 映像，請在「nginx」目錄中建立一個名為「Dockerfile」的新檔案：
```yaml
# Official Nginx image (Yes, in this article I always use the `latest`. Kill me!)
FROM nginx:latest

# Copy Nginx configuration file to the container
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80
```

再一次：

* 該 Dockerfile 取得官方 Nginx 運行時映像作為 parent image。
* 將 Nginx 設定檔複製到容器中。
* port 80。
對於建置和執行，您可以使用與 Golang API 容器相同的命令，但使用名稱服務web。
```
docker-compose build web
docker-compose up api web
```
如果您執行 golang 和 nginx 服務，您可以在瀏覽器中檢查。`http://localhost/api/hello?name=john`


下面的指令可以全部執行

```Bash
docker-compose up -d
```
-d 參數會使命令在背景中執行。

您應該會看到類似下列結果的輸出。
```
✔ Network docker-1_default  Created                                                                                                                            0.1s 
 ✔ Container docker-1-api-1  Started                                                                                                                            0.7s 
 ✔ Container docker-1-web-1  Started  
 ```

使用完這些容器時，只要移除這些容器即可。

請在命令列執行 
```
docker-compose down
```
容器隨即停止。 網路已移除。


---
## Next : [Todo with docker-compose](./docker-compose-2.md)