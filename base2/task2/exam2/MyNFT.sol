// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";

contract MyNFT is ERC721URIStorage {
    constructor() ERC721("MyNFT", "MNFT") {}

    Counters.Counter private _tokenIds;
    using Counters for Counters.Counter;

    function mintNFT(address recipient, string memory tokenURI) public returns (uint256) {
        // Counters.increment(_tokenIds);
        // uint256 currentId = Counters.current(_tokenIds);
        _tokenIds.increment();
        uint256 currentId = _tokenIds.current();

        _mint(recipient, currentId);
        _setTokenURI(currentId, tokenURI);

        return currentId;
    }

    // 图片 CID
    // bafybeibggbooitt42yb7lsycoypz6ppe77ocqmzae4brdc52wpw32avnvy

    // 图片查看地址
    // https://moccasin-causal-gull-565.mypinata.cloud/ipfs/bafybeibggbooitt42yb7lsycoypz6ppe77ocqmzae4brdc52wpw32avnvy

    // tokenURI
    // ipfs://bafkreiehgpzylcotrnnchfn35a2h6shl3pmzaiawl3rk5sklmldjiwojlm

    // 合约地址：
    // Sepolia：0xc9b03B7B6E740Ef35791E7287e90F4E3DEC3BC7B
}
