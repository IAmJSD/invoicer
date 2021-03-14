package main

import (
	"bytes"
	"context"
	"github.com/Masterminds/sprig"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/mandolyte/mdtopdf"
	"invoicer/db"
	"io/ioutil"
	"os"
	"strconv"
	"text/template"
)

// Get the value. Returns nil if it doesn't exist.
func getValue(jobId int, key string) *string {
	value, err := db.GetJobVariable(context.TODO(), jobId, key)
	if err == pgx.ErrNoRows {
		return nil
	} else if err == nil {
		return &value
	}
	panic(err)
}

// Return the bytes of a created PDF or a error.
func pdfGenerator(jobId int, content string) ([]byte, error) {
	// Process the templates in the markdown.
	m := template.FuncMap{
		"PersistentCounter": func(key string) int {
			val := getValue(jobId, key)
			if val == nil {
				// Insert into the database and return 0.
				if err := db.SetJobVariable(context.TODO(), jobId, key, "0"); err != nil {
					panic(err)
				}
				return 0
			}
			i, err := strconv.Atoi(*val)
			if err != nil {
				panic(err)
			}
			i++
			if err := db.SetJobVariable(context.TODO(), jobId, key, strconv.Itoa(i)); err != nil {
				panic(err)
			}
			return i
		},
	}
	tmpl, err := template.New("markdown").Funcs(m).Funcs(sprig.TxtFuncMap()).Parse(content)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, nil); err != nil {
		return nil, err
	}

	// From this, create the PDF.
	fp := uuid.NewString() + ".pdf"
	renderer := mdtopdf.NewPdfRenderer("", "", fp, "")
	if err := renderer.Process(buf.Bytes()); err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	_ = os.Remove(fp)
	return b, nil
}
