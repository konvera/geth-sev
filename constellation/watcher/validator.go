/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package watcher

import (
	"context"
	"encoding/asn1"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/konvera/geth-sev/constellation/atls"
	"github.com/konvera/geth-sev/constellation/attestation/choose"
	"github.com/konvera/geth-sev/constellation/attestation/idkeydigest"
	"github.com/konvera/geth-sev/constellation/attestation/measurements"
	"github.com/konvera/geth-sev/constellation/config"
	"github.com/konvera/geth-sev/constellation/constants"
	"github.com/konvera/geth-sev/constellation/file"
	"github.com/konvera/geth-sev/constellation/logger"
	"github.com/konvera/geth-sev/constellation/variant"
	"github.com/spf13/afero"
)

// Updatable implements an updatable atls.Validator.
type Updatable struct {
	log         *logger.Logger
	mux         sync.Mutex
	fileHandler file.Handler
	variant     variant.Variant
	atls.Validator
}

// NewValidator initializes a new updatable validator.
func NewValidator(log *logger.Logger, variant variant.Variant, fileHandler file.Handler) (*Updatable, error) {
	u := &Updatable{
		log:         log,
		fileHandler: fileHandler,
		variant:     variant,
	}

	if err := u.Update(); err != nil {
		return nil, err
	}
	return u, nil
}

// Validate calls the validators Validate method, and prevents any updates during the call.
func (u *Updatable) Validate(ctx context.Context, attDoc []byte, nonce []byte) ([]byte, error) {
	u.mux.Lock()
	defer u.mux.Unlock()
	return u.Validator.Validate(ctx, attDoc, nonce)
}

// OID returns the validators Object Identifier.
func (u *Updatable) OID() asn1.ObjectIdentifier {
	u.mux.Lock()
	defer u.mux.Unlock()
	return u.Validator.OID()
}

// Update switches out the underlying validator.
func (u *Updatable) Update() error {
	u.mux.Lock()
	defer u.mux.Unlock()

	u.log.Infof("Updating expected measurements")

	var measurements measurements.M
	if err := u.fileHandler.ReadJSON(filepath.Join(constants.ServiceBasePath, constants.MeasurementsFilename), &measurements); err != nil {
		return err
	}
	u.log.Debugf("New measurements: %+v", measurements)

	// Read ID Key config
	var idKeyCfg config.SNPFirmwareSignerConfig
	if u.variant.Equal(variant.AzureSEVSNP{}) {
		u.log.Infof("Updating SEV-SNP ID Key config")

		err := u.fileHandler.ReadJSON(filepath.Join(constants.ServiceBasePath, constants.IDKeyConfigFilename), &idKeyCfg)
		if err != nil {
			if !errors.Is(err, afero.ErrFileNotFound) {
				return fmt.Errorf("reading ID Key config: %w", err)
			}

			u.log.Warnf("ID Key config file not found, falling back to old format (v2.6 or earlier)")

			// v2.6 fallback
			// TODO: Remove after v2.7 release
			var digest idkeydigest.List
			var enforceIDKeyDigest idkeydigest.Enforcement
			enforceRaw, err := u.fileHandler.Read(filepath.Join(constants.ServiceBasePath, constants.EnforceIDKeyDigestFilename))
			if err != nil {
				return err
			}
			enforceIDKeyDigest = idkeydigest.EnforcePolicyFromString(string(enforceRaw))
			if err != nil {
				return fmt.Errorf("parsing content of EnforceIdKeyDigestFilename: %s: %w", enforceRaw, err)
			}

			idkeydigestRaw, err := u.fileHandler.Read(filepath.Join(constants.ServiceBasePath, constants.IDKeyDigestFilename))
			if err != nil {
				return err
			}
			if err = json.Unmarshal(idkeydigestRaw, &digest); err != nil {
				return fmt.Errorf("unmarshaling content of IDKeyDigestFilename: %s: %w", idkeydigestRaw, err)
			}

			idKeyCfg.AcceptedKeyDigests = digest
			idKeyCfg.EnforcementPolicy = enforceIDKeyDigest
		}

		u.log.Debugf("New ID Key config: %+v", idKeyCfg)
	}

	validator, err := choose.Validator(u.variant, measurements, idKeyCfg, u.log)
	if err != nil {
		return fmt.Errorf("updating validator: %w", err)
	}
	u.Validator = validator

	return nil
}
