// Copyright 2018 The agentx authors
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package value

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/errgo.v1"
)

// OID defines an OID.
type OID []uint32

// ParseOID parses the provided string and returns a valid oid. If one of the
// subidentifers canot be parsed to an uint32, the function will panic.
func ParseOID(text string) (OID, error) {
	var result OID

	parts := strings.Split(text, ".")
	for _, part := range parts {
		subidentifier, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return nil, errgo.Notef(err, "could not subidentifier [%s] to uint32", part)
		}
		result = append(result, uint32(subidentifier))
	}

	return result, nil
}

// MustParseOID works like ParseOID expect it panics on a parsing error.
func MustParseOID(text string) OID {
	result, err := ParseOID(text)
	if err != nil {
		panic(err)
	}
	return result
}

// First returns the first n subidentifiers as a new oid.
func (o OID) First(count int) OID {
	return o[:count]
}

// CommonPrefix compares the oid with the provided one and
// returns a new oid containing all matching prefix subidentifiers.
func (o OID) CommonPrefix(other OID) OID {
	matchCount := 0

	for index, subidentifier := range o {
		if index >= len(other) || subidentifier != other[index] {
			break
		}
		matchCount++
	}

	return o[:matchCount]
}

func (o OID) String() string {
	var parts []string

	for _, subidentifier := range o {
		parts = append(parts, fmt.Sprintf("%d", subidentifier))
	}

	return strings.Join(parts, ".")
}

// This function performs a lexical compare of two OIDs a and b. It will return
// 0 if A = B
// 1 if A is DEEPER in the tree then B
// -1 of A is HIGHER in the tree then B
func (o OID) Compare(b OID) int {
	a := o
	// Check if a is longer then b
	if len(a) > len(b) {
		// Find the first value in b which does not match with a.
		// If the value/digit is greater in a then it is in b
		// a must bemust be deeper in the tree
		for i := 0; i < len(b); i++ {
			if a[i] != b[i] {
				if a[i] > b[i] {
					return 1
				} else {
					return -1
				}
			}
		}
		// If a is longer then b, and has no value different
		// from b up to the length of b, then a must be deeper
		return 1
		// Check if a is shorter then b
	} else if len(a) < len(b) {
		// Find the first value in b which does not match with a.
		// If the value/digit is less in a then it is in b
		// a must be higher in the tree
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				if a[i] < b[i] {
					return -1
				} else {
					return 1
				}
			}
		}
		// If a is shorter then b and and has no value different
		// from b up to the length of a, then a must be higher
		return -1
		// a and b are equal length.
	} else {
		for i := 0; i < len(a); i++ {
			// if a has any digit > b then it must be deeper in the tree
			if a[i] > b[i] {
				return 1
			}
			// if b has any digit > a then it must be deeper in the tree
			if b[i] > a[i] {
				return -1
			}
		}
		// If we get here the two must be equal
		return 0
	}
}

// Returns true if oid is between from and to in the tree. If includeFrom is
// true, the function will also return true if oid matches from
func (o OID) Within(from OID, includeFrom bool, to OID) bool {
	// We have to determine if the "oid" comes after "from"
	// eg. Is further down then from in the tree AND if it comes before "to"
	// eg. Is above to in the tree

	// If includeFrom is true and compare shows that from
	// and oid match then the result must be true
	fromCompare := o.Compare(from)
	toCompare := o.Compare(to)

	return (fromCompare == 0 && includeFrom == true) || (fromCompare == 1 && toCompare == -1)
}
