package main

import (
	"bytes"
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"text/template"
)

const (
	extendSwitchRolesBaseAccountTemplate = `
[{{.Name}}]
aws_account_id = {{.Data.baseAccountId}}

`
)

const (
	extendSwitchRolesTemplate = `
[{{.Name}}-{{.Data.role}}]
source_profile = {{.Data.sourceProfile}}
color = {{.Data.color}}
role_arn = arn:aws:iam::{{.Data.accountId}}:role/{{.Data.role}}

`
)

func (extension *ExtensionAwsExtendSwitchRoles) Exec(globalConfig schema.TerralessGlobalConfig, data schema.TerralessData) error {
	logger.Info(fmt.Sprintf("[%s] Executing", ExtensionName))
	buffer := bytes.Buffer{}

	for _, team := range globalConfig.Teams {
		buffer = renderTemplateToBuffer(team, buffer, extendSwitchRolesBaseAccountTemplate, "baseAccount")

		for _, provider := range team.Providers {
			if provider.Type == "aws" {
				provider.Data["sourceProfile"] = team.Name

				for _, role := range provider.Roles {
					provider.Data["role"] = role
					buffer = renderTemplateToBuffer(provider, buffer, extendSwitchRolesTemplate, "profile")
				}
			}
		}
	}

	logger.Info(fmt.Sprintf("AWS Extend Switch Roles configuration:\n%s\n", buffer.String()))

	return nil
}

func renderTemplateToBuffer(config interface{}, buffer bytes.Buffer, tpl string, name string) bytes.Buffer {
	tmpl := template.Must(template.New(name).Parse(tpl))
	err := tmpl.Execute(&buffer, config)

	if err != nil {
		fatal("Failed writing to Buffer: ", err)
	}

	return buffer
}
