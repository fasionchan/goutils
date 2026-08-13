package qrcode

import (
	sqrcode "github.com/skip2/go-qrcode"
	rqr "rsc.io/qr"
)

// skip2 把本包纠错级别映射为 skip2/go-qrcode 的级别。
func (l RecoveryLevel) skip2() sqrcode.RecoveryLevel {
	switch l {
	case LevelL:
		return sqrcode.Low
	case LevelQ:
		return sqrcode.High
	case LevelH:
		return sqrcode.Highest
	default:
		return sqrcode.Medium
	}
}

// terminal 把本包纠错级别映射为 rsc.io/qr 的级别（qrterminal 底层使用）。
func (l RecoveryLevel) terminal() rqr.Level {
	switch l {
	case LevelL:
		return rqr.L
	case LevelQ:
		return rqr.Q
	case LevelH:
		return rqr.H
	default:
		return rqr.M
	}
}
