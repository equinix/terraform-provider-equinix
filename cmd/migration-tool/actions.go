package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Individual file io strategies for different operations
type FileAction func(string, string) error

// Traverse a directory, executing the supplied FileAction on each file
func ProcessDirectory(targetDir string, backupDir string, fileActionFn FileAction, targetExtns ...string) (err error) {
	_, err = os.Stat(targetDir)

	if err != nil {
		return fmt.Errorf("error reading directory\n %s", err)
	}

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("error reading directory contents \n %s", err)
	}

	for _, res := range files {
		targetRes := path.Join(targetDir, res.Name())
		backupRes := path.Join(backupDir, res.Name())

		if res.IsDir() {
			err = ProcessDirectory(targetRes, backupRes, fileActionFn, targetExtns...)

			if err != nil {
				return err
			}
		} else {
			if len(targetExtns) == 0 {
				err = fileActionFn(targetRes, backupRes)

				if err != nil {
					return err
				}
			} else {
				if contains(targetExtns, filepath.Ext(res.Name())) {
					err = fileActionFn(targetRes, backupRes)

					if err != nil {
						return err
					}
				} else {
					fmt.Println("Skipping: ", targetDir)
				}
			}
		}
	}

	return
}

// Copy file from targetFile path to backupFile path
func CopyFile(targetFile string, backupFile string) (err error) {
	// make sure directory structure exists
	bkDir := path.Dir(backupFile)
	_, err = os.Stat(bkDir)

	if err != nil {
		if os.IsNotExist(err) {
			oDir := path.Dir(targetFile)
			fi, err := os.Stat(oDir)
			if err != nil {
				return fmt.Errorf("error reading original directory %s", err)
			}

			err = os.MkdirAll(bkDir, fi.Mode())

			if err != nil {
				return fmt.Errorf("error creating directory for file %s", err)
			}
		} else {
			return fmt.Errorf("unexpected error reading original directory %s", err)
		}
	}

	src, err := os.Open(targetFile)
	if err != nil {
		return fmt.Errorf("error reading original file\n %s", err)
	}

	defer src.Close()

	dst, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("error creating backup file\n %s", err)
	}

	defer dst.Close()

	fmt.Printf("Copying %s --> %s", targetFile, backupFile)
	size, err := io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error writing file\n %s", err)
	}

	fmt.Printf(", %d bytes\n", size)
	return
}

// Read file from backup location, apply transforms and overwrite original file
func MigratePlanFile(targetFile string, backupFile string) (err error) {
	src, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("error reading file\n %s", err)
	}

	defer src.Close()

	dst, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("error creating write location\n %s", err)
	}

	defer dst.Close()

	wrtr := bufio.NewWriter(dst)

	var replaceStrategy func(string) string
	if filepath.Ext(backupFile) == ".tf" {
		replaceStrategy = replaceTemplateTokens
	} else {
		replaceStrategy = replaceStatefileTokens
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		str := scanner.Text()
		str = replaceStrategy(str)
		fmt.Fprintln(wrtr, str)
	}
	wrtr.Flush()

	return
}

// Scan TF files for terraform:required_providers and provider blocks and define or update Equinix provider
func TransformProvider(targetFile string, backupFile string) error {
	fmt.Printf("Scanning %s\n", targetFile)

	fileInfo, err := os.Stat(targetFile)
	if err != nil {
		return fmt.Errorf("error while updating provider\n %s", err)
	}

	const maxSize = 1024 * 1024
	if fileInfo.Size() > maxSize {
		return fmt.Errorf("file too large to process")
	}

	fileBytes, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("error updating terraform:required_providers block\n %s", err)
	}

	content := string(fileBytes)

	content, err = scanAndUpdateRequiredProvider(content)
	if err != nil {
		return fmt.Errorf("error updating terraform:required_providers block\n %s", err)
	}

	content, err = scanAndUpdateProvider(content)
	if err != nil {
		return fmt.Errorf("error updating provider block\n %s", err)
	}

	ioutil.WriteFile(targetFile, []byte(content), fileInfo.Mode())

	return err
}

// find a string in a slice of strings
func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
