package find

import (
	"debricked/pkg/client"
	"debricked/pkg/file"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var debClient *client.DebClient
var finder *file.Finder

func NewFindCmd(debrickedClient *client.DebClient) *cobra.Command {
	debClient = debrickedClient
	finder, _ = file.NewFinder(debClient)
	cmd := &cobra.Command{
		Use:   "find [path]",
		Short: "Find all dependency files in inputted path",
		Long: `Find all dependency files in inputted path. Related files are grouped together. 
For example ` + "`package.json`" + ` with ` + "`package-lock.json`.",
		Args: validateArgs,
		RunE: run,
	}

	return cmd
}

func run(_ *cobra.Command, args []string) error {
	directoryPath := args[0]

	return find(directoryPath, []string{})
}

func find(path string, exclusions []string) error {
	fileGroups, err := finder.GetGroups(path, exclusions)
	if err != nil {
		return err
	}
	for _, group := range fileGroups {
		group.Print()
	}

	return nil
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires path")
	}
	if isValidFilepath(args[0]) {
		return nil
	}
	return fmt.Errorf("invalid path specified: %s", args[0])
}

func isValidFilepath(path string) bool {
	_, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}

	return true
}