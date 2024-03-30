package transaction

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

// GetMonthlyData_GenExpResult generates the result of the monthly data
func GetMonthlyData_GenExpResult(monthlyData domain.MonthDayToTransactionType, maxDay int, err error) []domain.TransactionType {
	result := make([]domain.TransactionType, 0, maxDay)

	if err != nil {
		return result
	}

	// loop from 1 to the last day of the month(30 or 31)
	for i := 1; i <= maxDay; i++ {
		if t, ok := monthlyData[i]; ok {
			result = append(result, t)
		} else {
			result = append(result, domain.TransactionTypeUnSpecified)
		}
	}

	return result
}
