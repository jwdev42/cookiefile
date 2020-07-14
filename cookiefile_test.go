/* This file is part of the "cookiefile" library, ©2020 Jörg Walter
 *  This software is licensed under the "GNU Lesser General Public License version 3" */
package cookiefile

import (
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
	_, err := LoadJar("test/valid.txt")
	if err != nil {
		t.Error(err)
	}
}
