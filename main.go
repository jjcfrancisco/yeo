package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func openConfigs(dbName string) (*database, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	raw, err := os.ReadFile("./databases.json")
	if err != nil {
		raw, err = os.ReadFile(homeDir + "/databases.json")
		if err != nil {
			fmt.Println("databases.json file is not present")
			os.Exit(1)
		}
	}

	var dbs databases

	err = json.Unmarshal([]byte(raw), &dbs)
	if err != nil {
		log.Fatalf("There's an issue with the databases.json file %s", err)
	}

	for _, db := range dbs.Dbs {
		if db.Name == dbName {
			return &db, nil
		}
	}

	return nil, fmt.Errorf("Could not find such database")

}

// Cleans up a database if user ctrl+c
func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

func main() {

	app := startCli()

	// Allows suggestions if mispellings
	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
