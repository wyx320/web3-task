// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

/*
✅  用 solidity 实现罗马数字转数整数
题目描述在 https://leetcode.cn/problems/roman-to-integer/description/3.
*/

contract RomanToInteger {
    mapping(bytes1 => int256) private romanValues;

    constructor() {
        romanValues["I"] = 1;
        romanValues["V"] = 5;
        romanValues["X"] = 10;
        romanValues["L"] = 50;
        romanValues["C"] = 100;
        romanValues["D"] = 500;
        romanValues["M"] = 1000;
    }

    function romanToInteger(string memory roman)
        external
        view
        returns (int256)
    {
        bytes memory romanBytes = bytes(roman);
        uint256 length = romanBytes.length;

        // 检查长度限制
        require(1 <= length && length <= 15, "Invalid length");

        int256 result = 0;

        for (uint256 i = 0; i < length; i++) {
            int256 current = romanValues[romanBytes[i]];
            int256 next = i == romanBytes.length - 1
                ? int256(0)
                : romanValues[romanBytes[i + 1]];

            // 检查字符有效性
            require(current != 0, "Invalid Roman character");

            // 核心逻辑：当前字符值小于下一个字符值时需减去当前值
            if (current < next) {
                result -= current;
            } else {
                result += current;
            }
        }

        // 检查结果限制
        require(1 <= result && result <= 3999, "Out of valid range");

        return result;
    }
}
