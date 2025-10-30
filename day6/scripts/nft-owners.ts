import { ethers } from "ethers";
import * as fs from "fs";
import * as dotenv from "dotenv";
import fetch from "node-fetch"; // 需要安装 node-fetch

dotenv.config();

// 已部署合约地址
const CONTRACT_ADDRESS = "0x01f51446f51F0a44135fD471C7a0C228E698903B";

// 读取 ABI
const artifact = JSON.parse(fs.readFileSync("./artifacts/contracts/MyNFT.sol/MyNFT.json", "utf8"));
const abi = artifact.abi;

async function main() {
    const provider = new ethers.JsonRpcProvider(process.env.SEPOLIA_RPC_URL);
    const wallet = new ethers.Wallet(process.env.SEPOLIA_PRIVATE_KEY!, provider);

    const contract = new ethers.Contract(CONTRACT_ADDRESS, abi, wallet);

    // 获取当前 tokenCounter，表示总 NFT 数量
    const tokenCounter: bigint = await contract.tokenCounter();

    console.log(`Total NFTs minted: ${tokenCounter}`);

    for (let tokenId = 0n; tokenId < tokenCounter; tokenId++) {
        // 获取拥有者
        const owner = await contract.ownerOf(tokenId);
        // 获取 tokenURI
        const tokenURI: string = await contract.tokenURI(tokenId);

        // 获取 JSON 元数据
        let imageUrl = "";
        try {
            const response = await fetch(tokenURI);
            const metadata = await response.json();
            imageUrl = metadata.image || "";
        } catch (err) {
            console.error(`Failed to fetch metadata for tokenId ${tokenId}:`, err);
        }

        console.log(`NFT tokenId: ${tokenId}`);
        console.log(` Tokenurl:${tokenURI}`);
        console.log(`  Owner: ${owner}`);
        console.log(`  Image URL: ${imageUrl}`);
    }
}

main().catch((err) => {
    console.error(err);
    process.exit(1);
});
