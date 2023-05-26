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

	return executor.ExecuteOutput{
		//Data: data,
		Message: api.Message{
			Sections: []api.Section{
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
					Base: api.Base{
						Header:      "This a stupid test for multiselect",
						Description: "Number 1, this is a xxx multiselect!!",
					},
					MultiSelect: api.MultiSelect{
						Name: "multi select t1",
						Description: api.Body{
							CodeBlock: "c1",
							Plaintext: "p1",
						},
						Command: "echo",
						Options: []api.OptionItem{
							{
								Name:  "on1",
								Value: "ov1",
							},
							{
								Name:  "on2",
								Value: "ov2",
							},
						},
						InitialOptions: []api.OptionItem{
							{
								Name:  "on11",
								Value: "ov11",
							},
							{
								Name:  "on22",
								Value: "ov22",
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
							{
								Type:    api.ExternalSelect,
								Name:    "si2",
								Command: "echo",
							},
						},
					},
				},
				{
					Base: api.Base{
						Header:      "This a stupid test for PlaintextInputs",
						Description: "Number 5, this is a xxxx PlaintextInputs!",
					},
					PlaintextInputs: []api.LabelInput{
						{
							Command:          "echo",
							Text:             "plt1",
							Placeholder:      "plp1",
							DispatchedAction: api.NoDispatchInputAction,
						},
						{
							Command:          "echo",
							Text:             "plt2",
							Placeholder:      "plp2",
							DispatchedAction: api.DispatchInputActionOnEnter,
						},
						{
							Command:          "echo",
							Text:             "plt3",
							Placeholder:      "plp3",
							DispatchedAction: api.DispatchInputActionOnCharacter,
						},
					},
				},
				{
					Base: api.Base{
						Header:      "This a stupid test for TextFields",
						Description: "Number 1, this is a xxx TextFields!",
					},
					TextFields: []api.TextField{
						{
							Key:   "ttk1",
							Value: "ttv1",
						},
						{
							Key:   "ttk2",
							Value: "ttv2",
						},
					},
				},
				{
					Base: api.Base{
						Header:      "This a stupid test for BulletLists",
						Description: "Number 4, this is a xxx BulletLists!",
					},
					BulletLists: []api.BulletList{
						{
							Title: "bbt1",
							Items: []string{
								"bbi1",
								"bbi2",
							},
						},
					},
				},
				{
					Base: api.Base{
						Header:      "This a stupid test for Context",
						Description: "Number !, this is a xxx Context!",
					},
					Context: []api.ContextItem{
						{
							Text: "cct1",
						},
						{
							Text: "cct2",
						},
					},
				},
			},
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
