package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
	"github.com/kubeshop/botkube/pkg/pluginx"
	"strings"
)

// version is set via ldflags by GoReleaser.
var version = "dev"

const (
	pluginName  = "echo"
	description = "Echo is an example Botkube executor plugin used during e2e tests. It's not meant for production usage."
)

// Config holds executor configuration.
type Config struct {
	ChangeResponseToUpperCase *bool `yaml:"changeResponseToUpperCase,omitempty"`
}

// EchoExecutor implements Botkube executor plugin.
type EchoExecutor struct {
	btnBuilder *api.ButtonBuilder
}

var _ executor.Executor = &EchoExecutor{}

// Metadata returns details about Echo plugin.
func (e *EchoExecutor) Metadata(context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     version,
		Description: description,
		JSONSchema:  jsonSchema(),
	}, nil
}

// Execute returns a given command as response.
func (e *EchoExecutor) Execute(_ context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	var cfg Config
	err := pluginx.MergeExecutorConfigs(in.Configs, &cfg)
	if err != nil {
		return executor.ExecuteOutput{}, fmt.Errorf("while merging input configuration: %w", err)
	}

	data := in.Command
	if strings.Contains(data, "@fail") {
		return executor.ExecuteOutput{}, errors.New("The @fail label was specified. Failing execution.")
	}

	if cfg.ChangeResponseToUpperCase != nil && *cfg.ChangeResponseToUpperCase {
		data = strings.ToUpper(data)
	}

	cmdPrefix := func(cmd string) string {
		return fmt.Sprintf("%s %s %s", api.MessageBotNamePlaceholder, pluginName, cmd)
	}

	secOut := []api.Section{
		{
			Base: api.Base{
				Header:      "This a stupid test for button",
				Description: "Number 1, this is a danger button!",
			},
			Buttons: []api.Button{
				e.btnBuilder.ForCommandWithDescCmd("echo", "echo", api.ButtonStyleDanger),
			},
		},
		{
			Selects: api.Selects{
				ID: "select-id",
				Items: []api.Select{
					{
						Name:    "two-groups",
						Command: cmdPrefix("selects two-groups"),
						OptionGroups: []api.OptionGroup{
							{
								Name: cmdPrefix("selects two-groups/1"),
								Options: []api.OptionItem{
									{Name: "BAR", Value: "BAR"},
									{Name: "BAZ", Value: "BAZ"},
									{Name: "XYZ", Value: "XYZ"},
								},
							},
							{
								Name: cmdPrefix("selects two-groups/2"),
								Options: []api.OptionItem{
									{Name: "123", Value: "123"},
									{Name: "456", Value: "456"},
									{Name: "789", Value: "789"},
								},
							},
						},
						// MUST be defined also under OptionGroups.Options slice.
						InitialOption: &api.OptionItem{
							Name: "789", Value: "789",
						},
					},
				},
			},
		},
		{
			Base: api.Base{
				Header:      "This a stupid test for selects",
				Description: "Number 4, this is a xxx selects!",
			},
			Selects: api.Selects{
				ID: "S1",
				Items: []api.Select{
					{
						Type:    api.StaticSelect,
						Name:    "si1",
						Command: "echo",
						OptionGroups: []api.OptionGroup{
							{
								Name: "sio1",
								Options: []api.OptionItem{
									{
										Name:  "sioon1",
										Value: "sioov1",
									},
									{
										Name:  "sioon2",
										Value: "sioov2",
									},
								},
							},
						},
						InitialOption: &api.OptionItem{
							Name:  "siin1",
							Value: "siin1",
						},
					},
				},
			},
		},
	}

	//out := api.Message{
	//	Sections: []api.Section{
	//		{
	//			Base: api.Base{
	//				Header:      "This a stupid test for button",
	//				Description: "Number 1, this is a danger button!",
	//			},
	//			Buttons: []api.Button{
	//				e.btnBuilder.ForCommandWithDescCmd("echo", "echo", api.ButtonStyleDanger),
	//			},
	//		},
	//	},
	//}
	//outStr, _ := jsoniter.Marshal(out)

	return executor.ExecuteOutput{
		//Data: data,
		Message: api.Message{
			Type:     api.DefaultMessage,
			Sections: secOut,
		},
	}, nil
}

// Help returns help message
func (e *EchoExecutor) Help(context.Context) (api.Message, error) {
	return api.Message{
		Sections: []api.Section{
			{
				Base: api.Base{
					Header:      "This so bad!",
					Description: "114514",
				},
				Context: api.ContextItems{api.ContextItem{Text: "用个嘚的帮助命令，没屌用的"}},
			},
		},
	}, nil
}

func main() {
	executor.Serve(map[string]plugin.Plugin{
		pluginName: &executor.Plugin{
			Executor: &EchoExecutor{},
		},
	})
}

func jsonSchema() api.JSONSchema {
	return api.JSONSchema{
		Value: heredoc.Docf(`{
			"$schema": "http://json-schema.org/draft-07/schema#",
			"title": "botkube/auth",
			"description": "%s",
			"type": "object",
			"properties": {
				"changeResponseToUpperCase": {
					"description": "When changeResponseToUpperCase is true, the echoed string will be in upper case",
					"type": "boolean"
				}
			},
			"required": []
		}`, description),
	}
}
