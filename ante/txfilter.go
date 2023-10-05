package ante

import (
	"fmt"
	"regexp"
	"strings"

	servertypes "github.com/Finschia/finschia-sdk/server/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/spf13/cast"
)

const (
	FlagTxFilter        = "tx-filter.allowed-targets"
	FlagInitHeight      = "tx-filter.initial-block-height"
	FlagAllowedContract = "tx-filter.allowed-contracts"
	FlagDisableFilter   = "tx-filter.disable-contract-filter"
)

type TxFilterOptions struct {
	AllowedMsgRegex    string
	InitialBlockHeight int64
	AllowedContract    map[string]bool
}

func NewTxFilterOptions(appOpts servertypes.AppOptions) TxFilterOptions {
	return TxFilterOptions{
		AllowedMsgRegex:    GenAllowedMsgRegex(cast.ToStringSlice(appOpts.Get(FlagTxFilter))),
		InitialBlockHeight: cast.ToInt64(appOpts.Get(FlagInitHeight)),
		AllowedContract:    GenAllowedContractMap(cast.ToStringSlice(appOpts.Get(FlagAllowedContract)), cast.ToBool(appOpts.Get(FlagDisableFilter))),
	}
}

// TxFilterDecorator blocks the sdk.Msg if it's not registered in tx-filter AllowedTargets
type TxFilterDecorator struct {
	Options TxFilterOptions
}

func NewTxFilterDecorator(opts TxFilterOptions) TxFilterDecorator {
	return TxFilterDecorator{
		Options: opts,
	}
}

func (tfd TxFilterDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	} else if tfd.Options.InitialBlockHeight == ctx.BlockHeight() {
		return next(ctx, tx, simulate)
	}

	if ctx.IsCheckTx() {
		if err := tfd.CheckIfAllowed(tx.GetMsgs()); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

func (tfd TxFilterDecorator) CheckIfAllowed(msgs []sdk.Msg) error {
	MsgVerifier := regexp.MustCompile(tfd.Options.AllowedMsgRegex)
	for _, msg := range msgs {
		m := sdk.MsgTypeURL(msg)
		if !MsgVerifier.MatchString(m) {
			return fmt.Errorf("%s is not allowed on proxy node", m)
		}

		if strings.Split(m, ".")[1] == "wasm" {
			if err := tfd.filterWasmLogics(msg); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tfd TxFilterDecorator) filterWasmLogics(msg sdk.Msg) error {
	if wasmExec, ok := msg.(*wasmtypes.MsgExecuteContract); ok {
		if tfd.Options.AllowedContract == nil {
			return nil
		} else if !tfd.Options.AllowedContract[wasmExec.Contract] {
			return fmt.Errorf("%s is not allowed contract", wasmExec.Contract)
		}
	} else {
		return fmt.Errorf("%s is not allowed on proxy node", sdk.MsgTypeURL(msg))
	}
	return nil
}

func GenAllowedMsgRegex(targets []string) string {
	if targets == nil {
		return "^$"
	}

	regex := ``
	for _, mod := range targets {
		switch len(strings.Split(mod, ".")) {
		case 1, 2, 3:
			regex += fmt.Sprintf(`(^\/%s\.)|`, strings.ReplaceAll(mod, ".", `\.`))
		case 4:
			regex += fmt.Sprintf(`(^\/%s$)|`, strings.ReplaceAll(mod, ".", `\.`))
		default:
			panic("invalid tx-filter type")
		}
	}

	return regex[:len(regex)-1]
}

func GenAllowedContractMap(contracts []string, isOff bool) map[string]bool {
	if isOff {
		return nil
	}

	m := make(map[string]bool)
	for i := 0; i < len(contracts); i++ {
		m[contracts[i]] = true
	}
	return m
}
