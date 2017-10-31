package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/aliate/pump/fluentd"
	"github.com/aliate/pump/journald"
)

const (
	DefaultSystemdUnit = "docker.service"
	DefaultBootIdFile = "bootid.current"
)

func WriteBootId(bootId string) error {
	f, err := os.OpenFile(DefaultBootIdFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(bootId)
	if err != nil {
		return err
	}
	return nil
}

func ReadBootId() (string, error) {
	f, err := os.Open(DefaultBootIdFile)
	if err != nil {
		return "", nil
	}
	defer f.Close()
	buf := make([]byte, 64)
	_, err = f.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func main() {

	flag.Parse()

	bootId, err := ReadBootId()
	if err != nil {
		panic(err)
	}

	fmt.Println("BootId: ", bootId)

	collector := journald.NewCollector(DefaultSystemdUnit, bootId)
	sender := fluentd.NewSender("localhost", 8888, DefaultSystemdUnit)

	c := make(chan journald.JournalEntry)
	go collector.CollectJournal(c)

	for {
		select {
		case entry := <-c:
			log := fluentd.DockerLogEntry{
				Hostname:      entry.Hostname,
				ContainerTag:  entry.ContainerTag,
				ContainerId:   entry.ContainerFullId,
				ContainerName: entry.ContainerName,
				Source:        "stdout",
				Log:           entry.Message,
				Timestamp:     entry.SourceRealtimeTimestamp,
			}

			fmt.Printf("[%s] %s\n", log.ContainerName, log.Log)
			WriteBootId(entry.BootId)

			if err := sender.Send(&log); err != nil {
				panic(err)
			}
		}
	}
}
