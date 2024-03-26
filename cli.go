package main

import (
	"fmt"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"log"
	"os"
	"syscall"
	"os/signal"

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
	app := &cli.App{
		Name:  "Yeo",
		Usage: "Backup utilities for PostgreSQL databases",
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
				Action: func(cCtx *cli.Context) error {
					database := cCtx.Args().Get(0)
					filename := cCtx.Args().Get(1)
					if err := validateFilename(filename); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					if err := prepareDb(database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
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
				Action: func(cCtx *cli.Context) error {
					temp_filename := "temp.dump"
					os.Setenv("FILENAME", temp_filename)
					og_database := cCtx.Args().Get(0)
					target_database := cCtx.Args().Get(1)
					if err := validateTargetDb(target_database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Println()
					w := wow.New(os.Stdout, spin.Get(spin.Dots), fmt.Sprintf(" Cloning %s database", og_database))
					w.Start()
					err := backup(og_database, temp_filename, true)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					if err := prepareDb(target_database); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					if err = revive(target_database, temp_filename, false); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					cleanup()

					return nil
				},
			},
		},
	}

	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
