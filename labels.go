package pretty

import (
	"sync"
)

type Labels interface {
	Set(name string, value string)
	Clear()
	Get(name string) string
	Current() []Label
	Exists(name string) bool
}

type labels struct {
	mu         sync.RWMutex
	labelsMap  map[string]string
	labelNames []string
}

func (l *labels) Exists(name string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.labelsMap[name]
	return ok
}

func (l *labels) Set(name string, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.labelsMap[name]
	if ok {
		l.labelsMap[name] = value
	}
}

func (l *labels) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, k := range l.labelNames {
		l.labelsMap[k] = ""
	}
}

func (l *labels) Get(name string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.labelsMap[name]
}

func (l *labels) Current() []Label {
	l.mu.RLock()
	defer l.mu.RUnlock()
	result := make([]Label, len(l.labelsMap))
	for _, name := range l.labelNames {
		result = append(result, Label{
			Name:  name,
			Value: l.labelsMap[name],
		})
	}
	return result
}

func NewLabels(names ...string) Labels {
	m := make(map[string]string)
	for _, n := range names {
		m[n] = ""
	}
	return &labels{
		mu:         sync.RWMutex{},
		labelsMap:  m,
		labelNames: names,
	}
}
