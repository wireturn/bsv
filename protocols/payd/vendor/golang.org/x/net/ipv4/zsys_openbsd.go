// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs defs_openbsd.go

package ipv4

const (
	sysIP_RECVDSTADDR = 0x7
	sysIP_RECVIF      = 0x1e
	sysIP_RECVTTL     = 0x1f

	sizeofIPMreq = 0x8
)

type ipMreq struct {
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
}
