package txnbuild

import (
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

// AccountFlag represents the bitmask flags used to set and clear account authorization options.
type AccountFlag uint32

// AuthRequired is a flag that requires the issuing account to give other accounts
// permission before they can hold the issuing account's credit.
const AuthRequired = AccountFlag(xdr.AccountFlagsAuthRequiredFlag)

// AuthRevocable is a flag that allows the issuing account to revoke its credit
// held by other accounts.
const AuthRevocable = AccountFlag(xdr.AccountFlagsAuthRevocableFlag)

// AuthImmutable is a flag that if set prevents any authorization flags from being
// set, and prevents the account from ever being merged (deleted).
const AuthImmutable = AccountFlag(xdr.AccountFlagsAuthImmutableFlag)

// Threshold is the datatype for MasterWeight, Signer.Weight, and Thresholds.
type Threshold *xdr.Uint32

// NewThreshold returns a pointer to the provided value, used by SetOptions, so that it
// can distinguish between a zero-valued weight or threshold and one that is unset.
func NewThreshold(t uint8) Threshold {
	toXDRType := xdr.Uint32(t)
	return &toXDRType
}

// SetOptions represents the Stellar set options operation. See
// https://www.stellar.org/developers/guides/concepts/list-of-operations.html
type SetOptions struct {
	destAccountID        xdr.AccountId
	InflationDestination string
	SetAuthorization     []AccountFlag
	ClearAuthorization   []AccountFlag
	MasterWeight         Threshold
	LowThreshold         Threshold
	MediumThreshold      Threshold
	HighThreshold        Threshold
	xdrOp                xdr.SetOptionsOp
}

// BuildXDR for SetOptions returns a fully configured XDR Operation.
func (so *SetOptions) BuildXDR() (xdr.Operation, error) {
	err := so.handleInflation()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "Failed to set inflation destination address")
	}

	so.handleClearFlags()
	so.handleSetFlags()
	so.handleMasterWeight()
	so.handleLowThreshold()
	so.handleMediumThreshold()
	so.handleHighThreshold()

	opType := xdr.OperationTypeSetOptions
	body, err := xdr.NewOperationBody(opType, so.xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "Failed to build XDR OperationBody")
	}

	return xdr.Operation{Body: body}, nil
}

// handleInflation for SetOptions sets the XDR inflation destination.
// Once set, a new address can be set, but there's no way to ever unset.
func (so *SetOptions) handleInflation() (err error) {
	if so.InflationDestination != "" {
		err = so.destAccountID.SetAddress(so.InflationDestination)
		if err != nil {
			return
		}
		so.xdrOp.InflationDest = &so.destAccountID
	}
	return
}

// handleSetFlags for SetOptions sets XDR account flags (represented as a bitmask).
// See https://www.stellar.org/developers/guides/concepts/accounts.html
func (so *SetOptions) handleSetFlags() {
	var flags xdr.Uint32
	for _, flag := range so.SetAuthorization {
		flags = flags | xdr.Uint32(flag)
	}
	if len(so.SetAuthorization) > 0 {
		so.xdrOp.SetFlags = &flags
	}
}

// handleClearFlags for SetOptions unsets XDR account flags (represented as a bitmask).
// See https://www.stellar.org/developers/guides/concepts/accounts.html
func (so *SetOptions) handleClearFlags() {
	var flags xdr.Uint32
	for _, flag := range so.ClearAuthorization {
		flags = flags | xdr.Uint32(flag)
	}
	if len(so.ClearAuthorization) > 0 {
		so.xdrOp.ClearFlags = &flags
	}
}

// handleMasterWeight for SetOptions set the XDR weight of the master signing key.
// See https://www.stellar.org/developers/guides/concepts/multi-sig.html
func (so *SetOptions) handleMasterWeight() {
	if so.MasterWeight != nil {
		so.xdrOp.MasterWeight = so.MasterWeight
	}
}

// handleLowThreshold for SetOptions set the XDR value of the account's "low" threshold.
// See https://www.stellar.org/developers/guides/concepts/multi-sig.html
func (so *SetOptions) handleLowThreshold() {
	if so.LowThreshold != nil {
		so.xdrOp.LowThreshold = so.LowThreshold
	}
}

// handleLowThreshold for SetOptions set the XDR value of the account's "medium" threshold.
// See https://www.stellar.org/developers/guides/concepts/multi-sig.html
func (so *SetOptions) handleMediumThreshold() {
	if so.MediumThreshold != nil {
		so.xdrOp.MedThreshold = so.MediumThreshold
	}
}

// handleLowThreshold for SetOptions set the XDR value of the account's "high" threshold.
// See https://www.stellar.org/developers/guides/concepts/multi-sig.html
func (so *SetOptions) handleHighThreshold() {
	if so.HighThreshold != nil {
		so.xdrOp.HighThreshold = so.HighThreshold
	}
}
