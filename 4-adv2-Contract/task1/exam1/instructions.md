这是一份操作指南，指导如何部署和使用代币合约，包括如何进行代币交易、添加和移动流动性等操作。

## 1. 部署代币合约

1. 使用 vscode 打开 .\\web3-task\4-adv2-Contract\task1\exam1 文件夹

2. Terminal 输入 remixd 打开 websocket

3. 网页端 remix 使用 remixd 插件连接 vscode

4. 编译 ShibStyleToken 合约

5. remix 部署页面选择 ShibStyleToken 合约

6. 填写部署参数，提交合约部署交易。部署参数中文注释：

   ```solidity
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
   ```

## 2. 使用代币合约

1. 合约使用固定税率。transfer 转账、transferFrom 代理转账，都会根据转账金额，额外从账户中扣除 转账金额百分之二 的税率金额。
2. addLiquidity、removeLiquidity 函数可以从 IUniswap 中添加/移除流动性。