package DiskStruct

import (
	"fmt"
)

type MRB struct {
	MbrSize      int32    // 4 bytes
	CreationDate [10]byte // YYYY-MM-DD
	Signature    int32    // 4 bytes
	Fit          [1]byte  // 1 byte: B, F, W
	Partitions   [4]Partition
}

func PrintMBR(data MRB) {
	fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d", string(data.CreationDate[:]), string(data.Fit[:]), data.MbrSize))
	for i := 0; i < 4; i++ {
		PrintPartition(data.Partitions[i])
	}
}

type Partition struct {
	Status      [1]byte // Mounted, Unmounted
	Type        [1]byte // P, E
	Fit         [1]byte // B, F, W
	Start       int32   // Where the partition starts
	Size        int32
	Name        [16]byte
	Correlative int32
	Id          [4]byte
}

func PrintPartition(data Partition) {
	fmt.Println(fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s", string(data.Name[:]), string(data.Type[:]), data.Start, data.Size, string(data.Status[:]), string(data.Id[:])))
}

type EBR struct {
	PartMount byte //Mounted, Unmounted
	PartFit   byte //B, F, W
	PartStart int32
	PartSize  int32
	PartNext  int32 //EBR next byte, -1 if it's the last one
	PartName  [16]byte
}

func PrintEBR(data EBR) {
	fmt.Println(fmt.Sprintf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount))
}
