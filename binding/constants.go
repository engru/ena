package binding

// Source is the binding source
type Source string

const (
	// SourceQuery indicates the value from Query
	SourceQuery = "Query"

	// SourceHeader indicates the value from Header
	SourceHeader = "Header"

	// SourceBody indicates the value from Body
	SourceBody = "Body"

	// SourcePath indicates the value from Path
	SourcePath = "Path"

	// SourceAuto indicates the value should populate from sub-field tag.
	// NOTE: this should only been used with Struct
	SourceAuto = "Auto"
)

var (
	emptyStruct     = struct{}{}
	availableSource = map[Source]struct{}{
		SourceQuery:  emptyStruct,
		SourceHeader: emptyStruct,
		SourceBody:   emptyStruct,
		SourcePath:   emptyStruct,
		SourceAuto:   emptyStruct,
	}
)
