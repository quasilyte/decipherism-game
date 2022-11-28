package main

type componentSchema struct {
	entry *schemaElem

	elems []*schemaElem

	numKeywords     int
	keywords        []string
	encodedKeywords []string

	hasCondTransform bool
	hasPolygraphic   bool
	hasAtbash        bool
	hasRot13         bool
	hasIncDec        bool
	hasShift         bool
	hasNegation      bool
}
