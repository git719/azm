# zman
`zman` is a [CLI](https://en.wikipedia.org/wiki/Command-line_interface) utility for **managing** Azure objects. It builds upon the functionality of the [zls](https://github.com/git719/zls) utility, which is strictly a reader/listing program. This is also a little _Swiss Army knife_ that can very **quickly** do the following:

- Perform all the same **listing** functions that `zls` can do
- Delete/Create/Update the following [Azure Resources Services](https://que.tips/azure/#azure-resource-services) objects in your tenant:
  - RBAC Role Definitions
  - RBAC Role Assignments
- Can output a sample RBAC Role definition or assignment __specification file__ in either JSON or YAML, that can then be used to create a new role or assignment
- Update the following [Azure Security Services](https://que.tips/azure/#azure-security-services) objects:
  - Service Principals: Can only create or delete SP secrets (Cannot yet create SPs)
  - Applications: Can only create or delete App secrets (Cannot yet create Apps)
- Create a UUID

## Quick Example
A quick example is adding a secret to an Application object: 

```
$ zman -apas 51afab9e-0225-4c36-81f0-f42289c1a57a "My Secret"
App_Object_Id: 51afab9e-0225-4c36-81f0-f42289c1a57a
New_Secret_Id: 7c140771-c547-43f9-8525-d08bd234e267
New_Secret_Name: My Secret
New_Secret_Expiry: 2025-01-06
New_Secret_Text: 8p68Q~Ab7OxR2nj.YOrXtOLwq1BT4bDy6wNebaYn
```

As the **usage** section shows, the secret Expiry defaults to 366 days if none is given. 

- Note that you have to use the **Objectd ID**, not the App ID or Client ID of the application
- The name could have been nulled with `""`
- To remove above secret, you can simply do: `zman -aprs 51afab9e-0225-4c36-81f0-f42289c1a57a 7c140771-c547-43f9-8525-d08bd234e267`

Another quick example is generating a random [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier), which can always be handy. To do so, simply do: `zman -uuid`

## Introduction
Everything else that applies to the [zls](https://github.com/git719/zls) utility also applies to this `zman` utility.

## Usage
```
zman Azure Resource RBAC and MS Graph MANAGER v2.4.0
    MANAGER FUNCTIONS
    -rm UUID|Specfile|"role name"     Delete role definition or assignment based on specifier
    -up Specfile                      Create or update definition or assignment based on specfile (YAML or JSON)
    -kd[j]                            Create a skeleton role-definition.yaml specfile (JSON option)
    -ka[j]                            Create a skeleton role-assignment.yaml specfile (JSON option)
    -spas SP_UUID "name" [Expiry]     Add secret to SP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)
    -sprs SP_UUID SecretID            Remove secret from Service Principal
    -apas APP_UUID "name" [Expiry]    Add secret to APP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)
    -aprs APP_UUID SecretID           Remove secret from application
    -uuid                             Generate new UUID

    READER FUNCTIONS
    UUID                              Show object for given UUID
    -vs Specfile                      Compare YAML or JSON specfile to what's in Azure (only for d and a objects)
    -X[j] [Specifier]                 List all X objects tersely, with option for JSON output and/or match on Specifier
    -Xx                               Delete X object local file cache

      Where 'X' can be any of these object types:
      d  = RBAC Role Definitions   a  = RBAC Role Assignments   s  = Azure Subscriptions
      m  = Management Groups       u  = Azure AD Users          g  = Azure AD Groups
      sp = Service Principals      ap = Applications            ad = Azure AD Roles

    -xx                               Delete ALL cache local files
    -ar                               List all RBAC role assignments with resolved names
    -mt                               List Management Group and subscriptions tree
    -pags                             List all Azure AD Privileged Access Groups
    -st                               List local cache count and Azure count of all objects
    -tmg                              Dump current token string for MS Graph API
    -taz                              Dump current token string for Azure Resource API
    -tc "TokenString"                 Dump token claims

    -id                               Display configured login values
    -id TenantId Username             Set up user for interactive login
    -id TenantId ClientId Secret      Set up ID for automated login
    -tx                               Delete current configured login values and token
    -v                                Print this usage page
```
