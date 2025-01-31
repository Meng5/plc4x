//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
	"encoding/xml"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"io"
)

// Code generated by build-utils. DO NOT EDIT.

// The data-structure of this message
type ApduControlContainer struct {
	ControlApdu *ApduControl
	Parent      *Apdu
}

// The corresponding interface
type IApduControlContainer interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
	xml.Unmarshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ApduControlContainer) Control() uint8 {
	return 1
}

func (m *ApduControlContainer) InitializeParent(parent *Apdu, numbered bool, counter uint8) {
	m.Parent.Numbered = numbered
	m.Parent.Counter = counter
}

func NewApduControlContainer(controlApdu *ApduControl, numbered bool, counter uint8) *Apdu {
	child := &ApduControlContainer{
		ControlApdu: controlApdu,
		Parent:      NewApdu(numbered, counter),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastApduControlContainer(structType interface{}) *ApduControlContainer {
	castFunc := func(typ interface{}) *ApduControlContainer {
		if casted, ok := typ.(ApduControlContainer); ok {
			return &casted
		}
		if casted, ok := typ.(*ApduControlContainer); ok {
			return casted
		}
		if casted, ok := typ.(Apdu); ok {
			return CastApduControlContainer(casted.Child)
		}
		if casted, ok := typ.(*Apdu); ok {
			return CastApduControlContainer(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *ApduControlContainer) GetTypeName() string {
	return "ApduControlContainer"
}

func (m *ApduControlContainer) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Simple field (controlApdu)
	lengthInBits += m.ControlApdu.LengthInBits()

	return lengthInBits
}

func (m *ApduControlContainer) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ApduControlContainerParse(io *utils.ReadBuffer) (*Apdu, error) {

	// Simple Field (controlApdu)
	controlApdu, _controlApduErr := ApduControlParse(io)
	if _controlApduErr != nil {
		return nil, errors.Wrap(_controlApduErr, "Error parsing 'controlApdu' field")
	}

	// Create a partially initialized instance
	_child := &ApduControlContainer{
		ControlApdu: controlApdu,
		Parent:      &Apdu{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *ApduControlContainer) Serialize(io utils.WriteBuffer) error {
	ser := func() error {

		// Simple Field (controlApdu)
		_controlApduErr := m.ControlApdu.Serialize(io)
		if _controlApduErr != nil {
			return errors.Wrap(_controlApduErr, "Error serializing 'controlApdu' field")
		}

		return nil
	}
	return m.Parent.SerializeParent(io, m, ser)
}

func (m *ApduControlContainer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var token xml.Token
	var err error
	token = start
	for {
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			switch tok.Name.Local {
			case "controlApdu":
				var dt *ApduControl
				if err := d.DecodeElement(&dt, &tok); err != nil {
					return err
				}
				m.ControlApdu = dt
			}
		}
		token, err = d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (m *ApduControlContainer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeElement(m.ControlApdu, xml.StartElement{Name: xml.Name{Local: "controlApdu"}}); err != nil {
		return err
	}
	return nil
}

func (m ApduControlContainer) String() string {
	return string(m.Box("ApduControlContainer", utils.DefaultWidth*2))
}

func (m ApduControlContainer) Box(name string, width int) utils.AsciiBox {
	if name == "" {
		name = "ApduControlContainer"
	}
	boxes := make([]utils.AsciiBox, 0)
	boxes = append(boxes, utils.BoxAnything("ControlApdu", m.ControlApdu, width-2))
	return utils.BoxBox(name, utils.AlignBoxes(boxes, width-2), 0)
}
