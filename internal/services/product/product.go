package product

type Product struct {
	ProductID       int     `json:"productId"`
	PricePerUnit    float64 `json:"pricePerUnit"`
	ProductName     string  `json:"productName"`
	ProductBrand    string  `json:"productBrand"`
	Description     string  `json:"description"`
	StockQuantity   int     `json:"stockQuantity"`
	PricePerUnitUSD float64 `json:"price_usd"`
	PricePerUnitEUR float64 `json:"price_eur"`
	PricePerUnitGBP float64 `json:"price_gbp"`
}
