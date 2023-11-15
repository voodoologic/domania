package dig

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/table"
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

func InitDomain(domain string) ([]table.Row, error) {
	rows := []table.Row{}
	for _, recordType := range []string{"A", "MX", "TXT", "CNAME", "NS", "SOA"} {
		digRows, err := DigDomain(domain, recordType)
		if err != nil {

		}
		rows = append(rows, digRows...)
	}
	return rows, nil
}

// digDomain function to run dig command on a domain and returns DigResult struct
func DigDomain(domain string, recordType string) ([]table.Row, error) {
	// ips, err := net.LookupIP("google.com")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
	// 	os.Exit(1)
	// }
	// for _, ip := range ips {
	// 	fmt.Printf("google.com. IN A %s\n", ip.String())
	// }
	cmd := fmt.Sprintf("dig +answer %s %s | jc --dig", domain, recordType)
	var result []DNSLookupResult
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal")
	}
	results := processResults(result)
	return results, nil
}

func processResults(results []DNSLookupResult) []table.Row {
	rows := []table.Row{}
	for _, digReport := range results {
		for _, answer := range digReport.Answer {
			switch answer.Type {
			case "MX":
				re := regexp.MustCompile(`\d+ `)
				cleanMXData := re.ReplaceAllString(answer.Data, "")
				rows = append(rows, table.Row{"?", answer.Type, answer.Name, cleanMXData})
			case "A":
				//noop
			case "SOA":
				data := strings.Fields(answer.Data)
				re := regexp.MustCompile(`^\D*$`) // matches only if string contains no digits

				var filteredEntries []string
				for _, entry := range data {
					if re.MatchString(entry) {
						filteredEntries = append(filteredEntries, entry)
					}
				}
				for _, datum := range filteredEntries {
					rows = append(rows, table.Row{"?", "A", answer.Name, datum})
				}
			default:
				rows = append(rows, table.Row{"?", answer.Type, answer.Name, answer.Data})
			}
		}
	}
	return rows
}
