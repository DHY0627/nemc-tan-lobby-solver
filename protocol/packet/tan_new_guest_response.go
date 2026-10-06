package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

// TanNewGuestResponse is sent by the transfer server to the room host
// when a new player enters the room.
//
// NOTE: the exact wire layout of this packet has NOT been confirmed against a
// capture. It is modelled after its sibling TanEnterRoomResponse: an int8 error
// code followed by a uint8-length list of uint32 player IDs, with the body
// present only when the error code is non-zero… i.e. only when it IS zero.
//
// IMPORTANT: the early return below is required, not cosmetic. This library's
// reader has no bounds checking and panics on EOF, so reading the list after a
// non-zero error code walks off the end of the buffer and panics. Do not
// "fix" the early return by removing it.
//
// 重要：这里「ErrorCode 非 0 就 return」是必须的，不是多余的。
// 本库的 reader 没有边界检查，越界读取会直接 panic；实测把 return 去掉后
// 网关会在收到该包时崩溃（panic: Uint32: EOF）。
// 另外 ErrorCode 打印出负数（如 -74）说明字段解析本身就存疑，
// 目前先保持原样，待抓到真实报文后再校正结构。
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
	if t.ErrorCode != 0 {
		return
	}
	encoding.FuncSliceUint8Length(io, &t.PlayerIDList, io.Uint32)
}
