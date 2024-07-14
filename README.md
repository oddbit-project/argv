# Argv

github.com/oddbit-project/argv

[![GitHub tag](https://img.shields.io/github/tag/oddbit-project/argv.svg?style=flat)](https://github.com/oddbit-project/argv/releases)

Parse string slices into structs

## Usage scenario

argv attempts to parse a string slice, composed of parameters and values, and fill a struct based on its annotation
information and converting values to the specific field type.

The typical use case is to convert an arbitrary list of arguments and values to a well-defined struct:

Sample main.go:

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oddbit-project/argv"
	"os"
	"strings"
)

type CertInfo struct {
	CommonName         string `argv:"CN"`
	OrganizationalUnit string `argv:"OU"`
	Organization       string `argv:"O,optional"` // optional parameter
	Algorithm          string `argv:"alg"`
	KeyLen             uint   `argv:"bits"`
	Days               uint32 `argv:"days"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command> [options]\n", os.Args[0])
		os.Exit(0)
	}
	cmd := os.Args[1]

	switch cmd {
	case "gencert":
		record := &CertInfo{}

		// attempt to parse & serialize values from command line
		err := argv.ParseArgv(record, os.Args[2:])
		if err == nil {
			str, _ := json.Marshal(record)
			fmt.Printf("Parsed parameters:\n %s\n", string(str))
			os.Exit(0)
		}

		// empty arg list, show available names
		if errors.Is(err, argv.ErrEmptyArgs) {
			names, err := argv.ParseNames(record)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Println("Available args: ", strings.Join(names, ", "))
			os.Exit(0)
		}

		fmt.Println(err)
		os.Exit(-1)
	default:
		fmt.Printf("'%s' is not a supported command\n", cmd)
		os.Exit(-1)
	}
}
```

Running:

```shell
$ go run main.go gencert
empty argument list
exit status 255

$ go run main.go gencert 
Available args:  CN, OU, O,optional, alg, bits, days

$ go run main.go gencert -CN "certificate name" -OU "SomeOrg"
value for arg 'alg' is missing
exit status 255

$  go run main.go gencert -CN "certificate name" -OU "SomeOrg" -alg "rsa" -bits 4096 -days 365
Parsed parameters:
 {"CommonName":"certificate name","OrganizationalUnit":"SomeOrg","Organization":"","Algorithm":"rsa","KeyLen":4096,"Days":365}
```

## Supported field types

| type      | description                                 |
|-----------|---------------------------------------------|
| int       | Int value                                   |
| int8      | Int8 value                                  |
| int32     | Int32 value                                 |
| int64     | Int64 value                                 |
| uint      | Uit value                                   |
| uint8     | Uint8 value                                 |
| uint32    | Uint32 value                                |
| uint64    | Uint64 value                                |
| float32   | float32 value; supports scientific notation |
| float64   | float64 value; supports scientific notation |
| time.Time | RFC3339 time string                         |
| bool      | true/false or 1/0 string                    |
| string    | arbitrary string                            |
| []string  | list of strings, such as "value1,value2"    |


