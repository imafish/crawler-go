package main

// GroupContext represents a group of tasks that has similar properties.
type GroupContext struct {
	groupBy string

	name    string
	counter *counter
	dir     string
}

// IsMatch returns if given task context is a match to a group context
func (g GroupContext) IsMatch(taskContext *TaskContext) bool {
	match := g.name == FormatString(g.groupBy, taskContext)
	return match
}
