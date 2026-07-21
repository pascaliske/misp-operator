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
	"fmt"
	"strconv"
	"strings"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"
	"github.com/pascaliske/misp-operator/internal/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsv1apply "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1apply "k8s.io/client-go/applyconfigurations/core/v1"
	metav1apply "k8s.io/client-go/applyconfigurations/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	RestartAnnotation         = "misp.k8s.pascaliske.dev/restart-at"
	InternalRestartAnnotation = "kubectl.kubernetes.io/restartedAt"
)

func (r *MispInstanceReconciler) createCoreContainer(mispInstance *mispv1alpha1.MispInstance) *corev1apply.ContainerApplyConfiguration {
	adminSecretName := mispInstance.GetNameWithSuffix("admin")

	if mispInstance.Spec.Admin != nil && mispInstance.Spec.Admin.CredentialsSecretRef != nil {
		adminSecretName = mispInstance.Spec.Admin.CredentialsSecretRef.Name
	}

	coreEnv := []*corev1apply.EnvVarApplyConfiguration{
		corev1apply.
			EnvVar().
			WithName("TZ").
			WithValue(mispInstance.Spec.TimeZone),
		corev1apply.
			EnvVar().
			WithName("BASE_URL").
			WithValue(mispInstance.Spec.BaseUrl),
		corev1apply.
			EnvVar().
			WithName("ADMIN_EMAIL").
			WithValueFrom(
				corev1apply.
					EnvVarSource().
					WithSecretKeyRef(
						corev1apply.
							SecretKeySelector().
							WithName(adminSecretName).
							WithKey("email"),
					),
			),
		corev1apply.
			EnvVar().
			WithName("ADMIN_PASSWORD").
			WithValueFrom(
				corev1apply.
					EnvVarSource().
					WithSecretKeyRef(
						corev1apply.
							SecretKeySelector().
							WithName(adminSecretName).
							WithKey("password"),
					),
			),
		corev1apply.
			EnvVar().
			WithName("ADMIN_KEY").
			WithValueFrom(
				corev1apply.
					EnvVarSource().
					WithSecretKeyRef(
						corev1apply.
							SecretKeySelector().
							WithName(adminSecretName).
							WithKey("apiKey"),
					),
			),
		corev1apply.
			EnvVar().
			WithName("MYSQL_HOST").
			WithValue(mispInstance.Spec.Database.Host),
		corev1apply.
			EnvVar().
			WithName("MYSQL_PORT").
			WithValue(strconv.Itoa(int(mispInstance.Spec.Database.Port))),
		corev1apply.
			EnvVar().
			WithName("MYSQL_DATABASE").
			WithValue(mispInstance.Spec.Database.Name),
		corev1apply.
			EnvVar().
			WithName("MYSQL_USER").
			WithValueFrom(
				corev1apply.
					EnvVarSource().
					WithSecretKeyRef(
						corev1apply.
							SecretKeySelector().
							WithName(mispInstance.Spec.Database.CredentialsSecretRef.Name).
							WithKey("username"),
					),
			),
		corev1apply.
			EnvVar().
			WithName("MYSQL_PASSWORD").
			WithValueFrom(
				corev1apply.
					EnvVarSource().
					WithSecretKeyRef(
						corev1apply.
							SecretKeySelector().
							WithName(mispInstance.Spec.Database.CredentialsSecretRef.Name).
							WithKey("password"),
					),
			),
		corev1apply.
			EnvVar().
			WithName("REDIS_HOST").
			WithValue(mispInstance.Spec.Cache.Host),
		corev1apply.
			EnvVar().
			WithName("REDIS_PORT").
			WithValue(strconv.Itoa(int(mispInstance.Spec.Cache.Port))),
		corev1apply.
			EnvVar().
			WithName("DISABLE_PRINTING_PLAINTEXT_CREDENTIALS").
			WithValue("true"),
	}

	// inject optional instance uuid variable
	if mispInstance.Spec.Uuid != "" {
		coreEnv = append(coreEnv,
			corev1apply.
				EnvVar().
				WithName("UUID").
				WithValue(mispInstance.Spec.Uuid),
		)
	}

	// inject optional admin organisation related variables
	if mispInstance.Spec.Admin != nil && mispInstance.Spec.Admin.Organisation != nil {
		if mispInstance.Spec.Admin.Organisation.Name != "" {
			coreEnv = append(coreEnv,
				corev1apply.
					EnvVar().
					WithName("ADMIN_ORG").
					WithValue(mispInstance.Spec.Admin.Organisation.Name),
			)
		}

		if mispInstance.Spec.Admin.Organisation.Uuid != "" {
			coreEnv = append(coreEnv,
				corev1apply.
					EnvVar().
					WithName("ADMIN_ORG_UUID").
					WithValue(mispInstance.Spec.Admin.Organisation.Uuid),
			)
		}
	}

	// inject optional redis related variables
	if mispInstance.Spec.Cache.PasswordSecretRef != nil && mispInstance.Spec.Cache.PasswordSecretRef.Name != "" {
		coreEnv = append(coreEnv,
			corev1apply.
				EnvVar().
				WithName("REDIS_PASSWORD").
				WithValueFrom(
					corev1apply.
						EnvVarSource().
						WithSecretKeyRef(
							corev1apply.
								SecretKeySelector().
								WithName(mispInstance.Spec.Cache.PasswordSecretRef.Name).
								WithKey("password"),
						),
				),
		)
	} else if mispInstance.Spec.Cache.EnableEmptyPassword {
		coreEnv = append(coreEnv,
			corev1apply.
				EnvVar().
				WithName("ENABLE_REDIS_EMPTY_PASSWORD").
				WithValue("true"),
		)
	}

	// inject optional oidc related variables
	if mispInstance.Spec.Oidc != nil && mispInstance.Spec.Oidc.Enabled {
		coreEnv = append(coreEnv,
			corev1apply.
				EnvVar().
				WithName("OIDC_ENABLE").
				WithValue("true"),
			corev1apply.
				EnvVar().
				WithName("OIDC_PROVIDER_URL").
				WithValueFrom(
					corev1apply.
						EnvVarSource().
						WithSecretKeyRef(
							corev1apply.
								SecretKeySelector().
								WithName(mispInstance.Spec.Oidc.CredentialsSecretRef.Name).
								WithKey("OIDC_PROVIDER_URL"),
						),
				),
			corev1apply.
				EnvVar().
				WithName("OIDC_CLIENT_ID").
				WithValueFrom(
					corev1apply.
						EnvVarSource().
						WithSecretKeyRef(
							corev1apply.
								SecretKeySelector().
								WithName(mispInstance.Spec.Oidc.CredentialsSecretRef.Name).
								WithKey("OIDC_CLIENT_ID"),
						),
				),
			corev1apply.
				EnvVar().
				WithName("OIDC_CLIENT_SECRET").
				WithValueFrom(
					corev1apply.
						EnvVarSource().
						WithSecretKeyRef(
							corev1apply.
								SecretKeySelector().
								WithName(mispInstance.Spec.Oidc.CredentialsSecretRef.Name).
								WithKey("OIDC_CLIENT_SECRET"),
						),
				),
			corev1apply.
				EnvVar().
				WithName("OIDC_LOGOUT_URL").
				WithValueFrom(
					corev1apply.
						EnvVarSource().
						WithSecretKeyRef(
							corev1apply.
								SecretKeySelector().
								WithName(mispInstance.Spec.Oidc.CredentialsSecretRef.Name).
								WithKey("OIDC_LOGOUT_URL"),
						),
				),
			corev1apply.
				EnvVar().
				WithName("OIDC_AUTH_METHOD").
				WithValue(mispInstance.Spec.Oidc.Settings.Method.String()),
			corev1apply.
				EnvVar().
				WithName("OIDC_SCOPES").
				WithValue(utils.ToJsonString(mispInstance.Spec.Oidc.Settings.Scopes)),
			corev1apply.
				EnvVar().
				WithName("OIDC_DEFAULT_ORG").
				WithValue(mispInstance.Spec.Oidc.Settings.DefaultOrg),
			corev1apply.
				EnvVar().
				WithName("OIDC_ROLES_PROPERTY").
				WithValue(mispInstance.Spec.Oidc.Settings.RolesProperty),
			corev1apply.
				EnvVar().
				WithName("OIDC_ROLES_MAPPING").
				WithValue(utils.ToJsonString(mispInstance.Spec.Oidc.Settings.RolesMapping)),
			corev1apply.
				EnvVar().
				WithName("OIDC_OFFLINE_ACCESS").
				WithValue(strconv.FormatBool(mispInstance.Spec.Oidc.Settings.OfflineAccess)),
			corev1apply.
				EnvVar().
				WithName("OIDC_MIXEDAUTH").
				WithValue(strconv.FormatBool(mispInstance.Spec.Oidc.Settings.MixedAuth)),
			corev1apply.
				EnvVar().
				WithName("OIDC_ALLOW_EMAIL_LINKING").
				WithValue(strconv.FormatBool(mispInstance.Spec.Oidc.Settings.AllowEmailLinking)),
			corev1apply.
				EnvVar().
				WithName("OIDC_REQUIRE_EMAIL_VERIFIED").
				WithValue(strconv.FormatBool(mispInstance.Spec.Oidc.Settings.RequireEmailVerified)),
			corev1apply.
				EnvVar().
				WithName("OIDC_DISABLE_REQUEST_OBJECT").
				WithValue(strconv.FormatBool(mispInstance.Spec.Oidc.Settings.DisableRequestObject)),
			corev1apply.
				EnvVar().
				WithName("OIDC_DISABLE_PUSHED_AUTHORIZATION_REQUEST").
				WithValue(strconv.FormatBool(mispInstance.Spec.Oidc.Settings.DisablePushedAuthorizationRequest)),
		)
	}

	// inject optional email related variables
	if mispInstance.Spec.Email != nil && mispInstance.Spec.Email.Smtp != nil {
		coreEnv = append(coreEnv,
			corev1apply.
				EnvVar().
				WithName("SMTP_FQDN").
				WithValue(mispInstance.Spec.Email.Smtp.Host),
		)

		if mispInstance.Spec.Email.Smtp.Port > 0 {
			coreEnv = append(coreEnv,
				corev1apply.
					EnvVar().
					WithName("SMTP_PORT").
					WithValue(strconv.Itoa(int(mispInstance.Spec.Email.Smtp.Port))),
			)
		}
	}

	// inject modules fqdn variable
	if mispInstance.Spec.Modules != nil && mispInstance.Spec.Modules.Enabled {
		coreEnv = append(coreEnv,
			corev1apply.
				EnvVar().
				WithName("MISP_MODULES_FQDN").
				WithValue(fmt.Sprintf("http://%s-modules", mispInstance.Name)),
		)
	}

	// inject custom extra environment variables
	if len(mispInstance.Spec.ExtraEnvs) > 0 {
		for _, e := range mispInstance.Spec.ExtraEnvs {
			coreEnv = append(coreEnv, utils.EnvVarToApplyConfiguration(e))
		}
	}

	return corev1apply.
		Container().
		WithName("misp-core").
		WithImage(mispInstance.GetCoreImage()).
		WithImagePullPolicy(utils.Try(mispInstance.Spec.ImagePullPolicy, corev1.PullIfNotPresent)).
		WithEnv(coreEnv...).
		WithVolumeMounts(
			corev1apply.
				VolumeMount().
				WithName("misp-storage").
				WithMountPath("/var/www/MISP/.gnupg").
				WithSubPath("gnupg"),
			corev1apply.
				VolumeMount().
				WithName("misp-storage").
				WithMountPath("/var/www/MISP/app/files").
				WithSubPath("files"),
			corev1apply.
				VolumeMount().
				WithName("misp-configs").
				WithMountPath("/var/www/MISP/app/Config"),
			corev1apply.
				VolumeMount().
				WithName("misp-logs").
				WithMountPath("/var/www/MISP/app/tmp/logs"),
		).
		WithSecurityContext(
			corev1apply.
				SecurityContext().
				WithAllowPrivilegeEscalation(false).
				// WithReadOnlyRootFilesystem(true).
				WithSeccompProfile(
					corev1apply.
						SeccompProfile().
						WithType(corev1.SeccompProfileTypeRuntimeDefault),
				),
		).
		WithStartupProbe(
			corev1apply.
				Probe().
				WithTCPSocket(corev1apply.TCPSocketAction().WithPort(intstr.FromInt(9002))).
				WithInitialDelaySeconds(15).
				WithPeriodSeconds(10).
				WithTimeoutSeconds(5).
				WithSuccessThreshold(1).
				WithFailureThreshold(30),
		).
		WithLivenessProbe(
			corev1apply.
				Probe().
				WithTCPSocket(corev1apply.TCPSocketAction().WithPort(intstr.FromInt(9002))).
				WithInitialDelaySeconds(15).
				WithPeriodSeconds(5).
				WithTimeoutSeconds(5).
				WithSuccessThreshold(1).
				WithFailureThreshold(2),
		).
		WithReadinessProbe(
			corev1apply.
				Probe().
				WithTCPSocket(corev1apply.TCPSocketAction().WithPort(intstr.FromInt(9002))).
				WithInitialDelaySeconds(5).
				WithPeriodSeconds(10).
				WithTimeoutSeconds(1).
				WithSuccessThreshold(1).
				WithFailureThreshold(3),
		).
		WithResources(
			corev1apply.
				ResourceRequirements().
				WithLimits(
					corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("2048Mi"),
					},
				).
				WithRequests(
					corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
				),
		)
}

func (r *MispInstanceReconciler) createNginxContainer(mispInstance *mispv1alpha1.MispInstance) *corev1apply.ContainerApplyConfiguration {
	nginxEnv := []*corev1apply.EnvVarApplyConfiguration{
		corev1apply.
			EnvVar().
			WithName("TZ").
			WithValue(mispInstance.Spec.TimeZone),
		corev1apply.
			EnvVar().
			WithName("BASE_URL").
			WithValue(mispInstance.Spec.BaseUrl),
		corev1apply.
			EnvVar().
			WithName("FASTCGI_LISTEN").
			WithValue("127.0.0.1:9002"),
		corev1apply.
			EnvVar().
			WithName("FASTCGI_LISTEN_STATUS").
			WithValue("127.0.0.1:9003"),
	}

	// inject optional nginx related variables
	if mispInstance.Spec.Nginx != nil {
		// inject optional client max body size
		if mispInstance.Spec.Nginx.ClientMaxBodySize != "" {
			nginxEnv = append(nginxEnv,
				corev1apply.
					EnvVar().
					WithName("NGINX_CLIENT_MAX_BODY_SIZE").
					WithValue(mispInstance.Spec.Nginx.ClientMaxBodySize),
			)
		}

		// inject optional forwarded headers
		if mispInstance.Spec.Nginx.ForwardedHeaders != nil && mispInstance.Spec.Nginx.ForwardedHeaders.Enabled {
			nginxEnv = append(nginxEnv,
				corev1apply.
					EnvVar().
					WithName("NGINX_X_FORWARDED_FOR").
					WithValue("true"),
			)

			if len(mispInstance.Spec.Nginx.ForwardedHeaders.TrustedProxies) > 0 {
				nginxEnv = append(nginxEnv,
					corev1apply.
						EnvVar().
						WithName("NGINX_SET_REAL_IP_FROM").
						WithValue(strings.Join(mispInstance.Spec.Nginx.ForwardedHeaders.TrustedProxies, ",")),
				)
			}
		}

		// inject optional security headers
		if mispInstance.Spec.Nginx.SecurityHeaders != nil {
			if mispInstance.Spec.Nginx.SecurityHeaders.FrameOptions != "" {
				nginxEnv = append(nginxEnv,
					corev1apply.
						EnvVar().
						WithName("NGINX_X_FRAME_OPTIONS").
						WithValue(mispInstance.Spec.Nginx.SecurityHeaders.FrameOptions),
				)
			}

			if mispInstance.Spec.Nginx.SecurityHeaders.ContentSecurityPolicy != "" {
				nginxEnv = append(nginxEnv,
					corev1apply.
						EnvVar().
						WithName("NGINX_CONTENT_SECURITY_POLICY").
						WithValue(mispInstance.Spec.Nginx.SecurityHeaders.ContentSecurityPolicy),
				)
			}

			if mispInstance.Spec.Nginx.SecurityHeaders.HstsMaxAge >= 0 {
				nginxEnv = append(nginxEnv,
					corev1apply.
						EnvVar().
						WithName("NGINX_HSTS_MAX_AGE").
						WithValue(strconv.Itoa(int(mispInstance.Spec.Nginx.SecurityHeaders.HstsMaxAge))),
				)
			}
		}

		// inject optional fastcgi values
		if mispInstance.Spec.Nginx.FastCGI != nil {
			if mispInstance.Spec.Nginx.FastCGI.ReadTimeout != "" {
				nginxEnv = append(nginxEnv,
					corev1apply.
						EnvVar().
						WithName("FASTCGI_READ_TIMEOUT").
						WithValue(mispInstance.Spec.Nginx.FastCGI.ReadTimeout),
				)
			}

			if mispInstance.Spec.Nginx.FastCGI.SendTimeout != "" {
				nginxEnv = append(nginxEnv,
					corev1apply.
						EnvVar().
						WithName("FASTCGI_SEND_TIMEOUT").
						WithValue(mispInstance.Spec.Nginx.FastCGI.SendTimeout),
				)
			}

			if mispInstance.Spec.Nginx.FastCGI.ConnectTimeout != "" {
				nginxEnv = append(nginxEnv,
					corev1apply.
						EnvVar().
						WithName("FASTCGI_CONNECT_TIMEOUT").
						WithValue(mispInstance.Spec.Nginx.FastCGI.ConnectTimeout),
				)
			}
		}
	}

	return corev1apply.
		Container().
		WithName("misp-nginx").
		WithImage(mispInstance.GetNginxImage()).
		WithImagePullPolicy(utils.Try(mispInstance.Spec.ImagePullPolicy, corev1.PullIfNotPresent)).
		WithPorts(
			corev1apply.
				ContainerPort().
				WithName("http").
				WithContainerPort(8080),
		).
		WithEnv(nginxEnv...).
		WithVolumeMounts(
			corev1apply.
				VolumeMount().
				WithName("misp-cache").
				WithMountPath("/var/www/MISP/cache"),
			corev1apply.
				VolumeMount().
				WithName("nginx-config").
				WithMountPath("/etc/nginx/conf.d"),
			corev1apply.
				VolumeMount().
				WithName("nginx-cache").
				WithMountPath("/var/cache/nginx"),
			corev1apply.
				VolumeMount().
				WithName("nginx-tmp").
				WithMountPath("/tmp"),
		).
		WithSecurityContext(
			corev1apply.
				SecurityContext().
				WithRunAsUser(101).
				WithRunAsGroup(101).
				WithAllowPrivilegeEscalation(false).
				WithReadOnlyRootFilesystem(true).
				WithSeccompProfile(
					corev1apply.
						SeccompProfile().
						WithType(corev1.SeccompProfileTypeRuntimeDefault),
				).
				WithCapabilities(
					corev1apply.
						Capabilities().
						WithDrop(
							corev1.Capability("ALL"),
						),
				),
		).
		WithReadinessProbe(
			corev1apply.
				Probe().
				WithTCPSocket(corev1apply.TCPSocketAction().WithPort(intstr.FromInt(8080))).
				WithInitialDelaySeconds(1).
				WithPeriodSeconds(5).
				WithTimeoutSeconds(5).
				WithSuccessThreshold(1).
				WithFailureThreshold(2),
		).
		WithResources(
			corev1apply.
				ResourceRequirements().
				WithLimits(
					corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("250m"),
						corev1.ResourceMemory: resource.MustParse("2048Mi"),
					},
				).
				WithRequests(
					corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
				),
		)
}

func (r *MispInstanceReconciler) createModulesContainer(mispInstance *mispv1alpha1.MispInstance) *corev1apply.ContainerApplyConfiguration {
	return corev1apply.
		Container().
		WithName("misp-modules").
		WithImage(mispInstance.GetModulesImage()).
		WithImagePullPolicy(utils.Try(mispInstance.Spec.ImagePullPolicy, corev1.PullIfNotPresent)).
		WithPorts(
			corev1apply.
				ContainerPort().
				WithName("http").
				WithContainerPort(6666),
		).
		WithEnv(
			corev1apply.
				EnvVar().
				WithName("TZ").
				WithValue(mispInstance.Spec.TimeZone),
		).
		WithReadinessProbe(
			corev1apply.
				Probe().
				WithTCPSocket(corev1apply.TCPSocketAction().WithPort(intstr.FromInt(6666))).
				WithInitialDelaySeconds(1).
				WithPeriodSeconds(5).
				WithTimeoutSeconds(5).
				WithSuccessThreshold(1).
				WithFailureThreshold(2),
		).
		WithResources(
			corev1apply.
				ResourceRequirements().
				WithLimits(
					corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("2048Mi"),
					},
				).
				WithRequests(
					corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("20m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				),
		)
}

func (r *MispInstanceReconciler) createInstanceDeployment(mispInstance *mispv1alpha1.MispInstance) *appsv1apply.DeploymentApplyConfiguration {
	podSpec := corev1apply.
		PodSpec().
		WithServiceAccountName(mispInstance.GetName()).
		WithContainers(
			r.createCoreContainer(mispInstance),
			r.createNginxContainer(mispInstance),
		).
		WithSecurityContext(
			corev1apply.
				PodSecurityContext().
				WithFSGroup(33),
		).
		WithVolumes(
			corev1apply.
				Volume().
				WithName("misp-storage").
				WithPersistentVolumeClaim(
					corev1apply.PersistentVolumeClaimVolumeSource().WithClaimName(mispInstance.GetNameWithSuffix("storage")),
				),
			corev1apply.
				Volume().
				WithName("misp-configs").
				WithEmptyDir(
					corev1apply.EmptyDirVolumeSource(),
				),
			corev1apply.
				Volume().
				WithName("misp-logs").
				WithEmptyDir(
					corev1apply.EmptyDirVolumeSource(),
				),
			corev1apply.
				Volume().
				WithName("misp-cache").
				WithEmptyDir(
					corev1apply.EmptyDirVolumeSource(),
				),
			corev1apply.
				Volume().
				WithName("nginx-config").
				WithEmptyDir(
					corev1apply.EmptyDirVolumeSource(),
				),
			corev1apply.
				Volume().
				WithName("nginx-cache").
				WithEmptyDir(
					corev1apply.EmptyDirVolumeSource(),
				),
			corev1apply.
				Volume().
				WithName("nginx-tmp").
				WithEmptyDir(
					corev1apply.EmptyDirVolumeSource(),
				),
		)

	// inject image pull secrets if supplied
	if len(mispInstance.Spec.ImagePullSecrets) > 0 {
		secrets := make([]*corev1apply.LocalObjectReferenceApplyConfiguration, 0, len(mispInstance.Spec.ImagePullSecrets))

		for _, secret := range mispInstance.Spec.ImagePullSecrets {
			secrets = append(secrets, corev1apply.LocalObjectReference().WithName(secret.Name))
		}

		podSpec = podSpec.WithImagePullSecrets(secrets...)
	}

	return appsv1apply.
		Deployment(mispInstance.Name, mispInstance.Namespace).
		WithLabels(utils.BuildAppLabels(mispInstance.Name, utils.AppLabelComponentMisp)).
		WithOwnerReferences(
			metav1apply.
				OwnerReference().
				WithAPIVersion(mispInstance.APIVersion).
				WithKind(mispInstance.Kind).
				WithName(mispInstance.Name).
				WithUID(mispInstance.UID).
				WithController(true).
				WithBlockOwnerDeletion(true),
		).
		WithSpec(
			appsv1apply.
				DeploymentSpec().
				WithStrategy(appsv1apply.DeploymentStrategy().WithType(appsv1.RecreateDeploymentStrategyType)).
				WithSelector(metav1apply.LabelSelector().WithMatchLabels(utils.BuildSelectorLabels(mispInstance.Name, utils.AppLabelComponentMisp))).
				WithTemplate(
					corev1apply.
						PodTemplateSpec().
						WithLabels(utils.BuildAppLabels(mispInstance.Name, utils.AppLabelComponentMisp)).
						WithSpec(podSpec),
				),
		)
}

func (r *MispInstanceReconciler) createModulesDeployment(mispInstance *mispv1alpha1.MispInstance) *appsv1apply.DeploymentApplyConfiguration {
	podSpec := corev1apply.
		PodSpec().
		WithServiceAccountName(mispInstance.GetNameWithSuffix("modules")).
		WithContainers(
			r.createModulesContainer(mispInstance),
		).
		WithSecurityContext(
			corev1apply.
				PodSecurityContext().
				WithFSGroup(33),
		)

	// inject image pull secrets if supplied
	if len(mispInstance.Spec.ImagePullSecrets) > 0 {
		secrets := make([]*corev1apply.LocalObjectReferenceApplyConfiguration, 0, len(mispInstance.Spec.ImagePullSecrets))

		for _, secret := range mispInstance.Spec.ImagePullSecrets {
			secrets = append(secrets, corev1apply.LocalObjectReference().WithName(secret.Name))
		}

		podSpec = podSpec.WithImagePullSecrets(secrets...)
	}

	return appsv1apply.
		Deployment(mispInstance.GetNameWithSuffix("modules"), mispInstance.Namespace).
		WithLabels(utils.BuildAppLabels(mispInstance.Name, utils.AppLabelComponentModules)).
		WithOwnerReferences(
			metav1apply.
				OwnerReference().
				WithAPIVersion(mispInstance.APIVersion).
				WithKind(mispInstance.Kind).
				WithName(mispInstance.Name).
				WithUID(mispInstance.UID).
				WithController(true).
				WithBlockOwnerDeletion(true),
		).
		WithSpec(
			appsv1apply.
				DeploymentSpec().
				WithStrategy(appsv1apply.DeploymentStrategy().WithType(appsv1.RecreateDeploymentStrategyType)).
				WithSelector(metav1apply.LabelSelector().WithMatchLabels(utils.BuildSelectorLabels(mispInstance.Name, utils.AppLabelComponentModules))).
				WithTemplate(
					corev1apply.
						PodTemplateSpec().
						WithLabels(utils.BuildAppLabels(mispInstance.Name, utils.AppLabelComponentModules)).
						WithSpec(podSpec),
				),
		)
}

func (r *MispInstanceReconciler) reconcileDeployments(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) error {
	options := &client.ApplyOptions{
		FieldManager: applyFieldManagerKeyInstance,
		Force:        new(true),
	}

	// reconcile misp modules deployment if enabled
	if mispInstance.Spec.Modules != nil && mispInstance.Spec.Modules.Enabled {
		if err := r.Apply(ctx, r.createModulesDeployment(mispInstance), options); err != nil {
			return err
		}
	}

	// reconcile misp instance deployment
	return r.Apply(ctx, r.createInstanceDeployment(mispInstance), options)
}
