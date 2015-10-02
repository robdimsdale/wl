package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	urlLongFlag = "url"
)

var (
	// Flags
	url string

	// Commands
	cmdWebhooks = &cobra.Command{
		Use:   "webhooks",
		Short: "gets all webhooks",
		Long: `webhooks gets the user's webhooks.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// Currently sending completed=false is the same as not sending completed
			// Checking for whether the flag has changed protects us from potentially
			// breaking changes, i.e. if the webhooks endpoint decides to return all webhooks,
			// not just non-completed ones.

			if listID == 0 {
				renderOutput(newClient(cmd).Webhooks())
			} else {
				renderOutput(newClient(cmd).WebhooksForListID(listID))
			}
		},
	}

	cmdWebhook = &cobra.Command{
		Use:   "webhook <webhook-id>",
		Short: "gets the webhook for the provided webhook id",
		Long: `webhook gets a webhook specified by <webhook-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(webhook(cmd, args))
		},
	}

	cmdCreateWebhook = &cobra.Command{
		Use:   "create-webhook",
		Short: "creates a webhook with the specified args",
		Long: `create-webhook creates a webhook with the specified args
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(newClient(cmd).CreateWebhook(
				listID,
				url,
				"generic",
				"",
			))
		},
	}

	cmdDeleteWebhook = &cobra.Command{
		Use:   "delete-webhook <webhook-id>",
		Short: "deletes the webhook for the provided webhook id",
		Long: `delete-webhook deletes the webhook specified by <webhook-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			webhook, err := webhook(cmd, args)
			if err != nil {
				fmt.Printf("error getting webhook: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).DeleteWebhook(webhook)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("webhook %d deleted successfully\n", webhook.ID)
		},
	}
)

func init() {
	cmdWebhooks.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")

	cmdCreateWebhook.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "id of list to which webhook will belong")
	cmdCreateWebhook.Flags().StringVar(&url, urlLongFlag, "", "url of webhook")
}

func webhook(cmd *cobra.Command, args []string) (wl.Webhook, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing webhookID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).Webhook(id)
}
