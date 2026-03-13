package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Side is the side of the order (BUY or SELL)
type Side string

const (
	BUY  Side = "BUY"
	SELL Side = "SELL"
)

// Order represents a single trade order
type Order struct {
	ID        string    `json:"id"`
	Side      Side      `json:"side"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}

// Trade represents a successful match between two orders
type Trade struct {
	BuyerOrderID  string  `json:"buyer_order_id"`
	SellerOrderID string  `json:"seller_order_id"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
}

// OrderBook maintains the state of bids and asks
type OrderBook struct {
	Bids []*Order
	Asks []*Order
	mu   sync.Mutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: []*Order{},
		Asks: []*Order{},
	}
}

// ProcessOrder adds an order to the book and attempts matching
func (ob *OrderBook) ProcessOrder(order *Order) []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var trades []Trade

	if order.Side == BUY {
		trades = ob.matchBuyOrder(order)
	} else {
		trades = ob.matchSellOrder(order)
	}

	return trades
}

func (ob *OrderBook) matchBuyOrder(order *Order) []Trade {
	var trades []Trade

	// Sort asks by price (ascending) and then timestamp (ascending) for FIFO
	sort.Slice(ob.Asks, func(i, j int) bool {
		if ob.Asks[i].Price == ob.Asks[j].Price {
			return ob.Asks[i].Timestamp.Before(ob.Asks[j].Timestamp)
		}
		return ob.Asks[i].Price < ob.Asks[j].Price
	})

	remainingQty := order.Quantity
	var updatedAsks []*Order

	for i := 0; i < len(ob.Asks); i++ {
		ask := ob.Asks[i]
		if remainingQty <= 0 || order.Price < ask.Price {
			updatedAsks = append(updatedAsks, ob.Asks[i:]...)
			break
		}

		matchQty := remainingQty
		if ask.Quantity < remainingQty {
			matchQty = ask.Quantity
		}

		trades = append(trades, Trade{
			BuyerOrderID:  order.ID,
			SellerOrderID: ask.ID,
			Price:         ask.Price, // Trade happens at the existing order's price
			Quantity:      matchQty,
		})

		remainingQty -= matchQty
		ask.Quantity -= matchQty

		if ask.Quantity > 0 {
			updatedAsks = append(updatedAsks, ask)
			updatedAsks = append(updatedAsks, ob.Asks[i+1:]...)
			break
		}
	}

	ob.Asks = updatedAsks

	if remainingQty > 0 {
		order.Quantity = remainingQty
		ob.Bids = append(ob.Bids, order)
		// Sort bids by price (descending) and then timestamp (ascending)
		sort.Slice(ob.Bids, func(i, j int) bool {
			if ob.Bids[i].Price == ob.Bids[j].Price {
				return ob.Bids[i].Timestamp.Before(ob.Bids[j].Timestamp)
			}
			return ob.Bids[i].Price > ob.Bids[j].Price
		})
	}

	return trades
}

func (ob *OrderBook) matchSellOrder(order *Order) []Trade {
	var trades []Trade

	// Sort bids by price (descending) and then timestamp (ascending)
	sort.Slice(ob.Bids, func(i, j int) bool {
		if ob.Bids[i].Price == ob.Bids[j].Price {
			return ob.Bids[i].Timestamp.Before(ob.Bids[j].Timestamp)
		}
		return ob.Bids[i].Price > ob.Bids[j].Price
	})

	remainingQty := order.Quantity
	var updatedBids []*Order

	for i := 0; i < len(ob.Bids); i++ {
		bid := ob.Bids[i]
		if remainingQty <= 0 || order.Price > bid.Price {
			updatedBids = append(updatedBids, ob.Bids[i:]...)
			break
		}

		matchQty := remainingQty
		if bid.Quantity < remainingQty {
			matchQty = bid.Quantity
		}

		trades = append(trades, Trade{
			BuyerOrderID:  bid.ID,
			SellerOrderID: order.ID,
			Price:         bid.Price,
			Quantity:      matchQty,
		})

		remainingQty -= matchQty
		bid.Quantity -= matchQty

		if bid.Quantity > 0 {
			updatedBids = append(updatedBids, bid)
			updatedBids = append(updatedBids, ob.Bids[i+1:]...)
			break
		}
	}

	ob.Bids = updatedBids

	if remainingQty > 0 {
		order.Quantity = remainingQty
		ob.Asks = append(ob.Asks, order)
		// Sort asks by price (ascending) and then timestamp (ascending)
		sort.Slice(ob.Asks, func(i, j int) bool {
			if ob.Asks[i].Price == ob.Asks[j].Price {
				return ob.Asks[i].Timestamp.Before(ob.Asks[j].Timestamp)
			}
			return ob.Asks[i].Price < ob.Asks[j].Price
		})
	}

	return trades
}

func main() {
	fmt.Println("Matching Engine: Initialized")
	ob := NewOrderBook()

	// Sample orders
	order1 := &Order{ID: "1", Side: SELL, Price: 100.5, Quantity: 10, Timestamp: time.Now()}
	order2 := &Order{ID: "2", Side: BUY, Price: 101.0, Quantity: 5, Timestamp: time.Now().Add(time.Second)}
	order3 := &Order{ID: "3", Side: BUY, Price: 100.5, Quantity: 10, Timestamp: time.Now().Add(time.Second * 2)}

	fmt.Printf("Processing Order 1: %+v\n", order1)
	trades1 := ob.ProcessOrder(order1)
	fmt.Printf("Trades: %+v\n", trades1)

	fmt.Printf("Processing Order 2: %+v\n", order2)
	trades2 := ob.ProcessOrder(order2)
	fmt.Printf("Trades: %+v\n", trades2)

	fmt.Printf("Processing Order 3: %+v\n", order3)
	trades3 := ob.ProcessOrder(order3)
	fmt.Printf("Trades: %+v\n", trades3)

	fmt.Printf("Final Order Book - Bids: %d, Asks: %d\n", len(ob.Bids), len(ob.Asks))
}
