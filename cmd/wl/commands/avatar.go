package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	avatarSizeLongFlag = "size"

	avatarFallbackLongFlag = "fallback"
)

var (
	// Flags
	avatarSize     int
	avatarFallback bool

	// Commands
	cmdAvatarURL = &cobra.Command{
		Use:   "avatar-url <user-id> [flags]",
		Short: "gets a user's avatar",
		Long: `avatar-url gets the URL of a user's avatar>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("incorrect number of arguments provided\n\n")
				cmd.Usage()
				os.Exit(2)
			}

			userIDInt, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("error parsing userID: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}
			userID := uint(userIDInt)

			renderOutput(newClient(cmd).AvatarURL(
				userID,
				avatarSize,
				avatarFallback,
			))
		},
	}
)

func init() {
	cmdAvatarURL.Flags().IntVar(&avatarSize, avatarSizeLongFlag, 0, "specify size")
	cmdAvatarURL.Flags().BoolVar(&avatarFallback, avatarFallbackLongFlag, true, "use fallback avatar if user's avatar is not present")
}
