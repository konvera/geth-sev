/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package es

import (
	"github.com/konvera/geth-sev/constellation/attestation"
	"github.com/konvera/geth-sev/constellation/attestation/gcp"
	"github.com/konvera/geth-sev/constellation/attestation/variant"
	"github.com/konvera/geth-sev/constellation/attestation/vtpm"
	tpmclient "github.com/google/go-tpm-tools/client"
)

// Issuer for GCP confidential VM attestation.
type Issuer struct {
	variant.GCPSEVES
	*vtpm.Issuer
}

// NewIssuer initializes a new GCP Issuer.
func NewIssuer(log attestation.Logger) *Issuer {
	return &Issuer{
		Issuer: vtpm.NewIssuer(
			vtpm.OpenVTPM,
			tpmclient.GceAttestationKeyRSA,
			gcp.GCEInstanceInfo(gcp.MetadataClient{}),
			log,
		),
	}
}
