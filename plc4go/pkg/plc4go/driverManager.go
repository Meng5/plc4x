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
package plc4go

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
)

// This is the main entry point for PLC4Go applications
type PlcDriverManager interface {
	// Manually register a new driver
	RegisterDriver(driver PlcDriver)
	// List the names of all drivers registered in the system
	ListDriverNames() []string
	// Get access to a driver instance for a given driver-name
	GetDriver(driverName string) (PlcDriver, error)

	// Manually register a new driver
	RegisterTransport(transport transports.Transport)
	// List the names of all drivers registered in the system
	ListTransportNames() []string
	// Get access to a driver instance for a given driver-name
	GetTransport(transportName string, connectionString string, options map[string][]string) (transports.Transport, error)

	// Get a connection to a remote PLC for a given plc4x connection-string
	GetConnection(connectionString string) <-chan PlcConnectionConnectResult

	// Execute all available discovery methods on all available drivers using all transports
	Discover(func(event model.PlcDiscoveryEvent)) error
}

type PlcDriverManger struct {
	drivers    map[string]PlcDriver
	transports map[string]transports.Transport
}

func NewPlcDriverManager() PlcDriverManager {
	log.Trace().Msg("Creating plc driver manager")
	return PlcDriverManger{
		drivers:    map[string]PlcDriver{},
		transports: map[string]transports.Transport{},
	}
}

func (m PlcDriverManger) RegisterDriver(driver PlcDriver) {
	if driver == nil {
		panic("driver must not be nil")
	}
	log.Debug().Str("protocolName", driver.GetProtocolName()).Msg("Registering driver")
	// If this driver is already registered, just skip resetting it
	for driverName := range m.drivers {
		if driverName == driver.GetProtocolCode() {
			log.Warn().Str("protocolName", driver.GetProtocolName()).Msg("Already registered")
			return
		}
	}
	m.drivers[driver.GetProtocolCode()] = driver
	log.Info().Str("protocolName", driver.GetProtocolName()).Msgf("Driver for %s registered", driver.GetProtocolName())
}

func (m PlcDriverManger) ListDriverNames() []string {
	log.Trace().Msg("Listing driver names")
	var driverNames []string
	for driverName := range m.drivers {
		driverNames = append(driverNames, driverName)
	}
	log.Trace().Msgf("Found %d driver(s)", len(driverNames))
	return driverNames
}

func (m PlcDriverManger) GetDriver(driverName string) (PlcDriver, error) {
	if val, ok := m.drivers[driverName]; ok {
		return val, nil
	}
	return nil, errors.Errorf("couldn't find driver %s", driverName)
}

func (m PlcDriverManger) RegisterTransport(transport transports.Transport) {
	if transport == nil {
		panic("transport must not be nil")
	}
	log.Debug().Str("transportName", transport.GetTransportName()).Msg("Registering transport")
	// If this transport is already registered, just skip resetting it
	for transportName := range m.transports {
		if transportName == transport.GetTransportCode() {
			log.Warn().Str("transportName", transport.GetTransportName()).Msg("Transport already registered")
			return
		}
	}
	m.transports[transport.GetTransportCode()] = transport
	log.Info().Str("transportName", transport.GetTransportName()).Msgf("Transport for %s registered", transport.GetTransportName())
}

func (m PlcDriverManger) ListTransportNames() []string {
	log.Trace().Msg("Listing transport names")
	var transportNames []string
	for transportName := range m.transports {
		transportNames = append(transportNames, transportName)
	}
	log.Trace().Msgf("Found %d transports", len(transportNames))
	return transportNames
}

func (m PlcDriverManger) GetTransport(transportName string, _ string, _ map[string][]string) (transports.Transport, error) {
	if val, ok := m.transports[transportName]; ok {
		log.Debug().Str("transportName", transportName).Msg("Returning transport")
		return val, nil
	}
	return nil, errors.Errorf("couldn't find transport %s", transportName)
}

func (m PlcDriverManger) GetConnection(connectionString string) <-chan PlcConnectionConnectResult {
	log.Debug().Str("connectionString", connectionString).Msgf("Getting connection for %s", connectionString)
	// Parse the connection string.
	connectionUrl, err := url.Parse(connectionString)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing connection")
		ch := make(chan PlcConnectionConnectResult)
		go func() {
			ch <- NewPlcConnectionConnectResult(nil, errors.Wrap(err, "error parsing connection string"))
		}()
		return ch
	}
	log.Debug().Stringer("connectionUrl", connectionUrl).Msg("parsed connection URL")

	// The options will be used to configure both the transports as well as the connections/drivers
	configOptions := connectionUrl.Query()

	// Find the driver specified in the url.
	driverName := connectionUrl.Scheme
	driver, err := m.GetDriver(driverName)
	if err != nil {
		log.Err(err).Str("driverName", driverName).Msgf("Couldn't get driver for %s", driverName)
		ch := make(chan PlcConnectionConnectResult)
		go func() {
			ch <- NewPlcConnectionConnectResult(nil, errors.Wrap(err, "error getting driver for connection string"))
		}()
		return ch
	}

	// If a transport is provided alongside the driver, the URL content is decoded as "opaque" data
	// Then we have to re-parse that to get the transport code as well as the host & port information.
	var transportName string
	var transportConnectionString string
	if len(connectionUrl.Opaque) > 0 {
		log.Trace().Msg("we handling a opaque connectionUrl")
		connectionUrl, err := url.Parse(connectionUrl.Opaque)
		if err != nil {
			log.Err(err).Str("connectionUrl.Opaque", connectionUrl.Opaque).Msg("Couldn't get transport due to parsing error")
			ch := make(chan PlcConnectionConnectResult)
			go func() {
				ch <- NewPlcConnectionConnectResult(nil, errors.Wrap(err, "error parsing connection string"))
			}()
			return ch
		}
		transportName = connectionUrl.Scheme
		transportConnectionString = connectionUrl.Host
	} else {
		log.Trace().Msg("we handling a non-opaque connectionUrl")
		// If no transport was provided the driver has to provide a default transport.
		transportName = driver.GetDefaultTransport()
		transportConnectionString = connectionUrl.Host
	}
	log.Debug().
		Str("transportName", transportName).
		Str("transportConnectionString", transportConnectionString).
		Msgf("got a transport %s", transportName)
	// If no transport has been specified explicitly or per default, we have to abort.
	if transportName == "" {
		log.Error().Msg("got a empty transport")
		ch := make(chan PlcConnectionConnectResult)
		go func() {
			ch <- NewPlcConnectionConnectResult(nil, errors.New("no transport specified and no default defined by driver"))
		}()
		return ch
	}

	// Assemble a correct transport url
	transportUrl := url.URL{
		Scheme: transportName,
		Host:   transportConnectionString,
	}
	log.Debug().Stringer("transportUrl", &transportUrl).Msg("Assembled transport url")

	// Create a new connection
	return driver.GetConnection(transportUrl, m.transports, configOptions)
}

// TODO: Currently all network devices are used as well as all transports and all protocols. It would be cool if we had some sort of DiscoveryRequestBuilder instead of only this single method.
func (m PlcDriverManger) Discover(callback func(event model.PlcDiscoveryEvent)) error {
	for _, driver := range m.drivers {
		if driver.SupportsDiscovery() {
			err := driver.Discover(callback)
			if err != nil {
				return errors.Wrapf(err, "Error running Discover on driver %s", driver.GetProtocolName())
			}
		}
	}
	return nil
}
