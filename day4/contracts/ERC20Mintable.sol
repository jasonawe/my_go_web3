// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ERC20Mintable is ERC20, Ownable {

    uint256 public cap;

    event Minted(address indexed to, uint256 amount, address indexed operator);

    constructor(
        string memory name_,
        string memory symbol_,
        uint256 initialSupply_,
        uint256 cap_
    ) ERC20(name_, symbol_) Ownable(msg.sender) {   // ← 这里传入 msg.sender
        cap = cap_;
        if (initialSupply_ > 0) {
            _mint(msg.sender, initialSupply_);
            emit Minted(msg.sender, initialSupply_, msg.sender);
        }
    }

    function mint(address to, uint256 amount) external onlyOwner {
        if (cap != 0) {
            require(totalSupply() + amount <= cap, "Cap exceeded");
        }
        _mint(to, amount);
        emit Minted(to, amount, msg.sender);
    }

    function batchMint(address[] memory recipients, uint256[] memory amounts) external onlyOwner {
        require(recipients.length == amounts.length, "Length mismatch");
        for (uint i = 0; i < recipients.length; i++) {
            if (cap != 0) {
                require(totalSupply() + amounts[i] <= cap, "Cap exceeded");
            }
            _mint(recipients[i], amounts[i]);
            emit Minted(recipients[i], amounts[i], msg.sender);
        }
    }
}
