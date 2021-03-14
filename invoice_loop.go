package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	"invoicer/db"
	"strings"
	"time"
)

func invoiceLoop() {
	for {
		jobs, err := db.GetPendingInvoiceJobs(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, v := range jobs {
			// Get the from e-mail.
			email, err := db.SelectOneEmailAddress(context.TODO(), v.FromEmail, v.FromName)
			if err != nil {
				// Record doesn't exist at the moment.
				continue
			}

			// Format who this is to properly.
			var to string
			if v.ToName == "" {
				to = "<" + v.ToName + ">"
			} else {
				to = v.ToName + " <" + v.ToEmail + ">"
			}

			// Defines who this email is from.
			var from string
			if v.FromName == "" {
				from = "<" + v.FromEmail + ">"
			} else {
				from = v.FromName + " <" + v.FromEmail + ">"
			}

			// Defines the writer.
			buf := &bytes.Buffer{}
			boundary := strings.ReplaceAll(uuid.NewString(), "-", "")

			// Write the content start.
			buf.WriteString("From: " + from + "\r\nTo: " + to + "\r\nSubject: " + v.Subject + "\r\nContent-Type: multipart/mixed; " +
				"boundary=" + boundary + "\r\n\r\n--" + boundary + "\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n")

			// Write the content.
			buf.WriteString(strings.TrimRight(v.EmailContent, "\r\n") + "\r\n")

			// Write the attachment header.
			buf.WriteString("\r\n--" + boundary + "\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n" +
				"Content-Transfer-Encoding: base64\r\nContent-Disposition: attachment;filename=\"invoice.pdf\"\r\n\r\n")

			// Get the PDF.
			pdfBytes, err := pdfGenerator(v.JobID, v.Markdown)
			if err != nil {
				if err = db.InsertFailure(context.TODO(), err); err != nil {
					panic(err)
				}
				continue
			}

			// Encode and write the invoice.
			buf.WriteString(base64.StdEncoding.EncodeToString(pdfBytes))

			// Send the e-mail.
			if err = email.Send([]string{v.ToEmail}, buf.Bytes()); err != nil {
				if err = db.InsertFailure(context.TODO(), err); err != nil {
					panic(err)
				}
				continue
			}

			// Mark the job as complete so it doesn't run again.
			if err = v.MarkAsComplete(context.TODO()); err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second * 5)
	}
}
