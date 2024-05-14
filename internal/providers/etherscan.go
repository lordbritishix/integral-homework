package providers

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lordbritishix/integral/internal/config"
	"github.com/lordbritishix/integral/internal/dto/etherscan"
	"strings"
)

type EtherscanClient struct {
	httpClient *resty.Client
	config     *config.Config
}

func NewEtherscanClient(config *config.Config) *EtherscanClient {
	return &EtherscanClient{
		httpClient: resty.New().SetBaseURL("https://api.etherscan.io/api"),
		config:     config,
	}
}

func (e *EtherscanClient) GetEthTransactions(address string, startBlock int, endBlock int, offset int) ([]etherscan.Transaction, error) {
	urlFragment := fmt.Sprintf("?module=account&action=txlist&address=%s&startBlock=%d&endBlock=%d&offset=%d&sort=desc&apikey=%s",
		address, startBlock, endBlock, offset, e.config.EtherscanApiKey)

	resp, err := e.httpClient.R().Get(urlFragment)
	if err != nil {
		return nil, err
	}

	var data etherscan.Response
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}

	var result []etherscan.Transaction

	if data.Status == "1" {
		transactions := data.Result
		for _, tx := range transactions {
			// Check if the transaction is an Ethereum transfer
			if strings.TrimSpace(tx.Input) == "0x" {
				result = append(result, tx)
			}
		}
	} else if data.Status == "0" {
		return result, nil
	}

	return result, nil
}

func (e *EtherscanClient) GetTokenTransactions(address string, startBlock int, endBlock int, offset int) ([]etherscan.Transaction, error) {
	urlFragment := fmt.Sprintf("?module=account&action=tokentx&address=%s&startBlock=%d&endBlock=%d&offset=%d&sort=desc&apikey=%s",
		address, startBlock, endBlock, offset, e.config.EtherscanApiKey)

	resp, err := e.httpClient.R().Get(urlFragment)
	if err != nil {
		return nil, err
	}

	var data etherscan.Response
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}

	var result []etherscan.Transaction

	if data.Status == "1" {
		transactions := data.Result
		for _, tx := range transactions {
			result = append(result, tx)
		}
	} else {
		return result, nil
	}

	return result, nil
}

func (e *EtherscanClient) IsTokenSpam(contractAddress string) (bool, error) {
	// if source code is verified by etherscan, there's a high chance that it is not a spam
	urlFragment := fmt.Sprintf("?module=contract&action=getabi&address=%s&apikey=%s", contractAddress, e.config.EtherscanApiKey)
	resp, err := e.httpClient.R().Get(urlFragment)
	if err != nil {
		return false, err
	}

	var data etherscan.ResponseContract
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return false, err
	}

	if data.Result == "Contract source code not verified" {
		return true, nil
	}

	return false, nil
}
