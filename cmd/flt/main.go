package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const accountFlag = "account"

const fromFlag = "from"
const toFlag = "to"
const valueFlag = "value"

func main() {
	var cmd = &cobra.Command{
		Use:   "flt",
		Short: "FluteChain CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.AddCommand(versionCmd)
	cmd.AddCommand(balancesCmd())
	cmd.AddCommand(txCmd())

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(accountFlag, "", "account to get balance")
	cmd.MarkFlagRequired(accountFlag)
}
