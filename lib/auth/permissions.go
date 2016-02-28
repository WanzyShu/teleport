/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"fmt"

	"github.com/gravitational/teleport"

	"github.com/gravitational/trace"
)

// PermissionChecker interface verifies that clients have permissions
// to execute any action of the auth server
type PermissionChecker interface {
	// HasPermission checks if the given role has a permission to execute
	// the action
	HasPermission(role teleport.Role, action string) error
}

// NewStandardPermissions returns permission checker with hardcoded roles
// that are built in when auth server starts in standard mode
func NewStandardPermissions() PermissionChecker {
	sp := standardPermissions{}
	sp.permissions = make(map[teleport.Role](map[string]bool))

	sp.permissions[teleport.RoleUser] = map[string]bool{
		ActionSignIn:             true,
		ActionCreateWebSession:   true,
		ActionGenerateUserCert:   true,
		ActionGetCertAuthorities: true,
		ActionGetSession:         true,
		ActionGetSessions:        true,
		ActionGetEvents:          true,
	}

	sp.permissions[teleport.RoleProvisionToken] = map[string]bool{
		ActionRegisterUsingToken:    true,
		ActionRegisterNewAuthServer: true,
	}

	sp.permissions[teleport.RoleNode] = map[string]bool{
		ActionUpsertServer:       true,
		ActionGetCertAuthorities: true,
		ActionGetLocalDomain:     true,
		ActionGetUserKeys:        true,
		ActionGetServers:         true,
		ActionUpsertParty:        true,
		ActionLogEntry:           true,
		ActionGetChunkWriter:     true,
	}

	sp.permissions[teleport.RoleWeb] = map[string]bool{
		//ActionSignIn:           true,
		ActionCreateWebSession: true,
		ActionGetWebSession:    true,
		ActionDeleteWebSession: true,
		ActionGetSession:       true,
		ActionGetSessions:      true,
		ActionGetEvents:        true,
	}

	sp.permissions[teleport.RoleSignup] = map[string]bool{
		ActionGetSignupTokenData:  true,
		ActionCreateUserWithToken: true,
	}

	return &sp
}

// NewHangoutPermissions is a set of permissions allowed to various
// roles when auth server is started in hangout mode on user's computer
func NewHangoutPermissions() PermissionChecker {
	sp := standardPermissions{}
	sp.permissions = make(map[teleport.Role](map[string]bool))

	sp.permissions[teleport.RoleUser] = map[string]bool{
		ActionGenerateUserCert:   true,
		ActionGetCertAuthorities: true,
	}

	sp.permissions[teleport.RoleProvisionToken] = map[string]bool{
		ActionRegisterUsingToken:    true,
		ActionRegisterNewAuthServer: true,
		ActionGenerateUserCert:      true,
	}

	sp.permissions[teleport.RoleHangoutRemoteUser] = map[string]bool{
		ActionGenerateUserCert: true,
	}

	sp.permissions[teleport.RoleNode] = map[string]bool{
		ActionUpsertServer:        true,
		ActionGetCertAuthorities:  true,
		ActionGetLocalDomain:      true,
		ActionGetUserKeys:         true,
		ActionGetServers:          true,
		ActionUpsertParty:         true,
		ActionLogEntry:            true,
		ActionGetChunkWriter:      true,
		ActionUpsertCertAuthority: true,
		ActionUpsertSession:       true,
		ActionGetAuthServers:      true,
	}

	sp.permissions[teleport.RoleWeb] = map[string]bool{
		ActionGetWebSession:    true,
		ActionDeleteWebSession: true,
	}

	sp.permissions[teleport.RoleSignup] = map[string]bool{
		ActionGetSignupTokenData:  true,
		ActionCreateUserWithToken: true,
	}

	return &sp
}

type standardPermissions struct {
	permissions map[teleport.Role](map[string]bool)
}

func (sp *standardPermissions) HasPermission(role teleport.Role, action string) error {
	if role == teleport.RoleAdmin {
		return nil
	}
	if permissions, ok := sp.permissions[role]; ok {
		if permissions[action] {
			return nil
		}
		return trace.Wrap(
			teleport.AccessDenied(
				fmt.Sprintf(
					"role '%v' doesn't have permission for action '%v'",
					role, action)))
	}
	return trace.Wrap(
		teleport.AccessDenied(
			fmt.Sprintf("role '%v' is not allowed", role)))
}

type allowAllPermissions struct {
}

func NewAllowAllPermissions() PermissionChecker {
	aap := allowAllPermissions{}
	return &aap
}

func (aap *allowAllPermissions) HasPermission(role teleport.Role, action string) error {
	return nil
}

var StandardRoles = []teleport.Role{
	teleport.RoleAuth,
	teleport.RoleUser,
	teleport.RoleWeb,
	teleport.RoleNode,
	teleport.RoleAdmin,
	teleport.RoleProvisionToken,
	teleport.RoleSignup,
}

var HangoutRoles = []teleport.Role{
	teleport.RoleAdmin,
	teleport.RoleProvisionToken,
	teleport.RoleHangoutRemoteUser,
}

const (
	ActionGetSessions                   = "GetSessions"
	ActionGetSession                    = "GetSession"
	ActionDeleteSession                 = "DeleteSession"
	ActionUpsertSession                 = "UpsertSession"
	ActionUpsertParty                   = "UpsertParty"
	ActionUpsertCertAuthority           = "UpsertCertAuthority"
	ActionGetCertAuthorities            = "GetCertAuthorities"
	ActionGetLocalDomain                = "GetLocalDomain"
	ActionDeleteCertAuthority           = "DeleteCertAuthority"
	ActionGenerateToken                 = "GenerateToken"
	ActionRegisterUsingToken            = "RegisterUsingToken"
	ActionRegisterNewAuthServer         = "RegisterNewAuthServer"
	ActionLog                           = "Log"
	ActionLogEntry                      = "LogEntry"
	ActionGetEvents                     = "GetEvents"
	ActionGetChunkWriter                = "GetChunkWriter"
	ActionGetChunkReader                = "GetChunkReader"
	ActionUpsertServer                  = "UpsertServer"
	ActionGetServers                    = "GetServers"
	ActionGetAuthServers                = "GetAuthServers"
	ActionUpsertPassword                = "UpsertPassword"
	ActionCheckPassword                 = "CheckPassword"
	ActionSignIn                        = "SignIn"
	ActionCreateWebSession              = "CreateWebSession"
	ActionGetWebSession                 = "GetWebSession"
	ActionGetWebSessionsKeys            = "GetWebSessionsKeys"
	ActionDeleteWebSession              = "DeleteWebSession"
	ActionGetUsers                      = "GetUsers"
	ActionDeleteUser                    = "DeleteUser"
	ActionUpsertUserKey                 = "UpsertUserKey"
	ActionGetUserKeys                   = "GetUserKeys"
	ActionDeleteUserKey                 = "DeleteUserKey"
	ActionGenerateKeyPair               = "GenerateKeyPair"
	ActionGenerateHostCert              = "GenerateHostCert"
	ActionGenerateUserCert              = "GenerateUserCert"
	ActionResetHostCertificateAuthority = "ResetHostCertificateAuthority"
	ActionResetUserCertificateAuthority = "ResetUserCertificateAuthority"
	ActionGenerateSealKey               = "GenerateSealKey"
	ActionGetSealKeys                   = "GetSeakKeys"
	ActionGetSealKey                    = "GetSealKey"
	ActionDeleteSealKey                 = "DeleteSealKey"
	ActionAddSealKey                    = "AddSealKey"
	ActionCreateSignupToken             = "CreateSignupToken"
	ActionGetSignupTokenData            = "GetSignupTokenData"
	ActionCreateUserWithToken           = "CreateUserWithToken"
	ActionUpsertUser                    = "UpsertUser"
)
