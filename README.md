# blockchain
一个简单区块链模仿品

运行 `demo` 后, 默认监听 `8080` 端口

`http://127.0.0.1:8080/mine` 执行一次挖矿操作

`http://127.0.0.1:8080/transactions/new` 通过 `post` 请求, 进行一次交易

格式为下

    {
        "sender":"d4ee26eee15148ee92c6cd394edd974e",
        "recipient": "someone-other-address",
        "amount": 5
    }

`http://127.0.0.1:8080/chain` 查看所有区块链
