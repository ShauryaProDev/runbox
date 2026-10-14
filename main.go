package main

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var sourceFile *os.File
var items = make(map[string]string, 0)

var rootCmd = &cobra.Command{
	Use:   "box",
	Short: "A box tool",
}

var setCmd = &cobra.Command{
	Use:   "set <item> <command>",
	Short: "Set an item to runbox",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		readJson()
		items[args[0]] = args[1]
		writeJson()
	},
}

var runCmd = &cobra.Command{
	Use:   "run <item>",
	Short: "Run a command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		readJson()

		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("failed to get current dir")
		}

		command := items[args[0]]
		if command == "" {
			log.Fatalf("item %q not found", args[0])
		}

		_cmd := exec.Command("cmd", "/C", command)
		_cmd.Dir = currentDir
		_cmd.Stdout = os.Stdout
		_cmd.Stderr = os.Stderr

		if err := _cmd.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	currentUser, err := user.Current()

	if err != nil {
		log.Fatal("failed to read current user")
	}

	sourceDir := filepath.Join(currentUser.HomeDir, "runbox")

	err = os.Mkdir(sourceDir, 0755)

	sourceFilePath := filepath.Join(sourceDir, "runbox.json")

	sourceFile, err = os.OpenFile(sourceFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("failed to create source file")
	}

	fileInfo, err := sourceFile.Stat()

	if err != nil {
		log.Fatal("failed to init json file")
	}

	if fileInfo.Size() == 0 {
		sourceFile.WriteString("{}")
	}

	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(runCmd)
}

func initJson() {
	sourceFile.WriteString("{}")
}

func readJson() {
	absPath, err := filepath.Abs(sourceFile.Name())

	if err != nil {
		log.Fatal("failed to read abs path")
	}

	data, err := os.ReadFile(absPath)

	if err != nil {
		log.Fatal("failed to open source file")
	}

	err = json.Unmarshal(data, &items)

	if err != nil {
		log.Fatal("failed to convert json to data")
	}
}

func writeJson() {
	absPath, err := filepath.Abs(sourceFile.Name())

	if err != nil {
		log.Fatal("failed to read abs path")
	}

	data, err := json.Marshal(items)

	if err != nil {
		log.Fatal("failed to convert data to json")
	}

	err = os.WriteFile(absPath, data, 0644)

	if err != nil {
		log.Fatal("failed to write data to source file")
	}
}

func main() {
	rootCmd.Execute()
}
