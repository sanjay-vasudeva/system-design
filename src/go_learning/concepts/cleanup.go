package concepts

import (
	"os"
	"runtime"
)

type MemoryMappedFile struct {
	data []byte
}

func NewMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Stat()
	if err != nil {
		return nil, err
	}

	conn, err := f.SyscallConn()
	if err != nil {
		return nil, err
	}

	var data []byte
	connError := conn.Control(func(fd uintptr) {
		//syscall.Mmap(int(fd), 0 int(stat.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	})

	if connError != nil {
		return nil, connError
	}

	if err != nil {
		return nil, err
	}

	mf := &MemoryMappedFile{data: data}

	cleanup := func(data []byte) {
		//syscall.Munmap(data)
	}

	runtime.AddCleanup(mf, cleanup, data)
	return mf, nil
}
