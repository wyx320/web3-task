// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 透明代理
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
// UUPS代理
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
// 角色权限管理
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
// 数学工具
import "@openzeppelin/contracts/utils/math/Math.sol";

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
        // Staking token amount that user provided 用户质押的代币数量
        uint256 stAmount;
        // Finished distrubution MetaNode tokens to user 已领取的 MetaNode 数量
        uint256 finishedMetaNode;
        // Pending to claim MetaNode tokens 待领取的 MetaNode 数量
        uint256 pendingMetaNode;
        // Withdraw request list 解质押请求列表
        UnstakeRequest[] requests;
    }

    // 质押池结构体
    struct StakePool {
        // Address of stake token 质押代币的地址
        address stTokenAddress;
        // Weight of pool 质押池的权重，影响奖励分配
        uint256 poolWeight;
        // Last block number that MetaNode token distribution occurs for pool 最后一次计算奖励的区块号
        uint256 lastRewardBlock;
        // Accumulated MetaNode tokens for per staking token of pool 每个质押代币累积的 RCC 数量
        uint256 accMetaNodePerST;
        // Staking token amount 质押池中的总质押代币量
        uint256 stTokenAmount;
        // Min staking amount 最小质押金额
        uint256 minDepositAmount;
        // Withdrow locked blocks 解除质押的锁定区块数
        uint256 unstakeLockedBlocks;
    }

    // *************************************** STAKE VERIABLES ***************************************

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

    bool public withdrawPaused;
    bool public claimPaused;

    // pool id => user address => user info
    mapping(uint256 => mapping(address => User)) public users;

    // *************************************** Events ***************************************

    event PauseWithdraw();

    event UnpauseWithdraw();

    event PauseClaim();

    event UnpauseClaim();

    event SetStartBlock(uint256 startBlock);

    event SetEndBlock(uint256 endBlock);

    event SetMetaNodePerBlock(uint256 metaNodePerBlock);

    event AddPool(
        address stTokenAddress,
        uint256 poolWeight,
        uint256 lastRewardBlock,
        uint256 minDepositAmount,
        uint256 unstakeLockedBlocks
    );

    // *************************************** MODIFIER ***************************************

    modifier checkPid(uint256 _pid) {
        require(_pid < pool.length, "Invalid pool id");
        _;
    }

    // 合约初始化
    function initialize(
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _metaNodePerBlock
    ) public initializer {
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);

        startBlock = _startBlock;
        endBlock = _endBlock;
        metaNodePerBlock = _metaNodePerBlock;
    }

    // 合约升级自定义授权逻辑
    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyRole(UPGRADE_ROLE) {}

    // *************************************** ADMIN FUNCTION ***************************************

    /**
     * @notice Pause withdraw. Can only be called by admin
     */
    function pauseWithdraw() public onlyRole(ADMIN_ROLE) {
        require(!withdrawPaused, "Withdraw is already paused");

        withdrawPaused = true;

        emit PauseWithdraw();
    }

    /**
     * @notice Unpause withdraw. Can only be called by admin
     */
    function unpauseWithdraw() public onlyRole(ADMIN_ROLE) {
        require(withdrawPaused, "Withdraw is already unpaused");

        withdrawPaused = false;

        emit UnpauseWithdraw();
    }

    /**
     * @notice Pause claim. Can only be called by admin
     */
    function pauseClaim() public onlyRole(ADMIN_ROLE) {
        require(!claimPaused, "Claim is already paused");

        claimPaused = true;

        emit PauseClaim();
    }

    /**
     * @notice Unpause claim. Can only be called by admin
     */
    function unpauseClaim() public onlyRole(ADMIN_ROLE) {
        require(claimPaused, "Claim is already unpaused");

        claimPaused = false;

        emit UnpauseClaim();
    }

    /**
     * @notice Set start block. Can only be called by admin
     */
    function setStartBlock(uint256 _startBlock) public onlyRole(ADMIN_ROLE) {
        require(
            _startBlock <= endBlock,
            "Start block cannot be greater than end block"
        );

        startBlock = _startBlock;

        emit SetStartBlock(_startBlock);
    }

    /**
     * @notice Set end block. Can only be called by admin
     */
    function setEndBlock(uint256 _endBlock) public onlyRole(ADMIN_ROLE) {
        require(
            startBlock <= _endBlock,
            "End block cannot be less than start block"
        );

        endBlock = _endBlock;

        emit SetEndBlock(_endBlock);
    }

    /**
     * @notice Set MetaNode per block. Can only be called by admin
     */
    function setMetaNodePerBlock(
        uint256 _metaNodePerBlock
    ) public onlyRole(ADMIN_ROLE) {
        require(
            _metaNodePerBlock > 0,
            "MetaNode per block must be greater than 0"
        );

        metaNodePerBlock = _metaNodePerBlock;

        emit SetMetaNodePerBlock(_metaNodePerBlock);
    }

    /**
     * @notice Add a new staking to pool. Can only be called by admin
     * DO NOT add the same staking more than once. MetaNode rewards will be massed up if you do
     */
    function addPool(
        address _stTokenAddress,
        uint256 _poolWeight,
        uint256 _minDepositAmount,
        uint256 _unstakeLockedBlocks,
        bool _massUpdatePools
    ) public onlyRole(ADMIN_ROLE) {
        // Default the first pool to be ETH pool. so the first pool must be added with stTokenAddress = address(0x00)
        if (pool.length == 0) {
            require(_stTokenAddress == address(0), "First pool must be ETH");
        } else {
            require(_stTokenAddress != address(0), "Pool must be token");
        }

        // allow the min deposit amount equal to 0
        // require(_minDepositAmount > 0, "Min deposit amount must be greater than 0");
        require(
            _unstakeLockedBlocks > 0,
            "Unstake locked blocks must be greater than 0"
        );
        require(endBlock > block.timestamp, "End block must be in the future");

        if (_massUpdatePools) {}

        uint256 lastRewardBlock = block.number > startBlock
            ? block.number
            : startBlock;
        totalPoolWeight += _poolWeight;

        pool.push(
            StakePool({
                stTokenAddress: _stTokenAddress,
                poolWeight: _poolWeight,
                lastRewardBlock: lastRewardBlock,
                accMetaNodePerST: 0,
                stTokenAmount: 0,
                minDepositAmount: _minDepositAmount,
                unstakeLockedBlocks: _unstakeLockedBlocks
            })
        );

        emit AddPool(
            _stTokenAddress,
            _poolWeight,
            lastRewardBlock,
            _minDepositAmount,
            _unstakeLockedBlocks
        );
    }

    // *************************************** PUBLISH FUNCTION ***************************************

    /**
     * @notice Update reward variables for all pools. Be careful of gas spending!
     */
    function massUpdatePools() public {
        uint256 length = pool.length;
        for (uint256 pid = 0; pid < length; pid++) {
            updatePool(pid);
        }
    }

    /**
     * @notice Update reward variables of the given pool to be up-to-date. 更新给定池的奖励变量，使其保持最新。
     */
    function updatePool(uint256 _pid) public checkPid(_pid) {

    }
}
