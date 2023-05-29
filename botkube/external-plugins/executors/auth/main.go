package main

import (
	"ChatOps/botkube/external-plugins/executors/auth/pkg"
	"context"
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
		Description: pkg.Description,
		JSONSchema:  jsonSchema(),
	}, nil
}

// Execute returns a given command as response.
func (e *EchoExecutor) Execute(_ context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	var cfg Config
	var sectionOut []api.Section
	err := pluginx.MergeExecutorConfigs(in.Configs, &cfg)
	if err != nil {
		return executor.ExecuteOutput{}, fmt.Errorf("while merging input configuration: %w", err)
	}

	if in.Command == pkg.PluginAuth {
		return executor.ExecuteOutput{
			//Data: data,
			Message: api.Message{
				Type:     api.DefaultMessage,
				Sections: pkg.InitMessage(),
			},
		}, nil
	}

	commands := strings.Split(in.Command, " ")
	if len(commands) > 1 {
		if commands[1] == pkg.PluginEcho {
			sectionOut = pkg.VerbsMessage(pkg.PluginEcho)
			if len(commands) > 2 {
				sectionOut = pkg.ApprovalMessage(sectionOut)
			}
		}
		if commands[1] == pkg.PluginKubectl {
			sectionOut = pkg.VerbsMessage(pkg.PluginKubectl)
			if len(commands) > 2 {
				sectionOut = pkg.ApprovalMessage(sectionOut)
			}
		}
	}

	return executor.ExecuteOutput{
		//Data: data,
		Message: api.Message{
			Type:     api.DefaultMessage,
			Sections: sectionOut,
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
		pkg.PluginAuth: &executor.Plugin{
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
		}`, pkg.Description),
	}
}
