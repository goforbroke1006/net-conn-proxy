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
net-conn-proxy -p tcp -d 127.0.0.1:8080 -u example.com:80 -b 64
```

Open downstream address [http://127.0.0.1:8080](http://127.0.0.1:8080) in browser.

##### Proxy for L4D2 server

On machine for proxy and game server:

1. Run L4D2 server on local machine on non-default port **47015**:
    
    ```shell
    docker run --rm -it \
      -p 47015:27015/tcp \
      -p 47015:27015/udp \
      -e L4D2_SERVER_HOSTNAME="L4D2.server.docker-compose.local" \
      -e L4D2_SERVER_SV_REGION=3 \
      -e L4D2_SERVER_MAX_PLAYERS=4 \
      --name l4d2 \
      goforbroke1006/l4d2-server:latest
    
    ```

2. Run separate proxies for TCP and for UDP connections:
    
    ```shell
    net-conn-proxy -p tcp -d 0.0.0.0:27015 -u 127.0.0.1:47015 -b 2048
    ```
    
    ```shell
    net-conn-proxy -p udp -d 0.0.0.0:27015 -u 127.0.0.1:47015 -b 64 >> ./udp-27015-47015.log
    ```

3. Check IP in LAN

```shell
ifconfig
# my IP in wlan is 192.168.0.9
```

On machine with Steam and client of L4D2:

1. Run L4D2 game
2. Press ~ button to open developer console
3. Type **connect 192.168.0.9:27015**
