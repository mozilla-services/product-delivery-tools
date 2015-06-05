package mozlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

// DefaultLogger is the logger that will be set if UseMozLogger is run
var DefaultLogger = &MozLogger{
	Output: os.Stdout,
	Logger: "MozLog",
}

// MozLogger implements the io.Writer interface
type MozLogger struct {
	Output io.Writer
	Logger string
}

// Write converts the log to AppLog
func (m *MozLogger) Write(l []byte) (int, error) {

	log := New(m.Logger)
	log.Fields["msg"] = string(bytes.TrimSpace(l))

	out, err := log.ToJSON()
	if err != nil {
		// Need someway to notify that this happened.
		fmt.Fprintln(os.Stderr, err)
		return 0, err
	}

	_, err = m.Output.Write(append(out, '\n'))
	return len(l), err
}

// UseMozLogger sets the log.std to DefaultLogger
func UseMozLogger(logger string) {
	DefaultLogger.Logger = logger
	log.SetOutput(DefaultLogger)
	log.SetFlags(log.Lshortfile)
}

// AppLog implements Mozilla logging standard
type AppLog struct {
	Timestamp  int64
	Type       string
	Logger     string
	Hostname   string `json:",omitempty"`
	EnvVersion string
	Pid        int `json:",omitempty"`
	Severity   int `json:",omitempty"`
	Fields     map[string]string
}

// New returns an AppLog
func New(logger string) *AppLog {
	return &AppLog{
		Timestamp:  time.Now().UnixNano(),
		Type:       "app.log",
		Logger:     logger,
		Hostname:   hostname,
		EnvVersion: "2.0",
		Pid:        os.Getpid(),
		Fields:     make(map[string]string),
	}
}

// ToJSON converts a logline to JSON
func (a *AppLog) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}
