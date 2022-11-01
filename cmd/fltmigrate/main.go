package main

import (
	"flutechain/database"
	"fmt"
	"os"
	"time"
)

func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer state.Close()

	state.Add(database.NewTx(database.Account("dan"), database.Account("andre"), 400, ""))
	state.Add(database.NewTx(database.Account("dan"), database.Account("dan"), 1500, "reward"))

	hash, err := state.Persist()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	block0 := database.NewBlock(
		hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx(database.Account("dan"), database.Account("andre"), 400, ""),
			database.NewTx(database.Account("dan"), database.Account("dan"), 1500, "reward"),
		},
	)

	err = state.AddBlock(block0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	blockHash0 := state.LatestHash()

	block1 := database.NewBlock(
		blockHash0,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx(database.Account("andre"), database.Account("dan"), 100, ""),
			database.NewTx(database.Account("dan"), database.Account("sam"), 500, ""),
			database.NewTx(database.Account("dan"), database.Account("dan"), 1000, "reward"),
		},
	)

	err = state.AddBlock(block1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
