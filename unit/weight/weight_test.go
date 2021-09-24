package weight_test

import (
	"encoding/json"
	"testing"

	"github.com/richardwilkes/gcs/unit/weight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type embeddedWeight struct {
	Field weight.Weight
}

func TestWeightConversion(t *testing.T) {
	assert.Equal(t, "1 lb", weight.FromInt64(1, weight.Pound).Format(weight.Pound))
	assert.Equal(t, "15 lb", weight.FromInt64(15, weight.Pound).Format(weight.Pound))
	assert.Equal(t, "0.5 kg", weight.FromInt64(1, weight.Pound).Format(weight.Kilogram))
	assert.Equal(t, "7.5 kg", weight.FromInt64(15, weight.Pound).Format(weight.Kilogram))

	w, err := weight.FromString("1", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "1 lb", w.String())
	w, err = weight.FromString("1", weight.Kilogram)
	assert.NoError(t, err)
	assert.Equal(t, "2 lb", w.String())
	w, err = weight.FromString("22.34 lb", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "22.34 lb", w.String())
	w, err = weight.FromString(" +22.34   lb  ", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "22.34 lb", w.String())
	w, err = weight.FromString("0.5kg", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "0.5 kg", w.Format(weight.Kilogram))
	w, err = weight.FromString("15.25kg", weight.Pound)
	assert.NoError(t, err)
	assert.Equal(t, "15.25 kg", w.Format(weight.Kilogram))
}

func TestWeightJSON(t *testing.T) {
	inc := weight.FromFloat64(1.0/3.0, weight.Pound)
	max := weight.FromInt64(5, weight.Pound)
	for i := weight.Weight(0); i <= max; i += inc {
		e1 := embeddedWeight{Field: i}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embeddedWeight
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}
