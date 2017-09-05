package functions

import (
	"bytes"
	"fmt"
)

// GetAvatarURL get avatar url by img name
func GetAvatarURL(img string) string {
	if img == "" {
		return "/assets/image/avatar/default.png"
	}

	return fmt.Sprintf("/assets/image/avatar/%s", img)
}

// GeneratePager GeneratePager
func GeneratePager(page, total, limit int, sortBy, URL string, filter map[string]interface{}) string {
	var pagerTotal int
	var pagerCount = 8
	var prevPage = "上一页"
	var nextPage = "下一页"
	if total%page == 0 {
		pagerTotal = total / limit
	} else {
		pagerTotal = total/limit + 1
	}

	filterStr := new(bytes.Buffer)
	if filter != nil {
		for k, v := range filter {
			if k == "" {
				continue
			}
			filterStr.WriteString("&")
			filterStr.WriteString(k)
			filterStr.WriteString("=")
			filterStr.WriteString(v.(string))
		}
	}

	buf := new(bytes.Buffer)
	buf.WriteString(`<div class="pager"><ul>`)
	if page-1 >= 1 {
		prevHref := fmt.Sprintf("%s?page=%d&sort=%s%s", URL, page-1, sortBy, filterStr.String())
		buf.WriteString(`<li><a href="`)
		buf.WriteString(prevHref)
		buf.WriteString(`">`)
		buf.WriteString(prevPage)
		buf.WriteString(`</a></li>`)
	}

	var href = ""
	if pagerTotal <= pagerCount {
		pagerCount = pagerTotal
	}
	for i := 0; i < pagerCount; i++ {
		href = fmt.Sprintf("%s?page=%d&sort=%s%s", URL, page+i, sortBy, filterStr.String())

		buf.WriteString(`<li><a href="`)
		buf.WriteString(href)
		buf.WriteString(`">`)
		buf.WriteString(prevPage)
		buf.WriteString(`</a></li>`)
	}

	if page+1 < pagerTotal {
		nextHref := fmt.Sprintf("%s?page=%d&sort=%s%s", URL, page+1, sortBy, filterStr.String())
		buf.WriteString(`<li><a href="`)
		buf.WriteString(nextHref)
		buf.WriteString(`">`)
		buf.WriteString(nextPage)
		buf.WriteString(`</a></li>`)
	}
	buf.WriteString("</uL></div>")

	return buf.String()
}
