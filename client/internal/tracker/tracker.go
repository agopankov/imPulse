package tracker

import (
	"sync"
	"time"
)

type SymbolChange struct {
	Symbol             string
	PriceChange        string
	FirstPriceChange   string
	PriceChangePct     float64
	AddedAt            time.Time
	IsNew              bool
	NotificationOfPump bool
}

type Tracker struct {
	mu             sync.Mutex
	trackedSymbols map[string]SymbolChange
}

func NewTracker() *Tracker {
	return &Tracker{
		trackedSymbols: make(map[string]SymbolChange),
	}
}

func (t *Tracker) GetTrackedSymbols() map[string]SymbolChange {
	t.mu.Lock()
	defer t.mu.Unlock()
	copiedSymbols := make(map[string]SymbolChange)
	for k, v := range t.trackedSymbols {
		copiedSymbols[k] = v
	}
	return copiedSymbols
}

func (t *Tracker) UpdateTrackedSymbol(symbolChange SymbolChange) {
	t.mu.Lock()
	defer t.mu.Unlock()
	symbolChange.IsNew = true
	t.trackedSymbols[symbolChange.Symbol] = symbolChange
}

func (t *Tracker) RemoveTrackedSymbol(symbol string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.trackedSymbols, symbol)
}

func (t *Tracker) IsTracked(symbol string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, exists := t.trackedSymbols[symbol]
	return exists
}
