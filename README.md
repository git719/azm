# azm
`azm` is a [CLI](https://en.wikipedia.org/wiki/Command-line_interface) utility for managing [Indentity and Access Management (IAM)](https://www.nist.gov/identity-access-management) related Azure objects. It is a little _Swiss Army knife_ that can very quickly do the following:

1. Read-Only Functions
    > **Note**<br>
    Of course these Read-Only functions are *only* available if you setup the tool to logon with an account that has the respective Read-Only privileges.
    - List the following [Azure Resources Services](https://que.tips/azure/#azure-resource-services) objects in your tenant:
        - RBAC Role Definitions
        - RBAC Role Assignments
        - Azure Subscriptions
        - Azure Management Groups
    - List the following [Azure Security Services](https://que.tips/azure/#azure-security-services) objects:
        - Azure AD Users
        - Azure AD Groups
        - Applications
        - Service Principals
        - Azure AD Roles that have been **activated**
        - Azure AD Roles standard definitions
    - Compare RBAC role definitions and assignments that are defined in a YAML __specification file__ to what that object currently looks like in the Azure tenant.
    - Dump the current Resources or Security JWT token being used (which can be used as a [simple Azure REST API caller](https://github.com/git719/azm/tree/main/pman) for testing purposes) 
    - Perform *many* other related listing functions.
2. Read-Write Functions
    > **Note**<br>
    Again, these Read-Write functions are **only** available if you setup the tool to logon with an account that has the respective Read-Write privileges
    - Delete/Create/Update the following [Azure Resources Services](https://que.tips/azure/#azure-resource-services) objects in your tenant:
        - RBAC Role Definitions
        - RBAC Role Assignments
    - Can output a sample RBAC Role definition or assignment YAML __specification file__, that can then be used to create a new role or assignment
    - Update the following [Azure Security Services](https://que.tips/azure/#azure-security-services) objects:
        - Service Principals: Can only create or delete SP secrets (Cannot yet create SPs)
        - Applications: Can only create or delete App secrets (Cannot yet create Apps)
    - Create a UUID

## Quick Example
A quick example is adding a secret to an Application object: 

```
$ azm -apas 51afab9e-0225-4c36-81f0-f42289c1a57a "My Secret"
App_Object_Id: 51afab9e-0225-4c36-81f0-f42289c1a57a
New_Secret_Id: 7c140771-c547-43f9-8525-d08bd234e267
New_Secret_Name: My Secret
New_Secret_Expiry: 2025-01-06
New_Secret_Text: 8p68Q~Ab7OxR2nj.YOrXtOLwq1BT4bDy6wNebaYn
```

As the **usage** section shows, the secret Expiry defaults to 366 days if none is given. 

- Note that you have to use the **Objectd ID**, not the App ID (Client ID) of the application
- The name could have been nulled with `""`
- To remove above secret, you can simply do: `azm -aprs 51afab9e-0225-4c36-81f0-f42289c1a57a 7c140771-c547-43f9-8525-d08bd234e267`

Another quick example is generating a random [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier), which can always be handy. To do so, simply do: `azm -uuid`

## Usage
```
azm Azure IAM utility v2.4.8
    Read-Only Functions
    UUID                              Show object for given UUID
    -vs Specfile                      Compare specfile (YAML or JSON) to what's in Azure (only for d and a objects)
    -X[j] [Specifier]                 List all X objects tersely, with option for JSON output and/or match on Specifier
    -Xx                               Delete X object local file cache

      Where 'X' can be any of these object types:
      d  = RBAC Role Definitions   a  = RBAC Role Assignments   s  = Azure Subscriptions
      m  = Management Groups       u  = Azure AD Users          g  = Azure AD Groups
      sp = Service Principals      ap = Applications            ad = Azure AD Roles

    Read-Write Functions
    -rm[f] UUID|Specfile|"role name"  Delete role definition or assignment based on specifier; force option
    -up[f] Specfile                   Create or update definition or assignment based on specfile (YAML or JSON); force option
    -kd[j]                            Create a skeleton role-definition.yaml specfile (JSON option)
    -ka[j]                            Create a skeleton role-assignment.yaml specfile (JSON option)
    -spas SP_UUID "name" [Expiry]     Add secret to SP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)
    -sprs SP_UUID SecretID            Remove secret from Service Principal
    -apas APP_UUID "name" [Expiry]    Add secret to APP; Expiry in YYYY-MM-DD format or X number of days (defaults to 366)
    -aprs APP_UUID SecretID           Remove secret from application
    -uuid                             Generate new UUID

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
