package main

import utils "github.com/Amenokal-Labs/mev-mantis.git/pkg/utils"

func main() {
	utils.HandleTransaction("POST", "0x936ed523582d8d531eeef8c2ae7e478858f1b90c568badc95a581ea61267110f")
	utils.HandleTransaction("DELETE", "0x936ed523582d8d531eeef8c2ae7e478858f1b90c568badc95a581ea61267110f")

	utils.HandleAddress("POST", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	utils.HandleAddress("DELETE", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	utils.GetAddresses("ethereum", "main")
	utils.GetTransactions("ethereum", "main")
}
