package main

import (
	"context"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
	"github.com/kubeshop/botkube/pkg/bot/interactive"
)

// EchoExecutor implements the Botkube executor plugin interface.
type EchoExecutor struct{}

func (*EchoExecutor) Metadata(context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     "1.0.0",
		Description: "Echo sends back the command that was specified.",
		JSONSchema: api.JSONSchema{
			Value: heredoc.Doc(`{
		   "$schema": "http://json-schema.org/draft-04/schema#",
		   "title": "botkube/echo",
		   "description": "example echo plugin",
		   "type": "object",
		   "properties": {
			 "formatOptions": {
			   "description": "options to format echoed string",
			   "type": "array",
			   "items": {
				 "type": "string",
				 "enum": [
				   "bold",
				   "italic"
				 ]
			   }
			 }
		   },
		   "additionalProperties": false,
		   "required": []
		 }`),
		},
	}, nil
}

func (*EchoExecutor) Execute(_ context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	return executor.ExecuteOutput{
		Data: fmt.Sprintf("Echo: %s", in.Command),
	}, nil
}

func (*EchoExecutor) Help(_ context.Context) (interactive.Message, error) {
	return interactive.Message{
		Base: interactive.Base{
			Body: interactive.Body{
				CodeBlock: "Echo prints out given input string.",
			},
		},
	}, nil
}

func main() {
	executor.Serve(map[string]plugin.Plugin{
		"echo": &executor.Plugin{
			Executor: &EchoExecutor{},
		},
	})
}
