package main

import (
	"flutechain/database"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func txCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (add, )",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	txsCmd.AddCommand(txAddCmd())

	return txsCmd
}

func txAddCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "Adds new Tx to database",
		Run: func(cmd *cobra.Command, args []string) {

			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			from, _ := cmd.Flags().GetString(fromFlag)
			to, _ := cmd.Flags().GetString(toFlag)
			value, _ := cmd.Flags().GetUint(valueFlag)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			tx := database.NewTx(fromAcc, toAcc, value, "")

			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			defer state.Close()

			err = state.Add(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			_, err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			fmt.Println("TX successfully added to ledger")
		},
	}

	addDefaultRequiredFlags(cmd)

	cmd.Flags().String(fromFlag, "", "Account to send tokens from")
	cmd.Flags().String(toFlag, "", "Account to send tokens to")
	cmd.Flags().Uint(valueFlag, 0, "Amount of tokens to send")

	return cmd
}
