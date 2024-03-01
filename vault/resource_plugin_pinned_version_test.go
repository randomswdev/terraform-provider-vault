// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-vault/internal/consts"
	"github.com/hashicorp/terraform-provider-vault/internal/provider"
	"github.com/hashicorp/terraform-provider-vault/testutil"
)

func TestPluginPinnedVersion(t *testing.T) {
	const (
		typ     = "auth"
		version = "v1.0.0"
	)

	destName := acctest.RandomWithPrefix("tf-plugin-pinned-version")

	resourceName := "vault_plugin_pinned_version.test"

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck: func() {
			testutil.TestAccPreCheck(t)
			SkipIfAPIVersionLT(t, testProvider.Meta(), provider.VaultVersion116)
		},
		Steps: []resource.TestStep{
			{
				Config: testPluginPinnedVersionConfig(typ, destName, version),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resourceName, consts.FieldType, typ),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, destName),
					resource.TestCheckResourceAttr(resourceName, consts.FieldVersion, version),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil),
		},
	})
}

func testPluginPinnedVersionConfig(pluginType, name, version string) string {
	sha256 := strings.Repeat("01234567", 8)

	ret := fmt.Sprintf(`
%s

resource "vault_plugin_pinned_version" "test" {
  type      = vault_plugin.test.type
  name      = vault_plugin.test.name
  version   = vault_plugin.test.version
}
`, testPluginConfig(pluginType, name, version, sha256, "command", `["--arg"]`, `["foo=bar"]`))

	return ret
}
