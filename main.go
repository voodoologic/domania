package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

type DNSLookupResult struct {
	ID               int           `json:"id"`
	Opcode           string        `json:"opcode"`
	Status           string        `json:"status"`
	Flags            []string      `json:"flags"`
	QueryNum         int           `json:"query_num"`
	AnswerNum        int           `json:"answer_num"`
	AuthorityNum     int           `json:"authority_num"`
	AdditionalNum    int           `json:"additional_num"`
	OptPseudosection Pseudosection `json:"opt_pseudosection"`
	Question         DNSEntity     `json:"question"`
	Answer           []DNSEntity   `json:"answer"`
	QueryTime        int           `json:"query_time"`
	Server           string        `json:"server"`
	When             string        `json:"when"`
	Rcvd             int           `json:"rcvd"`
	WhenEpoch        int64         `json:"when_epoch"`
	// Don't forget to add other fields here...
}

type Pseudosection struct {
	EDNS EDNSInfo `json:"edns"`
	// Don't forget to add other fields here...
}

type EDNSInfo struct {
	Version int      `json:"version"`
	Flags   []string `json:"flags"`
	UDP     int      `json:"udp"`
}

type DNSEntity struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Type  string `json:"type"`
	TTL   int    `json:"ttl,omitempty"`
	Data  string `json:"data,omitempty"`
}

// digDomain function to run dig command on a domain and returns DigResult struct
func digDomain(domain string) (*[]DNSLookupResult, error) {

	cmd := fmt.Sprintf("dig %s | jc --dig", domain)
	var result []DNSLookupResult
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal")
	}
	fmt.Printf("the result is: %v", output)
	return &result, nil
}

var digCmd = &cobra.Command{
	Use:   "dig",
	Short: "Executes dig command on a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := digDomain(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Println(result)
	},
}

func main() {
	fmt.Println("Domania")
	if err := digCmd.Execute(); err != nil {
		panic(err)

	}

}
