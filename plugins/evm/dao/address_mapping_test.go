package dao

import (
	"context"
	"errors"
	"testing"

	"github.com/itering/subscan/util/ss58"
)

// Public third-party sample from subscan-issue-tracker#164.
func TestReviveMappedAccountResolution(t *testing.T) {
	const h160 = "0x2ee6ace1cdc464a61d3f3910a6853ab5d093a107"
	native := ss58.Decode("12u9Ha4PxyyQPvJgq3BghnqNXDwLqTnnJFuXV7aZQoiregT2")
	if len(native) != 64 {
		t.Fatalf("invalid sample AccountId: %q", native)
	}
	failure := errors.New("RPC unavailable")
	tests := []struct {
		name, mapped, want string
		lookupErr          error
		wantErr            bool
	}{
		{name: "registered public sample", mapped: "0x" + native, want: native},
		{name: "unmapped Ethereum fallback", want: reviveAccount(h160)},
		{name: "RPC error must not become fallback", lookupErr: failure, wantErr: true},
		{name: "malformed mapping must not become fallback", mapped: h160, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got, err := resolveH160Account(context.Background(), h160, false, true, func(ctx context.Context, input string) (string, error) {
				calls++
				if input != h160 {
					t.Fatalf("lookup key: %s", input)
				}
				return tt.mapped, tt.lookupErr
			})
			if calls != 1 || (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("got %q, err %v, calls %d; want %q, error %v", got, err, calls, tt.want, tt.wantErr)
			}
			if tt.lookupErr != nil && !errors.Is(err, failure) {
				t.Fatalf("lost lookup error: %v", err)
			}
		})
	}
	// The old padded lookup misses the indexed native balance; mapping finds it.
	balances := map[string]int64{native: 53700000000}
	resolved, err := resolveH160Account(context.Background(), h160, false, true, func(context.Context, string) (string, error) { return native, nil })
	if err != nil || balances[resolved] == 0 || balances[reviveAccount(h160)] != 0 {
		t.Fatalf("mapped balance lookup failed: %s, %v", resolved, err)
	}
}

func TestAccountResolutionSkipsMappingForOtherNetworks(t *testing.T) {
	lookup := func(context.Context, string) (string, error) { t.Fatal("unexpected Revive lookup"); return "", nil }
	const h160 = "0x2ee6ace1cdc464a61d3f3910a6853ab5d093a107"
	got, err := resolveH160Account(context.Background(), h160, true, true, lookup)
	if err != nil || got != h160 {
		t.Fatalf("EVM chain: %s %v", got, err)
	}
	got, err = resolveH160Account(context.Background(), h160, false, false, lookup)
	if err != nil || got == "" || got == reviveAccount(h160) {
		t.Fatalf("Frontier: %s %v", got, err)
	}
	got, err = resolveH160Account(context.Background(), "invalid", false, true, lookup)
	if err != nil || got != "" {
		t.Fatalf("invalid address: %s %v", got, err)
	}
}

func TestReviveLookupHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := readReviveOriginalAccount(ctx, "0x2ee6ace1cdc464a61d3f3910a6853ab5d093a107")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
}

func TestReviveLookupNormalizesH160(t *testing.T) {
	const input = "0x2EE6ACe1CDC464a61D3F3910a6853Ab5D093A107"
	const want = "0x2ee6ace1cdc464a61d3f3910a6853ab5d093a107"
	got, err := resolveH160Account(context.Background(), input, false, true, func(_ context.Context, key string) (string, error) {
		if key != want {
			t.Fatalf("noncanonical lookup: %s", key)
		}
		return "", nil
	})
	if err != nil || got != reviveAccount(want) {
		t.Fatalf("got %s %v", got, err)
	}
}
