# `goaccessor` Design Doc

Basic the core code can be separated into:

- `task`
- parse
- generate

The code of `goaccessor` will be executed as follow:

1. parse the 

## Core type

``` go
type task
```

## Parse

## Generate

## Limitation

### Multiple File Generation

In the current design of `goaccessor`, each `//go:generate goaccessor ...` directive generates a separate output file. This is due to the way the `go generate` command works: it executes each directive independently. Therefore, if a single `.go` file contains multiple `goaccessor` directives, each directive will result in a separate output file.

For example, if we have a file `types.go` with three `goaccessor` directives, the output would be three separate files.

### Potential Solutions

1. **Intermediate Output**: One possible solution to merge the outputs into a single file involves creating intermediate outputs for each `goaccessor` directive. A separate tool or script (`mergeaccessors`) could then be used to consolidate these intermediate outputs into a single file.

```go
//go:generate goaccessor Book --getter --setter --output intermediate1
type Book struct {
    Title  string
    Author string
}

//go:generate mergeaccessors
```

In this case, `mergeaccessors` would be responsible for taking all of the `intermediateN` files, merging them into a single Go source file, and then deleting the intermediate files.

### Conclusion

Given the goal of delivering a Minimum Viable Product (MVP), adding complex mechanisms such as intermediate outputs and merging is currently out of scope. As such, we have to tolerate this limitation for now. In future versions, depending on user feedback and the need for such functionality, we may consider implementing a solution to overcome this limitation. For now, users are advised to structure their `goaccessor` directives with this limitation in mind.
