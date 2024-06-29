package maincateg

import (
	"fmt"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func BluePrint(i int, last MainCateg) MainCateg {
	return MainCateg{
		Name: "test" + fmt.Sprint(i),
		Type: domain.TransactionTypeIncome.ToModelValue(),
	}
}
