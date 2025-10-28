import { ethers } from "ethers";
import * as dotenv from "dotenv";
dotenv.config();

async function main() {
  console.log("🚀 Transfer script starting...");

  const args = process.argv.slice(2);
  if (args.length < 3) {
    console.error("❌ Usage: npx tsx scripts/transfer.ts <from> <to> <amount>");
    process.exit(1);
  }

  const [from, to, amountStr] = args;
  console.log(`🔹 Params received:
  From: ${from}
  To: ${to}
  Amount: ${amountStr}`);

  // 创建provider + wallet
  const provider = new ethers.JsonRpcProvider(process.env.SEPOLIA_RPC_URL);
  const wallet = new ethers.Wallet(process.env.SEPOLIA_PRIVATE_KEY!, provider);

  console.log(`🔑 Using wallet: ${wallet.address}`);

  // 你的部署合约地址
  const tokenAddress = "0x71aeAE43e50dFCD8A322459beA3f907749e428D7";

  // 只需要最小ABI
  const abi = [
    "function symbol() view returns (string)",
    "function balanceOf(address) view returns (uint256)",
    "function transfer(address to, uint256 amount) returns (bool)"
  ];

  const token = new ethers.Contract(tokenAddress, abi, wallet);
  const symbol = await token.symbol();
  const decimals = 18;

  console.log(`💰 Connected to token: ${symbol}`);

  const amount = ethers.parseUnits(amountStr, decimals);
  console.log(`📦 Sending ${amountStr} ${symbol} to ${to}...`);

  const tx = await token.transfer(to, amount);
  console.log(`⏳ Tx sent: ${tx.hash}`);
  console.log("⏱ Waiting for confirmation...");

  const receipt = await tx.wait();
  console.log(`✅ Tx confirmed in block ${receipt.blockNumber}`);

  const balance = await token.balanceOf(wallet.address);
  console.log(`🏦 New balance of sender: ${ethers.formatUnits(balance, decimals)} ${symbol}`);
}

main().catch((err) => {
  console.error("🔥 Error occurred:", err);
  process.exitCode = 1;
});
