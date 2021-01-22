// Copyright 2018. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"github.com/akamai/cli/pkg/app"
	"github.com/akamai/cli/pkg/io"
	"github.com/akamai/cli/pkg/packages"
	"github.com/akamai/cli/pkg/stats"
	"github.com/akamai/cli/pkg/tools"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
)

var (
	// ThirdPartyDisclaimer is the message to be used when third party packages are installed
	ThirdPartyDisclaimer = color.CyanString("Disclaimer: You are installing a third-party package, subject to its own terms and conditions. Akamai makes no warranty or representation with respect to the third-party package.")
)

func cmdInstall(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.NewExitError(color.RedString("You must specify a repository URL"), 1)
	}

	oldCmds := getCommands()

	for _, repo := range c.Args() {
		repo = tools.Githubize(repo)
		err := installPackage(repo, c.Bool("force"))
		if err != nil {
			// Only track public github repos
			if isPublicRepo(repo) {
				stats.TrackEvent("package.install", "failed", repo)
			}
			return err
		}

		if isPublicRepo(repo) {
			stats.TrackEvent("package.install", "success", repo)
		}
	}

	packageListDiff(oldCmds)

	return nil
}

func packageListDiff(oldcmds []subcommands) {
	cmds := getCommands()

	var old []command
	for _, oldcmd := range oldcmds {
		for _, cmd := range oldcmd.Commands {
			old = append(old, cmd)
		}
	}

	var newCmds []command
	for _, newcmd := range cmds {
		for _, cmd := range newcmd.Commands {
			newCmds = append(newCmds, cmd)
		}
	}

	var added = make(map[string]bool)
	var removed = make(map[string]bool)

	for _, newCmd := range newCmds {
		found := false
		for _, oldCmd := range old {
			if newCmd.Name == oldCmd.Name {
				found = true
				break
			}
		}

		if !found {
			added[newCmd.Name] = true
		}
	}

	for _, oldCmd := range old {
		found := false
		for _, newCmd := range newCmds {
			if newCmd.Name == oldCmd.Name {
				found = true
				break
			}
		}

		if !found {
			removed[oldCmd.Name] = true
		}
	}

	listInstalledCommands(added, removed)
}

func isPublicRepo(repo string) bool {
	return !strings.Contains(repo, ":") || strings.HasPrefix(repo, "https://github.com/")
}

func installPackage(repo string, forceBinary bool) error {
	srcPath, err := tools.GetAkamaiCliSrcPath()
	if err != nil {
		return err
	}

	term := terminal.Standard()

	spin := term.Spinner()

<<<<<<< HEAD
	s := io.StartSpinner(fmt.Sprintf("Attempting to fetch command from %s...", repo), fmt.Sprintf("Attempting to fetch command from %s...", repo)+"... ["+color.GreenString("OK")+"]\n")
=======
	spin.Start("Attempting to fetch command from %s", repo)
>>>>>>> TFP-541

	dirName := strings.TrimSuffix(filepath.Base(repo), ".git")
	packageDir := filepath.Join(srcPath, dirName)
	if _, err = os.Stat(packageDir); err == nil {
<<<<<<< HEAD
		io.StopSpinnerFail(s)

=======
		spin.Stop(terminal.SpinnerStatusFail)
>>>>>>> TFP-541
		return cli.NewExitError(color.RedString("Package directory already exists (%s)", packageDir), 1)
	}

	_, err = git.PlainClone(packageDir, false, &git.CloneOptions{
		URL:      repo,
		Progress: spin,
		Depth:    1,
	})

	if err != nil {
		spin.Stop(terminal.SpinnerStatusFail)

		os.RemoveAll(packageDir)

<<<<<<< HEAD
		io.StopSpinnerFail(s)
		return cli.NewExitError(color.RedString("Unable to clone repository: "+err.Error()), 1)
	}

	io.StopSpinnerOk(s)

	if strings.HasPrefix(repo, "https://github.com/akamai/cli-") != true && strings.HasPrefix(repo, "git@github.com:akamai/cli-") != true {
		fmt.Fprintln(app.App.Writer, color.CyanString("Disclaimer: You are installing a third-party package, subject to its own terms and conditions. Akamai makes no warranty or representation with respect to the third-party package."))
=======
		return cli.NewExitError(color.RedString("Unable to clone repository: "+err.Error()), 1)
	}

	spin.Stop(terminal.SpinnerStatusOK)

	if strings.HasPrefix(repo, "https://github.com/akamai/cli-") != true && strings.HasPrefix(repo, "git@github.com:akamai/cli-") != true {
		term.Writef(ThirdPartyDisclaimer)
>>>>>>> TFP-541
	}

	if !installPackageDependencies(packageDir, forceBinary) {
		os.RemoveAll(packageDir)
		return cli.NewExitError("", 1)
	}

	return nil
}

func installPackageDependencies(dir string, forceBinary bool) bool {
<<<<<<< HEAD
	s := io.StartSpinner("Installing...", "Installing...... ["+color.GreenString("OK")+"]\n")
=======
	term := terminal.Standard()

	spin := term.Spinner()

	spin.Start("Installing...")
>>>>>>> TFP-541

	cmdPackage, err := readPackage(dir)

	if err != nil {
		io.StopSpinnerFail(s)
		fmt.Fprintln(app.App.Writer, err.Error())
		return false
	}

	lang := determineCommandLanguage(cmdPackage)

	var success bool
	switch lang {
	case "php":
		success, err = packages.InstallPHP(dir, cmdPackage.Requirements.Php)
	case "javascript":
		success, err = packages.InstallJavaScript(dir, cmdPackage.Requirements.Node)
	case "ruby":
		success, err = packages.InstallRuby(dir, cmdPackage.Requirements.Ruby)
	case "python":
		success, err = packages.InstallPython(dir, cmdPackage.Requirements.Python)
	case "go":
		var commands []string
		for _, cmd := range cmdPackage.Commands {
			commands = append(commands, cmd.Name)
		}
		success, err = packages.InstallGolang(dir, cmdPackage.Requirements.Go, commands)
	default:
<<<<<<< HEAD
		io.StopSpinnerWarnOk(s)
		fmt.Fprintln(app.App.Writer, color.CyanString("Package installed successfully, however package type is unknown, and may or may not function correctly."))
=======
		spin.Stop(terminal.SpinnerStatusWarnOK)

		term.Writef(color.CyanString("Package installed successfully, however package type is unknown, and may or may not function correctly."))

>>>>>>> TFP-541
		return true
	}

	if success && err == nil {
<<<<<<< HEAD
		io.StopSpinnerOk(s)
=======
		spin.Stop(terminal.SpinnerStatusOK)
>>>>>>> TFP-541
		return true
	}

	first := true
	for _, cmd := range cmdPackage.Commands {
		if cmd.Bin != "" {
			if first {
				first = false
<<<<<<< HEAD
				io.StopSpinnerWarn(s)
				fmt.Fprintln(app.App.Writer, color.CyanString(err.Error()))
=======

				spin.Stop(terminal.SpinnerStatusWarnOK)

				term.Writef(color.CyanString(err.Error()))

>>>>>>> TFP-541
				if !forceBinary {
					if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
						return false
					}

<<<<<<< HEAD
					fmt.Fprint(app.App.Writer, "Binary command(s) found, would you like to download and install it? (Y/n): ")
					answer := ""
					fmt.Scanln(&answer)
					if answer != "" && strings.ToLower(answer) != "y" {
=======
					answer, _ := term.Confirm("Binary command(s) found, would you like to download and install it?", true)

					if !answer {
>>>>>>> TFP-541
						return false
					}
				}

				os.MkdirAll(filepath.Join(dir, "bin"), 0700)

<<<<<<< HEAD
				s = io.StartSpinner("Downloading binary...", "Downloading binary...... ["+color.GreenString("OK")+"]\n")
			}

			if !downloadBin(filepath.Join(dir, "bin"), cmd) {
				io.StopSpinnerFail(s)
				fmt.Fprintln(app.App.Writer, color.RedString("Unable to download binary: "+err.Error()))
=======
				spin.Start("Downloading binary...")
			}

			if !downloadBin(filepath.Join(dir, "bin"), cmd) {
				spin.Stop(terminal.SpinnerStatusFail)

				term.Writef(color.RedString("Unable to download binary: " + err.Error()))

>>>>>>> TFP-541
				return false
			}
		}

		if first {
			first = false
<<<<<<< HEAD
			io.StopSpinnerFail(s)
			fmt.Fprintln(app.App.Writer, color.RedString(err.Error()))
=======

			spin.Stop(terminal.SpinnerStatusFail)

			term.Writef(color.RedString(err.Error()))

>>>>>>> TFP-541
			return false
		}
	}

<<<<<<< HEAD
	io.StopSpinnerOk(s)
=======
	spin.Stop(terminal.SpinnerStatusOK)

>>>>>>> TFP-541
	return true
}
