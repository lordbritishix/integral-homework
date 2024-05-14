package etherscan

type Transaction struct {
	BlockNumber     string  `json:"blockNumber"`
	TimeStamp       string  `json:"timeStamp"`
	Hash            string  `json:"hash"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Value           string  `json:"value"`
	Gas             string  `json:"gas"`
	GasPrice        string  `json:"gasPrice"`
	Input           string  `json:"input"`
	TokenName       *string `json:"tokenName,omitempty"`
	TokenSymbol     *string `json:"tokenSymbol,omitempty"`
	TokenDecimal    *string `json:"tokenDecimal,omitempty"`
	ContractAddress *string `json:"contractAddress,omitempty"`
}

type Response struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

type ResponseContract struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
