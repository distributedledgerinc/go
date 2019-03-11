package txnbuild

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestNativeAssetToXDR(t *testing.T) {
	asset := Asset{}

	received, err := asset.ToXDR()
	assert.Nil(t, err)

	expected := xdr.Asset{Type: xdr.AssetTypeAssetTypeNative}
	assert.Equal(t, expected, received)
}

func TestAlphaNum4AssetToXDR(t *testing.T) {

	asset := Asset{
		Code:   "USD",
		Issuer: newKeypair0().Address(),
	}
	var xdrAssetCode [4]byte
	copy(xdrAssetCode[:], asset.Code)
	var xdrIssuer xdr.AccountId
	require.NoError(t, xdrIssuer.SetAddress(asset.Issuer))

	received, err := asset.ToXDR()
	assert.Nil(t, err)

	expected := xdr.Asset{Type: xdr.AssetTypeAssetTypeCreditAlphanum4,
		AlphaNum4: &xdr.AssetAlphaNum4{
			AssetCode: xdrAssetCode,
			Issuer:    xdrIssuer,
		}}
	assert.Equal(t, expected, received)
}
