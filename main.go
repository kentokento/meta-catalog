package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type client struct {
	ctx        context.Context
	url        string
	httpClient *http.Client
}

func NewClient(apiVersion string, catalogID uint64) *client {
	return &client{
		url: fmt.Sprintf("https://graph.facebook.com/%s/%d/items_batch", apiVersion, catalogID),
	}
}

func (c *client) SetContext(ctx context.Context) *client {
	c.ctx = ctx
	return c
}

func (c *client) SetHttpClient(cli *http.Client) *client {
	c.httpClient = cli
	return c
}

func (c *client) SendUpsert(requests []Request, token string) error {
	req := newBatchRequest(requests, token)
	req.AllowUpsert = true
	return c.send(req)
}

func (c *client) Send(requests []Request, token string) error {
	req := newBatchRequest(requests, token)
	return c.send(req)
}

func (c *client) send(batchRequest BatchRequest) error {
	ctx := c.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	httpClient := c.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	requestBody, err := json.Marshal(batchRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	v := Response{}
	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}
	if len(v.Handles) == 0 {
		er := Error{}
		if err := json.Unmarshal(body, &er); err != nil {
			return err
		}
		return er
	}
	for _, vv := range v.ValidationStatus {
		for _, e := range vv.Errors {
			if err != nil {
				err = fmt.Errorf("%w, %w", err, e)
			} else {
				err = e
			}
		}
	}
	return err
}

func newBatchRequest(requests []Request, token string) BatchRequest {
	return BatchRequest{
		AccessToken: token,
		ItemType:    "PRODUCT_ITEM",
		Requests:    requests,
	}
}

type BatchRequest struct {
	AccessToken string    `json:"access_token"`
	ItemType    string    `json:"item_type"`
	Requests    []Request `json:"requests"`

	AllowUpsert bool `json:"allow_upsert"`
}

type Request struct {
	Method string  `json:"method"`
	Data   Product `json:"data"`
}

type Applink struct {
	IOSUrl              string `json:"ios_url"`
	IOSAppStoreID       string `json:"ios_app_store_id"`
	IOSAppName          string `json:"ios_app_name"`
	IPhoneUrl           string `json:"iphone_url"`
	IPhoneAppStoreID    string `json:"iphone_app_store_id"`
	IPhoneAppName       string `json:"iphone_app_name"`
	IPadUrl             string `json:"ipad_url"`
	IPadAppStoreID      string `json:"ipad_app_store_id"`
	IPadName            string `json:"ipad_app_name"`
	AndroidUrl          string `json:"android_url"`
	AndroidPackage      string `json:"android_package"`
	AndroidClass        string `json:"android_class"`
	AndroidName         string `json:"android_app_name"`
	WindowsPhoneUrl     string `json:"windows_phone_url"`
	WindowsPhoneAppID   string `json:"windows_phone_app_id"`
	WindowsPhoneAppName string `json:"windows_phone_app_name"`
}

type Media struct {
	Url string   `json:"url"`
	Tag []string `json:"tag"`
}

type Product struct {
	AdditionalImageLink        []string            `json:"additional_image_link,omitempty"`        // Optional.
	AdditionalVariantAttribute []map[string]string `json:"additional_variant_attribute,omitempty"` // Optional.
	AgeGroup                   string              `json:"age_group,omitempty"`                    // Optional. Group of people who are the same age or a similar age. Accepted values are newborn, infant, toddler, kids, adult.
	Applink                    *Applink            `json:"applink,omitempty"`                      // Optional.
	Availability               string              `json:"availability"`                           // Required. in stock / out of stock / available for order / discontinued
	Brand                      string              `json:"brand,omitempty"`                        // Optional.
	Color                      string              `json:"color,omitempty"`                        // Optional.
	Condition                  string              `json:"condition"`                              // Required. Product condition: new, refurbished, or used.
	CustomLabel0               string              `json:"custom_label_0,omitempty"`               // Optional.
	CustomLabel1               string              `json:"custom_label_1,omitempty"`               // Optional.
	CustomLabel2               string              `json:"custom_label_2,omitempty"`               // Optional.
	CustomLabel3               string              `json:"custom_label_3,omitempty"`               // Optional.
	CustomLabel4               string              `json:"custom_label_4,omitempty"`               // Optional.
	Description                string              `json:"description"`                            // Required.
	DisabledCapabilities       []string            `json:"disabled_capabilities,omitempty"`        // Optional.
	Gender                     string              `json:"gender,omitempty"`                       // Optional. Gender for sizing. Values include male, female, unisex.
	GoogleProductCategory      string              `json:"google_product_category,omitempty"`      // Optional.
	Gtin                       string              `json:"gtin,omitempty"`                         // Optional.
	Id                         string              `json:"id"`                                     // Required.
	Image                      []Media             `json:"image,omitempty"`                        // Optional.
	ImageLink                  string              `json:"image_link,omitempty"`                   // Not required if image is provided.
	Video                      []Media             `json:"video,omitempty"`                        // Optional.
	Inventory                  int                 `json:"inventory,omitempty"`                    // Optional
	ItemGroupId                string              `json:"item_group_id,omitempty"`                // Optional.
	Link                       string              `json:"link"`                                   // Required.
	ManufacturerPartNumber     string              `json:"manufacturer_part_number,omitempty"`     // Optional.
	Mpn                        string              `json:"mpn,omitempty"`                          // Optional.
	Material                   string              `json:"material,omitempty"`                     // Optional.
	Pattern                    string              `json:"pattern,omitempty"`                      // Optional.
	Price                      string              `json:"price"`                                  // Required.
	ProductTags                []string            `json:"product_tags,omitempty"`                 // Optional.
	SalePrice                  string              `json:"sale_price,omitempty"`                   // Optional, but required to use the Overlay feature for Advantage+ catalog ads.
	SalePriceEffectiveDate     string              `json:"sale_price_effective_date,omitempty"`    // Optional.
	Shipping                   string              `json:"shipping,omitempty"`                     // Optional.
	Size                       string              `json:"size,omitempty"`                         // Optional.
	Title                      string              `json:"title"`                                  // Required.
}

type Error struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FbTraceId string `json:"fbtrace_id"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s %d %s", e.Type, e.Code, e.Message)
}

type ValidationStatus struct {
	Errors     []Error `json:"errors"`
	RetailerId string  `json:"retailer_id"`
	Warnings   []Error `json:"warnings"`
}

type Response struct {
	Handles          []string           `json:"handles"`
	ValidationStatus []ValidationStatus `json:"validation_status"`
}
