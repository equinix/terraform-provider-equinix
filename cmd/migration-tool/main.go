package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing required command. One of [migrate, backup, version]")
		os.Exit(1)
	}

	if os.Args[1] == "-v" || os.Args[1] == "-version" || os.Args[1] == "--version" || os.Args[1] == "version" {
		fmt.Println(Version)
		os.Exit(0)
	}

	if os.Args[1] == "backup" {
		backup := flag.NewFlagSet("backup", flag.PanicOnError)
		backup.Usage = func() {
			backup.PrintDefaults()
			os.Exit(0)
		}
		dir := backup.String("dir", "", "Required, specify the plan directory to operate on")
		purge := backup.Bool("purge", false, "Optional, whether to purge the backup directory")
		restore := backup.Bool("restore", false, "Optional, whether to restore from the backup directory")

		err := backup.Parse(os.Args[2:])

		if *dir == "" {
			fmt.Println("Missing required directory flag\nCommand flags:")
			backup.PrintDefaults()
			os.Exit(1)
		}

		if err != nil {
			panic(err)
		}

		targetDir := path.Clean(*dir)
		backupDir := targetDir + ".backup"

		fmt.Println(targetDir)

		if *purge {
			err := DeleteBackup(backupDir)
			if err != nil {
				panic(err)
			}

			return
		}

		if *restore {
			err := RestoreBackup(backupDir, targetDir)
			if err != nil {
				panic(err)
			}

			return
		}

		err = CreateBackup(targetDir, backupDir)

		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	if os.Args[1] == "migrate" {
		migrate := flag.NewFlagSet("migrate", flag.PanicOnError)
		migrate.Usage = func() {
			migrate.PrintDefaults()
			os.Exit(0)
		}
		dir := migrate.String("dir", "", "Required, specify the plan directory to operate on")
		err := migrate.Parse(os.Args[2:])

		if *dir == "" {
			fmt.Println("Missing required directory flag\nCommand flags:")
			migrate.PrintDefaults()
			os.Exit(1)
		}

		if err != nil {
			panic(err)
		}

		targetDir := path.Clean(*dir)
		backupDir := targetDir + ".backup"

		err = Migrate(targetDir, backupDir)

		if err != nil {
			panic(err)
		}

		err = MigrateProvider(targetDir, backupDir)

		if err != nil {
			panic(err)
		}

		fmt.Println(`Migration Successful!`)
		os.Exit(0)
	}

	fmt.Println("Unknown command")
	os.Exit(1)
}
