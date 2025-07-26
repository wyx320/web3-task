// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

/*
✅  用 solidity 实现整数转罗马数字
题目描述在 https://leetcode.cn/problems/integer-to-roman/description/
*/

contract IntegerToRoman {
    uint256[] private values = [
        1000,
        900,
        500,
        400,
        100,
        90,
        50,
        40,
        10,
        9,
        5,
        4,
        1
    ];
    string[] private symbol = [
        "M",
        "CM",
        "D",
        "CD",
        "C",
        "XC",
        "L",
        "XL",
        "X",
        "IX",
        "V",
        "IV",
        "I"
    ];

    function ToRoman(uint256 s) external view returns (string memory) {
        require(s >= 1 && s <= 3999, "Number out of range (1-3999)");
        bytes memory roman = "";

        for (uint256 i = 0; i < values.length; i++) {
            while (s >= values[i]) {
                roman =abi.encodePacked(roman, symbol[i]);
                s-=values[i];
            }
        }

        return string(roman);
    }
}
