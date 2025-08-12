const { ethers, network } = require("hardhat");
const { expect } = require("chai");

// 快速挖到指定区块高度
// async function mineToBlock(targetBlock) {
//   const currentBlock = await ethers.provider.getBlockNumber();
//   if (targetBlock > currentBlock) {
//     const diff = targetBlock - currentBlock;
//     await network.provider.send("hardhat_mine", [`0x${diff.toString(16)}`]);
//   }
// }
async function mineToBlock(num) {
  for (let i = 0; i < num; i++) {
    await network.provider.send("evm_mine");
  }
}

describe("MetaNodeStake (Hardhat JS tests)", function () {
  let metaNode, metaNodeStake;
  let owner;

  beforeEach(async function () {
    [owner] = await ethers.getSigners();
    // 部署代币合约
    var metaNodeFactory = await ethers.getContractFactory("MetaNode");
    metaNode = await metaNodeFactory.deploy();
    await metaNode.waitForDeployment();

    // 部署质押合约
    var metaNodeStakeFactory = await ethers.getContractFactory("MetaNodeStake");
    metaNodeStake = await metaNodeStakeFactory.deploy();
    await metaNodeStake.waitForDeployment();

    // 初始化质押合约
    await metaNodeStake.initialize(
      100,
      10000000,
      ethers.parseEther("3"),
      metaNode.getAddress()
    );
  });

  async function addNativePool() {
    await metaNodeStake.addPool(
      ethers.ZeroAddress, // native currency pool
      100, // poolWeight
      100, // minDepositAmount
      100, // withDrawLockedBlocks
      true // withUpdate
    );
  }

  it("AddPool: 添加原生代币迟并检查初始状态", async function () {
    // console.log(await ethers.provider.getBlockNumber());

    await addNativePool();
    const pool = await metaNodeStake.pools(0);
    console.log("pool:", pool);
    expect(pool.stTokenAddress).to.equal(ethers.ZeroAddress);
    expect(pool.poolWeight).to.equal(100n);
    expect(pool.minDepositAmount).to.equal(100n);
    expect(pool.unstakeLockedBlocks).to.equal(100n);
    expect(pool.lastRewardBlock).to.equal(100n); // 质押合约从高度 100 开始
    expect(pool.accMetaNodePerSt).to.equal(0n);
    expect(pool.stTokenAmount).to.equal(0n);
  });

  it("massUpdatePoolsRewards: 在区块推进后更新 lastRewardBlock", async function () {
    await addNativePool();
    let pool = await metaNodeStake.pools(0);
    expect(pool.lastRewardBlock).to.equal(100n);

    // 获取当前区块高度
    const currentBlockNumber = await ethers.provider.getBlockNumber();
    // console.log("Current block number:", currentBlockNumber);

    // 模拟区块推进到 200
    await mineToBlock(200);
    // console.log(
    //   "Block number after mining:",
    //   await ethers.provider.getBlockNumber()
    // );

    const tx = await metaNodeStake.massUpdatePoolsRewards();
    const { blockNumber } = await tx.wait();
    pool = await metaNodeStake.pools(0);
    // console.log("Last reward block after update:", pool.lastRewardBlock);

    // expect(pool.lastRewardBlock).to.equal(200n);

    // 把区块高度推到 200 后, 再发起一笔交易, 用来调用 massUpdatePoolsRewards, 这笔交易会被打进"下一个"区块. 也就是说:
    // · mineToBlock(200) 当前链头是 200.
    // · 随后发起的写交易会被挖矿在下一个新区块里, 也就是 201.
    // · 合约里调用 block.number 是 "承载这笔交易的区块号", 因此 lastRewardBlock 被更新到 201 是符合预期的.
    expect(pool.lastRewardBlock).to.gt(200n);
    // 改进建议:
    // 在测试里用交易回执来断言，避免“我以为是 300，但其实交易在 301 才执行”的错觉：
    expect(pool.lastRewardBlock).to.equal(BigInt(blockNumber));
  });

  it("SetPoolWeight: 修改池权重并检查 totalWeight", async function () {
    await addNativePool();
    let preTotalWeight = await metaNodeStake.totalPoolWeight();

    await metaNodeStake.updatePoolWeight(0, 200, false);
    let pool = await metaNodeStake.pools(0);
    let totalPoolWeight = await metaNodeStake.totalPoolWeight();

    // console.log("totalPoolWeight:", totalPoolWeight);
    // console.log("preTotalWeight:", preTotalWeight);

    expect(pool.poolWeight).to.equal(200n);
    expect(totalPoolWeight).to.equal(preTotalWeight + 100n); // 之前是 100，现在变成 200，所以增加了 100
  });

  it("Deposit native currency: 存入原生代币并检查用户与质押池状态", async function () {
    await addNativePool();

    // 读取初始状态
    let pool = await metaNodeStake.pools(0);
    let preStTotalAmount = pool.stTokenAmount;
    let user = await metaNodeStake.users(0, owner.address);
    let preStAmount = user.stAmount;
    let preFinishedMetaNode = user.finishedMetaNode;

    // 第一次小额存入 100 wei
    await metaNodeStake.depositETH({ value: 100 });
    pool = await metaNodeStake.pools(0);
    user = await metaNodeStake.users(0, owner.address);
    expect(pool.stTokenAmount).to.equal(preStTotalAmount + 100n);
    expect(user.stAmount).to.equal(preStAmount + 100n);
    expect(user.finishedMetaNode).to.equal(preFinishedMetaNode);

    // 更多存入与解押序列
    await metaNodeStake.depositETH({ value: 200 });
    await mineToBlock(2000);
    await metaNodeStake.unstake(0, 100);

    await metaNodeStake.depositETH({ value: 300 });
    await mineToBlock(2000);
    await metaNodeStake.unstake(0, 100);

    await metaNodeStake.depositETH({ value: 400 });
    await mineToBlock(2000);
    await metaNodeStake.unstake(0, 100);

    await metaNodeStake.depositETH({ value: 500 });
    await mineToBlock(2000);
    await metaNodeStake.unstake(0, 100);

    await metaNodeStake.depositETH({ value: 600 });
    await mineToBlock(2000);
    await metaNodeStake.unstake(0, 100);

    await metaNodeStake.depositETH({ value: 700 });
    await mineToBlock(2000);
    await metaNodeStake.unstake(0, 100);

    await metaNodeStake.withdraw(0);
  });

  it("Unstake: 解押后检查用户与质押池状态", async function () {
    await addNativePool();
    await metaNodeStake.depositETH({ value: ethers.parseEther("1") });

    // 推进区块高度
    await mineToBlock(1000);
    await metaNodeStake.unstake(0, ethers.parseEther("0.3"));

    // 检查用户状态
    const user = await metaNodeStake.users(0, owner.address);
    expect(user.stAmount).to.equal(ethers.parseEther("0.7"));
    // expect(user.finishedMetaNode).to.equal(
    console.log("user.finishedMetaNode:", user.finishedMetaNode);
    expect(user.pendToClaimRewards).to.be.gt(0n);

    // 检查质押池状态
    const pool = await metaNodeStake.pools(0);
    expect(pool.stTokenAmount).to.equal(ethers.parseEther("0.7"));
  });

  it("Withdraw: 解质押后提取原生代币", async function () {
    await addNativePool();
    await metaNodeStake.depositETH({ value: ethers.parseEther("1") });
    await metaNodeStake.unstake(0, ethers.parseEther("0.3"));

    // 推进区块高度，超过 withdrawLockedBlocks
    await mineToBlock(2000);

    const preUserBalance = await ethers.provider.getBalance(owner.address);

    //提取原生代币
    const tx = await metaNodeStake.withdraw(0);
    await tx.wait();

    const stakeAddress = (await metaNodeStake.getAddress)
      ? await metaNodeStake.getAddress()
      : metaNodeStake.address;
    const contractBalance = await ethers.provider.getBalance(stakeAddress);
    const userBalance = await ethers.provider.getBalance(owner.address);

    expect(contractBalance).to.be.equal(ethers.parseEther(("0.7")));
    // 因为用户支付了 gas，这里用 大于等于 不严谨；为稳妥检查差值正向
    expect(userBalance).to.be.gt((preUserBalance)-ethers.parseEther("0.01"));
  });

  it("Claim after deposit: 存入后领取奖励", async function () {
    await addNativePool();
    await metaNodeStake.depositETH({ value: ethers.parseEther("1") });

    // 推进区块高度
    await mineToBlock(1000);

    // 给质押合约充足的 MetaNode 代币
  const stakeAddress = metaNodeStake.getAddress
    ? await metaNodeStake.getAddress()
    : metaNodeStake.address;

    await metaNode.transfer(stakeAddress, ethers.parseEther("1000000"));

    const preUserBalance = await ethers.provider.getBalance(owner.address);
    const tx = await metaNodeStake.claimRewards(0);
    await tx.wait();
    const userBalance = await ethers.provider.getBalance(owner.address);

    const balance = await metaNode.balanceOf(owner.address);
    console.log("MetaNode balance after claim:", ethers.formatEther(balance));
    expect(balance).to.be.gt(ethers.parseEther("0")); // 领取奖励后，用户余额应增加
  });

  it("Claim after withdraw: 提取后领取奖励", async function () {
    await addNativePool();
    await metaNodeStake.depositETH({ value: ethers.parseEther("1") });
    await metaNodeStake.unstake(0, ethers.parseEther("1"));

    // 给质押合约充足的 MetaNode 代币
    const stakeAddress = (await metaNodeStake.getAddress)
      ? metaNodeStake.getAddress()
      : metaNodeStake.address;
    await metaNode.transfer(stakeAddress, ethers.parseEther("1000"));

    // 推进区块高度
    await mineToBlock(1000);

    // 提取原生代币
    const tx = await metaNodeStake.withdraw(0);
    await tx.wait();
    // 检查用户余额
    const preUserBalance = await ethers.provider.getBalance(owner.address);

    // 推进区块高度
    await mineToBlock(1000);

    // 领取奖励
    const claimTx = await metaNodeStake.claimRewards(0);
    await claimTx.wait();
    // 检查用户余额
    const userBalance = await ethers.provider.getBalance(owner.address);

    // expect(userBalance).to.be.gt(preUserBalance); // 领取奖励后，用户余额应增加

    const balance = await metaNode.balanceOf(owner.address);
    console.log("MetaNode balance after claim:", ethers.formatEther(balance));
    expect(balance).to.be.gt(ethers.parseEther("0")); // 领取奖励后，用户余额应增加
  });
});
