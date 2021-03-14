package db

import "context"

// GetJobVariable is used to get the job variable from the database.
func GetJobVariable(ctx context.Context, JobID int, Key string) (string, error) {
	var value string
	err := conn.QueryRow(ctx, "SELECT value FROM job_variables WHERE job_id = $1 AND key = $2", JobID, Key).Scan(&value)
	return value, err
}

// SetJobVariable is used to set the job variable in the database.
func SetJobVariable(ctx context.Context, JobID int, Key, Value string) error {
	_, err := conn.Exec(ctx, "INSERT INTO job_variables (job_id, key, value) VALUES ($1, $2, $3) ON CONFLICT "+
		"(job_id, key) DO UPDATE SET value = $3", JobID, Key, Value)
	return err
}
