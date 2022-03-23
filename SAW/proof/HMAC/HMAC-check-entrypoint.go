/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	utility "aws-lc-verification/proof/common"
	"log"
	"os"
)

func main() {
	log.Printf("Started HMAC check.")
	// When 'HMAC_SELECTCHECK' is undefined, quickcheck is executed.
	env_var := os.Getenv("HMAC_SELECTCHECK")
	if len(env_var) == 0 {
		utility.RunSawScript("verify-HMAC-SHA384-quickcheck.saw")
		return
	}
	// When 'HMAC_SELECTCHECK' is defined, selectcheck is executed.
	start_indx := utility.ParseSelectCheckRange("HMAC_SELECTCHECK_START_INDX", 0)
	end_indx := utility.ParseSelectCheckRange("HMAC_SELECTCHECK_END_INDX", 127)
	if !(start_indx >= 0 && end_indx < 128) {
		log.Fatal("The HMAC select check range should be within [0, 127], but got [%s, %s]", start_indx, end_indx)
	}
	saw_template := "verify-HMAC-SHA384-selectcheck-template.txt"
	placeholder_map := map[string]int{
		"RANGE_START_INDX_PLACEHOLDER": start_indx,
		"RANGE_END_INDX_PLACEHOLDER":   end_indx,
	}
	utility.CreateAndRunSawScript(saw_template, placeholder_map, nil)

	log.Printf("Completed HMAC check.")
}
