/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package ancestry

import (
	"io/fs"
	"strings"
	"unicode/utf8"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/rand"
)

const (
	nameGeneratorTypeKey         = "type"
	nameGeneratorTrainingDataKey = "training_data"
)

type charThreshold struct {
	ch        rune
	threshold int
}

// NameGenerator holds the data for name generation.
type NameGenerator struct {
	generationType NameGenerationType
	trainingData   []string
	min            int
	max            int
	entries        map[string][]charThreshold
}

// NewNameGenerator creates a new NameGenerator.
func NewNameGenerator(generationType NameGenerationType, trainingData []string) *NameGenerator {
	g := &NameGenerator{
		generationType: generationType,
		trainingData:   make([]string, 0, len(trainingData)),
	}
	for _, one := range trainingData {
		one = strings.ToLower(strings.TrimSpace(one))
		if utf8.RuneCountInString(one) >= 2 {
			g.trainingData = append(g.trainingData, one)
		}
	}
	if generationType == MarkovChain {
		g.min = 20
		g.max = 2
		builders := make(map[string]map[rune]int)
		for _, one := range g.trainingData {
			runes := []rune(one)
			length := len(runes)
			if g.min > length {
				g.min = length
			}
			if g.max < length {
				g.max = length
			}
			for i := 2; i < length; i++ {
				charGroup := string(runes[i-2 : i])
				occurrences, exists := builders[charGroup]
				if !exists {
					occurrences = make(map[rune]int)
					builders[charGroup] = occurrences
				}
				occurrences[runes[i]]++
			}
		}
		g.entries = make(map[string][]charThreshold)
		for k, v := range builders {
			g.entries[k] = makeCharThresholdEntry(v)
		}
	}
	return g
}

// NewNameGeneratorFromFS creates a new NameGenerator from a file.
func NewNameGeneratorFromFS(fileSystem fs.FS, filePath string) (*NameGenerator, error) {
	data, err := encoding.LoadJSONFromFS(fileSystem, filePath)
	if err != nil {
		return nil, err
	}
	return NewNameGeneratorFromJSON(encoding.Object(data)), nil
}

// NewNameGeneratorFromJSON creates a new NameGenerator from a JSON object.
func NewNameGeneratorFromJSON(data map[string]interface{}) *NameGenerator {
	array := encoding.Array(data[nameGeneratorTrainingDataKey])
	trainingData := make([]string, len(array))
	for i, one := range array {
		trainingData[i] = encoding.String(one)
	}
	return NewNameGenerator(NameGenerationTypeFromKey(encoding.String(data[nameGeneratorTypeKey])), trainingData)
}

// ToJSON emits this object as JSON.
func (g *NameGenerator) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(nameGeneratorTypeKey, g.generationType.Key(), false, false)
	encoder.Key(nameGeneratorTrainingDataKey)
	encoder.StartArray()
	for _, one := range g.trainingData {
		encoder.String(one)
	}
	encoder.EndArray()
	encoder.EndObject()
}

// Generate a name.
func (g *NameGenerator) Generate() string {
	rnd := rand.NewCryptoRand()
	switch g.generationType {
	case Simple:
		if len(g.trainingData) == 0 {
			return ""
		}
		return txt.FirstToUpper(g.trainingData[rnd.Intn(len(g.trainingData))])
	case MarkovChain:
		var buffer strings.Builder
		var sub []rune
		for k := range g.entries {
			buffer.WriteString(txt.FirstToUpper(k))
			sub = []rune(k)
			break // Only want one, which is random
		}
		targetSize := g.min + rnd.Intn(g.max+1-g.min)
		for i := 2; i < targetSize; i++ {
			entry, exists := g.entries[string(sub)]
			if !exists {
				break
			}
			next := chooseCharacter(entry)
			if next == 0 {
				break
			}
			buffer.WriteRune(next)
			sub[0] = sub[1]
			sub[1] = next
		}
		return buffer.String()
	default:
		return ""
	}
}

func makeCharThresholdEntry(occurrences map[rune]int) []charThreshold {
	ct := make([]charThreshold, len(occurrences))
	i := 0
	for k, v := range occurrences {
		ct[i].ch = k
		if i > 0 {
			v += ct[i-1].threshold
		}
		ct[i].threshold = v
		i++
	}
	return ct
}

func chooseCharacter(ct []charThreshold) rune {
	threshold := rand.NewCryptoRand().Intn(ct[len(ct)-1].threshold + 1)
	for i := range ct {
		if ct[i].threshold >= threshold {
			return ct[i].ch
		}
	}
	return 0
}
