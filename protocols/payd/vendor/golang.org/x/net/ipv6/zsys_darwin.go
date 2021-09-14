// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs defs_darwin.go

package ipv6

const (
	sysIPV6_TCLASS   = 0x24
	sysIPV6_PATHMTU  = 0x2c
	sysIPV6_PKTINFO  = 0x2e
	sysIPV6_HOPLIMIT = 0x2f
	sysIPV6_NEXTHOP  = 0x30

	sizeofSockaddrStorage = 0x80
	sizeofSockaddrInet6   = 0x1c
	sizeofInet6Pktinfo    = 0x14
	sizeofIPv6Mtuinfo     = 0x20

	sizeofIPv6Mreq       = 0x14
	sizeofGroupReq       = 0x84
	sizeofGroupSourceReq = 0x104

	sizeofICMPv6Filter = 0x20
)

type sockaddrStorage struct {
	Len         uint8
	Family      uint8
	X__ss_pad1  [6]int8
	X__ss_align int64
	X__ss_pad2  [112]int8
}

type sockaddrInet6 struct {
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type inet6Pktinfo struct {
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
}

type ipv6Mtuinfo struct {
	Addr sockaddrInet6
	Mtu  uint32
}

type ipv6Mreq struct {
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
}

type icmpv6Filter struct {
	Filt [8]uint32
}

type groupReq struct {
	Interface uint32
	Pad_cgo_0 [128]byte
}

type groupSourceReq struct {
	Interface uint32
	Pad_cgo_0 [128]byte
	Pad_cgo_1 [128]byte
}
