package main

import (
	"fmt"

	dig "github.com/voodoologic/domania/cmd"

	"github.com/spf13/cobra"
)

var digCmd = &cobra.Command{
	Use:   "dig",
	Short: "Executes dig command on a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dig.DigDomain(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Println(result)
	},
}

var simpleList = &cobra.Command{
	Use:   "dig",
	Short: "Executes dig command on a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dig.DigDomain(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Println(result)
	},
}

func main() {
	// if err := digCmd.Execute(); err != nil {
	// 	panic(err)
	// }
	x := dig.StartProgram()
	x = nil
	dig.Testing()
}
