package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"spack/config"
	"spack/parser"
	"spack/printer"
	"spack/solidity"
)

func main() {
	if err := newSpackApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type SpackApp struct {
	readFromFile bool
	unpacked     bool
	printer      *printer.Printer
	args         cli.Args
}

func newSpackApp() *cli.App {
	newPrinter, err := printer.NewPrinter(config.GetDefaultPrintingConfig())
	if err != nil {
		panic(err)
	}

	var readFromFile, unpacked bool
	app := &cli.App{
		Name:  "Spack",
		Usage: "pack Solidity structs",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "loads a Solidity struct from a file",
				Destination: &readFromFile,
			},
			&cli.BoolFlag{
				Name:        "unpacked",
				Aliases:     []string{"u"},
				Usage:       "does not pack the struct",
				Destination: &unpacked,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "pack",
				Aliases: []string{"p"},
				Usage:   "packs a Solidity struct",
				Action: func(c *cli.Context) error {
					AppConfig := &SpackApp{
						printer:      &newPrinter,
						args:         c.Args(),
						readFromFile: readFromFile,
						unpacked:     unpacked,
					}

					result, err := pack(AppConfig)
					if err != nil {
						return err
					}
					fmt.Println(result)
					return nil
				},
			},
			{
				Name:    "count",
				Aliases: []string{"c"},
				Usage:   "count the slots of the given struct",
				Action: func(c *cli.Context) error {
					AppConfig := &SpackApp{
						printer:      &newPrinter,
						args:         c.Args(),
						readFromFile: readFromFile,
						unpacked:     unpacked,
					}

					result, err := count(AppConfig)
					if err != nil {
						return err
					}
					fmt.Println(result)
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}

func pack(settings *SpackApp) (string, error) {
	solidityStruct, err := getStruct(settings)
	if err != nil {
		return "", errors.Wrap(err, "Error parsing struct")
	}
	if settings.unpacked {
		solidityStruct.StorageSlots = packStructCurrentFieldOrder(solidityStruct.Fields)
		return settings.printer.PrintSolidityStruct(solidityStruct), nil
	}

	solidityStruct.StorageSlots = packStructOptimal(solidityStruct.Fields)

	return settings.printer.PrintSolidityStruct(solidityStruct), nil
}

func count(settings *SpackApp) (int, error) {
	structDef, err := getStruct(settings)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing struct")
	}
	if settings.unpacked {
		structDef.StorageSlots = packStructCurrentFieldOrder(structDef.Fields)
		return len(structDef.StorageSlots), nil
	}

	structDef.StorageSlots = packStructOptimal(structDef.Fields)
	return len(structDef.StorageSlots), nil
}

func getStruct(settings *SpackApp) (solidity.Struct, error) {
	input := settings.args.Get(0)
	if input == "" {
		return solidity.Struct{}, errors.New("No input specified")
	}

	structString := input
	if settings.readFromFile {
		fileByes, err := os.ReadFile(input)
		if err != nil {
			panic(err)
		}
		structString = string(fileByes)
	}

	structDef, err := parser.ParseStruct(structString)
	if err != nil {
		return solidity.Struct{}, errors.Wrap(err, "Error parsing struct")
	}
	return structDef, nil
}
