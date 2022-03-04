package k8sbridge

import (
	"fmt"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/grafana/grafana/pkg/schema"
)

// Clientset
type Clientset struct {
	*kubernetes.Clientset
	clients map[k8schema.GroupVersion]*rest.RESTClient
}

// NewClientset
func NewClientset(cfg *rest.Config, schemas schema.CoreSchemaList) (*Clientset, error) {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	clients := make(map[k8schema.GroupVersion]*rest.RESTClient, len(schemas))

	for _, s := range schemas {
		ver := k8schema.GroupVersion{
			Group:   s.GroupName(),
			Version: s.GroupVersion(),
		}

		c := *cfg
		c.NegotiatedSerializer = clientscheme.Codecs.WithoutConversion()
		c.GroupVersion = &ver

		cli, err := rest.RESTClientFor(&c)
		if err != nil {
			return nil, err
		}

		clients[ver] = cli
	}

	return &Clientset{
		clientset,
		clients,
	}, nil
}

// ClientForVersion
func (c *Clientset) ClientForSchema(schema schema.ObjectSchema) (*rest.RESTClient, error) {
	k := k8schema.GroupVersion{
		Group:   schema.GroupName(),
		Version: schema.GroupVersion(),
	}

	v, ok := c.clients[k]
	if !ok {
		return nil, fmt.Errorf("no client registered for schema: %s/%s", schema.GroupName(), schema.GroupVersion())
	}

	return v, nil
}

// NewCRD
// TODO: use these to automatically register CRDs to the server.
func NewCRD(
	objectKind, groupName, groupVersion string, schema apiextensionsv1beta1.JSONSchemaProps,
) apiextensionsv1beta1.CustomResourceDefinition {
	return apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", groupName, groupVersion),
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   groupName,
			Version: groupVersion,
			Scope:   apiextensionsv1beta1.NamespaceScoped, // TODO: make configurable?
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural:   objectKind + "s", // TODO: figure out better approach?
				Singular: objectKind,
				Kind:     objectKind,
			},
			Validation: &apiextensionsv1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &schema,
			},
		},
	}
}
