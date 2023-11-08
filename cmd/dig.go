package dig

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	err := checkDependencies()
	if err != nil {
		log.Fatal(err)
	}
}

func checkDependencies() error {
	if err := checkBinary("dig"); err != nil {
		return err
	}
	if err := checkBinary("jc"); err != nil {
		return err
	}
	return nil
}

func checkBinary(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("Binary %s is not installed", name)
	}
	return nil
}

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
func DigDomain(domain string) (*[]DNSLookupResult, error) {
	cmd := fmt.Sprintf("dig +answer %s | jc --dig", domain)
	var result []DNSLookupResult
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal")
	}
	return &result, nil
}
