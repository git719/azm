## azm
`azm` is a [CLI](https://en.wikipedia.org/wiki/Command-line_interface) utility for managing [Indentity and Access Management (IAM)](https://www.nist.gov/identity-access-management) related Azure objects.


### Why?
Microsoft already has an official [Azure CLI tool](https://learn.microsoft.com/en-us/cli/azure/) called `az`, so **why** this?
- Only focuses on a smaller set of Azure object types related to IAM 
- This is written in [Go](https://go.dev/), so it's much faster than `az`
- Compiles into a binary executable that can easily be used from Windows, Linux, or macOS operating systems
- Automatic Linux [releases](https://github.com/git719/azm/releases/tag/v2.4.8) can easily be used in Github workflows
- Supports leveraging OIDC Github Action workflows with no passwords for configured Service Principal

**PS:** Every repo should have a **Why?** entry at the top 😊


### Functions
This is a little _Swiss Army knife_ that can very quickly perform the following functions:
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
    - Dump the current Resources or Security JWT tokens (see **pman** below)
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


### Quick Examples
1. List any Azure RBAC role, like the Built-in "Billing Reader" role for example:

```
$ azm -d "Billing Reader"
id: fa23ad8b-c56e-40d8-ac0c-ce449e1d2c64
properties:
  roleName: Billing Reader
  description: Allows read access to billing data
  assignableScopes:
    - /
  permissions:
    - actions:
        - Microsoft.Authorization/*/read
        - Microsoft.Billing/*/read
        - Microsoft.Commerce/*/read
        - Microsoft.Consumption/*/read
        - Microsoft.Management/managementGroups/read
        - Microsoft.CostManagement/*/read
        - Microsoft.Support/*
      notActions:
      dataActions:
      notDataActions:
```

- Another way of listing the same role is to call it by its UUID: `azm -d fa23ad8b-c56e-40d8-ac0c-ce449e1d2c64`
- The YAML listing format is more human-friendly and easier to read, and only displays the attributes that are most relevant to Azure systems engineers
- You can also display it in JSON format by calling: `azm -dj fa23ad8b-c56e-40d8-ac0c-ce449e1d2c64`
- One advantage of the JSON format is that it displays every single attribute in the Azure object

2. Add a secret to an Application object: 

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

3. Generate a random [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier). This can be very handy sometimes. Simply do: `azm -uuid`


### Usage
```
azm Azure IAM utility v2.4.8
    Read-Only Functions
    UUID                              Show object for given UUID
    -vs Specfile                      Compare YAML specfile to what's in Azure (only for d and a objects)
    -X[j] [Specifier]                 List all X objects tersely, with option for JSON output and/or match on Specifier
    -Xx                               Delete X object local file cache

      Where 'X' can be any of these object types:
      d  = RBAC Role Definitions   a  = RBAC Role Assignments   s  = Azure Subscriptions
      m  = Management Groups       u  = Azure AD Users          g  = Azure AD Groups
      sp = Service Principals      ap = Applications            ad = Azure AD Roles

    Read-Write Functions
    -rm[f] UUID|Specfile|"role name"  Delete role definition or assignment based on specifier; force option
    -up[f] Specfile                   Create or update definition or assignment based on YAML specfile; force option
    -kd[j]                            Create a skeleton role-definition.yaml specfile
    -ka[j]                            Create a skeleton role-assignment.yaml specfile
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
Instead of documenting individual examples of all of the above switches, it is best for you to play around with the utility to see the different listing functionality that it offers.


### pman
[pman](https://github.com/git719/azm/tree/main/pman) is a poor man's REST API Postman BASH script, which leverages `azm`'s `-tmg` and `-taz` arguments to get the current user's token to make other generic REST API calls against those 2 Microsoft APIs.


### Token Sample Code
Included in this repo are examples of how to get a Microsoft token with 3 different languages. The token can be for any API. The examples use Docker Compose.
1. [Python Example](https://github.com/git719/azm/tree/main/token-python)
2. [PowerShell Example](https://github.com/git719/azm/tree/main/token-powershell)
3. [Node Example](https://github.com/git719/azm/tree/main/token-node)


### Container
There is also a Docker Compose file and a Dockerfile for an example of how to use this program within a container.


### To-Do and Known Issues
The program is stable enough to be relied on as a small utility, but there are a number of little niggly things that could be improved. Will put a list together at some point.

At any rate, no matter how stable any code is, it is always worth remembering computer scientist [Tony Hoare](https://en.wikipedia.org/wiki/Tony_Hoare)'s famous quote:
> "Inside every large program is a small program struggling to get out."


### Coding Philosophy and Feedback
The primary goal of this utility is to serve as a study aid for coding Azure utilities in Go, as well as to serve as a quick _Swiss Army knife* utility for managin tenant IAM objects. If you look through the code I think you will find that is relatively straightforward. There is a deliberate effor to keep the code as clear as possible, and simple to understand and maintain.

Note that the bulk of the code is actually in the [maz](https://github.com/git719/maz) library, and other packages. Please visit that repo for more info.

This utility along with the required libraries are obviously very useful to me. I don't think I'm going to formalize the way to contribute to this project, but if you find it useful and have some improvement suggestion please let me know. Anyway, this is published as an open source project, so feel free to clone and use on your own, with proper attributions.
