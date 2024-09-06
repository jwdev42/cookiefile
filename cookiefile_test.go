/* This file is part of the "cookiefile" library, ©2020 Jörg Walter
 *  This software is licensed under the "GNU Lesser General Public License version 3" */
package cookiefile

import (
	"fmt"
	"net/url"
	"testing"
)

func TestLoad(t *testing.T) {
	cookies, err := Load("test/valid.txt")
	if err != nil {
		t.Error(err)
	}
	for i, cookie := range cookies {
		t.Logf("Cookie %3d: %s", i, cookie)
	}
}

func TestLoadJar(t *testing.T) {
	jar, err := LoadJar("test/valid.txt")
	if err != nil {
		t.Error(err)
	}

	for _, host := range []string{"example.net", "httponly.net"} {
		addr, err := url.Parse(fmt.Sprintf("http://%s", host))
		if err != nil {
			t.Error(err)
		}
		if len(jar.Cookies(addr)) != 1 {
			t.Errorf("Expected 1 cookie for host %q, got %d", addr.Hostname(), len(jar.Cookies(addr)))
		}
	}
}
