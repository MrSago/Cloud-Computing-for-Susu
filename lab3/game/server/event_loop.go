package main

import (
	"sync"
	"time"
)

type TickEntry interface {
	Tick(difference time.Duration)
}

type EventLoop struct {
	entries  []TickEntry
	lastTick time.Time
	mutex    sync.Mutex
}

func NewEventLoop() *EventLoop {
	return &EventLoop{
		lastTick: time.Now(),
	}
}

func (l *EventLoop) Tick() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	difference := time.Since(l.lastTick)
	l.lastTick = time.Now()

	for _, t := range l.entries {
		t.Tick(difference)
	}
}

func (l *EventLoop) Add(t TickEntry) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.entries = append(l.entries, t)
}

func (l *EventLoop) Remove(t TickEntry) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for i, v := range l.entries {
		if v == t {
			l.entries = append(l.entries[:i], l.entries[i+1:]...)
		}
	}
}

func (l *EventLoop) Start() {
	for {
		time.Sleep(time.Second)
		l.Tick()
	}
}
