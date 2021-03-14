package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"net/smtp"
	"strconv"
	"strings"
)

// EmailAddress is used to define a email in the database.
type EmailAddress struct {
	// FromEmail is used to define an e-mail address this is from.
	FromEmail string `json:"from_email"`

	// FromName is used to define a name this is from.
	FromName string `json:"from_name"`

	// Username is the login username which is being used.
	Username string `json:"username"`

	// Password is the login password which is being used.
	Password string `json:"password,omitempty"`

	// SMTPHost is the SMTP host which is being used.
	SMTPHost string `json:"smtp_host"`
}

// Returns a split version of the hostname.
func (e *EmailAddress) splitHostname() (hostname string, port int, err error) {
	// Make a copy of the hostname.
	hostname = e.SMTPHost

	// If there's no colon, return defaults.
	if !strings.Contains(hostname, ":") {
		port = 25
		return
	}

	// Split by the colon and parse the port.
	split := strings.SplitN(hostname, ":", 2)
	if port, err = strconv.Atoi(split[1]); err != nil {
		return
	}
	hostname = split[0]

	// Return now we are done.
	return
}

// SendEmail is used to send an email with the credentials in the struct.
func (e *EmailAddress) Send(Recipients []string, Body []byte) error {
	// Formats the hostname to make sure it's fine.
	hostname, port, err := e.splitHostname()
	if err != nil {
		return err
	}

	// Handle authentication.
	auth := smtp.PlainAuth("", e.Username, e.Password, hostname)

	// Send the e-mail.
	return smtp.SendMail(hostname+":"+strconv.Itoa(port), auth, e.FromEmail, Recipients, Body)
}

// InsertEmailAddress is used to insert into the database.
func InsertEmailAddress(ctx context.Context, email *EmailAddress) error {
	var username *string
	if email.Username != "" {
		username = &email.Username
	}
	_, err := conn.Exec(ctx,
		"INSERT INTO email_addresses (from_email, from_name, username, password, smtp_host) VALUES ($1, $2, $3, $4, $5)",
		email.FromEmail, email.FromName, username, email.Password, email.SMTPHost)
	return err
}

// DeleteEmailAddress is used to delete a e-mail address from the table.
func DeleteEmailAddress(ctx context.Context, FromEmail, FromName string) error {
	_, err := conn.Exec(ctx, "DELETE FROM email_addresses WHERE from_email = $1 AND from_name = $2", FromEmail, FromName)
	return err
}

// ChangeEmailAddressUsername is used to change the username used for an e-mail address.
func ChangeEmailAddressUsername(ctx context.Context, FromEmail, FromName, Username string) error {
	var u *string
	if Username != "" {
		u = &Username
	}
	_, err := conn.Exec(ctx, "UPDATE email_addresses SET username = $1 WHERE from_email = $2 AND from_name = $3", u, FromEmail, FromName)
	return err
}

// ChangeEmailAddressPassword is used to change the password used for an e-mail address.
func ChangeEmailAddressPassword(ctx context.Context, FromEmail, FromName, Password string) error {
	_, err := conn.Exec(ctx, "UPDATE email_addresses SET password = $1 WHERE from_email = $2 AND from_name = $3", Password, FromEmail, FromName)
	return err
}

// EmailAddressRows is used to get the rows from the e-mail addresses table.
type EmailAddressRows struct {
	rows pgx.Rows
}

// Close is used to close the pointer.
func (e EmailAddressRows) Close() {
	e.rows.Close()
}

// Next passes through the Next function from pgx.Rows.
func (e EmailAddressRows) Next() bool {
	return e.rows.Next()
}

// Err passes through the Err function from pgx.Rows.
func (e EmailAddressRows) Err() error {
	return e.rows.Err()
}

// Scan is used to scan the row into a EmailAddress object.
func (e EmailAddressRows) Scan() (*EmailAddress, error) {
	var x EmailAddress
	var username *string
	if err := e.rows.Scan(&x.FromEmail, &x.FromName, &username, &x.SMTPHost); err != nil {
		return nil, err
	}
	if username != nil {
		x.Username = *username
	}
	return &x, nil
}

// All is used to get all of the e-mail addresses from the database.
func (e EmailAddressRows) All() ([]*EmailAddress, error) {
	a := make([]*EmailAddress, 0)
	for e.Next() {
		email, err := e.Scan()
		if err != nil {
			return nil, err
		}
		a = append(a, email)
	}
	return a, nil
}

// SelectAllEmailAddresses is used to create a pointer to select all e-mail addresses from the table.
// Note this does NOT include the password for security reasons. Changing the password is expected to be another (blind) request.
func SelectAllEmailAddresses(ctx context.Context) (*EmailAddressRows, error) {
	rows, err := conn.Query(ctx, "SELECT from_email, from_name, username, smtp_host FROM email_addresses")
	if err != nil {
		return nil, err
	}
	return &EmailAddressRows{rows: rows}, nil
}

// SelectOneEmailAddress is used to select a single record from the database.
func SelectOneEmailAddress(ctx context.Context, FromEmail, FromName string) (*EmailAddress, error) {
	e := EmailAddress{
		FromEmail: FromEmail,
		FromName:  FromName,
	}
	if err := conn.QueryRow(ctx,
		"SELECT username, password, smtp_host FROM email_addresses"+
			" WHERE from_email = $1 AND from_name = $2", FromEmail, FromName,
	).Scan(&e.Username, &e.Password, &e.SMTPHost); err != nil {
		return nil, err
	}
	return &e, nil
}
