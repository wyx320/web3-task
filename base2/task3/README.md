项目结构：
--contracts                 合约文件
    --Auction               拍卖合约
        --NftAuction.sol    V1 版本
        --NftAuctionV2.sol  V2 版本
    --NFT                   NFT 合约
        --MyNFT.sol         V1 版本
    --MockV3Aggregator.sol  Mock 合约
--deploy                    部署脚本
    --.cache                缓存文件
    --deploy.js             部署脚本
    --upgrade.js            升级脚本
--test                      测试文件
    --localohst.js          本地测试
    --sepolia.js            Sepolia测试

功能说明：
1. 拍卖合约。允许用户创建拍卖、出价、结束拍卖、转移NFT、转移ETH、转移ERC20代币。
2. NFT合约。允许用户铸造NFT、转移NFT。
3. Mock合约。用于模拟Chainlink价格。
4. 部署脚本。用于部署合约。
5. 测试脚本。用于测试合约。
6. 支持部署到 Sepolia 测试网进行测试。
7. 支持UUPS升级合约。

未完成功能：
1. 跨链拍卖
2. 工厂模式
3. 动态手续费

部署步骤：
1. 本地部署：npx hardhat test .\test\localhost.js
2. 部署到Sepolia测试网：
    1. developer.metamask.io 创建账户 获取 Sepolia 测试网 RPC URL 填入
    2. 填入 hardhat.config.js =》sepolia =》url
    3. MetaMask 钱包获取私钥
    4. 填入 hardhat.config.js =》sepolia =》accounts
    5. npx hardhat deploy --network sepolia
    6. 拿到拍卖合约地址与NFT地址，填入 .\test\sepolia.js
    7. npx hardhat test .\test\sepolia.js