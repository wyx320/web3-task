// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

/*
✅  合并两个有序数组 (Merge Sorted Array)
题目描述：将两个有序数组合并为一个有序数组。
*/
contract MergeSortedArray {

    function merge(uint256[] memory arr1, uint256[] memory arr2)
        public
        pure
        returns (uint256[] memory)
    {
        uint256 len1 = arr1.length;
        uint256 len2 = arr2.length;

        uint256[] memory result = new uint256[](len1 + len2);

        uint256 i = 0; // 指向 arr1
        uint256 j = 0; // 指向 arr2
        uint256 k = 0; // 指向 result

        // 比较两个数组的元素，将较小的放入结果数组
        while (i < len1 && j < len2) {
            if (arr1[i] <= arr2[j]) {
                result[k] = arr1[i];
                i++;
            } else {
                result[k] = arr2[j];
                j++;
            }
            k++;
        }

        // 将 arr1 剩余元素全部返给结果数组
        while (i < arr1.length) {
            result[k] = arr1[i];
            i++;
            k++;
        }

        // 将 arr2 剩余元素全部返给结果数组
        while (j < arr2.length) {
            result[k] = arr2[j];
            j++;
            k++;
        }

        return result;
    }
}

/*
测试数据：
[2,2,3,7,8]
[1,2,3,4,5,9]
*/

/*
测试步骤演示：
2 2 3 7 8
1 2 3 4 5 9

i=0
j=0
k=0

1
i=0
j=1
k=1

1 2
i=1
j=1
k=2

1 2 2
i=2
j=1
k=3

1 2 2 2
i=2
j=2
k=4

1 2 2 2 3
i=3
j=2
k=5

2 2 3 7 8
1 2 3 4 5 9

1 2 2 2 3 3
i=3
j=3
k=6

1 2 2 2 3 3 4
i=3
j=4
k=6

1 2 2 2 3 3 4 5
i=3
j=5
k=7

1 2 2 2 3 3 4 5 7
i=4
j=4
k=8

1 2 2 2 3 3 4 5 7 8
i=5
j=4
k=8

2 2 3 7 8
1 2 3 4 5 9

1 2 2 2 3 3 4 5 7 8 9

*/