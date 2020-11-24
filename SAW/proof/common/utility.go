/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// A utility function to terminate this program when err exists.
func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// A function to create a saw script from saw template, replace placeholder with target value, and then execute the script.
func CreateAndRunSawScript(path_to_template string, placeholder_map map[string]string, wg *sync.WaitGroup) {
	log.Printf("Start creating saw script based on template %s.", path_to_template)
	// Create a new saw script.
	file_name := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	file, err := os.Create(file_name)
	CheckError(err)
	// Read file content of verification template.
	content, err := ioutil.ReadFile(path_to_template)
	CheckError(err)
	verification_script := string(content)
	// Replace some placeholders of the file content with target values.
	for placeholder_key, placeholder_value := range placeholder_map {
		verification_script = strings.Replace(verification_script, placeholder_key, placeholder_value, 1)
	}
	defer file.Close()
	file.WriteString(verification_script)
	defer os.Remove(file_name)
	// Run saw script.
	defer wg.Done()
	RunSawScript(file_name)
}

// A function to run saw script.
func RunSawScript(path_to_saw_file string) {
	log.Printf("Running saw script %s", path_to_saw_file)
	cmd := exec.Command("saw", path_to_saw_file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	CheckError(err)
}

// A function to limit number of concurrent processes.
func Wait(process_count *int, limit int, wg *sync.WaitGroup) {
	if *process_count >= limit {
		log.Printf("Count [%d] reached process limit [%d].", *process_count, limit)
		wg.Wait()
		*process_count = 0
	} else {
		*process_count = (*process_count) + 1
	}
}
