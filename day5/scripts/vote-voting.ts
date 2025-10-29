import { ethers } from "ethers";
import 'dotenv/config';

const { SEPOLIA_PRIVATE_KEY, SEPOLIA_RPC_URL } = process.env;

const contractAddress = "0x1039e272d697eA067f42096adab0A53E9cCb8b6A";
const abi = [
  "function vote(uint256 proposalIndex) external",
  "function winnerName() external view returns (string memory)",
  "function getProposalsCount() external view returns (uint256)"
];

async function main() {
    const provider = new ethers.JsonRpcProvider(SEPOLIA_RPC_URL);
    const wallet = new ethers.Wallet(SEPOLIA_PRIVATE_KEY, provider);
    const contract = new ethers.Contract(contractAddress, abi, wallet);

    // 投票给 proposal 1 (Bob)
    const tx = await contract.vote(2);
    await tx.wait();
    console.log("Voted for proposal 2");

    // 查询当前获胜者
    const winner = await contract.winnerName();
    console.log("Current winner:", winner);
}

main();
