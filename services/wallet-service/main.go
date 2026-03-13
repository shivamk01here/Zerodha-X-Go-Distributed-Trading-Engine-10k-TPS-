package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrWalletNotFound    = errors.New("wallet not found")
)

type TransactionType string

const (
	CREDIT TransactionType = "CREDIT"
	DEBIT  TransactionType = "DEBIT"
)

// Wallet represents a user's financial account
type Wallet struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

// LedgerEntry records a transaction in the wallet
type LedgerEntry struct {
	ID          string          `json:"id"`
	WalletID    string          `json:"wallet_id"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}

// WalletService handles balance operations
type WalletService struct {
	wallets map[string]*Wallet
	ledger  []LedgerEntry
	mu      sync.Mutex
}

func NewWalletService() *WalletService {
	return &WalletService{
		wallets: make(map[string]*Wallet),
		ledger:  []LedgerEntry{},
	}
}

// CreateWallet initializes a new wallet for a user
func (s *WalletService) CreateWallet(userID string, initialBalance float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wallets[userID] = &Wallet{UserID: userID, Balance: initialBalance}
}

// GetBalance returns the current balance of a wallet
func (s *WalletService) GetBalance(userID string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, ok := s.wallets[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}
	return wallet.Balance, nil
}

// DeductFunds removes funds from a wallet atomically
func (s *WalletService) DeductFunds(userID string, amount float64, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, ok := s.wallets[userID]
	if !ok {
		return ErrWalletNotFound
	}

	if wallet.Balance < amount {
		return ErrInsufficientFunds
	}

	wallet.Balance -= amount
	s.addLedgerEntry(userID, amount, DEBIT, description)
	return nil
}

// AddFunds adds funds to a wallet atomically
func (s *WalletService) AddFunds(userID string, amount float64, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, ok := s.wallets[userID]
	if !ok {
		return ErrWalletNotFound
	}

	wallet.Balance += amount
	s.addLedgerEntry(userID, amount, CREDIT, description)
	return nil
}

func (s *WalletService) addLedgerEntry(userID string, amount float64, tType TransactionType, desc string) {
	entry := LedgerEntry{
		ID:          fmt.Sprintf("tx-%d", len(s.ledger)+1),
		WalletID:    userID,
		Amount:      amount,
		Type:        tType,
		Description: desc,
		CreatedAt:   time.Now(),
	}
	s.ledger = append(s.ledger, entry)
}

func main() {
	fmt.Println("Wallet Service: Initialized")
	ws := NewWalletService()

	// Initial setup
	ws.CreateWallet("user1", 1000.0)
	
	bal, _ := ws.GetBalance("user1")
	fmt.Printf("Initial Balance for user1: %.2f\n", bal)

	err := ws.DeductFunds("user1", 100.0, "Order #123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	newBal, _ := ws.GetBalance("user1")
	fmt.Printf("Balance after deduction: %.2f\n", newBal)

	ws.AddFunds("user1", 50.0, "Refund #123")
	finalBal, _ := ws.GetBalance("user1")
	fmt.Printf("Final Balance for user1: %.2f\n", finalBal)
}
