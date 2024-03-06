package ante_test

import (
	"testing"

	cryptotypes "github.com/Finschia/finschia-sdk/crypto/types"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	stakingtypes "github.com/Finschia/finschia-sdk/x/staking/types"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	lbmwasmtypes "github.com/Finschia/wasmd/x/wasmplus/types"
	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-proxy/v3/ante"
	linkhelper "github.com/Finschia/finschia-proxy/v3/app/helpers"
)

func TestGenWhiteRegex(t *testing.T) {
	tcs := []struct {
		allowed  []string
		expected string
		ispanic  bool
	}{
		{[]string{"cosmos.bank", "lbm.token.v1", "lbm.collection.v1.MsgTransferFT"}, `(^\/cosmos\.bank\.)|(^\/lbm\.token\.v1\.)|(^\/lbm\.collection\.v1\.MsgTransferFT$)`, false},
		{[]string{"wrong.prefix.non-standard.namespace.test", "cosmos.bank"}, "", true},
	}

	for _, tc := range tcs {
		tc := tc
		if tc.ispanic {
			require.Panics(t, func() { ante.GenAllowedMsgRegex(tc.allowed) }, "this type is not supported")
		} else {
			require.Equal(t, tc.expected, ante.GenAllowedMsgRegex(tc.allowed))
		}
	}
}

func TestTxFilterDecorator_CheckIfAllowed(t *testing.T) {
	tcs := []struct {
		msgs             []sdk.Msg
		allowedTargets   []string
		allowedContracts []string
		ErrStr           string
		pass             bool
		disableFilter    bool
	}{
		{
			[]sdk.Msg{&banktypes.MsgSend{}},
			[]string{"cosmos.bank.v1beta1.MsgSend"},
			nil,
			"",
			true,
			false,
		},
		{
			[]sdk.Msg{&banktypes.MsgSend{}},
			[]string{"cosmos.bank.v1beta1"},
			nil,
			"",
			true,
			false,
		},
		{
			[]sdk.Msg{&banktypes.MsgSend{}, &stakingtypes.MsgDelegate{}},
			[]string{"cosmos.bank", "cosmos.staking"},
			nil,
			"",
			true,
			false,
		},
		{
			[]sdk.Msg{&banktypes.MsgSend{}, &stakingtypes.MsgDelegate{}},
			[]string{"cosmos"},
			nil,
			"",
			true,
			false,
		},
		{
			[]sdk.Msg{&banktypes.MsgSend{}, &stakingtypes.MsgDelegate{}},
			[]string{"cosmos.bank", "cosmos.sta"},
			nil,
			"/cosmos.staking.v1beta1.MsgDelegate is not allowed on proxy node",
			false,
			false,
		},
		{
			[]sdk.Msg{&banktypes.MsgSend{}, &testdata.TestMsg{}},
			[]string{"cosmos.bank"},
			nil,
			"/testdata.TestMsg is not allowed on proxy node",
			false,
			false,
		},
		{
			[]sdk.Msg{},
			[]string{"cosmos.bank"},
			nil,
			"",
			true,
			false,
		},
		{
			[]sdk.Msg{&testdata.TestMsg{}},
			nil,
			nil,
			"/testdata.TestMsg is not allowed on proxy node",
			false,
			false,
		},
		{
			[]sdk.Msg{&wasmtypes.MsgExecuteContract{Contract: "allowedAddress"}},
			[]string{"cosmwasm.wasm"},
			[]string{"allowedAddress"},
			"",
			true,
			false,
		},
		{
			[]sdk.Msg{&wasmtypes.MsgExecuteContract{Contract: "forbiddenAddress"}},
			[]string{"cosmwasm.wasm"},
			[]string{"allowedAddress"},
			"forbiddenAddress is not allowed contract",
			false,
			false,
		},
		{
			[]sdk.Msg{&wasmtypes.MsgExecuteContract{Contract: "forbiddenAddress"}},
			[]string{"cosmwasm.wasm"},
			[]string{"allowedAddress"},
			"",
			true,
			true,
		},
		{
			[]sdk.Msg{&wasmtypes.MsgExecuteContract{Contract: "forbiddenAddress"}},
			[]string{"cosmwasm.wasm"},
			[]string{},
			"forbiddenAddress is not allowed contract",
			false,
			false,
		},
		{
			[]sdk.Msg{&lbmwasmtypes.MsgStoreCodeAndInstantiateContract{}},
			[]string{"cosmwasm.wasm"},
			nil,
			"/lbm.wasm.v1.MsgStoreCodeAndInstantiateContract is not allowed on proxy node",
			false,
			false,
		},
	}

	for _, tc := range tcs {
		opts := ante.TxFilterOptions{AllowedMsgRegex: ante.GenAllowedMsgRegex(tc.allowedTargets), AllowedContract: ante.GenAllowedContractMap(tc.allowedContracts, tc.disableFilter)}
		deco := ante.NewTxFilterDecorator(opts)
		err := deco.CheckIfAllowed(tc.msgs)
		if tc.pass {
			require.NoError(t, err)
		} else {
			require.Contains(t, err.Error(), tc.ErrStr)
		}
	}
}

// TestTxFilterDecorator tests only TP case because other unit tests cover all error cases.
func (suite *AnteTestSuite) TestTxFilterDecorator() {
	suite.SetupTest(false, 1)
	app, ctx := suite.app, suite.ctx
	accounts := suite.CreateTestAccounts(2)
	acc1, acc2 := accounts[0].acc, accounts[1].acc
	priv1 := accounts[0].priv
	feeAmount := sdk.NewCoins(sdk.NewInt64Coin("foo", 0))
	gasLimit := testdata.NewTestGasLimit()
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}

	amt := sdk.NewCoins(sdk.NewInt64Coin("foo", 7777))
	suite.Require().NoError(linkhelper.FundAccount(app, ctx, acc1.GetAddress(), amt))
	sendAmt := sdk.NewCoins(sdk.NewInt64Coin("foo", 777))
	msgs := []sdk.Msg{banktypes.NewMsgSend(acc1.GetAddress(), acc2.GetAddress(), sendAmt)}

	tc := TestCase{
		desc:     "test a TxFilterDecorator logic around other decorators",
		malleate: func() {},
		simulate: false,
		expPass:  true,
		expErr:   nil,
	}

	suite.Run(tc.desc, func() {
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
		suite.RunTestCase(privs, msgs, feeAmount, gasLimit, accNums, accSeqs, suite.ctx.ChainID(), tc)
	})
}
