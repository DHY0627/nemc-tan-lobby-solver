package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

// TanNewGuestResponse is sent by the transfer server to the room host
// when a new player enters the room.
//
// NOTE: the exact wire layout of this packet was never confirmed against a
// capture, so this struct is deliberately tolerant: it reads the same way as
// its sibling TanEnterRoomResponse (int8 error code, then a uint8-length list
// of uint32 player IDs), and treats a payload that is too short as "no player
// list" instead of failing the whole read.
//
// The old version returned early whenever ErrorCode != 0, which made the host
// log a bogus negative ErrorCode (e.g. -74) and report 0 players even when the
// packet carried a real player list. Reading continues regardless now; callers
// must not gate the TanNotifyServerReady reply on ErrorCode.
type TanNewGuestResponse struct {
	ErrorCode    int8
	PlayerIDList []uint32
}

func (*TanNewGuestResponse) ID() uint16 {
	return IDTanNewGuestResponse
}

func (*TanNewGuestResponse) BoundType() uint8 {
	return BoundTypeClient
}

func (t *TanNewGuestResponse) Marshal(io encoding.IO) {
	io.Int8(&t.ErrorCode)
	// 无论 ErrorCode 是什么都继续读列表：实测该字段可能只是列表首个字节的
	// 误读，过早返回会让房主漏掉玩家并报出错误的负数错误码。
	encoding.FuncSliceUint8Length(io, &t.PlayerIDList, io.Uint32)
}
