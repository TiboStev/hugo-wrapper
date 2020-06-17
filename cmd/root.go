/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/TiboStev/cobra"
	"github.com/TiboStev/hugo-wrapper/versionmanager"
	"github.com/TiboStev/pflag"
	homedir "github.com/mitchellh/go-homedir"
)

// rootCmd represents the base command when called without any subcommands

var hugoVersion string

//var onWrapper bool
var rootCmd = &cobra.Command{
	Use:                          "hugo-wrapper",
	Short:                        "Wrap hugo command",
	Long:                         `This is a wrapper for the hugo command, it allows to use different version of hugo without struggle.`,
	PersistentFParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Args:                         cobra.ArbitraryArgs,
	PersistentPreRun:             persistantPreRun,
	PersistentPostRun:            persistantPostRun,
	Run:                          func(cmd *cobra.Command, args []string) {},
}

// Execute ads all child commands to the root command and sets flags appropriately.
// This is called y main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&hugoVersion, "hugo-version", "latest", "use this specific hugo version")
}

var wrappedArgs []string
var wrappedFlags []*pflag.Flag
var hugoCommand *exec.Cmd

func persistantPreRun(cmd *cobra.Command, args []string) {
	wrappedArgs = []string{}
	if cmd.HasParent() {
		wrappedArgs = append(wrappedArgs, cmd.Name())
	}

	cmd.VisitParents(func(parent *cobra.Command) {
		if parent.HasParent() {
			wrappedArgs = append([]string{parent.Name()}, wrappedArgs...)
		}
	})
	wrappedArgs = append(wrappedArgs, args...)
	// For every set flag, append its name and value
	cmd.Flags().VisitUnknowns(func(flag *pflag.Flag) {
		wrappedFlags = append(wrappedFlags, flag)
	})
}

func getHugoCommand(args []string) (*exec.Cmd, error) {

	homePath, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	hugoVersionManagerPath := path.Join(homePath, ".hugo-wrapper")
	if _, err := os.Stat(hugoVersionManagerPath); err != nil {
		fmt.Println("creation of ~/.hugo-wrapper")
		os.Mkdir(hugoVersionManagerPath, os.ModeDir)
	}

	versionManager, err := versionmanager.NewVersionManager(hugoVersionManagerPath)
	if err != nil {
		return nil, err
	}

	command := new(exec.Cmd)
	path, selectedVersion, err := versionManager.GetExecPath(hugoVersion)
	if err != nil {
		return nil, err
	}
	command.Path = path
	fmt.Printf("selected version: %s\n", selectedVersion)
	command.Args = []string{"hugo"}
	command.Args = append(command.Args, args...)

	return command, nil
}

func persistantPostRun(cmd *cobra.Command, args []string) {
	for _, flag := range wrappedFlags {
		stringFlag := "-" + flag.Name
		if len(flag.Name) > 1 {
			stringFlag = "-" + stringFlag
		}
		wrappedArgs = append(wrappedArgs, stringFlag)
		if flag.NoOptDefVal == "" {
			wrappedArgs = append(wrappedArgs, flag.Value.String())
		}
	}
	command, err := getHugoCommand(wrappedArgs)
	if err != nil {
		fmt.Println(err)
		return
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}
