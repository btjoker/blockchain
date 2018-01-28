# blockchain
一个简单区块链模仿品

根据 [传送门](https://hackernoon.com/learn-blockchains-by-building-one-117428612f46) 粗略的翻译为 `go` 语言版本.

同时感谢 [掘金翻译计划](https://github.com/xitu/gold-miner) 提供的该 [文章](https://juejin.im/entry/59faa0ed51882576ea3507de) 翻译

`blockchain.py` 是原作者的 `python` 版本

构建 `demo` 运行后, 默认监听 `8080` 端口, 可用 demo -port :6666 指定端口

`http://127.0.0.1:8080/mine` 执行一次挖矿操作

`http://127.0.0.1:8080/transactions/new` 通过 `post` 请求, 进行一次交易

格式为下:

    {
        "sender":"d4ee26eee15148ee92c6cd394edd974e",
        "recipient": "someone-other-address",
        "amount": 5
    }

`http://127.0.0.1:8080/chain` 查看所有区块链


`http://127.0.0.1:8080/nodes/register` 通过 `post` 请求, 注册一个节点

格式为下:

    {
        "nodes": ["http://127.0.0.1:8081", "http://127.0.0.1:8082"]
    }


`http://127.0.0.1:8080/nodes/resolve` 解决冲突, 同步最长的正确区块链


以上就是照猫画虎的过程....To Be Continue