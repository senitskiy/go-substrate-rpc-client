// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	// "rpc/signature"
	// "rpc/types"

	// "rpc/types/codec"

	"github.com/JFJun/go-substrate-crypto/ss58"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/ed25519"
)

type KeyringPair struct {
	// URI is the derivation path for the private key in subkey
	URI string
	// Address is an SS58 address
	Address string
	// PublicKey
	PublicKey []byte
}

func KeyringPairFromSecreted25519(seedOrPhrase string, network uint8) (KeyringPair, error) {
	scheme := ed25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, seedOrPhrase)
	if err != nil {
		return KeyringPair{}, err
	}

	ss58Address := kyr.SS58Address(uint16(network))
	if err != nil {
		return KeyringPair{}, err
	}

	var pk = kyr.Public()

	return KeyringPair{
		URI:       seedOrPhrase,
		Address:   ss58Address,
		PublicKey: pk,
	}, nil
}

func main() {

	// The following example shows how to instantiate a Substrate API and use it to connect to a node
	url := "wss://westend.api.onfinality.io/public-ws" //"wss://kusama-rpc.polkadot.io"
	// "wss://westend.api.onfinality.io/public-ws" //"wss://westend-rpc.polkadot.io"

	api, err := gsrpc.NewSubstrateAPI(url) //wss://westend.api.onfinality.io/public-ws")
	if err != nil {
		panic(err)
	}

	chain, err := api.RPC.System.Chain()
	if err != nil {
		panic(err)
	}
	nodeName, err := api.RPC.System.Name()
	if err != nil {
		panic(err)
	}
	nodeVersion, err := api.RPC.System.Version()
	if err != nil {
		panic(err)
	}

	fmt.Printf("You are connected to chain %v using %v v%v\n", chain, nodeName, nodeVersion)

	var pub []byte

	pub, err = ss58.DecodeToPub("5Dc3EjbRGfafWUwu3igkHqyrTvzRZ5hcixrWCscRda5GL1CB")
	if err != nil {
		fmt.Errorf("ss58 decode address error: %v", err)
	}
	fmt.Println("pub", pub)

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	var storage []byte

	storage, err = types.CreateStorageKey(meta, "System", "Account", pub, nil)
	if err != nil {
		fmt.Errorf("create System.Account storage error: %v", err)
	}
	fmt.Println("storage", storage)

	// var accountInfo1 types.AccountInfo
	// var ok1 bool

	// var accountInfoProviders expand.AccountInfoWithProviders
	// ok1, err = api.RPC.State.GetStorageLatest(storage, &accountInfoProviders)
	// if err != nil || !ok1 {
	// 	fmt.Errorf("get account info error: %v", err)
	// }
	// accountInfo1.Nonce = accountInfoProviders.Nonce
	// // accountInfo.Refcount = types.U8(accountInfoProviders.Consumers)
	// accountInfo1.Data.Free = accountInfoProviders.Data.Free
	// accountInfo1.Data.FreeFrozen = accountInfoProviders.Data.FreeFrozen
	// accountInfo1.Data.MiscFrozen = accountInfoProviders.Data.MiscFrozen
	// accountInfo1.Data.Reserved = accountInfoProviders.Data.Reserved

	// fmt.Println("accountInfoProviders: ", accountInfoProviders.Data.Reserved)
	// fmt.Println("ok: ", ok1)

	// accountInfoProviders.Data.Free

	sub1, err := api.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer sub1.Unsubscribe()

	count := 0

	for {
		head := <-sub1.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		count++

		if count == 1 {
			sub1.Unsubscribe()
			break
		}
	}

	// var accountInfo types.AccountInfo
	// ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	// if err != nil || !ok {
	// 	panic(err)
	// }

	// previous := accountInfo.Data.Free
	// fmt.Printf("has a balance of %v\n", previous)

	// acc, err := api.GetAccountInfo("3nbzmMR1aufpfRcmueqatRpBzStubHKxporVw5xz1KNy1KnY")

	// Create a call, transferring 12345 units to Bob //NewAddressFromHexAccountID 0xe143f23803ac50e8f6f8e62695d1ce9e4e1d68aa36c1cd2cfd15340213f3423e
	bob, err := types.NewMultiAddressFromHexAccountID("0xe143f23803ac50e8f6f8e62695d1ce9e4e1d68aa36c1cd2cfd15340213f3423e") //0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")

	if err != nil {
		panic(err)
	}

	fmt.Println("bob.AsID:", bob.AsID)
	// fmt.Println("bob.AsID:", bob.AsID.MarshalJSON())
	// fmt.Println("bob.AsID:", bob.AsID.Equal())
	fmt.Println("bob.AsID.ToBytes:", bob.AsID.ToBytes())
	fmt.Println("bob.AsID.ToHexString:", bob.AsID.ToHexString())
	// fmt.Println("bob.AsID:", bob.AsID.UnmarshalJSON())
	fmt.Println("bob.IsID:", bob.IsID)
	// fmt.Println("bob.AsAccountID:", bob.AsAddress20)
	// fmt.Println("bob.AsAccountID.ToHexString:", bob.AsAddress32)
	// fmt.Println("bob.AsID:", bob.AsID)
	// fmt.Println("bob.AsAccountID.ToHexString:", bob.AsIndex)
	// fmt.Println("bob.AsAccountID:", bob.AsRaw)
	// fmt.Println("bob.AsAccountID.ToHexString:", bob.IsAddress20)
	// fmt.Println("bob.AsAccountID:", bob.IsAddress32)

	// fmt.Println("bob.AsAccountID:", bob.IsIndex)
	// fmt.Println("bob.AsAccountID.ToHexString:", bob.IsRaw)
	// fmt.Println("bob.AsAccountID:", bob.)
	// fmt.Println("bob.AsAccountID.ToHexString:", bob.AsAddress32)

	// fmt.Println("bob.AsAccountID.ToHexString:", bob. .ToHexString())

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	// return KeyringPair{
	// 	URI:       seedOrPhrase,
	// 	Address:   ss58Address,
	// 	PublicKey: pk,
	// }, nil

	// Get the nonce for Alice
	seedPhrase := "case attend monitor normal record garbage raven code teach news moon vanish"
	johnKey, err := signature.KeyringPairFromSecret(seedPhrase, 42)
	if err != nil {
		panic(err)
	}

	johnKey2, err := KeyringPairFromSecreted25519(seedPhrase, 42)
	if err != nil {
		panic(err)
	}

	var john types.MultiAddress
	// var johnAcc types.AccountID
	// johnAcc = newKey.PublicKey [AccountIDLen]byte
	// johnAcc, err := types.NewAccountID(newKey.PublicKey)
	if err != nil {
		panic(err)
	}

	a := types.AccountID{}
	// copy(a[:], bob.AsID.ToBytes())
	copy(a[:], johnKey.PublicKey)

	// var jfds types.AccountID{}
	// jfds = johnAcc
	john.AsID = a
	john.IsID = true

	fmt.Println("john:", johnKey)
	fmt.Println("john2:", johnKey2)

	fmt.Println("newKey.Address:", johnKey.Address)
	fmt.Println("newKey.PublicKey:", johnKey.PublicKey)
	fmt.Println("newKey.URI:", johnKey.URI)
	// fmt.Printf("newKey: %v\n", newKey)
	// fmt.Println("newKey.signature:", newKey.)

	key1, err := types.CreateStorageKey(meta, "System", "Account", johnKey.PublicKey)
	if err != nil {
		panic(err)
	}

	var accountInfo1 types.AccountInfo
	ok1, err := api.RPC.State.GetStorageLatest(key1, &accountInfo1)
	if err != nil || !ok1 {
		panic(err)
	}
	fmt.Println("NEW accountInfo.Data.Free: ", accountInfo1.Data.Free)

	//  429 496 729 600 000 000 000
	// 1373898004 051 167 019 008
	newKey, err := types.CreateStorageKey(meta, "System", "Account", johnKey.PublicKey) //  signature.TestKeyringPairAlice.PublicKey) //bob.
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(newKey, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	fmt.Println("accountInfo.Data.Free: ", accountInfo.Data.Free)

	amount := types.NewUCompactFromUInt(600000000000)

	// fmt.Println("c:", c)
	fmt.Println("bob:", bob)
	fmt.Println("newKey:", newKey)

	// signature.Verify()
	ivan, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")

	c, err := types.NewCall(meta, "Balances.transfer", ivan, amount) //bob
	if err != nil {
		panic(err)
	}
	fmt.Println("c:", c)
	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	// fmt.Println("ext:", ext)

	// acc, err := api.GetAccountInfo("5DskLcrX2digY5Qy2E5qVWhVYCYKkhNK8N7QvUAkF27Z9kGo") accountInfo1

	// key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey)
	// if err != nil {
	// 	panic(err)
	// }

	nonce := uint32(accountInfo.Nonce)
	fmt.Println("nonce:", nonce)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(100),
		TransactionVersion: rv.TransactionVersion,
	}
	fmt.Println("o: ", o)
	fmt.Println("amount: ", amount)

	fmt.Printf("Sending %v from %#x to %#x with nonce %v\n", amount, john.AsID, ivan.AsID, nonce)

	// Sign the transaction using Alice's default account
	// err = ext.Sign(signature.TestKeyringPairAlice, o)
	// if err != nil {
	// 	panic(err)
	// }

	// err = ext.Sign(signature.TestKeyringPairAlice, o)
	err = ext.Sign(johnKey, o)
	if err != nil {
		panic(err)
	}

	// Send the extrinsic
	_, err = api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Balance transferred from Alice to Bob: %v\n", amount)
	// Output: Balance transferred from Alice to Bob: 100000000000000

	// d, _ := json.Marshal(acc)
	// fmt.Println(string(d))
	// fmt.Println(acc.Data.Free.String())

	// Do the transfer and track the actual status
	// sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	// if err != nil {
	// 	panic(err)
	// }
	// defer sub.Unsubscribe()

	// for {
	// 	status := <-sub.Chan()
	// 	fmt.Printf("Transaction status: %#v\n", status)

	// 	if status.IsInBlock {
	// 		fmt.Printf("Completed at block hash: %#x\n", status.AsInBlock)
	// 		return
	// 	}
	// }

}
