/*
 * Author: fasion
 * Created time: 2026-03-06 23:47:17
 * Last Modified by: fasion
 * Last Modified time: 2026-03-09 00:45:34
 */

package baseutils

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/fasionchan/goutils/stl"
)

const (
	SpillBufferMemoryWaterMark = 1024 * 1024 * 1 // 1MB 内存水位线
)

type SpillBuffer struct {
	memory          *bytes.Buffer
	rfile           *os.File
	file            *os.File
	mutex           *sync.RWMutex
	memoryWaterMark int
}

func NewSpillBuffer(memoryWaterMark int, withMutex bool) *SpillBuffer {
	if memoryWaterMark <= 0 {
		memoryWaterMark = SpillBufferMemoryWaterMark
	}

	var mutex *sync.RWMutex
	if withMutex {
		mutex = &sync.RWMutex{}
	}

	return &SpillBuffer{
		memory:          bytes.NewBuffer(nil),
		memoryWaterMark: memoryWaterMark,
		mutex:           mutex,
	}
}

func (sb *SpillBuffer) Read(p []byte) (n int, err error) {
	if mutex := sb.mutex; mutex != nil {
		mutex.RLock()
		defer mutex.RUnlock()
	}

	if rfile := sb.rfile; rfile != nil {
		return rfile.Read(p)
	}

	return sb.memory.Read(p)
}

func (sb *SpillBuffer) Write(p []byte) (n int, err error) {
	if mutex := sb.mutex; mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	if file := sb.file; file != nil {
		return file.Write(p)
	}

	n, err = sb.memory.Write(p)
	if err != nil {
		return
	}

	err = sb.switchToFile()

	return
}

func (sb *SpillBuffer) switchToFile() error {
	if total := sb.memory.Len(); total > sb.memoryWaterMark {
		tempFile, err := os.CreateTemp("", "spill-buffer-*.tmp")
		if err != nil {
			return err
		}

		written, err := sb.memory.WriteTo(tempFile)
		if err != nil {
			return err
		}

		if written != int64(total) {
			return fmt.Errorf("write to file failed: %d != %d", written, int64(total))
		}

		rfile, err := os.Open(tempFile.Name())
		if err != nil {
			return err
		}

		sb.file = tempFile
		sb.rfile = rfile

		sb.memory = nil
	}

	return nil
}

func (sb *SpillBuffer) Close() error {
	if mutex := sb.mutex; mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	if memory := sb.memory; memory != nil {
		sb.memory = nil
		return nil
	}

	var errs stl.Errors
	if rfile := sb.rfile; rfile != nil {
		if err := rfile.Close(); err != nil {
			errs = errs.Append(err)
		} else {
			sb.rfile = nil
		}
	}

	if file := sb.file; file != nil {
		if err := os.Remove(file.Name()); err != nil {
			errs = errs.Append(err)
		}

		if err := file.Close(); err != nil {
			errs = errs.Append(err)
		} else {
			sb.file = nil
		}
	}

	return errs.Simplify()
}

func (sb *SpillBuffer) WithMemoryWaterMark(waterMark int) *SpillBuffer {
	if waterMark <= 0 {
		waterMark = SpillBufferMemoryWaterMark
	}

	sb.memoryWaterMark = waterMark
	return sb
}

func (sb *SpillBuffer) WithMutex() *SpillBuffer {
	sb.mutex = &sync.RWMutex{}
	return sb
}
