# Login example with redis

Launch redis cluster

```shell
docker run --rm --name=rediboard -p 6379:6379 redis
```

Then run `go run main.go`.

Then run curl like

```shell
curl -X POST localhost:8080/setgo
curl -X POST localhost:8080/set/go
curl localhost:8080/get
curl -X POST localhost:8080/set/java
curl localhost:8080/get
```
