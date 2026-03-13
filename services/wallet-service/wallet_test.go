package main

import (
	"sync"
	"testing"
)

func TestWalletService(t *testing.T) {
	ws := NewWalletService()
	userID := "test-user"
	initialBal := 500.0

	ws.CreateWallet(userID, initialBal)

	t.Run("Check Balance", func(t *testing.T) {
		bal, err := ws.GetBalance(userID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if bal != initialBal {
			t.Errorf("Expected balance %.2f, got %.2f", initialBal, bal)
		}
	})

	t.Run("Deduct Funds Success", func(t *testing.T) {
		err := ws.DeductFunds(userID, 100.0, "test deduction")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		bal, _ := ws.GetBalance(userID)
		if bal != 400.0 {
			t.Errorf("Expected balance 400.0, got %.2f", bal)
		}
	})

	t.Run("Deduct Funds Insufficient", func(t *testing.T) {
		err := ws.DeductFunds(userID, 1000.0, "overdraft")
		if err != ErrInsufficientFunds {
			t.Errorf("Expected ErrInsufficientFunds, got %v", err)
		}
	})

	t.Run("Concurrent Add/Deduct", func(t *testing.T) {
		var wg sync.WaitGroup
		iterations := 100
		amount := 1.0

		for i := 0; i < iterations; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				ws.AddFunds(userID, amount, "concurrent add")
			}()
			go func() {
				defer wg.Done()
				ws.DeductFunds(userID, amount, "concurrent deduct")
			}()
		}
		wg.Wait()

		bal, _ := ws.GetBalance(userID)
		if bal != 400.0 {
			t.Errorf("Expected balance 400.0 after concurrent ops, got %.2f", bal)
		}
	})
}
