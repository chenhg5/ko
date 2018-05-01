# ko

dapp skeleton based on go-kit

## test

```
curl -XPOST -d'{"param":{"orderId":"123"}}' http://localhost:4001/svc/order/v1/order
{"code":0,"msg":"ok","data":{"order":"123"}}

curl http://localhost:4001/svc/ucenter/v1/user/122552323
{"code":0,"msg":"ok","data":{"user":"122552323"}}
```