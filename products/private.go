package products

func (product *Product) GenerateJson() ProductJSON {

	return ProductJSON{
		Name:        product.Name,
		Brand:       product.Brand,
		Price:       product.Price,
		Barcode:     product.Barcode,
		ImageUri:    product.ImageUri,
		Description: product.Description,
	}
}