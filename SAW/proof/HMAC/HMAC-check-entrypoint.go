/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// A utility function to terminate this program when err exists.
func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// A function to create and run saw script.
func createAndRunSawScript(path_to_template string, placeholder_key string, value int, wg *sync.WaitGroup) {
	log.Printf("Start creating saw script for target value %s based on template %s.", value, path_to_template)
	// Create a new saw script.
	file_name := fmt.Sprint(value, ".saw")
	file, err := os.Create(file_name)
	checkError(err)
	// Read file content of verification template.
	content, err := ioutil.ReadFile(path_to_template)
	checkError(err)
	verification_template := string(content)
	// Replace some placeholders of the file content with target values.
	text := strings.Replace(verification_template, placeholder_key, strconv.Itoa(value), 1)
	defer file.Close()
	file.WriteString(text)
	defer os.Remove(file_name)
	// Run saw script.
	defer wg.Done()
	runSawScript(file_name)
}

// A function to run saw script.
func runSawScript(path_to_saw_file string) {
	log.Printf("Running saw script %s", path_to_saw_file)
	cmd := exec.Command("saw", path_to_saw_file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	checkError(err)
}

func main() {
	log.Printf("Started HMAC check.")
	// When 'HMAC_SELECTCHECK' is undefined, quickcheck is executed.
	env_var := os.Getenv("HMAC_SELECTCHECK")
	if len(env_var) == 0 {
		runSawScript("verify-HMAC-SHA384-quickcheck.saw")
		return
	}

	// When 'HMAC_SELECTCHECK' is defined, run 'HMAC-Init-ex' with diff parameters concurrently.
	var wg sync.WaitGroup
	for num := 0; num <= 129; num++ {
		wg.Add(1)
		saw_template := "verify-HMAC-Init-ex-selectcheck-template.txt"
		placeholder_name := "HMAC_TARGET_KEY_LEN_PLACEHOLDER"
		go createAndRunSawScript(saw_template, placeholder_name, num, &wg)
	}
	wg.Wait()

	// When 'HMAC_SELECTCHECK' is defined, run 'HMAC-Final' with diff parameters concurrently.
	for num := 0; num <= 127; num++ {
		wg.Add(1)
		saw_template := "verify-HMAC-Final-selectcheck-template.txt"
		placeholder_name := "HMAC_TARGET_NUM_PLACEHOLDER"
		go createAndRunSawScript(saw_template, placeholder_name, num, &wg)
	}

	wg.Wait()
	log.Printf("Completed HMAC check.")
}
