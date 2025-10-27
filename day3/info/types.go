package info

type TokenInfo struct {
	Address string `json:"address"`
	Symbol  string `json:"symbol"`
	Balance string `json:"balance"`
}

type WalletResponse struct {
	Address    string      `json:"address"`
	EthBalance string      `json:"eth_balance"`
	Tokens     []TokenInfo `json:"tokens"`
}
