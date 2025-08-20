 // SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MathOptV3 {

    // 纯函数：没有存储写入，Gas更低
    function add(int256 a, int256 b) public pure returns (int256) {
        int256 newVal = a + b;
        return newVal;
    }

    function subtract(int256 a, int256 b) public pure returns (int256) {
        int256 newVal = a - b;
        return newVal;
    }
}
