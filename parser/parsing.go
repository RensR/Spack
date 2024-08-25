package parser

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rensr/spack/solidity"
)

type SolidityAST struct {
	Nodes []ASTNode `json:"nodes"`
}

type ASTNode struct {
	NodeType string `json:"nodeType"`
	Name    string `json:"name"`
	Nodes []ContractNode `json:"nodes"`
}

type ContractNode struct {
	NodeType string `json:"nodeType"`
	Name   string `json:"name"`
	Nodes []Member `json:"members"`
}

type Member struct {
	NodeType string `json:"nodeType"`
	TypeName TypeName `json:"typeName"`
	Name string `json:"name"`
}

type TypeName struct {
	NodeType string `json:"nodeType"`
	TypeDescriptions TypeDescription `json:"typeDescriptions"`
}

type TypeDescription struct {
	Type string `json:"typeString"`
}

// Parses the Solidity AST object and returns a list of structs
func (s *SolidityAST) ParseStructs(solidityStruct string) ([]solidity.Struct, error) {
	var structs []solidity.Struct
		for _, node := range s.Nodes {
			if node.NodeType == "ContractDefinition" {
				singleLookup := solidityStruct != ""
				for _, contractNode := range node.Nodes {
					if contractNode.NodeType == "StructDefinition" {
						if singleLookup && contractNode.Name != solidityStruct {
							continue
						}
						ss := solidity.Struct{Name: contractNode.Name}
						for _, member := range contractNode.Nodes {
							size, err := member.ComputeSize()
							if err != nil  {
								return nil, errors.Wrap(err, "Failed to compute size")
							}
							if strings.Contains(member.TypeName.TypeDescriptions.Type, "struct") {
								member.TypeName.TypeDescriptions.Type = strings.Trim(member.TypeName.TypeDescriptions.Type, "struct " + node.Name + ".")
							}
						
							ss.Fields = append(ss.Fields, solidity.DataDef{Name: member.Name, Comment: "", Type: solidity.DataType(member.TypeName.TypeDescriptions.Type), Size: size})
						}
						structs = append(structs, ss)
					}
				}
			}
		}
	return structs, nil
}

// Computes the size of a member
// For dynamic types such as strings, bytes and arrays, it returns 32
func (m *Member) ComputeSize() (uint8, error) {
	if m.NodeType == "UserDefinedTypeName" {
		return 32, nil
	}

	switch m.TypeName.TypeDescriptions.Type {
		case "bool": return 1, nil
		case "address": return 20, nil
		case "string": return 32, nil
		case "bytes": return 32, nil
		default:
			if strings.Contains(m.TypeName.TypeDescriptions.Type, "[]") {
				return 32, nil
			}
			if strings.Contains(m.TypeName.TypeDescriptions.Type, "uint") {
				sizeString := strings.Trim(m.TypeName.TypeDescriptions.Type, "uint")
				size, err := strconv.ParseUint(sizeString, 10, 64)
				if err != nil {
					return 0, err
				}
				return uint8(size) / 8, nil
			}
			if strings.Contains(m.TypeName.TypeDescriptions.Type, "int") {
				sizeString := strings.Trim(m.TypeName.TypeDescriptions.Type, "int")
				size, err := strconv.ParseUint(sizeString, 10, 64)
				if err != nil {
					return 0, err
				}
				return uint8(size) / 8, nil
			}
			if strings.Contains(m.TypeName.TypeDescriptions.Type, "bytes") {
				sizeString := strings.Trim(m.TypeName.TypeDescriptions.Type, "bytes")
				size, err := strconv.ParseUint(sizeString, 10, 64)
				if err != nil {
					return 0, err
				}
				return uint8(size), nil
			}
			return 32, nil
	}
}