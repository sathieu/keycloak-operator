package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ClientSecret(cr *v1alpha1.KeycloakClient) *v1.Secret {
	escapedSecretName := SanitizeResourceNameWithAlphaNum(ClientSecretName + "-" + cr.Spec.Client.ClientID)
	secret := &v1.Secret{
		ObjectMeta: v12.ObjectMeta{
			Name:      escapedSecretName,
			Namespace: cr.Namespace,
			Labels:    make(map[string]string),
		},
		Data: map[string][]byte{
			ClientSecretClientIDProperty:     []byte(cr.Spec.Client.ClientID),
			ClientSecretClientSecretProperty: []byte(cr.Spec.Client.Secret),
		},
	}
	if cr.Spec.Client.SecretTemplate != nil && cr.Spec.Client.SecretTemplate.Metadata != nil {
		if cr.Spec.Client.SecretTemplate.Metadata.Labels != nil {
			secret.Labels = cr.Spec.Client.SecretTemplate.Metadata.Labels
		}
		if cr.Spec.Client.SecretTemplate.Metadata.Annotations != nil {
			secret.Annotations = cr.Spec.Client.SecretTemplate.Metadata.Annotations
		}
	}
	secret.Labels["app"] = ApplicationName
	return secret
}

func ClientSecretSelector(cr *v1alpha1.KeycloakClient) client.ObjectKey {
	escapedSelectorName := SanitizeResourceNameWithAlphaNum(ClientSecretName + "-" + cr.Spec.Client.ClientID)
	return client.ObjectKey{
		Name:      escapedSelectorName,
		Namespace: cr.Namespace,
	}
}

func ClientSecretReconciled(cr *v1alpha1.KeycloakClient, currentState *v1.Secret) *v1.Secret {
	reconciled := currentState.DeepCopy()
	// Since the client is synced upon update, we always override what's there...
	reconciled.Data = map[string][]byte{
		ClientSecretClientIDProperty:     []byte(cr.Spec.Client.ClientID),
		ClientSecretClientSecretProperty: []byte(cr.Spec.Client.Secret),
	}
	return reconciled
}
