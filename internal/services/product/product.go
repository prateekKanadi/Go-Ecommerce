package product

type Product struct {
	ProductID     int            `json:"productId"`
	PricePerUnit  float64        `json:"pricePerUnit"`
	ProductName   string         `json:"productName"`
	ProductBrand  string         `json:"productBrand"`
	Description   string         `json:"description"`
	StockQuantity int            `json:"stockQuantity"`
	Category      string         `json:"category"`
	SubCategory   string         `json:"subCategory"`
	ImageURL      string         `json:"imageURL"`
	Prices        []ProductPrice `json:"prices"`
}

type ProductPrice struct {
	CurrencyCode string  `json:"currencyCode"`
	Amount       float64 `json:"price"`
}

type VariantProduct struct {
	VariantID   string `json:"variantId"`
	ProductID   int    `json:"productId"`
	VariantName string `json:"variantName"`
	Color       string `json:"color"`
	Product
}
