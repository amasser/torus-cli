package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/arigatomachine/cli/api"
	"github.com/arigatomachine/cli/apitypes"
	"github.com/arigatomachine/cli/config"
	"github.com/urfave/cli"
)

func init() {
	signup := cli.Command{
		Name:      "signup",
		Usage:     "Create a new Arigato account which, while in alpha, requires an invite code",
		ArgsUsage: "[email] [code]",
		Category:  "ACCOUNT",
		Action:    Chain(EnsureDaemon, signup),
	}
	Cmds = append(Cmds, signup)
}

func signup(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) > 0 && len(args) != 2 {
		var text string
		if len(args) > 2 {
			text = "Too many arguments supplied.\n\n"
		} else {
			text = "Too few arguments supplied.\n\n"
		}
		text += usageString(ctx)
		return cli.NewExitError(text, -1)
	}

	name, err := FullNamePrompt()
	if err != nil {
		return err
	}

	username, err := UsernamePrompt()
	if err != nil {
		return err
	}

	defaultEmail := ""
	defaultInvite := ""
	if len(args) == 2 {
		defaultEmail = args[0]
		defaultInvite = args[1]
	}

	email, err := EmailPrompt(defaultEmail)
	if err != nil {
		return err
	}

	var inviteCode string
	if os.Getenv("AG_DEBUG") == "" {
		inviteCode, err = InviteCodePrompt(defaultInvite)
		if err != nil {
			return err
		}
	}

	password, err := PasswordPrompt(true)
	if err != nil {
		return err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg)

	var output api.ProgressFunc
	output = func(event *api.Event, err error) {
		if event != nil {
			fmt.Println(event.Message)
		}
	}

	signup := apitypes.Signup{
		Name:       name,
		Username:   username,
		Passphrase: password,
		Email:      email,
		InviteCode: inviteCode,
	}

	c := context.Background()

	fmt.Println("")
	user, err := client.Users.Signup(c, &signup, &output)
	if err != nil {
		if strings.Contains(err.Error(), "resource exists") {
			return cli.NewExitError("Email address in use, please try again.", -1)
		}
		return cli.NewExitError("Signup failed, please try again.", -1)
	}

	// Log the user in
	err = performLogin(c, client, user.Body.Email, password)
	if err != nil {
		return err
	}

	// Generate keypairs, look up the user's org
	err = generateKeypairsForOrg(ctx, c, client, nil, true)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("Your account has been created!")
	return nil
}