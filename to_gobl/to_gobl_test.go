package to_gobl_test

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl.xinvoice/to_gobl"
	"github.com/invopop/gobl.xinvoice/xinvoice/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateOut = flag.Bool("update", false, "Update the JSON files in the test/data/out directory")

func TestNewDocumentGOBL(t *testing.T) {
	examples, err := test.GetDataGlob("*.xml")
	require.NoError(t, err)

	for _, example := range examples {
		inName := filepath.Base(example)
		outName := strings.Replace(inName, ".xml", ".json", 1)

		t.Run(inName, func(t *testing.T) {
			doc := new(to_gobl.XMLDoc)

			goblEnv, err := to_gobl.NewDocumentGOBL(doc)
			require.NoError(t, err)

			data, err := json.MarshalIndent(goblEnv, "", "  ")
			require.NoError(t, err)

			output, err := test.LoadOutputFile(outName)
			assert.NoError(t, err)

			if *updateOut {
				err = test.SaveOutputFile(outName, data)
				require.NoError(t, err)
			} else {
				assert.JSONEq(t, string(output), string(data), "Output should match the expected JSON. Update with --update flag.")
			}
		})
	}
}

// LoadTestXMLDoc returns a to_gobl.XMLDoc from a file in the test data folder
func LoadTestXMLDoc(name string) (*to_gobl.XMLDoc, error) {
	src, err := os.Open(filepath.Join(test.GetDataPath(), name))
	if err != nil {
		return nil, err
	}
	defer src.Close()

	inData, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	doc := new(to_gobl.XMLDoc)
	if err := xml.Unmarshal(inData, doc); err != nil {
		return nil, err
	}

	return doc, nil
}
