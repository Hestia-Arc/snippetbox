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
