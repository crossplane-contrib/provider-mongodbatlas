// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha3

import v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"

// SetPasswordSecretRef sets the PasswordSecretRef field on the User's forProvider spec.
// This implements common.PasswordSecretRefSetter, allowing the PasswordGenerator
// initializer to point passwordSecretRef at the auto-generated secret.
func (u *User) SetPasswordSecretRef(ref *v1.LocalSecretKeySelector) {
	u.Spec.ForProvider.PasswordSecretRef = ref
}
