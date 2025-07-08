// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

/*
✅ 创建一个名为Voting的合约，包含以下功能：
一个mapping来存储候选人的得票数
一个vote函数，允许用户投票给某个候选人
一个getVotes函数，返回某个候选人的得票数
一个resetVotes函数，重置所有候选人的得票数
*/

import "@openzeppelin/contracts/access/Ownable.sol";

contract Voting is Ownable {
    // 存储每个候选人的得票数
    mapping(string => uint256) public candidateVotes;

    // 记录每个地址是否已经投过票
    mapping(address => bool) public hasVoted;

    // 记录候选人列表
    string[] public candidateList;

    constructor() Ownable(msg.sender) {}

    function vote(string memory candidate) external {
        require(!hasVoted[msg.sender], "You have already voted.");

        candidateVotes[candidate] += 1;
        candidateList.push(candidate);

        hasVoted[msg.sender] = true;
    }

    function getVotes(string memory candidate) external view returns (uint256) {
        return candidateVotes[candidate];
    }

    function resetVotes() external onlyOwner {
        for (uint256 i = 0; i < candidateList.length; i++) {
            delete candidateVotes[candidateList[i]];
        }
        delete candidateList;
    }
}
