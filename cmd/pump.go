package main

import (
	"fmt"
	"flag"

	"github.com/aliate/pump/fluentd"
	"github.com/aliate/pump/journald"
)

const (
	DefaultSystemdUnit = "docker.service"
)


func main() {

	flag.Parse()

	collector := journald.NewCollector(DefaultSystemdUnit)
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

			if err := sender.Send(&log); err != nil {
				panic(err)
			}
		}
	}
}
