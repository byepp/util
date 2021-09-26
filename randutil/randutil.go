package randutil

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func GenerateUUID(upper bool, withoutMinus bool) (ret string) {
	r, _ := uuid.NewRandom()
	ret = r.String()
	if upper {
		ret = strings.ToUpper(ret)
	}
	if withoutMinus {
		ret = strings.ReplaceAll(ret, "-", "")
	}
	return
}
