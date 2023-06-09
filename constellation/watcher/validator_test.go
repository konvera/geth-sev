/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package watcher

import (
	"context"
	"encoding/asn1"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/konvera/geth-sev/constellation/atls"
	"github.com/konvera/geth-sev/constellation/attestation/idkeydigest"
	"github.com/konvera/geth-sev/constellation/attestation/measurements"
	"github.com/konvera/geth-sev/constellation/config"
	"github.com/konvera/geth-sev/constellation/constants"
	"github.com/konvera/geth-sev/constellation/file"
	"github.com/konvera/geth-sev/constellation/logger"
	"github.com/konvera/geth-sev/constellation/variant"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		// https://github.com/census-instrumentation/opencensus-go/issues/1262
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
	)
}

func TestNewUpdateableValidator(t *testing.T) {
	testCases := map[string]struct {
		variant   variant.Variant
		writeFile bool
		wantErr   bool
	}{
		"azure": {
			variant:   variant.AzureSEVSNP{},
			writeFile: true,
		},
		"gcp": {
			variant:   variant.GCPSEVES{},
			writeFile: true,
		},
		"qemu": {
			variant:   variant.QEMUVTPM{},
			writeFile: true,
		},
		"no file": {
			variant:   variant.AzureSEVSNP{},
			writeFile: false,
			wantErr:   true,
		},
		"invalid provider": {
			variant:   fakeOID{1, 3, 9900, 9999, 9999},
			writeFile: true,
			wantErr:   true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			handler := file.NewHandler(afero.NewMemMapFs())
			if tc.writeFile {
				require.NoError(handler.WriteJSON(
					filepath.Join(constants.ServiceBasePath, constants.MeasurementsFilename),
					measurements.M{11: measurements.WithAllBytes(0x00, measurements.Enforce)},
				))

				require.NoError(handler.WriteJSON(
					filepath.Join(constants.ServiceBasePath, constants.IDKeyConfigFilename),
					config.SNPFirmwareSignerConfig{
						AcceptedKeyDigests: idkeydigest.DefaultList(),
						EnforcementPolicy:  idkeydigest.WarnOnly,
					},
				))
			}

			_, err := NewValidator(
				logger.NewTest(t),
				tc.variant,
				handler,
			)
			if tc.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := file.NewHandler(afero.NewMemMapFs())

	// create server
	validator := &Updatable{
		log:         logger.NewTest(t),
		variant:     variant.Dummy{},
		fileHandler: handler,
	}

	// Update should fail if the file does not exist
	assert.Error(validator.Update())

	// write measurement config
	require.NoError(handler.WriteJSON(
		filepath.Join(constants.ServiceBasePath, constants.MeasurementsFilename),
		measurements.M{11: measurements.WithAllBytes(0x00, measurements.Enforce)},
	))
	require.NoError(handler.WriteJSON(
		filepath.Join(constants.ServiceBasePath, constants.IDKeyConfigFilename),
		config.SNPFirmwareSignerConfig{
			AcceptedKeyDigests: idkeydigest.List{[]byte{0x00}},
			EnforcementPolicy:  idkeydigest.WarnOnly,
		},
	))

	// call update once to initialize the server's validator
	require.NoError(validator.Update())

	// create tls config and start the server
	serverConfig, err := atls.CreateAttestationServerTLSConfig(nil, []atls.Validator{validator})
	require.NoError(err)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello")
	}))
	server.TLS = serverConfig
	server.StartTLS()
	defer server.Close()

	// test connection to server
	clientOID := variant.Dummy{}
	resp, err := testConnection(require, server.URL, clientOID)
	require.NoError(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(err)
	assert.EqualValues("hello", body)

	// update the server's validator
	validator.variant = variant.QEMUVTPM{}
	require.NoError(validator.Update())

	// client connection should fail now, since the server's validator expects a different OID from the client
	resp, err = testConnection(require, server.URL, clientOID)
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Error(err)

	// test old ID Key Digest format
	require.NoError(handler.Write(
		filepath.Join(constants.ServiceBasePath, constants.IDKeyDigestFilename),
		[]byte{},
	))
	require.NoError(handler.Write(
		filepath.Join(constants.ServiceBasePath, constants.EnforceIDKeyDigestFilename),
		[]byte("false"),
	))

	assert.NoError(validator.Update())
}

func TestOIDConcurrency(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := file.NewHandler(afero.NewMemMapFs())
	require.NoError(handler.WriteJSON(
		filepath.Join(constants.ServiceBasePath, constants.MeasurementsFilename),
		measurements.M{11: measurements.WithAllBytes(0x00, measurements.Enforce)},
	))
	require.NoError(handler.WriteJSON(
		filepath.Join(constants.ServiceBasePath, constants.IDKeyConfigFilename),
		config.SNPFirmwareSignerConfig{
			AcceptedKeyDigests: idkeydigest.List{[]byte{0x00}},
			EnforcementPolicy:  idkeydigest.WarnOnly,
		},
	))

	// create server
	validator := &Updatable{
		log:         logger.NewTest(t),
		variant:     variant.Dummy{},
		fileHandler: handler,
	}

	// call update once to initialize the server's validator
	require.NoError(validator.Update())

	var wg sync.WaitGroup
	wg.Add(2 * 20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			assert.NoError(validator.Update())
		}()
		go func() {
			defer wg.Done()
			validator.OID()
		}()
	}
	wg.Wait()
}

func TestUpdateConcurrency(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := file.NewHandler(afero.NewMemMapFs())
	validator := &Updatable{
		log:         logger.NewTest(t),
		fileHandler: handler,
		variant:     variant.Dummy{},
	}
	require.NoError(handler.WriteJSON(
		filepath.Join(constants.ServiceBasePath, constants.MeasurementsFilename),
		measurements.M{11: measurements.WithAllBytes(0x00, measurements.Enforce)},
		file.OptNone,
	))
	require.NoError(handler.WriteJSON(
		filepath.Join(constants.ServiceBasePath, constants.IDKeyConfigFilename),
		config.SNPFirmwareSignerConfig{
			AcceptedKeyDigests: idkeydigest.List{[]byte{0x00}},
			EnforcementPolicy:  idkeydigest.WarnOnly,
		},
	))

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(validator.Update())
		}()
	}

	wg.Wait()
}

func testConnection(require *require.Assertions, url string, oid variant.Getter) (*http.Response, error) {
	clientConfig, err := atls.CreateAttestationClientTLSConfig(fakeIssuer{oid}, nil)
	require.NoError(err)
	client := http.Client{Transport: &http.Transport{TLSClientConfig: clientConfig}}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	require.NoError(err)
	return client.Do(req)
}

type fakeIssuer struct {
	variant.Getter
}

func (fakeIssuer) Issue(_ context.Context, userData []byte, nonce []byte) ([]byte, error) {
	return json.Marshal(fakeDoc{UserData: userData, Nonce: nonce})
}

type fakeOID asn1.ObjectIdentifier

func (o fakeOID) OID() asn1.ObjectIdentifier {
	return asn1.ObjectIdentifier(o)
}

func (o fakeOID) String() string {
	return o.OID().String()
}

func (o fakeOID) Equal(other variant.Getter) bool {
	return o.OID().Equal(other.OID())
}

type fakeDoc struct {
	UserData []byte
	Nonce    []byte
}
