const { ethers, upgrades } = require("hardhat");

const path = require("path");
const fs = require("fs");

module.exports = async ({ deployments, getNamedAccounts }) => {
  const factory = await ethers.getContractFactory("NftAuction");
  // const [deployer] = await getNamedAccounts();

  const proxy = await upgrades.deployProxy(factory, [], {
    kind: "uups",
    initializer: "initialize",
  });
  proxy.waitForDeployment();

  // 获取代理地址
  const proxyAddress = await proxy.getAddress();
  // 获取合约地址
  const implAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);

  fs.writeFileSync(
    path.resolve(__dirname, "./.cache/AuctionProxy.json"),
    JSON.stringify({
      proxyAddress: proxyAddress,
      implAddress: implAddress,
    })
  )

  const { save } = deployments;
  console.log("===deployments parameter===", Object.keys(deployments)); // 查看所有方法名
  await save("AuctionProxy", {
    address: proxyAddress,
    abi: factory.interface.format("json"),
    args: [],
    log: true,
  })
};

module.exports.tags = ["deployAuction"];
