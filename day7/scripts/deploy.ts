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
    const artifactPath = path.resolve("artifacts/contracts/OnchainSettle.sol/OnchainSettle.json");
    const artifact = JSON.parse(fs.readFileSync(artifactPath, "utf8"));

    const factory = new ethers.ContractFactory(artifact.abi, artifact.bytecode, wallet);

    const onchainSettle = await factory.deploy();

    await onchainSettle.waitForDeployment();  // ethers v6

    console.log("OnchainSettle deployed at:", onchainSettle.target);
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});
