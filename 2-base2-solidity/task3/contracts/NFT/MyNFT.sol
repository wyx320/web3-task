// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract MyNFT is ERC721, Ownable {
    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {}


    function mintNFT(address to, uint256 tokenId) public onlyOwner {
        _mint(to, tokenId);
    }
}
