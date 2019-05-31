// Copyright 2018 The agentx authors
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package value_test

import (
	"testing"

	. "github.com/eriksejr/go-agentx/test"
	. "github.com/eriksejr/go-agentx/value"
)

func TestCommonPrefix(t *testing.T) {
	oid := MustParseOID("1.3.6.1.2")
	result := oid.CommonPrefix(MustParseOID("1.3.6.1.4"))
	AssertEquals(t, MustParseOID("1.3.6.1"), result)
}
