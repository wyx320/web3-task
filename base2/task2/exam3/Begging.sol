// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Begging {

    // 捐款记录：address => amount
    mapping(address donater => uint256 amount) private _donations;
    // 捐款总金额
    uint256 private _totaldonations;

    // 捐款事件
    event Donated(address donot, uint256 value);

    // 捐赠函数：payable 允许转账
    function donate() external payable {
        require (msg.value > 0,"Donation amount must be greater than zero.");

        _donations[msg.sender] += msg.value;
        _totaldonations += msg.value;

        emit Donated(msg.sender, msg.value);
    }

    address private _owner;

    constructor(){
        _owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == _owner, "Only owner can call this function.");
        _;
    }

    // 合约所有者提取所有合约余额
    function withdraw() external onlyOwner {
        uint256 contractBalance = address(this).balance;
        require(contractBalance > 0,"No funds to withdraw.");

        payable(_owner).transfer(contractBalance);
    }

    // 查询某个地址的捐款金额
    function getDonation(address donor) external view returns (uint256) {
        return _donations[donor];
    }

    // 获取合约余额（辅助调试）
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}

// Sepolia 部署合约地址：0x1E9A30c75E9E5E354328B2C3f5f13591CaFC6d6a