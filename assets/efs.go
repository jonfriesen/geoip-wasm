package assets

import _ "embed"

//go:embed "db/GeoLite2-City.mmdb"
var GeoLite2City []byte
