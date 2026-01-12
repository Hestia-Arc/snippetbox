=== Configuration and error handling ===

Making improvements that’ll make it easier to manage our
application as it grows.

- Managing configuration settings
- Leveled logging
- Dependency injection
- Centralized error handling

1. Managing configuration settings

- Our web application’s main.go file currently contains a couple of hard-coded configuration
  settings:

The network address for the server to listen on (currently ":4000")
The file path for the static files directory (currently "./ui/static")

- Having these hard-coded isn’t ideal. There’s no separation between our configuration
  settings and code, and we can’t change the settings at runtime (which is important if you
  need different settings for development, testing and production environments).

=== Command-line flags ===

In Go, a common and idiomatic way to manage configuration settings is to use command-line
flags when starting an application. For example:
$ go run ./cmd/web -addr=":80"
The easiest way to accept and parse a command-line flag from your application is with a line
of code like this:
addr := flag.String("addr", ":4000", "HTTP network address")
This essentially defines a new command-line flag with the name addr, a default value of
":4000" and some short help text explaining what the flag controls. The value of the flag will
be stored in the addr variable at runtime.

Try using the -addr flag when you start the application. You should find that
the server now listens on whatever address you specify, like so:

$ go run ./cmd/web -addr=":9999"
2022/01/29 15:50:20 Starting server on :9999

Note: Ports 0-1023 are restricted and (typically) can only be used by services which have
root privileges. If you try to use one of these ports you should get a
bind: permission denied error message on start-up.

=== Environment variables ===

If you want, you can store your configuration settings in environment variables and access
them directly from your application by using the os.Getenv() function like so:

addr := os.Getenv("SNIPPETBOX_ADDR")

But this has some drawbacks compared to using command-line flags. You can’t specify a
default setting (the return value from os.Getenv() is the empty string if the environment
variable doesn’t exist), you don’t get the -help functionality that you do with command-line
flags, and the return value from os.Getenv() is always a string — you don’t get automatic
type conversions like you do with flag.Int() and the other command line flag functions.
Instead, you can get the best of both worlds by passing the environment variable as a
command-line flag when starting the application. For example:

$ export SNIPPETBOX_ADDR=":9999"
$ go run ./cmd/web -addr=$SNIPPETBOX_ADDR
2022/01/29 15:54:29 Starting server on :9999

2. Leveled logging

At the moment in our main.go file we’re outputting log messages using the log.Printf() and
log.Fatal() functions.

Both these functions output messages via Go’s standard logger, which — by default —
prefixes messages with the local date and time and writes them to the standard error stream
(which should display in your terminal window). The log.Fatal() function will also call
os.Exit(1) after writing the message, causing the application to immediately exit.

In our application, we can break apart our log messages into two distinct types — or levels.
The first type is informational messages (like "Starting server on :4000") and the second
type is error messages.

log.Printf("Starting server on %s", *addr) // Information message
err := http.ListenAndServe(*addr, mux)
log.Fatal(err) // Error message

Let’s improve our application by adding some leveled logging capability, so that information
and error messages are managed slightly differently. Specifically:

- We will prefix informational messages with "INFO" and output the message to standard
  out (stdout).

- We will prefix error messages with "ERROR" and output them to standard error (stderr),
  along with the relevant file name and line number that called the logger (to help with
  debugging).

There are a couple of different ways to do this, but a simple and clear approach is to use the
log.New() function to create two new custom loggers.

Tip: If you want to include the full file path in your log output, instead of just the file
name, you can use the log.Llongfile flag instead of log.Lshortfile when creating
your custom logger. You can also force your logger to use UTC datetimes (instead of
local ones) by adding the log.LUTC flag.

=== Decoupled logging ===

A big benefit of logging your messages to the standard streams (stdout and stderr) like we are
is that your application and logging are decoupled.

Your application itself isn’t concerned
with the routing or storage of the logs, and that can make it easier to manage the logs
differently depending on the environment.

During development, it’s easy to view the log output because the standard streams are
displayed in the terminal.

In staging or production environments, you can redirect the streams to a final destination for
viewing and archival. This destination could be on-disk files, or a logging service such as
Splunk. Either way, the final destination of the logs can be managed by your execution
environment independently of the application.

For example, we could redirect the stdout and stderr streams to on-disk files when starting
the application like so:

$ go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log

Note: Using the double arrow >> will append to an existing file, instead of truncating it
when starting the application.
