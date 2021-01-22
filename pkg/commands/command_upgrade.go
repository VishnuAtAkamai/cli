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
	"github.com/akamai/cli/pkg/stats"
	"github.com/akamai/cli/pkg/version"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func cmdUpgrade(c *cli.Context) error {
	s := io.StartSpinner("Checking for upgrades...", "Checking for upgrades...... ["+color.GreenString("OK")+"]\n")

	if latestVersion := CheckUpgradeVersion(true); latestVersion != "" {
		io.StopSpinnerOk(s)
		fmt.Fprintf(app.App.Writer, "Found new version: %s (current version: %s)\n", color.BlueString("v"+latestVersion), color.BlueString("v"+version.Version))
		os.Args = []string{os.Args[0], "--version"}
		success := UpgradeCli(latestVersion)
		if success {
			stats.TrackEvent("upgrade.user", "success", "to: "+latestVersion+" from:"+version.Version)
		} else {
			stats.TrackEvent("upgrade.user", "failed", "to: "+latestVersion+" from:"+version.Version)
		}
	} else {
		io.StopSpinnerWarnOk(s)
		fmt.Fprintf(app.App.Writer, "Akamai CLI (%s) is already up-to-date", color.CyanString("v"+version.Version))
	}

	return nil
}
