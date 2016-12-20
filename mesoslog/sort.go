package mesoslog

type SortTasksByLatestTimestamp []*mstateTask

func (s SortTasksByLatestTimestamp) Len() int {
	return len(s)
}
func (s SortTasksByLatestTimestamp) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortTasksByLatestTimestamp) Less(i, j int) bool {
	return s[j].LastTimestamp < s[i].LastTimestamp
}
