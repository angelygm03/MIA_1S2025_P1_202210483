package main

import (
	"Proyecto1/backend/DiskCommands"
	"Proyecto1/backend/DiskControl"
	"Proyecto1/backend/FileSystem"
	"Proyecto1/backend/UserManagement"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ====== JSON Request ======
type MKDISKRequest struct {
	Path string `json:"path"`
	Size int    `json:"size"`
	Unit string `json:"unit"`
	Fit  string `json:"fit"`
}

type RMDISKRequest struct {
	Path string `json:"path"`
}

type FDISKRequest struct {
	Size int    `json:"size"`
	Path string `json:"path"`
	Name string `json:"name"`
	Unit string `json:"unit"`
	Type string `json:"type"`
	Fit  string `json:"fit"`
}

type MountRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type ReportRequest struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Id         string `json:"id"`
	PathFileLs string `json:"pathFileLs"`
}

type MkfsRequest struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"type"`
	Id       string `json:"id"`
}

type MkusrRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Grp  string `json:"grp"`
}

type MkgrpRequest struct {
	Name string `json:"name"`
}

// ====== CORS ======
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ====== Handlers ======
func createDisk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Decodify the JSON
	var req MKDISKRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para crear disco:", req)

	// Call the function to create the disk
	DiskControl.Mkdisk(req.Size, req.Fit, req.Unit, req.Path)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Disk created successfully at %s", req.Path)))
}

func removeDisk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req RMDISKRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Verify if the file exists
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		fmt.Println("Error: El archivo no existe en la ruta especificada")
		http.Error(w, "Error: Disk not found", http.StatusNotFound)
		return
	}

	// Remove the file
	err := os.Remove(req.Path)
	if err != nil {
		fmt.Println("Error al eliminar el archivo:", err)
		http.Error(w, "Error deleting disk", http.StatusInternalServerError)
		return
	}

	fmt.Println("Archivo eliminado exitosamente:", req.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Disk removed successfully at %s", req.Path)))
}

func createPartition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req FDISKRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para crear partición:", req)

	// Call the function to create the partition
	DiskControl.Fdisk(req.Size, req.Path, req.Name, req.Unit, req.Type, req.Fit)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Partition created successfully at %s", req.Path)))
}

func mountPartition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req MountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para crear partición:", req)

	// Call the function to create the partition
	DiskControl.Mount(req.Path, req.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Partition mounted successfully at %s", req.Path)))
}

func generateReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para generar reporte:")
	fmt.Printf("Name: %s, Path: %s, ID: %s, PathFileLs: %s\n", req.Name, req.Path, req.Id, req.PathFileLs)

	if req.Id == "" {
		http.Error(w, "Error: 'id' es un parámetro obligatorio", http.StatusBadRequest)
		return
	}

	reportCommand := fmt.Sprintf("-name=%s -path=%s -id=%s", req.Name, req.Path, req.Id)
	if req.PathFileLs != "" {
		reportCommand += fmt.Sprintf(" -path_file_ls=%s", req.PathFileLs)
	}

	DiskCommands.Fn_Rep(reportCommand)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Reporte generado exitosamente en %s", req.Path)))
}

func formatMkfs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req MkfsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para formatear partición:", req)
	fsType := "2fs"

	FileSystem.Mkfs(req.Id, req.Type, fsType)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Partition formatted successfully with id %s", req.Id)))
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para loggear usuario:", req)

	// Call the function to login the user
	UserManagement.Login(req.User, req.Password, req.Id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User logged in successfully with id %s", req.Id)))
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("Solicitud recibida para desloguear usuario")

	// Call the function to logout the user
	UserManagement.Logout()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out successfully"))
}

func getMountedPartitionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	mountedPartitions := DiskControl.GetMountedPartitions()

	// Convert the data to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mountedPartitions); err != nil {
		http.Error(w, "Error al generar JSON", http.StatusInternalServerError)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decodify the JSON
	var req MkusrRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para crear usuario:", req)

	// Validate the request
	if req.User == "" || req.Pass == "" || req.Grp == "" {
		http.Error(w, "Error: Los parámetros 'user', 'pass' y 'grp' son obligatorios.", http.StatusBadRequest)
		return
	}

	if len(req.User) > 10 || len(req.Pass) > 10 || len(req.Grp) > 10 {
		http.Error(w, "Error: Los valores de 'user', 'pass' y 'grp' no pueden exceder los 10 caracteres.", http.StatusBadRequest)
		return
	}

	command := fmt.Sprintf("-user=%s -pass=%s -grp=%s", req.User, req.Pass, req.Grp)
	DiskCommands.AnalyzeCommand("mkusr", command)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Usuario '%s' creado exitosamente en el grupo '%s'", req.User, req.Grp)))
	UserManagement.PrintUsersFile()
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req MkgrpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para crear grupo:", req)

	if req.Name == "" {
		http.Error(w, "Error: El parámetro 'name' es obligatorio.", http.StatusBadRequest)
		return
	}

	// Call the function to create the group
	UserManagement.Mkgrp(req.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Grupo '%s' creado exitosamente.", req.Name)))
}

func removeUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("Solicitud recibida para eliminar usuario:", req)

	if req.User == "" {
		http.Error(w, "Error: El parámetro 'user' es obligatorio.", http.StatusBadRequest)
		return
	}

	UserManagement.Rmusr(req.User)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Usuario '%s' eliminado exitosamente.", req.User)))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/mkdisk", createDisk)
	mux.HandleFunc("/rmdisk", removeDisk)
	mux.HandleFunc("/fdisk", createPartition)
	mux.HandleFunc("/mount", mountPartition)
	mux.HandleFunc("/report", generateReport)
	mux.HandleFunc("/mkfs", formatMkfs)
	mux.HandleFunc("/login", loginUser)
	mux.HandleFunc("/logout", logoutUser)
	mux.HandleFunc("/list-mounted", getMountedPartitionsHandler)
	mux.HandleFunc("/mkusr", createUserHandler)
	mux.HandleFunc("/mkgrp", createGroupHandler)
	mux.HandleFunc("/rmusr", removeUserHandler)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", enableCORS(mux))
}
