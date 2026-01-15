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

=== The http.Server error log ===

There is one more change we need to make to our application. By default, if Go’s HTTP server
encounters an error it will log it using the standard logger. For consistency it’d be better to
use our new errorLog logger instead.

To make this happen we need to initialize a new http.Server struct containing the
configuration settings for our server, instead of using the http.ListenAndServe() shortcut.

=== Concurrent logging ===

Custom loggers created by log.New() are concurrency-safe. You can share a single logger and
use it across multiple goroutines and in your handlers without needing to worry about race
conditions.
That said, if you have multiple loggers writing to the same destination then you need to be
careful and ensure that the destination’s underlying Write() method is also safe for
concurrent use.

=== Logging to a file ===

As I said above, my general recommendation is to log your output to standard streams and
redirect the output to a file at runtime. But if you don’t want to do this, you can always open a
file in Go and use it as your log destination. As a rough example:

f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
if err != nil {
log.Fatal(err)
}
defer f.Close()
infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

3. Dependency injection

You’ll notice that the home handler function is still writing error messages
using Go’s standard logger, not the errorLog logger that we want to be using.

This raises a good question: how can we make our new errorLog logger available to our home
function from main()?
And this question generalizes further. Most web applications will have multiple dependencies
that their handlers need to access, such as a database connection pool, centralized error
handlers, and template caches. What we really want to answer is: how can we make any
dependency available to our handlers?

The simplest being to just put the dependencies in
global variables. But in general, it is good practice to inject dependencies into your handlers. It
makes your code more explicit, less error-prone and easier to unit test than if you use global
variables.
For applications where all your handlers are in the same package, like ours, a neat way to
inject dependencies is to put them into a custom application struct, and then define your
handler functions as methods against application.

Let’s try this out by quickly adding a deliberate error to our application.

Open your terminal and rename the ui/html/pages/home.tmpl to ui/html/pages/home.bak.
When we run our application and make a request for the home page, this now should result in
an error because the ui/html/pages/home.tmpl no longer exists.
Go ahead and make the change:

$ cd $HOME/code/snippetbox
$ mv ui/html/pages/home.tmpl ui/html/pages/home.bak

Then run the application and make a request to http://localhost:4000. You should get an
Internal Server Error HTTP response in your browser, and see a corresponding error
message in your terminal similar to this:

$ go run ./cmd/web
INFO  
2022/01/29 16:12:36 Starting server on :4000
ERROR  
2022/01/29 16:12:40 handlers.go:29: open ./ui/html/pages/home.tmpl: no such file or directory

Notice how the log message is now prefixed with ERROR and originated from line 25 of the
handlers.go file? This demonstrates nicely that our custom errorLog logger is being passed
through to our home handler as a dependency, and is working as expected.

=== Closures for dependency injection ===

The pattern that we’re using to inject dependencies won’t work if your handlers are spread
across multiple packages. In that case, an alternative approach is to create a config package
exporting an Application struct and have your handler functions close over this to form a
closure. Very roughly:

func main() {
app := &config.Application{
ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}
mux.Handle("/", examplePackage.ExampleHandler(app))
}

func ExampleHandler(app *config.Application) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
...
ts, err := template.ParseFiles(files...)
if err != nil {
app.ErrorLog.Println(err.Error())
http.Error(w, "Internal Server Error", 500)
return
}
...
}
}

4. Centralized error handling

Let’s neaten up our application by moving some of the error handling code into helper
methods. This will help:

- separate our concerns and
- stop us repeating code as we progress through the build.

Add a new helpers.go file under the cmd/web directory:

$ cd $HOME/code/snippetbox
$ touch cmd/web/helpers.go

There’s not a huge amount of new code here, but it does introduce a couple of features which
are worth discussing.

- In the serverError() helper we use the debug.Stack() function to get a stack trace for the
  current goroutine and append it to the log message. Being able to see the execution path
  of the application via the stack trace can be helpful when you’re trying to debug errors.

- In the clientError() helper we use the http.StatusText() function to automatically
  generate a human-friendly text representation of a given HTTP status code. For example,
  http.StatusText(400) will return the string "Bad Request".
  Once that’s done, head back to your handlers.go file and update it to use the new helpers.

When that’s updated, restart your application and make a request to http://localhost:4000
in your browser.

Again, this should result in our (deliberate) error being raised and you should see the
corresponding error message and stack trace in your terminal.

---

If you look closely at this you’ll notice a small problem: the file name and line number being
reported in the ERROR log line is now helpers.go:13 — because this is where the log message
is now being written from.
What we want to report is the file name and line number one step back in the stack trace,
which would give us a clearer idea of where the error actually originated from.
We can do this by changing the serverError() helper to use our logger’s Output() function
and setting the frame depth to 2. Reopen your helpers.go file and update:

    app.errorLog.Output(2, trace)

And if you try again now, you should find that the appropriate file name and line number
(handlers.go:25) is being reported in the ERROR log line

---

Revert the deliberate error

$ mv ui/html/pages/home.bak ui/html/pages/home.tmp
