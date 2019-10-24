package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var filename string

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVarP(&filename, "filename", "f", "", "path to the file to read speckle resources from")
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply the contents of a file to a speckle server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		requestObjects, err := parseResourceFile(filename)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// fmt.Println(requestObjects)
		requestObjects.MakeRequests(ctx, &speckleClient)
	},
}
