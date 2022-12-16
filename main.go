package main

import (
	"bot/models"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

// var dev = false
// var token = "eyJhbGciOiJIUzI1NiJ9.dNt1CWz3chfApmj5VW7EtnA054wLHUxPbm-4JmULr4TPxjlpbka9shM4msYYj2UF9WTnFV_ZJcjZxKboWAzbLG8ZXP_boyD2CMapkUV645ee6rvYEjCVhEBsRFF3A9qLPfCZqegdKqy5rDGIz4dzNaOqEOuVUpgyQAh8GGTnUJY.OPVfo0Tb8rH8SrPNWjaEIRI9z87xRuaQBgO-lEKT9VQ"

// ##### BABY Lovet ######
// var shopName = "@babylovett"
// var shopId = "16992"
// var shopShipmentChannelId = "30507"
// var targets = []string{"3t", "3-t", "2-3t", "3-4t"}
// var lowestProductId = 1003472347
//##########

// ##### TILLY ######
var shopName = "@tillymilly"
var shopId = "233173"
var shopShipmentChannelId = "302545"
var lowestProductId = 1003355532
var targets = []string{"3t", "3-t", "2-3t", "3-4t"}

//##########

func main() {
	token := flag.String("token", "", "")
	dev := flag.Bool("dev", false, "")
	flag.Parse()

	log("")
	log("")
	log("------------------------------------")
	log("dev mode: ", *dev)
	log("token: ", *token)
	log("shop: ", shopName, " id:", shopId)
	log("lowestProductId: ", 1003411793)
	log("------------------------------------")
	log("")
	log("")

	getProducts(*token, *dev)

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func getProducts(token string, dev bool) {
	url := "https://customer-api.line-apps.com/search/graph"
	method := "POST"

	reqBody := fmt.Sprintf(`{
		"operationName": "shopProductListQuery",
		"variables": {
			"limit": 100,
			"page": 1,
			"shopId": %s,
			"sortType": "CREATED_TIME",
			"sortSoldOutType": "NONE"
		},
		"query": "query shopProductListQuery($limit: Int!, $page: Int!, $shopId: Int!, $sortType: SortType!, $sortSoldOutType: SortSoldOutType!) {\n  shopProductList(\n    limit: $limit\n    page: $page\n    shopId: $shopId\n    sortType: $sortType\n    sortSoldOutType: $sortSoldOutType\n  ) {\n    totalPage\n    totalProduct\n    products {\n      id\n      productName: name\n      imgUrl: imageUrl\n      instantDiscount: discountPercent\n      price\n      discountedPrice: salePrice\n      isInStock\n      __typename\n    }\n    __typename\n  }\n}\n"
	}`, shopId)
	payload := strings.NewReader(reqBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Host", "customer-api.line-apps.com")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	req.Header.Add("accept", "*/*")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("origin", "https://shop.line.me")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://shop.line.me/")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th;q=0.8")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var p *models.ProductList
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&p); err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("total products: ", p.Data.ShopProductList.TotalProduct)

	hasNewItem := false
	for _, v := range p.Data.ShopProductList.Products {
		// fmt.Println("product-id: ", v.ID, "product-name: ", v.ProductName)
		pId, err := strconv.Atoi(v.ID)
		if err != nil {
			fmt.Println("err: ", err)
		} else {
			if pId > lowestProductId {
				hasNewItem = true
				fmt.Println("🎉 new id:", pId, " name: ", v.ProductName)
				go getProductDetail(v.ID, token, dev)
			}
		}
	}
	if !hasNewItem {
		fmt.Println("🥶 waiting for new product....")
		getProducts(token, dev)
	}
}

func getProductDetail(productId string, token string, dev bool) {
	fmt.Println("--------------------------")
	fmt.Println("----Product Detail----")
	fmt.Println("--------------------------")
	fmt.Println("product-id: ", productId)
	url := fmt.Sprintf("https://sc-oms-api.line-apps.com/api/v1/shopend/%s/product/%s", shopName, productId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Host", "sc-oms-api.line-apps.com")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("authorization", token)
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("origin", "https://shop.line.me")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://shop.line.me/")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Cookie", "JSESSIONID=9D35C4CA951F924DE3E2561112604E2C")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var p *models.Product
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&p); err != nil {
		fmt.Println(err)
		return
	}
	productName := p.Data.ProductName
	pId := p.Data.ID
	fmt.Println("--------> Product detail RESPONSE")
	fmt.Println("--------> ")
	fmt.Println("product-id: ", pId)
	fmt.Println("product: ", productName)
	fmt.Println("has-varient: ", p.Data.HasVariant)
	if p.Data.HasVariant {
		fmt.Println("##varients##")
		for _, v := range p.Data.ProductVariants {
			fmt.Println("varient-id: ", v.ID)
			fmt.Println("available: ", v.Available)
			fmt.Println("option1: ", v.VariantOptionValue1)
			fmt.Println("option2: ", v.VariantOptionValue2)
			fmt.Println("------------")
			sort.Strings(targets)
			lowerCase := strings.ReplaceAll(strings.ToLower(v.VariantOptionValue1), " ", "")
			if contains(targets, lowerCase) {
				if v.Available > 0 {
					fmt.Println("🎯 meet target: ", v.VariantOptionValue1, " available: ", v.Available)
					go placeOrder(productName, pId, strconv.Itoa(v.ID), token, dev)
				} else {
					fmt.Println("😩 meet target: ", v.VariantOptionValue1, " but not available")
				}
			} else {
				fmt.Println("🥵 non target")
			}
		}
		fmt.Println("############")
	} else {
		if p.Data.Available > 0 {
			go placeOrder(productName, pId, "null", token, dev)
		}
	}
}

func placeOrder(name string, productId int, varientId string, token string, dev bool) {
	fmt.Println("🚕-----Place Order----")
	fmt.Println("productId: ", productId, " - productName: ", name, " - varientId: ", varientId)

	if dev {
		fmt.Println("🚕Place order in development")
		return
	}
	url := fmt.Sprintf("https://sc-oms-api.line-apps.com/api/v5/shopend/%s/order/place-order-bank", shopName)
	method := "POST"

	reqBody := fmt.Sprintf(`{
		"customer_address_id": "2816333",
		"shop_shipment_channel_id": %s,
		"shipping_address": {
			"id": 2816333,
			"recipient_name": "กชนิภา อิสรั่น",
			"address": "12 ซอยประชาร่วมใจ11. แขวงทรายกองดินใต้ เขตคลองสามวา กรุงเทพฯ",
			"province": "กรุงเทพ",
			"postal_code": "10510",
			"country": "TH",
			"country_id": 1,
			"phone_number": "0990948476",
			"email": "rnut.ist@gmail.com"
		},
		"is_android": false,
		"is_liff": false,
		"remark_buyer": "",
		"point": null,
		"ignore_point": null,
		"ignore_l_m_coupon": null,
		"promotion_id": "",
		"items": [
			{
				"product_id": %d,
				"product_variant_id": %s,
				"is_promotion": false,
				"quantity": 1
			}
			],
			"cart_id": "",
			"coupon_id": null
			}`, shopShipmentChannelId, productId, varientId)

	fmt.Println("reqBody: ", reqBody)

	payload := strings.NewReader(reqBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("authority", "sc-oms-api.line-apps.com")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en")
	req.Header.Add("authorization", token)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://shop.line.me")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("referer", "https://shop.line.me/")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "JSESSIONID=45B2AD99BECB83DD6DDA0E40C66FDD11")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("--------------------------")
	fmt.Println("PLACE ORDER RESPONSE")
	fmt.Println(string(body))
	fmt.Println("--------------------------")
}

func contains(s []string, searchterm string) bool {
	for _, v := range s {
		if strings.Contains(searchterm, v) {
			return true
		}
	}
	return false
}

func log(a ...interface{}) {
	fmt.Println(a...)
}
