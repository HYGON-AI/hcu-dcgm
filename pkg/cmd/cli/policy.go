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
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	getFlag       bool
	regFlag       bool
	setCsv        string
	clearFlag     bool
	maxPages      string
	maxTemp       string
	maxPower      string
	eccErrorFlag  bool
	pcieErrorFlag bool
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Used to control policies for groups of HCUs.",
	Long: `policy -- Used to control policies for groups of HCUs. Policies control actions
 which are triggered by specific events.`,
	Example: `  dcgmi policy --get -j
  dcgmi policy -g <groupId> --reg
  dcgmi policy -g <groupId> --set <actn,val> -M <maxRetiredPages> -T
        <maxTemp> -P <maxPower> -e -p 
  dcgmi policy -g <groupId> --clear`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main dispatcher logic
		switch {
		case getFlag:
			handleGet()
		case verboseFlag && !getFlag:
			fmt.Println("Verbose option only available with the get policy arg (--get)")
		case clearFlag:
			handleClear()
		case regFlag:
			handleReg()
		case setCsv != "":
			handleSet()
		default:
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)

	policyCmd.Flags().BoolVar(&getFlag, "get", false, "Get the current violation policy.")
	policyCmd.Flags().BoolVar(&regFlag, "reg", false, "Register this process for policy updates.  This\n"+
		"process will sit in an infinite loop waiting for\n updates from the policy manager.")
	policyCmd.Flags().StringVar(&setCsv, "set", "", "Set the current violation policy. Use csv action\n"+
		",validation (ie. 1,2) \n----- \nAction to take when any of the violations\n specified occur. \n0 - None \n1 - HCU Reset \n"+
		"----- \nValidation to take after the violation action has\nbeen performed. \n0 - None \n1 - System Validation (short) \n2 - System Validation (medium) \n3 - System Validation (long)")
	policyCmd.Flags().StringVarP(&groupId, "group", "g", "", "The group ID to query.")
	policyCmd.Flags().BoolVar(&clearFlag, "clear", false, "Clear the current violation policy.")
	policyCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Display policy information per HCU.")
	policyCmd.Flags().StringVarP(&maxPages, "maxpages", "M", "", "Specify the maximum number of retired pages that\nwill trigger a violation.")
	policyCmd.Flags().StringVarP(&maxTemp, "maxtemp", "T", "", "Specify the maximum temperature a group's HCUs can\nreach before triggering a violation.")
	policyCmd.Flags().StringVarP(&maxPower, "maxpower", "P", "", "Specify the maximum power a group's HCUs can reach\nbefore triggering a violation.")
	policyCmd.Flags().BoolVarP(&eccErrorFlag, "eccerrors", "e", false, "Add ECC errors to the policy\nconditions.")
	policyCmd.Flags().BoolVarP(&pcieErrorFlag, "pcierrors", "p", false, "Add PCIe replay errors to the policy conditions.")
	policyCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Print the output in a json format.")
}

func printPolicyInfo(header string, conditions []string, action string, validation string) {
	fmt.Printf("+-----------------------------+------------------------------------------------+\n"+
		"| Policy Information                                                           |\n"+
		"| %-28s                                                 |\n"+
		"+=============================+================================================+\n", header)
	for i := 0; i < len(conditions); i++ {
		if i == 0 {
			fmt.Printf("| Violation conditions        | %-47s|\n", conditions[i])
		} else {
			fmt.Printf("|                             | %-47s|\n", conditions[i])
		}
	}
	fmt.Printf("| Action on violation         | %-47s|\n", action)
	fmt.Printf("| Validation after action     | %-47s|\n", validation)
	fmt.Println("+-----------------------------+------------------------------------------------+")
}

func formatHcuList(groupId string) (hcuList []int, headerStr string, err error) {
	if groupId != "" {
		groupIdInt, err := strconv.Atoi(groupId)
		if err != nil {
			return hcuList, headerStr, err
		}
		hcuList, headerStr, err = dcgm.GetHCUListFromGroup(groupIdInt)
		if err != nil {
			return nil, "", err
		}
	} else {
		numHcus, err := dcgm.NumMonitorDevices()
		if err != nil {
			return hcuList, headerStr, fmt.Errorf("Error getting device count: %v\n", err)
		}
		hcuList = make([]int, numHcus)
		for i := 0; i < numHcus; i++ {
			hcuList[i] = i
		}
		headerStr = "DCGM_ALL_SUPPORTED_HCUS"
	}
	return hcuList, headerStr, nil
}

func formatConditionList(policyConditions dcgm.PolicyConditions) (conditionList []string) {
	if !(policyConditions.EccErrorsEnable || policyConditions.PcieErrorsEnable || policyConditions.MaxPagesEnable || policyConditions.MaxTempEnable || policyConditions.MaxPowerEnable) {
		return []string{"None"}
	}
	if policyConditions.EccErrorsEnable {
		conditionList = append(conditionList, "ECC errors")
	}
	if policyConditions.PcieErrorsEnable {
		conditionList = append(conditionList, "PCI errors and replays")
	}
	if policyConditions.MaxPagesEnable {
		conditionList = append(conditionList, fmt.Sprintf("Max retired pages threshold - %d", policyConditions.MaxPages))
	}
	if policyConditions.MaxTempEnable {
		conditionList = append(conditionList, fmt.Sprintf("Max temperature threshold - %.0f", policyConditions.MaxTemp))
	}
	if policyConditions.MaxPowerEnable {
		conditionList = append(conditionList, fmt.Sprintf("Max power threshold - %d", policyConditions.MaxPower))
	}
	return conditionList
}

func handleGet() {
	hcuList, headerStr, err := formatHcuList(groupId)
	if err != nil {
		fmt.Printf("Error getting hcu list: %v\n", err)
		return
	}
	if len(hcuList) == 0 {
		fmt.Printf("Failed to query group: no entity found for group %v\n", groupId)
		return
	}
	policyList, err := dcgm.GetPolicy(hcuList)
	if err != nil {
		fmt.Printf("Error getting policy: %v\n", err)
		return
	}
	if jsonFlag {
		json.NewEncoder(os.Stdout).Encode(policyList)
		return
	}

	fmt.Println("Policy information")
	if verboseFlag {
		for _, policyInfo := range policyList {
			hcuIndex := policyInfo.DvInd
			headerStr = fmt.Sprintf("HCU ID: %d", hcuIndex)
			policyConditions := policyInfo.Conditions
			conditionList := formatConditionList(policyConditions)
			action := policyInfo.ActionIndex.String()
			validation := policyInfo.ValidationIndex.String()
			printPolicyInfo(headerStr, conditionList, action, validation)
		}
	} else {
		tmpPolicyInfo := policyList[0]
		tmpPolicyConditions := tmpPolicyInfo.Conditions
		tmpConditionList := formatConditionList(tmpPolicyConditions)
		tmpAction := tmpPolicyInfo.ActionIndex.String()
		tmpValidation := tmpPolicyInfo.ValidationIndex.String()
		if len(policyList) > 1 {
			for i := 1; i < len(policyList); i++ {
				policyInfo := policyList[i]
				policyConditions := policyInfo.Conditions
				conditionList := formatConditionList(policyConditions)
				action := policyInfo.ActionIndex.String()
				validation := policyInfo.ValidationIndex.String()
				if tmpConditionList[0] != NOT_APPLICABLE && !slices.Equal(tmpConditionList, conditionList) {
					tmpConditionList = []string{NOT_APPLICABLE}
				}
				if tmpAction != NOT_APPLICABLE && tmpAction != action {
					tmpAction = NOT_APPLICABLE
				}
				if tmpValidation != NOT_APPLICABLE && tmpValidation != validation {
					tmpValidation = NOT_APPLICABLE
				}
			}
		}
		printPolicyInfo(headerStr, tmpConditionList, tmpAction, tmpValidation)
		fmt.Println("**** Non-homogenous settings across group. Use with –v flag to see details.")
	}
}

func handleClear() {
	hcuList, _, err := formatHcuList(groupId)
	if err != nil {
		fmt.Printf("Error getting hcu list: %v\n", err)
	}
	err = dcgm.ClearPolicy(hcuList)
	if err != nil {
		fmt.Printf("Error clearing policy: %v\n", err)
	}
	fmt.Println("Policy successfully set.")
}

func parseSetStr(input string) (action dcgm.PolicyAction, validation dcgm.PolicyValidation, err error) {
	inputSlice := []rune(input)
	errorStr := "Must use action,validation (csv) format for arg set."
	if len(inputSlice) != 3 {
		return 0, 0, fmt.Errorf(errorStr)
	}
	items := strings.Split(input, ",")
	actionInt, err := strconv.Atoi(items[0])
	if err != nil {
		return 0, 0, fmt.Errorf(errorStr)
	}
	action = dcgm.PolicyAction(actionInt)
	if action < dcgm.ACTION_NONE || action >= dcgm.ACTION_COUNT {
		return 0, 0, fmt.Errorf("The action must be between %d and %d (inclusive) for arg set.", dcgm.ACTION_NONE, dcgm.ACTION_COUNT-1)
	}

	validationInt, err := strconv.Atoi(items[1])
	if err != nil {
		return 0, 0, fmt.Errorf(errorStr)
	}
	validation = dcgm.PolicyValidation(validationInt)
	if validation < dcgm.VALIDATION_NONE || validation >= dcgm.VALIDATION_COUNT {
		return 0, 0, fmt.Errorf("The validation must be between %d and %d (inclusive) for arg set.", dcgm.VALIDATION_NONE, dcgm.VALIDATION_COUNT-1)
	}
	return action, validation, nil
}

func handleSet() {
	if maxPages == "" && maxTemp == "" && maxPower == "" && !eccErrorFlag && !pcieErrorFlag {
		fmt.Println("Error: No conditions specified for arg set.")
		return
	}
	var maxPagesInt, maxPowerInt int
	var maxTempFloat float64
	var maxPagesEnable, maxTempEnable, maxPowerEnable bool
	action, validation, err := parseSetStr(setCsv)
	if err != nil {
		fmt.Printf("Error parsing set string: %v\n", err)
		return
	}
	if maxPages != "" {
		maxPagesEnable = true
		maxPagesInt, err = strconv.Atoi(maxPages)
		if err != nil {
			fmt.Printf("Error parsing max pages: %v\n", err)
			return
		}
	} else {
		maxPagesEnable = false
	}
	if maxTemp != "" {
		maxTempEnable = true
		maxTempFloat, err = strconv.ParseFloat(maxTemp, 64)
		if err != nil {
			fmt.Printf("Error parsing max temp: %v\n", err)
			return
		}
	} else {
		maxTempEnable = false
	}
	if maxPower != "" {
		maxPowerEnable = true
		maxPowerInt, err = strconv.Atoi(maxPower)
		if err != nil {
			fmt.Printf("Error parsing max power: %v\n", err)
			return
		}
	} else {
		maxPowerEnable = false
	}
	hcuList, _, err := formatHcuList(groupId)
	if err != nil {
		fmt.Printf("Error getting hcu list: %v\n", err)
		return
	}

	for _, hcuIndex := range hcuList {
		conditions := dcgm.PolicyConditions{
			MaxPagesEnable:   maxPagesEnable,
			MaxPages:         maxPagesInt,
			MaxTempEnable:    maxTempEnable,
			MaxTemp:          maxTempFloat,
			MaxPowerEnable:   maxPowerEnable,
			MaxPower:         maxPowerInt,
			EccErrorsEnable:  eccErrorFlag,
			PcieErrorsEnable: pcieErrorFlag,
		}
		policyInfo := dcgm.Policy{
			DvInd:           hcuIndex,
			ActionIndex:     action,
			ValidationIndex: validation,
			Conditions:      conditions,
		}
		err := dcgm.SetPolicy(policyInfo, hcuIndex)
		if err != nil {
			fmt.Printf("Error setting policy: %v\n", err)
			return
		}
	}
	fmt.Println("Policy successfully set.")
}

func handleReg() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-stop
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		done <- true
	}()

	fmt.Println("Listening for violations.")
	hcuList, _, err := formatHcuList(groupId)
	if err != nil {
		fmt.Printf("Error getting hcu list: %v\n", err)
		return
	}

loop:
	for {
		select {
		case <-done:
			break loop
		default:
			hcuIndex, err := dcgm.JudgePolicyConditions(hcuList)
			if err != nil {
				now := time.Now()
				fmt.Printf("Timestamp: %s\n", now.Format("Mon Jan 2 15:04:05 2006"))
				fmt.Println(err)
				fmt.Printf("Action will be performed on HCU %d in 10 seconds.\n", hcuIndex)
				time.Sleep(10 * time.Second)
				err = dcgm.TakePolicyAction(hcuIndex)
				if err != nil {
					fmt.Printf("Error taking policy action: %v\n", err)
					break loop
				}
				time.Sleep(5 * time.Second)
				//break loop
			}
			time.Sleep(10 * time.Second)
		}
	}
}
