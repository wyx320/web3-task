// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MathOptV2 {
    int256 public resultAdd = 3;
    int256 public resultSub = 2;

    // 原理：EVM 的 SSTORE 操作非常耗 Gas，如果写入的值没有变化，可以声调存储写，减少 Gas 消耗
    // 优化：只有当新结果和旧结果不同时才写入存储
    function add(int256 a, int256 b) public returns (int256) {
        int256 newVal = a + b;
        if (newVal != resultAdd) {
            resultAdd = newVal;
        }
        return newVal;
    }

    function subtract(int256 a, int256 b) public returns (int256) {
        int256 newVal = a - b;
        if (newVal != resultSub) {
            resultSub = newVal;
        }
        return newVal;
    }
}
