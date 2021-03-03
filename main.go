package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag {
			&cli.StringSliceFlag{
				Name: "repo",
				Aliases: []string{"r", "repos"},
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.StringSlice("repos")) > 0 {
				handleRepos(c.StringSlice("repo"))
			}

			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err)
	}

}

func handleRepos(repos []string) {
	version := "latest"

	current, _ := os.Getwd()

	versionFile, err := ioutil.ReadFile(path.Join(current, "VERSION"))

	if err == nil {
		version = string(versionFile)

		version = strings.TrimRight(strings.TrimRight(version, "\n"), "\r")
	}

	if _, err := os.Stat(path.Join(current, ".git")); !os.IsNotExist(err) {
		// git rev-list -1 HEAD
		cmd := exec.Command("git", "rev-list", "-1", "HEAD")
		stdout, err := cmd.Output()

		if err != nil {}

		version = fmt.Sprintf("%s-%s", version, string(stdout))
	}

	for _, repo := range repos {
		fmt.Println(repo, version)
	}
}
