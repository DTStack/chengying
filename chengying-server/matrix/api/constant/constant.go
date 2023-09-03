package constant

type TemplateFile struct {
	FileName string
	//it means that if the template is available before k8s version 1.15
	Old bool
}

type GenerateType string

var (
	RuntimeBinDir   = "/matrix/easyagent/dtstack-runtime/bin/"
	ClusterStoreDir = "/matrix/easyagent/dtstack-runtime/cluster/"
	ShellLogDir     = "/matrix/easyagent/shell_history/"

	TPL_SELF_BUILD = TemplateFile{
		FileName: "templates/selfbuild-component.yml",
		Old:      false,
	}
	TPL_SELF_BUILD_V1BETA1 = TemplateFile{
		FileName: "templates/selfbuild-component-v1beta1.yml",
		Old:      true,
	}
	TPL_CLUSTER_RESOURCE = TemplateFile{
		FileName: "templates/import-cluster-resource.yml",
		Old:      false,
	}
	TPL_CLUSTER_RESOURCE_V1BETA1 = TemplateFile{
		FileName: "templates/import-cluster-resource-v1beta1.yml",
		Old:      true,
	}
	TPL_NS_RESOURCE = TemplateFile{
		FileName: "templates/import-ns-resouce.yml",
		Old:      false,
	}
)

const (
	DATE_FORMAT     = "2006-01-02 15:04:05"
	TEMPLATE_SUFFIX = ".yaml"

	TYPE_SELF_BUILD        GenerateType = "type_self_build"
	TYPE_IMPORT_CLUSTER    GenerateType = "type_import_cluster"
	TYPE_IMPORT_CLUSTER_NS GenerateType = "type_import_cluster_ns"

	CLUSTER_TEMPLATE_DIR         = "/matrix/easyagent/dtstack-runtime/templates/"
	TEMPLATES_FILE_SERVER_PREFIX = "/easyagent/dtstack-runtime/templates/"

	NAMESPACE_VALID       = "valid"
	NAMESPACE_INVALID     = "invalid"
	NAMESPACE_NOT_CONNECT = "not_connect"
)
