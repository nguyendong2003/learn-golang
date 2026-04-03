package util

import "fmt"

func AddILIKECondition(
	conditions *[]string,
	args *[]any,
	column string,
	value *string,
) {
	if value != nil && *value != "" {
		*conditions = append(*conditions, fmt.Sprintf("%s ILIKE ?", column))
		*args = append(*args, "%"+*value+"%")
	}
}

func AddEqualCondition(
	conditions *[]string,
	args *[]any,
	column string,
	value *string,
) {
	if value != nil && *value != "" {
		*conditions = append(*conditions, fmt.Sprintf("%s = ?", column))
		*args = append(*args, *value)
	}
}

func AddEqualBoolCondition(
	conditions *[]string,
	args *[]any,
	column string,
	value *bool,
) {
	if value != nil {
		*conditions = append(*conditions, fmt.Sprintf("%s = ?", column))
		*args = append(*args, *value)
	}
}
