package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/evmos/ethermint/encoding"
)

func TestEthermintAppExport(t *testing.T) {
	db := dbm.NewMemDB()
	app := SetupWithDB(false, nil, db)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewEthermintApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		encoding.MakeConfig(ModuleBasics),
		simtestutil.NewAppOptionsWithFlagHome(DefaultNodeHome),
		baseapp.SetChainID(ChainID))
	_, err := app2.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
