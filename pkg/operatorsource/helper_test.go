package operatorsource_test

import (
	"fmt"

	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func helperNewOperatorSource(namespace, name, endpointType string) *v1alpha1.OperatorSource {
	return &v1alpha1.OperatorSource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: fmt.Sprintf("%s/%s",
				v1alpha1.SchemeGroupVersion.Group, v1alpha1.SchemeGroupVersion.Version),
			Kind: v1alpha1.OperatorSourceKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},

		Spec: v1alpha1.OperatorSourceSpec{
			Type:     endpointType,
			Endpoint: "http://localhost:5000/cnr",
		},
	}
}

func helperNewCatalogSourceConfig(namespace, name string) *v1alpha1.CatalogSourceConfig {
	return &v1alpha1.CatalogSourceConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: fmt.Sprintf("%s/%s",
				v1alpha1.SchemeGroupVersion.Group, v1alpha1.SchemeGroupVersion.Version),
			Kind: v1alpha1.CatalogSourceConfigKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
