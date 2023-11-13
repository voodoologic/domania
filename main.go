package main

import (
	dig "github.com/voodoologic/domania/cmd"
)

// var digCmd = &cobra.Command{
// 	Use:   "dig",
// 	Short: "Executes dig command on a domain",
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		result, err := dig.DigDomain(args[0])
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(result)
// 	},
// }

// var simpleList = &cobra.Command{
// 	Use:   "dig",
// 	Short: "Executes dig command on a domain",
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		result, err := dig.DigDomain(args[0])
// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Println(result)
// 	},
// }

func main() {
	// if err := digCmd.Execute(); err != nil {
	// 	panic(err)
	// }
	// DomainList := dig.StartProgram()
	// selectedDomain := dig.ListDomains(DomainList)
	// fmt.Println(selectedDomain)
	// if *selectedDomain == "" {
	// 	panic("blank")
	// }
	// dig.DomainDetails(selectedDomain)
	dig.DomainDetails("dougheadley.com")
}
