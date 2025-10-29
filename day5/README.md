# Mini Voting DApp / 小型投票去中心化应用

## 项目简介 / Project Overview
这是一个基于 Ethereum 的小型投票 DApp 项目，使用 Hardhat + TypeScript + Solidity 实现。项目支持智能合约部署、投票操作、投票列表查询，以及合约交互脚本。  
This is a small voting DApp project based on Ethereum, implemented with Hardhat + TypeScript + Solidity. The project supports smart contract deployment, voting operations, vote list query, and contract interaction scripts.

## 功能 / Features
- 部署投票智能合约 / Deploy voting smart contracts
- 列出投票候选 / List voting candidates
- 投票操作 / Cast votes
- 脚本化操作 / Scripted contract interactions
- 测试合约逻辑 / Test contract logic

## 技术栈 / Tech Stack
- Node.js + TypeScript
- Hardhat
- Solidity 0.8.x
- Ethers.js
- Mocha + Chai（合约测试）

## 环境安装 / Environment Setup

### 1. 克隆项目 / Clone the repository
```bash
git clone [<my_go_web3l>](https://github.com/jasonawe/my_go_web3.git)
cd day5
```

### 2. 初始化 / Initialize

```bash
npm init -y
npx hardhat --init
npm install dotenv
```

### 3. 配置.env / Config .env
```
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/779e8a98dfb1417a9ad6002e4d1faa90
SEPOLIA_PRIVATE_KEY={your_eth_wallet_private_key}
```

### 4. 编译部署合约 / compile contract
编译合约
```bash
npx hardhat compile
```

部署合约
``` bash 
npx hardhat run scripts/deploy.ts --network sepolia
```

输出示例
```bash
Nothing to compile
Nothing to compile

Deploying with account: 0xaFa20294f278FEf6682000fAc922545542b7990c
MiniVoting deployed at: 0x1039e272d697eA067f42096adab0A53E9cCb8b6A
```

## 使用示例 / Usage Example
### 列出投票候选 / List Voting Candidates

输入

```bash
npx hardhat run scripts/vote-list.ts --network sepolia
```
输出示例
```bash
共有 3 个候选人
候选人 0: Alice, 票数: 0
候选人 1: Bob, 票数: 1
候选人 2: Charlie, 票数: 0
当前获胜者: Bob

```
### 投票 / Cast Vote
输入
```bash
npx hardhat run scripts/vote-voting.ts --network sepolia
```
输出示例
```bash
Voted for proposal 1
Current winner: Bob
```



## 项目结构 / Project Structure
```
├── contracts/             # Solidity 智能合约 / Solidity smart contracts
│   └── MinVoting.sol      # 投票合约 / Voting contract
├── scripts/               # 合约部署与交互脚本 / Deployment & interaction scripts
├── test/                  # 测试脚本 / Test scripts
├── artifacts/             # Hardhat 编译产物 / Hardhat artifacts
├── cache/                 # 编译缓存 / Compilation cache
├── ignition/modules/      # 自定义模块 / Custom modules
├── package.json           # Node.js 项目配置 / Node.js project config
├── tsconfig.json          # TypeScript 配置 / TypeScript config
└── hardhat.config.ts      # Hardhat 配置 / Hardhat config
```