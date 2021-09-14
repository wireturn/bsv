package bob

// InputAddresses returns the Bitcoin addresses for the transaction inputs
func (t *Tx) InputAddresses() (addresses []string) {
	for _, i := range t.In {
		if i.E.A != "" && i.E.A != "false" {
			addresses = append(addresses, i.E.A)
		}
	}
	return
}

// OutputAddresses returns the Bitcoin addresses for the transaction outputs
func (t *Tx) OutputAddresses() (addresses []string) {
	for _, i := range t.In {
		if i.E.A != "" && i.E.A != "false" {
			addresses = append(addresses, i.E.A)
		}
	}
	return
}
