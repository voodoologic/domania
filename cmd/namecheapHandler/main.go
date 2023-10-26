package namecheapHandler

import (
	"fmt"
	"os"

	namecheap "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

// Client holds Namecheap API client instance
type Client struct {
	ncClient *namecheap.Client
}

// NewClient returns a new client to communicate with Namecheap API
func NewClient() (*Client, error) {
	options := &namecheap.ClientOptions{
		UserName:   os.Getenv("USER_NAME"),
		ApiUser:    os.Getenv("USER_NAME"),
		ApiKey:     os.Getenv("API_KEY"),
		ClientIp:   os.Getenv("CLIENT_IP"),
		UseSandbox: os.Getenv("SANDBOX") == "true",
	}
	ncClient := namecheap.NewClient(options)
	return &Client{ncClient: ncClient}, nil
}

// GetDomainList fetches and returns all domains
func (cl *Client) GetDomainList() ([]string, error) {
	domainListArgs := &namecheap.DomainsGetListArgs{}
	domainList, err := cl.ncClient.Domains.GetList(domainListArgs)
	if err != nil {
		return nil, err
	}
	fmt.Println(domainList)

	var domainNames []string
	for _, domain := range *domainList.Domains {
		domainNames = append(domainNames, *domain.Name)
	}

	return domainNames, nil
}
