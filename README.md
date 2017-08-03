# mchain

A super tiny go package that handles middleware chaining in it's most minimal form. 

**Godoc:** https://godoc.org/github.com/prasannavl/mchain  
PS: `Read the source, Luke` - it's tiny.

## Get

`go get -u github.com/prasannavl/mchain`

## Standard middlewares

```go
type HttpMiddleware func(http.Handler) http.Handler

type HttpChain struct {
	Middlewares []HttpMiddleware
}
```

That's about it. It's even simpler than the very neat `alice` package. However, the HttpChain provides no `Append`, `Extend` like methods. They are cleanly separated into a builder - `HttpChainBuilder`, that provides all the composition. So, now the `Middlewares` field is public, and `HttpChain` can be transparently passed around, cloned, extended at will.

## mchain middlewares

The standard middleware pattern looks fine, however proves very difficult to chain error handling cleanly. So `mchain` provides this as the alternative middleware:

```go
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type Middleware func(Handler) Handler

type Chain struct {
	Middlewares []Middleware
}

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}
```

Actually, that's the entire `mchain.go` file. Very simple. This allows clean error handling. Errors can be used to communicate error status code as well. A pattern that can be easily achieved with `HttpError` from simple error composers like [`prasannavl/goerror`](https://www.github.com/prasannavl/goerror)

This aligns with Go's idiomatic way of error handling.

```go
err := h.ServeHTTP(w, r)
if err != nil {
// handle error
}
```

I personally think somewhere on it's way - the standard library team got stuck in the choice between simplicity and consistency - and they seem to have chosen the former. And now it's stuck - you can't just go back and change the standard way even if the other is deemed better.

But thankfully, you don't have to choose. You can combine both :)

`mchain` brings this pattern with almost no overhead. And it has a set of conversation functions that provide two way conversions between the standard `net/http` package, and `mchain`, like `MiddlewareFromHttp` and `MiddlewareToHttp`, `HandlerFromHttp`, and `HandlerToHttp` - that allows both to coexist, and mix and match both types of handlers.

### Example

```go

func newAppHandler(host string) http.Handler {
	c := appcontext.AppContext{Services: appcontext.Services{}}

	return mchain.NewBuilder(
		// An existing http handler based middleware
		mchain.MiddlewareFromHttp(c.HandlerWithContext, nil),
		middleware.RequestContextInitHandler,
		middleware.RequestLogHandler,
		middleware.RequestDurationHandler,
	).
	Handler(mchain.HandlerFromHttp(CreateActionHandler(host))).
	BuildHttp(nil)
}

func newHttpAppHandler(host string) http.Handler {
	c := appcontext.AppContext{Services: appcontext.Services{}}

	return mchain.NewHttpBuilder(
		c.HandlerWithContext,
		standardmiddleware.RequestContextInitHandler,
		standardmiddleware.RequestLogHandler,
		standardmiddleware.RequestDurationHandler,
	).
	Handler(CreateActionHandler(host)).
	Build()
}

func CreateActionHandler(host string) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Message string
			Date    time.Time
		}{
			fmt.Sprintf("Hello world from %s", host),
			time.Now(),
		}
		render.JSON(w, r, &data)
	}
	return http.HandlerFunc(f)
}

```

## Why return errors along with the handler?

Consider a similar middleware setup to above example

With `net/http` middleware chain:

```go
func RequestIDMustInitHandler(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		c := FromRequest(r)
		if _, ok := r.Header[RequestIDHeaderKey]; ok {
			http.Error(w, fmt.Sprintf("error: illegal header (%s)", RequestIDHeaderKey), 400)
			return
		}
		var uid uuid.UUID
		mustNewUUID(&uid)
		c.RequestID = uid
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
```

The problem? `http.Error` writes directly. What if this was a JSON api, or a gRPC based API? Writing a plain text error is a problem. Or alternatively, you need to write an exclusive error handling method that's used across everywhere that has intimate knowledge of the pipeline path.

Now, using `mchain` handlers:

```go
func RequestIDMustInitHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		c := FromRequest(r)
		if _, ok := r.Header[RequestIDHeaderKey]; ok {
			msg := fmt.Sprintf("error: illegal header (%s)", RequestIDHeaderKey)
			return errors.New(msg)
			// However, a better way would be to use the
			// goerror package that communicates error
			// along with status codes, in a clean way.
			//
			// return httperror.New(400, msg, true)
		}
		var uid uuid.UUID
		mustNewUUID(&uid)
		c.RequestID = uid
		return next.ServeHTTP(w, r)
	}
	return mchain.HandlerFunc(f)
}
```

Now, the errors can be handled up the middleware chain with an error handler that knows how to format the error the way it has to. Works naturally with the chain, without thinking about how to handle the error in every aspect of the middleware - when in doubt, pass it up the chain.
