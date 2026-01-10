1. Project setup and creating a module

open your terminal and create a new project directory called
snippetbox anywhere on your computer.

I’m going to locate my project directory under
$HOME/code, but you can choose a different location if you wish.

$ mkdir -p $HOME/code/snippetbox

2. Creating a module

$ cd $HOME/code/snippetbox
$ go mod init snippetbox.platinumhestia
go: creating new go.mod: module snippetbox.platinumhestia

3. Hello world!

create a new main.go in your project directory containing the following code:

$ touch main.go

---

File: main.go
package main

import "fmt"

func main() {
fmt.Println("Hello world!")
}

---

Save this file, then use the go run . command in your terminal to compile and execute the
code in the current directory. All being well, you will see the following output:

$ go run .
Hello world!

Additional information
Module paths for downloadable packages
If you’re creating a project which can be downloaded and used by other people and
programs, then it’s good practice for your module path to equal the location that the code
can be downloaded from.
For instance, if your package is hosted at https://github.com/foo/bar then the module path
for the project should be github.com/foo/bar.

4.  Web application basics (handler, router, server)

    Now that everything is set up correctly let’s make the first iteration of our web application.
    We’ll begin with the three absolute essentials:

    a. The first thing we need is a handler. If you’re coming from an MVC-background, you can
    think of handlers as being a bit like controllers. They’re responsible for executing your
    application logic and for writing HTTP response headers and bodies.

    b.The second component is a router (or servemux in Go terminology). This stores a mapping
    between the URL patterns for your application and the corresponding handlers. Usually
    you have one servemux for your application containing all your routes.

    c. The last thing we need is a web server. One of the great things about Go is that you can
    establish a web server and listen for incoming requests as part of your application itself.
    You don’t need an external third-party server like Nginx or Apache.

    Let’s put these components together in the main.go file to make a working application

    Note:
    -- The home handler function is just a regular Go function with two parameters.
    -- The http.ResponseWriter parameter provides methods for assembling a HTTP response
    and sending it to the user, and
    -- the \*http.Request parameter is a pointer to a struct which holds information about the
    current request (like the HTTP method and the URL
    being requested).

    When you run this code in main.go,
    -- it should start a web server listening on port 4000 of your local machine.
    -- Each time the server receives a new HTTP request it will pass the request on to the
    servemux and — in turn — the servemux will check the URL path and dispatch the request to
    the matching handler.

    server (listen on port 4000), receives request and pass request to the servemux (router)
    servemux (check url path) and dispatch request to matching handler

---

Let’s give this a whirl. Save your main.go file and then try running it from your terminal using
the go run command.
$ cd $HOME/code/snippetbox
$ go run .
2022/01/29 11:13:26 Starting server on :4000
While the server is running, open a web browser and try visiting http://localhost:4000. If
everything has gone to plan you should see a page which looks a bit like this:

---

If you head back to your terminal window, you can stop the server by pressing Ctrl+c on your
keyboard

5. Routing requests

6.Project structure and organization

a. The cmd directory will contain the application-specific code for the executable applications
in the project. For now we’ll have just one executable application — the web application —
which will live under the cmd/web directory.

b. The internal directory will contain the ancillary non-application-specific code used in the
project. We’ll use it to hold potentially reusable code like validation helpers and the SQL
database models for the project.

c. The ui directory will contain the user-interface assets used by the web application.
Specifically, the ui/html directory will contain HTML templates, and the ui/static
directory will contain static files (like CSS and images).

So why are we using this structure?
There are two big benefits:

1. It gives a clean separation between Go and non-Go assets. All the Go code we write will live exclusively under the cmd and internal directories, leaving the project root free to
   hold non-Go assets like UI files, makefiles and module definitions (including our go.mod
   file). This can make things easier to manage when it comes to building and deploying your
   application in the future.
2. It scales really nicely if you want to add another executable application to your project.
   For example, you might want to add a CLI (Command Line Interface) to automate some
   administrative tasks in the future. With this structure, you could create this CLI application
   under cmd/cli and it will be able to import and reuse all the code you’ve written under
   the internal directory

   7. HTML templating and inheritance

   a. creating a template file
   $ cd $HOME/code/snippetbox
   $ mkdir ui/html/pages
   $ touch ui/html/pages/home.tmpl

   b. how do we get our home handler to render it?

   - For this we need to use Go’s html/template package, which provides a family of functions for
     safely parsing and rendering HTML templates. We can use the functions in this package to
     parse the template file and then execute the template.

restart the application:

$ cd $HOME/code/snippetbox
$ go run ./cmd/web
2022/01/29 12:06:02 Starting server on :4000

Then open http://localhost:4000 in your web browser. You should find that the HTML homepage is shaping up nicely

=== Template composition ====

As we add more pages to this web application there will be some shared, boilerplate, HTML
markup that we want to include on every page — like the header, navigation and metadata
inside the <head> HTML element.

To save us typing and prevent duplication, it’s a good idea to create a base (or master)
template which contains this shared content, which we can then compose with the page
specific markup for the individual pages.

create a new ui/html/base.tmpl file…
$ touch ui/html/base.tmpl

- Here we’re using the {{define "base"}}...{{end}} action to define a distinct named
  template called base, which contains the content we want to appear on every page.
- Inside this we use the {{template "title" .}} and {{template "main" .}} actions to
  denote that we want to invoke other named templates (called title and main) at a particular
  point in the HTML.
  Note: If you’re wondering, the dot at the end of the {{template "title" .}} action
  represents any dynamic data that you want to pass to the invoked template.

---

So now, instead of containing HTML directly, our template set contains 3 named templates —
base, title and main. We use the ExecuteTemplate() method to tell Go that we specifically
want to respond using the content of the base template (which in turn invokes our title and
main templates).

Feel free to restart the server and give this a try. You should find that it renders the same
output as before (although there will be some extra whitespace in the HTML source where the
actions are)

---

=== Embedding partials (reusable component) ===

For some applications you might want to break out certain bits of HTML into partials that can
be reused in different pages or layouts. To illustrate, let’s create a partial containing the
primary navigation bar for our web application

a. Create a new ui/html/partials/nav.tmpl file containing a named template called "nav"
$ mkdir ui/html/partials
$ touch ui/html/partials/nav.tmpl

b. Then update the base template so that it invokes the navigation partial using the
{{template "nav" .}} action.

c. Finally, we need to update the home handler to include the new ui/html/partials/nav.tmpl
file when parsing the template files

d. Once you restart the server, the base template should now invoke the nav template.

8. Serving static files

- Now let’s improve the look and feel of the home page by adding some static CSS and image
  files to our project,
- along with a tiny bit of JavaScript to highlight the active navigation item.

you can grab the necessary files and extract them into the
ui/static folder that we made earlier with the following commands:

$ cd $HOME/code/snippetbox
$ curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static/

=== The http.Fileserver handler ===

- Go’s net/http package ships with a built-in http.FileServer handler which you can use to
  serve files over HTTP from a specific directory.
- Let’s add a new route to our application so that
  all requests which begin with "/static/" are handled using this.

- Remember: The pattern "/static/" is a subtree path pattern, so it acts a bit like there
  is a wildcard at the end.

- To create a new http.FileServer handler, we need to use the http.FileServer() function
  like this:

  fileServer := http.FileServer(http.Dir("./ui/static/"))

  - When this handler receives a request, it will remove the leading slash from the URL path and
    then search the ./ui/static directory for the corresponding file to send to the user.
  - So, for this to work correctly, we must strip the leading "/static" from the URL path before
    passing it to http.FileServer. Otherwise it will be looking for a file which doesn’t exist and
    the user will receive a 404 page not found response. Fortunately Go includes a
    http.StripPrefix() helper specifically for this task.

Go’s file server has a few really nice features that are worth mentioning:

- It sanitizes all request paths by running them through the path.Clean() function before
  searching for a file. This removes any . and .. elements from the URL path, which helps to
  stop directory traversal attacks. This feature is particularly useful if you’re using the
  fileserver in conjunction with a router that doesn’t automatically sanitize URL paths.

Disabling directory listings
If you want to disable directory listings there are a few different approaches you can take.
The simplest way? Add a blank index.html file to the specific directory that you want to
disable listings for. This will then be served instead of the directory listing, and the user will
get a 200 OK response with no body. If you want to do this for all directories under
./ui/static you can use the command:

$ find ./ui/static -type d -exec touch {}/index.html \;

9. The http.Handler interface

Strictly speaking, what we mean by handler is an object which satisfies the http.Handler interface:

type Handler interface {
ServeHTTP(ResponseWriter, \*Request)
}

In simple terms, this basically means that to be a handler an object must have a ServeHTTP()
method with the exact signature:

ServeHTTP(http.ResponseWriter, \*http.Request)

So in its simplest form a handler might look something like this:

type home struct {}
func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
w.Write([]byte("This is my home page"))
}

Here we have an object (in this case it’s a home struct, but it could equally be a string or
function or anything else), and we’ve implemented a method with the signature
ServeHTTP(http.ResponseWriter, \*http.Request) on it. That’s all we need to make a handler.

You could then register this with a servemux using the Handle method like so:

mux := http.NewServeMux()
mux.Handle("/", &home{})

When this servemux receives a HTTP request for "/", it will then call the ServeHTTP() method
of the home struct — which in turn writes the HTTP response.

=== Handler functions ===

Now, creating an object just so we can implement a ServeHTTP() method on it is long-winded
and a bit confusing. Which is why in practice it’s far more common to write your handlers as a
normal function (like we have been so far). For example:

func home(w http.ResponseWriter, r \*http.Request) {
w.Write([]byte("This is my home page"))
}

But this home function is just a normal function; it doesn’t have a ServeHTTP() method. So in
itself it isn’t a handler.
Instead we can transform it into a handler using the http.HandlerFunc() adapter, like so:

mux := http.NewServeMux()
mux.Handle("/", http.HandlerFunc(home))

The http.HandlerFunc() adapter works by automatically adding a ServeHTTP() method to
the home function. When executed, this ServeHTTP() method then simply calls the content of
the original home function. It’s a roundabout but convenient way of coercing a normal function
into satisfying the http.Handler interface.

Throughout this project so far we’ve been using the HandleFunc() method to register our
handler functions with the servemux. This is just some syntactic sugar that transforms a
function to a handler and registers it in one step, instead of having to do it manually. The
code above is functionality equivalent to this:

mux := http.NewServeMux()
mux.HandleFunc("/", home)

=== Chaining handlers ===

The eagle-eyed of you might have noticed something interesting right at the start of this
project. The http.ListenAndServe() function takes a http.Handler object as the second
parameter…

func ListenAndServe(addr string, handler Handler) error

… but we’ve been passing in a servemux.

We were able to do this because the servemux also has a ServeHTTP() method, meaning that
it too satisfies the http.Handler interface.

For me it simplifies things to think of the servemux as just being a special kind of handler,
which instead of providing a response itself passes the request on to a second handler. This
isn’t as much of a leap as it might first sound. Chaining handlers together is a very common
idiom in Go, and something that we’ll do a lot of later in this project.

In fact, what exactly is happening is this: When our server receives a new HTTP request, it calls
the servemux’s ServeHTTP() method. This looks up the relevant handler based on the
request URL path, and in turn calls that handler’s ServeHTTP() method. You can think of a Go
web application as a chain of ServeHTTP() methods being called one after another.

=== Requests are handled concurrently ===

There is one more thing that’s really important to point out: all incoming HTTP requests are
served in their own goroutine. For busy servers, this means it’s very likely that the code in or
called by your handlers will be running concurrently. While this helps make Go blazingly fast,
the downside is that you need to be aware of (and protect against) race conditions when
accessing shared resources from your handlers.
