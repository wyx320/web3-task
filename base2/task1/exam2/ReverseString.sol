// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

/*
✅ 反转字符串 (Reverse String)
题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"
*/

contract ReverseString {
    function reverse(string memory str) external pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        uint256 strBytesLength = strBytes.length;

        if (strBytesLength == 1) {
            return str;
        }

        for (uint256 i = 0; i < strBytesLength / 2; i++) {
            bytes1 temp = strBytes[i];
            strBytes[i] = strBytes[strBytesLength - i - 1];
            strBytes[strBytesLength - i - 1] = temp;
        }

        return string(strBytes);
    }
}
