package journald

import (
	"bufio"
	"encoding/json"
	"log"
	//"fmt"
	"os/exec"
)

type Collector struct {
	SystemdUnit	string
	BootId		string
}

func NewCollector(systemdUnit string, bootId string) *Collector {
	return &Collector{
		SystemdUnit: systemdUnit,
		BootId: bootId,
	}
}

func (c *Collector) GetCmdArgs() []string {
	args := []string{
		"--output",
		"json",
		"--follow",
		"--no-pager",
		"-u",
		c.SystemdUnit,
	}
	//if len(c.BootId) > 0 {
	//	args = append(args, "-b")
	//	args = append(args, c.BootId)
	//}
	return args
}

func (c *Collector) CollectJournal(ch chan JournalEntry) {
	args := c.GetCmdArgs()
	log.Printf("%v\n", args)
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
		}
	}
}

