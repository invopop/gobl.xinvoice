package main

import (
	"fmt"
	"io"

	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/spf13/cobra"
)

type convertOpts struct {
	*rootOpts
	format string
}

func convert(o *rootOpts) *convertOpts {
	return &convertOpts{rootOpts: o}
}

func (c *convertOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert <infile> <outfile>",
		Short: "Convert a GOBL JSON into a XRechnung & Factur-X XML document",
		RunE:  c.runE,
	}

	f := cmd.Flags()
	f.StringVar(&c.format, "format", "", "The format to convert to, either 'xrechnung', 'zugferd' or 'facturx'")

	return cmd
}

func (c *convertOpts) runE(cmd *cobra.Command, args []string) error {
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

	doc, err := xinvoice.Convert(inData, c.format)
	if err != nil {
		return fmt.Errorf("building document: %w", err)
	}

	if _, err = out.Write(doc); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
