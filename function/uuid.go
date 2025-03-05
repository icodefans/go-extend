package function

import (
	"github.com/satori/go.uuid"
)

func UUID() string {
	u1 := uuid.NewV4()
	// fmt.Println("UUID V4:", u1)
	return u1.String()
}
