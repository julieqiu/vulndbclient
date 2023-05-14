// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	tcs := []struct {
		endpoint string
	}{
		{
			endpoint: "index/db",
		},
		{
			endpoint: "index/modules",
		},
		{
			endpoint: "ID/GO-2021-0068",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.endpoint, func(t *testing.T) {
			test := func(t *testing.T, s source) {
				got, err := s.get(context.Background(), tc.endpoint)
				if err != nil {
					t.Fatal(err)
				}
				want, err := os.ReadFile(testVulndb + "/" + tc.endpoint + ".json")
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("mismatch(-want, +got): %s", diff)
				}
			}
			testAllSourceTypes(t, test)
		})
	}
}

// testAllSourceTypes runs a given test for all source types.
func testAllSourceTypes(t *testing.T, test func(t *testing.T, s source)) {
	t.Run("http", func(t *testing.T) {
		srv := newTestServer(testVulndb)
		hs := newHTTPSource(srv.URL, &Options{HTTPClient: srv.Client()})
		test(t, hs)
	})

	t.Run("local", func(t *testing.T) {
		uri, err := url.Parse(testVulndbFileURL)
		if err != nil {
			t.Fatal(err)
		}

		fs, err := newLocalSource(uri)
		if err != nil {
			t.Fatal(err)
		}

		test(t, fs)
	})

	t.Run("in-memory", func(t *testing.T) {
		testEntries, err := entries(testIDs)
		if err != nil {
			t.Fatal(err)
		}

		ms, err := newInMemorySource(testEntries)
		if err != nil {
			t.Fatal(err)
		}

		test(t, ms)
	})
}
