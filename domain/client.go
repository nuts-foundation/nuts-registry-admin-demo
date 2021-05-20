package domain

import "net"

func UnwrapAPIError(err error) error {
	if _, ok := err.(net.Error); ok {
		return ErrNutsNodeUnreachable
	}
	return err
}
