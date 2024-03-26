package main

import (
	"bytes"
	"fmt"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"os"
	"os/exec"
)

func backup(database string, filename string, clone bool) error {
	configs, err := openConfigs(database)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	os.Setenv("PGPASSWORD", configs.Password)
	cmd := exec.Command("pg_dump", "-h", configs.Host, "-p", configs.Port, "-Fc", "-U", configs.User, "-d", configs.Database)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Stdout = file

	if !clone {
		fmt.Println()
		w := wow.New(os.Stdout, spin.Get(spin.Dots), fmt.Sprintf(" Backing up '%s' database", configs.Database))
		w.Start()
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(stderr.String())
	}

	fmt.Println("\n\n  Database backed up.")

	return nil
}

func revive(database string, filename string, clone bool) error {
	configs, err := openConfigs(database)
	if err != nil {
		return err
	}

	cmd := exec.Command("pg_restore", "-d", configs.Database, "-U", configs.User, "-h", configs.Host, "-p", configs.Port, "--no-owner", "--no-privileges", filename)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if !clone {
		fmt.Println()
		w := wow.New(os.Stdout, spin.Get(spin.Dots), fmt.Sprintf(" Reviving %s into database '%s'", filename, configs.Database))
		w.Start()
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(stderr.String())
	}

	if clone {
		fmt.Println("\n\n  Database cloned.")
	} else {
		fmt.Println("\n\n  Database revived.")
	}

	return nil

}

func prepareDb(database string) error {
	configs, err := openConfigs(database)
	if err != nil {
		return err
	}

	os.Setenv("PGPASSWORD", configs.Password)
	cmd := exec.Command("dropdb", "-U", configs.User, "-h", configs.Host, "-p", configs.Port, configs.Database, "--if-exists")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(stderr.String())
	}

	cmd = exec.Command("createdb", "-U", configs.User, "-h", configs.Host, "-p", configs.Port, configs.Database)
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(stderr.String())
	}

	return nil

}
