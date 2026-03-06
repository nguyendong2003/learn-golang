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
