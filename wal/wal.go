package wal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"

	"os"
	"reflect"
	"sync"
	"time"
)

type Wal struct {
	f  *os.File
	mu sync.Mutex
}

type WalOption func(w *Wal)

// |record len (4 bytes) | timestamp (8 bytes) | hashed key (32 bytes) | compressed value len (4 bytes) | value... | crc32 (4 bytes)
func OpenWAL(opts ...WalOption) (*Wal, error) {

	f, err := os.OpenFile("wal.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	w := Wal{
		f: f,
	}

	for _, opt := range opts {
		opt(&w)
	}

	return &w, nil
}

func (w *Wal) Close() error {
	return w.f.Close()
}

func (w *Wal) Add(key, value []byte) error {
	var buf bytes.Buffer

	crc := crc32.NewIEEE()

	valueLen := len(value)
	walLen := 4 + 32 + 8 + 4 + uint32(valueLen) + 4

	if err := binary.Write(&buf, binary.LittleEndian, walLen); err != nil {
		return err
	}

	if err := binary.Write(&buf, binary.LittleEndian, time.Now().UTC().UnixMicro()); err != nil {
		return err
	}

	buf.Write(key)

	if err := binary.Write(&buf, binary.LittleEndian, uint32(valueLen)); err != nil {
		return err
	}
	buf.Write(value)

	crc.Write(buf.Bytes()[4:])
	buf.Write(crc.Sum(nil))

	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := w.f.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

type memtable interface {
	Put(key, value []byte)
}

func (w *Wal) FillMemtable(memtable memtable) error {
	f, err := os.OpenFile("wal.txt", os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	var n int
	recordLenBuf := make([]byte, 4)
	var buf []byte

	for err != io.EOF || n != 0 {
		clear(recordLenBuf)
		// n, err = io.ReadAtLeast(f, buf, 4)

		n, err = io.ReadFull(f, recordLenBuf)
		if err != nil {
			if n == 0 {
				return nil
			}
			return fmt.Errorf("read recordLen: %w", err)
		}

		recordLen := binary.LittleEndian.Uint32(recordLenBuf)

		buf = make([]byte, int(recordLen)-4)

		bufReqLen := int(recordLen) - 4
		n, err = io.ReadFull(f, buf)
		if err != nil && err != io.EOF {
			if n == 0 {
				return nil
			}
			return fmt.Errorf("read full record buf %w", err)
		}

		crc := crc32.NewIEEE()
		crc.Write(buf[:bufReqLen-4])

		if !reflect.DeepEqual(buf[bufReqLen-4:bufReqLen], crc.Sum(nil)) {
			return errors.New("crc not equal")
		}

		key := buf[8:40]
		val := buf[44 : bufReqLen-4]
		memtable.Put(key, val)
	}

	return nil
}

func (w *Wal) FileName() string {
	return w.f.Name()
}
