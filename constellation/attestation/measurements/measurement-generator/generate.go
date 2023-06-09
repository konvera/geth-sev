/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/konvera/geth-sev/constellation/attestation/measurements"
	"github.com/konvera/geth-sev/constellation/cloud/cloudprovider"
	"github.com/konvera/geth-sev/constellation/config"
	"github.com/konvera/geth-sev/constellation/constants"
	"github.com/konvera/geth-sev/constellation/sigstore"
	"github.com/konvera/geth-sev/constellation/versionsapi"
	"golang.org/x/tools/go/ast/astutil"
)

// this tool is used to generate hardcoded measurements for the enterprise build.
// Measurements are embedded in the constellation cli.

func main() {
	defaultConf := config.Default()
	log.Printf("Generating measurements for %s\n", defaultConf.Image)

	const filePath = "./measurements_enterprise.go"

	ctx := context.Background()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	rekor, err := sigstore.NewRekor()
	if err != nil {
		log.Fatal(err)
	}

	var returnStmtCtr int

	newFile := astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		n := cursor.Node()

		// find all switch cases for the CSPs of the form:
		// switch provider {
		// case cloudprovider.XYZ:
		// 	return M{...}

		if clause, ok := n.(*ast.CaseClause); ok && len(clause.List) > 0 && len(clause.Body) > 0 {
			sel, ok := clause.List[0].(*ast.SelectorExpr)
			if !ok {
				return true
			}
			returnStmt, ok := clause.Body[0].(*ast.ReturnStmt)
			if !ok || len(returnStmt.Results) == 0 {
				return true
			}

			provider := cloudprovider.FromString(sel.Sel.Name)
			if provider == cloudprovider.Unknown {
				log.Fatalf("unknown provider %s", sel.Sel.Name)
			}
			log.Println("Found", provider)
			returnStmtCtr++
			// retrieve and validate measurements for the given CSP and image
			measuremnts := mustGetMeasurements(ctx, rekor, []byte(constants.CosignPublicKey), http.DefaultClient, provider, defaultConf.Image)
			// replace the return statement with a composite literal containing the validated measurements
			returnStmt.Results[0] = measurementsCompositeLiteral(measuremnts)
		}
		return true
	}, nil,
	)

	if returnStmtCtr == 0 {
		log.Fatal("no measurements updated")
	}

	var buf bytes.Buffer
	printConfig := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

	if err = printConfig.Fprint(&buf, fset, newFile); err != nil {
		log.Fatalf("error formatting file %s: %s", filePath, err)
	}
	if err := os.WriteFile(filePath, buf.Bytes(), 0o644); err != nil {
		log.Fatalf("error writing file %s: %s", filePath, err)
	}
	log.Println("Successfully generated hashes.")
}

// mustGetMeasurements fetches the measurements for the given image and CSP and verifies them.
func mustGetMeasurements(ctx context.Context, verifier rekorVerifier, cosignPublicKey []byte, client *http.Client, provider cloudprovider.Provider, image string) measurements.M {
	measurementsURL, err := measurementURL(provider, image, "measurements.json")
	if err != nil {
		panic(err)
	}
	signatureURL, err := measurementURL(provider, image, "measurements.json.sig")
	if err != nil {
		panic(err)
	}

	log.Println("Fetching measurements from", measurementsURL, "and signature from", signatureURL)
	var fetchedMeasurements measurements.M
	hash, err := fetchedMeasurements.FetchAndVerify(
		ctx, client,
		measurementsURL,
		signatureURL,
		cosignPublicKey,
		measurements.WithMetadata{
			CSP:   provider,
			Image: image,
		},
	)
	if err != nil {
		panic(err)
	}
	if err := verifyWithRekor(ctx, verifier, hash); err != nil {
		panic(err)
	}
	return fetchedMeasurements
}

// measurementURL returns the URL for the measurements file for the given image and CSP.
func measurementURL(provider cloudprovider.Provider, image, file string) (*url.URL, error) {
	version, err := versionsapi.NewVersionFromShortPath(image, versionsapi.VersionKindImage)
	if err != nil {
		return nil, fmt.Errorf("parsing image name: %w", err)
	}

	return url.Parse(
		version.ArtifactsURL() + path.Join("/image", "csp", strings.ToLower(provider.String()), file),
	)
}

// verifyWithRekor verifies that the given hash is present in rekor and is valid.
func verifyWithRekor(ctx context.Context, verifier rekorVerifier, hash string) error {
	uuids, err := verifier.SearchByHash(ctx, hash)
	if err != nil {
		return fmt.Errorf("searching Rekor for hash: %w", err)
	}

	if len(uuids) == 0 {
		return fmt.Errorf("no matching entries in Rekor")
	}

	// We expect the first entry in Rekor to be our original entry.
	// SHA256 should ensure there is no entry with the same hash.
	// Any subsequent hashes are treated as potential attacks and are ignored.
	// Attacks on Rekor will be monitored from other backend services.
	artifactUUID := uuids[0]

	return verifier.VerifyEntry(
		ctx, artifactUUID,
		base64.StdEncoding.EncodeToString([]byte(constants.CosignPublicKey)),
	)
}

// byteArrayCompositeLit returns a *ast.CompositeLit representing a byte array literal.
// The returned literal is of the form:
// [32]byte{ 0x01, 0x02, 0x03, ... }.
func byteArrayCompositeLit(hex [32]byte) *ast.CompositeLit {
	var elts []ast.Expr
	// create list of byte literals
	for _, b := range hex {
		elts = append(elts, &ast.BasicLit{
			Kind:  token.INT,
			Value: fmt.Sprintf("0x%02x", b),
		})
	}
	// create the composite literal
	// containing the byte literals as part of an [32]byte array
	return &ast.CompositeLit{
		Type: &ast.ArrayType{
			Len: ast.NewIdent("32"),
			Elt: ast.NewIdent("byte"),
		},
		Elts: elts,
	}
}

// measurementsEntryKeyValueExpr returns a *ast.KeyValueExpr representing a measurements.Measurement entry.
// The returned expression is of the form:
//
//	0: {
//		  Expected: [32]byte{ 0x01, 0x02, 0x03, ... },
//		  WarnOnly: false,
//	},
func measurementsEntryKeyValueExpr(pcr uint32, measuremnt measurements.Measurement) *ast.KeyValueExpr {
	key := fmt.Sprintf("%d", pcr)

	var validationOptString string
	if measuremnt.ValidationOpt {
		validationOptString = "WarnOnly"
	} else {
		validationOptString = "Enforce"
	}

	return &ast.KeyValueExpr{
		Key: &ast.BasicLit{
			Kind:  token.INT,
			Value: key,
		},
		Value: &ast.CompositeLit{
			Elts: []ast.Expr{
				&ast.KeyValueExpr{
					Key:   &ast.Ident{Name: "Expected"},
					Value: byteArrayCompositeLit(measuremnt.Expected),
				},
				&ast.KeyValueExpr{
					Key:   &ast.Ident{Name: "ValidationOpt"},
					Value: &ast.Ident{Name: validationOptString},
				},
			},
		},
	}
}

// measurementsCompositeLiteral returns a *ast.CompositeLit representing a measurements.M literal.
// The returned literal is of the form:
//
//	M{
//		0: {
//			Expected: [32]byte{ 0x01, 0x02, 0x03, ... },
//			WarnOnly: false,
//		},
//		1: {
//			Expected: [32]byte{ 0x01, 0x02, 0x03, ... },
//			WarnOnly: false,
//		},
//		...
//	}
func measurementsCompositeLiteral(measurements measurements.M) *ast.CompositeLit {
	var elts []ast.Expr
	pcrs := make([]uint32, 0, len(measurements))
	for pcr := range measurements {
		pcrs = append(pcrs, pcr)
	}
	sort.Slice(pcrs, func(i, j int) bool { return pcrs[i] < pcrs[j] })
	for _, pcr := range pcrs {
		kvExpr := measurementsEntryKeyValueExpr(pcr, measurements[pcr])
		elts = append(elts, kvExpr)
	}
	return &ast.CompositeLit{
		Type: &ast.Ident{
			Name: "M",
		},
		Elts: elts,
	}
}

type rekorVerifier interface {
	SearchByHash(context.Context, string) ([]string, error)
	VerifyEntry(context.Context, string, string) error
}
