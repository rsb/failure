# Failure
The failure package builds upon the [errors](https://pkg.go.dev/errors) 
package, implementing a strategy called ``Opaque errors`` which I first learned 
about from an article [Don't just check errors handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
by [Dave Cheney](https://dave.cheney.net).

The failure package concentrates on an api that categories errors to better 
articulate your intent in code.

## Usage

### Server Failure
Describe general server errors. If I don't have a specific category for the 
error I use this. When mapping to `http` errors this would map to `500`

### System Failure
Means the same thing as `Server` but you prefer the name `System` instead. Note 
that you can not mix and match, meaning `IsServer` will not return true for a 
`System`

### Startup
Is used to indicate failures that prevented the system from starting up

### Shutdown
Is used to signal a shutdown of the system. 

### Defer
Categorize errors that originated inside a `defer` call

### BadRequest
Categories errors that occurred because the accepted request was bad in some way

### Validation
Describes an error that occurred because some validation failed

### Config
Describes a failure caused by invalid configuration

### InvalidParam
Describes a failure cause by a bad function param or struct field

### NotAuthorized, NotAuthenticated, Forbidden
Describe auth errors

### NotFound
Describes a failure due to the absence of a resource


### Multiple
This is a direct port of [hashicorp multierror](https://github.com/hashicorp/go-multierror). Many thanks
to the hard work and great code produced my the hashicorp team. I integrated their code into this 
codebase to produce a seamless api when using multiple errors in the failure package.

NOTE: if you're using this package just because of multiple errors then I would use the hashicorps instead.
I simply repurposed their code to fit into the failure system because I liked the way it reads.

The `Append` function is used to create a list of errors. This function
behaves a lot like the Go built-in `append` function: it doesn't matter
if the first argument is nil, a `failure.Error`, or any other `error`,
the function behaves as you would expect.

```go
var result error

if err := step1(); err != nil {
	result = failure.Append(result, err)
}
if err := step2(); err != nil {
	result = failure.Append(result, err)
}

return result
```

#### Customizing the formatting of the errors

By specifying a custom `ErrorFormat`, you can customize the format
of the `Error() string` function:

```go
var result *failure.Multi

// ... accumulate errors here, maybe using Append

if result != nil {
	result.Formatter = func([]error) string {
		return "errors!"
	}
}
```

#### Extracting an error

```go
// Assume err is a failure value
err := somefunc()

if failure.IsMultiple(err) { 
	// It has it, and now errRich is populated.
}

var result []error

result, ok := failure.MultipleResult(err)
if ok {
	// Result will be []error in the failure.Multi
}


```

### Timeout
Describes failures that occurred because something took too long


## General Usage
```go
  func(db *Client) Insert(ctx context.Context, model business.Model) error {
    ...

		if !db.Validatea(model) {
		  return failure.Validation("Model is not valid")	
		}
		
		out, err := db.api.Put(ctx, in)
		if err != nil {
			return failure.ToSystem(err, "db.api.Put failed for (%s)", in.ID)
		} 
		
		...
  }
	
```

```go 
	
func (h *Handler) Handle(ctx context.Context, w http.ResponseWrite, r http.Request) error {
		...
		if err = db.Insert(ctx, model); err != nil {
			switch {
			case: failure.IsValidation(err):
				// Return 400 
			default: 
				// return 500
			}
		}
	}
```