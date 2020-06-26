package controller

import (
	"hmxq.top/oilapp-operator/pkg/controller/oilappservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, oilappservice.Add)
}
