import { ethers } from "ethers";
import * as dotenv from "dotenv";
import ERC20MintableJson from "../artifacts/contracts/ERC20Mintable.sol/ERC20Mintable.json";

dotenv.config();

async function main() {
  // 从 .env 读取配置
  const rpcUrl = process.env.SEPOLIA_RPC_URL!;
  const privateKey = process.env.SEPOLIA_PRIVATE_KEY!;

  // 初始化 provider 和 signer
  const provider = new ethers.JsonRpcProvider(rpcUrl);
  const wallet = new ethers.Wallet(privateKey, provider);

  // 你的 ERC20 合约地址
  const tokenAddress = "0x71aeAE43e50dFCD8A322459beA3f907749e428D7";

  // 连接合约
  const token = new ethers.Contract(tokenAddress, ERC20MintableJson.abi, wallet);

  // 要查询的账户地址
  const target = "0xaFa20294f278FEf6682000fAc922545542b7990c";

  // 查询余额
  const balance = await token.balanceOf(target);

  console.log(`💰 Balance of ${target}: ${ethers.formatUnits(balance, 18)} MINI`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
