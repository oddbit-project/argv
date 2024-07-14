package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oddbit-project/argv"
	"os"
	"strings"
)

type AlgorithmDetails struct {
	Algorithm string `argv:"alg"`
	KeyLen    uint   `argv:"bits"`
}

type CertInfo struct {
	CommonName         string `argv:"CN"`
	OrganizationalUnit string `argv:"OU"`
	Organization       string `argv:"O,optional"`
	Algorithm          AlgorithmDetails
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
		// destination struct
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
