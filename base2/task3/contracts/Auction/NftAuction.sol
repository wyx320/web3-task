// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

// chainlink dataFeed
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
// ERC20 - MODIFIED: Using upgradeable interface
import "@openzeppelin/contracts-upgradeable/token/ERC20/IERC20Upgradeable.sol";
// // Ownable
// import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
// UUPS
// import "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";    // 包不能导入错了
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
// ERC721
import "@openzeppelin/contracts-upgradeable/token/ERC721/IERC721Upgradeable.sol";

contract NftAuction is UUPSUpgradeable, OwnableUpgradeable   {
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
    
    // 存储拍卖信息
    mapping(uint256 => Auction) public auctions;
    // 存储拍卖ID
    uint256 public nextAuctionId;

    // 创建拍卖
    function createAuction(uint256 _duration, uint256 _tokenId, uint256 _startPrice, address _nftContract) public {
        require(_nftContract != address(0), "nft contract address is zero");
        require(_startPrice > 0, "start price must be greater than zero");
        require(_duration > 0, "duration must be greater than zero");

        // 将 NFT 从卖家转移到拍卖合约
        IERC721Upgradeable(_nftContract).transferFrom(msg.sender, address(this), _tokenId);

        // TODO：这里可以用一个 mapping 存储正在拍卖的信息 来防止一个NFT被同时拍卖

        auctions[nextAuctionId] = Auction({
            seller: msg.sender,
            startPrice: _startPrice,
            startTime: block.timestamp,
            duration: _duration,
            ended: false,
            tokenId: _tokenId,
            highestBid: 0,
            highestBidder: address(0),
            paymentTokenAddress: address(0),
            nftContract: _nftContract
        });
        nextAuctionId++;
    }    

    // 初始化
    // constructor() Owner(msg.sender) {}   // 合约部署后 构造函数只会执行一次 UUPS可升级合约要求可进行多次部署
    function initialize() public initializer{
        __Ownable_init();  // 初始化 OwnableUpgradeable
        __UUPSUpgradeable_init();  // 添加这行
    }
    // constructor() {
    //     _disableInitializers();  // 禁用初始化函数，防止合约被再次初始化
    // }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    // 出价
    function bid(uint256 auctionId, uint256 amount, address _paymentTokenAddress) public payable {
        Auction storage auction = auctions[auctionId];
        // require(amount > auction.highestBid, "bid must be greater than the current highest bid");
        require(!auction.ended && block.timestamp < auction.startTime + auction.duration, "auction has ended");
        require(auction.seller != msg.sender, "seller cannot bid on their own auction");

        // 这里无论哪种方式出价，最终都会用 Chainlink 的 feedData 预言机获取 ERC20/以太坊 到美元 USD 的价格
        uint256 payValue = 0;
        if(_paymentTokenAddress == address(0)){
            // ETH 方式出价
            amount = msg.value;
            payValue = amount * getChainlinkDataFeedLatestAnswer(address(0));
        }
        else {
            // ERC20 方式出价
            payValue = amount * getChainlinkDataFeedLatestAnswer(_paymentTokenAddress);
        }

        uint256 _startPrice = auction.startPrice * getChainlinkDataFeedLatestAnswer(auction.paymentTokenAddress);
        uint256 _highestBid = auction.highestBid * getChainlinkDataFeedLatestAnswer(auction.paymentTokenAddress);    

        
        require(payValue >= _startPrice && payValue > _highestBid, "bid must be greater than the current highest bid");

        // 转移 ERC20 到合约
        if(_paymentTokenAddress != address(0)){
            IERC20Upgradeable erc20Token = IERC20Upgradeable(_paymentTokenAddress);
            erc20Token.transferFrom(msg.sender, address(this), amount);
        }
        // 注意：以太坊不需要转移 因为函数有 payable 修饰符 ETH 在调用时自动转入合约余额

        // 退还前最高出价
        if(auction.highestBidder != address(0)) { // Added check to prevent sending to address(0)
            if(auction.paymentTokenAddress == address(0)){
                // 前最高价是 ETH 方式出价
                payable(auction.highestBidder).transfer(auction.highestBid);
            }
            else {
                // 前最高价是 ERC20 出价
                IERC20Upgradeable erc20Token = IERC20Upgradeable(auction.paymentTokenAddress);
                erc20Token.transfer(auction.highestBidder, auction.highestBid);
            }
        }

        // 更新拍卖信息
        auction.paymentTokenAddress = _paymentTokenAddress;
        auction.highestBidder = msg.sender;
        auction.highestBid = amount;
    }

    // 存储 chainLink PriceFeed 兑换币 => USD  交易对地址信息信息
    mapping (address => AggregatorV3Interface) public priceFeeds;

    function setChainlinkDataFeedLatestAnswer(address tokenAddress, address feedAddress) public {
    priceFeeds[tokenAddress] = AggregatorV3Interface(feedAddress);
}

    /**
     * Returns the latest answer.
     */
    function getChainlinkDataFeedLatestAnswer(address paymentTokenAddress) public view returns (uint256) {
        AggregatorV3Interface priceFeed = priceFeeds[paymentTokenAddress];
        // prettier-ignore
        (
            /* uint80 roundId */,
            int256 answer,
            /*uint256 startedAt*/,
            /*uint256 updatedAt*/,
            /*uint80 answeredInRound*/
        ) = priceFeed.latestRoundData();
        return uint256(answer);
    }

    // 结束拍卖
    function endAuction(uint256 auctionId) external {
        Auction storage auction = auctions[auctionId];
        require(auction.seller == msg.sender, "only the seller can end the auction");
        require(!auction.ended && block.timestamp <= auction.startTime + auction.duration, "auction has already ended");
        
        auction.ended = true;   // 先标记结束，防止重入


        if(auction.highestBidder != address(0)) {
            // 有人出价
            if(auction.paymentTokenAddress == address(0)){
                // 转移 ETH 到卖家
                payable(auction.seller).transfer(auction.highestBid);
            }
            else{
                // 转移 ERC20 到卖家
                IERC20Upgradeable erc20Token = IERC20Upgradeable(auction.paymentTokenAddress);
                erc20Token.transfer(auction.seller, auction.highestBid);
            }
            // 转移 NFT 到最高出价者
            IERC721Upgradeable nft = IERC721Upgradeable(auction.nftContract);
            nft.safeTransferFrom(address(this), auction.highestBidder, auction.tokenId);
        }
        else
        {
            // 没人出价，退还 NFT 给卖家
            IERC721Upgradeable nft = IERC721Upgradeable(auction.nftContract);
            nft.safeTransferFrom(address(this), auction.seller, auction.tokenId);
        }

    }
}