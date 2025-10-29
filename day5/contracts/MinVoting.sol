// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MiniVoting {
    struct Proposal {
        string name;
        uint256 voteCount;
    }

    Proposal[] public proposals;
    mapping(address => bool) public voted;

    constructor(string[] memory proposalNames) {
        for (uint256 i = 0; i < proposalNames.length; i++) {
            proposals.push(Proposal({name: proposalNames[i], voteCount: 0}));
        }
    }

    function vote(uint256 proposalIndex) public {
        require(!voted[msg.sender], "Already voted");
        require(proposalIndex < proposals.length, "Invalid proposal index");
        proposals[proposalIndex].voteCount += 1;
        voted[msg.sender] = true;
    }

    function winningProposal() public view returns (uint256 winningIndex) {
        uint256 highest = 0;
        for (uint256 i = 0; i < proposals.length; i++) {
            if (proposals[i].voteCount > highest) {
                highest = proposals[i].voteCount;
                winningIndex = i;
            }
        }
    }

    function winnerName() external view returns (string memory) {
        return proposals[winningProposal()].name;
    }

    function getProposalsCount() external view returns (uint256) {
        return proposals.length;
    }
}
