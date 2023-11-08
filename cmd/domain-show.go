package dig

import (
	"fmt"
)

func DomainDetails(host string) {
	// I want to call dig on this and display the information in a table
	// I want to call namecheap and show the same information in a table

	// namecheap.DomainsDNSGetHostsCommandResponse()
	// client.Domains.GetInfo(name)
	lookupResult, _ := DigDomain(host)
	fmt.Println(host)
	for _, result := range *lookupResult {

		fmt.Printf("question: %s\n", result.Question.Type)
		for _, answer := range result.Answer {
			fmt.Printf("%v", answer)
			fmt.Printf("Type: %s | Name: %s | Data: %s\n", answer.Type, answer.Name, answer.Data)
		}

	}
}
