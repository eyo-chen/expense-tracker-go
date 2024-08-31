package transaction

import (
	"fmt"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/model/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/model/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/model/user"
)

var (
	mockLoc, _  = time.LoadLocation("")
	mockTimeNow = time.Unix(1629446406, 0).Truncate(24 * time.Hour).In(mockLoc)
)

func BluePrint(i int) Transaction {
	return Transaction{
		Type:  domain.TransactionTypeIncome.ToModelValue(),
		Price: float64(i*10.0 + 1.0),
		Note:  "test" + fmt.Sprint(i),
		Date:  mockTimeNow,
	}
}

// GetAll_GenExpResult generates expected transactions

func GetAll_GenExpResult(ts []Transaction, u user.User, ms []maincateg.MainCateg, ss []subcateg.SubCateg, is []icon.Icon, indexList ...int) []domain.Transaction {
	expResult := make([]domain.Transaction, 0, len(indexList))
	for _, i := range indexList {
		expResult = append(expResult, domain.Transaction{
			ID:     ts[i].ID,
			UserID: u.ID,
			Type:   domain.CvtToTransactionType(ts[i].Type),
			MainCateg: domain.MainCateg{
				ID:   ms[i].ID,
				Name: ms[i].Name,
				Type: domain.CvtToTransactionType(ms[i].Type),
				Icon: domain.Icon{
					ID:  is[i].ID,
					URL: is[i].URL,
				},
			},
			SubCateg: domain.SubCateg{
				ID:          ss[i].ID,
				Name:        ss[i].Name,
				MainCategID: ss[i].MainCategID,
			},
			Price: ts[i].Price,
			Note:  ts[i].Note,
			Date:  ts[i].Date,
		})
	}

	return expResult
}
