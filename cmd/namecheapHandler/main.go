package namecheapHandler

import (
	"os"

	"github.com/charmbracelet/bubbles/table"
	namecheap "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

// Client holds Namecheap API client instance
// NewClient returns a new client to communicate with Namecheap API
var cacheClient *namecheap.Client

func NewClient() (*namecheap.Client, error) {

	options := &namecheap.ClientOptions{
		UserName:   os.Getenv("USER_NAME"),
		ApiUser:    os.Getenv("USER_NAME"),
		ApiKey:     os.Getenv("API_KEY"),
		ClientIp:   os.Getenv("CLIENT_IP"),
		UseSandbox: os.Getenv("SANDBOX") == "true",
	}
	if cacheClient == nil {
		cacheClient = namecheap.NewClient(options)
	}
	return cacheClient, nil
}

// GetDomainList fetches and returns all domains
func GetDomainList() ([]string, error) {
	cl, _ := NewClient()
	domainListArgs := &namecheap.DomainsGetListArgs{}
	domainList, err := cl.Domains.GetList(domainListArgs)
	if err != nil {
		return nil, err
	}
	var domainNames []string
	for _, domain := range *domainList.Domains {
		domainNames = append(domainNames, *domain.Name)
	}

	return domainNames, nil
}

func GetDomainDetails(hostname string) ([]table.Row, error) {
	rows := []table.Row{}
	rows, _ = getNameServers(hostname, rows)
	rows, _ = getHosts(hostname, rows)
	return rows, nil
}

func getHosts(hostname string, rows []table.Row) ([]table.Row, error) {
	cl, err := NewClient()
	if err != nil {
		return nil, err
	}
	myHostResult, err := cl.DomainsDNS.GetHosts(hostname)
	if err != nil {
		return nil, err
	}
	hosts := *myHostResult.DomainDNSGetHostsResult.Hosts
	var myhost string
	for _, host := range hosts {
		if *host.Name == "@" {
			myhost = hostname + "."
		} else {
			myhost = *host.Name + "." + hostname + "."
		}
		var myType string
		if *host.Type == "ALIAS" {
			myType = "A"
		} else {
			myType = *host.Type
		}
		r := table.Row{
			"?",
			myType,
			myhost,
			*host.Address,
		}
		rows = append(rows, r)
	}
	return rows, nil
}

func getNameServers(hostname string, rows []table.Row) ([]table.Row, error) {
	cl, _ := NewClient()
	getInfoResponse, err := cl.Domains.GetInfo(hostname)
	if err != nil {
		return nil, err
	}

	for _, derp := range *getInfoResponse.DomainDNSGetListResult.DnsDetails.Nameservers {
		r := table.Row{
			"?",
			"NS",
			hostname + ".",
			derp + ".",
		}
		rows = append(rows, r)

	}
	return rows, nil
}
