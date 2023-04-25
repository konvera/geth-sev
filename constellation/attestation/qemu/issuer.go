/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package qemu

import (
	"context"
	"io"

	"github.com/konvera/geth-sev/constellation/attestation/vtpm"
	"github.com/konvera/geth-sev/constellation/variant"
	tpmclient "github.com/google/go-tpm-tools/client"
)

// Issuer for qemu TPM attestation.
type Issuer struct {
	variant.QEMUVTPM
	*vtpm.Issuer
}

// NewIssuer initializes a new QEMU Issuer.
func NewIssuer(log vtpm.AttestationLogger) *Issuer {
	return &Issuer{
		Issuer: vtpm.NewIssuer(
			vtpm.OpenVTPM,
			tpmclient.AttestationKeyRSA,
			func(context.Context, io.ReadWriteCloser, []byte) ([]byte, error) { return nil, nil },
			log,
		),
	}
}
