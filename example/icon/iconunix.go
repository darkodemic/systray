// Modifications Copyright 2026 Darko Demić.
// Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.

//go:build linux || darwin || freebsd || openbsd || netbsd

package icon

import (
    _ "embed"
)

//go:embed icon.png
var Data []byte
