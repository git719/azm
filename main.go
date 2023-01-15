// main.go

package main

import (
	"fmt"
	"github.com/git719/maz"
	"github.com/git719/utl"
	"os"
	"path/filepath"
)

const (
	prgname = "zman"
	prgver  = "0.1.0"
)

func PrintUsage() {
	fmt.Printf(prgname + " Azure Resource RBAC role definitions and assignments manager v" + prgver + "\n" +
		"    UUID                              List definition or assignment object given its UUID\n" +
		"    -vs Specfile                      Compare YAML or JSON specfile to what's in Azure\n" +
		"    -rm UUID|Specfile|\"role name\"     Delete definition or assignment based on specifier\n" +
		"    -up Specfile                      Create or update definition or assignment based on specfile (YAML or JSON)\n" +
		"    -kd[j]                            Create a skeleton role-definition.yaml specfile (JSON option)\n" +
		"    -ka[j]                            Create a skeleton role-assignment.yaml specfile (JSON option)\n" +
		"    -d[j] [Specifier]                 List all role definitions, with Specifier filter and JSON options\n" +
		"    -a[j] [Specifier]                 List all role assignments, with Specifier filter and JSON options\n" +
		"    -s[j] [Specifier]                 List all subscriptions, with Specifier filter and JSON options\n" +
		"\n" +
		"    -z                                Dump variables in running program\n" +
		"    -cr                               Dump values in credentials file\n" +
		"    -cr  TenantId ClientId Secret     Set up MSAL automated ClientId + Secret login\n" +
		"    -cri TenantId Username            Set up MSAL interactive browser popup login\n" +
		"    -tx                               Delete MSAL accessTokens cache file\n" +
		"    -v                                Print this usage page\n")
	os.Exit(0)
}

func SetupVariables(z *maz.Bundle) maz.Bundle {
	// Set up variable object struct
	*z = maz.Bundle{
		ConfDir:      "",
		CredsFile:    "credentials.yaml",
		TokenFile:    "accessTokens.json",
		TenantId:     "",
		ClientId:     "",
		ClientSecret: "",
		Interactive:  false,
		Username:     "",
		AuthorityUrl: "",
		MgToken:      "",
		MgHeaders:    map[string]string{},
		AzToken:      "",
		AzHeaders:    map[string]string{},
	}
	// Set up configuration directory
	z.ConfDir = filepath.Join(os.Getenv("HOME"), "."+prgname)
	if utl.FileNotExist(z.ConfDir) {
		if err := os.Mkdir(z.ConfDir, 0700); err != nil {
			panic(err.Error())
		}
	}
	return *z
}

func main() {
	//TestFunction()
	numberOfArguments := len(os.Args[1:]) // Not including the program itself
	if numberOfArguments < 1 || numberOfArguments > 4 {
		PrintUsage() // Don't accept less than 1 or more than 4 arguments
	}

	var z maz.Bundle
	z = SetupVariables(&z)

	switch numberOfArguments {
	case 1: // Process 1-argument requests
		arg1 := os.Args[1]
		// This first set of 1-arg requests do not require API tokens to be set up
		switch arg1 {
		case "-v":
			PrintUsage()
		case "-cr":
			maz.DumpCredentials(z)
		}
		z = maz.SetupApiTokens(&z) // The remaining 1-arg requests DO required API tokens to be set up
		switch arg1 {
		case "-xx":
			maz.RemoveCacheFile("all", z)
		case "-tx", "-dx", "-ax", "-sx", "-mx", "-ux", "-gx", "-spx", "-apx", "-adx":
			t := arg1[1 : len(arg1)-1] // Single out the object type
			maz.RemoveCacheFile(t, z)
		case "-dj", "-aj", "-sj", "-mj", "-uj", "-gj", "-spj", "-apj", "-adj":
			t := arg1[1 : len(arg1)-1]
			all := maz.GetObjects(t, "", false, z) // false means do not force Azure call, ok to use cache
			utl.PrintJson(all)                     // Print entire set in JSON
		case "-d", "-a", "-s", "-m", "-u", "-g", "-sp", "-ap", "-ad":
			t := arg1[1:]
			all := maz.GetObjects(t, "", false, z)
			for _, i := range all { // Print entire set tersely
				maz.PrintTersely(t, i)
			}
		case "-ar":
			maz.PrintRoleAssignmentReport(z)
		case "-mt":
			maz.PrintMgTree(z)
		case "-pags":
			maz.PrintPags(z)
		case "-st":
			maz.PrintCountStatus(z)
		case "-z":
			maz.DumpVariables(z)
		default:
			PrintUsage()
		}
	case 2: // Process 2-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		z = maz.SetupApiTokens(&z)
		switch arg1 {
		case "-vs":
			maz.CompareSpecfileToAzure(arg2, z)
		case "-dj", "-aj", "-sj", "-mj", "-uj", "-gj", "-spj", "-apj", "-adj":
			t := arg1[1 : len(arg1)-1] // Single out the object type
			if utl.ValidUuid(arg2) {   // Search/print single object, if it's valid UUID
				x := maz.GetObjectById(t, arg2, z)
				utl.PrintJson(x)
			} else {
				matchingObjects := maz.GetObjects(t, arg2, false, z)
				if len(matchingObjects) == 1 {
					utl.PrintJson(matchingObjects[0]) // Print single matching object in JSON
				} else if len(matchingObjects) > 1 {
					utl.PrintJson(matchingObjects) // Print all matching objects in JSON
				}
			}
		case "-d", "-a", "-s", "-m", "-u", "-g", "-sp", "-ap", "-ad":
			t := arg1[1:]            // Single out the object type
			if utl.ValidUuid(arg2) { // Search/print single object, if it's valid UUID
				x := maz.GetObjectById(t, arg2, z)
				maz.PrintObject(t, x, z)
			} else {
				matchingObjects := maz.GetObjects(t, arg2, false, z)
				if len(matchingObjects) == 1 {
					x := matchingObjects[0].(map[string]interface{})
					maz.PrintObject(t, x, z)
				} else if len(matchingObjects) > 1 {
					for _, i := range matchingObjects { // Print all matching object teresely
						x := i.(map[string]interface{})
						maz.PrintTersely(t, x)
					}
				}
			}
		default:
			PrintUsage()
		}
	case 3: // Process 3-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		switch arg1 {
		case "-cri":
			z.TenantId = arg2
			z.Username = arg3
			maz.SetupInterativeLogin(z)
		default:
			PrintUsage()
		}
	case 4: // Process 4-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		arg4 := os.Args[4]
		switch arg1 {
		case "-cr":
			z.TenantId = arg2
			z.ClientId = arg3
			z.ClientSecret = arg4
			maz.SetupAutomatedLogin(z)
		default:
			PrintUsage()
		}
	}
	os.Exit(0)
}

// if ( $args.Count -eq 1 ) {        # Process 1-argument requests
//     $arg1 = $args[0]
//     # These 1-arg requests don't need credentials and API tokens to be setup
//     if ( $arg1 -eq "-cr" ) {
//         DumpCredentials
//     } elseif ( $arg1 -eq "-tx" ) {
//         ClearTokenCache
//         exit
//     } elseif ( ($arg1 -eq "-kd") -or ($arg1 -eq "-kdj") -or ($arg1 -eq "-ka") -or ($arg1 -eq "-kaj") ) {
//         CreateSkeletonFile $arg1
//     } elseif ( $arg1 -eq "-v" ) {
//         PrintUsage
//     }
//     # The rest do need global credentials and API tokens available
//     SetupApiTokens
//     if ( ValidUuid $arg1 ) {
//         ShowObject $arg1
//     } elseif ( ($arg1 -eq "-dj") -or ($arg1 -eq "-aj") -or ($arg1 -eq "-sj") ) {
//         $t = $arg1.Substring(1,1)    # Get object type designator
//         $allObjects = GetAllAzObjects $t
//         PrintJson ($allObjects)
//         exit
//     } elseif ( ($arg1 -eq "-d") -or ($arg1 -eq "-a") -or ($arg1 -eq "-s") ) {
//         $t = $arg1.Substring(1,1)    # Get object type designator
//         PrintAllAzObjectsTersely $t
//         exit
//     } elseif ( $arg1 -eq "-z" ) {
//         DumpVariables
//     } else {
//         PrintUsage
//     }
// } elseif ( $args.Count -eq 2 ) {  # Process 2-argument requests
//     $arg1 = $args[0] ; $arg2 = $args[1]
//     SetupApiTokens
//     if ( $arg1 -eq "-vs" ) {
//         CompareSpecfile $arg2
//     } elseif ( $arg1 -eq "-rm" ) {
//         DeleteObject $arg2
//     } elseif ( $arg1 -eq "-up" ) {
//         UpsertAzObject $arg2  # Create or Update role definition or assignment
//     } elseif ( ($arg1 -eq "-dj") -or ($arg1 -eq "-aj") -or ($arg1 -eq "-sj") ) {
//         # Process request with JSON formatted output option
//         $t = $arg1.Substring(1,1)    # Get object type designator
//         $objects = GetMatching $t $arg2   # Get all matching objects
//         if ( $objects.Count -gt 1 ) {
//             PrintJson $objects
//         } elseif ( $objects.Count -gt 0 ) {
//             PrintJson $objects[0]
//         }
//         exit
//     } elseif ( ($arg1 -eq "-d") -or ($arg1 -eq "-a") -or ($arg1 -eq "-s") ) {
//         # Process request with reguarly, tersely formatted output option
//         $t = $arg1.Substring(1,1)    # Get object type designator
//         $objects = GetMatching $t $arg2
//         if ( $objects.Count -gt 1 ) {
//             foreach ($i in $objects) {
//                 PrintAzObjectTersely $t $i
//             }
//         } elseif ( $objects.Count -gt 0 ) {
//             PrintAzObject $t $objects[0]
//         }
//         exit
//     } else {
//         PrintUsage
//     }
// } elseif ( $args.Count -eq 3 ) {  # Process 3-argument requests
//     $arg1 = $args[0] ; $arg2 = $args[1] ; $arg3 = $args[2]
//     if ( $arg1 -eq "-cri" ) {
//         SetupInteractiveLogin $arg2 $arg3
//     } else {
//         PrintUsage
//     }
// } elseif ( $args.Count -eq 4 ) {  # Process 4-argument requests
//     $arg1 = $args[0] ; $arg2 = $args[1] ; $arg3 = $args[2] ; $arg4 = $args[3]
//     if ( $arg1 -eq "-cr" ) {
//         SetupAutomatedLogin $arg2 $arg3 $arg4
//     } else {
//         PrintUsage
//     }
// } else {
//     PrintUsage
// }
