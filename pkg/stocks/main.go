package stocks

type Portfolio struct {
	Name     string
	Products []Product
}

type Product struct {
	Name       string
	SymbolISIN string
	Quantity   float64
	ValueEUR   float64
}
