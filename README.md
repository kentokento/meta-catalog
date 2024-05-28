# meta-catalog
meta product catalog batch for Go


### Official documentations
- https://developers.facebook.com/docs/marketing-api/catalog-batch/guides/send-product-updates
- https://developers.facebook.com/docs/marketing-api/catalog-batch/reference/#product-item

### How to use

```go
list := []catalog.Request{
	catalog.Request{Method: "UPDATE", Data: catalog.Product{
		Id:           "1234",
		Availability: "in stock",
		Condition:    "new",
		Description:  "description",
		Link:         "https://link",
		Price:        "2000 JPY",
		Title:        "product title",
	}},
	catalog.Request{Method: "DELETE", Data: catalog.Product{Id: "1235"}},
}

const catalogID = 1234567
const apiVersion = "v20.0"
cli := catalog.NewClient(apiVersion, catalogID)
const accessToken = "Please generate an access token with scope of 「catalog_management」"
if err := cli.SendUpsert(list, accessToken); err != nil {
    return err
}
```
