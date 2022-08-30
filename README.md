# proxy

```shell
make
make install
```

```shell
goforbroke1006-proxy-tcp --upstream www.bing.com:80 --downstream 0.0.0.0:8080
```

```shell
echo -e "GET /search?q=hello-world HTTP/2\nUser-Agent: curl/7.54.0\nAccept: */*" | nc 127.0.0.1 8080
```
