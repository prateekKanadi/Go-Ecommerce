package product

type Product struct {
	ProductID    int    `json:"productId"`
	PricePerUnit string `json:"pricePerUnit"`
	ProductName  string `json:"productName"`
	ProductBrand string `json:"productBrand"`
}
