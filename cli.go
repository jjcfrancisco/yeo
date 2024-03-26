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
				Usage:     "Creates a backup of a database. This is also known as 'dump'.",
				UsageText: fmt.Sprintf("backup [database] [filename]\n\nExample: yeo backup my_db my_db_backup.dump"),
				Action: func(cCtx *cli.Context) error {
					database := cCtx.Args().Get(0)
					filename := cCtx.Args().Get(1)
					os.Setenv("FILENAME", filename)
					if err := validateFilename(filename); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					err := backup(database, filename, false)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:      "revive",
				Usage:     "Revives a database from a backup. This is also known as 'restore'.",
				UsageText: "revive [database] [filename]\n\nExample: yeo revive my_new_db my_db_backup.dump",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "unlock",
						Usage: "The unlock flag allows to revive into non-local databases",
					},
				},
				Action: func(cCtx *cli.Context) error {
					database := cCtx.Args().Get(0)
					filename := cCtx.Args().Get(1)

					// Checks '.dump' filename provided is valid 
					if err := validateFilename(filename); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					// Checks the target database is allowed
					isUnlocked := cCtx.Bool("unlock")
					if err := validateTargetDb(database, isUnlocked); err != nil {
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
				Usage:     "Clones a database into another database. This is also known as 'dump' and 'restore'.",
				UsageText: "clone [database] [filename]\n\nExample: yeo revive my_new_db my_db_backup.dump",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "unlock",
						Usage: "The unlock flag allows to clone into non-local databases",
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
					isUnlocked := cCtx.Bool("unlock")

					// Checks the target database is allowed
					if err := validateTargetDb(target_database, isUnlocked); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Println() // Used for padding
					// Spin wheel
					w := wow.New(os.Stdout, spin.Get(spin.Dots), fmt.Sprintf(" Cloning %s database", og_database))
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
					if err = revive(target_database, temp_filename, false); err != nil {
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
