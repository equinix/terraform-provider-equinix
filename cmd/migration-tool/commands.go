package main

import (
	"fmt"
	"os"
)

// Copy target directory and append .backup
func CreateBackup(targetDir string, backupDir string) (err error) {
	fmt.Println("creating backup...", targetDir, "-->", backupDir)

	fi, err := os.Stat(targetDir)
	if err != nil {
		return fmt.Errorf("error reading directory\n %s", err)
	}

	if !fi.IsDir() {
		return fmt.Errorf("file targeted for migration")
	}

	_, err = os.Stat(backupDir)

	if err == nil {
		return fmt.Errorf("attempting to overwrite backups")
	}

	fmt.Println("copying", targetDir, "-->", backupDir)

	err = ProcessDirectory(targetDir, backupDir, CopyFile)

	if err != nil {
		return err
	}

	bfi, err := os.Stat(backupDir)

	if fi.Size() != bfi.Size() {
		return fmt.Errorf("backup corrupt")
	}

	fmt.Println("complete")
	return
}

// Overwrite target directory with contents of .backup directory
func RestoreBackup(backupDir string, targetDir string) (err error) {
	fmt.Println("restoring from backup...")

	fi, err := os.Stat(backupDir)
	if err != nil {
		return fmt.Errorf("error reading backup\n %s", err)
	}

	err = os.RemoveAll(targetDir)

	if err != nil {
		return fmt.Errorf("error removing original directory\n %s", err)
	}

	os.MkdirAll(targetDir, fi.Mode())

	err = ProcessDirectory(backupDir, targetDir, CopyFile)

	if err != nil {
		return fmt.Errorf("error restoring from backup directory\n %s", err)
	}

	fmt.Println("complete")
	return
}

// Remove .backup directory
func DeleteBackup(backupDir string) (err error) {
	fmt.Println("Purging backup...")

	err = os.RemoveAll(backupDir)

	if err != nil {
		return fmt.Errorf("error removing backup directory\n %s", err)
	}

	fmt.Println("complete")
	return
}

// Traverse all .tf files and apply transforms
func Migrate(targetDir string, backupDir string) (err error) {
	fmt.Println("migrating plan directory...")
	err = CreateBackup(targetDir, backupDir)

	if err != nil {
		return fmt.Errorf("error backing up directory before migration\n %s", err)
	}

	err = ProcessDirectory(targetDir, backupDir, MigratePlanFile, ".tf", ".tfstate")

	if err != nil {
		return fmt.Errorf("error removing backup directory\n %s", err)
	}

	fmt.Println("complete")
	return
}

// Traverse all .tf files and migrate or update Equinix provider
func MigrateProvider(targetDir string, backupDir string) (err error) {
	fmt.Println("scanning tf files for provider...")

	err = ProcessDirectory(targetDir, backupDir, TransformProvider, ".tf")

	if err != nil {
		return fmt.Errorf("error scanning providers for missing region value\n %s", err)
	}

	fmt.Println("complete")
	return
}
