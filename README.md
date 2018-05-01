# ko

dapp skeleton based on go-kit

## test

```
curl -XPOST -d'{"param":{"orderId":"123"}}' http://localhost:4001/svc/order/v1/order
{"code":0,"msg":"ok","data":{"order":"123"}}

curl http://localhost:4001/svc/ucenter/v1/user/122552323
{"code":0,"msg":"ok","data":{"user":"122552323"}}
```

## todo

- [ ] gateway
    - [ ] rate limit
    - [ ] logging
        - [ ] access.log / error.log
        - [ ] ELK
    - [ ] instrumentation
    - [ ] auth
        - [ ] jwt
        - [ ] cookie
    - [ ] load balance
- [ ] svc 
    - [ ] connections
        - [ ] mysql
        - [ ] redis
    - [ ] logging
    - [ ] orm
    - [ ] mq
    
## 架构

![架构](https://ws3.sinaimg.cn/large/006tNc79gy1fqwe7f2kn6j31kw0v1dli.jpg)