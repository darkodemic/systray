// Modifications Copyright 2026 Darko Demić.
// Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.

//go:build windows

package icon

import (
	_ "embed"
)

//go:embed "iconwin.ico"
var Data []byte
