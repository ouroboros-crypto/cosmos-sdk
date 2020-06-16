package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetTxSign(t *testing.T) {
	clientCtx := client.Context{}
	// clientCtx = clientCtx.WithTxGenerator(simappparams.MakeEncodingConfig().TxGenerator).WithChainID("test")

	dir, clean := tests.NewTestCaseDir(t)
	t.Cleanup(clean)

	path := hd.CreateHDPath(118, 0, 0).String()
	kr, err := keyring.New(t.Name(), "test", dir, nil)
	require.NoError(t, err)

	var from = "test_sign"

	_, seed, err := kr.NewMnemonic(from, keyring.English, path, hd.Secp256k1)
	require.NoError(t, err)
	require.NoError(t, kr.Delete(from))

	info, err := kr.NewAccount(from, seed, "", path, hd.Secp256k1)
	require.NoError(t, err)

	viper.Set(flags.FlagGenerateOnly, true)
	viper.Set(flags.FlagFrom, info.GetAddress())

	clientCtx = clientCtx.WithTxGenerator(simappparams.MakeEncodingConfig().TxGenerator).WithChainID("test").WithKeyring(kr).WithFrom(from)

	cmd := GetSignCommand(clientCtx)
	encodingConfig := simappparams.MakeEncodingConfig()
	authtypes.RegisterCodec(encodingConfig.Amino)
	sdk.RegisterCodec(encodingConfig.Amino)

	txGen := encodingConfig.TxGenerator

	// Build a test transaction
	fee := authtypes.NewStdFee(50000, sdk.Coins{sdk.NewInt64Coin("atom", 150)})
	stdTx := authtypes.NewStdTx([]sdk.Msg{}, fee, []authtypes.StdSignature{}, "foomemo")
	txJSONEncoded, err := txGen.TxJSONEncoder()(stdTx)
	require.NoError(t, err)

	txFile, cleanup := tests.WriteToNewTempFile(t, string(txJSONEncoded))
	txFileName := txFile.Name()
	t.Cleanup(cleanup)

	err = cmd.RunE(cmd, []string{txFileName})
	require.NoError(t, err)
}