// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Counter {
    uint private count;

    event CountUpdated(uint newCount);

    function increment() public {
        count += 1;
        emit CountUpdated(count);
    }

    function decrement() public {
        require(count > 0, "Count cannot be negative");
        count -= 1;
        emit CountUpdated(count);
    }

    function reset() public {
        count = 0;
        emit CountUpdated(count);
    }

    function getCount() public view returns (uint) {
        return count;
    }
}
