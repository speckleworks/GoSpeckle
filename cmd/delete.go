package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteAPIClienttCmd.Flags().StringVarP(&id, "id", "i", "", "the ID of the resource to delete")
	deleteCommenttCmd.Flags().StringVarP(&id, "id", "i", "", "the ID of the resource to delete")
	deleteObjecttCmd.Flags().StringVarP(&id, "id", "i", "", "the ID of the resource to delete")
	deleteProjecttCmd.Flags().StringVarP(&id, "id", "i", "", "the ID of the resource to delete")
	deleteStreamtCmd.Flags().StringVarP(&id, "id", "i", "", "the ID of the resource to delete")

	deleteAPIClienttCmd.MarkFlagRequired("id")
	deleteCommenttCmd.MarkFlagRequired("id")
	deleteObjecttCmd.MarkFlagRequired("id")
	deleteProjecttCmd.MarkFlagRequired("id")
	deleteStreamtCmd.MarkFlagRequired("id")

	deleteCmd.AddCommand(deleteAPIClienttCmd)
	deleteCmd.AddCommand(deleteCommenttCmd)
	deleteCmd.AddCommand(deleteObjecttCmd)
	deleteCmd.AddCommand(deleteProjecttCmd)
	deleteCmd.AddCommand(deleteStreamtCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Speckle resource",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var deleteAPIClienttCmd = &cobra.Command{
	Use:   "client",
	Short: "Delete a Speckle client",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := speckleClient.APIClient.Delete(ctx, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Succesfully deleted client with ID: %s\n", id)
	},
}

var deleteCommenttCmd = &cobra.Command{
	Use:   "comment",
	Short: "Delete a Speckle comment",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := speckleClient.Comment.Delete(ctx, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Succesfully deleted comment with ID: %s\n", id)
	},
}

var deleteObjecttCmd = &cobra.Command{
	Use:   "object",
	Short: "Delete a Speckle object",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := speckleClient.Object.Delete(ctx, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Succesfully deleted object with ID: %s\n", id)
	},
}

var deleteProjecttCmd = &cobra.Command{
	Use:   "project",
	Short: "Delete a Speckle project",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := speckleClient.Project.Delete(ctx, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Succesfully deleted project with ID: %s\n", id)
	},
}

var deleteStreamtCmd = &cobra.Command{
	Use:   "stream",
	Short: "Delete a Speckle stream",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := speckleClient.Stream.Delete(ctx, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Succesfully deleted stream with ID: %s\n", id)
	},
}
