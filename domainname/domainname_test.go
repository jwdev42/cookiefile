/* This file is part of the "cookiefile" library, ©2020 Jörg Walter
 *  This software is licensed under the "GNU Lesser General Public License version 3" */
package domainname

import "testing"

func TestValidNames(t *testing.T) {
	names := []string{
		"example.net",
		".example.net",
		"subdomain.example.net",
		"a.b.c",
		"m-y-hyphen-ex--ampl-e.example.net",
		"BigAndSmaLL.example.net",
		"Rocket69",
		"xn--dierzte-7wa.de",
	}

	for _, name := range names {
		v := &Validator{name: name}
		if err := v.Validate(); err != nil {
			t.Errorf("%v", err)
		}
	}
}

func TestInvalidNames(t *testing.T) {
	names := []string{
		"example.net.",
		".",
		"t..t",
		"t.-.t",
		"t.-0-.t",
		"-subdomain.example.net",
		"subdomain-.example.net",
		"subdomain-.example.net-",
		"subdomain.-example.net",
		"樹林.example.net",
		"green樹林.example.net",
		";test.example.net",
	}
	for _, name := range names {
		v := &Validator{name: name}
		err := v.Validate()
		if err == nil {
			t.Errorf("Test string %q: Error expected!", name)
		} else {
			t.Logf("Error as expected: %v", err)
		}
	}
}
