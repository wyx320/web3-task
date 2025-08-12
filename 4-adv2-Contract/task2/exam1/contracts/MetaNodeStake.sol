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
// ERC20 代币
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
// 暂停/恢复合约 UUPS 更新
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";

import "hardhat/console.sol";

contract MetaNodeStake is
    Initializable,
    UUPSUpgradeable,
    AccessControlUpgradeable,
    PausableUpgradeable
{
    using Math for uint256;
    using SafeERC20 for IERC20;

    // **************************************** INVARIANT ***************************************

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant UPGRADE_ROLE = keccak256("UPGRADE_ROLE");

    uint256 public constant ETH_PID = 0;

    // *************************************** DATA STRUCTS ***************************************

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
        // Finished distrubution MetaNode tokens to user 已完成向用户分发奖励 MetaNode 数量 （用户应得但尚未领取的奖励的起点，可理解为已结算奖励的基准值）
        // 被更新场景： 1.用户质押/提取时 2.池子奖励更新时 （UpdatePool）
        uint256 finishedMetaNode;
        // Pending to claim MetaNode tokens 用户尚未领取的奖励 MetaNode 数量
        // 被更新场景：1.计算新增的奖励时 2.用户领取奖励时（claim）
        uint256 pendToClaimRewards;
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
        uint256 accMetaNodePerSt;
        // Staking token amount 质押池中的总质押代币量
        uint256 stTokenAmount;
        // Min staking amount 最小质押金额
        uint256 minDepositAmount;
        // Withdrow locked blocks 解除质押的锁定区块数
        uint256 unstakeLockedBlocks;
    }

    // *************************************** STAKE VERIABLES ***************************************

    // 质押池
    StakePool[] public pools;
    // 质押池总权重
    uint256 public totalPoolWeight;

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

    // MetaNode token
    IERC20 public MetaNode;

    // *************************************** Events ***************************************

    event PauseWithdraw();

    event UnpauseWithdraw();

    event PauseClaim();

    event UnpauseClaim();

    event SetStartBlock(uint256 indexed startBlock);

    event SetEndBlock(uint256 indexed endBlock);

    event SetMetaNodePerBlock(uint256 indexed metaNodePerBlock);

    event AddPool(
        address indexed stTokenAddress,
        uint256 indexed poolWeight,
        uint256 indexed lastRewardBlock,
        uint256 minDepositAmount,
        uint256 unstakeLockedBlocks
    );

    event UpdatePool(
        uint256 indexed pid,
        uint256 indexed lastRewardBlock,
        uint256 newTotalRewards
    );

    event SetMetaNode(IERC20 indexed metaNode);

    event UpdatePoolInfo(
        uint256 indexed pid,
        uint256 indexed minDepositAmount,
        uint256 indexed unstakeLockedBlocks
    );

    event UpdatePoolWeight(
        uint256 indexed pid,
        uint256 indexed poolWeight,
        uint256 totalPoolWeight
    );

    event Deposit(address indexed user, uint256 indexed pid, uint256 amount);

    event Withdraw(
        address indexed user,
        uint256 indexed pid,
        uint256 amount,
        uint256 blockNumber
    );

    event ClaimRewards(
        address indexed user,
        uint256 indexed pid,
        uint256 amount
    );

    event RequestStake(
        address indexed user,
        uint256 indexed pid,
        uint256 amount
    );

    // *************************************** MODIFIER ***************************************

    modifier checkPid(uint256 _pid) {
        require(_pid < pools.length, "Invalid pool id");
        _;
    }

    modifier whenNotWithdrawPaused() {
        require(!withdrawPaused, "Withdraw is paused");
        _;
    }

    modifier whenNotPausedClaimRewards() {
        require(!claimPaused, "Claim is paused");
        _;
    }

    /*
     * @notice Initialize MetaNodeStake 合约初始化
     */
    function initialize(
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _metaNodePerBlock,
        IERC20 _metaNode
    ) public initializer {
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);

        setMetaNode(_metaNode);

        startBlock = _startBlock;
        endBlock = _endBlock;
        metaNodePerBlock = _metaNodePerBlock;
    }

    /*
     * @notice Authorize upgrade 合约升级自定义授权逻辑
     */
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
        bool _massUpdatePoolsRewards
    ) public onlyRole(ADMIN_ROLE) {
        // Default the first pool to be ETH pool. so the first pool must be added with stTokenAddress = address(0x00)
        if (pools.length == 0) {
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
        require(endBlock > block.number, "End block must be in the future");

        if (_massUpdatePoolsRewards) {
            massUpdatePoolsRewards();
        }

        uint256 lastRewardBlock = block.number > startBlock
            ? block.number
            : startBlock;
        totalPoolWeight += _poolWeight;

        pools.push(
            StakePool({
                stTokenAddress: _stTokenAddress,
                poolWeight: _poolWeight,
                lastRewardBlock: lastRewardBlock,
                accMetaNodePerSt: 0,
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

    /**
     * @notice Set MetaNode token address. Can only be called by admin
     */
    function setMetaNode(IERC20 _metaNode) public onlyRole(ADMIN_ROLE) {
        MetaNode = _metaNode;

        emit SetMetaNode(_metaNode);
    }

    /**
     * @notice Update pool info. Can only be called by admin
     */
    function updatePoolInfo(
        uint256 _pid,
        uint256 _minDepositAmount,
        uint256 _unstakeLockedBlocks
    ) public onlyRole(ADMIN_ROLE) checkPid(_pid) {
        pools[_pid].minDepositAmount = _minDepositAmount;
        pools[_pid].unstakeLockedBlocks = _unstakeLockedBlocks;

        emit UpdatePoolInfo(_pid, _minDepositAmount, _unstakeLockedBlocks);
    }

    /**
     * @notice Update pool weight. Can only be called by admin
     */
    function updatePoolWeight(
        uint256 _pid,
        uint256 _poolWeight,
        bool _massUpdatePoolsRewards
    ) public onlyRole(ADMIN_ROLE) checkPid(_pid) {
        require(_poolWeight > 0, "invalid pool weight");

        if (_massUpdatePoolsRewards) {
            massUpdatePoolsRewards();
        }

        totalPoolWeight =
            totalPoolWeight -
            pools[_pid].poolWeight +
            _poolWeight;
        pools[_pid].poolWeight = _poolWeight;

        emit UpdatePoolWeight(_pid, _poolWeight, totalPoolWeight);
    }

    // *************************************** PUBLIC FUNCTION ***************************************

    /**
     * @notice Update reward variables for all pools. Be careful of gas spending!
     */
    function massUpdatePoolsRewards() public {
        uint256 length = pools.length;
        for (uint256 pid = 0; pid < length; pid++) {
            updatePoolRewards(pid);
        }
    }

    /**
     * @notice Update reward variables of the given pool to be up-to-date. 更新给定池的奖励变量，使其保持最新。
     */
    function updatePoolRewards(uint256 _pid) public checkPid(_pid) {
        StakePool storage pool_ = pools[_pid];

        if (block.number <= pool_.lastRewardBlock) {
            return;
        }

        (bool success1, uint256 totalMetaNode) = getRewardsByBolcks(
            pool_.lastRewardBlock,
            block.number
        ).tryMul(pool_.poolWeight);
        require(success1, "overflow");

        (success1, totalMetaNode) = totalMetaNode.tryDiv(totalPoolWeight);
        require(success1, "overflow");

        uint256 stSupply = pool_.stTokenAmount;
        if (stSupply > 0) {
            (bool success2, uint256 totalMetaNode_) = totalMetaNode.tryMul(
                1 ether
            );
            require(success2, "overflow");

            (success2, totalMetaNode_) = totalMetaNode_.tryDiv(stSupply);
            require(success2, "overflow");

            (bool success3, uint256 accMetaNodePerST) = pool_
                .accMetaNodePerSt
                .tryAdd(totalMetaNode_);
            require(success3, "overflow");
            pool_.accMetaNodePerSt = accMetaNodePerST;
        }

        pool_.lastRewardBlock = block.number;

        emit UpdatePool(_pid, pool_.lastRewardBlock, totalMetaNode);

        // StakePool storage pool = pools[_pid];

        // if (block.number <= pool.lastRewardBlock) {
        //     return;
        // }

        // (bool success1, uint256 newTotalRewards) = getRewardsByBolcks(
        //     pool.lastRewardBlock,
        //     block.number
        // ).tryMul(pool.poolWeight);
        // require(success1, "overflow1");

        // (bool success2, uint256 poolRewards) = newTotalRewards.tryDiv(
        //     totalPoolWeight
        // );
        // require(success2, "overflow2");

        // if (pool.stTokenAmount > 0) {
        //     console.log(
        //         "updatePoolRewards:pool.stTokenAmount:",
        //         pool.stTokenAmount
        //     );
        //     (bool success3, uint256 metaNodePerSt) = poolRewards.tryMul(
        //         1 ether
        //     );
        //     require(success3, "overflow3");

        //     // (success3, metaNodePerBlock) = metaNodePerSt.tryDiv(
        //     //     pool.stTokenAmount
        //     // );
        //     // require(success3, "overflow4");

        //     (success3, metaNodePerSt) = metaNodePerSt.tryDiv(
        //         pool.stTokenAmount
        //     );
        //     require(success3, "overflow4");

        //     (bool success4, uint256 accMetaNodePerSt) = metaNodePerSt.tryAdd(
        //         pool.accMetaNodePerSt
        //     );
        //     require(success4, "overflow5");

        //     console.log(
        //         "updatePoolRewards:pool.accMetaNodePerSt:",
        //         accMetaNodePerSt
        //     );
        //     pool.accMetaNodePerSt = accMetaNodePerSt;
        // }

        // pool.lastRewardBlock = block.number;

        // emit UpdatePool(_pid, pool.lastRewardBlock, newTotalRewards);
    }

    /**
     * @notice Deposit ETH to MetaNodeStake for MetaNode allocation.
     */
    function depositETH() public payable whenNotPaused {
        require(msg.value > 0, "Deposit amount must be greater than 0");
        StakePool storage ethPool = pools[ETH_PID];
        require(
            ethPool.stTokenAddress == address(0),
            "ETH staking token address"
        );

        _deposit(ETH_PID, msg.value);
    }

    /**
     * @notice Deposit tokens to MetaNodeStake for MetaNode allocation.
     */
    function deposit(
        uint256 _pid,
        uint256 _amount
    ) public checkPid(_pid) whenNotPaused {
        require(_pid != 0, "Deposit to ETH pool must use depositETH");
        require(_amount > 0, "Deposit amount must be greater than 0");

        StakePool storage pool = pools[_pid];
        require(
            pool.minDepositAmount <= _amount,
            "Deposit amount is less than min deposit amount"
        );

        IERC20(pool.stTokenAddress).safeTransferFrom(
            msg.sender,
            address(this),
            _amount
        );

        _deposit(_pid, _amount);
    }

    /**
     * @notice Unstake tokens from a pool.
     */
    function unstake(
        uint256 _pid,
        uint256 _amount
    ) public checkPid(_pid) whenNotPaused whenNotWithdrawPaused {
        StakePool storage pool_ = pools[_pid];
        User storage user_ = users[_pid][msg.sender];

        require(user_.stAmount >= _amount, "Not enough staking token balance");

        updatePoolRewards(_pid);

        uint256 pendingMetaNode_ = (user_.stAmount * pool_.accMetaNodePerSt) /
            (1 ether) -
            user_.pendToClaimRewards;

        if (pendingMetaNode_ > 0) {
            user_.pendToClaimRewards =
                user_.pendToClaimRewards +
                pendingMetaNode_;
        }

        if (_amount > 0) {
            user_.stAmount = user_.stAmount - _amount;
            user_.requests.push(
                UnstakeRequest({
                    amount: _amount,
                    unlockBlock: block.number + pool_.unstakeLockedBlocks
                })
            );
        }

        pool_.stTokenAmount = pool_.stTokenAmount - _amount;
        user_.finishedMetaNode =
            (user_.stAmount * pool_.accMetaNodePerSt) /
            (1 ether);

        emit RequestStake(msg.sender, _pid, _amount);

        // require(_amount > 0, "Unstake amount must be greater than 0");
        // User storage user = users[_pid][msg.sender];
        // require(user.stAmount >= _amount, "Insufficient staked amount");

        // // 更新质押池状态 主要是更新池的 accMetaNodePerSt
        // updatePoolRewards(_pid);

        // StakePool storage pool = pools[_pid];

        // // 先减少用户质押金额和池总质押金额，防止重入攻击
        // uint256 prevStAmount = user.stAmount;
        // user.stAmount -= prevStAmount - _amount;
        // pool.stTokenAmount -= pool.stTokenAmount - _amount;

        // // 计算用户未领取的奖励 （pending rewards）
        // uint256 pendToClaimRewards = (prevStAmount * pool.accMetaNodePerSt) /
        //     1 ether -
        //     user.finishedMetaNode;
        // if (pendToClaimRewards > 0) {
        //     user.pendToClaimRewards += pendToClaimRewards;
        // }

        // // 添加解质押请求
        // user.requests.push(
        //     UnstakeRequest({
        //         amount: _amount,
        //         unlockBlock: block.number + pool.unstakeLockedBlocks
        //     })
        // );

        // emit RequestStake(msg.sender, _pid, _amount);
    }

    /**
     * @notice Withdraw tokens from MetaNodeStake.
     */
    function withdraw(
        uint256 _pid
    ) public checkPid(_pid) whenNotPaused whenNotWithdrawPaused {
        StakePool storage pool = pools[_pid];
        User storage user = users[_pid][msg.sender];

        uint256 pendingWithdrawAmount;
        uint256 popNum;
        if (user.requests.length > 0) {
            for (uint256 i = 0; i < user.requests.length; i++) {
                console.log("user.requests[i].unlockBlock:", user.requests[i].unlockBlock);
                console.log("block.number:", block.number);
                if (user.requests[i].unlockBlock <= block.number) {
                    // 解锁区块已到，允许提取
                    pendingWithdrawAmount += user.requests[i].amount;
                    popNum++;
                }
            }
            for (uint256 j = 0; j < user.requests.length - popNum; j++) {
                // 前移未处理的请求
                user.requests[j] = user.requests[j + popNum];
            }
            for (uint256 k = 0; k < popNum; k++) {
                // 删除尾部冗余请求
                user.requests.pop();
            }

            if (pool.stTokenAddress == address(0)) {
                // 如果是 ETH 池，直接转账 ETH
                payable(msg.sender).transfer(pendingWithdrawAmount);
            } else {
                // 如果是其他代币池，转账代币
                IERC20(pool.stTokenAddress).safeTransfer(
                    msg.sender,
                    pendingWithdrawAmount
                );
            }
        }

        emit Withdraw(msg.sender, _pid, pendingWithdrawAmount, block.number);
    }

    /**
     * @notice Claim rewards from a pool. 领取质押奖励 （MetaNode 代币）
     * @dev This function updates the pool's reward state and allows the user to claim their pending rewards.
     */
    function claimRewards(
        uint256 _pid
    ) public checkPid(_pid) whenNotPaused whenNotPausedClaimRewards {
        // StakePool storage pool_ = pools[_pid];
        // User storage user_ = users[_pid][msg.sender];

        // updatePoolRewards(_pid);

        // uint256 pendingMetaNode_ = user_.stAmount * pool_.accMetaNodePerSt / (1 ether) - user_.finishedMetaNode + user_.pendToClaimRewards;

        // if(pendingMetaNode_ > 0) {
        //     user_.pendToClaimRewards = 0;
        //     MetaNode.safeTransfer(msg.sender, pendingMetaNode_);
        // }

        // user_.finishedMetaNode = user_.stAmount * pool_.accMetaNodePerSt / (1 ether);

        // emit ClaimRewards(msg.sender, _pid, pendingMetaNode_);

        StakePool storage pool = pools[_pid];
        User storage user = users[_pid][msg.sender];

        // console.log("Stake Balance:", address(this).balance);

        // 更新质押池奖励状态 主要是更新池的 accMetaNodePerSt
        updatePoolRewards(_pid);

        // 计算用户为领取的最新奖励
        uint256 pendToClaimRewards = (user.stAmount * pool.accMetaNodePerSt) /
            1 ether -
            user.finishedMetaNode +
            user.pendToClaimRewards;

        // console.log("claimRewards call");
        if (pendToClaimRewards > 0) {
            // 更新用户已领取的奖励
            user.pendToClaimRewards = 0;

            // console.log(
            //     "claimRewards call: pendToClaimRewards:",
            //     pendToClaimRewards
            // );
            // 转账奖励给用户
            MetaNode.safeTransfer(msg.sender, pendToClaimRewards);

            // console.log("claimRewards call");
        }

        // console.log("claimRewards call");

        user.finishedMetaNode =
            (user.stAmount * pool.accMetaNodePerSt) /
            1 ether;

        emit ClaimRewards(msg.sender, _pid, pendToClaimRewards);
    }

    // *************************************** QUERY FUNCTION ***************************************

    /**
     * @notice Get MetaNode rewards by given block range. 获取指定区块区间的挖矿奖励
     */
    function getRewardsByBolcks(
        uint256 _from,
        uint256 _to
    ) public view returns (uint256) {
        // require(_from <= _to, "invalid block");
        // if (_from < startBlock) {_from = startBlock;}
        // if (_to > endBlock) {_to = endBlock;}
        // require(_from <= _to, "end block must be greater than start block");
        // (bool success, uint256 multiplier) = (_to - _from).tryMul(metaNodePerBlock);
        // require(success, "multiplier overflow");
        // return multiplier;

        if (_from < startBlock) {
            _from = startBlock;
        }
        if (_to > endBlock) {
            _to = endBlock;
        }
        require(_from < _to, "Start block must be less than end block");

        (bool success, uint256 rewards) = (_to - _from).tryMul(
            metaNodePerBlock
        );
        require(success, "multiply overflow");
        return rewards;
    }

    /**
     * @notice Get the length/number of pool. 获取池的长度/数量
     */
    function pollLength() public view returns (uint256) {
        return pools.length;
    }

    /**
     * @notice Get pending MetaNode amount of user in pool
     */
    function getPendingMetaNode(
        uint256 _pid,
        address _user
    ) public view checkPid(_pid) returns (uint256) {
        return getPendingMetaNodeByBlockNumber(_pid, _user, block.number);
    }

    /**
     * @notice Get pending MetaNode amount of user by block number in pool
     */
    function getPendingMetaNodeByBlockNumber(
        uint256 _pid,
        address _user,
        uint256 _blockNumber
    ) public view checkPid(_pid) returns (uint256) {
        StakePool storage pool = pools[_pid];
        User storage user = users[_pid][_user];

        uint256 accMetaNodePerSt = pool.accMetaNodePerSt;

        if (_blockNumber > pool.lastRewardBlock && pool.stTokenAmount != 0) {
            // 当前区块高度大于池最后结算区块高度。
            // 此时必须先计算， 从最后奖励区块高度到当前区块高度的奖励， pool.accMetaNodePerSt 才准确。
            (bool success, uint256 rewards) = getRewardsByBolcks(
                pool.lastRewardBlock,
                _blockNumber
            ).tryMul(pool.poolWeight);
            require(success, "multiply overflow");

            (success, rewards) = rewards.tryDiv(totalPoolWeight);
            require(success, "divide overflow");

            accMetaNodePerSt =
                accMetaNodePerSt +
                (rewards * (1 ether)) /
                pool.stTokenAmount;
        }

        return
            (user.stAmount * accMetaNodePerSt) /
            1 ether -
            user.finishedMetaNode +
            user.pendToClaimRewards;
    }

    /**
     * @notice Get the staking amount of user
     */
    function getStakingBalance(
        uint256 _pid,
        address _user
    ) public view checkPid(_pid) returns (uint256) {
        return users[_pid][_user].stAmount;
    }

    /**
     * @notice Get the withdraw amount info, including the locked unstake amount and the unlocked unstake amount
     */
    function getWithdrawAmount(
        uint256 _pid,
        address _user
    )
        public
        view
        checkPid(_pid)
        returns (uint256 requestAmount, uint256 pendingWithdrawAmount)
    {
        User storage user = users[_pid][_user];

        for (uint256 index = 0; index < user.requests.length; index++) {
            if (user.requests[index].unlockBlock > block.number) {
                pendingWithdrawAmount += user.requests[index].amount;
            }
            requestAmount += user.requests[index].amount;
        }
    }

    // *************************************** INTERNAL FUNCTION ***************************************

    /**
     * @notice Internal function to handle deposit logic.
     * @dev This function updates the user's stake amount, the pool's total stake amount, and the user's pending rewards.
     */
    function _deposit(uint256 _pid, uint256 _amount) internal checkPid(_pid) {
        require(_amount > 0, "Deposit amount must be greater than 0");

        StakePool storage pool = pools[_pid];
        User storage user = users[_pid][msg.sender];

        // 更新质押池奖励状态
        updatePoolRewards(_pid);

        // 计算用户之前未领取的奖励 （pending rewards）
        if (user.stAmount > 0) {
            // console.log("1user.stAmount:", user.stAmount);
            // console.log("1user.finishedMetaNode:", user.finishedMetaNode);
            // console.log("1pool.accMetaNodePerSt:", pool.accMetaNodePerSt);

            (bool success, uint256 _pendToClaimRewards) = user.stAmount.tryMul(
                pool.accMetaNodePerSt
            );
            // console.log("1_pool.accMetaNodePerSt:", pool.accMetaNodePerSt);
            require(success, "overflow_d_1");
            (success, _pendToClaimRewards) = _pendToClaimRewards.tryDiv(
                1 ether
            );
            require(success, "overflow_d_2");

            // console.log("_pendToClaimRewards:", _pendToClaimRewards);
            // console.log("user.finishedMetaNode:", user.finishedMetaNode);
            (success, _pendToClaimRewards) = _pendToClaimRewards.trySub(
                user.finishedMetaNode
            );
            require(success, "overflow_d_3");

            if (_pendToClaimRewards > 0) {
                (success, _pendToClaimRewards) = _pendToClaimRewards.tryAdd(
                    user.pendToClaimRewards
                );
                require(success, "overflow_d_4");
                user.pendToClaimRewards = _pendToClaimRewards;
            }
        }

        // 更新用户质押金额
        (bool success1, uint256 newStAmount) = user.stAmount.tryAdd(_amount);
        require(success1, "overflow_d_5");
        user.stAmount = newStAmount;

        // 更新质押池的总质押量
        (bool success2, uint256 stTokenAmount) = pool.stTokenAmount.tryAdd(
            _amount
        );
        require(success2, "overflow_d_6");
        pool.stTokenAmount = stTokenAmount;

        // 更新用户应得但尚未领取的奖励的起点
        (bool success3, uint256 finishedMetaNode) = user.stAmount.tryMul(
            pool.accMetaNodePerSt
        );

        // console.log("finishedMetaNode3:", finishedMetaNode);
        require(success3, "overflow_d_7");
        (success3, finishedMetaNode) = finishedMetaNode.tryDiv(1 ether);

        // console.log("finishedMetaNode4:", finishedMetaNode);
        require(success3, "overflow_d_8");

        // console.log("user.finishedMetaNode2:", user.finishedMetaNode);
        user.finishedMetaNode = finishedMetaNode;

        emit Deposit(msg.sender, _pid, _amount);
    }
    // function _deposit(uint256 _pid, uint256 _amount) internal {
    //     StakePool storage pool_ = pools[_pid];
    //     User storage user_ = users[_pid][msg.sender];

    //     updatePoolRewards(_pid);

    //     if (user_.stAmount > 0) {
    //         // uint256 accST = user_.stAmount.mulDiv(pool_.accMetaNodePerST, 1 ether);
    //         (bool success1, uint256 accST) = user_.stAmount.tryMul(
    //             pool_.accMetaNodePerSt
    //         );
    //         require(success1, "user stAmount mul accMetaNodePerST overflow");
    //         (success1, accST) = accST.tryDiv(1 ether);
    //         require(success1, "accST div 1 ether overflow");

    //         (bool success2, uint256 pendingMetaNode_) = accST.trySub(
    //             user_.finishedMetaNode
    //         );
    //         require(success2, "accST sub finishedMetaNode overflow");

    //         if (pendingMetaNode_ > 0) {
    //             (bool success3, uint256 _pendingMetaNode) = user_
    //                 .pendToClaimRewards
    //                 .tryAdd(pendingMetaNode_);
    //             require(success3, "user pendingMetaNode overflow");
    //             user_.pendToClaimRewards = _pendingMetaNode;
    //         }
    //     }

    //     if (_amount > 0) {
    //         (bool success4, uint256 stAmount) = user_.stAmount.tryAdd(_amount);
    //         require(success4, "user stAmount overflow");
    //         user_.stAmount = stAmount;
    //     }

    //     (bool success5, uint256 stTokenAmount) = pool_.stTokenAmount.tryAdd(
    //         _amount
    //     );
    //     require(success5, "pool stTokenAmount overflow");
    //     pool_.stTokenAmount = stTokenAmount;

    //     // user_.finishedMetaNode = user_.stAmount.mulDiv(pool_.accMetaNodePerST, 1 ether);
    //     (bool success6, uint256 finishedMetaNode) = user_.stAmount.tryMul(
    //         pool_.accMetaNodePerSt
    //     );
    //     require(success6, "user stAmount mul accMetaNodePerST overflow");

    //     (success6, finishedMetaNode) = finishedMetaNode.tryDiv(1 ether);
    //     require(success6, "finishedMetaNode div 1 ether overflow");

    //     user_.finishedMetaNode = finishedMetaNode;

    //     emit Deposit(msg.sender, _pid, _amount);
    // }

    /**
     * @notice Safe transfer ETH to address, if the amount is greater than the contract balance, transfer all balance.
     * @dev This function is used to safely transfer ETH to an address, ensuring that it does not fail due to insufficient balance.
     */
    function safeTransferETH(address _to, uint256 _amount) internal {
        // // 这种写法是查询以太币余额
        // uint256 balance = address(this).balance;
        // 这种写法是查询 ERC20 代币余额
        uint256 balance = MetaNode.balanceOf(address(this));

        if (_amount > balance) {
            IERC20(_to).safeTransfer(_to, balance);
        } else {
            IERC20(_to).safeTransfer(_to, _amount);
        }
    }
}
