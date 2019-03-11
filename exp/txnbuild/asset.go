package txnbuild

import (
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

// Asset represents assets on the Stellar network.
type Asset struct {
	Code   string
	Issuer string
}

// ToXDR for Asset produces a corresponding XDR asset.
func (a *Asset) ToXDR() (xdr.Asset, error) {

	// Native (Lumens) has no code or issuer, and is a no-op
	if a.Code == "" && a.Issuer == "" {
		return xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
	}

	var issuer xdr.AccountId
	err := issuer.SetAddress(a.Issuer)
	if err != nil {
		return xdr.Asset{}, err
	}

	length := len(a.Code)
	switch {
	case length >= 1 && length <= 4:
		var codeArray [4]byte
		byteArray := []byte(a.Code)
		copy(codeArray[:], byteArray[0:length])
		asset := xdr.AssetAlphaNum4{AssetCode: codeArray, Issuer: issuer}
		return xdr.NewAsset(xdr.AssetTypeAssetTypeCreditAlphanum4, asset)
	case length >= 5 && length <= 12:
		var codeArray [12]byte
		byteArray := []byte(a.Code)
		copy(codeArray[:], byteArray[0:length])
		asset := xdr.AssetAlphaNum12{AssetCode: codeArray, Issuer: issuer}
		return xdr.NewAsset(xdr.AssetTypeAssetTypeCreditAlphanum12, asset)
	default:
		return xdr.Asset{}, errors.New("Asset code length must be between 1 and 12 characters")
	}
}
