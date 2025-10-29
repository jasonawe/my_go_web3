import { ethers } from "ethers";
import fs from "fs";
import path from "path";
import 'dotenv/config';

const { SEPOLIA_PRIVATE_KEY, SEPOLIA_RPC_URL } = process.env;

async function main() {
    // provider 和钱包
    const provider = new ethers.JsonRpcProvider(SEPOLIA_RPC_URL);
    const wallet = new ethers.Wallet(SEPOLIA_PRIVATE_KEY, provider);

    console.log("Deploying with account:", wallet.address);

    // 读取编译好的合约 JSON（ABI + bytecode）
    const artifactPath = path.resolve("artifacts/contracts/MinVoting.sol/MiniVoting.json");
    const artifact = JSON.parse(fs.readFileSync(artifactPath, "utf8"));

    const factory = new ethers.ContractFactory(artifact.abi, artifact.bytecode, wallet);

    // 初始化提案
    const proposalNames = ["Alice", "Bob", "Charlie"];
    const voting = await factory.deploy(proposalNames);

    await voting.waitForDeployment();  // ethers v6

    console.log("MiniVoting deployed at:", voting.target);
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});
