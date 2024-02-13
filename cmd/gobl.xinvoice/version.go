package main

import (
	"encoding/json"

	"github.com/invopop/gobl"
	"github.com/spf13/cobra"
)

var versionOutput = struct {
	Version string `json:"version"`
	GOBL    string `json:"gobl"`
	Date    string `json:"date,omitempty"`
}{
	Version: version,
	GOBL:    string(gobl.VERSION),
	Date:    date,
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "\t") // always indent version
			return enc.Encode(versionOutput)
		},
	}
}
