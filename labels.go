package pretty

import (
	"strings"
	"sync"
)

type Labels interface {
	SetIfExists(level string, name string, value string)
	Clear(level string)
	Current(level string) []Label
}

type labels struct {
	mu              sync.RWMutex
	levelsLabelsMap map[string]map[string]string
	labelNames      []string
	labelNamesMap   map[string]struct{}
}

func (l *labels) SetIfExists(level string, name string, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.labelNamesMap[name]
	if ok {
		if l.levelsLabelsMap[level] == nil {
			l.levelsLabelsMap[level] = make(map[string]string)
			for _, n := range l.labelNames {
				l.levelsLabelsMap[level][n] = ""
			}
		}
		l.levelsLabelsMap[level][name] = value
	}
}

func (l *labels) Clear(level string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.levelsLabelsMap[level] != nil {
		delete(l.levelsLabelsMap, level)
	}
}

func (l *labels) Current(currentLevel string) []Label {
	l.mu.RLock()
	defer l.mu.RUnlock()
	result := make([]Label, len(l.labelNames))
	for i, lab := range l.labelNames {
		result[i].Name = lab
	}
	breadCrumbs := strings.Split(currentLevel, sep)
	for _, level := range breadCrumbs {
		if l.levelsLabelsMap[level] != nil {
			for i, name := range l.labelNames {
				val, ok := l.levelsLabelsMap[level][name]
				if ok {
					if val != "" {
						result[i] = Label{
							Name:  name,
							Value: val,
						}
					}
				}
			}
		}
	}
	return result
}

func NewLabels(names ...string) Labels {
	m := make(map[string]struct{})
	for _, n := range names {
		m[n] = struct{}{}
	}
	return &labels{
		mu:              sync.RWMutex{},
		levelsLabelsMap: make(map[string]map[string]string),
		labelNamesMap:   m,
		labelNames:      names,
	}
}
