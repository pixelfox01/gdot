package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func createSymlink(source, target string) error {
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("target already exists: %s", target)
	}
	if err := os.Symlink(source, target); err != nil {
		return err
	}
	fmt.Printf("Symlink created: %s -> %s\n", target, source)
	return nil
}

func createSymlinks() error {
	sourceDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	targetDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %v\n", err)
	}

	return filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == sourceDir {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		targetPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			fmt.Printf("Directory created: %s\n", targetPath)
		} else {
			targetDir := filepath.Dir(targetPath)
			if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create parent directory: %v", err)
			}
			if err := createSymlink(path, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink: %v", err)
			}
		}

		return nil
	})
}

func main() {
	if err := createSymlinks(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
