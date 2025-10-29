import { ethers } from "ethers";
import 'dotenv/config';

const { SEPOLIA_PRIVATE_KEY, SEPOLIA_RPC_URL } = process.env;

const contractAddress = "0x1039e272d697eA067f42096adab0A53E9cCb8b6A";

// 只需要 ABI 中涉及读取的部分
const abi = [
  "function proposals(uint256 index) external view returns (string name, uint256 voteCount)",
  "function getProposalsCount() external view returns (uint256)",
  "function winnerName() external view returns (string)"
];

async function main() {
    const provider = new ethers.JsonRpcProvider(SEPOLIA_RPC_URL);
    const wallet = new ethers.Wallet(SEPOLIA_PRIVATE_KEY, provider);

    const contract = new ethers.Contract(contractAddress, abi, wallet);

    // 获取候选人数
    const count = await contract.getProposalsCount();
    console.log(`共有 ${count} 个候选人`);

    // 遍历显示票数
    for (let i = 0; i < count; i++) {
        const p = await contract.proposals(i);
        console.log(`候选人 ${i}: ${p.name}, 票数: ${p.voteCount.toString()}`);
    }

    // 查询当前获胜者
    const winner = await contract.winnerName();
    console.log(`当前获胜者: ${winner}`);
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});
