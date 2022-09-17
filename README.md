# net-conn-proxy

Proxy application to print client-server communication details.

Supported protocols:
* HTTP
* TCP
* UDP

### Usage

At first, install util

```shell
make
make install
```

##### Proxy for HTTP site

Install and run:

```shell
net-conn-proxy http -d 127.0.0.1:8080 -o example.com:80 -b 64
```

Open downstream address [http://127.0.0.1:8080](http://127.0.0.1:8080) in browser.

Another examples [here](./examples).
