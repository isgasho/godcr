package walletcore

import (
	"fmt"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"math/rand"
	"strconv"
	"strings"
)

type SyncStatus uint8

const (
	SyncStatusNotStarted SyncStatus = iota
	SyncStatusSuccess
	SyncStatusError
	SyncStatusInProgress
)

// FormatTxType returns a more readable representation of a transaction type,
// returning Regular instead of REGULAR. Ticket Purchase instead of TICKET_PURCHASE, etc
func FormatTxType(txType string) string {
	txType = strings.Replace(txType, "_", " ", -1)
	txType = strings.ToLower(txType)
	return strings.Title(txType)
}

func SimpleBalance(balance *Balance, detailed bool) string {
	if detailed {
		return balance.Total.String()
	} else {
		return balance.String()
	}
}

// NormalizeBalance adds 0s the right of balance to make it x.xxxxxxxx DCR
func NormalizeBalance(balance float64) string {
	return fmt.Sprintf("%010.8f DCR", balance)
}

// GetChangeDestinationsWithRandomAmounts generates change destination(s) based on the number of change address the user want
func GetChangeDestinationsWithRandomAmounts(wallet Wallet, nChangeOutputs int, amountInAtom int64, sourceAccount uint32,
	nUtxoSelection int, sendDestinations []txhelper.TransactionDestination) (changeOutputDestinations []txhelper.TransactionDestination, err error) {

	var changeAddresses []string
	for i := 0; i < nChangeOutputs; i++ {
		address, err := wallet.GenerateNewAddress(sourceAccount)
		if err != nil {
			return nil, fmt.Errorf("error generating address: %s", err.Error())
		}
		changeAddresses = append(changeAddresses, address)
	}

	changeAmount, err := txhelper.EstimateChange(nUtxoSelection, amountInAtom, sendDestinations, changeAddresses)
	if err != nil {
		return nil, fmt.Errorf("error in getting change amount: %s", err.Error())
	}
	if changeAmount <= 0 {
		return
	}

	var portionRations []float64
	var rationSum float64
	for i := 0; i < nChangeOutputs; i++ {
		portion := rand.Float64()
		portionRations = append(portionRations, portion)
		rationSum += portion
	}

	for i, portion := range portionRations {
		portionPercentage := portion / rationSum
		amount := portionPercentage * float64(changeAmount)

		changeOutput := txhelper.TransactionDestination{
			Address: changeAddresses[i],
			Amount:  dcrutil.Amount(amount).ToCoin(),
		}
		changeOutputDestinations = append(changeOutputDestinations, changeOutput)
	}
	return
}

func BuildTxDestinations(destinationAddresses []string, destinationAmounts []string) (destinations []txhelper.TransactionDestination, err error) {
	destinations = make([]txhelper.TransactionDestination, len(destinationAddresses))
	for i := range destinationAddresses {
		amount, err := strconv.ParseFloat(destinationAmounts[i], 64)
		if err != nil {
			return destinations, err
		}
		destinations[i] = txhelper.TransactionDestination{
			Address: destinationAddresses[i],
			Amount:  amount,
		}
	}
	return
}
