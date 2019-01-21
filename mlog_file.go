package mlog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MlogFile m log file struct
type MlogFile struct {
	wg       sync.WaitGroup
	mux      sync.Mutex
	bufWrite *bufio.Writer
	fd       *os.File
	lastTime time.Time
	msgChan  chan string
	dir      string
	relPath  string
}

// NewMlogFile new mlog file
func NewMlogFile(dir, relPath string, buffer int) *MlogFile {
	if buffer <= 0 {
		buffer = 100
	}
	mlogFile := &MlogFile{
		dir:     dir,
		relPath: relPath,
		msgChan: make(chan string, buffer),
	}

	mlogFile.wg.Add(1)
	go mlogFile.write()
	return mlogFile
}

func (m *MlogFile) Write(msg string) {
	m.msgChan <- msg
}

func (m *MlogFile) write() {
	defer m.wg.Done()

	defer func() {
		m.Flush()
	}()

	if err := m.checkFile(); err != nil {
		fmt.Printf("dir:%s relpath:%s checkfile err:%s\n", m.dir, m.relPath, err.Error())
		panic(err)
	}

	tick := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-tick.C:
			m.Flush()

			if err := m.checkFile(); err != nil {
				fmt.Printf("dir:%s relpath:%s checkfile err:%s\n", m.dir, m.relPath, err.Error())
				break
			}
		case msg, ok := <-m.msgChan:
			if !ok {
				// 关掉关闭退出携程，防止携程泄露
				return
			}

			if err := m.writeString(msg + "\n"); err != nil {
				fmt.Printf("dir:%s relpath:%s write string err:%s\n", m.dir, m.relPath, err.Error())
				continue
			}
		}
	}
	return
}

func (m *MlogFile) checkFile() (err error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	now := time.Now()
	if m.fd == nil ||
		now.Hour() != m.lastTime.Hour() ||
		now.Sub(m.lastTime) > time.Hour*1 {

		fd, err := m.makeFile(now)
		if err != nil {
			return err
		}

		if m.bufWrite != nil {
			m.bufWrite.Flush()
		}
		if m.fd != nil {
			m.fd.Close()
		}

		m.fd = fd
		m.bufWrite = bufio.NewWriter(m.fd)
		m.lastTime = now
	}
	return
}

func (m *MlogFile) makeFile(t time.Time) (fd *os.File, err error) {
	relPath := t.Format(m.relPath)
	absPath := filepath.Join(m.dir, relPath)

	if err = EnsureDir(filepath.Dir(absPath), 0755); err != nil {
		return
	}

	return os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func (m *MlogFile) writeString(s string) (err error) {
	m.mux.Lock()
	if m.bufWrite != nil {
		_, err = m.bufWrite.WriteString(s)
	}
	m.mux.Unlock()
	return
}

func (m *MlogFile) writeByte(b byte) (err error) {
	m.mux.Lock()
	if m.bufWrite != nil {
		err = m.bufWrite.WriteByte(b)
	}
	m.mux.Unlock()
	return
}

// Flush flush
func (m *MlogFile) Flush() {
	m.mux.Lock()
	if m.bufWrite != nil {
		m.bufWrite.Flush()
	}
	m.mux.Unlock()
}

// Close close
func (m *MlogFile) Close() error {
	close(m.msgChan)
	m.wg.Wait()

	m.Flush()
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.fd != nil {
		return m.fd.Close()
	}
	return nil
}
