package DiskCommands

import (
	"Proyecto1/backend/DiskControl"
	"bufio"
	"flag"
	"fmt"
	"os"
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
	} else if strings.Contains(command, "rep") {
		fmt.Print("COMANDO REP")
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

	// Check the fit: bf, ff, or wf. If not, FF for default
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: Fit must be 'bf', 'ff', or 'wf'")
		return
	}

	// Check the unit: k or m. If not m for default
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be 'k' or 'm'")
		return
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
