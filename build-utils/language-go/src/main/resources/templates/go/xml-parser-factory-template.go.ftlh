<#--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

	  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
-->
<#-- Prevent freemarker from escaping stuff -->
<#outputformat "undefined">
<#-- Declare the name and type of variables passed in to the template -->
<#-- @ftlvariable name="languageName" type="java.lang.String" -->
<#-- @ftlvariable name="protocolName" type="java.lang.String" -->
<#-- @ftlvariable name="outputFlavor" type="java.lang.String" -->
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.go.GoLanguageTemplateHelper" -->
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
<#-- Declare the name and type of variables declared locally inside the template -->
<#-- @ftlvariable name="arrayField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField" -->
<#-- @ftlvariable name="checksumField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ChecksumField" -->
<#-- @ftlvariable name="constField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ConstField" -->
<#-- @ftlvariable name="discriminatorField" type="org.apache.plc4x.plugins.codegenerator.types.fields.DiscriminatorField" -->
<#-- @ftlvariable name="enumField" type="org.apache.plc4x.plugins.codegenerator.types.fields.EnumField" -->
<#-- @ftlvariable name="implicitField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ImplicitField" -->
<#-- @ftlvariable name="manualArrayField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ManualArrayField" -->
<#-- @ftlvariable name="manualField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ManualField" -->
<#-- @ftlvariable name="optionalField" type="org.apache.plc4x.plugins.codegenerator.types.fields.OptionalField" -->
<#-- @ftlvariable name="paddingField" type="org.apache.plc4x.plugins.codegenerator.types.fields.PaddingField" -->
<#-- @ftlvariable name="reservedField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ReservedField" -->
<#-- @ftlvariable name="simpleField" type="org.apache.plc4x.plugins.codegenerator.types.fields.SimpleField" -->
<#-- @ftlvariable name="switchField" type="org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField" -->
<#-- @ftlvariable name="virtualField" type="org.apache.plc4x.plugins.codegenerator.types.fields.VirtualField" -->
<#-- @ftlvariable name="simpleTypeReference" type="org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference" -->
<#-- @ftlvariable name="complexTypeReference" type="org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference" -->
${helper.fileName(protocolName, languageName, outputFlavor)?replace(".", "/")}/XmlParserHelper.go
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
package ${outputFlavor?replace("-","")}

import (
	"encoding/xml"
	"github.com/apache/plc4x/plc4go/internal/plc4go/${protocolName}/${outputFlavor?replace("-","")}/model"
	"github.com/pkg/errors"
)

// Code generated by build-utils. DO NOT EDIT.

type ${protocolName?cap_first}XmlParserHelper struct {
}

func (m ${protocolName?cap_first}XmlParserHelper) Parse(typeName string, xmlString string) (interface{}, error) {
	switch typeName {
<#list helper.getComplexTypeRootDefinitions() as typeDefinition>
	case "${typeDefinition.name}":
		var obj *model.${helper.getLanguageTypeNameForTypeReference(typeDefinition.typeReference)}
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.Wrap(err ,"error unmarshalling xml")
		}
		return obj, nil
</#list>
	}
	return nil, errors.Errorf("Unsupported type %s", typeName)
}
</#outputformat>