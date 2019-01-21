package mlog

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	strChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // 62 characters
)

// EnsureDir EnsureDir
func EnsureDir(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return fmt.Errorf("Could not create directory %v. %v", dir, err)
		}
	}
	return nil
}

var (
	forwardList = []string{"X-FORWARDED-FOR", "x-forwarded-for", "X-Forwarded-For"}
)

// GetRequestIP GetRequestIP
func GetRequestIP(r *http.Request) string {
	onlineip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		onlineip = r.RemoteAddr
	}
	onlineip = ParseIPTo4(onlineip)

	var res string
	for _, f := range forwardList {
		if ipProxy := r.Header.Get(f); len(ipProxy) > 0 {
			res = ipProxy
			break
		}
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return onlineip
	}

	var forwards []string
	if strings.Contains(res, ",") {
		forwards = strings.Split(res, ",")
	} else if strings.Contains(res, ":") {
		// ipv6会有问题 未兼容ipv6
		forwards = strings.Split(res, ":")
	} else {
		forwards = []string{res}
	}

	if IsValidateIP(forwards[0]) {
		// 获取X_FORWARDED_FOR的第一个ip
		return ParseIPTo4(forwards[0])
	} else if onlineip != "" && IsValidateIP(onlineip) {
		return onlineip
	} else {
		for index := len(forwards) - 1; index >= 0; index-- {
			if IsValidateIP(forwards[index]) {
				return ParseIPTo4(forwards[index])
			}
		}
	}
	return onlineip
}

// ParseIPTo4 ParseIPTo4
func ParseIPTo4(s string) string {
	return net.ParseIP(strings.TrimSpace(s)).To4().String()
}

// IsValidateIP IsValidateIP 是否为全球唯一的IP即广域网IP
func IsValidateIP(s string) bool {
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip.IsGlobalUnicast() {
		if ip4 := ip.To4(); ip4 != nil {
			switch true {
			case ip4[0] == 10:
				return false
			case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
				return false
			case ip4[0] == 192 && ip4[1] == 168:
				return false
			default:
				return true
			}
		}
	}
	return false
}

//----------------------------------------------------
// AllTrim AllTrim
func AllTrim(str string) string {
	searchs := []string{" ", "　", "\n", "\r", "\t"}
	replaces := []string{"+", "+", "+", "+", "+"}
	for i, search := range searchs {
		str = strings.Replace(str, search, replaces[i], -1)
	}
	return str
}
