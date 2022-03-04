package schema

import (
	"github.com/grafana/thema"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ThemaSchema contains a Grafana schema where the canonical schema expression
// is made with Thema and CUE.
// TODO: figure out what fields should be here
type ThemaSchema struct {
	lineage        thema.Lineage
	groupName      string
	groupVersion   string
	openapiSchema  apiextensionsv1beta1.JSONSchemaProps
	runtimeObjects []runtime.Object
}

// NewThemaSchema
// TODO: support multiple versions. Should be possible, since versions are in the lineage.
func NewThemaSchema(
	lineage thema.Lineage,
	groupName, groupVersion string, // TODO: somehow figure this out from the lineage
	openapiSchema apiextensionsv1beta1.JSONSchemaProps, // TODO: should be part of the lineage
	resource, list runtime.Object,
) ThemaSchema {
	return ThemaSchema{
		lineage:        lineage,
		groupName:      groupName,
		groupVersion:   groupVersion,
		openapiSchema:  openapiSchema,
		runtimeObjects: []runtime.Object{resource, list},
	}
}

// Name returns the canonical string that identifies the object being schematized.
func (ts ThemaSchema) Name() string {
	return ts.lineage.Name()
}

// GroupName
func (ts ThemaSchema) GroupName() string {
	return ts.groupName
}

// GroupName
func (ts ThemaSchema) GroupVersion() string {
	return ts.groupVersion
}

// GetRuntimeObjects returns a runtime.Object that will accurately represent
// the authorial intent of the Thema lineage to Kubernetes.
func (ts ThemaSchema) RuntimeObjects() []runtime.Object {
	return ts.runtimeObjects
}

// OpenAPISchema
func (ts ThemaSchema) OpenAPISchema() apiextensionsv1beta1.JSONSchemaProps {
	return ts.openapiSchema
}
