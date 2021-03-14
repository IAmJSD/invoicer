# Invoicer
A tool to manage invoices.

## Why?
I struggle to stay on top of admin stuff. Writing a tool to do something which has caused me a lot of problems due to not being able to schedule it mentally makes this a significantly easier task.

## How does it work?
Firstly, you deploy the tool somewhere. The tool expects you to create an application in the Cloudflare Access panel first if you want Cloudflare Access support. It isn't overly computationally expensive, a $5/month DigitalOcean droplet should be enough.

There are 2 ways you could go about this:
1. **Use the built-in Docker Compose file:** On a fresh droplet, clone this GitHub repository. From here, rename `Caddyfile.example` to `Caddyfile` and set your domain in the configuration where it says `domain_here` and e-mail where it says `email_here`. Next, copy `.env.example` to `.env` and set the `TEAM_DOMAIN` and `APPLICATION_AUD`. You should now just be able to do `docker-compose up -d`, remove the Cloudflare proxy for a bit so it can grab a SSL certificate, reapply the Cloudflare proxy and then this will be configured!
2. **Use your own Docker configuration:** You will want to set `CONNECTION_STRING` to the connection string for your Postgres database. If you want Cloudflare Access support, you will then want to set your `TEAM_DOMAIN` and `APPLICATION_AUD`. The application will listen on port 3000.

It's important to note the tool is not designed to be clustered. Attempting to do so can (and probably will) lead to SQL race conditions.

When you are writing the markdown which will be used in PDF, the application supports Go template syntax with some extensions. This includes built in Go template functions, [sprig](https://github.com/Masterminds/sprig) and `PersistentCounter <counter name>` (a counter which persists across invoices).
