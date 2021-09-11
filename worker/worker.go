package worker

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Using a weird name to be more clear within the code
type workerType struct {
	uuid      string
	pid       uint64
	status    string // Could make a status struct to act as an enum (Completed, Killed, Running)
	history   []string
	listeners []chan string
	mutex     sync.RWMutex
}

// Start creates a new worker and gives it the command
func Start(command string) (chan string, error) {
	if len(command) == 0 {
		return nil, errors.New("invalid command, must have length greater than 0")
	}

	workerChannel := make(chan string)

	newWorker := &workerType{
		uuid:      uuid.New().String(),
		pid:       0,
		status:    "Running",
		history:   make([]string, 0),
		listeners: []chan string{workerChannel},
	}

	workerMap.set(newWorker)

	newWorker.runCommand(command)

	return workerChannel, nil
}

func Stop(uuid string) string {
	worker, err := workerMap.get(uuid)
	if err != nil {
		return fmt.Sprintf("Could not find worker with UUID: %s", uuid)
	}

	if worker.status == "Running" {
		// send kill command
		return fmt.Sprintf("Killing process: %d", worker.pid)
	}

	return fmt.Sprintf("Could not kill worker because it has a status of %s", worker.status)
}

func Query(uuid string) (string, error) {
	worker, err := workerMap.get(uuid)
	if err != nil {
		return "", fmt.Errorf("could not find worker with UUID: %s", uuid)
	}

	workerReport := worker.report()

	return workerReport, nil
}

// Listen - Creates a new channel, loads it with history, appends to worker's listeners, then returns it for use.
func Listen(uuid string) (chan string, error) {
	newChannel := make(chan string)

	worker, err := workerMap.get(uuid)
	if err != nil {
		return nil, fmt.Errorf("could not find worker with UUID: %s", uuid)
	}

	worker.newListener(newChannel)

	return newChannel, nil
}

// runCommand - High level function that will run until the process is done
func (w *workerType) runCommand(command string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	// Get PID, search run command and get process ID
	w.pid = 1234

	// stays alive until process exits
	go func() {
		// Get terminal output
		terminalOutput := "meow meow meow"
		w.appendHistory(terminalOutput)
		w.messageListeners(terminalOutput)
	}()

	w.status = "Completed"
}

func (w *workerType) messageListeners(message string) {
	for _, channel := range w.listeners {
		channel <- message
	}
}

func (w *workerType) newListener(channel chan string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Load up channel
	for _, item := range w.history {
		channel <- item
	}

	w.listeners = append(w.listeners, channel)
}

func (w *workerType) report() string {
	return fmt.Sprintf("Worker has a PID of %d, the current process is %s, the history has a length of %d, and %d listeners", w.pid, w.status, len(w.history), len(w.listeners))
}

func (w *workerType) appendHistory(message string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.history = append(w.history, message)
}
