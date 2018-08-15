package asdf

import (
	"net/url"
	"strconv"
)

// param[0] is path
// param[1] is parameters
// param[2] is query
// param[3] is fragment
//
// scheme:[//[user:password@]host[:port]][/]path[?query][#fragment]
func NewUrl(scheme, host, port string, param ...string) string {
	url := scheme + "://" + host + ":" + port + "/"
	split := [4]string{"", ";", "?", "#"}

	count := len(param)
	for i := 1; i < count; i++ {
		if Empty != param[i] {
			url += split[i] + param[i]
		} else {
			break
		}
	}

	return url
}

func HttpBaseUrl(host, port, path string) string {
	return NewUrl("http", host, port, path)
}

func HttpsBaseUrl(host, port, path string) string {
	return NewUrl("https", host, port, path)
}

func UrlParamInt(params url.Values, key string) (int, bool, error) {
	if values, ok := params[key]; ok {
		v, err := strconv.Atoi(values[0])
		if nil != err {
			Log.Error("url param %s: %d error:%s", v, err)
		}

		return v, true, err
	} else {
		return 0, false, nil
	}
}
