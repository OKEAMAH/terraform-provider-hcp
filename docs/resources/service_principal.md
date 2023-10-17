---
page_title: "Resource hcp_service_principal - terraform-provider-hcp"
subcategory: "Cloud Platform"
description: |-
  The service principal resource manages a HCP Service Principal.
  The user or service account that is running Terraform when creating a hcp_service_principal resource must have roles/Admin on the parent resource; either the project or organization.
---

# hcp_service_principal (Resource)

The service principal resource manages a HCP Service Principal.

The user or service account that is running Terraform when creating a `hcp_service_principal` resource must have `roles/Admin` on the parent resource; either the project or organization.

## Example Usage: Create in provider configured project

```terraform
resource "hcp_service_principal" "example" {
  name = "example-sp"
}
```

## Example Usage: Create in new project

```terraform
resource "hcp_project" "my_proj" {
  name = "example"
}

resource "hcp_service_principal" "example" {
  name   = "example-sp"
  parent = hcp_project.my_proj.resource_name
}
```

## Example Usage: Create organization service principal

```terraform
data "hcp_organization" "my_org" {
}

resource "hcp_service_principal" "example" {
  name   = "example-sp"
  parent = data.hcp_organization.my_org.resource_name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The service principal's name.

### Optional

- `parent` (String) The parent location to create the service principal under. If unspecified, the service principal will be created in the project the provider is configured with. If specified, the accepted values are "project/<project_id>" or "organization/<organization_id>"

### Read-Only

- `resource_id` (String) The service principal's unique identitier
- `resource_name` (String) The service principal's resource name in the format `iam/project/<project_id>/service-principal/<name>` or `iam/organization/<organization_id>/service-principal/<name>`

## Import

Import is supported using the following syntax:

```shell
# Service Principals can be imported by specifying the resource name
terraform import hcp_service_principal.example iam/project/840e3701-55b6-4f86-8c17-b1fe397303c5/service-principal/my-sp
```