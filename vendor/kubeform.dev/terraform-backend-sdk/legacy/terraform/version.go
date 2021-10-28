package terraform

import (
	"kubeform.dev/terraform-backend-sdk/version"
)

// Deprecated: Providers should use schema.Provider.TerraformVersion instead
func VersionString() string {
	return version.String()
}
