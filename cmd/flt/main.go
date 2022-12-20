package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const accountFlag = "account"

const fromFlag = "from"
const toFlag = "to"
const valueFlag = "value"
const flagDataDir = "datadir"

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
	cmd.Flags().String(flagDataDir, "", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired(flagDataDir)
}
