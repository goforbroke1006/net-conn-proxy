# L4D2 server proxying

On machine for proxy and game server:

```shell
docker-compose up -d

ifconfig
# my IP in wlan is 192.168.0.9
```

On machine with Steam and client of L4D2:

1. Run L4D2 game
2. Press ~ button to open developer console (first time you should enable it in the settings screen)
3. Type **connect 192.168.0.9:27015**