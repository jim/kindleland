//go:generate script/generate_constants
package kindleland

// #include "linux/types.h"
// #include "linux/einkfb.h"
import "C"

const (
	FxFlash         = int(C.fx_flash)
	FxUpdatePartial = int(C.fx_update_partial)
	FxUpdateFull    = int(C.fx_update_full)
)

const FBIOEinkUpdateDisplay = int(C.FBIO_EINK_UPDATE_DISPLAY)
