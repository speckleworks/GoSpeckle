package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// const validArgs []string = []string{ "pod", "node", "service", "replicationcontroller" }

var id string
var ids []string
var search string
var streamID string
var admin bool
var object interface{}
var err error

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the ID of the resource to retrieve from")

	getAccountCmd.Flags().StringVarP(&search, "search", "s", "", "the search string to find accounts")
	getAccountCmd.Flags().BoolVar(&admin, "admin", false, "run this command as administrator")
	getCmd.AddCommand(getAccountCmd)

	getAPIClientCmd.Flags().StringVar(&streamID, "stream", "", "the search streamId of the stream to list clients for")
	getCmd.AddCommand(getAPIClientCmd)

	getProjectCmd.Flags().BoolVar(&admin, "admin", false, "run this command as administrator")
	getCmd.AddCommand(getProjectCmd)

	getStreamCmd.Flags().BoolVar(&admin, "admin", false, "run this command as administrator")
	getCmd.AddCommand(getStreamCmd)

	getObjectCmd.Flags().StringVarP(&search, "search", "s", "", "the search string to find speckle objects")
	getObjectCmd.Flags().StringVar(&streamID, "stream", "", "the search streamId of the stream to list objects for")
	getObjectCmd.Flags().StringSliceVar(&ids, "ids", ids, "a list of object ids to search through (must be in format --ids=\"v1,v2\")")
	getCmd.AddCommand(getObjectCmd)

	getCommentCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "the ID of the resource to retrieve comments from")
	getCmd.AddCommand(getCommentCmd)

	getProjectCommentCmd.MarkFlagRequired("id")
	getStreamCommentCmd.MarkFlagRequired("id")
	getObjectCommentCmd.MarkFlagRequired("id")

	getCommentCmd.AddCommand(getProjectCommentCmd)
	getCommentCmd.AddCommand(getStreamCommentCmd)
	getCommentCmd.AddCommand(getObjectCommentCmd)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get one or a list of speckle resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var getAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Retrieve accounts",
	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			object, _, err = speckleClient.Account.Get(ctx, id)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else if search != "" {
			object, _, err = speckleClient.Account.Search(ctx, search)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else if admin {
			object, _, err = speckleClient.Account.AdminGet(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else {
			object, _, err = speckleClient.Account.Me(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err := printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getAPIClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Retrieve api clients",
	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			object, _, err = speckleClient.APIClient.Get(ctx, id)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if streamID != "" {
			object, _, err = speckleClient.Stream.ListClients(ctx, streamID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			object, _, err = speckleClient.APIClient.List(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		err := printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Retrieve comments from different resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var getProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Retrieve projects",
	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			object, _, err = speckleClient.Project.Get(ctx, id)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if admin {
			object, _, err = speckleClient.Project.AdminGet(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			object, _, err = speckleClient.Project.List(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		err := printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getProjectCommentCmd = &cobra.Command{
	Use:   "project",
	Short: "Retrieve comments from a project",
	Run: func(cmd *cobra.Command, args []string) {

		object, _, err = speckleClient.Project.GetComments(ctx, id)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getStreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Retrieve streams",
	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			object, _, err = speckleClient.Stream.Get(ctx, id)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if admin {
			object, _, err = speckleClient.Stream.AdminGet(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			object, _, err = speckleClient.Stream.List(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err := printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getStreamCommentCmd = &cobra.Command{
	Use:   "stream",
	Short: "Retrieve comments from a stream",
	Run: func(cmd *cobra.Command, args []string) {

		object, _, err = speckleClient.Stream.GetComments(ctx, id)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getObjectCmd = &cobra.Command{
	Use:   "object",
	Short: "Retrieve objects",
	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			object, _, err = speckleClient.Object.Get(ctx, id)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if streamID != "" {
			object, _, err = speckleClient.Stream.ListObjects(ctx, streamID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if search != "" {
			if ids == nil {
				fmt.Println("List of IDs to search within must be provided if using search string")
			}

			object, _, err = speckleClient.Object.Search(ctx, search, ids)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Object ID, Stream ID or Search string must be provided to retrieve objects")
		}

		err := printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getObjectCommentCmd = &cobra.Command{
	Use:   "object",
	Short: "Retrieve comments from a stream",
	Run: func(cmd *cobra.Command, args []string) {

		object, _, err = speckleClient.Object.GetComments(ctx, id)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = printJSON(object)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
