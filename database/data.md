# Data

# 1.Endpoint 설계

## Endpoint Fields

1. seq : number
2. name : 이름
3. host : host
4. port : port
5. need-load-balancing : 로드밸런싱필요유무
6. is-encrypted : 암호화유무

```js
db.endpoint.insertOne({
    "seq": 0,
    "name": "test",
    "host": "123.456.789.00",
    "port": "80",
    "allow-method": ["GET", "POST", "PUT", "DELETE"],
    "load-balancing-enabled": true,
    "resp-transformation-enabled": true,
    "allo-header" : [],
    "timeout" : 50, // setting defallt

    "req-enc" : "",
    "resp-enc" : "",

    "key-reqired" : true,
    "rate-limit":  3000, //throttling 
    // second
    // "caching-enabled" : true,
    "caching-expire-time" : 0, 
    "enabled": true,
    "log-enabled" : false,
    "resp-compression-enabled" : true, // 응답 압축
    
    "allow-domain" : [],//허용된 클라이언트의 IP 주소나 도메인
    
    "createdAt": "",
    "updatedAt": "",
});
```

# 2. Allow API Url

```js
db.url_white_list.insertOne({
    "seq": 0,
    "path": "test/list/hello/?",
    "endpoint": "test",
    "price": 200,
    "allow-method": ["GET", "POST", "PUT", "DELETE"],
    "enabled": true,
});
```