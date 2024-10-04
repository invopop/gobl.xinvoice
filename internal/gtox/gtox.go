package xinvoice

// Date defines date in the UDT structure
type Date struct {
	Date   string `xml:",chardata"`
	Format string `xml:"format,attr,omitempty"`
}
