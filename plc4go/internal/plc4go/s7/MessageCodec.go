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
package s7

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type MessageCodec struct {
	*spi.DefaultCodec
	expectationCounter int32
}

func NewMessageCodec(transportInstance transports.TransportInstance) *MessageCodec {
	codec := &MessageCodec{
		DefaultCodec:       spi.NewDefaultCodec(transportInstance),
		expectationCounter: 1,
	}
	codec.DefaultCodecRequiredInterface = codec
	return codec
}

func (m *MessageCodec) Send(message interface{}) error {
	log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	tpktPacket := model.CastTPKTPacket(message)
	// Serialize the request
	wb := utils.NewWriteBuffer()
	err := tpktPacket.Serialize(*wb)
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.TransportInstance.Write(wb.GetBytes())
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive() (interface{}, error) {
	log.Trace().Msg("receiving")
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.TransportInstance.GetNumReadableBytes(); (err == nil) && (num >= 6) {
		log.Debug().Msgf("we got %d readable bytes", num)
		data, err := m.TransportInstance.PeekReadableBytes(6)
		if err != nil {
			log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet
		// TODO: wrong size for s7
		packetSize := (uint32(data[4]) << 8) + uint32(data[5]) + 6
		if num < packetSize {
			log.Debug().Msgf("Not enough bytes. Got: %d Need: %d\n", num, packetSize)
			return nil, nil
		}
		data, err = m.TransportInstance.Read(packetSize)
		if err != nil {
			// TODO: Possibly clean up ...
			return nil, nil
		}
		rb := utils.NewReadBuffer(data)
		tcpAdu, err := model.COTPPacketParse(rb, uint16(packetSize))
		if err != nil {
			log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return tcpAdu, nil
	} else if err != nil {
		log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
