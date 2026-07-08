/*
MISP-Operator - A Kubernetes operator for simplified deployments of MISP at scale.
Copyright (C) 2026 Pascal Iske

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package v1alpha1

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var mispinstancelog = logf.Log.WithName("mispinstance-resource")

// SetupMispInstanceWebhookWithManager registers the webhook for MispInstance in the manager.
func SetupMispInstanceWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &mispv1alpha1.MispInstance{}).
		WithValidator(&MispInstanceCustomValidator{}).
		WithDefaulter(&MispInstanceCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-misp-k8s-pascaliske-dev-v1alpha1-mispinstance,mutating=true,failurePolicy=fail,sideEffects=None,groups=misp.k8s.pascaliske.dev,resources=mispinstances,verbs=create;update,versions=v1alpha1,name=mmispinstance-v1alpha1.kb.io,admissionReviewVersions=v1

// MispInstanceCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind MispInstance when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type MispInstanceCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind MispInstance.
func (d *MispInstanceCustomDefaulter) Default(_ context.Context, obj *mispv1alpha1.MispInstance) error {
	mispinstancelog.Info("Defaulting for MispInstance", "name", obj.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-misp-k8s-pascaliske-dev-v1alpha1-mispinstance,mutating=false,failurePolicy=fail,sideEffects=None,groups=misp.k8s.pascaliske.dev,resources=mispinstances,verbs=create;update,versions=v1alpha1,name=vmispinstance-v1alpha1.kb.io,admissionReviewVersions=v1

// MispInstanceCustomValidator struct is responsible for validating the MispInstance resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type MispInstanceCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type MispInstance.
func (v *MispInstanceCustomValidator) ValidateCreate(_ context.Context, obj *mispv1alpha1.MispInstance) (admission.Warnings, error) {
	mispinstancelog.Info("Validation for MispInstance upon creation", "name", obj.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type MispInstance.
func (v *MispInstanceCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj *mispv1alpha1.MispInstance) (admission.Warnings, error) {
	mispinstancelog.Info("Validation for MispInstance upon update", "name", newObj.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type MispInstance.
func (v *MispInstanceCustomValidator) ValidateDelete(_ context.Context, obj *mispv1alpha1.MispInstance) (admission.Warnings, error) {
	mispinstancelog.Info("Validation for MispInstance upon deletion", "name", obj.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
