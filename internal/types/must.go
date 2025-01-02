package types

import (
	"fmt"
	"os"
)

type CliMust struct {
	exitCode int
}

func NewCliMust() *CliMust {
	return &CliMust{
		exitCode: 0,
	}
}

func (c *CliMust) Must(data any, err error) any {
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(c.exitCode)
	}

	return data
}

func (c *CliMust) MustWithMessage(data any, err error, message string) any {
	if err != nil {
		fmt.Printf("Error: %v", message)
	}

	return data
}
