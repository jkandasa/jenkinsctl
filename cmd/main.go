package main

import (
	"github.com/jkandasa/jenkinsctl/cmd/command"
	"github.com/jkandasa/jenkinsctl/pkg/model"
)

func main() {
	streams := model.NewStdStreams()
	command.Execute(streams)
}
