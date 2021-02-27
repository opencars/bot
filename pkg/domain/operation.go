package domain

// Operation represents public registrations of transport.
type Operation struct {
	Person      string
	RegAddress  *string
	RegCode     int16
	Reg         string
	Date        string
	DepCode     int32
	Dep         string
	Brand       string
	Model       string
	Year        int16
	Color       string
	Kind        string
	Body        string
	Purpose     string
	Fuel        *string
	Capacity    *int
	OwnWeight   *float64
	TotalWeight *float64
	Number      string
}
