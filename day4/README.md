# 🪙 ERC20 Token Project (Mini ERC20)

## 📘 项目简介 | Project Overview
本项目是一个基于 **Hardhat + TypeScript** 的 ERC20 代币项目，包含代币合约、部署脚本、转账脚本以及事件监听功能。  
The project is an **ERC20 token implementation** built with **Hardhat + TypeScript**, featuring token minting, transfer, and event tracking.

---

## 📁 项目结构 | Project Structure
```
├── contracts/              # Solidity 合约代码 / Smart Contracts
│   ├── ERC20Mintable.sol   # 可增发 ERC20 合约
│   ├── Counter.sol         # 示例计数器合约
│
├── scripts/                # 部署与交互脚本 / Deployment & Interaction Scripts
│   ├── deploy.ts           # 部署合约脚本
│   ├── transfer.ts         # 代币转账脚本
│   ├── check-balance.ts    # 查询余额脚本
│   ├── watch-transfer.ts   # 监听 Transfer 事件脚本
│
├── ignition/modules/       # Hardhat Ignition 模块化部署
│   ├── Counter.ts
│
├── test/                   # 测试文件 / Tests
│   ├── Counter.ts
│
├── .env                    # 环境变量配置（RPC、私钥等）
├── hardhat.config.ts       # Hardhat 配置文件
├── package.json
├── tsconfig.json
└── README.md
```

---

## ⚙️ 环境配置 | Environment Setup

### 1️⃣ 安装依赖 | Install Dependencies
```bash
npm install
```
依赖安装可能会不是那么一帆风顺，多尝试去解决


### 2️⃣ 创建 `.env` 文件 | Setup Environment Variables
```bash
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/{your_api_key}
SEPOLIA_PRIVATE_KEY={your_wallet_private_key}
SEPOLIA_WSS_URL=wss://sepolia.infura.io/ws/v3/{your_api_key}
```

---

## 🚀 部署合约 | Deploy Contract
```bash
npx hardhat compile
npx hardhat run scripts/deploy.ts --network sepolia
```
部署完成后会输出已部署的代币合约地址。

执行输出
```bash
Nothing to compile
Nothing to compile

[dotenv@17.2.3] injecting env (2) from .env -- tip: 👥 sync secrets across teammates & machines: https://dotenvx.com/ops
Deploying contracts with account: 0xaFa20294f278FEf6682000fAc922545542b7990c
Token deployed at: 0x71aeAE43e50dFCD8A322459beA3f907749e428D7
Minted 100 tokens to 0xaFa20294f278FEf6682000fAc922545542b7990c
Batch mint completed

```

---

## 🐟余额查询 | Token Balance
查询余额
```bash
npx tsx scripts/check-balance.ts
```
script/check-balance.ts代码，修改成自己的合约地址
```ts
  // 你的 ERC20 合约地址
  const tokenAddress = "0x71aeAE43e50dFCD8A322459beA3f907749e428D7";
```
运行输出
```bash
[dotenv@17.2.3] injecting env (2) from .env -- tip: ⚙️  override existing env vars with { override: true }
💰 Balance of 0xaFa20294f278FEf6682000fAc922545542b7990c: 1138.5 MINI
```

## 💸 转账操作 | Token Transfer
使用说明
```usage
npx tsx script/transfer.ts {from} {to} {amount}
```

使用实例
```bash
npx tsx scripts/transfer.ts 0xaFa20294f278FEf6682000fAc922545542b7990c 0x295BaE93a44B01920Ab4E2f243724958327CCa71  11.5
```

执行输出
```bash
[dotenv@17.2.3] injecting env (2) from .env -- tip: 🔄 add secrets lifecycle management: https://dotenvx.com/ops
🚀 Transfer script starting...
🔹 Params received:
  From: 0xaFa20294f278FEf6682000fAc922545542b7990c
  To: 0x295BaE93a44B01920Ab4E2f243724958327CCa71
  Amount: 11.5
🔑 Using wallet: 0xaFa20294f278FEf6682000fAc922545542b7990c
💰 Connected to token: MINI
📦 Sending 11.5 MINI to 0x295BaE93a44B01920Ab4E2f243724958327CCa71...
⏳ Tx sent: 0x85c31f1fc44a54d5a41b96e5d9de1e2b3b9f66bc252aaec8993198913e53232a
⏱ Waiting for confirmation...
✅ Tx confirmed in block 9501330
🏦 New balance of sender: 1138.5 MINI
```

---

## 🪙 铸币操作 | Mint Tokens
如果合约支持 `mint(address to, uint256 amount)`，则可直接调用：
```bash
npx hardhat console --network sepolia
> const token = await ethers.getContractAt("ERC20Mintable", "<TokenAddress>")
> await token.mint("<toAddress>", ethers.parseUnits("100", 18))
```

---

## 🎯 监听转账事件 | Listen to Transfer Events
运行监听脚本：
```bash
npx tsx scripts/watch-transfer.ts
```
控制台输出：
```
📤 Transfer detected:
  From:   0xaFa20294f278FEf6682000fAc922545542b7990c
  To:     0x295BaE93a44B01920Ab4E2f243724958327CCa71
  Amount: 0.1 MINI
```

该脚本使用 WebSocket 连接 (`wss://sepolia.infura.io/ws/v3/<your_api_key>`) 实时监听事件。

---

## 🧪 测试运行 | Run Tests
```bash
npx hardhat test
```

---

## 📚 技术栈 | Tech Stack
- **Solidity** – 智能合约语言  
- **Hardhat** – 合约编译、部署与测试框架  
- **TypeScript** – 脚本编写语言  
- **Ethers.js** – 与以太坊交互  
- **Infura / Sepolia** – 测试网络与节点服务  

