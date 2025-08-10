// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MetaNode is ERC20 {
    constructor() ERC20("MetaNode Token", "MNT") {
        // _mint(msg.sender, 1000000 * 10 ** decimals()); // Mint initial supply to the contract deployer
        _mint(msg.sender, 1000000 * 10 ** 18); // Mint initial supply to the contract deployer
    }
}