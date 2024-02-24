package domain

// Transaction is an enumeration of transaction types
type TransactionType int64

const (
	// UnSpecified is an enumeration of unspecified transaction type
	UnSpecified TransactionType = iota
	// Income is an enumeration of income transaction type
	Income
	// Expense is an enumeration of expense transaction type
	Expense
)

// ToString returns the string representation of TransactionType
func (t TransactionType) ToString() string {
	switch t {
	case Income:
		return "income"
	case Expense:
		return "expense"
	}
	return "unknown type"
}

// ToModelValue returns the string enum of mysql
func (t TransactionType) ToModelValue() string {
	switch t {
	case Income:
		return "1"
	case Expense:
		return "2"
	}
	return "0"
}

// IsValid checks if the TransactionType is valid
func (t TransactionType) IsValid() bool {
	switch t {
	case Income, Expense:
		return true
	}
	return false
}

// CvtToTransactionType converts string to TransactionType
func CvtToTransactionType(s string) TransactionType {
	switch s {
	case "income":
		return Income
	case "expense":
		return Expense
	case "1":
		return Income
	case "2":
		return Expense
	}
	return UnSpecified
}
