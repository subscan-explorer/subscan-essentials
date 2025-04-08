package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnumPickOne(t *testing.T) {
	m := map[string]string{
		"key1": "value1",
		"key2": "",
	}
	expected := "value1"
	result := EnumPickOne(m)
	assert.Equal(t, expected, result)
}

func TestEnumPickOneInt(t *testing.T) {
	m := map[string]int{
		"key1": 0,
	}
	expected := 0
	result := EnumPickOneInt(m)
	assert.Equal(t, expected, result)
}

func TestEnumStringKey(t *testing.T) {
	m := map[string]string{
		"key1": "value1",
	}
	expected := "key1"
	result := EnumStringKey(m)
	assert.Equal(t, expected, result)
}

func TestEnumKey(t *testing.T) {
	m := map[string]interface{}{
		"key1": "value1",
	}
	expected := "key1"
	result := EnumKey(m)
	assert.Equal(t, expected, result)
}
