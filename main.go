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
			&cli.BoolFlag{
				Name: "git",
			},
			&cli.StringFlag{
				Name: "version",
				Aliases: []string{"v"},
				DefaultText: "",
			},
			&cli.StringFlag{
				Name: "version-file",
				Aliases: []string{"vf"},
				DefaultText: "",
			},
			&cli.StringFlag{
				Name: "name",
				Aliases: []string{"n"},
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.StringSlice("repos")) > 0 {
				handleRepos(c.String("name"), c.StringSlice("repo"), c.String("version"), c.String("version-file"), c.Bool("git"))
			}

			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err)
	}

}

func handleRepos(name string, repos []string, version string, versionFileName string, git bool) {
	current, _ := os.Getwd()

	if versionFileName != "" {
		versionFile, err := ioutil.ReadFile(path.Join(current, versionFileName))

		if err == nil {
			version = string(versionFile)

			version = strings.TrimRight(strings.TrimRight(version, "\n"), "\r")
		}
	}

	if git {
		if _, err := os.Stat(path.Join(current, ".git")); !os.IsNotExist(err) {
			cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
			stdout, err := cmd.Output()

			if err != nil {}

			if len(version) > 0 {
				version = fmt.Sprintf("%s-%s", version, string(stdout))
			} else {
				version = string(stdout)
			}
		}
	}

	images := ""

	for _, repo := range repos {
		images = strings.TrimSuffix(images + " " + fmt.Sprintf("-t %s/%s:%s", repo, name, version), "\n")
	}

	run(fmt.Sprintf("docker build %s .", images))

	for _, image := range strings.Split(strings.ReplaceAll(strings.TrimPrefix(images, " "), "-t ", ""), " ") {
		run(fmt.Sprintf("docker push %s", image))
	}
}

func run(cmd string) {
	fmt.Println(cmd)

	parts := strings.Fields(cmd)

	head := parts[0]
  parts = parts[1:]

  out, err := exec.Command(head,parts...).Output()

	if err != nil {
			fmt.Println("error occured")
			fmt.Printf("%s", err)
	}

	fmt.Printf("%s", out)
}
