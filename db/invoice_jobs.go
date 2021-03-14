package db

import (
	"context"
	"errors"
	"github.com/jackc/pgtype"
)

// InvoiceJob defines the invoice job information.
type InvoiceJob struct {
	// Subject is used to define the e-mail subject.
	Subject string `json:"subject"`

	// FromEmail is used to define the e-mail address this is from.
	FromEmail string `json:"from_email"`

	// FromName is used to define the name this is from.
	FromName string `json:"from_name"`

	// JobID is used to define the ID of the job.
	JobID int `json:"job_id"`

	// Date defines the date which the job triggers.
	Date uint8 `json:"date"`

	// MonthsBetweenTrigger defines how many months between each trigger.
	MonthsBetweenTrigger uint8 `json:"months_between_trigger"`

	// ToEmail is used to define the e-mail address this is to.
	ToEmail string `json:"to_email"`

	// ToName is used to define the name this is to.
	ToName string `json:"to_name"`

	// Markdown is used to define the markdown content used.
	Markdown string `json:"markdown"`

	// LastTriggered defines when the job was last triggered.
	LastTriggered pgtype.Timestamptz `json:"last_triggered"`

	// EmailContent is used to define the e-mails content.
	EmailContent string `json:"email_content"`
}

// ZeroMonthsInvalid is used when a 0 month span is set for the inserted job.
var ZeroMonthsInvalid = errors.New("span between jobs cannot be zero months")

// MuteBeBefore28th is used when a job should be before the 28th.
var MuteBeBefore28th = errors.New("job must be before the 28th")

// NoBlankSubject ensures that the subject cannot be blank.
var NoBlankSubject = errors.New("subject cannot be blank")

// InsertInvoiceJob is used to insert a invoice job into the database.
// Note that the job ID is ignored in the specified object.
func InsertInvoiceJob(ctx context.Context, job *InvoiceJob) error {
	if job.MonthsBetweenTrigger == 0 {
		return ZeroMonthsInvalid
	}
	if job.Date > 28 {
		return MuteBeBefore28th
	}
	if job.Subject == "" {
		return NoBlankSubject
	}
	_, err := conn.Exec(ctx, "INSERT INTO invoice_jobs (subject, from_email, "+
		"from_name, date, months_between_trigger, to_email, to_name, markdown, "+
		"email_content) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", job.Subject,
		job.FromEmail, job.FromName, job.Date, job.MonthsBetweenTrigger, job.ToEmail,
		job.ToName, job.Markdown, job.EmailContent)
	return err
}

// SelectAllInvoiceJobs is used to get all of the invoice jobs.
func SelectAllInvoiceJobs(ctx context.Context) ([]*InvoiceJob, error) {
	ptr, err := conn.Query(ctx, "SELECT from_email, from_name, job_id, date, "+
		"months_between_trigger, to_email, to_name, markdown, last_triggered, subject, email_content FROM invoice_jobs")
	if err != nil {
		return nil, err
	}
	a := make([]*InvoiceJob, 0)
	for ptr.Next() {
		x := InvoiceJob{}
		err = ptr.Scan(&x.FromEmail, &x.FromName, &x.JobID, &x.Date, &x.MonthsBetweenTrigger,
			&x.ToEmail, &x.ToName, &x.Markdown, &x.LastTriggered, &x.Subject, &x.EmailContent)
		if err != nil {
			return nil, err
		}
		a = append(a, &x)
	}
	return a, nil
}

// GetPendingInvoiceJobs is used to get the invoice jobs which are pending.
func GetPendingInvoiceJobs(ctx context.Context) ([]*InvoiceJob, error) {
	ptr, err := conn.Query(ctx, "SELECT from_email, from_name, job_id, date, "+
		"months_between_trigger, to_email, to_name, markdown, last_triggered, subject, email_content FROM invoice_jobs "+
		"WHERE NOW() >= COALESCE(last_triggered, CASE WHEN EXTRACT(DAY FROM NOW()) = date THEN "+
		"TIMESTAMP WITH TIME ZONE '1970-01-01 00:00:00+00' ELSE NOW() END) + (months_between_trigger * INTERVAL '1 month')")
	if err != nil {
		return nil, err
	}
	a := make([]*InvoiceJob, 0)
	for ptr.Next() {
		x := InvoiceJob{}
		err = ptr.Scan(&x.FromEmail, &x.FromName, &x.JobID, &x.Date, &x.MonthsBetweenTrigger,
			&x.ToEmail, &x.ToName, &x.Markdown, &x.LastTriggered, &x.Subject, &x.EmailContent)
		if err != nil {
			return nil, err
		}
		a = append(a, &x)
	}
	return a, nil
}

// MarkAsComplete is mark a invoice job as complete.
func (i *InvoiceJob) MarkAsComplete(ctx context.Context) error {
	_, err := conn.Exec(ctx, "UPDATE invoice_jobs SET last_triggered = NOW() WHERE job_id = $1", i.JobID)
	return err
}

// DeleteInvoiceJob is used to delete a invoice job from the database.
func DeleteInvoiceJob(ctx context.Context, JobID int) error {
	_, err := conn.Exec(ctx, "DELETE FROM invoice_jobs WHERE job_id = $1", JobID)
	return err
}

// UpdateInvoiceInformation is used to update some of the information in an invoice.
func UpdateInvoiceInformation(ctx context.Context, JobID int, EmailContent, Markdown, Subject string) error {
	_, err := conn.Exec(ctx, "UPDATE invoice_jobs SET email_content = $1, markdown = $2, subject = $3 WHERE job_id = $4", EmailContent, Markdown, Subject, JobID)
	return err
}
