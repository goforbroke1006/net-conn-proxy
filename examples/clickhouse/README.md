# Proxy for Clickhouse

```shell
docker-compose down --volumes --remove-orphans
docker-compose up -d
```

```shell
clickhouse-client --host=localhost --port=9000 \
  --user=demo --password=demo \
  --database=demo \
  --query="SELECT version()"
```
