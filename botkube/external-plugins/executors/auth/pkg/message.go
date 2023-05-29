package pkg

import (
	"fmt"
	"github.com/kubeshop/botkube/pkg/api"
	//"ChatOps/botkube/external-plugins/executors/auth/"
)

const (
	PluginEcho    = "echo"
	PluginAuth    = "auth"
	PluginKubectl = "kubectl"
	Description   = "Echo is an example Botkube executor plugin used during e2e tests. It's not meant for production usage."
)

func cmdPrefix(plugin string) string {
	return fmt.Sprintf("%s %s", api.MessageBotNamePlaceholder, plugin)
}

func InitMessage() []api.Section {
	return []api.Section{
		{
			Base: api.Base{
				Header: "Choose command",
			},
			Selects: api.Selects{
				ID: "select-command",
				Items: []api.Select{
					{
						Name:    "echo",
						Command: cmdPrefix(PluginAuth),
						OptionGroups: []api.OptionGroup{
							{
								Name: cmdPrefix("selects command"),
								Options: []api.OptionItem{
									{Name: "echo", Value: "echo"},
									{Name: "kubectl", Value: "kubectl"},
								},
							},
						},
						// MUST be defined also under OptionGroups.Options slice.
						InitialOption: &api.OptionItem{
							Name: "echo", Value: "echo",
						},
					},
				},
			},
		},
	}
}

func VerbsMessage(cmd string) []api.Section {
	var resOut []api.Section
	initOut := InitMessage()
	if cmd == PluginEcho {
		resOut = repDefaultItems(initOut, PluginEcho)

	}
	if cmd == PluginKubectl {
		resOut = repDefaultItems(initOut, PluginKubectl)
	}

	return expandItems(cmd, resOut)
}

func ApprovalMessage(out []api.Section) []api.Section {
	newItems := api.Select{
		Name:    "approval user",
		Command: cmdPrefix("auth echo --user"),
		OptionGroups: []api.OptionGroup{
			{
				Name: cmdPrefix(""),
				Options: []api.OptionItem{
					{Name: "npos", Value: "npos"},
					{Name: "silence", Value: "silence"},
				},
			},
		},
		// MUST be defined also under OptionGroups.Options slice.
		InitialOption: &api.OptionItem{
			Name: "npos", Value: "npos",
		},
	}

	out[0].Selects.Items = append(out[0].Selects.Items, newItems)

	return out
}

func repDefaultItems(out []api.Section, cmd string) []api.Section {
	out[0].Selects.Items[0].InitialOption.Name = cmd
	out[0].Selects.Items[0].InitialOption.Value = cmd
	return out
}

func expandItems(cmd string, out []api.Section) []api.Section {
	var newItems api.Select
	if cmd == "echo" {
		newItems = api.Select{
			Name:    "echo verbs",
			Command: cmdPrefix("auth echo"),
			OptionGroups: []api.OptionGroup{
				{
					Name: cmdPrefix("selects echo verb"),
					Options: []api.OptionItem{
						{Name: "hallo", Value: "hallo"},
						{Name: "world", Value: "world"},
					},
				},
			},
			// MUST be defined also under OptionGroups.Options slice.
			InitialOption: &api.OptionItem{
				Name: "hallo", Value: "hallo",
			},
		}
	}

	if cmd == "kubectl" {
		newItems = api.Select{
			Name:    "kubectl verbs",
			Command: cmdPrefix("auth kubectl"),
			OptionGroups: []api.OptionGroup{
				{
					Name: cmdPrefix("selects kubectl verb"),
					Options: []api.OptionItem{
						{Name: "get", Value: "get"},
						{Name: "delete", Value: "delete"},
					},
				},
			},
			// MUST be defined also under OptionGroups.Options slice.
			InitialOption: &api.OptionItem{
				Name: "get", Value: "get",
			},
		}
	}
	out[0].Selects.Items = append(out[0].Selects.Items, newItems)

	return out
}
