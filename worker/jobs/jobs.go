package jobs

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"gorm.io/gorm"
)

func RunSummarizer(text string, dbConn *gorm.DB) (string, error){
	log.Println(text)
	pythonPath := "python3" 
	scriptPath := "summrizer.py"
	cmd := exec.Command(pythonPath, scriptPath)
	cmd.Stdin = bytes.NewBufferString(text) 
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running script: %s\nOutput: %s", err, out.String())
		return "", fmt.Errorf("Error running script: %s\nOutput: %s", err, out.String())
	}
	return out.String(), nil
}
