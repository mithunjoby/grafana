package versionpolicy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana/pkg/storage/unified/resource/kv"
)

// ReadGlobalLayer reads the per-cluster version policy from the KV store: one entry per group, keyed by
// group name. A nil store yields no layer (falls back to ini). Any error is returned rather than partially
// applied, so the caller keeps its last-known layer instead of silently dropping a group's cap.
func ReadGlobalLayer(ctx context.Context, store kv.KV) (map[string]VersionPolicy, error) {
	if store == nil {
		return nil, nil
	}

	var keys []string
	for key, err := range store.Keys(ctx, kv.VersionPolicySection, kv.ListOptions{}) {
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return map[string]VersionPolicy{}, nil
	}

	layer := make(map[string]VersionPolicy, len(keys))
	for kve, err := range store.BatchGet(ctx, kv.VersionPolicySection, keys) {
		if err != nil {
			return nil, err
		}
		var p VersionPolicy
		decodeErr := json.NewDecoder(kve.Value).Decode(&p)
		_ = kve.Value.Close()
		if decodeErr != nil {
			return nil, fmt.Errorf("version policy: decoding KV value for group %q: %w", kve.Key, decodeErr)
		}
		if kve.Key == "" {
			continue // defensive: no group to key on
		}
		layer[kve.Key] = p
	}
	return layer, nil
}
