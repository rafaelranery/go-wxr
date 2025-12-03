package wxr

// Filter defines the interface for filtering WXR items.
// Implementations determine whether an item should be included in the parsed results.
type Filter interface {
	ShouldInclude(item *item) bool
}

// DefaultFilter implements the default filtering strategy.
// It includes only items with post_type="post" and status="publish".
type DefaultFilter struct {
	PostType string
	Status   string
}

// ShouldInclude returns true if the item matches the filter criteria.
func (f *DefaultFilter) ShouldInclude(item *item) bool {
	if f.PostType != "" && item.PostType != f.PostType {
		return false
	}
	if f.Status != "" && item.Status != f.Status {
		return false
	}
	return true
}

// NewDefaultFilter creates a new DefaultFilter with standard WordPress post filtering.
func NewDefaultFilter() *DefaultFilter {
	return &DefaultFilter{
		PostType: "post",
		Status:   "publish",
	}
}
