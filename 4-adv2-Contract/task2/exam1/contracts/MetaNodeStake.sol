// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 透明代理
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
// UUPS代理
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
// 角色权限管理
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";

contract MetaNodeStake is
    Initializable,
    UUPSUpgradeable,
    AccessControlUpgradeable
{
    // 定义权限
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant UPGRADE_ROLE = keccak256("UPGRADE_ROLE");

    // 解质押请求结构体
    struct UnstakeRequest {
        // 解质押数量
        uint256 amount;
        // 解锁区块
        uint256 unlockBlock;
    }

    // 用户结构体
    struct User {
        // 用户质押的代币数量
        uint256 stAmount;
        // 已领取的 MetaNode 数量
        uint256 finishedMetaNode;
        // 待领取的 MetaNode 数量
        uint256 pendingMetaNode;
        // 解质押请求列表
        UnstakeRequest[] requests;
    }

    // 质押池结构体
    struct StakePool {
        // 质押代币的地址
        address stTokenAddress;
        // 质押池的权重，影响奖励分配
        uint256 poolWeight;
        // 最后一次计算奖励的区块好
        uint256 lastRewardBlock;
        // 每个质押代币累积的 RCC 数量
        // RCC = Reward Calculation Coin （奖励计算代币）
        // acc = accumulate (累积)
        uint256 accMetaNodePerST;
        // 质押池中的总质押代币量
        uint256 stTokenAmount;
        // 最小质押金额
        uint256 minDepositAmount;
        // 解除质押的锁定区块数
        uint256 unstakeLockedBlocks;
    }

    // *************************************** 质押变量 ***************************************

    // 质押池
    StakePool[] public pool;
    // 质押池总权重
    uint256 totalPoolWeight;

// MetaNode 质押合约开始区块
    uint256 public startBlock;
    // MetaNode 质押合约结束区块
    uint256 public endBlock;
    // MetaNode 每个区块的奖励
    uint256 public metaNodePerBlock;

    uint256 public withdrawPaused;
    uint256 public claimPaused;

    mapping(uint256 => mapping(address => User)) public users;

    function initialize() public initializer {
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
    }

    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyRole(UPGRADE_ROLE) {}
}
