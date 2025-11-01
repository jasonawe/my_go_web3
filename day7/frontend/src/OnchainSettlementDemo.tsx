import React, { useEffect, useState } from "react";
import { ethers } from "ethers";
import OnchainSettleABI from "./artifacts/contracts/OnchainSettle.sol/OnchainSettle.json"; // 从 Hardhat artifacts 导出或复制 ABI

const CONTRACT_ADDRESS = "0x394A2F8504e740a699E416c2a27307591eCfB979"; // <- 部署后替换

type TransferInstruction = {
  token: string;
  from: string;
  to: string;
  amount: string; // decimal string (wei)
};

export default function OnchainSettlementDemo() {
  const [provider, setProvider] = useState<ethers.BrowserProvider | null>(null);
  const [signer, setSigner] = useState<ethers.Signer | null>(null);
  const [account, setAccount] = useState<string>("");
  const [log, setLog] = useState<string[]>([]);
  const [tokenAddr, setTokenAddr] = useState<string>(""); // ERC20 token to deposit
  const [depositAmount, setDepositAmount] = useState<string>("0.0");
  const [contract, setContract] = useState<ethers.Contract | null>(null);
  const [contractOwner, setContractOwner] = useState<string>("");

  useEffect(() => {
    if ((window as any).ethereum) {
      const prov = new ethers.BrowserProvider((window as any).ethereum);
      setProvider(prov);
    } else {
      appendLog("Please install MetaMask.");
    }
  }, []);

  useEffect(() => {
    if (provider) {
      provider.getSigner().then(s => {
        setSigner(s);
      }).catch(()=>{});
      const c = new ethers.Contract(CONTRACT_ADDRESS, OnchainSettleABI.abi, provider);
      setContract(c);
      // fetch owner (if provider supports)
      (async () => {
        try {
          const owner = await c.owner();
          setContractOwner(owner);
        } catch (e) {}
      })();
    }
  }, [provider]);

  const appendLog = (s: string) => setLog(l => [new Date().toISOString() + " " + s, ...l].slice(0, 200));

  async function connect() {
    if (!provider) { appendLog("No provider"); return; }
    try {
      await (window as any).ethereum.request({ method: "eth_requestAccounts" });
      const signer = provider.getSigner();
      const addr = await signer.getAddress();
      setSigner(signer);
      setAccount(addr);
      appendLog("Connected: " + addr);
    } catch (e) { appendLog("connect failed: " + (e as any).message); }
  }

  async function approveAndDeposit() {
    if (!signer || !contract) { appendLog("not connected"); return; }
    if (!tokenAddr) { appendLog("token required"); return; }
    const erc20 = new ethers.Contract(tokenAddr, [
      "function approve(address spender, uint256 amount) public returns (bool)",
      "function decimals() view returns (uint8)"
    ], signer);
    const decimals = await erc20.decimals().catch(()=>18);
    const amt = ethers.parseUnits(depositAmount || "0", decimals);
    appendLog(`Approving ${depositAmount} tokens to settlement contract...`);
    const tx1 = await erc20.approve(CONTRACT_ADDRESS, amt);
    appendLog("approve txHash: " + tx1.hash);
    await tx1.wait();
    appendLog("approve confirmed. Now depositing...");
    // deposit via signer contract
    const cWithSigner = contract.connect(signer);
    const tx2 = await cWithSigner.deposit(tokenAddr, amt);
    appendLog("deposit txHash: " + tx2.hash);
    await tx2.wait();
    appendLog("deposit confirmed.");
  }

  async function refreshDepositedBalance() {
    if (!contract || !account || !tokenAddr) { appendLog("provide token and connected account"); return; }
    try {
      const bal: bigint = await contract.depositedBalance(tokenAddr, account);
      appendLog(`Deposited balance (wei): ${bal.toString()}`);
    } catch (e) {
      appendLog("error reading depositedBalance: " + (e as any).message);
    }
  }

  async function withdraw() {
    if (!contract || !signer || !tokenAddr) { appendLog("missing"); return; }
    const c = contract.connect(signer);
    const decimals = 18;
    const amt = ethers.parseUnits(depositAmount || "0", decimals);
    const tx = await c.withdraw(tokenAddr, amt);
    appendLog("withdraw txHash: " + tx.hash);
    await tx.wait();
    appendLog("withdraw confirmed");
  }

  // parse CSV-like multi-line input: token,from,to,amountDecimal
  function parseInstructions(text: string): TransferInstruction[] {
    const lines = text.split(/\r?\n/).map(l => l.trim()).filter(l => l.length);
    const arr: TransferInstruction[] = [];
    for (const ln of lines) {
      const parts = ln.split(",").map(p => p.trim());
      if (parts.length < 4) continue;
      arr.push({ token: parts[0], from: parts[1], to: parts[2], amount: parts[3] });
    }
    return arr;
  }

  async function settleFromTextarea(text: string) {
    if (!contract || !signer) { appendLog("connect and load contract"); return; }
    const inst = parseInstructions(text);
    if (inst.length === 0) { appendLog("no instructions parsed"); return; }
    // convert to arrays and amounts in wei (assume token decimals 18 for simplicity here)
    const tokens = inst.map(i => i.token);
    const froms = inst.map(i => i.from);
    const tos = inst.map(i => i.to);
    const amounts = inst.map(i => ethers.parseUnits(i.amount, 18)); // for prod, query decimals per token

    appendLog(`Submitting batch with ${inst.length} transfers as owner...`);
    const c = contract.connect(signer);
    try {
      const tx = await c.settleBatch(tokens, froms, tos, amounts);
      appendLog("settleBatch txHash: " + tx.hash);
      await tx.wait();
      appendLog("settleBatch confirmed");
    } catch (e: any) {
      appendLog("settleBatch failed: " + (e.message || e));
    }
  }

  return (
    <div style={{ padding: 20 }}>
      <h2>On-chain Settlement Demo</h2>
      <div>
        <button onClick={connect}>Connect MetaMask</button>
        <div>Connected: {account}</div>
        <div>Contract Owner: {contractOwner}</div>
      </div>

      <hr />
      <h3>Deposit (approve -| deposit)</h3>
      <div>
        <label>ERC20 token address: </label>
        <input value={tokenAddr} onChange={e => setTokenAddr(e.target.value)} style={{ width: 400 }} />
      </div>
      <div>
        <label>Amount (decimal): </label>
        <input value={depositAmount} onChange={e => setDepositAmount(e.target.value)} />
      </div>
      <div>
        <button onClick={approveAndDeposit}>Approve & Deposit</button>
        <button onClick={refreshDepositedBalance}>Show My Deposited Balance</button>
        <button onClick={withdraw}>Withdraw</button>
      </div>

      <hr />
      <h3>Operator: Batch settlement</h3>
      <p>Paste CSV lines: <code>token,from,to,amount</code> one per line (amount decimal)</p>
      <textarea id="batchText" defaultValue={`0xTOKEN,0xFrom,0xTo,1.0`} style={{ width: "100%", height: 160 }} />
      <div>
        <button onClick={() => {
          const t = (document.getElementById("batchText") as HTMLTextAreaElement).value;
          settleFromTextarea(t);
        }}>Submit Batch (settleBatch)</button>
      </div>

      <hr />
      <h3>Logs</h3>
      <div style={{ whiteSpace: "pre-wrap", maxHeight: 320, overflow: "auto", background: "#eee", padding: 8 }}>
        {log.map((l, i) => <div key={i}>{l}</div>)}
      </div>
    </div>
  );
}
