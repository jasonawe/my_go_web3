import { ethers } from "ethers";
import * as dotenv from "dotenv";
dotenv.config();

async function main() {
  console.log("👀 Starting to listen for Transfer events...");
console.log("👀 Starting to listen for Transfer events...",process.env.SEPOLIA_WSS_URL);
  // ✅ 初始化 provider
  const provider = new ethers.WebSocketProvider(process.env.SEPOLIA_WSS_URL!);
  // 例如：wss://sepolia.infura.io/ws/v3/<YOUR_PROJECT_ID>

  // ✅ 代币合约地址
  const tokenAddress = "0x71aeAE43e50dFCD8A322459beA3f907749e428D7";

  // ✅ ERC20 ABI（仅保留需要的）
  const abi = [
    "event Transfer(address indexed from, address indexed to, uint256 value)",
    "function symbol() view returns (string)"
  ];

  const contract = new ethers.Contract(tokenAddress, abi, provider);
  const symbol = await contract.symbol();

  console.log(`📡 Listening on token ${symbol} at ${tokenAddress}`);
  console.log("-----------------------------------------------------");

  // ✅ 监听事件
  contract.on("Transfer", (from, to, value, event) => {
    console.log(`
📤 Transfer detected:
  From:   ${from}
  To:     ${to}
  Amount: ${ethers.formatUnits(value, 18)} ${symbol}
  TxHash: ${event.transactionHash}
-----------------------------------------------------`);
  });

  // 防止程序退出
  process.stdin.resume();
}

main().catch((err) => {
  console.error("❌ Error:", err);
});
