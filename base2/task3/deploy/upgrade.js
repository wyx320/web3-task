const { ethers, upgrades } = require("hardhat");

const path = require("path");
const fs = require("fs");

module.exports = async({ deployments }) => {
    const data = await fs.readFileSync(
        path.resolve(__dirname, "./.cache/AuctionProxy.json"),
        "utf8"
    );
    const { proxyAddress, implAddress } = JSON.parse(data);

    const factory = await ethers.getContractFactory("NftAuctionV2");
    const contract = await upgrades.upgradeProxy(proxyAddress, factory);
    await contract.waitForDeployment();

    fs.writeFileSync(
        path.resolve(__dirname, "./.cache/AuctionProxy.json"),
        JSON.stringify({
            proxyAddress: proxyAddress,
            implAddress: implAddress,
        })
    )

    const { save } = deployments;
    await save("AuctionProxy", {
        address: proxyAddress,
        abi: factory.interface.format("json"),
        args: [],
        log: true,
    })
};

module.exports.tags = ["UpgradeAuction"];
