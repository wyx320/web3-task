// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import { AggregatorV3Interface } from "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract NftAuction {
    struct Auction {
        // 拍卖人
        address seller;
        // 起拍价
        uint256 startPrice;
        // 开始时间
        uint256 startTime;
        // 持续时间
        uint256 duration;
        // 是否结束
        bool ended;
        // NFT Identifier
        uint256 tokenId;
        // 最高出价
        uint256 highestBid;
        // 最高出价人
        address highestBidder;
        // 参与竞价的资产类型 0x地址表示eth 其他地址表示erc20
        address paymentTokenAddress;
        // NFT 拍卖合约地址
        address nftContract;
    }

    constructor() Owner(msg.sender) {}

    // 存储拍卖信息
    mapping(uint256 => Auction) public auctions;
    // 存储拍卖ID
    uint256 public nextAuctionId;

    // 创建拍卖
    function createAuction(uint256 _duration, uint256 _tokneId, uint256 _startPrice, address _nftContract) public {
        require(_nftContract != address(0), "nft contract address is zero");
        require(_startPrice > 0, "start price must be greater than zero");
        require(_duration > 0, "duration must be greater than zero");

        // TODO：这里可以用一个 mapping 存储正在拍卖的信息 来防止一个NFT被同时拍卖

        Auction memory auction = Auction({
            seller: msg.sender,
            startPrice: _startPrice,
            startTime: block.timestamp,
            duration: _duration,
            ended: false,
            tokenId: _tokneId,
            highestBid: _startPrice,
            highestBidder: address(0),
            paymentTokenAddress: address(0),
            nftContract: _nftContract
        });
        auctions[nextAuctionId] = auction;
        nextAuctionId++;
    }

    // 出价
    function bid(uint256 auctionId, uint256 amount, address paymentTokenAddress) public payable {
        Auction storage auction = auctions[auctionId];
        // require(amount > auction.highestBid, "bid must be greater than the current highest bid");
        require(!auction.ended && block.timestamp < auction.startTime + auction.duration, "auction has ended");
        require(auction.seller != msg.sender, "seller cannot bid on their own auction");

        // 这里无论哪种方式出价，最终都会用 Chainlink 的 feedData 预言机获取 ERC20/以太坊 到美元 USD 的价格
        uint256 payValue = 0;
        if(paymentTokenAddress == address(0)){
            // ETH 方式出价
            amount = msg.value;

            payValue = amount * getChainlinkDataFeedLatestAnswer(address(0));
        }
        else {
            // ERC20 方式出价

            payValue = amount * getChainlinkDataFeedLatestAnswer(paymentTokenAddress);
        }

        uint256 startPrice = auction.startPrice * getChainlinkDataFeedLatestAnswer(auction.paymentTokenAddress);
        uint256 highestBid = auction.highestBid * getChainlinkDataFeedLatestAnswer(auction.paymentTokenAddress);    

        require(payValue > startPrice && payValue > highestBid, "bid must be greater than the current highest bid");

        // 转移 ERC20 到合约
        if(paymentTokenAddress != address(0)){
                IERC20(paymentTokenAddress).approve(address(this),amount);  // 用户授权当前合约可以操作其代币
                IERC20(paymentTokenAddress).transferFrom(msg.sender, address(this), amount);
        }
        // 注意：以太坊不需要转移 因为函数有 payable 修饰符 ETH 在调用时自动转入合约余额

        // 退还前最高出价
        if(paymentTokenAddress == address(0)){
            // 前最高价是 ETH 方式出价
            payable(auction.highestBidder).transfer(auction.highestBid);
        }
        else {
            // 前最高价是 ERC20 出价
            IERC20(auction.paymentTokenAddress).transfer(auction.highestBidder, auction.highestBid);
        }

        // 更新拍卖信息
        auction.paymentTokenAddress = paymentTokenAddress;
        auction.highestBidder = msg.sender;
        auction.highestBid = amount;
    }

    // 存储 chainLink PriceFeed 兑换币 => USD  交易对地址信息信息
    mapping (address => AggregatorV3Interface) public priceFeeds;

    // 设置 Chainlink 喂价
    function setChainlinkDataFeedLatestAnswer(address _priceETHFeed) {
        priceFeeds[_priceETHFeed] = AggregatorV3Interface(_priceETHFeed);
    }

    /**
     * Returns the latest answer.
     */
    function getChainlinkDataFeedLatestAnswer(address paymentTokenAddress) public view returns (int) {
        priceFeed = priceFeeds[paymentTokenAddress];
        // prettier-ignore
        (
            /* uint80 roundId */,
            int256 answer,
            /*uint256 startedAt*/,
            /*uint256 updatedAt*/,
            /*uint80 answeredInRound*/
        ) = dataFeed.latestRoundData();
        return answer;
    }

    // 结束拍卖
    function endAuction(uint256 auctionId) external onlyOwner {
        Auction storage auction = auctions[auctionId];
        require(auction.seller == msg.sender, "only the seller can end the auction");
        require(!auction.ended && block.timestamp > auction.startTime + auction.duration, "auction has already ended");

        // 转移 ERC20 到卖家
        if(auction.paymentTokenAddress == address(0)){
            payable(auction.seller).transfer(auction.highestBid);
        }
        else{
            IERC20(auction.nftContract).safeTransferFrom(address(this), auction.highestBidder, auction.tokenId);
        }

        auction.ended = true;
    }
}
