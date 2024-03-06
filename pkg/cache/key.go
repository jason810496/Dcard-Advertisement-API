package cache

import (
	"fmt"
	"strconv"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
)

func PublicAdKey(req *schemas.PublicAdRequest) string {
	return fmt.Sprintf(
		"ad:%s:%s:%s:%s",
		defaultInt(req.Age),
		defaultStr(req.Country),
		defaultStr(req.Platform),
		defaultStr(req.Gender))
}

func defaultInt(val int) string {
	if val == 0 {
		return "*"
	}
	return strconv.Itoa(val)
}

func defaultStr(val string) string {
	if val == "" {
		return "*"
	}
	return val
}
