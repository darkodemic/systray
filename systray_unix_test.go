// Modifications Copyright 2026 Darko Demić.
// Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.

//go:build (linux || freebsd || openbsd || netbsd) && !android

package systray

import (
	"testing"

	"github.com/darkodemic/systray/internal/generated/notifier"
)

func TestPropSpecCoversIntrospection(t *testing.T) {
	served := instance.createPropSpec()["org.kde.StatusNotifierItem"]
	for _, p := range notifier.IntrospectDataStatusNotifierItem.Properties {
		if _, ok := served[p.Name]; !ok {
			t.Errorf("property %q is advertised via introspection but not served", p.Name)
		}
	}
}
