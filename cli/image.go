package cli

import (
	"fmt"
	Openapi "jcli/client"
	"net/url"

	"github.com/spf13/cobra"
)

var image_build_base_url = "ws://localhost:8085/images/build"

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

type ImageBuildOptions struct {
	Context    string
	Dockerfile string
	Tag        string
	Quiet      bool
}

func ImageBuildCommand() *cobra.Command {
	opts := ImageBuildOptions{}
	cmd := &cobra.Command{
		Use:                   "build [OPTIONS] PATH",
		Short:                 "Build an image from a Dockerfile",
		Long:                  `Build an image from a Dockerfile`,
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The image build command have been executed")
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Dockerfile, "file", "f", "Dockerfile", "Name of the Dockerfile (default: 'PATH/Dockerfile')")
	flags.StringVarP(&opts.Tag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	flags.BoolVarP(&opts.Quiet, "quiet", "q", false, "Suppress the build output and print image ID on success (default: false)")
	return cmd
}

func ImageRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rm [OPTIONS] IMAGE [IMAGE...]",
		Short:                 "Remove one or more images",
		Long:                  `Remove one or more images`,
		Args:                  cobra.MinimumNArgs(1),
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
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The image ls command have been executed")
		},
	}
	return cmd
}

func BuildImageAndListenForMessages(options ImageBuildOptions) {
	ws_url, _ := url.Parse(image_build_base_url)
	query := ws_url.Query()
	query.Set("context", options.Context)
	query.Set("dockerfile", options.Dockerfile)
	query.Set("tag", options.Tag)
	query.Set("quiet", fmt.Sprint(options.Quiet))
	ws_url.RawQuery = query.Encode()
	endpoint := ws_url.String()

	done, interrupt, ws := Dial(endpoint)
	go ListenForWSMessages(done, ws)
	BuildImage(NewHTTPClient())
	AwaitDoneOrUserInterrupt(done, interrupt, ws)
}

func BuildImage(client *Openapi.ClientWithResponses) {
}
