// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Begging {

    // 捐款记录：address => amount
    mapping(address donater => uint256 amount) private _donations;
    // 捐赠者列表
    address[] private _donors;
    // 捐款总金额
    uint256 private _totaldonations;

    // 捐款事件
    event Donated(address indexed donot, uint256 value);

    // 捐赠函数：payable 允许转账
    function donate() external payable onlyDuringDonation {
        require (msg.value > 0,"Donation amount must be greater than zero.");

        // 如果是第一次捐赠，加入列表
        if(_donations[msg.sender] == 0){
            _donors.push(msg.sender);
        }

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

    // 以下是额外功能的函数。

    // 获取捐赠排行榜前三地址
    function getRankingTop3() external view returns (address[3] memory) {
        address[3] memory topAddress;
        uint256[3] memory topAmount;

        for (uint256 i = 0; i < _donors.length; i++){

            // 循环获取当前捐赠额
            address donor = _donors[i];
            uint256 amount = _donations[donor];

            // 只处理排行榜前三
            for (uint256 j = 0; j < 3; j++) {
                if(amount > topAmount[j]) {
                    // 插入当前项并移后其他项
                    for (uint256 k = 2; k > j; k--) {
                        topAddress[k] = topAddress[k - 1];
                        topAmount[k] = topAmount[k - 1];
                    }
                    topAddress[j] = donor;
                    topAmount[j] = amount;
                    break;
                }
            }
        }
        return topAddress;
    }

    uint256 private _startTime;
    uint256 private _duration;

    function setDonationPeriod(uint256 startTime, uint256 duration) external onlyOwner {
        _startTime = startTime;
        _duration = duration;
    }

    modifier onlyDuringDonation() {
        require(block.timestamp >= _startTime && block.timestamp < _startTime + _duration, "Donation period has ended.");
        _;
    }
    
    // 获取开始时间（辅助调试）
    function getStartTime() external view returns (uint256) {
        return _startTime;
    }

    // 获取结束时间（辅助调试）
    function getEndTime() external view returns (uint256) {
        return _startTime + _duration;
    }
}

// 已实现作业 基础功能+额外挑战 全部功能

// Sepolia 部署合约地址：
// 0x5EB226F16121325b4Da2810dD138Fef9744E2EeA