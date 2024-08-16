package utils

import (
	"github.com/go-resty/resty/v2"
)

// GetFavicon fetches the favicon from the given uri
func GetFavicon(faviconUri string) ([]byte, error) {
	client := resty.New()

	res, err := client.R().Get(faviconUri)
	if err != nil {
		return nil, err
	}
	defer res.RawResponse.Body.Close()

	return res.Body(), nil
}
