package main

import (
	"github.com/jkandasa/jenkinsctl/cmd/command"
	types "github.com/jkandasa/jenkinsctl/pkg/types"
)

func main() {
	streams := types.NewStdStreams()
	command.Execute(streams)
}
