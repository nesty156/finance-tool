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

func PortfolioValue(portfolio Portfolio) float64 {
	total := 0.0

	// Loop through products and add up the amounts
	for _, product := range portfolio.Products {
		total += product.ValueEUR
	}

	return total
}

func MergePortfolios(destination Portfolio, source Portfolio) Portfolio {
	for _, sourceProduct := range source.Products {
		productIndex := -1
		for i, destinationProduct := range destination.Products {
			if destinationProduct.SymbolISIN == sourceProduct.SymbolISIN {
				productIndex = i
				break
			}
		}

		if productIndex == -1 {
			newProduct := Product{Name: sourceProduct.Name, SymbolISIN: sourceProduct.SymbolISIN, Quantity: sourceProduct.Quantity, ValueEUR: sourceProduct.ValueEUR}
			destination.Products = append(destination.Products, newProduct)
		} else {
			destination.Products[productIndex].Quantity += sourceProduct.Quantity
			destination.Products[productIndex].ValueEUR += sourceProduct.ValueEUR
		}
	}

	return destination
}
