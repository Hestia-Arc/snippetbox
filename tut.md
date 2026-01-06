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
