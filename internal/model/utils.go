package model

func cvtToModelType(t string) string {
	if t == "income" {
		return "1"
	}
	return "2"
}

func cvtToDomainType(t string) string {
	if t == "1" {
		return "income"
	}
	return "expense"
}
