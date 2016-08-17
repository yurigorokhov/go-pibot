package main

import (
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type PiCamFileSystem struct {
	imageFilePath     string
	listeningChannels []chan<- []byte
	channelLock       *sync.Mutex
}

func NewPiCamFileSystem(imageFilePath string) *PiCamFileSystem {
	return &PiCamFileSystem{
		imageFilePath:     imageFilePath,
		listeningChannels: make([]chan<- []byte, 0),
		channelLock:       &sync.Mutex{},
	}
}

func (p *PiCamFileSystem) Subscribe() <-chan []byte {
	newChannel := make(chan []byte, 1)
	p.channelLock.Lock()
	defer p.channelLock.Unlock()
	p.listeningChannels = append(p.listeningChannels, newChannel)
	return newChannel
}

func (p *PiCamFileSystem) Start(quitChan <-chan bool) {
	go func() {
		for {
			select {
			case <-quitChan:

				//TODO: send close on all channels!
				break
			default:
				p.channelLock.Lock()
				listeners := make([]chan<- []byte, len(p.listeningChannels), len(p.listeningChannels))
				listenerCount := copy(listeners, p.listeningChannels)
				p.channelLock.Unlock()
				if listenerCount == 0 {
					time.Sleep(100 * time.Millisecond)
					continue
				}

				// read image from disk
				file, err := os.Open(p.imageFilePath)
				if err != nil {
					//TODO: log.Fatal(err)

				}
				b, err := ioutil.ReadAll(file)
				if err != nil {
					//TODO: log.Fatal(err)
				}

				// send image on all channels
				for _, c := range listeners {
					select {
					case c <- b:
					default:
						//TODO: we are currently not deleting old subscribers, nor is there a limit on the number of subscribers!
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}
