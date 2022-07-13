package client_test

import (
	"encoding/json"
	"testing"

	"github.com/ewangplay/serval/io"
	sdk "github.com/ewangplay/serval/sdk/go"
)

const (
	addr    = "localhost:8099"
	did     = "did:example:a85c33daeda746f29c18d6386117d95d"
	ddoJson = `
	{
		"@context" : "https://www.w3.org/ns/did/v1",
		"id" : "did:example:a85c33daeda746f29c18d6386117d95d",
		"version" : 1,
		"publicKey" : [
		  {
			"id" : "did:example:a85c33daeda746f29c18d6386117d95d#keys-1",
			"type" : "ED25519",
			"publicKeyHex" : "6ded51fe5afec830ac4fac26cd832f5d44d9673ecedd8eda0fe0158ac5b4e283"
		  },
		  {
			"id" : "did:example:a85c33daeda746f29c18d6386117d95d#keys-2",
			"type" : "ED25519",
			"publicKeyHex" : "79df8dffeadd9dc0318a7a394f8dbd99062de48ddb1fa487eb13335da7e38e76"
		  }
		],
		"controller" : "did:example:a85c33daeda746f29c18d6386117d95d",
		"authentication" : [
		  "did:example:a85c33daeda746f29c18d6386117d95d#keys-1"
		],
		"recovery" : [
		  "did:example:a85c33daeda746f29c18d6386117d95d#keys-2"
		],
		"service" : null,
		"proof" : {
		  "type" : "ED25519",
		  "creator" : "did:example:a85c33daeda746f29c18d6386117d95d#keys-1",
		  "signatureValue" : "MXjoTf3r6WhMQRqsncpxZJ4kkM5760kEP0u7ml1FZMyrpuzrqnNwIcoTuWdQuS0L1hl/i2zWMiKo29dSeeDjDQ=="
		},
		"created" : "2022-07-13T15:17:27.352714+08:00",
		"updated" : "2022-07-13T15:17:27.352714+08:00"
	}
	`
	proofJson = `
	{
		"type" : "ED25519",
		"creator" : "did:example:a85c33daeda746f29c18d6386117d95d#keys-2",
		"signatureValue" : "plehewLPcgiCiOcE6oqV1xCdmkT+VTwhvwo3ywvvYrr+ArQCzNjb+A8a3RCaWwTxJmNXTWHLCkB0ChGkbiSbAQ=="
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
		err = json.Unmarshal([]byte(ddoJson), &document)
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

	t.Run("Ed25519", func(t *testing.T) {
		var d *io.DDO
		d, err = c.ResolveDid(did)
		if err != nil {
			t.Fatal(err)
		}
		if d == nil {
			t.Fatal("Revolve did failed:", did)
		}
	})

	t.Run("Ed25519", func(t *testing.T) {
		var proof io.Proof
		err = json.Unmarshal([]byte(proofJson), &proof)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Proof: ", proof)

		req := &io.RevokeDidReq{
			Did:   did,
			Proof: proof,
		}
		err = c.RevokeDid(req)
		if err != nil {
			t.Fatal(err)
		}
	})
}
