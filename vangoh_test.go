package vangoh

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"hash"
	"testing"
)

type testProvider struct {
	promptErr bool
	key       []byte
	secret    []byte
}

func (tp *testProvider) GetSecret(key []byte) ([]byte, error) {
	if tp.promptErr {
		return nil, errors.New("testing error")
	}
	if !bytes.Equal(tp.key, key) {
		return nil, nil
	}
	return tp.secret, nil
}

var tp1 = &testProvider{
	promptErr: false,
	key:       []byte("testIDOne"),
	secret:    []byte("secretOne"),
}

var tp2 = &testProvider{
	promptErr: false,
	key:       []byte("testIDTwo"),
	secret:    []byte("secretTwo"),
}

var tpErr = &testProvider{
	promptErr: true,
	key:       []byte("testIDErr"),
	secret:    []byte("secretErr"),
}

func TestNew(t *testing.T) {
	vg := New()

	if vg.includedHeaders == nil {
		t.Error("includeHeaders not properly intialized")
	}
	if vg.keyProviders == nil {
		t.Error("keyProviders not properly intialized")
	}
	if vg.singleProvider {
		t.Error("default constructor should not create a single provider instance")
	}
	if !checkAlgorithm(vg, crypto.SHA256.New) {
		t.Error("default constructor should instantiate the algorithm to SHA256")
	}
}

func TestNewSingleProvider(t *testing.T) {
	vg := NewSingleProvider(tp1)

	if vg.includedHeaders == nil {
		t.Error("includeHeaders not properly intialized")
	}
	if vg.keyProviders == nil {
		t.Error("keyProviders not properly intialized")
	}
	if !vg.singleProvider {
		t.Error("singleProvider constructor should create a single provider instance")
	}
}

func TestAddProvider(t *testing.T) {
	vg := New()

	if len(vg.keyProviders) != 0 {
		t.Error("Wrong number of key providers in the Vangoh instance")
	}

	err := vg.AddProvider("test", tp1)
	if err != nil {
		t.Error("Should not have encountered error when adding a new provider")
	}

	if len(vg.keyProviders) != 1 {
		t.Error("Wrong number of key providers in the Vangoh instance")
	}

	err = vg.AddProvider("test", tp2)
	if err == nil {
		t.Error("Should error when trying to add multiple providers for same org tag")
	}

	if len(vg.keyProviders) != 1 {
		t.Error("Wrong number of key providers in the Vangoh instance")
	}

	err = vg.AddProvider("notTest", tp2)
	if err != nil {
		t.Error("Should not error when trying to add multiple providers for different org tags")
	}

	if len(vg.keyProviders) != 2 {
		t.Error("Wrong number of key providers in the Vangoh instance")
	}

	spvg := NewSingleProvider(tp1)

	if len(spvg.keyProviders) != 1 {
		t.Error("Wrong number of key providers in the Vangoh instance")
	}

	err = spvg.AddProvider("test", tp2)
	if err == nil {
		t.Error("Should error when trying to add second provider to single provider instance")
	}

	if len(spvg.keyProviders) != 1 {
		t.Error("Wrong number of key providers in the Vangoh instance")
	}
}

func TestAlgorithm(t *testing.T) {
	vg := New()

	if !checkAlgorithm(vg, crypto.SHA256.New) {
		t.Error("default constructor should instantiate the algorithm to SHA256")
	}

	vg.SetAlgorithm(crypto.SHA1.New)
	if !checkAlgorithm(vg, crypto.SHA1.New) {
		t.Error("Algorithm not correctly updated with SetAlgorithm method")
	}
}

func checkAlgorithm(vg *Vangoh, algo func() hash.Hash) bool {
	vga := fmt.Sprintf("%T", vg.algorithm())
	toCheck := fmt.Sprintf("%T", algo())

	return vga == toCheck
}
