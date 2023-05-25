package kc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func writeLabelToFile(label string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	filePath := filepath.Join(homeDir, LABEL_FILE_NAME)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()
	}

	existingDomains, err := readLabelsFromFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if contains(existingDomains, label) {
		fmt.Println("Label", label, "already exists in the file.")
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	file.WriteString(label + "\n")

	fmt.Println("Label", label, "written to file successfully at:", filePath)
}

func readLabelsFromFile(filePath string) ([]string, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	labels := strings.Split(string(fileContent), "\n")
	return labels, nil
}

func deleteLableInFile(label string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	filePath := filepath.Join(homeDir, LABEL_FILE_NAME)

	// Read the existing file contents
	labels, err := readLabelsFromFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Check if the domain exists in the file
	index := findIndex(labels, label)
	if index == -1 {
		fmt.Println("Label", label, "does not exist in the file.")
		return
	}

	// Remove the domain from the slice
	labels = remove(labels, index)

	// Write the updated domains to the file
	err = writeLabelsToFile(filePath, labels)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Label", label, "has been deleted from the file.")
}

// Write domains to the file
func writeLabelsToFile(filePath string, labels []string) error {
	fileContent := strings.Join(labels, "\n")
	return ioutil.WriteFile(filePath, []byte(fileContent), 0644)
}
