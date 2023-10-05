package ante

type AppOptionsMock map[string]any

func MakeAppOptionsMock() AppOptionsMock {
	return map[string]any{
		FlagTxFilter:        DefaultTxFilter(),
		FlagInitHeight:      0,
		FlagAllowedContract: []string{},
		FlagDisableFilter:   false,
	}
}

func (msc AppOptionsMock) Get(key string) any {
	return msc[key]
}

func DefaultTxFilter() []string {
	return []string{"cosmos.bank"}
}
