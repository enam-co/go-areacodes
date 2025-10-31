# go-areacodes

`go-areacodes` is a Go port of [BlueRival's node-areacodes](https://github.com/BlueRival/node-areacodes) library.  
It provides an in-memory database of NANPA area codes (U.S. states plus toll-free numbers) along with
their canonical city, state, and centroid coordinates. The data set is embedded in the module, so
lookups are fast and do not require any network access.

## Getting started

```bash
go get github.com/josephfinlayson/go-areacodes
```

### Example

```go
package main

import (
	"fmt"

	"github.com/josephfinlayson/go-areacodes"
)

func main() {
	entry, err := areacodes.Lookup("+1-303-123-4567")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s, %s (%s) — %s\n", entry.City, entry.State, entry.StateCode, entry.Type)
}
```

### API surface

* `Lookup(phone string) (Entry, error)` – normalises a phone number, extracts the area code, and
  returns the matched record. Errors are comparable to the exported `ErrInvalidPhone` and `ErrNotFound`.
* `All() (map[string]Entry, error)` – returns a defensive copy of the entire database keyed by
  three-digit area-code strings.

Entries mirror the JavaScript version:

```go
type Entry struct {
	Type      string
	City      string
	State     string
	StateCode string
	Location  struct {
		Latitude  float64
		Longitude float64
	}
}
```

## Accuracy and data provenance

- Primary data originates from the NANPA CSV exports at <https://www.nationalnanpa.com>
- City names are sourced from NANPA's city reports
- Polygon centroid data is supplied by the USGS shapefile at <https://www.sciencebase.gov>
- Missing location details were back-filled using the OpenStreetMap Nominatim service

This port carries the same caveats as the original project: area code boundaries are approximate
and any coordinates should be treated as ±200 miles from the true location.

## Development

```bash
go test ./...
```

The update utilities found in the Node.js repository (`update.js`, `lib/reference.js`) have not yet
been ported. The embedded `data/data.json` file is intentionally kept compatible so it can be
regenerated with the original tooling.

## License

MIT – see [LICENSE](LICENSE).
