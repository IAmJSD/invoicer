-- Defines the email address schema.
CREATE TABLE IF NOT EXISTS email_addresses (
    -- Defines the e-mail this is from.
    from_email TEXT NOT NULL,

    -- Defines the name this is from.
    from_name TEXT NOT NULL,

    -- Defines the hostname.
    smtp_host TEXT NOT NULL,

    -- Defines the login username which is used.
    username TEXT,

    -- Defines the login password which is used.
    password TEXT NOT NULL
);

-- Create a unique index for the from e-mail/name.
CREATE UNIQUE INDEX IF NOT EXISTS email_addresses_from_name_email ON email_addresses (from_email, from_name);

-- Defines the invoice jobs.
CREATE TABLE IF NOT EXISTS invoice_jobs (
    -- Defines the e-mail address this is from. This maps to the email_addresses table.
    from_email TEXT NOT NULL,
    from_name TEXT NOT NULL,

    -- Defines the job ID.
    job_id SERIAL PRIMARY KEY,

    -- Defines the date this triggers.
    date SMALLINT NOT NULL CHECK ( date >= 1 AND 28 >= date ),

    -- Defines how many months between each trigger.
    months_between_trigger SMALLINT NOT NULL CHECK ( months_between_trigger >= 1 ),

    -- Defines who the email is to.
    to_email TEXT NOT NULL,
    to_name TEXT NOT NULL,

    -- Defines the markdown.
    markdown TEXT NOT NULL,

    -- Defines when this was last triggered.
    last_triggered TIMESTAMPTZ,

    -- Defines the e-mail subject.
    subject TEXT NOT NULL,

    -- Defines the e-mail contents.
    email_content TEXT NOT NULL
);

-- Defines any failures.
CREATE TABLE IF NOT EXISTS failures (
    -- When it failed.
    failure_time TIMESTAMPTZ NOT NULL,

    -- The error message.
    message TEXT NOT NULL
);

-- Defines any variables for a job.
CREATE TABLE IF NOT EXISTS job_variables (
    -- Defines the job ID this belongs to.
    job_id INTEGER NOT NULL,

    -- Defines the variable key.
    key TEXT NOT NULL,

    -- Defines the value.
    value TEXT NOT NULL
);

-- Create a unique index for the key and job ID.
CREATE UNIQUE INDEX IF NOT EXISTS job_variables_id_key ON job_variables (job_id, key);
