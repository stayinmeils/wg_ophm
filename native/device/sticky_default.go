//go:build !linux

package device

import (
	"wg/native/conn"
	"wg/native/rwcancel"
)

func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
