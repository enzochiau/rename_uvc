package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/voxelbrain/goptions"
)

var options struct {
	RootDir    string        `goptions:"-d, --dir,  description='Directory where the rename happens'"`
	OldPattern string        `goptions:"-f, --from, obligatory, description='Original pattern'"`
	NewPattern string        `goptions:"-t, --to,   obligatory, description='New pattern'"`
	Help       goptions.Help `goptions:"-h, --help, description='Show this help'"`
}

func renamePackage(fileName, old, new string) error {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	lines := strings.Split(string(input), "\n")

	NEW := strings.ToUpper(new)
	OLD := strings.ToUpper(old)
	for i := range lines {
		lines[i] = strings.Replace(lines[i], old, new, -1)
		lines[i] = strings.Replace(lines[i], OLD, NEW, -1)
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(fileName, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func rename(rootDir, old, new string) error {
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			fileName := f.Name()
			if strings.HasPrefix(fileName, old) {
				newName := strings.Replace(fileName, old, new, 1)
				err = os.Rename(path.Join(rootDir, fileName), path.Join(rootDir, newName))
				if err != nil {
					return err
				}
				err = renamePackage(path.Join(rootDir, newName), old+"_", new+"_")
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	goptions.ParseAndFail(&options)

	// Use current directory by default
	if options.RootDir == "" {
		options.RootDir = os.Getenv("PWD")
	}

	options.RootDir = strings.TrimRight(options.RootDir, "/")

	fmt.Printf("%+v\n", options)

	err := rename(options.RootDir, options.OldPattern, options.NewPattern)

	if err != nil {
		fmt.Println(err)
	}
}
