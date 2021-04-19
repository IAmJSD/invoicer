package main

import (
	"context"
	"embed"
	"github.com/gofiber/fiber/v2"
	"invoicer/db"
	"strconv"
	"strings"
)

//go:embed templates/*.html
var templates embed.FS

// Parses the from email into 2 different bits.
func parseEmail(raw string) (string, string, bool) {
	if raw == "" {
		return "", "", false
	}
	i := strings.IndexByte(raw, '<')
	if i == -1 {
		return raw, "", true
	}
	name := strings.TrimSpace(raw[:i])
	l := len(raw)
	if raw[l-1] != '>' {
		return "", "", false
	}
	return raw[i+1 : l-1], name, true
}

// Loads the routes involving templates.
func loadTemplateRoutes(app *fiber.App) {
	// Loads the index page.
	app.Get("/", func(ctx *fiber.Ctx) error {
		ptr, err := db.SelectAllEmailAddresses(ctx.Context())
		if err != nil {
			return err
		}
		emails, err := ptr.All()
		if err != nil {
			return err
		}
		a, err := db.SelectAllInvoiceJobs(ctx.Context())
		if err != nil {
			return err
		}
		return ctx.Render("templates/index", fiber.Map{
			"Jobs":   a,
			"Emails": emails,
			"Error":  "",
			"Title":  "Invoice Jobs",
		}, "templates/base")
	})

	// Edits the job.
	editJob := func(ctx *fiber.Ctx, jobId, emailContent, markdown, subject string) string {
		i, err := strconv.Atoi(jobId)
		if err != nil {
			return "Job ID must be a number."
		}
		if err = db.UpdateInvoiceInformation(ctx.Context(), i, emailContent, markdown, subject); err != nil {
			return err.Error()
		}
		return ""
	}

	// Deletes the job.
	deleteJob := func(ctx *fiber.Ctx, jobId string) string {
		i, err := strconv.Atoi(jobId)
		if err != nil {
			return "Job ID must be a number."
		}
		if err = db.DeleteInvoiceJob(ctx.Context(), i); err != nil {
			return err.Error()
		}
		return ""
	}

	// Inserts a job.
	insertJob := func(
		ctx *fiber.Ctx, fromName, fromEmail, toName,
		toEmail, markdown, emailContent, subject, date, monthsBetweenTrigger string,
	) string {
		// Parse the date.
		dateParsedInt, err := strconv.Atoi(date)
		if err != nil {
			return "The date must be a number."
		}
		dateParsed := uint8(dateParsedInt)
		if 1 > dateParsed {
			return "Invalid date."
		}
		if dateParsed > 28 {
			return "Date must be before 28th."
		}

		// Parse the months between trigger.
		monthsParsedInt, err := strconv.Atoi(monthsBetweenTrigger)
		if err != nil {
			return "The months between trigger must be a number."
		}
		monthsParsed := uint8(monthsParsedInt)
		if 1 > monthsParsed {
			return "Invalid months between trigger."
		}

		// Insert the job into the database.
		if err = db.InsertInvoiceJob(ctx.Context(), &db.InvoiceJob{
			Subject:              subject,
			FromEmail:            fromEmail,
			FromName:             fromName,
			Date:                 dateParsed,
			MonthsBetweenTrigger: monthsParsed,
			ToEmail:              toEmail,
			ToName:               toName,
			Markdown:             markdown,
			EmailContent:         emailContent,
		}); err != nil {
			return err.Error()
		}

		// Return no errors.
		return ""
	}

	// Process the jobs form.
	processJobsForm := func(ctx *fiber.Ctx) string {
		// Check if the job ID is present.
		jobId := ctx.FormValue("job_id")
		markdown := ctx.FormValue("markdown")
		emailContent := ctx.FormValue("email_content")
		subject := ctx.FormValue("subject")
		if jobId != "" {
			// The job ID is present. This means that this request is to edit a pre-existing job.

			// If there's no markdown, delete the job.
			if markdown == "" {
				return deleteJob(ctx, jobId)
			}

			// Edit the job.
			return editJob(ctx, jobId, emailContent, markdown, subject)
		}

		// Get the form contents.
		fromRaw := ctx.FormValue("from_email")
		toName := ctx.FormValue("to_name")
		toEmail := ctx.FormValue("to_email")
		date := ctx.FormValue("date")
		monthsBetweenTrigger := ctx.FormValue("months_between_trigger")
		if fromRaw == "" || markdown == "" || toEmail == "" {
			return "Required fields are missing."
		}
		if date == "" {
			date = "1"
		}
		if monthsBetweenTrigger == "" {
			monthsBetweenTrigger = "1"
		}

		// Process the from e-mail.
		fromEmail, fromName, success := parseEmail(fromRaw)
		if !success {
			return "Unable to parse from e-mail address."
		}

		// Processes a job insert.
		return insertJob(ctx, fromName, fromEmail, toName, toEmail, markdown, emailContent, subject, date, monthsBetweenTrigger)
	}
	app.Post("/", func(ctx *fiber.Ctx) error {
		errorMessage := processJobsForm(ctx)
		ptr, err := db.SelectAllEmailAddresses(ctx.Context())
		if err != nil {
			return err
		}
		emails, err := ptr.All()
		if err != nil {
			return err
		}
		a, err := db.SelectAllInvoiceJobs(ctx.Context())
		if err != nil {
			return err
		}
		return ctx.Render("templates/index", fiber.Map{
			"Jobs":   a,
			"Emails": emails,
			"Error":  errorMessage,
			"Title":  "Invoice Jobs",
		}, "templates/base")
	})

	// Lists all the failures.
	app.Get("/failures", func(ctx *fiber.Ctx) error {
		ptr, err := db.SelectAllFailures(ctx.Context())
		if err != nil {
			return err
		}
		a, err := ptr.All()
		if err != nil {
			return err
		}
		return ctx.Render("templates/failures", fiber.Map{
			"Failures": a,
			"Title":    "Failures",
		}, "templates/base")
	})

	// Gets all of the e-mail addresses.
	app.Get("/emails", func(ctx *fiber.Ctx) error {
		ptr, err := db.SelectAllEmailAddresses(ctx.Context())
		if err != nil {
			return err
		}
		a, err := ptr.All()
		if err != nil {
			return err
		}
		return ctx.Render("templates/emails", fiber.Map{
			"Records": a,
			"Error":   "",
			"Title":   "E-mail Accounts",
		}, "templates/base")
	})

	// Used to process password updates.
	processPasswordUpdate := func(ctx context.Context, name, email, password string) string {
		err := db.ChangeEmailAddressPassword(ctx, email, name, password)
		if err == nil {
			return ""
		}
		return err.Error()
	}

	// Used to process username updates.
	processUsernameUpdate := func(ctx context.Context, name, email, username string) string {
		if username == email {
			username = ""
		}
		err := db.ChangeEmailAddressUsername(ctx, email, name, username)
		if err == nil {
			return ""
		}
		return err.Error()
	}

	// Used to process inserts.
	processInsert := func(ctx context.Context, name, email, username, password, hostname string) string {
		if hostname == "" || email == "" {
			return "Invalid request. Please make a bug report if you hit this!"
		}
		err := db.InsertEmailAddress(ctx, &db.EmailAddress{
			FromEmail: email,
			FromName:  name,
			Username:  username,
			Password:  password,
			SMTPHost:  hostname,
		})
		if err == nil {
			return ""
		}
		return err.Error()
	}

	// Used to delete an e-mail address.
	deleteAddress := func(ctx context.Context, name, email string) string {
		err := db.DeleteEmailAddress(ctx, email, name)
		if err == nil {
			return ""
		}
		return err.Error()
	}

	// Used to process the e-mails form.
	processEmailsForm := func(ctx *fiber.Ctx) string {
		// Get the name and email (minimum required).
		name := ctx.FormValue("name")
		email := ctx.FormValue("email")
		if email == "" {
			return "Invalid parameters provided."
		}

		// Get the username and password. Figure out the request type from this.
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")
		if !ctx.Request().PostArgs().Has("username") {
			// Check if this is a password reset or a delete request.
			hasPassword := ctx.Request().PostArgs().Has("password")
			if hasPassword {
				return processPasswordUpdate(ctx.Context(), name, email, password)
			}
			return deleteAddress(ctx.Context(), name, email)
		}

		// If hostname doesn't exist, it's a username update. Else it's a creation.
		hostname := ctx.FormValue("hostname")
		if hostname == "" {
			// The combo has to be name, email and username.
			return processUsernameUpdate(ctx.Context(), name, email, username)
		}

		// Process the insert.
		return processInsert(ctx.Context(), name, email, username, password, hostname)
	}

	// Handles e-mail form submissions.
	app.Post("/emails", func(ctx *fiber.Ctx) error {
		errorMessage := processEmailsForm(ctx)
		ptr, err := db.SelectAllEmailAddresses(ctx.Context())
		if err != nil {
			return err
		}
		a, err := ptr.All()
		if err != nil {
			return err
		}
		return ctx.Render("templates/emails", fiber.Map{
			"Records": a,
			"Error":   errorMessage,
			"Title":   "E-mail Accounts",
		}, "templates/base")
	})
}
