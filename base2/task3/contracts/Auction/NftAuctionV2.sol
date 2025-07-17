// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "./NftAuction.sol";

contract NftAuctionV2 is NftAuction {
    // V2 版本合约新增内容
    string public value;
    function setValue(string memory _newValue) public {
        value = _newValue;
    }
    function getValue() public view returns(string memory) {
        return value;
    }
}