package validator

// GetPutObjectURL validates the file name.
func (v *Validator) GetPutObjectURL(fileName string) bool {
	v.Check(len(fileName) > 0, "file_name", "File name can't be empty")
	return v.Valid()
}
