package functions

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"
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
	var showCount = 8
	var prevPage = "上一页"
	var nextPage = "下一页"
	if total == 0 || limit == 0 || total <= limit {
		return ""
	}
	if total%limit == 0 {
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
	if pagerTotal < showCount {
		showCount = pagerTotal
	}
	var href string
	var pageFlag int
	if showCount%2 == 0 {
		pageFlag = showCount / 2
	} else {
		pageFlag = showCount/2 + 1
	}
	//当前页之前的页码
	for i := pageFlag; i > 0; i-- {
		if page-i < 1 {
			continue
		}
		fmt.Println(i, page-i)
		href = fmt.Sprintf("%s?page=%d&sort=%s%s", URL, page-i, sortBy, filterStr.String())

		buf.WriteString(`<li><a href="`)
		buf.WriteString(href)
		buf.WriteString(`">`)

		pageName := page - i
		buf.WriteString(strconv.Itoa(pageName))
		buf.WriteString(`</a></li>`)
	}

	// 当前页码
	href = fmt.Sprintf("%s?page=%d&sort=%s%s", URL, page, sortBy, filterStr.String())
	buf.WriteString(`<li><a style="background-color:#ddd;" href="`)
	buf.WriteString(href)
	buf.WriteString(`">`)
	buf.WriteString(strconv.Itoa(page))
	buf.WriteString(`</a></li>`)

	// 当前页码之后的页码
	for i := 1; ; i++ {
		if page+i > showCount || i > pageFlag {
			break
		}
		href = fmt.Sprintf("%s?page=%d&sort=%s%s", URL, page+i, sortBy, filterStr.String())

		buf.WriteString(`<li><a href="`)
		buf.WriteString(href)
		buf.WriteString(`">`)

		pageName := page + i
		buf.WriteString(strconv.Itoa(pageName))
		buf.WriteString(`</a></li>`)
	}

	if page+1 <= pagerTotal {
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

// GetFlashCookie 生成一次性cookie
func GetFlashCookie(name, value string) *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.Value = value
	fmt.Println(value)
	cookie.Expires = time.Now().Add(time.Second * 5)
	return cookie
}

// SetCookieExpired set cookie's expires is now
func SetCookieExpired(cookie *http.Cookie) *http.Cookie {
	cookie.Expires = time.Now().Add(time.Second * -5)
	return cookie
}
