# Etherscan API consumption in Golang

## Instructions to run

clone the repo: `git clone https://github.com/Amenokal-Labs/mev-mantis.git`.

navigate to the cloned folder, then run: `go run extract.go`.

## Docs
<<<<<<< HEAD

`func getBalance(tag, address string) string`:
=======
>>>>>>> 9d2334f (fix: updated README.md)

`func getBalance(tag, address string) string`: returns the Ether balance of a given address.

| Parameter index | Parameter | Description |
| --- | --- | --- |
| 0 | tag | the string pre-defined block parameter, either earliest, pending or latest |
| 1 | address | the string representing the address to check for balance |

`func getBalances(tag string, addresses ...string) string`: returns the balance of the accounts from a list of addresses.

| Parameter index | Parameter | Description |
| --- | --- | --- |
| 0 | tag | the string pre-defined block parameter, either earliest, pending or latest |
| 1 | address | the strings representing the addresses to check for balance. Up to 20 addresses per call. |