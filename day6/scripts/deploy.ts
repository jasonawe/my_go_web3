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
    const tx = await contract.mintNFT(wallet.address, "https://gateway.pinata.cloud/ipfs/QmRv9z1rJHNH81orPBFnFBQQWkfG5UrywQysPKshxvxGqk");
    await tx.wait();
    console.log("Minted NFT to:", wallet.address);
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});
