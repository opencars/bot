package gov

// Sorter sorts array by LastModified field.
type Sorter []Resource

func (e Sorter) Len() int      { return len(e) }
func (e Sorter) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Sorter) Less(i, j int) bool {
	return e[i].LastModified.Before(e[j].LastModified.Time)
}
