package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"my-ether/counter"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/e48f7acfd92045759d9243f878973e8a")
	if err != nil {
		log.Fatal("Failed to connect to the Ethereum client: ", err)
	}

	// 如果地址文件不存在，则开始部署合约
	if !addrFileExists("address.txt") {

		// 开始打印部署合约信息
		fmt.Println("================ Contract Deploy ================")

		// 获取私钥
		privateKey, err := crypto.HexToECDSA("92c4ef37898a9f64b32e1fd5868a946c81202d15788113175d7fcbb43aacbd8c")
		if err != nil {
			log.Fatal("Failed to parse private key: ", err)
		}
		// 获取公钥
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("Failed to cast public key to ECDSA")
		}
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		// 获取随机数
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal("Failed to get nonce: ", err)
		}
		// 获取链 ID
		chainID, err := client.ChainID(context.Background())
		if err != nil {
			log.Fatal("Failed to get chain ID: ", err)
		}
		// 获取 gas价格
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal("Failed to suggest gas price: ", err)
		}
		// 签名
		transOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			log.Fatal("Failed to create transactor: ", err)
		}
		// 参数填充
		transOpts.Nonce = big.NewInt(int64(nonce)) // 设置 nonce
		transOpts.Value = big.NewInt(0)            // 不发送以太币
		transOpts.GasPrice = gasPrice              // 设置 gas 价格
		transOpts.GasLimit = 3000000               // 设置 gas 限制
		// 部署合约
		address, trans, instance, err := counter.DeployCounter(transOpts, client)
		if err != nil {
			log.Fatal("Failed to deploy contract: ", err)
		}
		fmt.Println("Contract Address:", address.Hex())
		fmt.Println("Transaction Address:", trans.Hash().Hex())
		// 保存合约地址

		_ = instance // 这里可以使用 instance 调用合约方法
	}
}

func addrFileExists(fileName string) bool {
	dir := "./.cache"
	// 检查地址文件夹是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false // 文件夹不存在
	}
	// 检查地址文件是否存在
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false // 文件不存在
	}
	return true
}

func saveAddressToFile(address string, fileName string) error {
dir := "./.cache"
if _,err:=os.Stat(dir); os.IsNotExist(err) {
	// 创建目录
	if err := os.MkdirAll(dir, os.ModeAppend)
}

	// file, err := os.Create(fileName)
	// if err != nil {
	// 	return fmt.Errorf("failed to create file: %w", err)
	// }
	// defer file.Close()

	// _, err = file.WriteString(address)
	// if err != nil {
	// 	return fmt.Errorf("failed to write address to file: %w", err)
	// }

	// return nil
}
