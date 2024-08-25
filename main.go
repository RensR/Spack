package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/rensr/spack/config"
	"github.com/rensr/spack/parser"
	"github.com/rensr/spack/printer"
	"github.com/rensr/spack/solidity"
)

func main() {
	if err := newSpackApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newSpackApp() *cli.App {
	var contract, solidityStruct, configFile string
	var unpacked bool
	app := &cli.App{
		Name:  "Spack",
		Usage: "pack Solidity structs",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "contract",
				Aliases:     []string{"c"},
				Usage:       "loads all Solidity structs from a contract",
				Destination: &contract,
			},
			&cli.StringFlag{
				Name: "struct",
				Aliases: []string{"s"},
				Usage: "loads a single Solidity struct from a contract",
				Destination: &solidityStruct,
			},
			&cli.BoolFlag{
				Name:        "unpacked",
				Aliases:     []string{"u"},
				Usage:       "does not pack the struct",
				Destination: &unpacked,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"cfg"},
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
					appConfig, err := NewAppSettings(configFile, contract, solidityStruct, unpacked, c.Args())
					if err != nil {
						return err
					}
					result, err := pack(&appConfig)
					if err != nil {
						return err
					}

					for _, r := range result {
						fmt.Println(r)
					}

					return nil
				},
			},
			{
				Name:    "count",
				Aliases: []string{"c"},
				Usage:   "count the slots of the given struct",
				Action: func(c *cli.Context) error {
					appConfig, err := NewAppSettings(configFile, contract, solidityStruct, unpacked, c.Args())
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
	outDir	   string
	contract     string
	solidityStruct string
	unpacked     bool
	printer      *printer.Printer
	args         cli.Args
}

func NewAppSettings(configFile string, contract string, solidityStruct string, unpacked bool, args cli.Args) (AppSettings, error) {
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
		outDir: configuration.OutDir,
		printer:      &newPrinter,
		args:         args,
		contract: contract,
		solidityStruct: solidityStruct,
		unpacked:     unpacked,
	}, nil
}

func pack(settings *AppSettings) ([]string, error) {
	solidityStructs, err := getStructs(settings)
	if err != nil {
		return []string{}, errors.Wrap(err, "Error parsing struct")
	}

	var results []string

	for _, solidityStruct := range solidityStructs {
		if settings.unpacked {
			solidityStruct.StorageSlots = packStructCurrentFieldOrder(solidityStruct.Fields)
			results = append(results, settings.printer.PrintSolidityStruct(solidityStruct))
		}

		solidityStruct.StorageSlots = packStructOptimal(solidityStruct.Fields)

		results = append(results, settings.printer.PrintSolidityStruct(solidityStruct))
	}

	return results, nil
}

func count(settings *AppSettings) ([]int, error) {
	structDef, err := getStructs(settings)
	if err != nil {
		return []int{}, errors.Wrap(err, "Error parsing struct")
	}

	var slots []int

	for _, structDef := range structDef {
		if settings.unpacked {
			structDef.StorageSlots = packStructCurrentFieldOrder(structDef.Fields)
			slots = append(slots, len(structDef.StorageSlots))
		}

		structDef.StorageSlots = packStructOptimal(structDef.Fields)
		slots = append(slots, len(structDef.StorageSlots))
	}

	return slots, nil
}

func getStructs(settings *AppSettings) ([]solidity.Struct, error) {
	cmd := exec.Command("forge",  "build",  "--ast")

    _, err := cmd.Output()

	if err != nil {
		return []solidity.Struct{}, errors.Wrap(err, "Failed to run forge compile")
	}
	
	if settings.contract == "" {
		return []solidity.Struct{}, errors.New("No input specified")
	}

	rawData, err := os.ReadFile((settings.outDir + settings.contract + ".sol/" + settings.contract + ".json"))

	if err != nil {
		return []solidity.Struct{}, errors.Wrap(err, "Failed to read file")
	}

	// Create a map to extract the AST data
	var data map[string]interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return []solidity.Struct{}, errors.New("Failed to parse AST")
	}

	// Extract the AST data
	astData, err := json.Marshal(data["ast"])
	if err != nil {
		return []solidity.Struct{}, errors.New("Failed to extract AST data")
	}

	// Cast astData to SolidityAST
	var ast parser.SolidityAST
	err = json.Unmarshal(astData, &ast)

	if err != nil {
		return []solidity.Struct{}, errors.New("Failed to cast AST data")
	}

	s, err := ast.ParseStructs(settings.solidityStruct)
	if err != nil {
		return []solidity.Struct{}, errors.Wrap(err, "Failed to parse structs")
	}

	return s, nil
}
