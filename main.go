package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
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

func createSymlinks(dryRun bool) error {
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

		if info.Name() == ".git" || info.Name() == ".gitignore" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil 
		}

		targetPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			if dryRun {
				fmt.Printf("[Dry Run] Directory would be created: %s\n", targetPath)
				return nil
			}
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			fmt.Printf("Directory created: %s\n", targetPath)
		} else {
			targetDir := filepath.Dir(targetPath)
			if dryRun {
				fmt.Printf("[Dry Run] Symlink would be created: %s -> %s\n", targetPath, path)
				return nil
			}
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
	var dryRun bool
	app := &cli.App{
		Name:  "gdot",
		Usage: "Creates symlinks for files and directories within the current working directory to the user's home directory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "dry-run",
				Aliases:     []string{"n"},
				Usage:       "perform a dry run without creating symlinks",
				Destination: &dryRun,
			},
		},
		Action: func(cCtx *cli.Context) error {
			createSymlinks(dryRun)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
