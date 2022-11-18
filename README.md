# HTTP DUMP Server

## Listening port

Set listening port with environment variable HTTP_SERVER_PORT (default 8080)
```
HTTP_SERVER_PORT=4242 ./http-dump
```

## Info
/info print json message
```
curl -s http://localhost:8080/info
```

## Dump
/dump endpoint just dump incoming request in the response
```
curl -s http://localhost:8080/dump
```
