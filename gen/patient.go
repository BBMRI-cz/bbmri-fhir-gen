// Copyright © 2019 The Samply Development Community
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gen

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Store seen identifiers to allow duplication
var seenIdentifiers []string

func Patient(r *rand.Rand, idx int) Object {
	patient := make(Object)
	patient["resourceType"] = "Patient"
	patient["id"] = fmt.Sprintf("bbmri-%d", idx)
	patient["meta"] = meta("https://fhir.bbmri.de/StructureDefinition/Patient")

	// Set gender only if not null
	if gender := randGender(r); gender != "" {
		patient["gender"] = gender
	}

	birthDate := randDate(r, 1950, 2030)
	patient["birthDate"] = birthDate.Format("2006-01-02")

	deceasedDate := birthDate.Add(randAge(r))
	if deceasedDate.Before(time.Now()) {
		patient["deceasedDateTime"] = deceasedDate.Format("2006-01-02")
	}

	// Add an identifier with a small chance of duplication
	var identifierValue string
	if len(seenIdentifiers) > 0 && r.Float64() < 0.05 {
		// 5% chance: reuse a previously seen identifier
		identifierValue = seenIdentifiers[r.Intn(len(seenIdentifiers))]
	} else {
		// Generate a new identifier
		identifierValue = fmt.Sprintf("id-%06d", r.Intn(1_000_000))
		seenIdentifiers = append(seenIdentifiers, identifierValue)
	}

	patient["identifier"] = []interface{}{
		map[string]interface{}{
			"system": "https://fhir.bbmri.de/id/patient",
			"value":  identifierValue,
		},
	}

	return patient
}

var genders = []string{"male", "female"}

func randGender(r *rand.Rand) string {
	gender := genders[r.Intn(len(genders))]

	// 10% chance to introduce an error
	if r.Float64() < 0.1 {
		switch r.Intn(3) {
		case 0:
			// Return null-equivalent (empty string)
			return ""
		case 1:
			// Add a random character
			pos := r.Intn(len(gender) + 1)
			randomChar := string('a' + rune(r.Intn(26)))
			gender = gender[:pos] + randomChar + gender[pos:]
		case 2:
			// Replace a character
			if len(gender) > 0 {
				pos := r.Intn(len(gender))
				randomChar := string('a' + rune(r.Intn(26)))
				var sb strings.Builder
				sb.WriteString(gender[:pos])
				sb.WriteString(randomChar)
				if pos+1 < len(gender) {
					sb.WriteString(gender[pos+1:])
				}
				gender = sb.String()
			}
		}
	}

	return gender
}

func randDate(r *rand.Rand, startYear int, endYear int) time.Time {
	start := time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(endYear, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(start+r.Int63n(end-start), 0)
}

func randAge(r *rand.Rand) time.Duration {
	return time.Duration(r.Intn(80*365)+10*365) * 24 * time.Hour
}
