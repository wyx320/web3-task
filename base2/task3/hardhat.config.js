require("@nomicfoundation/hardhat-toolbox");
require("hardhat-deploy");  // 用于支持在测试用例中使用 hardhat.deployments.fixture
require('@openzeppelin/hardhat-upgrades');

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  namedAccounts: {
    deployer: 0,
    user1: 1,
    user2: 2,
  },

    // 配置所有支持的网络
    networks: {
      // 本地开发网络
      hardhat: {
        chainId: 31337,
      },
      // Sepolia 测试网配置
      sepolia: {
        url: "https://sepolia.infura.io/v3/e48f7acfd92045759d9243f878973e8a",     // RPC 地址
        accounts: ["92c4ef37898a9f64b32e1fd5868a946c81202d15788113175d7fcbb43aacbd8c"],  // 用于部署的账户私钥
        chainId: 11155111,        // Sepolia 的 Chain ID
        saveDeployments: true,    // 部署时保存部署信息
      },
    },
};
