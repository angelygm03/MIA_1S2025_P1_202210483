package main

import (
	"Proyecto1/backend/DiskControl"
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

func main() {
	http.HandleFunc("/mkdisk", createDisk)
	http.HandleFunc("/rmdisk", removeDisk)
	http.HandleFunc("/fdisk", createPartition)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
