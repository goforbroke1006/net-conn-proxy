# net-conn-proxy

Proxy application to print client-server communication details.

### Usage

At first, install to ${GOBIN}

```shell
make
make install
```

##### Proxy for HTTP site

Install  and run:

```shell
net-conn-proxy -p tcp -d 127.0.0.1:8080 -u example.com:80 -bs 64
```

Open downstream address [http://127.0.0.1:8080](http://127.0.0.1:8080) in browser.

##### Proxy for L4D2 server

Run L4D2 server on local machine:

```shell
docker run --rm -p 27015:27015/tcp -p 27015:27015/udp --name l4d2 left4devops/l4d2
```

```shell
ifconfig

# my IP in wlan (local WIFI network over home router)
# 192.168.0.9
```

Separate proxies for TCP and for UDP connections:

```shell
net-conn-proxy -p tcp -d 0.0.0.0:47015 -u 127.0.0.1:27015 -bs 2048
```

```shell
net-conn-proxy -p udp -d 0.0.0.0:47015 -u 127.0.0.1:27015 -bs 2048
```
