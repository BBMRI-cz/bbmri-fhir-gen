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
	"time"
)

func Condition(r *rand.Rand, patientIdx int, conditionIdx int, date time.Time) Object {
	condition := Object{
		"resourceType":  "Condition",
		"id":            fmt.Sprintf("bbmri-%d-condition-%d", patientIdx, conditionIdx),
		"meta":          meta("https://fhir.bbmri.de/StructureDefinition/Condition"),
		"subject":       patientReference(patientIdx),
		"onsetDateTime": date.Format("2006-01-02"),
	}

	code := randIcd10Code(r)
	if code != "empty" {
		// only generate coding + codeableConcept if code is non-empty
		coding := codingWithVersion("http://hl7.org/fhir/sid/icd-10", "2016", code)
		if coding != nil {
			cc := codeableConcept(coding)
			if cc != nil {
				condition["code"] = cc
			}
		}
	}

	return condition
}



func randIcd10Code(r *rand.Rand) string {
	// 5% chance to return a null-equivalent
	if r.Float64() < 0.05 {
		return "empty"
	}
	return fmt.Sprintf("%s%02d.%d", string(rune(65+r.Intn(26))), r.Intn(100), r.Intn(10))
}
