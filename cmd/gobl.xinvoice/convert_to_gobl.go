package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/spf13/cobra"
)

// type convertOpts struct {
// 	*rootOpts
// }

func convertToGobl(o *rootOpts) *convertOpts {
	return &convertOpts{rootOpts: o}
}

func (c *convertOpts) cmdNew() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert_to_gobl <infile> <outfile>",
		Short: "Convert a XRechnung to GOBL",
		RunE:  c.runENew,
	}

	return cmd
}

func (c *convertOpts) runENew(cmd *cobra.Command, args []string) error {
	// ctx := commandContext(cmd)

	if len(args) == 0 || len(args) > 2 {
		return fmt.Errorf("expected one or two arguments, the command usage is `gobl.cfdi convert <infile> [outfile]`")
	}

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	out, err := c.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	inData, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	doc := new(xinvoice.XMLDoc)
	if err := xml.Unmarshal(inData, doc); err != nil {
		return fmt.Errorf("parsing input document: %w", err)
	}

	env, err := xinvoice.NewDocumentGOBL(doc)
	if err != nil {
		return fmt.Errorf("building GOBL envelope: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(env); err != nil {
		return fmt.Errorf("generating XRechnung and Factur-X xml: %w", err)
	}

	return nil
}
