package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/annymsmthd/gogitver/pkg/git"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	gogit "gopkg.in/src-d/go-git.v4"
)

var rootCmd = &cobra.Command{
	Use:   "gogitver",
	Short: "gogitver is a semver generator that uses git history",
	Long:  ``,
	Run:   runRoot,
}

func init() {
	rootCmd.Flags().String("path", ".", "the path to the git repository")
	rootCmd.Flags().String("settings", "./.gogitver.yaml", "the file that contains the settings")
	rootCmd.Flags().Bool("forbid-behind-master", false, "error if the current branch's calculated version is behind the calculated version of refs/heads/master")
}

// Execute gogitver
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	f := cmd.Flag("path")
	sf := cmd.Flag("settings")
	fbm, err := strconv.ParseBool(cmd.Flag("forbid-behind-master").Value.String())
	if err != nil {
		fbm = false
	}

	var s *git.Settings
	_, err = os.Stat(sf.Value.String())
	if sf.Changed || err == nil {
		r, err := os.Open(sf.Value.String())
		if err != nil {
			panic(errors.Wrap(err, "cannot open settings file"))
		}

		s, err = git.GetSettingsFromFile(r)
		if err != nil {
			panic(err)
		}
	} else {
		s = git.GetDefaultSettings()
	}

	r, err := gogit.PlainOpen(f.Value.String())
	if err != nil {
		panic(err)
	}

	version, err := git.GetCurrentVersion(r, s, false, fbm)
	if err != nil {
		panic(err)
	}

	fmt.Println(version)
}
