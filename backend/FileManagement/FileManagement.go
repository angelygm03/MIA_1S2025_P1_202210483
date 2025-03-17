package FileManagement

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// Function to create a file
func CreateFile(name string) error {
	// Check if the file exists
	dir := filepath.Dir(name) // Get the directory
	// if the directory does not exist, create it
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error creating file dir==", err)
		return err
	}

	// Create the file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Error creating file==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Function to open a file
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644) // Permission 644: Read and write
	if err != nil {
		fmt.Println("Error opening file==", err)
		return nil, err
	}
	return file, nil
}

// Function to write an object to a file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0) // Move the pointer to the position
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Error writing object==", err)
		return err
	}
	return nil
}

// Function to read an object from a file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data) // Read the object
	if err != nil {
		fmt.Println("Error reading object==", err)
		return err
	}
	return nil
}
