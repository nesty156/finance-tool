package stocks

import (
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

type EURAmount struct {
	float64
}

type DegiroProduct struct {
	Name       string    `csv:"Produkt"`
	SymbolISIN string    `csv:"Symbol/ISIN"`
	Quantity   float64   `csv:"Množství"`
	ValueEUR   EURAmount `csv:"Hodnota v EUR"`
}

// Convert the CSV string to internal float64
func (f *EURAmount) UnmarshalCSV(csv string) (err error) {
	csv = strings.ReplaceAll(csv, " ", "")
	csv = strings.ReplaceAll(csv, ",", ".")
	f.float64, err = strconv.ParseFloat(csv, 64)
	return err
}

func CreateDegiroPortfolio(fileName string, portfolioName string) (Portfolio, error) {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	degiroProducts := []*DegiroProduct{}

	if err := gocsv.UnmarshalFile(csvFile, &degiroProducts); err != nil { // Load clients from file
		panic(err)
	}

	// Convert to internal format
	products := []Product{}
	for _, product := range degiroProducts {
		products = append(products, ConvertToProduct(*product))
	}

	portfolio := Portfolio{Name: portfolioName, Products: products}

	return portfolio, nil
}

func ConvertToProduct(degiroProduct DegiroProduct) Product {
	product := Product{
		Name:       degiroProduct.Name,
		SymbolISIN: degiroProduct.SymbolISIN,
		Quantity:   degiroProduct.Quantity,
		ValueEUR:   degiroProduct.ValueEUR.float64,
	}

	return product
}
