package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSliceSetElem(t *testing.T) {
	sliceInt := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	sliceIntNew := make([]int, 0)
	err := SliceSetElem(reflect.ValueOf(&sliceIntNew).Elem(), len(sliceInt), func(i int, elem reflect.Value) (bool, error) {
		if i >= len(sliceInt) {
			return false, nil
		}
		elem.SetInt(int64(sliceInt[i]))
		return true, nil
	})
	if assert.NoError(t, err) {
		assert.Equal(t, sliceInt, sliceIntNew)
	}

	arrayInt := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	arrayIntNew := [9]int{}
	err = SliceSetElem(reflect.ValueOf(&arrayIntNew).Elem(), len(arrayInt), func(i int, elem reflect.Value) (bool, error) {
		if i >= len(arrayInt) {
			return false, nil
		}
		elem.SetInt(int64(sliceInt[i]))
		return true, nil
	})
	if assert.NoError(t, err) {
		assert.Equal(t, arrayInt, arrayIntNew)
	}

	arrayIntNew = [9]int{}
	err = SliceSetElem(reflect.ValueOf(&arrayIntNew).Elem(), 50, func(i int, elem reflect.Value) (bool, error) {
		if i >= 50 {
			return false, nil
		}
		elem.SetInt(int64(i + 1))
		return true, nil
	})
	if assert.NoError(t, err) {
		assert.Equal(t, arrayInt, arrayIntNew)
	}

	a, b, c := 1, 2, 3
	slicePtr := []*int{&a, &b, &c}
	slicePtrNew := make([]*int, 0, 1)
	err = SliceSetElem(reflect.ValueOf(&slicePtrNew).Elem(), len(slicePtr), func(i int, elem reflect.Value) (bool, error) {
		if i >= len(slicePtr) {
			return false, nil
		}
		elem.SetInt(int64(*slicePtr[i]))
		return true, nil
	})
	if assert.NoError(t, err) {
		assert.Equal(t, slicePtr, slicePtrNew)
	}
}

func TestPtrValue(t *testing.T) {
	var a *int
	aValue := reflect.ValueOf(&a)
	aValue = PtrValue(aValue)
	aValue.SetInt(1)
	assert.Equal(t, *a, 1)

	var b int
	bValue := reflect.ValueOf(&b)
	bValue = PtrValue(bValue)
	bValue.SetInt(1)
	assert.Equal(t, b, 1)

	var c ***string
	cValue := reflect.ValueOf(&c)
	cValue = PtrValue(cValue)
	cValue.SetString("hello")
	assert.Equal(t, ***c, "hello")
}
