package models

type Product struct {
	Data struct {
		ID                 int    `json:"id"`
		ProductName        string `json:"product_name"`
		HasVariant         bool   `json:"has_variant"`
		VariantOptionName1 string `json:"variant_option_name1"`
		Available          int    `json:"available"`
		ProductVariants    []struct {
			ID                  int    `json:"id"`
			VariantOptionValue1 string `json:"variant_option_value1"`
			VariantOptionValue2 string `json:"variant_option_value2"`
			Available           int    `json:"available"`
			ProductSnapshotID   int    `json:"product_snapshot_id"`
		} `json:"product_variants"`
	} `json:"data"`
}

type ProductList struct {
	Data struct {
		ShopProductList struct {
			TotalPage    int `json:"totalPage"`
			TotalProduct int `json:"totalProduct"`
			Products     []struct {
				ID              string `json:"id"`
				ProductName     string `json:"productName"`
				ImgURL          string `json:"imgUrl"`
				InstantDiscount int    `json:"instantDiscount"`
				Price           string `json:"price"`
				DiscountedPrice string `json:"discountedPrice"`
				IsInStock       bool   `json:"isInStock"`
				Typename        string `json:"__typename"`
			} `json:"products"`
		} `json:"shopProductList"`
	} `json:"data"`
}
