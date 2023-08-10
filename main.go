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

func newSpackApp() *cli.App {
	var configFile string
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
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "location of the config file",
				Destination: &configFile,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "pack",
				Aliases: []string{"p"},
				Usage:   "packs a Solidity struct",
				Action: func(c *cli.Context) error {
					appConfig, err := NewAppSettings(configFile, readFromFile, unpacked, c.Args())
					if err != nil {
						return err
					}
					result, err := pack(&appConfig)
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
					appConfig, err := NewAppSettings(configFile, readFromFile, unpacked, c.Args())
					if err != nil {
						return err
					}
					result, err := count(&appConfig)
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

type AppSettings struct {
	readFromFile bool
	unpacked     bool
	printer      *printer.Printer
	args         cli.Args
}

func NewAppSettings(configFile string, readFromFile, unpacked bool, args cli.Args) (AppSettings, error) {
	configuration := config.GetDefaultConfig()
	// If the user specified a config file, load it
	if configFile != "" {
		globalConfig, err := config.LoadConfigFromFile(configFile)
		if err != nil {
			return AppSettings{}, err
		}
		configuration = globalConfig
	}

	newPrinter, err := printer.NewPrinter(configuration.PrintingConfig)
	if err != nil {
		return AppSettings{}, err
	}

	return AppSettings{
		printer:      &newPrinter,
		args:         args,
		readFromFile: readFromFile,
		unpacked:     unpacked,
	}, nil
}

func pack(settings *AppSettings) (string, error) {
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

func count(settings *AppSettings) (int, error) {
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

func getStruct(settings *AppSettings) (solidity.Struct, error) {
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
