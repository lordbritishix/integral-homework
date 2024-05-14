NOTES

`POST /accounts/{accountId}/sync` - starts syncing transactions in the wallet - this is an asynchronous operation. 
Do this initially before calling `GET /accounts/{accountId}/transactions`. This is done because some wallets have millions
of transactions and so it is ideal that we run syncs out of band.

`/accounts/{accountId}/transactions` - returns transactions for a given wallet.

`/steth/stats` - not enough time to implement, but it looks like this is the proxy contract that can be used to get info around shares https://etherscan.io/token/0xae7ab96520de3a18e5e111b5eaab095312d7fe84#readProxyContract

`/steth/last-deposits` - not enough time to implement, but I think I maybe able to use this etherscan api call
https://api.etherscan.io/api?module=account&action=tokentx&contractAddress=0xae7ab96520de3a18e5e111b5eaab095312d7fe84&to=0xae7ab96520de3a18e5e111b5eaab095312d7fe84&sort=desc&offset=10&page=1&apikey=YourApiKeyToken in order to
get transfers made on that smart contract. I just can't differentiate if it is a rewards event or a deposit to the pool cause there's no method id.

There is 1 hard coded account created by the code. Use the account id `account_abc` to see how this works
```
model.Account{
		AccountId:   "account_abc",
		AccountName: "Jim Quitevis",
		Wallet: model.Wallet{
			WalletId:   "wallet_abc",
			WalletName: "Jim's Eth Wallet",
			Address:    "0x912fd21d7a69678227fe6d08c64222db41477ba0",
			Network:    model.EthereumNetwork,
		},
	}
```
