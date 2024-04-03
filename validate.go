package main

import (
	"fmt"
	"os"
	"strings"
)

func cleanup() {
	filename := os.Getenv("FILENAME")
	if err := os.Remove("./" + filename); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validateFilename(fn string) error {

	if !strings.Contains(fn, ".dump") {
		return fmt.Errorf("\nThe file must contain its '.dump' extension i.e. 'my_file.dump'\n")
	}

	return nil
}

func validateBackupArgs(args []string, database string) error {

	if len(args) <= 1 {
		return fmt.Errorf("\nError: You must pass valid arguments.\n\nExample: yeo backup prod prod.dump\n")
	} else if len(args) > 2 {
		return fmt.Errorf("\nToo many arguments passed.\n\nExample: yeo backup prod prod.dump")
	}

	_, err := openConfigs(database)
	if err != nil {
		return fmt.Errorf("\nNo credentials for '%s' database.\n", database)
	}

	return nil

}

func validateReviveArgs(args []string, database string) error {

	if len(args) <= 1 {
		return fmt.Errorf("\nError: You must pass valid arguments.\n\nExample: yeo revive development.dump prod\n")
	} else if len(args) > 2 {
		return fmt.Errorf("\nToo many arguments passed.\n\nExample: yeo revive development.dump prod")
	}

	_, err := openConfigs(database)
	if err != nil {
		return fmt.Errorf("\nNo credentials for '%s' database.\n", database)
	}

	return nil

}

func validateTargetDb(origin string, targetDb string, op string, unlock bool) error {

	configs, err := openConfigs(targetDb)
	if err != nil {
		return err
	}

	if unlock {
		return nil
	}

	if !configs.IsLocal && op == "revive" {
		err := fmt.Sprintf("\nError: Your target database isn't local. This is a security lock. To remove it, use the '--allow' flag\n\nExample: yeo revive --allow %s %s\n", origin, targetDb)
		return fmt.Errorf(err)
	} else if !configs.IsLocal && op == "clone" {
		err := fmt.Sprintf("\nError: Your target database isn't local. This is a security lock. To remove it, use the '--allow' flag\n\nExample: yeo clone --allow %s %s\n", origin, targetDb)
		return fmt.Errorf(err)
	}

	return nil
}

func validateCloneArgs(args []string, originDb string, targetDb string) error {

	if len(args) <= 1 {
		return fmt.Errorf("\nError: You must pass valid arguments.\n\nExample: yeo clone development prod\n")
	} else if len(args) > 2 {
		return fmt.Errorf("\nToo many arguments passed.\n\nExample: yeo clone development prod")
	}

	_, err := openConfigs(originDb)
	if err != nil {
		return fmt.Errorf("\nNo credentials for '%s' database.\n", originDb)
	}
	_, err = openConfigs(targetDb)
	if err != nil {
		return fmt.Errorf("\nNo credentials for '%s' database.\n", targetDb)
	}

	return nil

}
