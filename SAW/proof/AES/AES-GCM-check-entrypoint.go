/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	utility "aws-lc-verification/proof/common"
	"log"
	"os"
	"strconv"
	"sync"
)

const aes_process_limit int = 20

func main() {
	log.Printf("Started AES check.")
	// When 'AES_SELECTCHECK' is undefined, quickcheck is executed.
	env_var := os.Getenv("AES_SELECTCHECK")
	if len(env_var) == 0 {
		utility.RunSawScript("verify-AES-GCM-quickcheck.saw")
		return
	}

	// Generate saw scripts based on above verification template and target gcm len.
	var wg sync.WaitGroup
	process_count := 0
	placeholder_map := make(map[string]string)
	saw_template := "verify-AES-GCM-selectcheck.txt"
	for gcm_len := 10; gcm_len <= 10; gcm_len++ {
		for update_len := 1; update_len <= 320; update_len++ {
			wg.Add(1)
			placeholder_map["GCM_LEN_PLACEHOLDER"] = strconv.Itoa(gcm_len)
			placeholder_map["GCM_LEN_PLACEHOLDER"] = strconv.Itoa(gcm_len)
			go utility.CreateAndRunSawScript(saw_template, placeholder_map, &wg)
			utility.Wait(&process_count, aes_process_limit, &wg)
		}
	}

	wg.Wait()

	log.Printf("Completed AES check.")
}
