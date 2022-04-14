package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	epochstoragetypes "github.com/lavanet/lava/x/epochstorage/types"
)

type SpecKeeper interface {
	// Methods imported from spec should be defined here
}

type EpochstorageKeeper interface {
	// Methods imported from epochStorage should be defined here
	// Methods imported from bank should be defined here
	GetEpochStart(ctx sdk.Context) uint64
	GetEarliestEpochStart(ctx sdk.Context) uint64
	UnstakeHoldBlocks(ctx sdk.Context) (res uint64)
	IsEpochStart(ctx sdk.Context) (res bool)
	BlocksToSave(ctx sdk.Context) (res uint64)
	PopUnstakeEntries(ctx sdk.Context, storageType string, block uint64) (value []epochstoragetypes.StakeEntry)
	AppendUnstakeEntry(ctx sdk.Context, storageType string, stakeEntry epochstoragetypes.StakeEntry)
	GetStakeStorageUnstake(ctx sdk.Context, storageType string) (epochstoragetypes.StakeStorage, bool)
	ModifyStakeEntry(ctx sdk.Context, storageType string, chainID string, stakeEntry epochstoragetypes.StakeEntry, removeIndex uint64)
	AppendStakeEntry(ctx sdk.Context, storageType string, chainID string, stakeEntry epochstoragetypes.StakeEntry)
	RemoveStakeEntry(ctx sdk.Context, storageType string, chainID string, idx uint64)
	StakeEntryByAddress(ctx sdk.Context, storageType string, chainID string, address sdk.AccAddress) (value epochstoragetypes.StakeEntry, found bool, index uint64)
	GetStakeStorageCurrent(ctx sdk.Context, storageType string, chainID string) (epochstoragetypes.StakeStorage, bool)
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error
	// Methods imported from bank should be defined here
}
