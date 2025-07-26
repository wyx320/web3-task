package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到 Sepolia 测试网络
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/e48f7acfd92045759d9243f878973e8a")
	if err != nil {
		log.Fatal("Failed to connect to the Ethereum client:", err)
	}

	// var chainID *big.Int
	// // 获取链 ID
	// chainID, err = client.ChainID(context.Background())
	// if err != nil {
	// 	log.Fatal("Failed to get chain ID:", err)
	// }
	// fmt.Printf("Connected to chain ID: %d", chainID)

	fmt.Println("===================== Block Information ====================")

	blockNumber := big.NewInt(5671744)
	blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")

	var block *types.Block
	block, err = client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal("Failed to retrieve block:", err)
	}
	fmt.Println("Block Number:", block.Number().Uint64())
	fmt.Println("Block Hash:", block.Hash().Hex())
	fmt.Println("Block transactionCount:", len(block.Transactions()))

	fmt.Println("Block Time:", block.Time())
	fmt.Println("Block Difficulty:", block.Difficulty().Uint64())
	fmt.Println("Block BaseFee:", block.BaseFee().Uint64())
	fmt.Println("Block GasLimit:", block.GasLimit())

	block, err = client.BlockByHash(context.Background(), blockHash)
	if err != nil {
		log.Fatal("Failed to retrieve block by hash:", err)
	}
	fmt.Println("Block Number by Hash:", block.Number().Uint64())

	fmt.Println("===================== Account Information ====================")

	privateKey, err := crypto.HexToECDSA("92c4ef37898a9f64b32e1fd5868a946c81202d15788113175d7fcbb43aacbd8c")
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to cast public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("Address:", address.Hex())

	fmt.Println("====================== Transaction Information ====================")

	// 接收方地址
	toAddress := common.HexToAddress("0x3d067293acC10F8bA4687e2A3270945c378F28A5")
	// 转账金额
	amount := big.NewInt(100000000000000) // 0.0001 ETH
	// 随机数
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}
	// 燃料价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Failed to suggest gas price:", err)
	}
	// // 燃料小费
	// gacTipPrice, err := client.SuggestGasTipCap(context.Background())
	// if err != nil {
	// 	log.Fatal("Failed to suggest gas tip price:", err)
	// }
	// 燃料上限
	gasLimmit := uint64(21000)
	// 使用私钥对交易进行签名
	trans := types.NewTransaction(nonce, toAddress, amount, gasLimmit, gasPrice, nil)
	signedTrans, err := types.SignTx(trans, types.NewEIP155Signer(block.Number()), privateKey)
	if err != nil {
		log.Fatal("Failed to sign transaction:", err)
	}
	err = client.SendTransaction(context.Background(), signedTrans)
	if err != nil {
		log.Fatal("Failed to send transaction:", err)
	}

	fmt.Println("trans sent: ", signedTrans.Hash().Hex())
}
