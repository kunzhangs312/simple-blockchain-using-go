package main

import (
	"github.com/blockchain_go/03-Persistence-and-CLI/blc"
)

func main() {
	bc := blc.NewBlockchain()
	defer bc.Db.Close()

	cli := blc.CLI{bc}
	cli.Run()
}
