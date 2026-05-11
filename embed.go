package magzhanme

import "embed"

//go:embed templates
var Templates embed.FS

//go:embed locales
var Locales embed.FS

//go:embed static
var Static embed.FS
