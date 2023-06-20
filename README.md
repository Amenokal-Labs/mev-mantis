# Etherscan API consumption in Golang

## Instructions to run

clone the repo: `git clone https://github.com/Amenokal-Labs/mev-mantis.git`.

navigate to the cloned folder, then run: `go run extract.go`.

## Docs

`func getBalance(tag, address string) string`: returns the Ether balance of a given address.

| Parameter | Description |
| --- | --- |
| address | the `string` representing the address to check for balance |
| tag | the `string` pre-defined block parameter, either `earliest`, `pending` or `latest` |

`func getBalances(tag string, addresses ...string) string`: returns the balance of the accounts from a list of addresses.

| Parameter | Description |
| --- | --- |
| address | the `strings` representing the addresses to check for balance, seperated by `,`. Up to 20 addresses per call|
| tag | the string pre-defined block parameter, either earliest, pending or latest |