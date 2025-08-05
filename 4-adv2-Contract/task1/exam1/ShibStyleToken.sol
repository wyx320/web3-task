// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./node_modules/@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./node_modules/@openzeppelin/contracts/token/ERC20/IERC20.sol";

// 简单模拟流动性池合约接口
interface IUniswapV2Pair {
    function mint(address to) external returns (uint256 liquidity);
    function burn(
        address to
    ) external returns (uint256 amount0, uint256 amount1);
}

contract ShibStyleToken is ERC20 {
    // constructor() ERC20("ShibStyleToken", "SST") {}

    // 代币税税率，例如 2% 表示 200（2 * 100）
    uint256 public taxRate = 200;
    // 代币税收接收地址
    address public taxReceiver;
    // 单笔交易最大金额
    uint256 public maxTransactionAmount;
    // 每日交易次数限制
    uint256 public dailyTransactionCountLimit;
    // 用户每日交易次数记录
    mapping(address => uint256) dailyTransactionCount;
    // 用户每日交易重置时间
    mapping(address => uint256) lastTransactionDay;

    // 流动性池地址
    address public liquidityPool;

    constructor(
        string memory name,
        string memory symbol,
        uint256 _initialSupply,
        uint256 _maxTransactionAmount,
        uint256 _dailyTransactionCountLimit,
        address _taxReceiver,
        address _liquidityPool
    ) ERC20(name, symbol) {
        _mint(msg.sender, _initialSupply);
        maxTransactionAmount = _maxTransactionAmount;
        dailyTransactionCountLimit = _dailyTransactionCountLimit;
        taxReceiver = _taxReceiver;
        liquidityPool = _liquidityPool;
    }

    // 重写 transfer 函数， 实现代币税和交易限制
    function transfer(
        address recipient,
        uint256 amount
    ) public override returns (bool) {
        address sender = msg.sender;
        // 检查交易限制
        _checkTransactionLimits(sender, amount);
        // 计算税额
        uint256 taxAmount = _calculateTax(amount);
        // 扣除税额
        uint256 totalAmount = amount + taxAmount;
        require(
            balanceOf(sender) >= totalAmount,
            "Insufficient balance including tax"
        );
        // 转账原始金额到接收者
        super.transfer(recipient, amount);
        // 转账税额到税收接收地址
        super.transfer(taxReceiver, taxAmount);

        return true;
    }

    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) public override returns (bool) {
        // 先检查 allowance
        uint256 currentAllowance = allowance(sender, msg.sender);
        require(
            currentAllowance >= amount,
            "ERC20: transfer amount exceeds allowance"
        );

        // 检查交易限制
        _checkTransactionLimits(sender, amount);
        // 计算税额
        uint256 taxAmount = _calculateTax(amount);
        // 计算需要扣除的总金额（包含税额）
        uint256 totalAmount = amount + taxAmount;
        require(
            balanceOf(sender) >= totalAmount,
            "Insufficient balance including tax"
        );
        // 使用内部 _transfer 函数避免重复检查 allowance
        // 转账原始金额到接收者
        _transfer(sender, recipient, amount);
        // 转账税额到税收接收地址
        _transfer(sender, taxReceiver, taxAmount);
        // 更新 allowance（只减少请求的金额）
        _approve(sender, msg.sender, currentAllowance - amount);

        return true;
    }

    function _calculateTax(uint256 amount) internal view returns (uint256) {
        return (amount * taxRate) / 10000; // 计算税额
    }

    // 检查交易限制
    function _checkTransactionLimits(address sender, uint256 amount) internal {
        // 检查单笔交易金额限制
        require(
            amount <= maxTransactionAmount,
            "Transaction amount exceeds limit"
        );

        uint256 currentDay = block.timestamp / 1 days;
        uint256 lastDay = lastTransactionDay[sender];
        if (currentDay > lastDay) {
            // 重置每日交易次数
            dailyTransactionCount[sender] = 0;
            lastTransactionDay[sender] = currentDay;
        }

        // 检查发送者每日交易次数限制
        require(
            dailyTransactionCount[sender] < dailyTransactionCountLimit,
            "Daily transaction limit exceeded"
        );

        // 更新每日交易次数
        dailyTransactionCount[sender]++;
    }

    // 添加流动性
    function addLiquidity(
        uint256 tokenAmount,
        uint256 ethAmount
    ) external payable {
        require(liquidityPool != address(0), "Liquidity pool not set");
        require(ethAmount == msg.value, "ETH amount mismatch");

        // 将代币转移到流动性池
        // 使用 _transfer 而不是 this.transfer 避免税收
        _transfer(msg.sender, liquidityPool, tokenAmount);

        // 将 ETH 发送到流动性池
        (bool success, ) = liquidityPool.call{value: ethAmount}("");
        require(success, "ETH transfer failed");

        // 调用流动性池合约的 mint 函数 添加流动性
        IUniswapV2Pair(liquidityPool).mint(msg.sender);
    }

    // 移除流动性
    function removeLiquidity(
        uint256 liquidityTokens
    ) external returns (uint256 amount0, uint256 amount1) {
        require(liquidityPool != address(0), "Liquidity pool not set");

        // 将流动池代币转移到流动性池（假设用户已经拥有这些代币）
        IERC20(liquidityPool).transferFrom(
            msg.sender,
            liquidityPool,
            liquidityTokens
        );

        // 调用流动性池合约的 burn 函数 移除流动性
        (amount0, amount1) = IUniswapV2Pair(liquidityPool).burn(msg.sender);
    }

    // // 设置流动性池地址
    // function setLiquidityPool(address _liquidityPool) external {
    //     require(msg.sender == taxReceiver, "Only tax receiver can set pool");
    //     require(_liquidityPool != address(0), "Invalid pool address");
    //     liquidityPool = _liquidityPool;
    // }
}
