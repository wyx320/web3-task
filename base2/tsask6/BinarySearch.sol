// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

/*
✅  二分查找 (Binary Search)
题目描述：在一个有序数组中查找目标值。
*/

contract BinarySearch {
    // 标准二分查找 - 返回目标值的索引，不如果不存在返回-1
    function binarySearch(uint256[] memory nums, uint256 target)
        external
        pure
        returns (int256)
    {
        // 处理空数组
        if (nums.length == 0) {
            return -1;
        }

        uint256 left = 0;
        uint256 right = nums.length - 1;

        while (left <= right) {
            uint256 mid = left + (right - left) / 2;

            if (nums[mid] == target) {
                return int256(mid);
            } else if (nums[mid] > target) {
                right = mid - 1;
            } else if (nums[mid] < target) {
                left = mid + 1;
            }
        }

        // 未找到
        return -1;
    }
}

/*
测试用例：
[2, 5, 8, 12, 16, 23, 38, 45, 56, 67, 78]
*/
