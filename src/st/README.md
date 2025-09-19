# Symbol Table implementation in Golang
This Go module implements a simple SymbolTable for the subject PDL at UPM.

## Download
You can add the module to your project with `go get github.com/strozz1/go-st`.

Or cloning the repo.

## Usage
The entrypoint of the module is the STManager. The manager will be in charge of
creating the SymbolTable, defining attributes, adding symbols to the Global table,
creating new Scopes with the Local table, writing to file,...

The specs of the assignment requires only 2 scopes: global and local, so the STManager
is wrote to fulfill those requirements.

## Writing to file
When you create a STManager, you will be asked for a `io.Writer`. This can either be the `os.Stdout`
or a file to persist the result.

Calling `STManager.Write(table)` will write the table to the output defined by the Writer.



