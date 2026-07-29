package versionpolicy

import (
	"context"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/storage/unified/resource/kv"
)

func seedKV(t *testing.T) kv.KV {
	t.Helper()
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return kv.NewBadgerKV(db)
}

func save(t *testing.T, store kv.KV, key, value string) {
	t.Helper()
	w, err := store.Save(context.Background(), kv.VersionPolicySection, key)
	require.NoError(t, err)
	_, err = w.Write([]byte(value))
	require.NoError(t, err)
	require.NoError(t, w.Close())
}

func TestReadGlobalLayer(t *testing.T) {
	t.Run("nil store yields no layer", func(t *testing.T) {
		layer, err := ReadGlobalLayer(context.Background(), nil)
		require.NoError(t, err)
		require.Nil(t, layer)
	})

	t.Run("empty section yields empty layer", func(t *testing.T) {
		layer, err := ReadGlobalLayer(context.Background(), seedKV(t))
		require.NoError(t, err)
		require.Empty(t, layer)
	})

	t.Run("reads one entry per group", func(t *testing.T) {
		store := seedKV(t)
		save(t, store, "dashboard.grafana.app", `{"preferredVersion":"v1beta1","maxAllowedVersion":"v1"}`)
		save(t, store, "folder.grafana.app", `{"preferredVersion":"v0alpha1"}`)

		layer, err := ReadGlobalLayer(context.Background(), store)
		require.NoError(t, err)
		require.Equal(t, VersionPolicy{PreferredVersion: "v1beta1", MaxAllowedVersion: "v1"}, layer["dashboard.grafana.app"])
		require.Equal(t, VersionPolicy{PreferredVersion: "v0alpha1"}, layer["folder.grafana.app"])
		require.Len(t, layer, 2)
	})

	t.Run("malformed value errors, so the caller keeps its last-known layer", func(t *testing.T) {
		store := seedKV(t)
		save(t, store, "dashboard.grafana.app", `{"maxAllowedVersion":"v1"}`)
		save(t, store, "bad.grafana.app", `{not json`)

		layer, err := ReadGlobalLayer(context.Background(), store)
		require.Error(t, err, "a malformed value must fail the read, not silently drop the group and lift its cap")
		require.Nil(t, layer)
	})
}
