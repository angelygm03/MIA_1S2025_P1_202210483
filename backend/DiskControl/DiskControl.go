package DiskControl

import (
	"Proyecto1/backend/DiskStruct"
	"Proyecto1/backend/FileManagement"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Second validation of the command
func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("======INICIO MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fmt.Println("Path:", path)

	// Validate fit bf/ff/wf
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe ser bf, wf or ff")
		return
	}

	// Validate size > 0
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayo a  0")
		return
	}

	// Validate k - m
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades validas son k o m")
		return
	}

	// Create file
	err := FileManagement.CreateFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Convert size to bytes
	if unit == "k" {
		size = size * 1024 // 1 KB = 1024
	} else {
		size = size * 1024 * 1024 // 1 MB = 1024 * 1024 bytes
	}

	// Open bin file
	file, err := FileManagement.OpenFile(path)
	if err != nil {
		return
	}

	//  === Write MBR ===
	// Create array of byte(0)
	for i := 0; i < size; i++ {
		err := FileManagement.WriteObject(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// Create MRB
	var newMRB DiskStruct.MRB       // Create a new MRB
	newMRB.MbrSize = int32(size)    // Set the size
	newMRB.Signature = rand.Int31() // Set the signature to a random number
	copy(newMRB.Fit[:], fit)        // Set the fit

	// Date format yyyy-mm-dd
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	copy(newMRB.CreationDate[:], formattedDate)

	// Write the MRB
	if err := FileManagement.WriteObject(file, newMRB, 0); err != nil {
		return
	}

	// === Read MBR ===
	var TempMBR DiskStruct.MRB
	// Leer el archivo
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	DiskStruct.PrintMBR(TempMBR)

	// Cerrar el archivo
	defer file.Close()

	fmt.Println("======FIN MKDISK======")

}

// Fuction to remove a disk
func Rmdisk(path string) {
	fmt.Println("======INICIO RMDISK======")
	fmt.Println("Path:", path)

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Error: El archivo no existe en la ruta especificada")
		return
	}

	// Confirm deletion
	fmt.Println("¿Está seguro de que desea eliminar el archivo? (yes/no):")
	var confirmation string
	fmt.Scanln(&confirmation)

	if strings.ToLower(confirmation) == "yes" {
		// Remove the file
		err := os.Remove(path)
		if err != nil {
			fmt.Println("Error al eliminar el archivo:", err)
			return
		}
		fmt.Println("Archivo eliminado exitosamente")
	} else {
		fmt.Println("Operación cancelada")
	}

	fmt.Println("======FIN RMDISK======")
}
