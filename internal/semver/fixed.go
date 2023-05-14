// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import "github.com/google/osv-scanner/pkg/models"

func LatestFixedVersion(ranges []models.Range) string {
	var latestFixed string
	for _, r := range ranges {
		if r.Type == "SEMVER" {
			for _, e := range r.Events {
				fixed := e.Fixed
				if fixed != "" && Less(latestFixed, fixed) {
					latestFixed = fixed
				}
			}
			// If the vulnerability was re-introduced after the latest fix
			// we found, there is no latest fix for this range.
			for _, e := range r.Events {
				introduced := e.Introduced
				if introduced != "" && introduced != "0" && Less(latestFixed, introduced) {
					latestFixed = ""
					break
				}
			}
		}
	}
	return latestFixed
}
