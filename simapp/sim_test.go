package simapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsimops "github.com/cosmos/cosmos-sdk/x/auth/simulation/operations"
	banksimops "github.com/cosmos/cosmos-sdk/x/bank/simulation/operations"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrsimops "github.com/cosmos/cosmos-sdk/x/distribution/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govsimops "github.com/cosmos/cosmos-sdk/x/gov/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsimops "github.com/cosmos/cosmos-sdk/x/params/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingsimops "github.com/cosmos/cosmos-sdk/x/slashing/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingsimops "github.com/cosmos/cosmos-sdk/x/staking/simulation/operations"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

func init() {
	GetSimulatorFlags()
}

// TODO: add description
func testAndRunTxs(app *SimApp, config simulation.Config) []simulation.WeightedOperation {

	cdc := MakeCodec()
	ap := make(simulation.AppParams)

	if config.ParamsFile != "" {
		bz, err := ioutil.ReadFile(config.ParamsFile)
		if err != nil {
			panic(err)
		}

		cdc.MustUnmarshalJSON(bz, &ap)
	}

	return []simulation.WeightedOperation{
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightDeductFee, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			authsimops.SimulateDeductFee(app.AccountKeeper, app.SupplyKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgSend, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			banksimops.SimulateMsgSend(app.AccountKeeper, app.BankKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightSingleInputMsgMultiSend, &v, nil,
					func(_ *rand.Rand) {
						v = 10
					})
				return v
			}(nil),
			banksimops.SimulateSingleInputMsgMultiSend(app.AccountKeeper, app.BankKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgSetWithdrawAddress, &v, nil,
					func(_ *rand.Rand) {
						v = 50
					})
				return v
			}(nil),
			distrsimops.SimulateMsgSetWithdrawAddress(app.DistrKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgWithdrawDelegationReward, &v, nil,
					func(_ *rand.Rand) {
						v = 50
					})
				return v
			}(nil),
			distrsimops.SimulateMsgWithdrawDelegatorReward(app.DistrKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgWithdrawValidatorCommission, &v, nil,
					func(_ *rand.Rand) {
						v = 50
					})
				return v
			}(nil),
			distrsimops.SimulateMsgWithdrawValidatorCommission(app.DistrKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightSubmitVotingSlashingTextProposal, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			govsimops.SimulateSubmittingVotingAndSlashingForProposal(app.GovKeeper, govsimops.SimulateTextProposalContent),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightSubmitVotingSlashingCommunitySpendProposal, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			govsimops.SimulateSubmittingVotingAndSlashingForProposal(app.GovKeeper, distrsimops.SimulateCommunityPoolSpendProposalContent(app.DistrKeeper)),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightSubmitVotingSlashingParamChangeProposal, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			govsimops.SimulateSubmittingVotingAndSlashingForProposal(app.GovKeeper, paramsimops.SimulateParamChangeProposalContent),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgDeposit, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			govsimops.SimulateMsgDeposit(app.GovKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgCreateValidator, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgCreateValidator(app.AccountKeeper, app.StakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgEditValidator, &v, nil,
					func(_ *rand.Rand) {
						v = 5
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgEditValidator(app.StakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgDelegate, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgDelegate(app.AccountKeeper, app.StakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgUndelegate, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgUndelegate(app.AccountKeeper, app.StakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgBeginRedelegate, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			stakingsimops.SimulateMsgBeginRedelegate(app.AccountKeeper, app.StakingKeeper),
		},
		{
			func(_ *rand.Rand) int {
				var v int
				ap.GetOrGenerate(cdc, OpWeightMsgUnjail, &v, nil,
					func(_ *rand.Rand) {
						v = 100
					})
				return v
			}(nil),
			slashingsimops.SimulateMsgUnjail(app.SlashingKeeper),
		},
	}
}

func invariants(app *SimApp) []sdk.Invariant {
	// TODO: fix PeriodicInvariants, it doesn't seem to call individual invariants for a period of 1
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/4631
	if flagPeriodValue == 1 {
		return app.CrisisKeeper.Invariants()
	}
	return simulation.PeriodicInvariants(app.CrisisKeeper.Invariants(), flagPeriodValue, 0)
}

// Pass this in as an option to use a dbStoreAdapter instead of an IAVLStore for simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

// Profile with:
// /usr/local/go/bin/go test -benchmem -run=^$ github.com/cosmos/cosmos-sdk/simapp -bench ^BenchmarkFullAppSimulation$ -Commit=true -cpuprofile cpu.out
func BenchmarkFullAppSimulation(b *testing.B) {
	logger := log.NewNopLogger()
	config := NewConfigFromFlags()

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)
	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()
	app := NewSimApp(logger, db, nil, true, 0)

	// Run randomized simulation
	// TODO: parameterize numbers, save for a later PR
	_, params, simErr := simulation.SimulateFromSeed(
		b, os.Stdout, app.BaseApp, AppStateFn,
		testAndRunTxs(app, config), invariants(app),
		app.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		fmt.Println("Exporting app state...")
		appState, _, err := app.ExportAppStateAndValidators(false, nil)
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}
		err = ioutil.WriteFile(config.ExportStatePath, []byte(appState), 0644)
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if config.ExportParamsPath != "" {
		fmt.Println("Exporting simulation params...")
		paramsBz, err := json.MarshalIndent(params, "", " ")
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}

		err = ioutil.WriteFile(config.ExportParamsPath, paramsBz, 0644)
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if simErr != nil {
		fmt.Println(simErr)
		b.FailNow()
	}

	if config.Commit {
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
}

func TestFullAppSimulation(t *testing.T) {
	if !flagEnabledValue {
		t.Skip("Skipping application simulation")
	}

	var logger log.Logger
	config := NewConfigFromFlags()

	if flagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()

	app := NewSimApp(logger, db, nil, true, 0, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", app.Name())

	// Run randomized simulation
	_, params, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn,
		testAndRunTxs(app, config), invariants(app),
		app.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		fmt.Println("Exporting app state...")
		appState, _, err := app.ExportAppStateAndValidators(false, nil)
		require.NoError(t, err)

		err = ioutil.WriteFile(config.ExportStatePath, []byte(appState), 0644)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		fmt.Println("Exporting simulation params...")
		fmt.Println(params)
		paramsBz, err := json.MarshalIndent(params, "", " ")
		require.NoError(t, err)

		err = ioutil.WriteFile(config.ExportParamsPath, paramsBz, 0644)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
}

func TestAppImportExport(t *testing.T) {
	if !flagEnabledValue {
		t.Skip("Skipping application import/export simulation")
	}

	var logger log.Logger
	config := NewConfigFromFlags()

	if flagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()

	app := NewSimApp(logger, db, nil, true, 0, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", app.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn,
		testAndRunTxs(app, config), invariants(app),
		app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	if config.ExportStatePath != "" {
		fmt.Println("Exporting app state...")
		appState, _, err := app.ExportAppStateAndValidators(false, nil)
		require.NoError(t, err)

		err = ioutil.WriteFile(config.ExportStatePath, []byte(appState), 0644)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		fmt.Println("Exporting simulation params...")
		simParamsBz, err := json.MarshalIndent(simParams, "", " ")
		require.NoError(t, err)

		err = ioutil.WriteFile(config.ExportParamsPath, simParamsBz, 0644)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}

	fmt.Printf("Exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)
	fmt.Printf("Importing genesis...\n")

	newDir, _ := ioutil.TempDir("", "goleveldb-app-sim-2")
	newDB, _ := sdk.NewLevelDB("Simulation-2", dir)

	defer func() {
		newDB.Close()
		_ = os.RemoveAll(newDir)
	}()

	newApp := NewSimApp(log.NewNopLogger(), newDB, nil, true, 0, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", newApp.Name())

	var genesisState GenesisState
	err = app.cdc.UnmarshalJSON(appState, &genesisState)
	require.NoError(t, err)

	ctxB := newApp.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, genesisState)

	fmt.Printf("Comparing stores...\n")
	ctxA := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	type StoreKeysPrefixes struct {
		A        sdk.StoreKey
		B        sdk.StoreKey
		Prefixes [][]byte
	}

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.keys[baseapp.MainStoreKey], newApp.keys[baseapp.MainStoreKey], [][]byte{}},
		{app.keys[auth.StoreKey], newApp.keys[auth.StoreKey], [][]byte{}},
		{app.keys[staking.StoreKey], newApp.keys[staking.StoreKey],
			[][]byte{
				staking.UnbondingQueueKey, staking.RedelegationQueueKey, staking.ValidatorQueueKey,
			}}, // ordering may change but it doesn't matter
		{app.keys[slashing.StoreKey], newApp.keys[slashing.StoreKey], [][]byte{}},
		{app.keys[mint.StoreKey], newApp.keys[mint.StoreKey], [][]byte{}},
		{app.keys[distr.StoreKey], newApp.keys[distr.StoreKey], [][]byte{}},
		{app.keys[supply.StoreKey], newApp.keys[supply.StoreKey], [][]byte{}},
		{app.keys[params.StoreKey], newApp.keys[params.StoreKey], [][]byte{}},
		{app.keys[gov.StoreKey], newApp.keys[gov.StoreKey], [][]byte{}},
	}

	for _, storeKeysPrefix := range storeKeysPrefixes {
		storeKeyA := storeKeysPrefix.A
		storeKeyB := storeKeysPrefix.B
		prefixes := storeKeysPrefix.Prefixes

		storeA := ctxA.KVStore(storeKeyA)
		storeB := ctxB.KVStore(storeKeyB)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		fmt.Printf("Compared %d key/value pairs between %s and %s\n", len(failedKVAs), storeKeyA, storeKeyB)
		require.Len(t, failedKVAs, 0, GetSimulationLog(storeKeyA.Name(), app.sm.StoreDecoders, app.cdc, failedKVAs, failedKVBs))
	}

}

func TestAppSimulationAfterImport(t *testing.T) {
	if !flagEnabledValue {
		t.Skip("Skipping application simulation after import")
	}

	var logger log.Logger
	config := NewConfigFromFlags()

	if flagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ := sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()

	app := NewSimApp(logger, db, nil, true, 0, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", app.Name())

	// Run randomized simulation
	stopEarly, params, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn,
		testAndRunTxs(app, config), invariants(app),
		app.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		fmt.Println("Exporting app state...")
		appState, _, err := app.ExportAppStateAndValidators(false, nil)
		require.NoError(t, err)

		err = ioutil.WriteFile(config.ExportStatePath, []byte(appState), 0644)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		fmt.Println("Exporting simulation params...")
		paramsBz, err := json.MarshalIndent(params, "", " ")
		require.NoError(t, err)

		err = ioutil.WriteFile(config.ExportParamsPath, paramsBz, 0644)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}

	if stopEarly {
		// we can't export or import a zero-validator genesis
		fmt.Printf("We can't export or import a zero-validator genesis, exiting test...\n")
		return
	}

	fmt.Printf("Exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err)

	fmt.Printf("Importing genesis...\n")

	newDir, _ := ioutil.TempDir("", "goleveldb-app-sim-2")
	newDB, _ := sdk.NewLevelDB("Simulation-2", dir)

	defer func() {
		newDB.Close()
		_ = os.RemoveAll(newDir)
	}()

	newApp := NewSimApp(log.NewNopLogger(), newDB, nil, true, 0, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", newApp.Name())
	newApp.InitChain(abci.RequestInitChain{
		AppStateBytes: appState,
	})

	// Run randomized simulation on imported app
	_, _, err = simulation.SimulateFromSeed(
		t, os.Stdout, newApp.BaseApp, AppStateFn,
		testAndRunTxs(newApp, config), invariants(newApp),
		newApp.ModuleAccountAddrs(), config,
	)

	require.NoError(t, err)
}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on the application.
func TestAppStateDeterminism(t *testing.T) {
	if !flagEnabledValue {
		t.Skip("Skipping application simulation")
	}

	config := NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			logger := log.NewNopLogger()
			db := dbm.NewMemDB()
			app := NewSimApp(logger, db, nil, true, 0)

			fmt.Printf(
				"Running non-determinism simulation; seed: %d/%d (%d), attempt: %d/%d\n",
				i+1, numSeeds, config.Seed, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t, os.Stdout, app.BaseApp, AppStateFn,
				testAndRunTxs(app, config), []sdk.Invariant{},
				app.ModuleAccountAddrs(), config,
			)
			require.NoError(t, err)

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(t, appHashList[0], appHashList[j], "appHash list: %v", appHashList)
			}
		}
	}
}

func BenchmarkInvariants(b *testing.B) {
	logger := log.NewNopLogger()
	config := NewConfigFromFlags()

	dir, _ := ioutil.TempDir("", "goleveldb-app-invariant-bench")
	db, _ := sdk.NewLevelDB("simulation", dir)

	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()

	app := NewSimApp(logger, db, nil, true, 0)

	// 2. Run parameterized simulation (w/o invariants)
	_, params, simErr := simulation.SimulateFromSeed(
		b, ioutil.Discard, app.BaseApp, AppStateFn,
		testAndRunTxs(app, config), []sdk.Invariant{},
		app.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		fmt.Println("Exporting app state...")
		appState, _, err := app.ExportAppStateAndValidators(false, nil)
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}
		err = ioutil.WriteFile(config.ExportStatePath, []byte(appState), 0644)
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if config.ExportParamsPath != "" {
		fmt.Println("Exporting simulation params...")
		paramsBz, err := json.MarshalIndent(params, "", " ")
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}

		err = ioutil.WriteFile(config.ExportParamsPath, paramsBz, 0644)
		if err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}

	if simErr != nil {
		fmt.Println(simErr)
		b.FailNow()
	}

	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})

	// 3. Benchmark each invariant separately
	//
	// NOTE: We use the crisis keeper as it has all the invariants registered with
	// their respective metadata which makes it useful for testing/benchmarking.
	for _, cr := range app.CrisisKeeper.Routes() {
		b.Run(fmt.Sprintf("%s/%s", cr.ModuleName, cr.Route), func(b *testing.B) {
			if res, stop := cr.Invar(ctx); stop {
				fmt.Printf("broken invariant at block %d of %d\n%s", ctx.BlockHeight()-1, config.NumBlocks, res)
				b.FailNow()
			}
		})
	}
}
