 // SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MathOptV1 {
    int256 public result;

    function add(int256 a, int256 b) public returns (int256) {
        int256 newVal = a + b;
        result = newVal;
        return newVal;
    }

    function subtract(int256 a, int256 b) public returns (int256) {
        int256 newVal = a - b;
        result = newVal;
        return newVal;
    }
}
