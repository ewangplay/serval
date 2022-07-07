package client_test

import (
	"encoding/json"
	"testing"

	"github.com/ewangplay/serval/io"
	sdk "github.com/ewangplay/serval/sdk/go"
)

const (
	addr = "localhost:8099"
	did  = "did:example:fafdecaa29934fde9dcc5adaea8ea82b"
	ddo  = `
	{
		"@context" : "https://www.w3.org/ns/did/v1",
		"id" : "did:example:fafdecaa29934fde9dcc5adaea8ea82b",
		"version" : 1,
		"publicKey" : [
		{
			"id" : "did:example:fafdecaa29934fde9dcc5adaea8ea82b#keys-1",
			"type" : "Ed25519",
			"publicKeyHex" : "4cd5d192e33f390d7f2a5cc948103a080c52f42fefe2d4d334f86e7ac78e0938"
		},
		{
			"id" : "did:example:fafdecaa29934fde9dcc5adaea8ea82b#keys-2",
			"type" : "Ed25519",
			"publicKeyHex" : "e1df3e6e58d51ba0217137224d6daef2a5d1a5790b6537d6f9830d59639e0826"
		}
		],
		"controller" : "did:example:fafdecaa29934fde9dcc5adaea8ea82b",
		"authentication" : [
		"did:example:fafdecaa29934fde9dcc5adaea8ea82b#keys-1"
		],
		"recovery" : [
		"did:example:fafdecaa29934fde9dcc5adaea8ea82b#keys-2"
		],
		"service" : null,
		"proof" : {
		"type" : "Ed25519",
		"creator" : "did:example:fafdecaa29934fde9dcc5adaea8ea82b#keys-1",
		"signatureValue" : "ak5hxLAke9UQ1ExSyYVu8r1GjAzlScP1DHHefQ9187c3MCBsoR8My5X5ex1IBHxexDYj9yBetIv24yu+Dh2sBg=="
		},
		"created" : "2022-07-06T13:52:26.793567+08:00",
		"updated" : "2022-07-06T13:52:26.793567+08:00"
	}
	`
)

func newClient(t *testing.T) *sdk.Client {
	c, err := sdk.NewClient(addr)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Ping()
	if err != nil {
		t.Skip()
	}
	return c
}

func TestCreateDid(t *testing.T) {
	var err error
	c := newClient(t)

	t.Run("Ed25519", func(t *testing.T) {
		var document io.DDO
		err = json.Unmarshal([]byte(ddo), &document)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Document: ", document)

		req := &io.CreateDidReq{
			Did:      did,
			Document: document,
		}
		err = c.CreateDid(req)
		if err != nil {
			t.Fatal(err)
		}
	})

}

func TestResolveDid(t *testing.T) {
	t.Skip()
}

func TestRevokeDid(t *testing.T) {
	t.Skip()
}
