package main

import (
	"testing"
	"time"
)

func TestMatchingEngine(t *testing.T) {
	ob := NewOrderBook()

	// 1. Add a SELL order
	sellOrder := &Order{
		ID:        "sell-1",
		Side:      SELL,
		Price:     100.0,
		Quantity:  10,
		Timestamp: time.Now(),
	}
	trades := ob.ProcessOrder(sellOrder)
	if len(trades) != 0 {
		t.Errorf("Expected 0 trades, got %d", len(trades))
	}
	if len(ob.Asks) != 1 {
		t.Errorf("Expected 1 ask, got %d", len(ob.Asks))
	}

	// 2. Add a BUY order that matches partially
	buyOrder1 := &Order{
		ID:        "buy-1",
		Side:      BUY,
		Price:     100.0,
		Quantity:  5,
		Timestamp: time.Now(),
	}
	trades = ob.ProcessOrder(buyOrder1)
	if len(trades) != 1 {
		t.Errorf("Expected 1 trade, got %d", len(trades))
	}
	if trades[0].Quantity != 5 {
		t.Errorf("Expected trade quantity 5, got %d", trades[0].Quantity)
	}
	if ob.Asks[0].Quantity != 5 {
		t.Errorf("Expected remaining ask quantity 5, got %d", ob.Asks[0].Quantity)
	}

	// 3. Add a BUY order that completes the match and leaves remaining
	buyOrder2 := &Order{
		ID:        "buy-2",
		Side:      BUY,
		Price:     101.0,
		Quantity:  10,
		Timestamp: time.Now(),
	}
	trades = ob.ProcessOrder(buyOrder2)
	if len(trades) != 1 {
		t.Errorf("Expected 1 trade, got %d", len(trades))
	}
	if trades[0].Quantity != 5 {
		t.Errorf("Expected trade quantity 5, got %d", trades[0].Quantity)
	}
	if len(ob.Asks) != 0 {
		t.Errorf("Expected 0 asks, got %d", len(ob.Asks))
	}
	if len(ob.Bids) != 1 {
		t.Errorf("Expected 1 bid, got %d", len(ob.Bids))
	}
	if ob.Bids[0].Quantity != 5 {
		t.Errorf("Expected remaining bid quantity 5, got %d", ob.Bids[0].Quantity)
	}
}
