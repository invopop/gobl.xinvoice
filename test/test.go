// Package test provides tools for testing the library
package test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl"
	xinvoice "github.com/invopop/gobl.xinvoice"
	"github.com/invopop/gobl/bill"
)

// NewDocumentFrom creates a xinvoice Document from a GOBL file in the `test/data` folder
func NewDocumentFrom(name string) (*xinvoice.Document, error) {
	env, err := LoadTestEnvelope(name)
	if err != nil {
		return nil, err
	}

	return xinvoice.NewDocument(env)
}

// LoadTestInvoice returns a GOBL Invoice from a file in the `test/data` folder
func LoadTestInvoice(name string) (*bill.Invoice, error) {
	env, err := LoadTestEnvelope(name)
	if err != nil {
		return nil, err
	}

	return env.Extract().(*bill.Invoice), nil
}

// LoadTestEnvelope returns a GOBL Envelope from a file in the `test/data` folder
func LoadTestEnvelope(name string) (*gobl.Envelope, error) {
	src, _ := os.Open(filepath.Join(GetDataPath(), name))

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(buf.Bytes(), env); err != nil {
		return nil, err
	}

	return env, nil
}

// GenerateXInvoiceFrom returns a XInvoice Document from a GOBL Invoice
func GenerateXInvoiceFrom(inv *bill.Invoice) (*xinvoice.Document, error) {
	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, err
	}

	return xinvoice.NewDocument(env)
}

// LoadOutputFile returns byte data from a file in the `test/data/out` folder
func LoadOutputFile(name string) ([]byte, error) {
	src, _ := os.Open(filepath.Join(GetOutPath(), name))

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// LoadSchemaFile returns byte data from a file in the `test/data/schema` folder
func LoadSchemaFile(name string) ([]byte, error) {
	src, _ := os.Open(filepath.Join(GetSchemaPath(), name))

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetSchemaPath returns the path to the `test/data/schema` folder
func GetSchemaPath() string {
	return filepath.Join(GetDataPath(), "schema")
}

// GetOutPath returns the path to the `test/data/out` folder
func GetOutPath() string {
	return filepath.Join(GetDataPath(), "out")
}

// GetDataPath returns the path to the `test/data` folder
func GetDataPath() string {
	return filepath.Join(GetTestPath(), "data")
}

// GetTestPath returns the path to the `test` folder
func GetTestPath() string {
	return filepath.Join(getRootFolder(), "test")
}

func getRootFolder() string {
	cwd, _ := os.Getwd()

	for !isRootFolder(cwd) {
		cwd = removeLastEntry(cwd)
	}

	return cwd
}

func isRootFolder(dir string) bool {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}

	return false
}

func removeLastEntry(dir string) string {
	lastEntry := "/" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}
