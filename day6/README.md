# NFT 项目使用说明 / NFT Project Guide

## 1️⃣ 项目概述 / Project Overview

本项目是一个基于 Ethereum/ERC-721 标准的 NFT 项目，支持：

- 部署 NFT 合约到测试网（Sepolia/Goerli）
- 铸造 NFT 给指定地址
- 查询 NFT 拥有者及图片信息
- 可批量铸造 NFT
- 支持 IPFS 存储 Metadata 和图片

This project is an Ethereum ERC-721 NFT project, supporting:

- Deploying NFT contract to testnets (Sepolia/Goerli)
- Minting NFTs to specified addresses
- Querying NFT owners and image URLs
- Batch minting
- Metadata and images stored on IPFS

---
## 2️⃣ IPFS准备 / IPFS Ready
IPFS是啥？为什么需要IPFS？ IPFS和NFT的关系是什么？带着这个三个疑问一起来准备

1. IPFS全称是星际文件系统，是一种去中心化分布式存储协议
2. IPFS是用来存储图片，音视频等文件的工具
3. 因为NFT合约是运行在区块链上，而多媒体等文件一般比较大，不适合存储在区块链上，区块链只存储文件的元信息地址，两者是相互协助的关系

操作步骤
1. 做NFT首先得准备一个开放的IFPS服务，这里可以使用https://pinata.cloud/
2. 注册完成后上传一个文件到IFPS
3. 编写你的metadata.json，文件{your_file_cid}在上传界面上有
```json
{
  "name": "Crypto Tiger #1",
  "description": "A brave digital tiger NFT.",
  "image": "ipfs://{your_file_cid}"
}
```
4. developer->API keys 生成一个access key（这里要注意权限，不懂得话就全部打开）
5. 得到bearer key后，组装curl，替换参数，执行
```curl
curl --location --request POST 'https://api.pinata.cloud/pinning/pinFileToIPFS' \
--header 'Authorization: Bearer {your_api_key}' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: api.pinata.cloud' \
--header 'Connection: keep-alive' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------' \
--form 'file=@"{your_path}/metadata.json"'
```
6. 执行结果出现IpfsHash属性既表明正常，其他情况需要自己排查
```json
{
    "IpfsHash": "{hash}",
    "PinSize": 174,
    "Timestamp": "2025-10-30T02:19:02.752Z",
    "ID": "66cf0f82-c844-40df-b771-2",
    "Name": "metadata.json",
    "NumberOfFiles": 1,
    "MimeType": "application/json",
    "GroupId": null,
    "Keyvalues": null
}
```
7. 最后进入https://pinata.cloud/查看文件情况，有两个文件即为正常
8. 在铸造NTF时使用类似https://gateway.pinata.cloud/ipfs/{metadata_file_cid}

----
## 3️⃣ 环境配置 / Environment Setup

1. 安装 Node.js 和 npm  
2. 安装依赖：

```bash
npm install ethers dotenv ts-node typescript @openzeppelin/contracts node-fetch
```

## 3️⃣ 创建 .env 文件
填写 RPC URL 和钱包私钥：
```env
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID
SEPOLIA_PRIVATE_KEY=0x你的钱包私钥
```
----
## 4️⃣ NFT合约编写&编译 / NFT Contract Code and Compile 
合约 MyNFT.sol 使用 ERC721URIStorage 和 Ownable：
```js
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721URIStorage, Ownable {
    uint256 public tokenCounter;

    // 给 Ownable 传入初始拥有者
    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {
        tokenCounter = 0;
    }

    function mintNFT(address recipient, string memory tokenURI) public onlyOwner returns (uint256) {
        uint256 newItemId = tokenCounter;
        _safeMint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);
        tokenCounter += 1;
        return newItemId;
    }
}
```

合约编译
```bash
# npx hardhat compile
Nothing to compile
Nothing to compile

```
----
5️⃣ NFT合约部署编写&部署 / Deploy NFT Contract Code And Deploy

scripts下编写delpoy.ts

```js
import { ethers } from "ethers";
import * as fs from "fs";
import * as dotenv from "dotenv";

dotenv.config();

// 读取 ABI 和 bytecode
const artifact = JSON.parse(fs.readFileSync("./artifacts/contracts/MyNFT.sol/MyNFT.json", "utf8"));
const abi = artifact.abi;
const bytecode = artifact.bytecode;

async function main() {
    const provider = new ethers.JsonRpcProvider(process.env.SEPOLIA_RPC_URL);
    const wallet = new ethers.Wallet(process.env.SEPOLIA_PRIVATE_KEY!, provider);

    // 创建合约工厂
    const factory = new ethers.ContractFactory(abi, bytecode, wallet);

    console.log("Deploying contract...");
    const contract = await factory.deploy();
    await contract.waitForDeployment(); // 等待部署完成

    console.log("Contract deployed at:", contract.target);

    // 铸造 NFT
    const tx = await contract.mintNFT(wallet.address, "{your_ifps_file_metadata_url}");
    await tx.wait();
    console.log("Minted NFT to:", wallet.address);
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});

```
NFT合约部署
```bash
# npx hardhat run scripts/deploy.ts --network sepolia
Nothing to compile
Nothing to compile

[dotenv@17.2.3] injecting env (2) from .env -- tip: ⚙️  write to custom object with { processEnv: myObject }
Deploying contract...
Contract deployed at: 0x01f51446f51F0a44135fD471C7a0C228E698903B
Minted NFT to: 0xaFa20294f278FEf6682000fAc922545542b7990c

```
----
## 6️⃣ 查询 NFT / Query NFT

在scripts下编写一个nft-owners.ts的查询脚本

```bash
# npx tsx scripts\nft-owners.ts
[dotenv@17.2.3] injecting env (2) from .env -- tip: 🛠️  run anywhere with `dotenvx run -- yourcommand`
Total NFTs minted: 1
NFT tokenId: 0
 Tokenurl:https://gateway.pinata.cloud/ipfs/QmRv9z1rJHNH81orPBFnFBQQWkfG5UrywQysPKshxvxGqk
  Owner: 0xaFa20294f278FEf6682000fAc922545542b7990c
  Image URL: ipfs://bafkreifrfeoapqykusbiltgalbzlur65xfhjvfcljux6l6grk65skaaoee

```
可以看到结果正常
----

## 7️⃣ 注意事项 / Notes

链上只存储 tokenId、拥有者和 tokenURI，实际资源（图片/音频）存储在 IPFS 等链下存储。

部署和铸造需要测试网 ETH。

批量铸造时注意 Gas 消耗。

tokenURI 可以使用 Pinata IPFS 网关 URL，例如：

https://gateway.pinata.cloud/ipfs/QmRv9z1rJHNH81orPBFnFBQQWkfG5UrywQysPKshxvxGqk

----

## 8️⃣ 参考 / References

[OpenZeppelin ERC721 Documentation](https://docs.openzeppelin.com/contracts/4.x/erc721)

[Ethers.js Documentation](https://docs.ethers.org/v6/)

[Pinata IPFS](https://pinata.cloud/)


