package beans

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type Option struct {
	Addr string
}

type Serve struct {
	Conn *beanstalk.Conn
}

func NewServe(option *Option) (*Serve, error) {
	conn, err := beanstalk.Dial("tcp", option.Addr)
	if err != nil {
		return nil, err
	}
	return &Serve{Conn: conn}, nil
}

// Close beanstalk
func (s *Serve) Close() {
	_ = s.Conn.Close()
}

// Publish message to pipeline
func (s *Serve) Release(tubeName string, body []byte, priority uint32, delay, ttr time.Duration) (uint64, error) {
	if tubeName == "" {
		return 0, errors.New("tube is empty")
	}

	s.Conn.Tube.Name = tubeName
	s.Conn.TubeSet.Name[tubeName] = true

	id, err := s.Conn.Put(body, priority, delay, ttr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Receive message from pipeline
func (s *Serve) Receive(tubeName string) {
	s.Conn.Tube.Name = tubeName
	s.Conn.TubeSet.Name[tubeName] = true

	for {
		// Prevent consumption from taking up for a long time, and set timeout
		id, body, err := s.Conn.Reserve(3 * time.Second)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				log.Printf("tube=%s;err=%s", s.Conn.Tube.Name, err)
				continue
			}
			return
		}
		if id > 0 {
			if err := s.Conn.Delete(id); err != nil {
				log.Printf("Handle task exceptions:id=%d;body=%s;err=%s", id, string(body), err)
			} else {
				// todo do something...
				log.Printf("Successfully:id=%d;body=%s", id, string(body))

			}
		}
	}
}

// Watch tube
func (s *Serve) WatchTubes() ([]string, error) {
	tubes, err := s.Conn.ListTubes()
	if err != nil {
		return nil, err
	}
	return tubes, nil
}

// Statistical tube
func (s *Serve) TubeStat() (map[string]string, error) {
	stats, err := s.Conn.Tube.Stats()
	if err != nil {
		return nil, err
	}
	return stats, nil
}
