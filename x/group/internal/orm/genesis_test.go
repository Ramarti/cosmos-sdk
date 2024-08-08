package orm

import (
	"testing"

	"github.com/stretchr/testify/require"

	coretesting "cosmossdk.io/core/testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
)

func TestImportExportTableData(t *testing.T) {
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	table, err := NewAutoUInt64Table(AutoUInt64TablePrefix, AutoUInt64TableSeqPrefix, &testdata.TableModel{}, cdc, address.NewBech32Codec("cosmos"))
	require.NoError(t, err)

	testCtx := testutil.DefaultContextWithDB(t, "test")
	store := coretesting.KVStoreService(testCtx.Ctx, "test").OpenKVStore(testCtx.Ctx)

	tms := []*testdata.TableModel{
		{
			Id:       1,
			Name:     "my test 1",
			Number:   123,
			Metadata: []byte("metadata 1"),
		},
		{
			Id:       2,
			Name:     "my test 2",
			Number:   456,
			Metadata: []byte("metadata 2"),
		},
	}

	err = table.Import(store, tms, 2)
	require.NoError(t, err)

	for _, g := range tms {
		var loaded testdata.TableModel
		_, err := table.GetOne(store, g.Id, &loaded)
		require.NoError(t, err)

		require.Equal(t, g, &loaded)
	}

	var exported []*testdata.TableModel
	seq, err := table.Export(store, &exported)
	require.NoError(t, err)
	require.Equal(t, seq, uint64(2))

	for i, g := range exported {
		require.Equal(t, g, tms[i])
	}
}
