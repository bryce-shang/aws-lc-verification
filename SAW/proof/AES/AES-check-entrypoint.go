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
	log.Printf("Started AES check.")
	// When 'AES_SELECTCHECK' is undefined, quickcheck is executed.
	env_var := os.Getenv("AES_SELECTCHECK")
	if len(env_var) == 0 {
		utility.RunSawScript("verify-AES-GCM.saw")
		return
	}

	// When 'AES_SELECTCHECK' is defined, selectcheck is executed.
	utility.RunSawScript("verify-HMAC-SHA384-selectcheck.saw")

	log.Printf("Completed AES check.")
}
