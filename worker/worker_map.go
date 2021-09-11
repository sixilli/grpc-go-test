package worker

import (
	"fmt"
	"sync"
)

type threadSafeWorkerMap struct {
	workerMap map[string]*workerType
	mutex     sync.RWMutex
}

var workerMap threadSafeWorkerMap = threadSafeWorkerMap{
	workerMap: make(map[string]*workerType),
}

func (wm *threadSafeWorkerMap) get(uuid string) (*workerType, error) {
	wm.mutex.RLock()
	defer wm.mutex.Unlock()

	worker, ok := wm.workerMap[uuid]
	if !ok {
		return worker, fmt.Errorf("could not find worker in map with ID: %s", uuid)
	}

	return worker, nil
}

func (wm *threadSafeWorkerMap) set(worker *workerType) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	wm.workerMap[worker.uuid] = worker
}
