// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type GenericClass struct {
	Module string
	Name   string
}

var _ PyNewable = &GenericClass{}
var _ DictSetter = &GenericObject{}
var _ PyDictSettable = &GenericObject{}

type GenericObject struct {
	Class           *GenericClass
	ConstructorArgs []interface{}
}

func NewGenericClass(module, name string) *GenericClass {
	return &GenericClass{Module: module, Name: name}
}

func (goj *GenericObject) Set(key, value interface{}) {
	return
}

func (goj *GenericObject) PyDictSet(key, value interface{}) error {
	return nil
}

func (g *GenericClass) PyNew(args ...interface{}) (interface{}, error) {
	return &GenericObject{
		Class:           g,
		ConstructorArgs: args,
	}, nil
}

func (g *GenericClass) Call(args ...interface{}) (interface{}, error) {
	return g, nil
}
