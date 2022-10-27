package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Represents a fungible set of tokens.
type Partition struct {
	Amount int
	// Partition Address
	Partition string
	// 여기 파티션이 주소 일단 스트링으로 적음.
}

// partition Token
type PartitionToken struct {
	Name      string    `json:"name"`
	MSPID     string    `json:"mspid"`
	Locked    bool      `json:"locked"`
	Partition Partition `json:"partition"`
}

// ClientAccountBalance returns the balance of the requesting client's account
func (s *SmartContract) BalanceOfByPartition(ctx contractapi.TransactionContextInterface, _partition string, _tokenHolder string) (int, error) {

	// Create allowanceKey
	clientKey, err := ctx.GetStub().CreateCompositeKey(clientPartitionPrefix, []string{_partition, _tokenHolder})
	if err != nil {
		return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", clientPartitionPrefix, err)
	}

	balanceByPartitionBytes, err := ctx.GetStub().GetState(clientKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}

	if balanceByPartitionBytes == nil {
		return 0, fmt.Errorf("the account %s does not exist", _tokenHolder)
	}

	// // 이런 조건절 함수를 패브릭에서는 어디에서 처리를 해주어야 할까?
	// if _validPartition(result._partition, result.owner) {
	// 	return partitions[owner][partitionToIndex[owner][_partition]-1].amount
	// } else {
	// 	return 0
	// }

	balance, _ := strconv.Atoi(string(balanceByPartitionBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	return balance, nil
}

// ClientAccountBalance returns the balance of the requesting client's account
func (s *SmartContract) ClientAccountBalanceByPartition(ctx contractapi.TransactionContextInterface, _partition string) (int, error) {

	id, err := _msgSender(ctx)
	if err != nil {
		return 0, err
	}

	// owner Address
	owner := getAddress([]byte(id))

	// Create allowanceKey
	clientPartitionKey, err := ctx.GetStub().CreateCompositeKey(clientPartitionPrefix, []string{_partition, owner})
	if err != nil {
		return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", clientPartitionPrefix, err)
	}

	tokenBytes, err := ctx.GetStub().GetState(clientPartitionKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		return 0, fmt.Errorf("the account %s does not exist", owner)
	}

	// if _validPartition(result._partition, result.owner) {
	// 	return partitions[owner][partitionToIndex[owner][_partition]-1].amount
	// } else {
	// 	return 0
	// }

	token := new(PartitionToken)
	err = json.Unmarshal(tokenBytes, token)
	if err != nil {
		return 0, fmt.Errorf("failed to obtain JSON decoding: %v", err)
	}

	return PartitionToken.Amount, nil
}
