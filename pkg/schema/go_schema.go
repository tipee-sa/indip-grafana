package schema

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GoSchema contains a Grafana schema where the canonical schema expression is made
// with Go types, in traditional Kubernetes style.
// TODO figure out what fields should be here
type GoSchema struct {
	objectName     string
	groupName      string
	groupVersion   string
	openapiSchema  apiextensionsv1beta1.JSONSchemaProps
	runtimeObjects []runtime.Object
}

// NewGoSchema
// TODO: support multiple versions.
func NewGoSchema(
	objectKind, groupName, groupVersion string,
	openapiSchema apiextensionsv1beta1.JSONSchemaProps,
	resource, list runtime.Object,
) GoSchema {
	return GoSchema{
		objectName:     objectKind,
		groupName:      groupName,
		groupVersion:   groupVersion,
		openapiSchema:  openapiSchema,
		runtimeObjects: []runtime.Object{resource, list},
	}
}

// Name returns the canonical string that identifies the object being schematized.
func (gs GoSchema) Name() string {
	return gs.objectName
}

// GroupName
func (gs GoSchema) GroupName() string {
	return gs.groupName
}

// GroupName
func (gs GoSchema) GroupVersion() string {
	return gs.groupVersion
}

// GetRuntimeObjects returns a runtime.Object for this object kind.
func (gs GoSchema) GetRuntimeObjects() []runtime.Object {
	return gs.runtimeObjects
}

// OpenAPISchema
func (gs GoSchema) OpenAPISchema() apiextensionsv1beta1.JSONSchemaProps {
	return gs.openapiSchema
}
