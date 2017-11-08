package fluentd

import (
	//"log"
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
)

type Sender struct {
	Host	string
	Port	int
	Tag		string
}


func NewSender(Host string, Port int, Tag string) *Sender {
	return &Sender{
		Host: Host,
		Port: Port,
		Tag: Tag,
	}
}

func (s *Sender) GetURL() string {
	return fmt.Sprintf("http://%s:%d/%s", s.Host, s.Port, s.Tag)
}

func (s *Sender) Send(entry *DockerLogEntry) error {
	entryBytes, err := json.Marshal(&entry)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(entryBytes)
	_, err = http.Post(s.GetURL(), "application/json", buffer)
	if err != nil {
		return err
	}
	//log.Println("Resp: ", resp)
	return nil
}

