package UserManagement

import (
	"Proyecto1/backend/DiskControl"
	"Proyecto1/backend/DiskStruct"
	"Proyecto1/backend/FileManagement"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
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

			} else if index == 13 {
				// ==== INDIRECT CASE ====
				fmt.Println("Indirect case: Simple Indirect Block")

				// Read the Pointerblock
				var pointerBlock DiskStruct.Pointerblock
				if err := FileManagement.ReadObject(file, &pointerBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(DiskStruct.Pointerblock{})))); err != nil {
					fmt.Println("Error reading Pointerblock:", err)
					return -1
				}

				// Iterate over the pointers in the Pointerblock
				for _, pointer := range pointerBlock.B_pointers {
					if pointer != -1 {
						var crrFolderBlock DiskStruct.Folderblock
						// Read the Folderblock pointed by the current pointer
						if err := FileManagement.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+pointer*int32(binary.Size(DiskStruct.Folderblock{})))); err != nil {
							fmt.Println("Error reading Folderblock:", err)
							return -1
						}

						// Iterate over the contents of the Folderblock
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
									// Read the next Inode
									if err := FileManagement.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
										fmt.Println("Error reading NextInode:", err)
										return -1
									}
									return SearchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
								}
							}
						}
					}
				}
			} else if index == 14 {
				// ==== DOUBLE INDIRECT CASE ====
				fmt.Println("Indirect case: Double Indirect Block")

				// Read the first-level Pointerblock
				var firstLevelPointerBlock DiskStruct.Pointerblock
				if err := FileManagement.ReadObject(file, &firstLevelPointerBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(DiskStruct.Pointerblock{})))); err != nil {
					fmt.Println("Error reading first-level Pointerblock:", err)
					return -1
				}

				// Iterate over the first-level pointers
				for _, firstPointer := range firstLevelPointerBlock.B_pointers {
					if firstPointer != -1 {
						// Read the second-level Pointerblock
						var secondLevelPointerBlock DiskStruct.Pointerblock
						if err := FileManagement.ReadObject(file, &secondLevelPointerBlock, int64(tempSuperblock.S_block_start+firstPointer*int32(binary.Size(DiskStruct.Pointerblock{})))); err != nil {
							fmt.Println("Error reading second-level Pointerblock:", err)
							return -1
						}

						// Iterate over the second-level pointers
						for _, secondPointer := range secondLevelPointerBlock.B_pointers {
							if secondPointer != -1 {
								var crrFolderBlock DiskStruct.Folderblock
								// Read the Folderblock pointed by the second-level pointer
								if err := FileManagement.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+secondPointer*int32(binary.Size(DiskStruct.Folderblock{})))); err != nil {
									fmt.Println("Error reading Folderblock:", err)
									return -1
								}

								// Iterate over the contents of the Folderblock
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
											// Read the next Inode
											if err := FileManagement.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
												fmt.Println("Error reading NextInode:", err)
												return -1
											}
											return SearchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
										}
									}
								}
							}
						}
					}
				}
			}
		} else if index == 15 {
			// ==== TRIPLE INDIRECT CASE ====
			fmt.Println("Indirect case: Triple Indirect Block")

			// Read the first-level Pointerblock
			var firstLevelPointerBlock DiskStruct.Pointerblock
			if err := FileManagement.ReadObject(file, &firstLevelPointerBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(DiskStruct.Pointerblock{})))); err != nil {
				fmt.Println("Error reading first-level Pointerblock:", err)
				return -1
			}

			// Iterate over the first-level pointers
			for _, firstPointer := range firstLevelPointerBlock.B_pointers {
				if firstPointer != -1 {
					// Read the second-level Pointerblock
					var secondLevelPointerBlock DiskStruct.Pointerblock
					if err := FileManagement.ReadObject(file, &secondLevelPointerBlock, int64(tempSuperblock.S_block_start+firstPointer*int32(binary.Size(DiskStruct.Pointerblock{})))); err != nil {
						fmt.Println("Error reading second-level Pointerblock:", err)
						return -1
					}

					// Iterate over the second-level pointers
					for _, secondPointer := range secondLevelPointerBlock.B_pointers {
						if secondPointer != -1 {
							// Read the third-level Pointerblock
							var thirdLevelPointerBlock DiskStruct.Pointerblock
							if err := FileManagement.ReadObject(file, &thirdLevelPointerBlock, int64(tempSuperblock.S_block_start+secondPointer*int32(binary.Size(DiskStruct.Pointerblock{})))); err != nil {
								fmt.Println("Error reading third-level Pointerblock:", err)
								return -1
							}

							// Iterate over the third-level pointers
							for _, thirdPointer := range thirdLevelPointerBlock.B_pointers {
								if thirdPointer != -1 {
									var crrFolderBlock DiskStruct.Folderblock
									// Read the Folderblock pointed by the third-level pointer
									if err := FileManagement.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+thirdPointer*int32(binary.Size(DiskStruct.Folderblock{})))); err != nil {
										fmt.Println("Error reading Folderblock:", err)
										return -1
									}

									// Iterate over the contents of the Folderblock
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
												// Read the next Inode
												if err := FileManagement.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
													fmt.Println("Error reading NextInode:", err)
													return -1
												}
												return SearchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
											}
										}
									}
								}
							}
						}
					}
				}
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
	// read existing data
	existingData := GetInodeFileData(*inode, file, superblock)
	fmt.Println("Datos existentes en el bloque:", existingData)

	// Concatenate the new data
	fullData := existingData + newData
	fmt.Println("Datos completos después de agregar:", fullData)

	// Verify that the file size does not exceed the total capacity
	blockSize := binary.Size(DiskStruct.Fileblock{})
	totalCapacity := len(inode.I_block) * blockSize
	if len(fullData) > totalCapacity {
		fmt.Printf("Error: El tamaño del archivo (%d bytes) excede la capacidad total asignada (%d bytes).\n", len(fullData), totalCapacity)
		return fmt.Errorf("el tamaño del archivo excede la capacidad total asignada y no se ha implementado la creación de bloques adicionales")
	}

	// Split the full data into blocks
	remainingData := fullData
	blockIndex := 0

	for len(remainingData) > 0 {
		// If no block is assigned, find a free block
		if blockIndex >= len(inode.I_block) || inode.I_block[blockIndex] == -1 {
			newBlockIndex := FindFreeBlock(superblock, file)
			if newBlockIndex == -1 {
				return fmt.Errorf("no hay bloques libres disponibles")
			}
			inode.I_block[blockIndex] = int32(newBlockIndex)
			fmt.Printf("Asignando nuevo bloque: %d\n", newBlockIndex)
		}

		// Create a new file block with the data
		var updatedFileBlock DiskStruct.Fileblock
		copy(updatedFileBlock.B_content[:], remainingData[:min(len(remainingData), blockSize)])

		// Write the updated block to the file
		position := int64(superblock.S_block_start + inode.I_block[blockIndex]*int32(blockSize))
		fmt.Printf("Escribiendo bloque en la posición: %d\n", position)
		if err := FileManagement.WriteObject(file, updatedFileBlock, position); err != nil {
			return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
		}

		// Update the rest of the data
		remainingData = remainingData[min(len(remainingData), blockSize):]
		blockIndex++
	}

	// Update inode size
	inode.I_size = int32(len(fullData))
	inodePosition := int64(superblock.S_inode_start + inode.I_block[0]*int32(binary.Size(DiskStruct.Inode{})))
	fmt.Printf("Actualizando inodo en la posición: %d\n", inodePosition)
	fmt.Printf("Nuevo tamaño del inodo (I_size): %d\n", inode.I_size)
	if err := FileManagement.WriteObject(file, *inode, inodePosition); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	fmt.Println("Bloque e inodo actualizados correctamente.")
	return nil
}

// Aux func to find a free block of the superblock and
func FindFreeBlock(superblock DiskStruct.Superblock, file *os.File) int {
	// Read the bitmap of blocks
	bitmap := make([]byte, superblock.S_blocks_count)
	if _, err := file.ReadAt(bitmap, int64(superblock.S_bm_block_start)); err != nil {
		fmt.Println("Error al leer el bitmap de bloques:", err)
		return -1
	}

	// Find the first free block
	for i, b := range bitmap {
		if b == 0 {
			// Mark the block as used
			bitmap[i] = 1
			if _, err := file.WriteAt(bitmap, int64(superblock.S_bm_block_start)); err != nil {
				fmt.Println("Error al actualizar el bitmap de bloques:", err)
				return -1
			}
			return i
		}
	}

	return -1 //If there is no free block
}

// Aux fun copy the new data to the block
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Mkusr(user string, pass string, grp string) {
	fmt.Printf("Parámetros recibidos: user=%s, pass=%s, grp=%s\n", user, pass, grp)

	// Validate that the user is root
	if !IsRootUser() {
		fmt.Println("Error: Solo el usuario root puede ejecutar este comando.")
		fmt.Println("====== End MKUSR ======")
		return
	}

	// Validate the length of the parameters
	if len(user) > 10 || len(pass) > 10 || len(grp) > 10 {
		fmt.Println("Error: Los parámetros 'user', 'pass' y 'grp' no pueden exceder los 10 caracteres.")
		fmt.Println("====== End MKUSR ======")
		return
	}

	// Get mounted partitions and find the active partition
	mountedPartitions := DiskControl.GetMountedPartitions()
	var filepath string
	var partitionFound bool

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn { // Find the active partition
				filepath = partition.Path
				partitionFound = true
				fmt.Printf("Partición activa encontrada: %s\n", filepath)
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No hay ninguna partición activa.")
		fmt.Println("====== End MKUSR ======")
		return
	}

	// Open bin file
	file, err := FileManagement.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		fmt.Println("====== End MKUSR ======")
		return
	}
	defer file.Close()

	// Read the MBR
	var TempMBR DiskStruct.MRB
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Read the Superblock
	var tempSuperblock DiskStruct.Superblock
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Status[0] == '1' { // Active partition
			if err := FileManagement.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[i].Start)); err != nil {
				fmt.Println("Error: No se pudo leer el Superblock:", err)
				return
			}
			break
		}
	}

	// Find the users.txt file
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		fmt.Println("====== End MKUSR ======")
		return
	}

	var crrInode DiskStruct.Inode
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo del archivo users.txt:", err)
		return
	}

	// Read the content of the users.txt file
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	fmt.Println("Contenido actual del archivo users.txt:")
	fmt.Println(data)

	// Verify if the group exists
	lines := strings.Split(data, "\n")
	groupExists := false
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 3 && words[1] == "G" && words[2] == grp {
			groupExists = true
			fmt.Printf("Grupo encontrado: %s\n", grp)
			break
		}
	}

	if !groupExists {
		fmt.Println("Error: El grupo especificado no existe.")
		fmt.Println("====== End MKUSR ======")
		return
	}

	// If the user already exists, return an error
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 5 && words[1] == "U" && words[3] == user {
			fmt.Println("Error: El usuario especificado ya existe.")
			fmt.Println("====== End MKUSR ======")
			return
		}
	}

	// Global counter for the user IDs
	if nextUserID == 0 {
		// Init the counter
		if err := InitializeUserIDCounter(file, tempSuperblock); err != nil {
			fmt.Println(err)
			fmt.Println("====== End MKUSR ======")
			return
		}
	}

	newUserID := nextUserID
	nextUserID++ // Increase the counter for the next user
	newUser := fmt.Sprintf("%d,U,%s,%s,%s\n", newUserID, user, grp, pass)
	fmt.Printf("Nuevo usuario a agregar: %s\n", newUser)

	// Add the new user to the users.txt file
	if err := AppendToFileBlock(&crrInode, newUser, file, tempSuperblock); err != nil {
		fmt.Println("Error: No se pudo agregar el nuevo usuario al archivo users.txt:", err)
		fmt.Println("====== End MKUSR ======")
		return
	}

	fmt.Println("Usuario creado exitosamente.")
	fmt.Println("====== End MKUSR ======")
}

// Aux fun to verify if the user is root or not
func IsRootUser() bool {
	// Get mounted partitions
	mountedPartitions := DiskControl.GetMountedPartitions()

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn {
				// If there is an active session, verify if the user is root
				file, err := FileManagement.OpenFile(partition.Path)
				if err != nil {
					fmt.Println("Error: No se pudo abrir el archivo:", err)
					return false
				}
				defer file.Close()

				var TempMBR DiskStruct.MRB
				if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
					fmt.Println("Error: No se pudo leer el MBR:", err)
					return false
				}

				var tempSuperblock DiskStruct.Superblock
				if err := FileManagement.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[0].Start)); err != nil {
					fmt.Println("Error: No se pudo leer el Superblock:", err)
					return false
				}

				// Find the users.txt file
				indexInode := InitSearch("/users.txt", file, tempSuperblock)
				if indexInode == -1 {
					fmt.Println("Error: No se encontró el archivo users.txt.")
					return false
				}

				var crrInode DiskStruct.Inode
				if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
					fmt.Println("Error: No se pudo leer el Inodo del archivo users.txt:", err)
					return false
				}

				// Read the content of the users.txt file
				data := GetInodeFileData(crrInode, file, tempSuperblock)

				// Verify if the logged user is root
				lines := strings.Split(data, "\n")
				for _, line := range lines {
					words := strings.Split(line, ",")
					if len(words) == 5 && words[1] == "U" && words[3] == "root" {
						return true
					}
				}
			}
		}
	}

	// If not active session, return false
	return false
}

func PrintUsersFile() {
	fmt.Println("====== Start Print Users File ======")

	// Get mounted partitions
	mountedPartitions := DiskControl.GetMountedPartitions()
	var filepath string
	var partitionFound bool

	// Find the active partition
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn { // Active partition
				filepath = partition.Path
				partitionFound = true
				fmt.Printf("Partición activa encontrada: %s\n", filepath)
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No hay ninguna partición activa.")
		fmt.Println("====== End Print Users File ======")
		return
	}

	// Open bin file
	file, err := FileManagement.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		fmt.Println("====== End Print Users File ======")
		return
	}
	defer file.Close()

	// Read the Superblock
	var TempMBR DiskStruct.MRB
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	var tempSuperblock DiskStruct.Superblock
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Status[0] == '1' { // Active partition
			if err := FileManagement.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[i].Start)); err != nil {
				fmt.Println("Error: No se pudo leer el Superblock:", err)
				return
			}
			break
		}
	}

	// Find the users.txt file
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		fmt.Println("====== End Print Users File ======")
		return
	}

	var crrInode DiskStruct.Inode
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo del archivo users.txt:", err)
		return
	}

	// Read the content of the users.txt file
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	fmt.Println("Contenido del archivo users.txt:")
	fmt.Println(data)

	fmt.Println("====== End Print Users File ======")
}

// ==== GLOBAL COUNTERS ====
var nextUserID int = 0
var nextGroupID int = 0

// Func to initialize the user ID counter
func InitializeUserIDCounter(file *os.File, tempSuperblock DiskStruct.Superblock) error {
	// Find the users.txt file
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Printf("error: No se encontró el archivo users.txt.")

	}

	// Read the Inode of the users.txt file
	var crrInode DiskStruct.Inode
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Printf("error al leer el Inodo del archivo users.txt: %v", err)
	}

	// Read the data from the file
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	lines := strings.Split(data, "\n")

	// Calculating the max ID
	maxID := 0
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) > 0 {
			// Get the ID from the first column
			if id, err := strconv.Atoi(strings.TrimSpace(words[0])); err == nil {
				if id > maxID {
					maxID = id
				}
			}
		}
	}

	// Update global counter
	nextUserID = maxID + 1
	fmt.Printf("Contador de IDs inicializado en: %d\n", nextUserID)
	return nil
}

// Func to init the grup counter
func InitializeGroupIDCounter(file *os.File, tempSuperblock DiskStruct.Superblock) error {
	// Find the users.txt file
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Printf("Error: No se encontró el archivo users.txt.\n")
		return fmt.Errorf("archivo users.txt no encontrado")
	}

	// Read the Inode of the users.txt file
	var crrInode DiskStruct.Inode
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Printf("Error al leer el Inodo del archivo users.txt: %v\n", err)
		return err
	}

	// Read the data from the file
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	lines := strings.Split(data, "\n")

	// Calculating the max group ID
	maxGroupID := 0
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) > 0 && len(words) >= 3 && strings.TrimSpace(words[1]) == "G" {
			// Get the ID from the first column
			if id, err := strconv.Atoi(strings.TrimSpace(words[0])); err == nil {
				if id > maxGroupID {
					maxGroupID = id
				}
			}
		}
	}

	// Update global counter
	nextGroupID = maxGroupID + 1
	fmt.Printf("Contador de IDs de grupos inicializado en: %d\n", nextGroupID)
	return nil
}

func Mkgrp(name string) {
	fmt.Printf("Parámetro recibido: name=%s\n", name)

	// User must be root
	if !IsRootUser() {
		fmt.Println("Error: Solo el usuario root puede ejecutar este comando.")
		fmt.Println("====== End MKGRP ======")
		return
	}

	// Get mounted partitions
	mountedPartitions := DiskControl.GetMountedPartitions()
	var filepath string
	var partitionFound bool

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn { // Find the active partition
				filepath = partition.Path
				partitionFound = true
				fmt.Printf("Partición activa encontrada: %s\n", filepath)
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No hay ninguna partición activa.")
		fmt.Println("====== End MKGRP ======")
		return
	}

	// Open bin file
	file, err := FileManagement.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		fmt.Println("====== End MKGRP ======")
		return
	}
	defer file.Close()

	// Read the MBR
	var TempMBR DiskStruct.MRB
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Read the Superblock
	var tempSuperblock DiskStruct.Superblock
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Status[0] == '1' { // Active partition
			if err := FileManagement.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[i].Start)); err != nil {
				fmt.Println("Error: No se pudo leer el Superblock:", err)
				return
			}
			break
		}
	}

	// Find the users.txt file
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		fmt.Println("====== End MKGRP ======")
		return
	}

	var crrInode DiskStruct.Inode
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo del archivo users.txt:", err)
		return
	}

	// Read the content of the users.txt file
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	fmt.Println("Contenido actual del archivo users.txt:")
	fmt.Println(data)

	// Verify if the group already exists
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 3 && words[1] == "G" && words[2] == name {
			fmt.Println("Error: El grupo especificado ya existe.")
			fmt.Println("====== End MKGRP ======")
			return
		}
	}

	// Init the global counter for the group IDs
	if nextGroupID == 0 {
		if err := InitializeGroupIDCounter(file, tempSuperblock); err != nil {
			fmt.Println(err)
			fmt.Println("====== End MKGRP ======")
			return
		}
	}

	// Create the new group
	newGroupID := nextGroupID
	nextGroupID++ //Increase the counter for the next group
	newGroup := fmt.Sprintf("%d,G,%s\n", newGroupID, name)
	fmt.Printf("Nuevo grupo a agregar: %s\n", newGroup)

	// Add the new group to the users.txt file
	if err := AppendToFileBlock(&crrInode, newGroup, file, tempSuperblock); err != nil {
		fmt.Println("Error: No se pudo agregar el nuevo grupo al archivo users.txt:", err)
		fmt.Println("====== End MKGRP ======")
		return
	}

	fmt.Println("Grupo creado exitosamente.")
	fmt.Println("====== End MKGRP ======")
}

func Rmusr(user string) {
	fmt.Printf("Parámetro recibido: user='%s'\n", user)

	// Validate that the user is root
	if !IsRootUser() {
		fmt.Println("Error: Solo el usuario root puede ejecutar este comando.")
		fmt.Println("====== End RMUSR ======")
		return
	}

	// Get mounted partitions and find the active partition
	mountedPartitions := DiskControl.GetMountedPartitions()
	var filepath string
	var partitionFound bool

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn { //active sesion
				filepath = partition.Path
				partitionFound = true
				fmt.Printf("Partición activa encontrada: %s\n", filepath)
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No hay ninguna partición activa.")
		fmt.Println("====== End RMUSR ======")
		return
	}

	// Open bin file
	file, err := FileManagement.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		fmt.Println("====== End RMUSR ======")
		return
	}
	defer file.Close()

	// Read the MBR
	var TempMBR DiskStruct.MRB
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Read the superblock
	var tempSuperblock DiskStruct.Superblock
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Status[0] == '1' { // active session
			if err := FileManagement.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[i].Start)); err != nil {
				fmt.Println("Error: No se pudo leer el Superblock:", err)
				return
			}
			break
		}
	}

	// Find the users.txt file
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		fmt.Println("====== End RMUSR ======")
		return
	}

	var crrInode DiskStruct.Inode
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo del archivo users.txt:", err)
		return
	}

	// Read the content of the users.txt file
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	fmt.Println("Contenido actual del archivo users.txt:")
	fmt.Println(data)

	// Find the user to remove
	lines := strings.Split(data, "\n")
	var updatedLines []string
	userFound := false

	// Clean the user parameter
	cleanedUser := strings.TrimSpace(user)
	cleanedUser = strings.ReplaceAll(cleanedUser, "\u200B", "") // Remove invisible characters

	for _, line := range lines {
		// Eliminar espacios en blanco adicionales
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "\u200B", "") // Remove invisible characters
		if line == "" {
			continue // Ignorar líneas vacías
		}

		words := strings.Split(line, ",")
		fmt.Printf("Campos de la línea: %v\n", words)

		if len(words) == 5 {
			// Clean the user field
			for i := range words {
				words[i] = strings.TrimSpace(words[i])
				words[i] = strings.ReplaceAll(words[i], "\u200B", "")
			}

			// Compare the user
			fmt.Printf("Comparando usuario: '%s' con '%s'\n", words[2], cleanedUser)
			if words[1] == "U" && words[2] == cleanedUser {
				// Change the status of the user to 0
				words[0] = "0"
				userFound = true
			}
			updatedLines = append(updatedLines, strings.Join(words, ","))
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	if !userFound {
		fmt.Printf("Error: El usuario '%s' no existe.\n", cleanedUser)
		fmt.Println("====== End RMUSR ======")
		return
	}

	// Update the content of the users.txt file
	newData := strings.Join(updatedLines, "\n")

	// Overwrite the file block with the updated data
	if err := OverwriteFileBlock(&crrInode, newData, file, tempSuperblock, indexInode); err != nil {
		fmt.Println("Error: No se pudo actualizar el archivo users.txt:", err)
		fmt.Println("====== End RMUSR ======")
		return
	}

	// Updated content of the users.txt file
	fmt.Println("Contenido actualizado del archivo users.txt:")
	fmt.Println(newData)

	fmt.Println("Usuario eliminado exitosamente.")
	fmt.Println("====== End RMUSR ======")
}

func OverwriteFileBlock(inode *DiskStruct.Inode, newData string, file *os.File, superblock DiskStruct.Superblock, indexInode int32) error {
	// Split the new data into blocks
	blockSize := binary.Size(DiskStruct.Fileblock{})
	remainingData := newData
	blockIndex := 0

	for len(remainingData) > 0 {
		// If no block is assigned, find a free block
		if blockIndex >= len(inode.I_block) || inode.I_block[blockIndex] == -1 {
			newBlockIndex := FindFreeBlock(superblock, file)
			if newBlockIndex == -1 {
				return fmt.Errorf("no hay bloques libres disponibles")
			}
			inode.I_block[blockIndex] = int32(newBlockIndex)
		}

		// New file block with the data
		var updatedFileBlock DiskStruct.Fileblock
		copy(updatedFileBlock.B_content[:], remainingData[:min(len(remainingData), blockSize)])

		// Write the updated block to the file
		position := int64(superblock.S_block_start + inode.I_block[blockIndex]*int32(blockSize))
		if err := FileManagement.WriteObject(file, updatedFileBlock, position); err != nil {
			return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
		}

		// Update the rest of the data
		remainingData = remainingData[min(len(remainingData), blockSize):]
		blockIndex++
	}

	//Update inode size
	inode.I_size = int32(len(newData))
	inodePosition := int64(superblock.S_inode_start + indexInode*int32(binary.Size(DiskStruct.Inode{})))
	if err := FileManagement.WriteObject(file, *inode, inodePosition); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}

func Rmgrp(name string) {
	fmt.Printf("Parámetro recibido: name='%s'\n", name)

	// VUser must be root
	if !IsRootUser() {
		fmt.Println("Error: Solo el usuario root puede ejecutar este comando.")
		fmt.Println("====== End RMGRP ======")
		return
	}

	// Get mounted partitions
	mountedPartitions := DiskControl.GetMountedPartitions()
	var filepath string
	var partitionFound bool

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn { // active session
				filepath = partition.Path
				partitionFound = true
				fmt.Printf("Partición activa encontrada: %s\n", filepath)
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No hay ninguna partición activa.")
		fmt.Println("====== End RMGRP ======")
		return
	}

	// Open bin file
	file, err := FileManagement.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		fmt.Println("====== End RMGRP ======")
		return
	}
	defer file.Close()

	// Read the MBR
	var TempMBR DiskStruct.MRB
	if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Read the Superblock
	var tempSuperblock DiskStruct.Superblock
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Status[0] == '1' { // active partition
			if err := FileManagement.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[i].Start)); err != nil {
				fmt.Println("Error: No se pudo leer el Superblock:", err)
				return
			}
			break
		}
	}

	// Find the users.txt file
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		fmt.Println("====== End RMGRP ======")
		return
	}

	var crrInode DiskStruct.Inode
	if err := FileManagement.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(DiskStruct.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo del archivo users.txt:", err)
		return
	}

	// Read the content of the users.txt file
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	fmt.Println("Contenido actual del archivo users.txt:")
	fmt.Println(data)

	// find the group to remove
	lines := strings.Split(data, "\n")
	var updatedLines []string
	groupFound := false

	// Clean the group parameter
	cleanedName := strings.TrimSpace(name)
	cleanedName = strings.ReplaceAll(cleanedName, "\u200B", "") // Delete invisible characters

	for _, line := range lines {
		// Eliminar espacios en blanco adicionales
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "\u200B", "") // Delete invisible characters
		if line == "" {
			continue // Ignorar líneas vacías
		}

		words := strings.Split(line, ",")
		fmt.Printf("Campos de la línea: %v\n", words)

		if len(words) == 3 {
			// Clean the group field
			for i := range words {
				words[i] = strings.TrimSpace(words[i])
				words[i] = strings.ReplaceAll(words[i], "\u200B", "")
			}

			// Compare the group
			fmt.Printf("Comparando grupo: '%s' con '%s'\n", words[2], cleanedName)
			if words[1] == "G" && words[2] == cleanedName {
				// Change the id of the group to 0
				words[0] = "0"
				groupFound = true
			}
			updatedLines = append(updatedLines, strings.Join(words, ","))
		} else {
			updatedLines = append(updatedLines, line)
		}
	}

	if !groupFound {
		fmt.Printf("Error: El grupo '%s' no existe.\n", cleanedName)
		fmt.Println("====== End RMGRP ======")
		return
	}

	// Update the content of the users.txt file
	newData := strings.Join(updatedLines, "\n")

	// Overwrite the file block with the updated data
	if err := OverwriteFileBlock(&crrInode, newData, file, tempSuperblock, indexInode); err != nil {
		fmt.Println("Error: No se pudo actualizar el archivo users.txt:", err)
		fmt.Println("====== End RMGRP ======")
		return
	}

	fmt.Println("Contenido actualizado del archivo users.txt:")
	fmt.Println(newData)

	fmt.Println("Grupo eliminado exitosamente.")
	fmt.Println("====== End RMGRP ======")
}
