/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"encoding/json"
	"fmt"
	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var (
	deleteFlag   bool
	fieldIds     string
	fieldGroupId string
)

var fieldGroupCmd = &cobra.Command{
	Use:   "fieldgroup",
	Short: "Used to create and maintain groups of field IDs.",
	Long: `fieldgroup -- Used to create and maintain groups of field IDs. Groups of field
 IDs can then be uniformly controlled through other DCGMI subsystems.`,
	Example: `  dcgmi fieldgroup -l -j
  dcgmi fieldgroup -c <fieldGroupName> -f <fieldIds>
  dcgmi fieldgroup -i -g <fieldGroupId> -j
  dcgmi fieldgroup -d -g <fieldGroupId>`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main dispatcher logic
		switch {
		case listFlag:
			handleFieldGroupList()
		case createName != "":
			handleFieldGroupCreate()
		case deleteFlag:
			handleFieldGroupDelete()
		case infoFlag:
			handleFieldGroupInfo()
		default:
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(fieldGroupCmd)

	fieldGroupCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List the field groups that currently exist for the host.")
	fieldGroupCmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "Delete a field group on the host.")
	fieldGroupCmd.Flags().StringVarP(&createName, "create", "c", "", "Create a field group on the host.")
	fieldGroupCmd.Flags().StringVarP(&fieldIds, "fieldids", "f", "", "Comma-separated list of the field ids to add to a field group when creating a new one.")
	fieldGroupCmd.Flags().StringVarP(&fieldGroupId, "fieldgroup", "g", "", "Specify the field group ID to operate on.")
	fieldGroupCmd.Flags().BoolVarP(&infoFlag, "info", "i", false, "Display the information for the specified field group ID.")
	fieldGroupCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Print the output in a json format.")
}

func formatFieldIds(fieldIds []int) {
	if len(fieldIds) == 0 {
		fmt.Printf("|    -> Field IDs   | None                                                     |\n")
	} else {
		fieldIdStr := ""
		for i, fieldId := range fieldIds {
			if i > 0 {
				fieldIdStr += ", "
			}
			fieldIdStr += fmt.Sprintf("%d", fieldId)
		}
		maxWidth := 56
		if len(fieldIdStr) <= maxWidth {
			fmt.Printf("|    -> Field IDs   | %-56s |\n", fieldIdStr)
		} else {
			fmt.Printf("|    -> Field IDs   | %-56s |\n", fieldIdStr[:maxWidth])
			remaining := fieldIdStr[maxWidth:]
			for len(remaining) > 0 {
				line := remaining
				if len(line) > maxWidth {
					line = remaining[:maxWidth]
				}
				fmt.Printf("|                   | %-56s |\n", line)
				if len(remaining) > maxWidth {
					remaining = remaining[maxWidth:]
				} else {
					break
				}
			}
		}
	}
}

// Handlers for subcommands
func handleFieldGroupList() {
	fieldGroups, err := dcgm.ListAllFieldGroups()
	if err != nil {
		fmt.Printf("Error retrieving field group list: %v\n", err)
		return
	}
	if jsonFlag {
		json.NewEncoder(os.Stdout).Encode(fieldGroups)
		return
	}

	fmt.Println("+-------------------+----------------------------------------------------------+")
	fmt.Println("| FIELD GROUPS                                                                 |")
	if len(fieldGroups) == 1 {
		fmt.Printf("| 1 field group found.                                                         |\n")
	} else {
		fmt.Printf("| %-3dfield groups found.                                                       |\n", len(fieldGroups))
	}
	fmt.Println("+===================+==========================================================+")
	fmt.Printf("| Field Groups      |                                                          |\n")

	for _, fieldGroup := range fieldGroups {
		fmt.Printf("| -> %-15d|                                                          |\n", fieldGroup.FieldGroupId)
		fmt.Printf("|    -> ID          | %-56d |\n", fieldGroup.FieldGroupId)
		fmt.Printf("|    -> Name        | %-56s |\n", fieldGroup.FieldGroupName)
		formatFieldIds(fieldGroup.FieldIds)
	}
	fmt.Println("+-------------------+----------------------------------------------------------+")
}

// parseFieldIds parses a CSV string of ids into []int
func parseFieldIds(input string) (fieldIdList []int, err error) {
	items := strings.Split(input, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		fieldId, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("invalid field ID in '%s': %v\n", item, err)
		}
		fieldIdList = append(fieldIdList, fieldId)
	}
	return fieldIdList, nil
}

func handleFieldGroupCreate() {
	if fieldIds == "" {
		fmt.Println("Error: No field IDs given (specify with -f or --fieldids) for arg create.")
		return
	}
	fieldIdList, err := parseFieldIds(fieldIds)
	if err != nil {
		fmt.Printf("Failed to create field group: %v\n", err)
		return
	}
	fieldGroupId, err := dcgm.CreateFieldGroup(createName, fieldIdList)
	if err != nil {
		fmt.Printf("Failed to create field group: %v\n", err)
		return
	}
	fmt.Printf("Successfully created field group \"%s\" with a field group ID of %d\n", createName, fieldGroupId)
}

func handleFieldGroupDelete() {
	if fieldGroupId == "" {
		fmt.Println("Error: No fieldGroupId given (specify with -g or --fieldgroup) for arg delete.")
		return
	}
	fieldGroupId, err := strconv.Atoi(fieldGroupId)
	if err != nil {
		fmt.Printf("Invalid field group ID\n")
		return
	}
	err = dcgm.DestroyFieldGroup(fieldGroupId)
	if err != nil {
		fmt.Printf("Failed to delete field group: %v\n", err)
		return
	}
	fmt.Printf("Successfully removed field group %d\n", fieldGroupId)
}

func handleFieldGroupInfo() {
	if fieldGroupId == "" {
		fmt.Println("Error: No fieldGroupId given (specify with -g or --fieldgroup) for arg info.")
		return
	}
	fieldGroupId, err := strconv.Atoi(fieldGroupId)
	fieldGroupInfo, err := dcgm.GetFieldGroupInfo(fieldGroupId)
	if err != nil {
		fmt.Printf("Failed to get field group info: %v\n", err)
		return
	}

	if jsonFlag {
		json.NewEncoder(os.Stdout).Encode(fieldGroupInfo)
		return
	}
	if infoFlag {
		fmt.Println("+-------------------+----------------------------------------------------------+")
		fmt.Println("| FIELD GROUP INFO                                                             |")
		fmt.Println("+===================+==========================================================+")
		fmt.Printf("| -> %-15d|                                                          |\n", fieldGroupId)
		fmt.Printf("|    -> ID          | %-56d |\n", fieldGroupInfo.FieldGroupId)
		fmt.Printf("|    -> Name        | %-56s |\n", fieldGroupInfo.FieldGroupName)
		formatFieldIds(fieldGroupInfo.FieldIds)
		fmt.Println("+-------------------+----------------------------------------------------------+")
		return
	}
}
