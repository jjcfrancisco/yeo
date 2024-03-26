package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func openConfigs(dbName string) (*database, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	raw, err := os.ReadFile("./databases.json")
	if err != nil {
		raw, err = os.ReadFile(homeDir+"/databases.json")
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

func validateFilename(fn string) error {

	if !strings.Contains(fn, ".dump") {
		return fmt.Errorf(`You must provide a valid file i.e. "my_file.dump"`)
	}

	return nil
}

func validateTargetDb(targetDb string) error {

	configs, err := openConfigs(targetDb)
	if err != nil {
		return err
	}

	if !strings.Contains(configs.Host, "localhost") {
		return fmt.Errorf(`Your target database isn't local.`)
	}

	return nil
}
