const { ethers } = require("hardhat");
const { expect } = require("chai");

// =======================================================
// ============== Sepolia 测试网配置 ======================
// =======================================================
const SEPOLIA_NFT_CONTRACT_ADDRESS = "0x...YourDeployedMyNFTContractAddress..."; // 替换成你部署的 NFT 合约地址
const SEPOLIA_AUCTION_PROXY_ADDRESS = "0x...YourDeployedAuctionProxyAddress..."; // 替换成你部署的拍卖合约代理地址

// Sepolia 上真实的 ETH/USD Price Feed 地址
const SEPOLIA_ETH_USD_PRICE_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
// =======================================================

describe("Test Auction on Sepolia", async () => {
    let auction, erc721;
    let seller, buyer1;

    before(async () => {
        // 获取签名者
        [seller, buyer1] = await ethers.getSigners();
        console.log("Seller address:", seller.address);

        // === 不再部署，而是连接到已部署的合约 ===
        erc721 = await ethers.getContractAt("MyNFT", SEPOLIA_NFT_CONTRACT_ADDRESS);
        auction = await ethers.getContractAt("NftAuction", SEPOLIA_AUCTION_PROXY_ADDRESS);

        console.log("Connected to NFT Contract at:", await erc721.getAddress());
        console.log("Connected to Auction Contract at:", await auction.getAddress());
    });

    it("Should be able to run a full auction cycle on Sepolia", async function () {
        this.timeout(300000); // 将测试超时时间延长到 5 分钟，因为测试网很慢

        // === 1. 设置 Chainlink Price Feed (如果尚未设置) ===
        console.log("Setting price feed...");
        const setFeedTx = await auction.setChainlinkDataFeedLatestAnswer(
            ethers.ZeroAddress,
            SEPOLIA_ETH_USD_PRICE_FEED
        );
        await setFeedTx.wait(1); // 等待交易被确认
        console.log("Price feed set.");

        // === 2. 准备 NFT (mint 和授权) ===
        const tokenId = Math.floor(Math.random() * 100000); // 使用一个随机的 tokenId 避免冲突
        console.log(`Minting NFT with tokenId: ${tokenId}...`);
        const mintTx = await erc721.mintNFT(seller.address, tokenId);
        await mintTx.wait(1);
        console.log("NFT minted. Owner:", await erc721.ownerOf(tokenId));

        console.log("Approving auction contract...");
        const approvalTx = await erc721.connect(seller).setApprovalForAll(await auction.getAddress(), true);
        await approvalTx.wait(1);
        const isApproved = await erc721.isApprovedForAll(seller.address, await auction.getAddress());
        expect(isApproved).to.be.true;
        console.log("Approval successful.");

        // === 3. 创建拍卖 ===
        const auctionDuration = 120; // 设置一个短的拍卖时间，比如 120 秒
        console.log(`Creating auction for tokenId ${tokenId} with duration ${auctionDuration} seconds...`);
        const createTx = await auction.connect(seller).createAuction(
            auctionDuration,
            tokenId,
            ethers.parseEther("0.0001"), // 使用一个非常低的价格
            await erc721.getAddress()
        );
        await createTx.wait(1);
        const newAuction = await auction.auctions(0); // 假设这是第一个拍卖
        console.log("Auction created. End time:", new Date(Number(newAuction.startTime) * 1000 + auctionDuration * 1000).toLocaleString());

        // === 4. 出价 ===
        console.log("Buyer 1 is bidding...");
        const bidTx = await auction.connect(buyer1).bid(
            0, // auctionId
            0, // amount for ERC20, ignored for ETH
            ethers.ZeroAddress,
            { value: ethers.parseEther("0.0001") }
        );
        await bidTx.wait(1);
        console.log("Bid placed successfully.");
        const auctionStateAfterBid = await auction.auctions(0);
        expect(auctionStateAfterBid.highestBidder).to.equal(buyer1.address);

        // === 5. 等待拍卖结束 ===
        console.log(`Waiting for ${auctionDuration} seconds for auction to end...`);
        await new Promise(resolve => setTimeout(resolve, auctionDuration * 1000));
        console.log("Wait time finished.");

        // === 6. 结束拍卖 ===
        // 注意：在真实网络上，可能需要等待下一个区块才能让 `block.timestamp` 更新
        // 为了稳定，我们再多等一小会儿
        await new Promise(resolve => setTimeout(resolve, 15000)); // 等待 15 秒

        console.log("Ending auction...");
        const endTx = await auction.connect(seller).endAuction(0); // 卖家结束拍卖
        await endTx.wait(1);
        console.log("Auction ended.");

        // === 7. 验证结果 ===
        const finalOwner = await erc721.ownerOf(tokenId);
        console.log(`Final owner of NFT ${tokenId} is: ${finalOwner}`);
        expect(finalOwner).to.equal(buyer1.address);
        console.log("✅ Test successful! NFT was correctly transferred to the winner.");
    });
});