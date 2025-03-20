package DiskCommands

import (
	"Proyecto1/backend/DiskControl"
	"Proyecto1/backend/DiskStruct"
	"Proyecto1/backend/FileManagement"
	"Proyecto1/backend/UserManagement"
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

// Fuction to get the command and its parameters
func GetCommand(input string) (string, string) {
	parts := strings.Fields(input) // Split the input into parts
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])   // Get the command in lowercase
		params := strings.Join(parts[1:], " ") // Join the rest of the parts into a string
		return command, params
	}
	return "", input
}

func Analyze() {

	for true {
		var input string
		fmt.Println("======================")
		fmt.Println("Ingrese comando: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := GetCommand(input)

		fmt.Println("Comando: ", command, " - ", "Parametro: ", params)

		AnalyzeCommand(command, params)

		//mkdisk -size=3000 -unit=K -fit=BF -path="/home/angely-gmartinez/Disks/disk1.bin"
	}
}

func AnalyzeCommand(command string, params string) {
	// Check the command
	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params) // Call the function mkdisk
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params) // Call the function rmdisk
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params) // Call the function fdisk
	} else if strings.Contains(command, "mount") {
		fn_mount(params) // Call the function mount
	} else if strings.Contains(command, "rep") {
		Fn_Rep(params) // Call the function rep
	} else if strings.Contains(command, "login") {
		fn_login(params) // Call the function login
	} else {
		fmt.Println("Error: Comando inválido o no encontrado")
	}
}

func fn_mkdisk(params string) {
	// Definir flag
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Find all the flags
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Get the flag name: size, fit, unit, or path
		flagValue := strings.ToLower(match[2]) // match[2]: Get the flag value in lowercase

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// ====== Check the flags ======

	// Check the size: positive and greater than 0
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// Check the fit: bf, ff, or wf
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: Fit must be 'bf', 'ff', or 'wf'")
		return
	}

	//If fit is empty, set it to "ff"
	if *fit == "" {
		*fit = "ff"
	}

	// Check the unit: k or m
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be 'k' or 'm'")
		return
	}

	//If unit is empty, set it to "m"
	if *unit == "" {
		*unit = "m"
	}

	// Check the path: not empty
	if *path == "" {
		fmt.Println("Error: Path is required")
		return
	}

	DiskControl.Mkdisk(*size, *fit, *unit, *path)
}

// Function to remove a disk
func fn_rmdisk(params string) {
	// Define flag
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Find all the flags
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Get the flag name: path
		flagValue := strings.ToLower(match[2]) // match[2]: Get the flag value in lowercase

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Check the path: not empty
	if *path == "" {
		fmt.Println("Error: Path is required")
		return
	}

	DiskControl.Rmdisk(*path)
}

func fn_fdisk(input string) {
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "k", "Unidad")
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "wf", "Ajuste")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	if *path == "" {
		fmt.Println("Error: Path is required")
		return
	}

	// If fit is empty, set it to "w"
	if *fit == "" {
		*fit = "wf"
	}

	//Fit must be 'bf', 'ff', or 'ww'
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: Fit must be 'bf', 'ff', or 'wf'")
		return
	}

	// If unit is empty, set it to "k"
	if *unit == "" {
		*unit = "k"
	}

	//Unit must be 'k', 'm'or 'b'
	if *unit != "k" && *unit != "m" && *unit != "b" {
		fmt.Println("Error: Unit must be 'b', 'm' or 'k'")
		return
	}

	// If type is empty, set it to "p"
	if *type_ == "" {
		*type_ = "p"
	}

	// Type must be 'p', 'e', or 'l'
	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		fmt.Println("Error: Type must be 'p', 'e', or 'l'")
		return
	}

	// Call the function
	DiskControl.Fdisk(*size, *path, *name, *unit, *type_, *fit)
}

func fn_mount(params string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre de la partición")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		fs.Set(flagName, flagValue)
	}

	if *path == "" || *name == "" {
		fmt.Println("Error: Path y Name son obligatorios")
		return
	}

	// Convertir el nombre a minúsculas antes de pasarlo al Mount
	lowercaseName := strings.ToLower(*name)
	DiskControl.Mount(*path, lowercaseName)
}

func Fn_Rep(input string) {
	fmt.Println("======Start REP======")
	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	name := fs.String("name", "", "Nombre del reporte a generar (mbr, disk, inode, block, bm_inode, bm_block, sb, file, ls)")
	path := fs.String("path", "", "Ruta donde se generará el reporte")
	id := fs.String("id", "", "ID de la partición")
	pathFileLs := fs.String("path_file_ls", "", "Nombre del archivo o carpeta para reportes file o ls")

	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.Trim(match[2], "\"")

		switch flagName {
		case "name", "path", "id", "path_file_ls":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag no encontrada:", flagName)
			fmt.Println("======FIN REP======")
		}
	}

	// Name, path and id are required
	if *name == "" || *path == "" || *id == "" {
		fmt.Println("Error: 'name', 'path' y 'id' son parámetros obligatorios.")
		fmt.Println("======FIN REP======")
		return
	}

	// Verifying if the partition is mounted
	mounted := false
	var diskPath string
	for _, partitions := range DiskControl.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.ID == *id {
				mounted = true
				diskPath = partition.Path
				break
			}
		}
	}

	if !mounted {
		fmt.Println("Error: La partición con ID", *id, "no está montada.")
		fmt.Println("======FIN REP======")
		return
	}

	// Creating the reports directory if it doesn't exist
	reportsDir := filepath.Dir(*path)
	err := os.MkdirAll(reportsDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error al crear la carpeta:", reportsDir)
		fmt.Println("======FIN REP======")
		return
	}

	switch *name {

	// ===== MBR REPORT =====
	case "mbr":
		// Create the directory if it doesn't exist
		dir := filepath.Dir(*path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755) // Crear el directorio con permisos 0755
			if err != nil {
				fmt.Printf("Error al crear el directorio: %v\n", err)
				fmt.Println("======FIN REP======")
				return
			}
		}

		// Open the binary file of the mounted disk
		file, err := FileManagement.OpenFile(diskPath)
		if err != nil {
			fmt.Println("Error: No se pudo abrir el archivo en la ruta:", diskPath)
			fmt.Println("======FIN REP======")
			return
		}
		defer file.Close()

		// Read the MBR object from the binary file
		var TempMBR DiskStruct.MRB
		if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
			fmt.Println("Error: No se pudo leer el MBR desde el archivo")
			fmt.Println("======FIN REP======")
			return
		}

		// Read and process the EBRs if there are extended partitions
		var ebrs []DiskStruct.EBR
		for i := 0; i < 4; i++ {
			if string(TempMBR.Partitions[i].Type[:]) == "e" { // Extended partition: e
				fmt.Println("Partición extendida encontrada: ", string(TempMBR.Partitions[i].Name[:]))

				// First EBR position
				ebrPosition := TempMBR.Partitions[i].Start
				ebrCounter := 1

				// Read all the EBRs in the extended partition
				for ebrPosition != -1 {
					fmt.Printf("Leyendo EBR en posición: %d\n", ebrPosition)
					var tempEBR DiskStruct.EBR
					if err := FileManagement.ReadObject(file, &tempEBR, int64(ebrPosition)); err != nil {
						fmt.Println("Error: No se pudo leer el EBR desde el archivo")
						fmt.Println("======FIN REP======")
						break
					}

					// Add the EBR to the slice
					ebrs = append(ebrs, tempEBR)
					fmt.Printf("EBR %d leído. Start: %d, Size: %d, Next: %d, Name: %s\n", ebrCounter, tempEBR.PartStart, tempEBR.PartSize, tempEBR.PartNext, string(tempEBR.PartName[:]))
					DiskStruct.PrintEBR(tempEBR)

					// Move to the next EBR
					ebrPosition = tempEBR.PartNext
					ebrCounter++

					if ebrPosition == -1 {
						fmt.Println("No hay más EBRs en esta partición extendida.")
					}
				}
			}
		}

		// Generate the .dot file of the MBR
		reportPath := *path
		if err := FileManagement.GenerateMBRReport(TempMBR, ebrs, reportPath, file); err != nil {
			fmt.Println("Error al generar el reporte MBR:", err)
			fmt.Println("======FIN REP======")
		} else {
			fmt.Println("Reporte MBR generado exitosamente en:", reportPath)

			dotFile := strings.TrimSuffix(reportPath, filepath.Ext(reportPath)) + ".dot"
			outputJpg := reportPath
			cmd := exec.Command("dot", "-Tjpg", dotFile, "-o", outputJpg)
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error al renderizar el archivo .dot a imagen:", err)
				fmt.Println("======FIN REP======")
			} else {
				fmt.Println("Imagen generada exitosamente en:", outputJpg)
			}
		}

	// ===== DISK REPORT =====
	case "disk":
		dir := filepath.Dir(*path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755) // Crear el directorio con permisos 0755
			if err != nil {
				fmt.Printf("Error al crear el directorio: %v\n", err)
				fmt.Println("======FIN REP======")
				return
			}
		}
		// Open the binary file of the mounted disk
		file, err := FileManagement.OpenFile(diskPath)
		if err != nil {
			fmt.Println("Error: No se pudo abrir el archivo en la ruta:", diskPath)
			fmt.Println("======FIN REP======")
			return
		}
		defer file.Close()

		// Read the MBR object from the binary file
		var TempMBR DiskStruct.MRB
		if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
			fmt.Println("Error: No se pudo leer el MBR desde el archivo")
			fmt.Println("======FIN REP======")
			return
		}

		// Read and process the EBRs if there are extended partitions
		var ebrs []DiskStruct.EBR
		for i := 0; i < 4; i++ {
			if string(TempMBR.Partitions[i].Type[:]) == "e" { // Partición extendida
				ebrPosition := TempMBR.Partitions[i].Start
				for ebrPosition != -1 {
					var tempEBR DiskStruct.EBR
					if err := FileManagement.ReadObject(file, &tempEBR, int64(ebrPosition)); err != nil {
						break
					}
					ebrs = append(ebrs, tempEBR)   // Add the EBR to the slice
					ebrPosition = tempEBR.PartNext // Move to the next EBR
				}
			}
		}

		// Calculate the total disk size
		totalDiskSize := TempMBR.MbrSize

		// Generates the .dot file
		reportPath := *path
		if err := FileManagement.GenerateDiskReport(TempMBR, ebrs, reportPath, file, totalDiskSize); err != nil {
			fmt.Println("Error al generar el reporte DISK:", err)
			fmt.Println("======FIN REP======")
		} else {
			fmt.Println("Reporte DISK generado exitosamente en:", reportPath)

			dotFile := strings.TrimSuffix(reportPath, filepath.Ext(reportPath)) + ".dot"
			outputJpg := reportPath
			cmd := exec.Command("dot", "-Tjpg", dotFile, "-o", outputJpg)
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error al renderizar el archivo .dot a imagen:", err)
				fmt.Println("======FIN REP======")
			} else {
				fmt.Println("Imagen generada exitosamente en:", outputJpg)
			}
		}
	case "bm_inode":
		dir := filepath.Dir(*path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Error al crear el directorio: %v\n", err)
				fmt.Println("======FIN REP======")
				return
			}
		}

		// Verify if the partition is mounted
		mounted := false
		var diskPath string
		for _, partitions := range DiskControl.GetMountedPartitions() {
			for _, partition := range partitions {
				if partition.ID == *id {
					mounted = true
					diskPath = partition.Path
					break
				}
			}
			if mounted {
				break
			}
		}

		if !mounted {
			fmt.Printf("No se encontró la partición con el ID: %s.\n", *id)
			fmt.Println("======FIN REP======")
			return
		}

		// Open the binary file of the mounted disk
		file, err := FileManagement.OpenFile(diskPath)
		if err != nil {
			fmt.Printf("No se pudo abrir el archivo en la ruta: %s\n", diskPath)
			fmt.Println("======FIN REP======")
			return
		}
		defer file.Close()

		// Read the MBR object from the bin file
		var TempMBR DiskStruct.MRB
		if err := FileManagement.ReadObject(file, &TempMBR, 0); err != nil {
			fmt.Println("No se pudo leer el MBR desde el archivo.")
			fmt.Println("======FIN REP======")
			return
		}

		// Find the partition with the given ID
		var index int = -1
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				if strings.Contains(string(TempMBR.Partitions[i].Id[:]), *id) {
					if TempMBR.Partitions[i].Status[0] == '1' {
						index = i
					} else {
						fmt.Printf("La partición con el ID:%s no está montada.\n", *id)
						fmt.Println("======FIN REP======")
						return
					}
					break
				}
			}
		}

		if index == -1 {
			fmt.Printf("No se encontró la partición con el ID: %s.\n", *id)
			fmt.Println("======FIN REP======")
			return
		}

		// Read the SuperBlock
		var TemporalSuperBloque DiskStruct.Superblock
		if err := FileManagement.ReadObject(file, &TemporalSuperBloque, int64(TempMBR.Partitions[index].Start)); err != nil {
			fmt.Println("Error al leer el SuperBloque.")
			fmt.Println("======FIN REP======")
			return
		}

		// Check the values of the SuperBlock
		if TemporalSuperBloque.S_inodes_count <= 0 || TemporalSuperBloque.S_bm_inode_start <= 0 {
			fmt.Println("Valores inválidos en el SuperBloque.")
			fmt.Println("======FIN REP======")
			return
		}

		// Read the bitmap of inodes
		BitMapInode := make([]byte, TemporalSuperBloque.S_inodes_count)
		if _, err := file.ReadAt(BitMapInode, int64(TemporalSuperBloque.S_bm_inode_start)); err != nil {
			fmt.Println("No se pudo leer el bitmap de inodos:", err)
			fmt.Println("======FIN REP======")
			return
		}

		// Create the report file
		SalidaArchivo, err := os.Create(*path)
		if err != nil {
			fmt.Println("No se pudo crear el archivo de reporte:", err)
			fmt.Println("======FIN REP======")
			return
		}

		// Close the file
		defer SalidaArchivo.Close()

		// Write the bitmap of inodes to the report file
		for i, bit := range BitMapInode {
			if bit != 0 && bit != 1 {
				fmt.Printf("Advertencia: Valor inesperado en el bitmap de inodos: %d\n", bit)
				fmt.Println("======FIN REP======")
				continue
			}
			if i > 0 && i%20 == 0 {
				fmt.Fprintln(SalidaArchivo)
			}
			fmt.Fprintf(SalidaArchivo, "%d ", bit)
		}

		fmt.Printf("Reporte de BITMAP INODE de la partición:%s generado con éxito en la ruta: %s\n", *id, *path)
		fmt.Println("======FIN REP======")

	// ===== FILE -LS REPORT =====
	case "file", "ls":
		// Parameter 'path_file_ls' is required for these reports
		if *pathFileLs == "" {
			fmt.Println("Error: 'path_file_ls' es obligatorio para los reportes 'file' y 'ls'.")
			fmt.Println("======FIN REP======")
			return
		}

		fmt.Println("Generando reporte", *name, "con archivo/carpeta:", *pathFileLs)
		// ------ TO DO ------
		// Implement the report generation for 'file' and 'ls'

	default:
		fmt.Println("Error: Tipo de reporte no válido.")
		fmt.Println("======FIN REP======")
	}
}

func fn_login(input string) {
	fmt.Println("======Start LOGIN======")
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "Id")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			fmt.Println("======FIN LOGIN======")
		}
	}

	UserManagement.Login(*user, *pass, *id)

}
