package assignment02IBC

import (
	"crypto/sha256"
	"fmt"
)

const miningReward = 100
const rootUser = "Satoshi"

type BlockData struct {
	Title    string
	Sender   string
	Receiver string
	Amount   int
}
type Block struct {
	Data        []BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

func CalculateBalance(userName string, chainHead *Block) int {
	mainHead := chainHead
	var totalAmount int = 0
	for mainHead != nil {
		transactions := mainHead.Data
		for i := 0; i < len(transactions); i++ {
			if transactions[i].Sender == userName {
				totalAmount -= transactions[i].Amount
			} else if transactions[i].Receiver == userName {
				totalAmount += transactions[i].Amount
			}
		}
		mainHead = mainHead.PrevPointer
	}
	return totalAmount
}
func CalculateHash(inputBlock *Block) string {
	data := fmt.Sprintf("%v", inputBlock.Data)
	calHash := fmt.Sprintf("%x\n", sha256.Sum256([]byte(data)))
	return calHash
}
func VerifyTransaction(transaction *BlockData, chainHead *Block) bool {
	balance := CalculateBalance(transaction.Sender, chainHead)
	if balance >= transaction.Amount {
		return true
	} else {
		fmt.Print("ERROR: ", transaction.Sender, " has ", balance, " coins - ", transaction.Amount, " were needed !\n")
		return false
	}
}
func InsertBlock(blockData []BlockData, chainHead *Block) *Block {
	miningTransaction := BlockData{
		Title:    "Mining Reward",
		Sender:   "System",
		Receiver: rootUser,
		Amount:   miningReward,
	}
	blockData = append(blockData, miningTransaction)
	for i := 0; i < len(blockData); i++ {
		if blockData[i].Sender != "System" {
			if !VerifyTransaction(&blockData[i], chainHead) {
				return chainHead
			}
		}
	}
	if chainHead == nil {
		temp := &Block{
			Data:        blockData,
			PrevPointer: nil,
			CurrentHash: "",
			PrevHash:    "",
		}
		currentHash := CalculateHash(temp)
		chainHead = &Block{
			Data:        blockData,
			PrevPointer: nil,
			CurrentHash: currentHash,
			PrevHash:    "",
		}

	} else {
		temp := &Block{
			Data:        blockData,
			PrevPointer: chainHead,
			CurrentHash: "",
			PrevHash:    chainHead.CurrentHash,
		}
		currentHash := CalculateHash(temp)
		chainHead = &Block{
			Data:        blockData,
			PrevPointer: chainHead,
			CurrentHash: currentHash,
			PrevHash:    chainHead.CurrentHash,
		}
	}
	for i := 0; i < len(chainHead.Data); i++ {
		if chainHead.Data[i].Sender != "System" {
			if CalculateBalance(chainHead.Data[i].Sender, chainHead) < 0 {
				return chainHead.PrevPointer
			}
		}
	}
	return chainHead
}
func ListBlocks(chainHead *Block) {
	tempHead := chainHead
	fmt.Printf("Head\n")
	for tempHead != nil {
		fmt.Printf("⬇\n")
		for i := 0; i < len(tempHead.Data); i++ {
			fmt.Print("Title:", tempHead.Data[i].Title, " Sender:", tempHead.Data[i].Sender, " Receiver:", tempHead.Data[i].Receiver, " Amount:", tempHead.Data[i].Amount, "\n")
		}
		tempHead = tempHead.PrevPointer
	}
	fmt.Print("\n")
}
func VerifyChain(chainHead *Block) {
	tempHead := chainHead
	for tempHead.PrevPointer != nil {
		VerHash := CalculateHash(tempHead)
		if tempHead.CurrentHash != VerHash {
			fmt.Println("Block Chain Crompromised")
			return
		}
		tempHead = tempHead.PrevPointer
	}
	VerHash := CalculateHash(tempHead)
	if tempHead.CurrentHash != VerHash {
		fmt.Println("Block Chain Crompromised")
		return
	} else {
		fmt.Println("Block Chian OK")
	}
}
func PremineChain(chainHead *Block, numBlocks int) *Block {
	// bolackData := [] BlockData {}
	for numBlocks != 0 {
		chainHead = InsertBlock([]BlockData{}, chainHead)
		numBlocks--
	}
	return chainHead
}
