package driver

import (
	"fmt"
	"strings"
)

// Filter is a filter for a query.
type Filter struct {
	Condition Condition
	Key       string
	Value     string
}

// GetCondition returns the condition of the filter.
func (f *Filter) GetCondition() Condition {
	return f.Condition
}

// GetKey returns the key of the filter.
func (f *Filter) GetKey() string {
	return f.Key
}

// GetValue returns the value of the filter.
func (f *Filter) GetValue() string {
	return f.Value
}

// FilterableConnection is a connection that can be filtered.
type FilterableConnection interface {
	Where(key string, value string) FilterableConnection
	WhereCondition(condition Condition, key string, value string) FilterableConnection
	GetFilters() []*Filter
	ResetFilters()
}

// Condition is a connection to a database.
type Condition uint8

const (
	// Eq is the equal condition.
	Eq Condition = iota
	// Ne is the not equal condition.
	Ne
	// Gt is the greater than condition.
	Gt
	// Lt is the less than condition.
	Lt
	// Ge is the greater than or equal condition.
	Ge
	// Le is the less than or equal condition.
	Le
)

var operators = map[Condition]string{
	Eq: "=",
	Ne: "!=",
	Gt: ">",
	Lt: "<",
	Ge: ">=",
	Le: "<=",
}

func (c Condition) String() string {
	return operators[c]
}

func GetConditionFromString(condition string) Condition {
	for k, v := range operators {
		if v == condition {
			return k
		}
	}
	return Eq
}

type DefaultFilterBuilder struct {
	filters []*Filter
}

func (fb *DefaultFilterBuilder) Where(key string, value string) FilterableConnection {
	return fb.WhereCondition(Eq, key, value)
}

func (fb *DefaultFilterBuilder) WhereCondition(condition Condition, key string, value string) FilterableConnection {
	fb.filters = append(fb.filters, &Filter{
		Condition: condition,
		Key:       key,
		Value:     value,
	})
	return fb
}

func (fb *DefaultFilterBuilder) GetFilters() []*Filter {
	return fb.filters
}

func (fb *DefaultFilterBuilder) ResetFilters() {
	fb.filters = []*Filter{}
}

func (fb *DefaultFilterBuilder) BuildFilterSQL() string {
	if len(fb.filters) == 0 {
		return ""
	}
	var sql strings.Builder
	sql.WriteString(" WHERE ")
	for i, filter := range fb.filters {
		if i > 0 {
			sql.WriteString(" AND ")
		}
		sql.WriteString(fmt.Sprintf("%s %s %s", filter.Key, filter.Condition, filter.Value))
	}
	return sql.String()
}
