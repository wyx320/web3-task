require("@nomicfoundation/hardhat-toolbox");
require("hardhat-deploy");  // 用于支持在测试用例中使用 hardhat.deployments.fixture
require("@openzeppelin/hardhat-upgrades");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  namedAccounts: {
    deployer: 0,
    user1: 1,
    user2: 2,
  }
};
