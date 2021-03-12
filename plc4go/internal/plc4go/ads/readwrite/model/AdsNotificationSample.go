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
	"encoding/base64"
	"encoding/xml"
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"io"
)

// The data-structure of this message
type AdsNotificationSample struct {
	NotificationHandle uint32
	SampleSize         uint32
	Data               []int8
	IAdsNotificationSample
}

// The corresponding interface
type IAdsNotificationSample interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

func NewAdsNotificationSample(notificationHandle uint32, sampleSize uint32, data []int8) *AdsNotificationSample {
	return &AdsNotificationSample{NotificationHandle: notificationHandle, SampleSize: sampleSize, Data: data}
}

func CastAdsNotificationSample(structType interface{}) *AdsNotificationSample {
	castFunc := func(typ interface{}) *AdsNotificationSample {
		if casted, ok := typ.(AdsNotificationSample); ok {
			return &casted
		}
		if casted, ok := typ.(*AdsNotificationSample); ok {
			return casted
		}
		return nil
	}
	return castFunc(structType)
}

func (m *AdsNotificationSample) GetTypeName() string {
	return "AdsNotificationSample"
}

func (m *AdsNotificationSample) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Simple field (notificationHandle)
	lengthInBits += 32

	// Simple field (sampleSize)
	lengthInBits += 32

	// Array field
	if len(m.Data) > 0 {
		lengthInBits += 8 * uint16(len(m.Data))
	}

	return lengthInBits
}

func (m *AdsNotificationSample) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func AdsNotificationSampleParse(io *utils.ReadBuffer) (*AdsNotificationSample, error) {

	// Simple Field (notificationHandle)
	notificationHandle, _notificationHandleErr := io.ReadUint32(32)
	if _notificationHandleErr != nil {
		return nil, errors.New("Error parsing 'notificationHandle' field " + _notificationHandleErr.Error())
	}

	// Simple Field (sampleSize)
	sampleSize, _sampleSizeErr := io.ReadUint32(32)
	if _sampleSizeErr != nil {
		return nil, errors.New("Error parsing 'sampleSize' field " + _sampleSizeErr.Error())
	}

	// Array field (data)
	// Count array
	data := make([]int8, sampleSize)
	for curItem := uint16(0); curItem < uint16(sampleSize); curItem++ {
		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'data' field " + _err.Error())
		}
		data[curItem] = _item
	}

	// Create the instance
	return NewAdsNotificationSample(notificationHandle, sampleSize, data), nil
}

func (m *AdsNotificationSample) Serialize(io utils.WriteBuffer) error {

	// Simple Field (notificationHandle)
	notificationHandle := uint32(m.NotificationHandle)
	_notificationHandleErr := io.WriteUint32(32, (notificationHandle))
	if _notificationHandleErr != nil {
		return errors.New("Error serializing 'notificationHandle' field " + _notificationHandleErr.Error())
	}

	// Simple Field (sampleSize)
	sampleSize := uint32(m.SampleSize)
	_sampleSizeErr := io.WriteUint32(32, (sampleSize))
	if _sampleSizeErr != nil {
		return errors.New("Error serializing 'sampleSize' field " + _sampleSizeErr.Error())
	}

	// Array Field (data)
	if m.Data != nil {
		for _, _element := range m.Data {
			_elementErr := io.WriteInt8(8, _element)
			if _elementErr != nil {
				return errors.New("Error serializing 'data' field " + _elementErr.Error())
			}
		}
	}

	return nil
}

func (m *AdsNotificationSample) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var token xml.Token
	var err error
	for {
		token, err = d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch token.(type) {
		case xml.StartElement:
			tok := token.(xml.StartElement)
			switch tok.Name.Local {
			case "notificationHandle":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.NotificationHandle = data
			case "sampleSize":
				var data uint32
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.SampleSize = data
			case "data":
				var _encoded string
				if err := d.DecodeElement(&_encoded, &tok); err != nil {
					return err
				}
				_decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
				_len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
				if err != nil {
					return err
				}
				m.Data = utils.ByteArrayToInt8Array(_decoded[0:_len])
			}
		}
	}
}

func (m *AdsNotificationSample) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	className := "org.apache.plc4x.java.ads.readwrite.AdsNotificationSample"
	if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
		{Name: xml.Name{Local: "className"}, Value: className},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.NotificationHandle, xml.StartElement{Name: xml.Name{Local: "notificationHandle"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.SampleSize, xml.StartElement{Name: xml.Name{Local: "sampleSize"}}); err != nil {
		return err
	}
	_encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(m.Data)))
	base64.StdEncoding.Encode(_encodedData, utils.Int8ArrayToByteArray(m.Data))
	if err := e.EncodeElement(_encodedData, xml.StartElement{Name: xml.Name{Local: "data"}}); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}