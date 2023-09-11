start metrics
```shell
cd cadvisor 
docker-compose up -d
```
start app (build on arm64)
```shell
docker-compose up -d
```


wrk
app
```shell
wrk -t4 -c4 -d60 --latency -s post.lua http://localhost:8080/jwt/generate
```

app-withgomaxprocs
```shell
wrk -t4 -c4 -d60 --latency -s post.lua http://localhost:8082/jwt/generate
```
