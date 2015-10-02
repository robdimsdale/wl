package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robdimsdale/wl"
	"github.com/spf13/cobra"
)

const (
	userIDLongFlag       = "userID"
	emailAddressLongFlag = "emailAddress"
	mutedLongFlag        = "muted"
)

var (
	// Flags
	userID       uint
	emailAddress string
	muted        bool

	// Commands
	cmdMemberships = &cobra.Command{
		Use:   "memberships",
		Short: "gets all memberships",
		Long: `memberships gets the user's memberships.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// Currently sending completed=false is the same as not sending completed
			// Checking for whether the flag has changed protects us from potentially
			// breaking changes, i.e. if the memberships endpoint decides to return all memberships,
			// not just non-completed ones.

			if listID == 0 {
				renderOutput(newClient(cmd).Memberships())
			} else {
				renderOutput(newClient(cmd).MembershipsForListID(listID))
			}
		},
	}

	cmdMembership = &cobra.Command{
		Use:   "membership <membership-id>",
		Short: "gets the membership for the provided membership id",
		Long: `membership gets a membership specified by <membership-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			renderOutput(membership(cmd, args))
		},
	}

	cmdInviteMember = &cobra.Command{
		Use:   "invite-member [flags]",
		Short: "invites the member identified by the provided flags to the specified list",
		Long: `invite-member invites the member specified by the flags to the specified list.
User must be identified by providing either ` + userIDLongFlag + ` or ` + emailAddressLongFlag + `
        `,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed(userIDLongFlag) {
				renderOutput(newClient(cmd).AddMemberToListViaUserID(
					userID,
					listID,
					muted,
				))
			} else if cmd.Flags().Changed(emailAddressLongFlag) {
				renderOutput(newClient(cmd).AddMemberToListViaEmailAddress(
					emailAddress,
					listID,
					muted,
				))
			} else {
				cmd.Usage()
				os.Exit(2)
			}
		},
	}

	cmdAcceptMembership = &cobra.Command{
		Use:   "accept-membership <membership-id>",
		Short: "accepts the membership invite for the provided membership id",
		Long: `accept-member rejects the invite for the membership specified by <membership-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			membership, err := membership(cmd, args)
			if err != nil {
				fmt.Printf("error getting membership: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			renderOutput(newClient(cmd).AcceptMember(membership))
		},
	}

	cmdRemoveMembership = &cobra.Command{
		Use:   "remove-membership <membership-id>",
		Short: "removes membership of list for the provided membership id",
		Long: `remove-membership removes membership specified by <membership-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			membership, err := membership(cmd, args)
			if err != nil {
				fmt.Printf("error getting membership: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).RemoveMemberFromList(membership)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("membership %d removed successfully\n", membership.ID)
		},
	}

	cmdRejectMembership = &cobra.Command{
		Use:   "reject-membership <membership-id>",
		Short: "rejects the membership invite for the provided membership id",
		Long: `reject-membership rejects the invite for the membership specified by <membership-id>
        `,
		Run: func(cmd *cobra.Command, args []string) {
			membership, err := membership(cmd, args)
			if err != nil {
				fmt.Printf("error getting membership: %v\n\n", err)
				cmd.Usage()
				os.Exit(2)
			}

			err = newClient(cmd).RejectInvite(membership)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("invite %d rejected successfully\n", membership.ID)
		},
	}
)

func init() {
	cmdMemberships.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "filter by listID")

	cmdInviteMember.Flags().UintVarP(&listID, listIDLongFlag, listIDShortFlag, 0, "list to which membership will belong")
	cmdInviteMember.Flags().UintVar(&userID, userIDLongFlag, 0, "identify user by userID")
	cmdInviteMember.Flags().StringVar(&emailAddress, emailAddressLongFlag, "", "identify user by emailAddress")
	cmdInviteMember.Flags().BoolVar(&muted, mutedLongFlag, false, "user is muted by default")
}

func membership(cmd *cobra.Command, args []string) (wl.Membership, error) {
	if len(args) != 1 {
		fmt.Printf("incorrect number of arguments provided\n\n")
		cmd.Usage()
		os.Exit(2)
	}

	idInt, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("error parsing membershipID: %v\n\n", err)
		cmd.Usage()
		os.Exit(2)
	}
	id := uint(idInt)

	return newClient(cmd).Membership(id)
}
