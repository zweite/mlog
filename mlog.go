package mlog

import (
	"bytes"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

const (
	debugType   = "debug"
	accessType  = "access"
	noticeType  = "notice"
	recordType  = "record"
	warningType = "warning"
	errorType   = "error"
	statType    = "stat"
)

// Mlog m log struct
type Mlog struct {
	mux     sync.RWMutex
	logCfg  *LogConfig
	logFile map[string]*MlogFile
}

// NewMlog NewMlog
func NewMlog(logCfg *LogConfig) *Mlog {
	return &Mlog{
		logCfg:  logCfg,
		logFile: make(map[string]*MlogFile),
	}
}

// EnsureStatDir ensure stat dir
func (m *Mlog) EnsureStatDir(logType string) error {
	return EnsureDir(filepath.Join(m.logCfg.Dir, statType, logType), 0755)
}

func (m *Mlog) write(relPath, msg string) {
	msgDir := filepath.Join(m.logCfg.Dir, relPath)

	m.mux.Lock()
	logFile, ok := m.logFile[msgDir]
	if !ok {
		logFile = NewMlogFile(msgDir, m.logCfg.SubRelPath, m.logCfg.Buffer)
		m.logFile[msgDir] = logFile
	}

	m.mux.Unlock()
	logFile.Write(msg)
}

// Access access
func (m *Mlog) Access(logType, str string, req *http.Request) {
	relPath := filepath.Join(accessType, logType)
	msg := time.Now().Format("2006-01-02 15:04:05") + "\t" +
		str + "\t" + GetRequestIP(req) + "\t" + req.URL.String()
	m.write(relPath, msg)
}

// Error error
func (m *Mlog) Error(logType, excp, desc string, req ...*http.Request) {
	relPath := filepath.Join(errorType, logType)

	msg := m.packLogMsg(excp, desc, req...)
	m.write(relPath, msg)
}

// Warning warning
func (m *Mlog) Warning(logType, excp, desc string, req ...*http.Request) {
	relPath := filepath.Join(warningType, logType)

	msg := m.packLogMsg(excp, desc, req...)
	m.write(relPath, msg)
}

// Notice notice
func (m *Mlog) Notice(logType, excp, desc string, req ...*http.Request) {
	relPath := filepath.Join(noticeType, logType)
	msg := m.packLogMsg(excp, desc, req...)
	m.write(relPath, msg)
}

func (m *Mlog) packLogMsg(excp, desc string, reqs ...*http.Request) string {
	desc = AllTrim(desc)
	var req *http.Request
	if len(reqs) > 0 {
		req = reqs[0]
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buf.WriteByte('\t')
	buf.WriteString(excp)
	buf.WriteByte('\t')
	buf.WriteString(desc)
	buf.WriteByte('\t')
	if req != nil {
		buf.WriteString(GetRequestIP(req))
	}

	buf.WriteByte('\t')
	if req != nil {
		buf.WriteString(req.URL.String())
	}

	buf.WriteByte('\t')
	buf.WriteString(m.logCfg.ServerIP)
	return buf.String()
}

// Record record
func (m *Mlog) Record(logType, msg string) {
	relPath := filepath.Join(recordType, logType)
	m.write(relPath, time.Now().Format("2006-01-02 15:04:05")+"\t"+msg)
}

// Debug debug日志
func (m *Mlog) Debug(logType, msg string) {
	relPath := filepath.Join(debugType, logType)
	m.write(relPath, time.Now().Format("2006-01-02 15:04:05")+"\t"+msg)
}

// Stat 统计日志
func (m *Mlog) Stat(logType, msg string) {
	relPath := filepath.Join(statType, logType)
	m.write(relPath, msg)
}

// Flush Flush
func (m *Mlog) Flush() {
	m.mux.RLock()
	for _, mlogFile := range m.logFile {
		mlogFile.Flush()
	}
	m.mux.RUnlock()
}

// Close Close
func (m *Mlog) Close() {
	m.mux.RLock()
	for _, mlogFile := range m.logFile {
		mlogFile.Flush()
		mlogFile.Close()
	}
	m.mux.RUnlock()
}
