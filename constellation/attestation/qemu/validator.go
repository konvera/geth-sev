/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package qemu

import (
	"context"
	"crypto"

	"github.com/konvera/geth-sev/constellation/attestation/vtpm"
	"github.com/konvera/geth-sev/constellation/config"
	"github.com/konvera/geth-sev/constellation/variant"
	"github.com/google/go-tpm-tools/proto/attest"
	"github.com/google/go-tpm/tpm2"
)

// Validator for QEMU VM attestation.
type Validator struct {
	variant.QEMUVTPM
	*vtpm.Validator
}

// NewValidator initializes a new QEMU validator with the provided PCR values.
func NewValidator(cfg config.QEMUVTPM, log vtpm.AttestationLogger) *Validator {
	return &Validator{
		Validator: vtpm.NewValidator(
			cfg.Measurements,
			unconditionalTrust,
			func(vtpm.AttestationDocument, *attest.MachineState) error { return nil },
			log,
		),
	}
}

// unconditionalTrust returns the given public key as the trusted attestation key.
func unconditionalTrust(_ context.Context, attDoc vtpm.AttestationDocument, _ []byte) (crypto.PublicKey, error) {
	pubArea, err := tpm2.DecodePublic(attDoc.Attestation.AkPub)
	if err != nil {
		return nil, err
	}
	return pubArea.Key()
}
