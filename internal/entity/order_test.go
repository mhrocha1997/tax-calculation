package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_If_I_Get_An_Error_If_Id_Is_Blank(t *testing.T) {
	order := Order{}
	assert.Error(t, order.Validate(), "invalid id")
}
func Test_If_Gets_An_Error_If_Price_Is_Blank(t *testing.T) {
	order := Order{ID: "123"}
	assert.Error(t, order.Validate(), "invalid price")
}

func Test_If_Gets_An_Error_If_Tax_Is_Blank(t *testing.T) {
	order := Order{ID: "123", Price: 10}
	assert.Error(t, order.Validate(), "invalid price")
}

func Test_AllParams(t *testing.T) {
	order := Order{ID: "123", Price: 10, Tax: 10}
	assert.NoError(t, order.Validate())
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 10.0, order.Tax)
	assert.NoError(t, order.CalculateFinalPrice())
	assert.Equal(t, 20.0, order.FinalPrice)

}
