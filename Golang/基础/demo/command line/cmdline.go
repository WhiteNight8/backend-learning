package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Command line arguments
type CommandFlags struct {
	Create  bool
	Read    bool
	Write   bool
	Copy    bool
	Delete  bool
	List    bool
	Rename  bool
	Append  bool
	Help    bool
	Path    string
	Content string
	Dest    string
}

func main() {
	// initialize command line arguments
	cmdFlags := parseFlags()

	//display help message if -help flag is set
	if cmdFlags.Help {
		printHelp()
		return
	}

	//execute command based on flags
	switch {
	case cmdFlags.Create:
		// create a new file
		if cmdFlags.Path == "" {
			fmt.Println("Path is required for creating a file.")
			return
		}
		if err := createFile(cmdFlags.Path); err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		} else {
			fmt.Printf("File created successfully: %s\n", cmdFlags.Path)
		}
	case cmdFlags.Read:
		// read a file
		if cmdFlags.Path == "" {
			fmt.Println("Path is required for reading a file.")
			return
		}
		content, err := readFile(cmdFlags.Path)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		} else {
			fmt.Printf("File content:\n%s\n", content)
		}
	case cmdFlags.Write:
		// write to a file
		if cmdFlags.Path == "" {
			fmt.Println("Path is required for writing to a file.")
			return
		}
		err := writeFile(cmdFlags.Path, cmdFlags.Content)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		} else {
			fmt.Printf("File written successfully: %s\n", cmdFlags.Path)
		}
	case cmdFlags.Append:
		// append to a file
		if cmdFlags.Path == "" {
			fmt.Println("Path is required for appending to a file.")
			return
		}
		err := appendToFile(cmdFlags.Path, cmdFlags.Content)
		if err != nil {
			fmt.Printf("Error appending to file: %v\n", err)
			return
		} else {
			fmt.Printf("File appended successfully: %s\n", cmdFlags.Path)
		}
	case cmdFlags.Copy:
		// copy a file
		if cmdFlags.Path == "" || cmdFlags.Dest == "" {
			fmt.Println("Path and destination are required for copying a file.")
			return
		}
		err := copyFile(cmdFlags.Path, cmdFlags.Dest)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			return
		} else {
			fmt.Printf("File copied successfully from %s to %s\n", cmdFlags.Path, cmdFlags.Dest)
		}
	case cmdFlags.Delete:
		// delete a file
		if cmdFlags.Path == "" {
			fmt.Println("Path is required for deleting a file.")
			return
		}
		err := deleteFile(cmdFlags.Path)
		if err != nil {
			fmt.Printf("Error deleting file: %v\n", err)
			return
		} else {
			fmt.Printf("File deleted successfully: %s\n", cmdFlags.Path)
		}
	case cmdFlags.List:
		// list files in a directory
		if cmdFlags.Path == "" {
			fmt.Println("Path is required for listing files in a directory.")
			return
		}
		files, err := listFiles(cmdFlags.Path)
		if err != nil {
			fmt.Printf("Error listing files: %v\n", err)
			return
		} else {
			fmt.Println("Files in directory:")
			for _, file := range files {
				fmt.Println(file)
			}
		}
	case cmdFlags.Rename:
		// rename a file
		if cmdFlags.Path == "" || cmdFlags.Dest == "" {
			fmt.Println("Path and destination are required for renaming a file.")
			return
		}
		err := renameFile(cmdFlags.Path, cmdFlags.Dest)
		if err != nil {
			fmt.Printf("Error renaming file: %v\n", err)
			return
		} else {
			fmt.Printf("File renamed successfully from %s to %s\n", cmdFlags.Path, cmdFlags.Dest)
		}
	default:
		// if no flags are set, show help message
		printHelp()

	}
}

// parse command line arguments
func parseFlags() CommandFlags {
	var cmdFlags CommandFlags

	flag.BoolVar(&cmdFlags.Create, "create", false, "Create a new file")
	flag.BoolVar(&cmdFlags.Read, "read", false, "Read a file")
	flag.BoolVar(&cmdFlags.Write, "write", false, "Write to a file")
	flag.BoolVar(&cmdFlags.Copy, "copy", false, "Copy a file")
	flag.BoolVar(&cmdFlags.Delete, "delete", false, "Delete a file")
	flag.BoolVar(&cmdFlags.List, "list", false, "List files in a directory")
	flag.BoolVar(&cmdFlags.Rename, "rename", false, "Rename a file")
	flag.BoolVar(&cmdFlags.Append, "append", false, "Append to a file")
	flag.BoolVar(&cmdFlags.Help, "help", false, "Show help message")
	flag.StringVar(&cmdFlags.Path, "path", "", "Path to the file or directory")
	flag.StringVar(&cmdFlags.Content, "content", "", "Content to write to the file")
	flag.StringVar(&cmdFlags.Dest, "dest", "", "Destination path for copy or rename")

	flag.Parse()
	return cmdFlags
}

// show help message
func printHelp() {
	helpText := `
Usage: fileutil [options]
Options:
	-create   Create a new file		
	-read     Read a file
	-write    Write to a file
	-copy     Copy a file
	-delete   Delete a file
	-list     List files in a directory
	-rename   Rename a file
	-append   Append to a file
	-help     Show help message
	-path     Path to the file or directory
	-content  Content to write to the file
	-dest    Destination path for copy or rename


Examples:
	fileutil -create -path /path/to/file.txt -content "Hello, World!"
	fileutil -read -path /path/to/file.txt
	fileutil -write -path /path/to/file.txt -content "New content"
	fileutil -copy -path /path/to/file.txt -dest /path/to/copy.txt
	fileutil -delete -path /path/to/file.txt
	fileutil -list -path /path/to/directory
	fileutil -rename -path /path/to/file.txt -dest /path/to/newfile.txt
	fileutil -append -path /path/to/file.txt -content "Appended content"
`
	fmt.Println(helpText)
}

// create a new file
func createFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// read a file
func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// write to a file
func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// append to a file
func appendToFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

// copy a file
func copyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// delete a file
func deleteFile(path string) error {
	return os.Remove(path)
}

// list files in a directory
func listFiles(path string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		fileInfo := entry.Name()
		if entry.IsDir() {
			fileInfo += "/"
		}
		files = append(files, fileInfo)
	}

	return files, nil
}

// rename a file
func renameFile(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}
