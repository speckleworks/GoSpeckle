package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GitCommit string
	OsArch    string
	BuildDate string
	Version   string
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information about this build",
	Long:  `All software has versions. This is GoSpeckle's`,
	Run: func(cmd *cobra.Command, args []string) {
		printLogo()
		fmt.Println("")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Built:      %s\n", BuildDate)
		fmt.Printf("OS/Arch:    %s\n", OsArch)
	},
}
