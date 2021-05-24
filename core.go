package dbb

const (
	DefInitVer string = "0.0"
)

type VerHandler struct {
	FromVer     string
	ToVer       string
	HandlerFunc interface{}
}

type VerHandlers []*VerHandler

func NewVerHandlers() *VerHandlers {
	return &VerHandlers{}
}

func (vhs *VerHandlers) AddVerHandler(fromV, toV string, handlerFunc interface{}) error {
	for _, vh := range *vhs {
		if vh.FromVer == fromV {
			return ErrVerHandlerExists
		}
	}

	*vhs = append(*vhs, &VerHandler{
		FromVer:     fromV,
		ToVer:       toV,
		HandlerFunc: handlerFunc,
	})

	return nil
}

func (vhs *VerHandlers) getVerHandler(fromV string) *VerHandler {
	for _, vh := range *vhs {
		if vh.FromVer == fromV {
			return vh
		}
	}

	return nil
}

//================================================================
//
//================================================================
type Schema struct {
	DBName string
	User   string
	Host   string
}
