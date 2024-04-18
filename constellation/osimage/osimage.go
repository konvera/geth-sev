/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

// package osimage is used to handle osimages in the CI (uploading and maintenance).
package osimage

import (
	"io"
	"time"

	"github.com/konvera/geth-sev/constellation/api/versionsapi"
	"github.com/konvera/geth-sev/constellation/cloud/cloudprovider"
)

// UploadRequest is a request to upload an os image.
type UploadRequest struct {
	Provider           cloudprovider.Provider
	Version            versionsapi.Version
	AttestationVariant string
	Timestamp          time.Time
	ImageReader        func() (io.ReadSeekCloser, error)
	ImagePath          string
}
