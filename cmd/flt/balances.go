package main

import (
	"flutechain/database"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var balancesListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all wallet balances.",
	Run: func(cmd *cobra.Command, args []string) {

		dataDir, _ := cmd.Flags().GetString(flagDataDir)
		state, err := database.NewStateFromDisk(dataDir)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer state.Close()

		fmt.Println("Account Balances:")
		fmt.Println("__________________")
		fmt.Println("")

		for account, balance := range state.Balances {
			fmt.Println(fmt.Sprintf("%s: %d", account, balance))
		}
	},
}

func balancesCmd() *cobra.Command {
	var balanceCmd = &cobra.Command{
		Use:   "balances",
		Short: "interact with balances",
		Run: func(cmd *cobra.Command, args []string) {

			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			defer state.Close()

			account, _ := cmd.Flags().GetString(accountFlag)

			if account == "" {
				fmt.Errorf("incomplete request")
				os.Exit(1)
			}

			fmt.Println(fmt.Sprintf("%s: %d", account, state.Balances[database.Account(account)]))
		},
	}

	balanceCmd.Flags().String(accountFlag, "", "account to get balance")
	balanceCmd.MarkFlagRequired(accountFlag)

	addDefaultRequiredFlags(balanceCmd)

	balanceCmd.AddCommand(balancesListCommand)

	return balanceCmd
}
