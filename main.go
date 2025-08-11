package main

import (
	"fmt"
	"os"

	"github.com/wjseele/gator/internal/config"
)

func main() {
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	err = dbConfig.SetUser("wjseele")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	dbConfig, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(dbConfig)
}
