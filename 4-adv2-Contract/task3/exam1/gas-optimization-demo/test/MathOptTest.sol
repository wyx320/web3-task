// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "../src/MathOptV1.sol";
import "../src/MathOptV2.sol";
import "../src/MathOptV3.sol";

contract MathOptTest is Test {
    MathOptV1 v1;
    MathOptV2 v2;
    MathOptV3 v3;

    function setUp() public {
        v1 = new MathOptV1();
        v2 = new MathOptV2();
        v3 = new MathOptV3();
    }

    function test_v1_add() public {
        v1.add(1, 2);   // 会记录在 forge test --gas-report
        v1.add(1, 2);
    }
    function test_v1_subtract() public {
        v1.subtract(5, 3);
        v1.subtract(5, 3);
    }

    function test_v2_add() public {
        v2.add(1, 2);
        v2.add(1, 2);
    }
    function test_v2_subtract() public {
        v2.subtract(5, 3);
        v2.subtract(5, 3);
    }

    function test_v3_add() public view {
        v3.add(1, 2);
        v3.add(1, 2);
    }
    function test_v3_subtract() public view {
        v3.subtract(5, 3);
        v3.subtract(5, 3);
    }
}
