/* This file is part of the "cookiefile" library, ©2020 Jörg Walter
 *  This software is licensed under the "GNU Lesser General Public License version 3" */

package cookiefile

import (
	"bufio"
	"fmt"
	"github.com/jwdev42/cookiefile/domainname"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseBool(input string) (bool, error) {
	if input == "TRUE" {
		return true, nil
	} else if input == "FALSE" {
		return false, nil
	}
	return false, fmt.Errorf("Invalid boolean expression: %s", input)
}

func parseTime(input string) (time.Time, error) {
	t, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(t, 0), nil
}

func validateDomain(input string) error {
	v := domainname.NewValidator(input)
	return v.Validate()
}

func validateName(input string) error {
	separators := "()<>@,;:\"\\/[]?={}"
	for i, c := range input {
		if c < 0x21 || c > 0x7e || strings.IndexRune(separators, c) > -1 {
			return fmt.Errorf("Name %q, illegal character \"%c\" at index %d", input, c, i)
		}
	}
	return nil
}

func validatePath(input string) error {
	for i, c := range input {
		if c < 0x21 || c == 0x3b || c > 0x7e {
			return fmt.Errorf("Path %q, illegal character \"%c\" at index %d", input, c, i)
		}
	}
	return nil
}

func validateValue(input string) error {
	excluded := "\",;\\"
	for i, c := range input {
		if c < 0x21 || c > 0x7e || strings.IndexRune(excluded, c) > -1 {
			return fmt.Errorf("Value %q, illegal character \"%c\" at index %d", input, c, i)
		}
	}
	return nil
}

func isWhitespaceAscii(r rune) bool {
	if r == 0x9 || r == 0x20 {
		return true
	}
	return false
}

func parseLine(line string) (*http.Cookie, error) {
	if "" == strings.TrimFunc(line, isWhitespaceAscii) || '#' == strings.TrimLeftFunc(line, isWhitespaceAscii)[0] {
		return nil, nil
	}
	entries := strings.Split(line, "\t")
	if len(entries) != 7 {
		return nil, fmt.Errorf("Invalid amount of fields")
	}
	cookie := &http.Cookie{}

	/*set domain*/
	if err := validateDomain(entries[0]); err != nil {
		return nil, err
	}
	cookie.Domain = entries[0]

	/*set path*/
	if err := validatePath(entries[2]); err != nil {
		return nil, err
	} else {
		cookie.Path = entries[2]
	}

	/*set "secure" flag*/
	if secure, err := parseBool(entries[3]); err != nil {
		return nil, err
	} else {
		cookie.Secure = secure
	}

	/*set expiration date*/
	if t, err := parseTime(entries[4]); err != nil {
		return nil, err
	} else {
		cookie.Expires = t
	}

	/*set name*/
	if err := validateName(entries[5]); err != nil {
		return nil, err
	} else {
		cookie.Name = entries[5]
	}

	/*set value*/
	if err := validateValue(entries[6]); err != nil {
		return nil, err
	} else {
		cookie.Value = entries[6]
	}
	return cookie, nil
}

func Load(path string) ([]*http.Cookie, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	cookies := make([]*http.Cookie, 0, 10)
	readbuf := bufio.NewScanner(f)
	for i := 1; readbuf.Scan(); i++ {
		line := readbuf.Text()
		cookie, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("Line %d: %w", i, err)
		}
		if cookie != nil {
			cookies = append(cookies, cookie)
		}
	}
	if err := readbuf.Err(); err != nil {
		return nil, err
	}
	return cookies, nil
}

func LoadJar(path string) (http.CookieJar, error) {
	getHost := func(cookie *http.Cookie) string {
		h := cookie.Domain
		if len(h) > 0 && h[0] == '.' {
			return h[1:]
		}
		return h
	}
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	cookiemap := make(map[string][]*http.Cookie)
	readbuf := bufio.NewScanner(f)
	for i := 1; readbuf.Scan(); i++ {
		line := readbuf.Text()
		cookie, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("Line %d: %w", i, err)
		}
		if cookie != nil {
			host := getHost(cookie)
			_, ok := cookiemap[host]
			if !ok {
				cookiemap[host] = make([]*http.Cookie, 0, 5)
			}
			cookiemap[host] = append(cookiemap[host], cookie)
		}
	}
	if err := readbuf.Err(); err != nil {
		return nil, err
	}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	for k, v := range cookiemap {
		u, err := url.Parse("http://" + k)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(u, v)
	}
	return jar, nil
}
