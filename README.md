# net-conn-proxy

Proxy application to print client-server communication details.

### Usage

Install to ${GOBIN} and run:

```shell
make
make install

net-conn-proxy -p tcp -d 127.0.0.1:8080 -u www.bing.com:80
```

Open downstream address [https://127.0.0.1:8080](https://127.0.0.1:8080) in browser.
