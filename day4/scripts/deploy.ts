// scripts/deploy-mint.ts
import { ethers } from "ethers";
import * as dotenv from "dotenv";
import ERC20MintableJson from "../artifacts/contracts/ERC20Mintable.sol/ERC20Mintable.json";

dotenv.config();

async function main() {
  const rpcUrl = process.env.SEPOLIA_RPC_URL!;
  const privateKey = process.env.SEPOLIA_PRIVATE_KEY!;

  // provider + signer
  const provider = new ethers.JsonRpcProvider(rpcUrl);
  const wallet = new ethers.Wallet(privateKey, provider);

  console.log("Deploying contracts with account:", wallet.address);

  // 部署合约
  const factory = new ethers.ContractFactory(
    ERC20MintableJson.abi,
    ERC20MintableJson.bytecode,
    wallet
  );

  const token = await factory.deploy(
    "MiniToken",
    "MINI",
    ethers.parseUnits("1000", 18),
    ethers.parseUnits("1000000", 18)
  );

  await token.waitForDeployment();
  console.log("Token deployed at:", token.target);

  // 单个铸币
  const recipient = "0xaFa20294f278FEf6682000fAc922545542b7990c"; // 替换
  const tx1 = await token.mint(recipient, ethers.parseUnits("100", 18));
  await tx1.wait();
  console.log(`Minted 100 tokens to ${recipient}`);

  // 批量铸币
  const recipients = [
    "0xaFa20294f278FEf6682000fAc922545542b7990c",
    "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238"
  ];
  const amounts = [ethers.parseUnits("50", 18), ethers.parseUnits("30", 18)];

  const tx2 = await token.batchMint(recipients, amounts);
  await tx2.wait();
  console.log("Batch mint completed");
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});
