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
	for gcm_len := 0; gcm_len <= 32; gcm_len++ {
		wg.Add(1)
		saw_template := "verify-AES-GCM-selectcheck.txt"
		placeholder_name := "GCM_LEN_PLACEHOLDER"
		go utility.CreateAndRunSawScript(saw_template, placeholder_name, gcm_len, &wg)
		utility.Wait(&process_count, aes_process_limit, &wg)
	}

	wg.Wait()

	log.Printf("Completed AES check.")
}
