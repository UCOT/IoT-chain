// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package dbft

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"bytes"
	// "fmt"
	"github.com/ethereum/go-ethereum/rlp"
)

// API is a user facing RPC API to allow controlling the signer
// mechanisms of dbft scheme.
type API struct {
	chain consensus.ChainReader
	dbft  *Dbft
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*SnapShot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dbft.SnapShot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*SnapShot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dbft.SnapShot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSigners retrieves the list of authorized signers at the specified block.
// Here the signer is the node who broadcasts the newblock
func (api *API) GetSigners(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.dbft.SnapShot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.signers(), nil
}

// GetSignersAtHash retrieves the state snapshot at a given block.
// Here the signer is the node who broadcasts the newblock
func (api *API) GetSignersAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.dbft.SnapShot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.signers(), nil
}

type GroupSigFormatter struct {
	Sig    string `json:"sigs"`
	IIndex uint64 `json:"i"`
}

// GetGroupSigsAtNumber retrieves the group signature at a given block.
func (api *API) GetGroupSigsAtNumber(number *rpc.BlockNumber) ([]GroupSigFormatter, error) { //***
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	sigsRLP, _ := api.chain.GetGroupSigRLP(header.Hash())
	sigs := new([]types.GroupSignature)
	if err := rlp.Decode(bytes.NewReader(sigsRLP), sigs); err != nil {
		return nil, err
	}
	sigstr := make([]GroupSigFormatter, len(*sigs))
	for i, v := range *sigs {
		sigstr[i].Sig = common.ToHex(v.Sig)
		sigstr[i].IIndex = v.IIndex
	}
	return sigstr, nil
	// return *sigs, nil
}

// GetGroupSigsAtHash retrieves the group signature at a given block.
func (api *API) GetGroupSigsAtHash(blockHash common.Hash) ([]GroupSigFormatter, error) { //***
	header := api.chain.GetHeaderByHash(blockHash)
	if header == nil {
		return nil, errUnknownBlock
	}
	sigsRLP, _ := api.chain.GetGroupSigRLP(blockHash)
	sigs := new([]types.GroupSignature)
	if err := rlp.Decode(bytes.NewReader(sigsRLP), sigs); err != nil {
		return nil, err
	}
	sigstr := make([]GroupSigFormatter, len(*sigs))
	for i, v := range *sigs {
		sigstr[i].Sig = common.ToHex(v.Sig)
		sigstr[i].IIndex = v.IIndex
	}
	return sigstr, nil
	// return *sigs, nil
}
