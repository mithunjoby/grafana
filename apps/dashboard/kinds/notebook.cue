package kinds

import (
	v2beta2 "github.com/grafana/grafana/sdkkinds/dashboard/v2beta2"
)

// Notebook is served at v2beta2 while its schema is experimental, on leaf types copied from
// dashboard v2. See v2beta2/notebook_spec.cue.
notebookV2beta2: {
	kind:       "Notebook"
	pluralName: "Notebooks"
	validation: {
		operations: ["CREATE", "UPDATE"]
	}
	mutation: {
		operations: ["CREATE", "UPDATE"]
	}
	schema: {
		spec: v2beta2.NotebookSpec
	}
}
