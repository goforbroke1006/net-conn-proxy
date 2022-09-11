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
