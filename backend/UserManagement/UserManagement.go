package UserManagement

import (
	"Proyecto1/backend/DiskControl"
	"Proyecto1/backend/DiskStruct"
	"Proyecto1/backend/FileManagement"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Login(user string, password string, id string) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Password:", password)
	fmt.Println("Id:", id)

	// Verify if the user is already logged in a partition
	mountedPartitions := DiskControl.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false // Nobody is logged in

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.ID == id && partition.LoggedIn { //Find the user in the mounted partitions
				fmt.Println("Ya existe un usuario logueado!")
				return
			}
			if partition.ID == id { // Find the partition with the given id
				filepath = partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No se encontró ninguna partición montada con el ID proporcionado")
		return
	}

	// Open bin file
	file, err := FileManagement.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR DiskStruct.MRB
	// Read the MBR from the binary file
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	DiskStruct.PrintMBR(TempMBR)
	fmt.Println("-------------")

	var index int = -1
	// Find the correct partition in the MBR
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("Partition found")
				if TempMBR.Partitions[i].Status[0] == '1' {
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		DiskStruct.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock DiskStruct.Superblock
	// Read the Superblock from the binary file
	if err := FileManagement.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Find users.txt and returns the index of the Inode
	indexInode := InitSearch("/users.txt", file, tempSuperblock)

	var crrInode DiskStruct.Inode
	// Read the Inode from the binary file
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo:", err)
		return
	}

	// Read the data from the file
	data := GetInodeFileData(crrInode, file, tempSuperblock)

	// Split the data by lines
	lines := strings.Split(data, "\n")

	// Iterate over the lines to find the user and password
	for _, line := range lines {
		words := strings.Split(line, ",")

		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], password)) {
				login = true
				break
			}
		}
	}

	fmt.Println("Inode", crrInode.I_block)

	// If the login was successful, mark the partition as logged in
	if login {
		fmt.Println("Usuario logueado con exito")
		DiskControl.MarkPartitionAsLoggedIn(id)
	}

	fmt.Println("======End LOGIN======")
}

// Returned value is the index of the Inode
func InitSearch(path string, file *os.File, tempSuperblock DiskStruct.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INICIAL ======")
	fmt.Println("path:", path)

	//Search and split the path (we need users.txt)
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 DiskStruct.Inode
	// Read object from bin file
	if err := FileManagement.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1
	}

	fmt.Println("======End BUSQUEDA INICIAL======")

	return SearchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

// stack (pila) to store logged in users
func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

// Search Inode by path
func SearchInodeByPath(StepsPath []string, Inode DiskStruct.Inode, file *os.File, tempSuperblock DiskStruct.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INODO POR PATH======")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {

				//==== DIRECT CASE ====
				var crrFolderBlock DiskStruct.Folderblock
				// Read object from bin file
				if err := FileManagement.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(DiskStruct.Folderblock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_content {
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if strings.Contains(string(folder.B_name[:]), SearchedName) {

						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder found======")
							return folder.B_inodo
						} else {
							fmt.Println("NextInode======")
							var NextInode DiskStruct.Inode
							// Read object from bin file
							if err := FileManagement.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
								return -1
							}
							return SearchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}

			} else {
				// ==== INDIRECT CASE ====
				// ----- TO DO	-----
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End BUSQUEDA INODO POR PATH======")
	return 0
}

// Logout function
func Logout() {
	fmt.Println("====== Start LOGOUT ======")

	// Get the mounted partitions
	mountedPartitions := DiskControl.GetMountedPartitions()
	var sessionActive bool
	var activePartitionID string

	// Verify if there is an active session
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn {
				sessionActive = true // There is an active session
				activePartitionID = partition.ID
				break
			}
		}
		if sessionActive {
			break
		}
	}

	// No logout if there is no active session
	if !sessionActive {
		fmt.Println("Error: No hay ninguna sesión activa.")
		fmt.Println("====== End LOGOUT ======")
		return
	}

	// Logout the active session
	DiskControl.MarkPartitionAsLoggedOut(activePartitionID)
	fmt.Println("Sesión cerrada con éxito en la partición:", activePartitionID)

	fmt.Println("====== End LOGOUT ======")
}

// Get the data from an Inode
func GetInodeFileData(Inode DiskStruct.Inode, file *os.File, tempSuperblock DiskStruct.Superblock) string {
	fmt.Println("======Start CONTENIDO DEL BLOQUE======")
	index := int32(0)

	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			//Inside of direct ones
			if index < 13 {
				var crrFileBlock DiskStruct.Fileblock
				// Read object from bin file
				if err := FileManagement.ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(DiskStruct.Fileblock{})))); err != nil {
					return ""
				}
				content += string(crrFileBlock.B_content[:])
			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End CONTENIDO DEL BLOQUE======")
	return content
}

//===== MKUSER =====

func AppendToFileBlock(inode *DiskStruct.Inode, newData string, file *os.File, superblock DiskStruct.Superblock) error {
	// Read the existing data from the file
	existingData := GetInodeFileData(*inode, file, superblock)

	// Join the existing data with the new data
	fullData := existingData + newData

	// Verify if the content exceeds the current block capacity
	if len(fullData) > len(inode.I_block)*binary.Size(DiskStruct.Fileblock{}) {
		// TO DO: Logic to create a new block if the file exceeds the current block capacity (64) !!!
		return fmt.Errorf("el tamaño del archivo excede la capacidad del bloque actual y no se ha implementado la creación de bloques adicionales")
	}

	// Write the new data to the block
	var updatedFileBlock DiskStruct.Fileblock
	copy(updatedFileBlock.B_content[:], fullData)
	if err := FileManagement.WriteObject(file, updatedFileBlock, int64(superblock.S_block_start+inode.I_block[0]*int32(binary.Size(DiskStruct.Fileblock{})))); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Update the inode size
	inode.I_size = int32(len(fullData))
	if err := FileManagement.WriteObject(file, *inode, int64(superblock.S_inode_start+inode.I_block[0]*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}
