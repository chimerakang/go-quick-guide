# Todo 重新開發，利用docker-compose 
by [@chimerakang](https://github.com/chimerakang)

---

[docker-compose.yml](../demos/docker/go-todo/docker-compose.yml)
```yaml
version: '3.8'

services:
  web:
    build: ./nginx
    ports:
      - "80:80"
    depends_on:
      - api
    networks:
      - "mynet"

  api:
    build:
      context: .
      dockerfile: Dockerfile.stage
    ports:
      - "9999:9999"
    environment:
      - DB_HOST=mysql
      - DB_USER=root
      - DB_PASS=secret
      - DB_NAME=gotodo
      - DB_PORT=3306
      - PORT=9999
    depends_on:
      - mysql
    networks:
      - "mynet"
    restart: on-failure

  mysql:
    platform: "linux/x86_64"
    image: "mysql:5.7"
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: gotodo
    ports:
      - "3306:3306"
    networks:
      - "mynet"
    volumes:
      - mysql_data:/var/lib/mysql      

networks:
  mynet:
    driver: "bridge"
volumes:
  mysql_data:    
```

## Nginx 設定
示範透過nginx這個服務做轉址的設定，如果透過將` todo `所有` request `轉換到 `api:9999` 的內部api

[nginx config](../demos/docker/go-todo/nginx/nginx.conf)
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

    location /todo/ {
      proxy_set_header X-Forwarded-For $remote_addr;
      proxy_set_header Host            $http_host;
      proxy_pass http://api:9999/;
      # 處理重定向
      #   proxy_redirect off;
      proxy_redirect / /todo/;
      # 移除 /todo 前缀
      rewrite ^/todo(.*)$ $1 break;
    }

 # 處理 /add, /complete, /delete 等路徑
    location ~ ^/(add|complete|delete|login|logout) {
        if ($request_method = POST) {
            proxy_pass http://api:9999$request_uri;
        }
        if ($request_method = GET) {
            return 301 /todo$request_uri;
        }
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header Host $http_host;
        proxy_redirect / /todo/;
    }    
  }
}

```



