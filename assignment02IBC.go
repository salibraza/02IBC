package assignment02IBC

import (
	"crypto/sha256"
	"fmt"
)

var p = fmt.Println
var p2 = fmt.Printf

const miningReward = 100
const rootUser = "Satoshi"

type Block struct {
	Spender     map[string]int
	skeys       []string
	Receiver    map[string]int
	rkeys       []string
	prevPointer *Block
	prevHash    string
	currentHash string
}

func newBlock() *Block {
	b := new(Block)
	b.prevPointer = nil
	b.prevHash = ""
	b.currentHash = ""
	b.Spender = make(map[string]int)
	b.Receiver = make(map[string]int)
	return b
}

//
func CalculateHash(inputBlock *Block) string {
	//fmt.Println("Calculate Hash")
	// value to hash = prev hash + string(prev pointer) + (spender + amt) + (receiver + amt)
	if inputBlock == nil {
		p("Blockchain Is Empty")
		return ""
	}

	var input string = inputBlock.prevHash + fmt.Sprintf("%x", inputBlock.prevPointer)

	for _, spender := range inputBlock.skeys {
		input += spender + fmt.Sprint(inputBlock.Spender[spender])
	}
	for _, receiver := range inputBlock.rkeys {
		input += receiver + fmt.Sprint(inputBlock.Receiver[receiver])
	}

	sum := sha256.Sum256([]byte(input))
	hash := fmt.Sprintf("%x", sum) // converting from type [32]uint8 to string

	//p("Hash Input: ", input)
	//p("Hash Output", hash)
	return hash
}

func InsertBlock(spendingUser string, receivingUser string, miner string, amount int, chainHead *Block) *Block {
	p("\nInserting")
	// New Block
	var i *Block = newBlock()
	if chainHead == nil {
		chainHead = newBlock()
	}
	// Coin Minting Case
	if spendingUser == "" && receivingUser == "" && miner == rootUser && amount == 0 {
		i.Receiver[miner] = miningReward
		i.rkeys = append(i.rkeys, miner)

		i.prevPointer = chainHead.prevPointer
		i.prevHash = chainHead.prevHash
		i.currentHash = CalculateHash(i)
		chainHead.prevPointer = i
		chainHead.prevHash = i.currentHash
		p("\nInserted")
		return chainHead
	}

	balance := CalculateBalance(spendingUser, chainHead)
	// Validity Checks
	if balance < amount {
		p("ERROR_S: ", spendingUser, " doesn't have enought balance")
		return chainHead
	}
	if miner != rootUser {
		p("ERROR_M: ", miner, " is not a valid miner")
		return chainHead
	}

	// inserting into map
	i.Spender[spendingUser] = amount
	i.skeys = append(i.skeys, spendingUser)
	i.Receiver[receivingUser] = amount
	i.Receiver[miner] = miningReward
	i.rkeys = append(i.rkeys, receivingUser, miner)

	i.prevPointer = chainHead.prevPointer
	i.prevHash = chainHead.prevHash
	i.currentHash = CalculateHash(i)
	chainHead.prevPointer = i
	chainHead.prevHash = i.currentHash

	p("\nInserted")
	return chainHead
}

//Completed
func CalculateBalance(userName string, chainHead *Block) int {
	var balance int = 0
	var current *Block = chainHead
	if current == nil {
		return balance
	}
	for current.prevPointer != nil {
		for spender := range current.prevPointer.Spender {
			if spender == userName {
				balance -= current.prevPointer.Spender[spender]
			}
		}
		for receiver := range current.prevPointer.Receiver {
			if receiver == userName {
				balance += current.prevPointer.Receiver[receiver]
			}
		}
		current = current.prevPointer
	}
	//p("Balance: ", balance)
	return balance
}

//Completed
func ListBlocks(chainHead *Block) {
	if chainHead == nil {
		p("Blockchain Is Empty")
		return
	}
	var current *Block = chainHead.prevPointer
	p("\nBLOCKS' TRANSACTION LIST")
	for current != nil {
		p("")
		for _, spender := range current.skeys {
			p(spender, " spent ", current.Spender[spender])
		}
		for _, receiver := range current.rkeys {
			p(receiver, " received ", current.Receiver[receiver])
		}
		current = current.prevPointer
	}
	p("")
}

func VerifyChain(chainHead *Block) {
	if chainHead == nil {
		p("Blockchain Is Empty")
		return
	}
	var current *Block = chainHead
	for current.prevPointer != nil {
		if current.prevHash != CalculateHash(current.prevPointer) {
			p("ERROR_V: Blockchain has been changed")
			return
		}
		current = current.prevPointer
	}
	p("Blockchain all Good :)")
}
