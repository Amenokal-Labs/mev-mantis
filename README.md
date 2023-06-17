# Etherscan API consumption in Golang

## Instructions to run

clone the repo: `git clone https://github.com/Amenokal-Labs/mev-mantis.git`.

navigate to the cloned folder, then run: `go run extract.go`.

## Docs

`func getBalance(tag, address string) string`:

Returns the Ether balance of a given address.

| Parameter index | Parameter | Description |
| --- | --- | --- |
| 0 | tag | the string pre-defined block parameter, either earliest, pending or latest |
| 1 | address | the string representing the address to check for balance |