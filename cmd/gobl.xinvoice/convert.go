package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/invopop/gobl"
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
	f.StringVar(&c.format, "format", "", "The format to convert to, either 'xrechnung-cii', 'xrechnung-ubl', 'zugferd' or 'facturx'")

	return cmd
}

func (c *convertOpts) runE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || len(args) > 2 {
		return fmt.Errorf("expected one or two arguments, the command usage is `gobl.xinvoice convert <infile> [outfile]`")
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

	in, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	j := json.Valid(in)
	var o []byte
	if j {
		env := new(gobl.Envelope)
		if err := json.Unmarshal(in, env); err != nil {
			return fmt.Errorf("parsing input as GOBL Envelope: %w", err)
		}
		switch c.format {
		case "xrechnung-cii":
			o, err = xinvoice.ConvertToXRechnungCII(env)
		case "xrechnung-ubl":
			o, err = xinvoice.ConvertToXRechnungUBL(env)
		case "facturx":
			o, err = xinvoice.ConvertToFacturX(env)
		case "zugferd":
			o, err = xinvoice.ConvertToZUGFeRD(env)
		default:
			return fmt.Errorf("unknown or missing format: %s", c.format)
		}
		if err != nil {
			return fmt.Errorf("generating XML: %w", err)
		}
	} else {
		// Assume XML if not JSON
		env, err := xinvoice.ConvertToGOBL(in)
		if err != nil {
			return fmt.Errorf("converting CII to GOBL: %w", err)
		}

		o, err = json.MarshalIndent(env, "", "  ")
		if err != nil {
			return fmt.Errorf("generating JSON output: %w", err)
		}
	}

	if _, err = out.Write(o); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
