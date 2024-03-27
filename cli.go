package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"

	"github.com/urfave/cli/v2"
)

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
	// Version Flag
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Print out version",
	}

	// Main CLI
	app := &cli.App{
		Name:    "Yeo!",
		Version: "v0.2.0",
		Usage:   "Backup utilities for PostgreSQL databases",
		Commands: []*cli.Command{
			{
				Name:      "backup",
				Usage:     "creates a backup of a database. This is also known as 'dump'.",
				UsageText: "yeo backup [database] [filename]\n\nExample: yeo backup my_db my_db_backup.dump",
				Action: func(cCtx *cli.Context) error {
					database := cCtx.Args().Get(0)
					filename := cCtx.Args().Get(1)

					// Validates input
					if err := validateBackupArgs(cCtx.Args().Slice(), database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Validates that there's a connection with database
					if err := checkDbCons(database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					os.Setenv("FILENAME", filename)
					// Checks '.dump' filename provided is valid
					if err := validateFilename(filename); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Backs up a database
					if err := backup(database, filename, false); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					return nil

				},
			},
			{
				Name:      "revive",
				Usage:     "revives a database from a backup. This is also known as 'restore'.",
				UsageText: "yeo revive [options] [dump filename] [database]\n\nExample 1: yeo revive db_backup.dump local_db\n\nExample 2: yeo revive --allow db_backup.dump prod",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "allow",
						Usage: "allows to revive into non-local databases",
					},
				},
				Action: func(cCtx *cli.Context) error {
					filename := cCtx.Args().Get(0)
					database := cCtx.Args().Get(1)

					// Validates input
					if err := validateReviveArgs(cCtx.Args().Slice(), database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Checks '.dump' filename provided is valid
					if err := validateFilename(filename); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Validates that there's a connection with database
					if err := checkDbCons(database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Checks the target database is allowed
					isUnlocked := cCtx.Bool("allow")
					if err := validateTargetDb(filename, database, "revive", isUnlocked); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Prepares database by dropping and creating a new database
					if err := prepareDb(database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Restores a database
					if err := revive(database, filename, false); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					return nil

				},
			},
			{
				Name:      "clone",
				Usage:     "clones a database into another database. This is also known as 'dump' and 'restore'.",
				UsageText: "yeo clone [options] [origin database] [target database]\n\nExample 1: yeo clone prod local_db\n\nExample 2: yeo clone --allow local_db prod",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "allow",
						Usage: "allows to clone into non-local databases",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Cloning doesn't need specific filename and so a temporal name is given
					temp_filename := "temp.dump"
					// This is set for the 'cleanup' function
					os.Setenv("FILENAME", temp_filename)

					// User args
					og_database := cCtx.Args().Get(0)
					target_database := cCtx.Args().Get(1)

					// Validates input
					if err := validateCloneArgs(cCtx.Args().Slice(), og_database, target_database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					isUnlocked := cCtx.Bool("allow")

					// Checks the target database is allowed
					if err := validateTargetDb(og_database, target_database, "clone", isUnlocked); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					fmt.Println() // Used for padding
					// Spin wheel
					w := wow.New(os.Stdout, spin.Get(spin.Dots), fmt.Sprintf(" Cloning '%s' database", og_database))
					w.Start()

					// Dumps a database
					err := backup(og_database, temp_filename, true)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Prepares database by dropping and creating a new database
					if err := prepareDb(target_database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Restores a database
					if err = revive(target_database, temp_filename, true); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Removes the 'temp.dump' file generated during the process
					cleanup()

					return nil
				},
			},
		},
	}

	// Allows suggestions if mispellings
	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
