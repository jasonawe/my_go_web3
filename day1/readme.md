## 这是eth相关的学习内容day1

### Part1-准备

基础知识准备，这里主要是要知道eth的网络，用于学习和测试的话，可以用sepolia.infura.io网络，sepolia.infura.io是eth的测试网络。今天主要的学习内容是通过rpc调用获取测试链上的相关信息。
测试网络RPC调用也是需要权限的，要去申请key才可以调用。

申请key地址`https://developer.metamask.io/key/all-endpoints`

然后是本地运行环境，需要go，这里是`go version go1.25.3 windows/amd64`，最好是安装最新版本

当然这里还需要科学上网工具

### Part2-获取当前块信息
先要配置.env
```
RPC_URL=https://sepolia.infura.io/v3/779e8a98dfb1417a9ad6002e4d1faa91
```

核心是通过RPC调用查询最新区块，就可以知道最新块号了，有助于咱们快速入门，增加成就感
当输入输出如下信息，即表示运行成功了，太棒了

```
go run get_block.go
```
```
✅ RPC 连接成功！当前区块号: 9484076
```
遇到问题也不要紧，一般都是包的问题，可以尝试使用提示来处理
### Part3-获取eth余额

先要配置.env
```
# 钱包地址
WALLET_ADDRESS=0xaFa20294f278FEf6682000fAc922545542b7990c
```
输入运行输出
```
go run eth_balance.go
```

```
钱包地址: 0xaFa20294f278FEf6682000fAc922545542b7990c
余额: 0.09798698 ETH
```

### Part3-获取代币合约的余额
对于什么是代币合约的余额，简单解释就是在以太坊里有中央银行发行的eth，也有商业银行发行的代币token，要查钱包的代币token余额就得去商业银行去查（合约地址）

非常好，第一二部分已经顺利完成了，现在整难一点的，虽然只是第一天，不过对于有基础的来说，其实原理都一样。

先要配置.env
```
# 钱包地址
WALLET_ADDRESS=0xaFa20294f278FEf6682000fAc922545542b7990c
# erc20合约地址-这个是USDC代币合约
ERC20_CONTRACT=0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
```

然后就是输入运行输出
```
go run erc20.go
```
```
✅ 钱包 0xaFa20294f278FEf6682000fAc922545542b7990c ERC20 余额: 0
```
太棒了，查到余额了，可惜是0，要是有个几百个就好了

### Part4-整合成http服务

既然咱们已经学了这么多了，何不整成一个http服务对外提供咱们的系统价值了，直接用gin整合起来

输入运行输出
```
go run .\gin_balance.go
```
```
- using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

2025/10/25 10:11:32 ✅ 已连接到 Sepolia 网络
[GIN-debug] GET    /balance/:address         --> main.main.func1 (3 handlers)
[GIN-debug] GET    /token_balance/:contract/:address --> main.main.func2 (3 handlers)
2025/10/25 10:11:32 🚀 服务已启动: http://localhost:8080
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
```

真不错，看来入门了

### 总结
1. 以太坊有多种网络，主网，测试等待
2. 要接入测试网络需要key，key申请要去developer
3. 链下与链上交互可以通过RPC调用
4. 钱包地址相当于个人唯一ID
5. 原生币只有eth，但是erc20代币有很多，要查代币数量需要erc20合约地址
