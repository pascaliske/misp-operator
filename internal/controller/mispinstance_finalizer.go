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

package controller

import (
	"context"
	"time"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const mispInstanceFinalizer = "misp.k8s.pascaliske.dev/finalizer"

//nolint:unparam // kept for interface consistency
func (r *MispInstanceReconciler) handleFinalizer(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) (ctrl.Result, error) {
	if controllerutil.ContainsFinalizer(mispInstance, mispInstanceFinalizer) {
		// TODO: handle deletion of external dependencies

		// remove finalizer
		controllerutil.RemoveFinalizer(mispInstance, mispInstanceFinalizer)
		if err := r.Update(ctx, mispInstance); err != nil {
			return ctrl.Result{}, err
		}
	}

	// stop reconciliation as the object is being deleted
	return ctrl.Result{}, nil
}

func (r *MispInstanceReconciler) setupFinalizer(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) (ctrl.Result, error) {
	// setup finalizer
	controllerutil.AddFinalizer(mispInstance, mispInstanceFinalizer)
	if err := r.Update(ctx, mispInstance); err != nil {
		return ctrl.Result{}, err
	}

	// requeue after finalizer was added
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}
