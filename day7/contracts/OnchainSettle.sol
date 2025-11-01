// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title OnchainSettle
/// @notice Users deposit ERC20 tokens to this contract. Operator (owner) can submit a batch of settlement transfers.
///         The contract will atomically execute all transfers from deposited balances to recipients.
contract OnchainSettle is Ownable {
    // token => user => balance (deposited into contract)
    mapping(address => mapping(address => uint256)) public deposits;

    event Deposited(address indexed token, address indexed user, uint256 amount);
    event Withdrawn(address indexed token, address indexed user, uint256 amount);
    event BatchSettled(address indexed operator, uint256 numTransfers);

    constructor() Ownable(msg.sender) {}

    /// @notice Deposit ERC20 tokens to the contract. User must approve this contract beforehand.
    function deposit(address token, uint256 amount) external {
        require(amount > 0, "amount=0");
        // transferFrom user -> this
        bool ok = IERC20(token).transferFrom(msg.sender, address(this), amount);
        require(ok, "transferFrom failed");
        deposits[token][msg.sender] += amount;
        emit Deposited(token, msg.sender, amount);
    }

    /// @notice Withdraw your deposited tokens (only available if you have balance)
    function withdraw(address token, uint256 amount) external {
        require(amount > 0, "amount=0");
        uint256 bal = deposits[token][msg.sender];
        require(bal >= amount, "insufficient deposited balance");
        deposits[token][msg.sender] = bal - amount;
        bool ok = IERC20(token).transfer(msg.sender, amount);
        require(ok, "transfer failed");
        emit Withdrawn(token, msg.sender, amount);
    }

    /// @notice Batch settlement executed by owner/operator.
    /// @param tokens array of ERC20 token addresses (len = N)
    /// @param froms array of payer addresses (len = N)
    /// @param tos array of recipient addresses (len = N)
    /// @param amounts array of amounts (len = N)
    ///
    /// Requirements:
    ///   - All arrays must be same length
    ///   - For each i, deposits[tokens[i]][froms[i]] >= amounts[i]
    /// The function will deduct deposits and transfer tokens from contract -> recipient for each instruction.
    function settleBatch(
        address[] calldata tokens,
        address[] calldata froms,
        address[] calldata tos,
        uint256[] calldata amounts
    ) external onlyOwner {
        uint256 n = tokens.length;
        require(n > 0, "empty batch");
        require(froms.length == n && tos.length == n && amounts.length == n, "length mismatch");

        // First pass: validate balances (to fail early if any insufficient)
        for (uint256 i = 0; i < n; ++i) {
            require(deposits[tokens[i]][froms[i]] >= amounts[i], "insufficient deposited balance for a transfer");
        }

        // Second pass: apply transfers
        for (uint256 i = 0; i < n; ++i) {
            address tok = tokens[i];
            address payer = froms[i];
            address recipient = tos[i];
            uint256 amt = amounts[i];

            // deduct internal balance
            deposits[tok][payer] -= amt;

            // transfer from contract to recipient
            bool ok = IERC20(tok).transfer(recipient, amt);
            require(ok, "erc20 transfer failed in batch");
        }

        emit BatchSettled(msg.sender, n);
    }

    // View helper: get deposited balance
    function depositedBalance(address token, address user) external view returns (uint256) {
        return deposits[token][user];
    }
}
