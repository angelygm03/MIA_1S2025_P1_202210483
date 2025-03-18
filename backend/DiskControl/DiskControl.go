package DiskControl

import (
	"Proyecto1/backend/DiskStruct"
	"Proyecto1/backend/FileManagement"
	"encoding/binary"
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

	// If fit is empty
	if fit == "" {
		fit = "ff"
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

	// If unit is empty
	if unit == "" {
		unit = "m"
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

func Fdisk(size int, path string, name string, unit string, type_ string, fit string) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Path:", path)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)

	// Fit bf, ff, wf
	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: Fit must be 'bf', 'ff', or 'wf'")
		return
	}
	// If fit is empty, set it to "w"
	if fit == "" {
		fit = "wf"
	}

	// Size must be greater than 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// Validate unit b, k or m
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be 'b', 'k', or 'm'")
		return
	}

	// If unit is empty, set it to "k"
	if unit == "" {
		unit = "k"
	}

	// Validate type p, e or l
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Type must be 'p', 'e', or 'l'")
		return
	}
	// If type is empty, set it to "p"
	if type_ == "" {
		type_ = "p"
	}

	// Size to bytes
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// Open file in correct path
	file, err := FileManagement.OpenFile(path)
	if err != nil {
		fmt.Println("Error: Could not open file at path:", path)
		return
	}

	var TempMBR DiskStruct.MRB
	// Read the object from the binary file
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file")
		return
	}

	// Print the object
	DiskStruct.PrintMBR(TempMBR)

	fmt.Println("-------------")

	// Partitions validation
	var primaryCount, extendedCount, totalPartitions int
	var usedSpace int32 = 0

	// Count the number of partitions (4 are allowed)
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			totalPartitions++
			usedSpace += TempMBR.Partitions[i].Size

			// Count the number of primary and extended partitions
			if TempMBR.Partitions[i].Type[0] == 'p' {
				primaryCount++
			} else if TempMBR.Partitions[i].Type[0] == 'e' {
				extendedCount++
			}
		}
	}

	// Validate that there are not more than 4 partitions
	if totalPartitions >= 4 {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		return
	}

	// Validate that exits an extended partition
	if type_ == "e" && extendedCount > 0 {
		fmt.Println("Error: Solo se permite una partición extendida por disco.")
		return
	}

	// If theres no extended partition, a logical partition can't be created
	if type_ == "l" && extendedCount == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		return
	}

	// Partition size can't be greater than the disk size
	if usedSpace+int32(size) > TempMBR.MbrSize {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return
	}

	// Starting position of the new partition
	var gap int32 = int32(binary.Size(TempMBR))
	if totalPartitions > 0 {
		gap = TempMBR.Partitions[totalPartitions-1].Start + TempMBR.Partitions[totalPartitions-1].Size
	}

	// Encontrar una posición vacía para la nueva partición
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size == 0 {
			if type_ == "p" || type_ == "e" {
				// Creating a primary or extended partition
				TempMBR.Partitions[i].Size = int32(size)
				TempMBR.Partitions[i].Start = gap
				copy(TempMBR.Partitions[i].Name[:], name)
				copy(TempMBR.Partitions[i].Fit[:], fit)
				copy(TempMBR.Partitions[i].Status[:], "0")
				copy(TempMBR.Partitions[i].Type[:], type_)
				TempMBR.Partitions[i].Correlative = int32(totalPartitions + 1)

				// if the partition is extended, initialize the EBR
				if type_ == "e" {
					ebrStart := gap // First EBR starts at the beginning of the extended partition
					ebr := DiskStruct.EBR{
						PartFit:   fit[0],
						PartStart: ebrStart,
						PartSize:  0,
						PartNext:  -1,
					}
					copy(ebr.PartName[:], "")
					FileManagement.WriteObject(file, ebr, int64(ebrStart)) // Write the EBR
				}
				break
			}
		}
	}

	// If the partition is logical
	if type_ == "l" {
		for i := 0; i < 4; i++ {
			// Find the extended partition
			if TempMBR.Partitions[i].Type[0] == 'e' {
				ebrPos := TempMBR.Partitions[i].Start
				var ebr DiskStruct.EBR
				//Find the EBR
				for {
					FileManagement.ReadObject(file, &ebr, int64(ebrPos))
					if ebr.PartNext == -1 {
						break
					}
					ebrPos = ebr.PartNext
				}

				// Starting position of the logical partition is calculated
				newEBRPos := ebr.PartStart + ebr.PartSize                    // The new EBR starts right after the last logical partition
				logicalPartitionStart := newEBRPos + int32(binary.Size(ebr)) // The logical partition starts right after the EBR

				// Adjust tbe next EBR
				ebr.PartNext = newEBRPos
				FileManagement.WriteObject(file, ebr, int64(ebrPos))

				// Create and write new EBR
				newEBR := DiskStruct.EBR{
					PartFit:   fit[0],
					PartStart: logicalPartitionStart,
					PartSize:  int32(size),
					PartNext:  -1,
				}
				copy(newEBR.PartName[:], name)
				FileManagement.WriteObject(file, newEBR, int64(newEBRPos))

				fmt.Println("Nuevo EBR creado:")
				DiskStruct.PrintEBR(newEBR)
				break
			}
		}
	}

	// Overwrite the MBR
	if err := FileManagement.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: Could not write MBR to file")
		return
	}

	var TempMBR2 DiskStruct.MRB
	// Verify the MBR was written correctly
	if err := FileManagement.ReadObject(file, &TempMBR2, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file after writing")
		return
	}

	DiskStruct.PrintMBR(TempMBR2)

	for i := 0; i < 4; i++ {
		if TempMBR2.Partitions[i].Type[0] == 'e' {
			fmt.Println("Leyendo EBRs dentro de la partición extendida...")
			ebrPos := TempMBR2.Partitions[i].Start
			var ebr DiskStruct.EBR
			for {
				err := FileManagement.ReadObject(file, &ebr, int64(ebrPos))
				if err != nil {
					fmt.Println("Error al leer un EBR:", err)
					break
				}
				fmt.Println("EBR encontrado en la posición:", ebrPos)
				DiskStruct.PrintEBR(ebr)
				if ebr.PartNext == -1 {
					break
				}
				ebrPos = ebr.PartNext
			}
		}
	}

	// Close file to avoid memory leaks
	defer file.Close()

	fmt.Println("======FIN FDISK======")
}
