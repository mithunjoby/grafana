package v2beta2

// Notebook-only schema types. They reuse the dashboard leaf types (PanelKind, LibraryPanelKind,
// TimeSettingsSpec, ElementReference) copied into this package from dashboard v2, WITHOUT being
// added to DashboardSpec's own element/layout unions. "Share the leaf types, diverge on the
// composition" — the dashboard schema never learns about Cell or NotebookLayout, and a notebook
// spec stays 1:1 with a dashboard v2 spec apart from the layout, so no NotebookSpec → DashboardSpec
// bridge is needed.
//
// Notebook is served at v2beta2 rather than v2 because the schema is still experimental: v2 is GA
// and must not take breaking schema changes. dashboard_spec.cue in this package is a generated copy
// of the v2 one (see the Makefile), because the SDK only resolves references within the CUE package
// named after the manifest version.

// A cell holds non-panel narrative content (markdown text, code) in a notebook layout.
// Panel cells are not represented here — they reuse PanelKind.
CellKind: {
	kind: "Cell"
	spec: CellSpec
}

CellSpec: {
	content: CellContentKind
}

// Pluggable cell content discriminated by `kind`. New content types are added
// by extending this union with another <Name>CellContentKind member.
CellContentKind: MarkdownCellContentKind | CodeCellContentKind

MarkdownCellContentKind: {
	kind: "Markdown"
	spec: MarkdownCellContentSpec
}

MarkdownCellContentSpec: {
	text: string
}

CodeCellContentKind: {
	kind: "Code"
	spec: CodeCellContentSpec
}

CodeCellContentSpec: {
	language: string
	code:     string
	highlight?: [...int]
	annotation?: string
}

NotebookLayoutKind: {
	kind: "NotebookLayout"
	spec: NotebookLayoutSpec
}

NotebookLayoutSpec: {
	cells: [...NotebookLayoutItemKind]
}

NotebookLayoutItemKind: {
	kind: "NotebookLayoutItem"
	spec: NotebookLayoutItemSpec
}

// One ordered item in a notebook layout. `element` references either a CellKind
// (markdown/code content) or a PanelKind in the notebook's elements map. `source`
// records who authored the cell; `collapsed` hides the body in the UI.
NotebookLayoutItemSpec: {
	element:    ElementReference
	source:     "assistant" | "user"
	collapsed?: bool
}

// A notebook element is a narrative cell, a panel, or a library panel. Unlike the dashboard
// Element union, this one includes CellKind — and it is referenced ONLY by NotebookSpec.
// CellKind is listed first so it is the generated default (a notebook is narrative-first).
NotebookElement: CellKind | PanelKind | LibraryPanelKind

// A notebook spec is a dashboard spec with narrative content. It has a title, optional description, tags, time
// settings, and a map of elements (panels and cells) referenced by the layout.
NotebookSpec: {
	title:        string
	description?: string
	tags: [...string] | *[]
	timeSettings: TimeSettingsSpec
	elements: [string]: NotebookElement
	layout: NotebookLayoutKind
}
