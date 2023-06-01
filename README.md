# htmxtools
`htmxtools` is a collection of constants and utilities for working with [htmx](https://htmx.org) from Go.

`htmx` plays REALLY nice with Go and (somewhat less nice but still nice) Go templates.


## General usage
There are different things in here for different use cases - http middleware as well as some constants and helpers for htmx requests


### Constants and such
There are quite a few constants and enums for some various bits of htmx that can help cut down on some error prone duplication. 

#### Response Headers
The response headers can be used as described here: https://htmx.org/reference/#response_headers and they are one of the coolest things about htmx.

You can do things like ensure the browser url bar and history point to a real html page (as opposed to a fragment which is the default when using hx-get and the like).

I use this exact pattern in another project when an add or delete call is made to the backend (inside my handler):
```go
w.Header().Add(htmxtools.LocationResponse.String(), `{"path":"delete-status-fragment", "target":"#content-div"}`)
w.Header().Add(htmxtools.ReplaceURLResponse.String(), "status.html")
w.WriteHeader(http.StatusAccepted)
```
When the handler is called via an `hx-` request, the contents `delete-status-fragment` will replace the contents of the div with id `content-div`.
However, unlike the default behaviour, the url bar will show `status.html` and be safe to reload while the default would show a path of `delete-status-fragment` which is not a valid full html page

_if you feel like the headers aren't working, make sure you actually wrote them to the http response. I make this mistake ALL THE TIME_

There's also a helper if you want for building the headers in a safer way:

```go
if htmxRequest := htmxtools.RequestFromContext(r.Context()); htmxRequest != nil {
        hxheaders := &htmxtools.HTMXResponse{
			ReplaceURL: htmxRequest.CurrentURL,
            Reswap: htmxtools.SwapOuterHTML,
		}
		if err := hxheaders.AddToResponse(w); err != nil {
			return
		}
}
```

### Middleware
```go
http.Handle("/", htmxtools.Wrap(myhandler))
```
or
```go
http.Handle("/",htmxtools.WrapFunc(myhandlerfunc))
```

This will detect htmx requests and inject the details into the context.
You can extract the details down the line via:
```go
func somefunc(w http.ResponseWriter, r *http.Request) {
    htmxrequest := htmxtools.RequestFromContext(r.Context())
    if htmxrequest == nil {
        // do something for non-htmx requests
    } else {
        // do something for htmx requests
        // note that not all fields will be populated so you'll want to check
    }
}
```

## Examples
The following contains a few different ways to use this library

### in-repo example
In the `examples/middleware` directory there's a small example:
```
go run examples/middleware/main.go
```

will start a small webserver on http://localhost:3000

- Loading the page will present a button, that when clicked, makes an htmx request to the backend which returns details about the htmx request:

![main image page with two buttons labeled "drink me" and "eat me"](images/middleware-initial.png?raw=true "Main page")

- Clicking "drink me" will ask for input and render a response from the server:

![drink me button dialog](images/drinkme-1.png?raw=true)![drink me button response from server](images/drinkme-2.png?raw=true)
visiting http://localhost:3000/other.html will show a similar page which doesn't go to the backend for data but passes through the middleware which injects the htmx request details in the template

- Clicking "eat me" will respond with a server-side generated htmx alert via headers:
![browser alert dialog](images/eatme.png?raw=true)

- Visting http://localhost:3000/other.html , will render a page that operates entirely client side (with the exception of the template values that are injected by the middleware)

![initial other page with a button labled "Generate htmx request"](images/other-1.png?raw=true)

- clicking the button will generate an htmx request for the same file (with `hx-confirm` client side) but will be passed an htmx request struct via the template execution:

![initial other page with confirmation dialog presented](images/other-confirm.png?raw=true)
![initial other page post confirm with an html table showing the htmx request data](images/other-after.png?raw=true)

### Building your own handler
If you wanted to, you can use the various bits to build your handler the way you want:

```go
// responds to htmx requests with the provided template updating the provided target
// replace is the path to replace in the url bar
func hxOnlyHandler(template, target string, next http.Handler, replace string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        htmxrequest := htmxtools.RequestFromContext(r.Context())
        if htmxrequest == nil {
            next.ServeHTTP(w, r)
            return
        }
        htmxloc := fmt.Sprintf(`{"path":"%s","target":"%s"}`, template, target)
        w.Header().Add(htmxtools.LocationResponse.String(), htmxloc)
        if strings.TrimSpace(replace) != "" {
            w.Header().Add(htmxtools.ReplaceURLResponse.String(), replace)
        }
        w.WriteHeader(http.StatusAccepted)
    }
}
```

## Templating tips
I am not a go template expert. I've tended to avoid them in the past because of the runtime implications however now that I've been working on browser content, I had to really dive back in.

### Use blocks
Blocks are a REALLY REALLY nice thing in go templates. They allow you to define a template inside of another templare. Note that blocks are still rendered (unless something inside the block says not to).

*tl;dr: wrap chunks of html content that you might want to reuse via htmx in a ```{{ block "unique-block-name" .}}blahblah{{end}}```* and call them by block name via `hx-get` (note that this requires your server to understand serving templates - example below)

Take the following template (`index.html`) from this repo:
```html
<!DOCTYPE html>
<html lang="en">
{{ block "head" . }}
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.2"
        integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h"
        crossorigin="anonymous"></script>
    <script>
        htmx.logAll();
    </script>
    <title>Middleware Demo</title>
</head>
{{ end }}

{{ block "body" . }}
<body>
    <script type="text/javascript">
        // for our post-swap alert
        document.body.addEventListener("showMessage", function (evt) {
            if (evt.detail.level === "info") {
                alert(evt.detail.message);
            }
        });
    </script>
    <div id="button-div"><button id="drink-me-button" hx-get="button-push" hx-target="#body-content"
            name="drink-me-button-name" hx-prompt="are you sure? type something to confirm">Drink me!</button>
    </div>
    <!-- set hx-swap to none so our button doesn't get replaced with empty results -->
    <div id="server-side-alert"><button id="server-side-alert" hx-get="alert" hx-swap="none">Eat me!</button>
    </div>
    <div id="body-content"></div>
</body>
{{ end }}
</html>
```

This is our index page and can be rendered as a whole html page as is. 
Now let's take a look at another file in the same directory:

```html
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.2"
        integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h"
        crossorigin="anonymous"></script>
    <script>
        htmx.logAll();
    </script>
    <title>Other Page</title>
</head>

<body>
    {{ block "other" . }}
    <div id="other-content">
        {{ if not .HtmxRequest }}
        <!-- content when not htmx -->
        {{ else }}
        <!-- we're passing in the htmxrequest struct we got from the context so we can refer to its values here-->
        <!-- other content here-->
        {{ end }}
    </div>
    {{ end }}
</body>
</html>
```

If we wanted to, we could rewrite the above template like so:
```html
<!DOCTYPE html>
<html lang="en">

{{ template "head" . }}

<body>
    {{ block "other" . }}
    <div id="other-content">
        {{ if not .HtmxRequest }}
        <!-- content when not htmx -->
        {{ else }}
        <!-- we're passing in the htmxrequest struct we got from the context so we can refer to its values here-->
        <!-- other content here-->
        {{ end }}
    </div>
    {{ end }}
</body>
</html>
```

This is non-htmx specific mind you but it DOES open up some fun tricks with htmx. Imagine we have the following html

```html
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.2"
        integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h"
        crossorigin="anonymous"></script>
    <script>
        htmx.logAll();
    </script>
    <title>Other Page</title>
</head>

<body>
    {{ block "list-items" . }}
    <table><!-- insert your table rows and such here - maybe with template placeholdrs --></table>
</body>
</html>
```

Now let's say we want the `list-items` html to be rendered somewhere else. We can do this now:

```html
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.2"
        integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h"
        crossorigin="anonymous"></script>
    <script>
        htmx.logAll();
    </script>
    <title>Third Page</title>
</head>

<body>
    <div id="other-list-items" hx-get="list-items"></div>
</body>
</html>
```

The above div (`other-list-items`) will be replaced with the contents of the `list-items` block from above. 

Note that doing so will, by default, update the url and location history to contain `http://hostname/list-items` which is NOT a valid html page and would just render the table only. To work around this, set the `HX-Replace-Url` header to set to a valid page or `false`. You can also set the attribute on the element as well via [`hx-replace-url`](https://htmx.org/attributes/hx-replace-url/)


## TODO

- add structs for htmx json (i.e. json passed to `HX-Location`)