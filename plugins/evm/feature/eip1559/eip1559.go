package eip1559

import (
	"github.com/shopspring/decimal"
)

// https://eips.ethereum.org/EIPS/eip-1559

// total fee: gas_usage * <base gas price + priority gas price>
// Burnt: 80% * total fee
// Treasury Fees: total fee - Burnt
// Savings: (max_price - base_price - priority price) * gas_usage

// func TreasuryFees(totalFee decimal.Decimal) decimal.Decimal {
// 	return totalFee.Mul(decimal.NewFromFloat(0.2))
// }

func BurntFee(maxPriorityFeePerGas, totalFee decimal.Decimal) decimal.Decimal {
	if maxPriorityFeePerGas.IsZero() {
		return decimal.Zero
	}
	return totalFee.Mul(decimal.NewFromFloat(0.8))
}

func SavingsFee(maxPrice, basePrice, priorityPrice, gasUsage decimal.Decimal) decimal.Decimal {
	return maxPrice.Sub(basePrice).Sub(priorityPrice).Mul(gasUsage)
}
