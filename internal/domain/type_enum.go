package domain

// Transaction is an enumeration of transaction types
type TransactionType int64

const (
	// TransactionTypeUnSpecified is an enumeration of unspecified transaction type
	TransactionTypeUnSpecified TransactionType = iota
	// TransactionTypeIncome is an enumeration of income transaction type
	TransactionTypeIncome
	// TransactionTypeExpense is an enumeration of expense transaction type
	TransactionTypeExpense
	// TransactionTypeBoth is an enumeration of both income and expense transaction type
	// it's only used in monthly data
	TransactionTypeBoth
)

// ToString returns the string representation of TransactionType
func (t TransactionType) ToString() string {
	switch t {
	case TransactionTypeIncome:
		return "income"
	case TransactionTypeExpense:
		return "expense"
	case TransactionTypeBoth:
		return "both"
	}
	return "unknown type"
}

// ToModelValue returns the string enum of mysql
func (t TransactionType) ToModelValue() string {
	switch t {
	case TransactionTypeIncome:
		return "1"
	case TransactionTypeExpense:
		return "2"
	}
	return "0"
}

// IsValid checks if the TransactionType is valid
func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeIncome, TransactionTypeExpense:
		return true
	}
	return false
}

// CvtToTransactionType converts string to TransactionType
func CvtToTransactionType(s string) TransactionType {
	switch s {
	case "income":
		return TransactionTypeIncome
	case "expense":
		return TransactionTypeExpense
	case "both":
		return TransactionTypeBoth
	case "1":
		return TransactionTypeIncome
	case "2":
		return TransactionTypeExpense
	}
	return TransactionTypeUnSpecified
}
