package generator

import (
	"container/heap"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	MaxClickTime = 60
	ImpTime      = 0.1
	ImpName      = "/tmp/impressions.pipe"
	ClickName    = "/tmp/clicks.pipe"
	ClickRate    = 0.05
)

// ClickData represents a click.
type ClickData struct {
	timestamp int64
	transID   string
}

// An ClickHeap is a min-heap of clicks.
type ClickHeap []ClickData

func (h *ClickHeap) Len() int {
	return len(*h)
}

func (h *ClickHeap) Less(i, j int) bool {
	return (*h)[i].timestamp < (*h)[j].timestamp
}

func (h *ClickHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *ClickHeap) Push(x interface{}) {
	*h = append(*h, x.(ClickData))
}

func (h *ClickHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func randomTransID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func cleanup() {
	_ = os.Remove(ImpName)
	_ = os.Remove(ClickName)
}

func FakeInfo() {
	// Remove existing pipes

	cleanup()
	defer cleanup()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	// Create impressions pipe

	err := syscall.Mkfifo(ImpName, 0600)
	if err != nil {
		log.Panicf("Error creating impressions pipe: %v", err)
	}

	impressionsFile, err := os.OpenFile(ImpName, os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("Error opening impressions pipe: %v", err)
	}

	// Create clicks pipe

	err = syscall.Mkfifo(ClickName, 0600)
	if err != nil {
		log.Panicf("Error creating clicks pipe: %v", err)
	}

	clicksFile, err := os.OpenFile(ClickName, os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("Error opening clicks pipe: %v", err)
	}

	// Main loop
	clickHeap := &ClickHeap{}

	for {
		transID, err := randomTransID()
		if err != nil {
			log.Panicf("Unable to generate random transID: %v", err)
		}

		_, err = impressionsFile.WriteString(fmt.Sprintf("%s %d\n", transID, rand.Intn(99)+1))
		if err != nil {
			log.Panicf("Unable to write line to impressions pipe: %v", err)
		}

		if rand.Float64() < ClickRate {
			heap.Push(clickHeap, ClickData{
				timestamp: time.Now().Add(time.Duration(rand.Intn(MaxClickTime*1000)) * time.Millisecond).Unix(),
				transID:   transID,
			})
		}

		for clickHeap.Len() > 0 && (*clickHeap)[0].timestamp <= time.Now().Unix() {
			click := heap.Pop(clickHeap).(ClickData)
			_, err = clicksFile.WriteString(fmt.Sprintf("%s\n", click.transID))
			if err != nil {
				log.Panicf("Unable to write line to clicks pipe: %v", err)
			}
		}

		time.Sleep(time.Duration(rand.Intn(ImpTime*1000)) * time.Millisecond)
	}
}
