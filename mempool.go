package main

func main() {
	HandleTransaction("POST", "0x936ed523582d8d531eeef8c2ae7e478858f1b90c568badc95a581ea61267110f")
	HandleTransaction("DELETE", "0x936ed523582d8d531eeef8c2ae7e478858f1b90c568badc95a581ea61267110f")

	HandleAddress("POST", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	HandleAddress("DELETE", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	GetAddresses("ethereum", "main")
	GetTransactions("ethereum", "main")
}
