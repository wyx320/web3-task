const { deployments, ethers } = require("hardhat");
const { expect } = require("chai");

describe("Test Auction Upgrade", async () => {
    it("Should be able to upgrade.", async()=>{
        // 部署 ERC721 合约
        const erc721Factory = await ethers.getContractFactory("MyNFT");
        const erc721 = await erc721Factory.deploy();
        erc721.waitForDeployment();
        // console.log("===erc721Contract.address===", await erc721.getAddress());

        // 部署业务合约 NftAuction
        await deployments.fixture("deployAuction");
        const proxy = await deployments.get("AuctionProxy");
        // console.log("===proxy===", proxy.address);
        // const allDeployments = await deployments.all();
        // console.log("📦 所有已部署合约:", Object.keys(allDeployments));
        const auction = await ethers.getContractAt("NftAuction", proxy.address);
        // console.log("===auction address===", await auction.getAddress());
        // console.log("Contract ABI:\n", auction.interface.formatJson());
        
        // 测试 ERC721 合约功能
        const [ seller, buyer1 ] = await ethers.getSigners();
        // 给自己 mint 10个 NFT
        for (let i = 0; i < 10; i++) {
            await erc721.mintNFT(seller, i);
          }

        // 部署 Mock Chainlink Aggregator（模拟 ETH/USD）
        const MockV3Aggregator = await ethers.getContractFactory("MockV3Aggregator");
        const mockFeed = await MockV3Aggregator.deploy(8, 2000e8); // $2000 USD / ETH
        await mockFeed.waitForDeployment();
    
        // 设置 ETH 的喂价器为这个 Mock 合约
        await auction.setChainlinkDataFeedLatestAnswer(
            ethers.ZeroAddress,         // 表示 ETH
            await mockFeed.getAddress() // Mock 喂价器地址
        );

        // 测试 V1 合约功能
        // NFT 所有者把 NFT 授权给业务合约
        const auctionAddress = await auction.getAddress();
        console.log("===auctionAddress===", auctionAddress);
        await erc721.connect(seller).setApprovalForAll(auctionAddress, true);
        // // 检查授权是否成功
        // const isApproved = await erc721.isApprovedForAll(seller, auctionAddress);
        // if (isApproved) {
        //     console.log("授权成功");
        // } else {
        //     console.log("授权失败");
        // }

        // 创建拍卖
        const tokenId = 3;
        const erc721Address = await erc721.getAddress();
        await auction.createAuction(1000*10, tokenId, ethers.parseEther("0.01"), erc721Address);
        // const result = await auction.createAuction(1000*10, tokenId, ethers.parseEther("0.01"), auctionAddress);
        // console.log(result);
        // 出价
        console.log("===0===");
        await auction.connect(buyer1).bid(0, 0, ethers.ZeroAddress,{
            value: ethers.parseEther("0.01"),
        })
        console.log("===1===");
        // console.log(await auction.auctions(0));
        // console.log(await auction.auctions(1));
        // 等待拍卖结束
        await new Promise((resolve) => setTimeout(resolve, 11000));
        console.log("===2===");
        await auction.endAuction(0);
        console.log("===3===");

        // 验证拍卖是否结束
        const auctionResult = await auction.auctions(0);
        expect(auctionResult.highestBidder).to.equal(buyer1.address);
        expect(auctionResult.highestBid).to.equal(ethers.parseEther("0.01"));

        // 验证 NFT 所有权
        const owner = await erc721.ownerOf(tokenId);
        expect(owner).to.equal(auctionResult.highestBidder);

        // 升级合约到 NftAuctionV2
        await deployments.fixture("UpgradeAuction");
        const auctionV2 = await ethers.getContractAt(
            "NftAuctionV2",
            proxy.address
          );
        // console.log("===auctionV2===", await auctionV2.getAddress());
        // 测试 V2 合约功能
        await auctionV2.setValue("Hello");
        expect(await auctionV2.getValue()).to.equal("Hello");
        console.log(await auctionV2.auctions(0));
    });
});