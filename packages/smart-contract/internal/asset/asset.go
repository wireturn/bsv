package asset

import (
	"context"
	"fmt"
	"time"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/smart-contract/internal/platform/db"
	"github.com/tokenized/smart-contract/internal/platform/node"
	"github.com/tokenized/smart-contract/internal/platform/state"
	"github.com/tokenized/specification/dist/golang/actions"
	"github.com/tokenized/specification/dist/golang/assets"
	"github.com/tokenized/specification/dist/golang/protocol"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

var (
	// ErrNotFound abstracts the standard not found error.
	ErrNotFound = errors.New("Asset not found")
)

// Retrieve gets the specified asset from the database.
func Retrieve(ctx context.Context, dbConn *db.DB, contractAddress bitcoin.RawAddress,
	assetCode *bitcoin.Hash20) (*state.Asset, error) {

	ctx, span := trace.StartSpan(ctx, "internal.asset.Retrieve")
	defer span.End()

	// Find asset in storage
	a, err := Fetch(ctx, dbConn, contractAddress, assetCode)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Create the asset
func Create(ctx context.Context, dbConn *db.DB, contractAddress bitcoin.RawAddress,
	assetCode *bitcoin.Hash20, nu *NewAsset, now protocol.Timestamp) error {

	ctx, span := trace.StartSpan(ctx, "internal.asset.Create")
	defer span.End()

	// Set up asset
	var a state.Asset

	// Get current state
	err := node.Convert(ctx, &nu, &a)
	if err != nil {
		return err
	}

	a.Code = assetCode
	a.Revision = 0
	a.CreatedAt = now
	a.UpdatedAt = now

	// a.Holdings = make(map[protocol.PublicKeyHash]*state.Holding)
	// a.Holdings[nu.AdministrationPKH] = &state.Holding{
	// 	PKH:              nu.AdministrationPKH,
	// 	PendingBalance:   nu.TokenQty,
	// 	FinalizedBalance: nu.TokenQty,
	// 	CreatedAt:        a.CreatedAt,
	// 	UpdatedAt:        a.UpdatedAt,
	// }

	if a.AssetPayload == nil {
		a.AssetPayload = []byte{}
	}

	return Save(ctx, dbConn, contractAddress, &a)
}

// Update the asset
func Update(ctx context.Context, dbConn *db.DB, contractAddress bitcoin.RawAddress,
	assetCode *bitcoin.Hash20, upd *UpdateAsset, now protocol.Timestamp) error {
	ctx, span := trace.StartSpan(ctx, "internal.asset.Update")
	defer span.End()

	// Find asset
	a, err := Fetch(ctx, dbConn, contractAddress, assetCode)
	if err != nil {
		return ErrNotFound
	}

	// Update fields
	if upd.Revision != nil {
		a.Revision = *upd.Revision
	}
	if upd.Timestamp != nil {
		a.Timestamp = *upd.Timestamp
	}

	if upd.AssetPermissions != nil {
		a.AssetPermissions = *upd.AssetPermissions
	}
	if upd.TradeRestrictions != nil {
		a.TradeRestrictions = *upd.TradeRestrictions
	}
	if upd.EnforcementOrdersPermitted != nil {
		a.EnforcementOrdersPermitted = *upd.EnforcementOrdersPermitted
	}
	if upd.VoteMultiplier != nil {
		a.VoteMultiplier = *upd.VoteMultiplier
	}
	if upd.AdministrationProposal != nil {
		a.AdministrationProposal = *upd.AdministrationProposal
	}
	if upd.HolderProposal != nil {
		a.HolderProposal = *upd.HolderProposal
	}
	if upd.AssetModificationGovernance != nil {
		a.AssetModificationGovernance = *upd.AssetModificationGovernance
	}
	if upd.AuthorizedTokenQty != nil {
		a.AuthorizedTokenQty = *upd.AuthorizedTokenQty
	}
	if upd.AssetPayload != nil {
		a.AssetPayload = *upd.AssetPayload
	}
	if upd.FreezePeriod != nil {
		a.FreezePeriod = *upd.FreezePeriod
	}

	a.UpdatedAt = now

	return Save(ctx, dbConn, contractAddress, a)
}

// ValidateVoting returns an error if voting is not allowed.
func ValidateVoting(ctx context.Context, as *state.Asset, initiatorType uint32,
	votingSystem *actions.VotingSystemField) error {

	switch initiatorType {
	case 0: // Administration
		if !as.AdministrationProposal {
			return errors.New("Administration proposals not allowed")
		}
	case 1: // Holder
		if !as.HolderProposal {
			return errors.New("Holder proposals not allowed")
		}
	}

	return nil
}

func timeString(t uint64) string {
	return time.Unix(int64(t)/1000000000, 0).String()
}

// IsTransferable returns an error if the asset is non-transferable.
func IsTransferable(ctx context.Context, as *state.Asset, now protocol.Timestamp) error {
	if as.FreezePeriod.Nano() > now.Nano() {
		return node.NewError(actions.RejectionsAssetFrozen,
			fmt.Sprintf("Asset frozen until %s", as.FreezePeriod.String()))
	}

	assetData, err := assets.Deserialize([]byte(as.AssetType), as.AssetPayload)
	if err != nil {
		return node.NewError(actions.RejectionsMsgMalformed, err.Error())
	}

	switch data := assetData.(type) {
	case *assets.Membership:
		if data.ExpirationTimestamp != 0 && data.ExpirationTimestamp < now.Nano() {
			return node.NewError(actions.RejectionsAssetNotPermitted,
				fmt.Sprintf("Membership expired at %s", timeString(data.ExpirationTimestamp)))
		}

	case *assets.ShareCommon:

	case *assets.CasinoChip:
		if data.ExpirationTimestamp != 0 && data.ExpirationTimestamp < now.Nano() {
			return node.NewError(actions.RejectionsAssetNotPermitted,
				fmt.Sprintf("CasinoChip expired at %s", timeString(data.ExpirationTimestamp)))
		}

	case *assets.Coupon:
		if data.ExpirationTimestamp != 0 && data.ExpirationTimestamp < now.Nano() {
			return node.NewError(actions.RejectionsAssetNotPermitted,
				fmt.Sprintf("Coupon expired at %s", timeString(data.ExpirationTimestamp)))
		}

	case *assets.LoyaltyPoints:
		if data.ExpirationTimestamp != 0 && data.ExpirationTimestamp < now.Nano() {
			return node.NewError(actions.RejectionsAssetNotPermitted,
				fmt.Sprintf("LoyaltyPoints expired at %s", timeString(data.ExpirationTimestamp)))
		}

	case *assets.TicketAdmission:
		if data.EventEndTimestamp != 0 && data.EventEndTimestamp < now.Nano() {
			return node.NewError(actions.RejectionsAssetNotPermitted,
				fmt.Sprintf("TicketAdmission expired at %s", timeString(data.EventEndTimestamp)))
		}
	}

	return nil
}
