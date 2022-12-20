package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type State struct {
	Balances        map[Account]uint
	TxMempool       []Tx
	latestBlockHash Hash
	dbFile          *os.File
}

func NewStateFromDisk(dataDir string) (*State, error) {

	err := initDataDirIfNotExists(dataDir)
	if err != nil {
		return nil, err
	}

	genesis, err := LoadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		return nil, err
	}

	// // to get working dir
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }

	// // get genesis file path
	// genesisFilePath := filepath.Join(cwd, "database", "genesis.json")
	// genesis, err := LoadGenesis(genesisFilePath)
	// if err != nil {
	// 	return nil, err
	// }

	balances := make(map[Account]uint)

	for account, balance := range genesis.Balances {
		balances[account] = balance
	}

	f, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{
		Balances:        balances,
		TxMempool:       make([]Tx, 0),
		dbFile:          f,
		latestBlockHash: Hash{},
	}

	// iterate over each of the db rows
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var blockFs BlockFS
		json.Unmarshal(scanner.Bytes(), &blockFs)

		if err := state.ApplyBlock(blockFs.Value); err != nil {
			return nil, err
		}
	}

	//err = state.takeSnapshot()
	//if err != nil {
	//	return nil, err
	//}

	return state, nil
}

func (s *State) ApplyBlock(b Block) error {
	for _, t := range b.Tx {
		err := s.ApplyTx(t)
		if err != nil {
			fmt.Println("Error Applying Tx")
			return err
		}
	}
	return nil
}

func (s *State) LatestHash() Hash {
	return s.latestBlockHash
}

func (s *State) ApplyTx(t Tx) error {

	if t.isReward() {
		s.Balances[t.To] += t.Value
		return nil
	}

	if s.Balances[t.From] < t.Value {
		return fmt.Errorf("insufficent balance")
	}

	s.Balances[t.From] -= t.Value
	s.Balances[t.To] += t.Value

	return nil
}

func (s *State) Add(tx Tx) error {
	s.TxMempool = append(s.TxMempool, tx)

	return nil
}

func (s *State) AddBlock(block Block) error {

	if err := s.ApplyBlock(block); err != nil {
		return err
	}

	blockHash, err := block.Hash()
	if err != nil {
		return err
	}

	bf := BlockFS{blockHash, block}

	blockFsJson, err := json.Marshal(bf)
	if err != nil {
		return err
	}

	//fmt.Printf("persisting new Block to disk:\n")
	fmt.Printf("%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return err
	}

	s.latestBlockHash = blockHash

	return nil
}

func (s *State) Persist() (Hash, error) {
	block := NewBlock(s.LatestHash(), uint64(time.Now().Unix()), s.TxMempool)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{blockHash, block}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	//fmt.Printf("persisting new Block to disk:\n")
	fmt.Printf("%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash

	s.TxMempool = []Tx{}

	return s.latestBlockHash, nil
}

func (s *State) takeSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}

	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}

	s.latestBlockHash = sha256.Sum256(txsData)

	return nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func NewAccount(value string) Account {
	return Account(value)
}
