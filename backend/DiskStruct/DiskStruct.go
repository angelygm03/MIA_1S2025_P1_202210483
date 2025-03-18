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

// ==== STRUCTURES FOR EXT2 FILE SYSTEM ====

type Superblock struct {
	S_filesystem_type   int32    // Number that identifies the file system
	S_inodes_count      int32    // Total number of inodes
	S_blocks_count      int32    // Total number of blocks
	S_free_blocks_count int32    // How many blocks are free
	S_free_inodes_count int32    // How many inodes are free
	S_mtime             [17]byte // Last mount time
	S_umtime            [17]byte // Last unmount time
	S_mnt_count         int32    // How many times the disk has been mounted
	S_magic             int32    // Id of the file system: 0xEF53
	S_inode_size        int32    // Inode size
	S_block_size        int32    // Block size
	S_fist_ino          int32    // First free inode
	S_first_blo         int32    // Fist free block
	S_bm_inode_start    int32    // Starting point of the bitmap of inodes
	S_bm_block_start    int32    // Starting point of the bitmap of blocks
	S_inode_start       int32    // Starting point of the table of inodes
	S_block_start       int32    // Starting point of the table of blocks
}

type Inode struct {
	I_uid   int32     // UID of the user
	I_gid   int32     // GID of the group
	I_size  int32     // Size of the file
	I_atime [17]byte  // Last access time
	I_ctime [17]byte  // Creation time
	I_mtime [17]byte  // Last modification time
	I_block [15]int32 // Pointers to the blocks
	I_type  [1]byte   // File type: 0 for folder, 1 for file
	I_perm  [3]byte   // Permissions
}

type Folderblock struct {
	B_content [4]Content // Array of content
}

type Content struct {
	B_name  [12]byte // Name of the file or folder
	B_inodo int32    // Inodo pointer to the file or folder
}

type Fileblock struct {
	B_content [64]byte // Array of content
}

type Pointerblock struct {
	B_pointers [16]int32 // Array of pointers
}
