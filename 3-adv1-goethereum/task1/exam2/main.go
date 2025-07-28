package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	counterDeploy "my-ether/counter-deploy"
	counterLoad "my-ether/counter-load"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/<-- Your Infura Project ID Here -->")
	if err != nil {
		log.Fatal("Failed to connect to the Ethereum client: ", err)
	}

	// 获取私钥
	privateKey, err := crypto.HexToECDSA("<-- Your PrivateKeyHere -->")
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	// 获取 gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Failed to suggest gas price: ", err)
	}
	// 获取链 ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("Failed to get chain ID: ", err)
	}

	addrExists := addrFileExists("address.txt")
	// 如果地址文件不存在，则开始部署合约
	if !addrExists {

		// 开始打印部署合约信息
		fmt.Println("================ Contract Deploy ================")

		// // 获取公钥
		// publicKey := privateKey.Public()
		// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		// if !ok {
		// 	log.Fatal("Failed to cast public key to ECDSA")
		// }
		// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		// // 获取随机数
		// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		// if err != nil {
		// 	log.Fatal("Failed to get nonce: ", err)
		// }
		// fmt.Println("Nonce:", nonce)
		// 签名
		transOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			log.Fatal("Failed to create transactor: ", err)
		}
		// 参数填充
		// nonce = 1
		// transOpts.Nonce = big.NewInt(int64(nonce))                      // 设置 nonce
		transOpts.Value = big.NewInt(0)                                 // 不发送以太币
		transOpts.GasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(2)) // 设置 gas 价格
		transOpts.GasLimit = 3000000                                    // 设置 gas 限制
		// 部署合约
		address, trans, instance, err := counterDeploy.DeployCounter(transOpts, client)
		if err != nil {
			log.Fatal("Failed to deploy contract: ", err)
		}

		fmt.Println("Contract Address:", address.Hex())
		fmt.Println("Transaction Address:", trans.Hash().Hex())

		// 等待交易被挖矿
		receipt, err := bind.WaitMined(context.Background(), client, trans)
		if err != nil {
			log.Fatal("Transaction not mined: ", err)
		}
		fmt.Println("Transaction mined in block:", receipt.BlockNumber.Uint64())
		fmt.Println("Gas used:", receipt.GasUsed)

		// 保存合约地址
		saveAddressToFile(address.Hex(), "address.txt")

		_ = instance // 这里可以使用 instance 调用合约方法
	}

	// 加载合约
	fmt.Println("================ Contract Load ================")

	addr, err := getAddrFromFile("address.txt")
	if err != nil {
		log.Fatal("Failed to read address file: ", err)
	}
	instance, err := counterLoad.NewCounter(common.HexToAddress(addr), client)
	if err != nil {
		log.Fatal("Failed to load contract instance: ", err)
	}
	//
	transOpt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("Failed to create transactor: ", err)
	}
	// // 获取随机数
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal("Failed to get nonce: ", err)
	// }
	// transOpt.GasLimit = 3000000  // 设置 gas 限制
	// transOpt.GasPrice = gasPrice // 设置 gas 价格
	// transOpt.Value = big.NewInt(0) // 不发送以太币
	// transOpt.data

	// 调用合约函数
	fmt.Println("================ Contract Test ================")

	callOpts := &bind.CallOpts{Context: context.Background()}

	// getCount
	count, err := instance.GetCount(callOpts)
	if err != nil {
		log.Fatal("Failed to call GetCount: ", err)
	}
	fmt.Println("Current Count:", count)

	// count++
	trans, err := instance.Increment(transOpt)
	if err != nil {
		log.Fatal("Failed to call Increment: ", err)
	}
	fmt.Println("Wait For Completed. Transaction Hex:", trans.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), client, trans)
	if err != nil {
		log.Fatal("Transaction not mined: ", err)
	}
	// 检查交易状态
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Fatal("Transaction failed")
	}
	fmt.Println("Incremented successfully, Incremented by 1")

	// count++
	trans, err = instance.Increment(transOpt)
	if err != nil {
		log.Fatal("Failed to call Increment: ", err)
	}
	fmt.Println("Wait For Completed. Transaction Hex:", trans.Hash().Hex())
	receipt, err = bind.WaitMined(context.Background(), client, trans)
	if err != nil {
		log.Fatal("Transaction not mined: ", err)
	}
	// 检查交易状态
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Fatal("Transaction failed")
	}
	fmt.Println("Incremented successfully, Incremented by 1")

	// getCount
	count, err = instance.GetCount(callOpts)
	if err != nil {
		log.Fatal("Failed to call GetCount: ", err)
	}
	fmt.Println("Current Count:", count)

	// count--
	trans, err = instance.Decrement(transOpt)
	if err != nil {
		log.Fatal("Failed to call Increment: ", err)
	}
	fmt.Println("Wait For Completed. Transaction Hex:", trans.Hash().Hex())
	receipt, err = bind.WaitMined(context.Background(), client, trans)
	if err != nil {
		log.Fatal("Transaction not mined: ", err)
	}
	// 检查交易状态
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Fatal("Transaction failed")
	}
	fmt.Println("Decremented successfully, Decremented by 1")

	// getCount
	count, err = instance.GetCount(callOpts)
	if err != nil {
		log.Fatal("Failed to call GetCount: ", err)
	}
	fmt.Println("Current Count:", count)

	// instance.Increment()
}

// addrFileExists checks if the address file exists in the specified directory.
func addrFileExists(fileName string) bool {
	dir := "./.cache"
	fullPath := filepath.Join(dir, fileName)
	// 检查地址文件夹是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false // 文件夹不存在
	}
	// 检查地址文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false // 文件不存在
	}
	return true
}

// saveAddressToFile saves the contract address to a file in the specified directory.
func saveAddressToFile(address string, fileName string) error {
	dir := "./.cache"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Sprintf("failed to create directory: %v", err))
		}
	}
	fullPath := filepath.Join(dir, fileName)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// 创建文件 写入内容
		err := os.WriteFile(fullPath, []byte(address), 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to create file: %v", err))
		}
	}
	// 如果文件已存在，直接返回
	return nil
}

// getAddrFromFile reads the contract address from a file.
// It returns the address as a string or an error if the file cannot be read.
func getAddrFromFile(fileName string) (string, error) {
	dir := "./.cache"
	fullPath := filepath.Join(dir, fileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(data), nil
}
