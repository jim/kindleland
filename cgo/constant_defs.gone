package kindleland

//go:generate script/generate_constants

// #include "linux/types.h"
// #include "linux/einkfb.h"
import "C"

const FrameBufferDevice = C.EINK_FRAME_BUFFER

type UpdateMode int

const (
	FxFlash         = UpdateMode(C.fx_flash)
	FxUpdatePartial = UpdateMode(C.fx_update_partial)
	FxUpdateFull    = UpdateMode(C.fx_update_full)
	FxUpdateFast    = UpdateMode(C.fx_update_fast)
	FxUpdateSlow    = UpdateMode(C.fx_update_slow)
)

// const FBIOEinkUpdateDisplay = int(C.FBIO_EINK_UPDATE_DISPLAY)
// const FBIOEinkUpdateDisplayArea = int(C.FBIO_EINK_UPDATE_DISPLAY_AREA)
const FBIOEinkUpdateDisplayFx = C.FBIO_EINK_UPDATE_DISPLAY_FX

// const FBIOEinkClearScreen = int(C.FBIO_EINK_CLEAR_SCREEN)

type UpdateArea C.update_area_t
