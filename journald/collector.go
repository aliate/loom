package journald

import (
	"bufio"
	"encoding/json"
	"log"
	//"fmt"
	"os/exec"
	"io/ioutil"
)

const DefaultCursorFile = "cursor.current"

type Collector struct {
	SystemdUnit	string
}

func NewCollector(systemdUnit string) *Collector {
	return &Collector{
		SystemdUnit: systemdUnit,
	}
}

func (c *Collector) GetCmdArgs() []string {
	args := []string{
		"--output",
		"json",
		"--follow",
		"--unit",
		c.SystemdUnit,
	}
	cursor := c.GetCursor()
	if len(cursor) > 0 {
		args = append(args, "--after-cursor")
		args = append(args, cursor)
	}
	return args
}

func (c *Collector) CollectJournal(ch chan JournalEntry) {
	args := c.GetCmdArgs()
	cmd := exec.Command("journalctl", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Could not run journalctl: %v\n", err)
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var entry JournalEntry
		err := json.Unmarshal(scanner.Bytes(), &entry)
		if err != nil {
			// Ignore blank lines
		} else {
			ch <- entry
			c.SaveCursor(entry.Cursor)
		}
	}
}

func (c *Collector) GetCursor() string {
	contents, err := ioutil.ReadFile(DefaultCursorFile)
	if err != nil {
		return ""
	}
	return string(contents)
}

func (c *Collector) SaveCursor(cursor string) error {
	b := []byte(cursor)
	return ioutil.WriteFile(DefaultCursorFile, b, 0755)
}

