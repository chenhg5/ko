# ko

dapp skeleton based on go-kit

## run

```
make
```

## test

```
curl -XPOST -d'{"param":{"orderId":"123"}}' http://localhost:4001/svc/order/v1/order
{"code":0,"msg":"ok","data":{"order":"123"}}

curl http://localhost:4001/svc/ucenter/v1/user/122552323
{"code":0,"msg":"ok","data":{"user":"122552323"}}
```

## architecture

### custom linus server architecture

![architecture](https://ws3.sinaimg.cn/large/006tNc79ly1fqwtlctza6j319i0p0afa.jpg)

### docker k8s architecture

...

## todo

- [ ] gateway
    - [ ] rate limit
    - [ ] logging
        - [ ] access.log / error.log
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
- [ ] code-generator-cmd-tool
- [ ] graceful restart/project hot update
- [ ] docker env