# Login example with redis

Launch redis cluster

```shell
docker run --rm --name=rediboard -p 6379:6379 redis
```

Then run `go run main.go`, then run `curl` commands like

```shell
# For the first example in first.go
curl -X POST localhost:8080/setgo
curl -X POST localhost:8080/set/go
curl localhost:8080/get
curl -X POST localhost:8080/set/java
curl localhost:8080/get
```

or

```shell
curl -H "Content-type: application/json" -d '{"username": "isa", "password": "test"}' localhost:8080/register
curl -H "Content-type: application/json" -d '{"username": "isa", "password": "test"}' localhost:8080/login
curl localhost:8080/welcome
```

reference: <https://blog.logrocket.com/how-to-use-redis-as-a-database-with-go-redis/>