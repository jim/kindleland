package kindleland

import "time"

type EventTime struct {
	Seconds      int32
	Microseconds int32
}

type Event struct {
	Time  EventTime
	Type  uint16
	Code  uint16
	Value int32
}

type KeyType int

const (
	KeyA KeyType = iota + 30
	KeyS
	KeyD
	KeyF
	KeyG
	KeyH
	KeyJ
	KeyK
	KeyL
)

const (
	KeyQ KeyType = iota + 16
	KeyW
	KeyE
	KeyR
	KeyT
	KeyY
	KeyU
	KeyI
	KeyO
	KeyP
)

const (
	KeyZ KeyType = iota + 44
	KeyX
	KeyC
	KeyV
	KeyB
	KeyN
	KeyM
)

const (
	KeyDelete KeyType = 14
	KeyReturn KeyType = 28
	KeyShift  KeyType = 42
	KeyPeriod KeyType = 52
	KeyAlt    KeyType = 56
	KeySpace  KeyType = 57
	KeyHome   KeyType = 102
	KeySym    KeyType = 126
	KeyMenu   KeyType = 139
	KeyBack   KeyType = 158
	KeyText   KeyType = 190
)

type KeyEventType int

const (
	KeyUp KeyEventType = iota
	KeyDown
	KeyHold
)

type KeyboardEvent struct {
	Time time.Time
	Type KeyEventType
	Key  KeyType
}
