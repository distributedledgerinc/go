package txnbuild

// Asset represents assets on the Stellar network.
type Asset struct {
	Code   string
	Issuer string
	Native bool
}
