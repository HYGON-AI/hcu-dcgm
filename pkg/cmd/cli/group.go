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
	listFlag     bool
	deleteId     string
	createName   string
	addEntity    string
	removeEntity string
	groupId      string
	infoFlag     bool
	jsonFlag     bool
	defaultFlag  bool
)

var entityGroupNameToId = map[string]dcgm.Field_Entity_Group{
	"hcu":    dcgm.FE_HCU,
	"vhcu":   dcgm.FE_VHCU,
	"switch": dcgm.FE_SWITCH,
	"gi":     dcgm.FE_HCU_GI,
	"ci":     dcgm.FE_HCU_CI,
	"link":   dcgm.FE_LINK,
	"cpu":    dcgm.FE_CPU,
	"core":   dcgm.FE_CPU_CORE,
}

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Used to create and maintain groups of HCUs.",
	Long: `group -- Used to create and maintain groups of HCUs. Groups of HCUs can then be
uniformly controlled through other DCGMI subsystems.`,
	Example: `  dcgmi group -l -j
  dcgmi group -c <groupName> --default
  dcgmi group -c <groupName> -a <entityId>
  dcgmi group -d <groupId>
  dcgmi group -g <groupId> -i -j
  dcgmi group -g <groupId> -a <entityId>
  dcgmi group -g <groupId> -r <entityId>`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main dispatcher logic
		switch {
		case listFlag:
			handleList()
		case createName != "":
			handleCreate()
		case deleteId != "":
			handleDelete()
		case groupId != "":
			handleGroupOperations()
		default:
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)

	groupCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List the groups that currently exist for the host.")
	groupCmd.Flags().StringVarP(&deleteId, "delete", "d", "", "Delete a group on the host.")
	groupCmd.Flags().StringVarP(&createName, "create", "c", "", "Create a group on the host.")
	groupCmd.Flags().StringVarP(&addEntity, "add", "a", "", "Add device(s) to group. (csv hcuIds like 0,1 or entityIds like hcu:0,hcu:1)")
	groupCmd.Flags().StringVarP(&removeEntity, "remove", "r", "", "Remove device(s) from group. (csv hcuIds like 0,1 or entityIds like hcu:0,hcu:1)")
	groupCmd.Flags().StringVarP(&groupId, "group", "g", "", "Specify the group ID to operate on.")
	groupCmd.Flags().BoolVarP(&infoFlag, "info", "i", false, "Display the information for the specified group ID.")
	groupCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Print the output in a json format.")
	groupCmd.Flags().BoolVar(&defaultFlag, "default", false, "Adds all available HCUs to the group being created.")
}

func formatEntities(entityList []dcgm.GroupEntityPair) {
	if len(entityList) == 0 {
		fmt.Printf("|    -> Entities    | None                                                     |\n")
	} else {
		entityStr := ""
		for i, entity := range entityList {
			if i > 0 {
				entityStr += ", "
			}
			entityStr += fmt.Sprintf("%s %d", entity.EntityGroupId.String(), entity.EntityId)
		}
		// If entity string is long, wrap it (optional, simple split)
		maxWidth := 56
		if len(entityStr) <= maxWidth {
			fmt.Printf("|    -> Entities    | %-56s |\n", entityStr)
		} else {
			fmt.Printf("|    -> Entities    | %-56s |\n", entityStr[:maxWidth])
			remaining := entityStr[maxWidth:]
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
func handleList() {
	groups, err := dcgm.ListAllGroups()
	if err != nil {
		fmt.Printf("Error retrieving group list: %v\n", err)
		return
	}
	if jsonFlag {
		json.NewEncoder(os.Stdout).Encode(groups)
		return
	}

	fmt.Println("+-------------------+----------------------------------------------------------+")
	fmt.Println("| GROUPS                                                                       |")
	if len(groups) == 1 {
		fmt.Printf("|  1 group found.                                                              |\n")
	} else {
		fmt.Printf("| %-3dgroups found.                                                             |\n", len(groups))
	}
	fmt.Println("+===================+==========================================================+")
	fmt.Printf("| Groups            |                                                          |\n")

	for _, group := range groups {
		fmt.Printf("| -> %-15d|                                                          |\n", group.GroupId)
		fmt.Printf("|    -> Group ID    | %-56d |\n", group.GroupId)
		fmt.Printf("|    -> Group Name  | %-56s |\n", group.GroupName)
		formatEntities(group.EntityList)
	}
	fmt.Println("+-------------------+----------------------------------------------------------+")
}

// parseEntityList parses a CSV string of entity specs into []GroupEntityPair
func parseEntityList(input string) (entities []dcgm.GroupEntityPair, err error) {
	items := strings.Split(input, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		var entityType string
		var entityIdStr string

		if strings.Contains(item, ":") {
			// Explicit type like hcu:0
			parts := strings.SplitN(item, ":", 2)
			entityType = strings.ToLower(parts[0])
			entityIdStr = parts[1]
		} else {
			entityType = "hcu"
			entityIdStr = item
		}

		entityId, err := strconv.Atoi(entityIdStr)
		if err != nil {
			return nil, fmt.Errorf("invalid entity ID in '%s': %v\n", item, err)
		}

		groupID, ok := entityGroupNameToId[entityType]
		if !ok {
			return nil, fmt.Errorf("invalid entity type '%s' in '%s'\n", entityType, item)
		}

		entities = append(entities, dcgm.GroupEntityPair{
			EntityGroupId: groupID,
			EntityId:      entityId,
		})
	}
	return entities, nil
}

func addEntityToGroup(addEntity string, groupId int) {
	entityList, err := parseEntityList(addEntity)
	if err != nil {
		fmt.Printf("Invalid --add input: %v\n", err)
		return
	}
	for _, entity := range entityList {
		if entity.EntityGroupId == dcgm.FE_HCU {
			_, err := dcgm.GetDeviceId(entity.EntityId)
			if err != nil {
				fmt.Printf("Invalid --add input: HCU %d, %v\n", entity.EntityId, err)
				return
			}
		}
	}
	err = dcgm.AddEntityToGroup(groupId, entityList)
	if err != nil {
		fmt.Printf("Failed to add entity to group: %v\n", err)
		return
	}
	fmt.Printf("Successfully added entities to group %d\n", groupId)
}

func handleCreate() {
	if defaultFlag {
		groupId, err := dcgm.CreateDefaultGroup(createName)
		if err != nil {
			fmt.Printf("Failed to create default group: %v\n", err)
		} else {
			fmt.Printf("Successfully created group \"%s\" with a group ID of %d\n", createName, groupId)
		}
		return
	}
	groupId, err := dcgm.CreateGroup(createName)
	if err != nil {
		fmt.Printf("Failed to create group: %v\n", err)
		return
	}
	fmt.Printf("Successfully created group \"%s\" with a group ID of %d\n", createName, groupId)

	if addEntity != "" {
		addEntityToGroup(addEntity, groupId)
	}
}

func handleDelete() {
	groupId, err := strconv.Atoi(deleteId)
	if err != nil {
		fmt.Printf("Invalid Group ID\n")
		return
	}
	err = dcgm.DestroyGroup(groupId)
	if err != nil {
		fmt.Printf("Failed to delete group: %v\n", err)
		return
	}
	fmt.Printf("Successfully removed group %s\n", deleteId)
}

func handleGroupOperations() {
	groupId, err := strconv.Atoi(groupId)
	if err != nil {
		fmt.Println("Invalid Group ID specified.")
		return
	}
	groupInfo, err := dcgm.GetGroupInfo(groupId)
	if err != nil {
		fmt.Printf("Failed to get group info: %v\n", err)
		return
	}

	if infoFlag && jsonFlag {
		json.NewEncoder(os.Stdout).Encode(groupInfo)
		return
	}
	if infoFlag {
		fmt.Println("+-------------------+----------------------------------------------------------+")
		fmt.Println("| GROUP INFO                                                                   |")
		fmt.Println("+===================+==========================================================+")
		fmt.Printf("| -> %-15d|                                                          |\n", groupId)
		fmt.Printf("|    -> Group ID    | %-56d |\n", groupInfo.GroupId)
		fmt.Printf("|    -> Group Name  | %-56s |\n", groupInfo.GroupName)
		formatEntities(groupInfo.EntityList)
		fmt.Println("+-------------------+----------------------------------------------------------+")
		return
	}
	if addEntity != "" {
		addEntityToGroup(addEntity, groupId)
		return
	}
	if removeEntity != "" {
		entityList, err := parseEntityList(removeEntity)
		if err != nil {
			fmt.Printf("Invalid --remove input: %v\n", err)
			return
		}
		err = dcgm.RemoveEntityFromGroup(groupId, entityList)
		if err != nil {
			fmt.Printf("Failed to remove entity from group: %v\n", err)
			return
		}
		fmt.Printf("Successfully removed entities from group %d\n", groupId)
	}
}
