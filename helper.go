package cas

import (
	"fmt"
	"strings"
)

func (c *Client) createURL(pathName string) string {
	if strings.HasSuffix(c.baseURL, "/") {
		return fmt.Sprintf("%s%s", c.baseURL, pathName)
	} else {
		return fmt.Sprintf("%s/%s", c.baseURL, pathName)
	}
}
