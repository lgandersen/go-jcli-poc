package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ImageCommand() *cobra.Command {
	containerCmd := &cobra.Command{
		Use:                   "image",
		Short:                 "Manage images",
		Long:                  `Manage images`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The image main command have been executed")
		},
	}

	containerCmd.AddCommand(ImageBuildCommand())
	containerCmd.AddCommand(ImageRemoveCommand())
	containerCmd.AddCommand(ImageListCommand())
	return containerCmd
}

func ImageBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "build",
		Short:                 "Build an image from a Dockerfile",
		Long:                  `Build an image from a Dockerfile`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The image build command have been executed")
		},
	}
	return cmd
}

func ImageRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rm",
		Short:                 "Remove one or more images",
		Long:                  `Remove one or more images`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The image rm command have been executed")
		},
	}
	return cmd
}

func ImageListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List images",
		Long:                  `List images`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The image ls command have been executed")
		},
	}
	return cmd
}
