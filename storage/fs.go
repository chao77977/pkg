package storage

import (
	"strconv"
	"syscall"
)

var _fsType = map[string]string{
	"1021994":  "TMPFS",
	"137d":     "EXT",
	"4244":     "HFS",
	"4d44":     "MSDOS",
	"52654973": "REISERFS",
	"5346544e": "NTFS",
	"58465342": "XFS",
	"61756673": "AUFS",
	"6969":     "NFS",
	"ef51":     "EXT2OLD",
	"ef53":     "EXT4",
	"f15f":     "ecryptfs",
	"794c7630": "overlayfs",
	"2fc12fc1": "zfs",
	"ff534d42": "cifs",
	"53464846": "wslfs",
}

func GetFSType(fsType int64) string {
	fsTypyHexStr := strconv.FormatInt(fsType, 16)
	fsTypeStr, ok := _fsType[fsTypyHexStr]
	if ok {
		return fsTypeStr
	}

	return "UNKNOWN"
}

type FSInfo struct {
	total  uint64 // total size of the volume / disk
	used   uint64 // free size of the volume / disk
	free   uint64 // free size of the volume / disk
	files  uint64 //  total inodes available
	ffree  uint64 // free inodes available
	fsType string // file system type
}

func GetFSInfo(path string) (*FSInfo, error) {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(path, &s)
	if err != nil {
		return nil, err
	}

	return &FSInfo{
		total:  s.Blocks * uint64(s.Bsize),
		free:   s.Bfree * uint64(s.Bsize),
		used:   (s.Blocks - s.Bfree) * uint64(s.Bsize),
		files:  s.Files,
		ffree:  s.Ffree,
		fsType: GetFSType(s.Type),
	}, nil
}

func (f *FSInfo) Total() string {
	return FormatByte(float64(f.total)).Show()
}

func (f *FSInfo) Free() string {
	return FormatByte(float64(f.free)).Show()
}

func (f *FSInfo) Used() string {
	return FormatByte(float64(f.used)).Show()
}

func (f *FSInfo) TotalBytes() uint64 {
	return f.total
}

func (f *FSInfo) FreeBytes() uint64 {
	return f.free
}

func (f *FSInfo) UsedBytes() uint64 {
	return f.used
}

func (f *FSInfo) TotalInodes() uint64 {
	return f.files
}

func (f *FSInfo) FreeInodes() uint64 {
	return f.ffree
}

func (f *FSInfo) Type() string {
	return f.fsType
}
