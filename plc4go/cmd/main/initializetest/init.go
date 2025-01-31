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
package initializetest

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func init() {
	onJenkins := os.Getenv("JENKINS_URL") != ""
	log.Logger = log.
		//// Enable below if you want to see the filenames
		//With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: onJenkins}).
		// TODO: set to INFO once log mining is done
		Level(zerolog.TraceLevel)
}
