package jobs

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func RunSummarizer(text string) {
	pythonPath := "python3" 
	scriptPath := "summrizer.py"

	// Create the command
	cmd := exec.Command(pythonPath, scriptPath)
	cmd.Stdin = bytes.NewBufferString(text) 

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running script: %s\nOutput: %s", err, out.String())
	}

	// Print the output
	fmt.Println("Summary:", out.String())

}
