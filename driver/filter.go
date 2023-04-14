package driver

import (
	"fmt"
	"strings"
)

// Filter is a filter for a query.
type Filter struct {
	Key       string
	Value     string
	Condition Condition
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
	Where(dataCollection, key string, value string) FilterableConnection
	WhereRoot(key, value string) FilterableConnection
	WhereCondition(dataCollection string, condition Condition, key string, value string) FilterableConnection
	WhereRootCondition(condition Condition, key string, value string) FilterableConnection
	GetFilters(dataCollection string) []*Filter
	GetAllFilters() map[string][]*Filter
	GetRootFilters() []*Filter
	ResetFilters(dataCollection string)
	ResetAllFilters()
	ResetRootFilters()
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

// DefaultFilterBuilder is a default implementation of FilterableConnection.
// That can be used as embedded struct in a driver to implement simple filterable connection.
// And also this builder can build sql query from the filters.
type DefaultFilterBuilder struct {
	filters     map[string][]*Filter
	rootFilters []*Filter // root filters is a general filter that will apply to all data collections.
}

func (fb *DefaultFilterBuilder) Where(dataCollection, key string, value string) FilterableConnection {
	return fb.WhereCondition(dataCollection, Eq, key, value)
}

// WhereRoot is a general filter that is not related to any data collection. and will apply to all data collections.
func (fb *DefaultFilterBuilder) WhereRoot(key string, value string) FilterableConnection {
	return fb.WhereRootCondition(Eq, key, value)
}

func (fb *DefaultFilterBuilder) WhereCondition(dataCollection string, condition Condition, key string, value string) FilterableConnection {
	fb.initFiltersIfIsNil() // allocate memory for the filters map if it is nil, for preventing panic.

	fb.filters[dataCollection] = append(fb.filters[dataCollection], &Filter{
		Condition: condition,
		Key:       key,
		Value:     value,
	})

	return fb
}

// WhereRootCondition is a general filter that is not related to any data collection. and will apply to all data collections.
func (fb *DefaultFilterBuilder) WhereRootCondition(condition Condition, key, value string) FilterableConnection {
	fb.initFiltersIfIsNil() // allocate memory for the filters map if it is nil, for preventing panic.

	fb.rootFilters = append(fb.rootFilters, &Filter{
		Condition: condition,
		Key:       key,
		Value:     value,
	})

	return fb
}

func (fb *DefaultFilterBuilder) GetFilters(dataCollection string) []*Filter {
	fb.initFiltersIfIsNil() // allocate memory for the filters map if it is nil, for preventing panic.
	return append(fb.rootFilters, fb.filters[dataCollection]...)
}

// GetAllFilters returns all data collections filters without applying root filters.
func (fb *DefaultFilterBuilder) GetAllFilters() map[string][]*Filter {
	return fb.filters
}

// GetRootFilters returns root filters.
func (fb *DefaultFilterBuilder) GetRootFilters() []*Filter {
	return fb.rootFilters
}

func (fb *DefaultFilterBuilder) ResetFilters(dataCollection string) {
	fb.initFiltersIfIsNil() // allocate memory for the filters map if it is nil, for preventing panic.
	fb.filters[dataCollection] = []*Filter{}
}
func (fb *DefaultFilterBuilder) ResetRootFilters() {
	fb.rootFilters = []*Filter{}
}

func (fb *DefaultFilterBuilder) ResetAllFilters() {
	fb.filters = make(map[string][]*Filter)
	fb.rootFilters = []*Filter{}
}

func (fb *DefaultFilterBuilder) BuildFilterSQL(dataCollection string) string {
	fb.initFiltersIfIsNil() // allocate memory for the filters map if it is nil, for preventing panic.

	if len(fb.rootFilters) == 0 && len(fb.filters[dataCollection]) == 0 {
		return ""
	}
	var sql strings.Builder
	sql.WriteString(" WHERE ")
	for i, filter := range append(fb.rootFilters, fb.filters[dataCollection]...) {
		if i > 0 {
			sql.WriteString(" AND ")
		}
		sql.WriteString(fmt.Sprintf("%s %s %s", filter.Key, filter.Condition, filter.Value))
	}
	return sql.String()
}

func (fb *DefaultFilterBuilder) initFiltersIfIsNil() {
	if fb.filters == nil {
		fb.filters = make(map[string][]*Filter)
	}
}
