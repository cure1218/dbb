# DB-Builder
```
import (
	"fmt"
	"github.com/cure1218/dbb"
)

func main() {
	vhs := dbb.NewVerHandlers()
	vhs.AddVerHandler("0.0", "1.0", func (db) error {
		...
	})
	vhs.AddVerHandler("1.0", "1.1", func (db) error {
		...
	})

	if dbi, err := dbb.AdminConn("mysql", "localhost", "3306", "admin_user", "password"); err != nil {
		fmt.Println(err)
	} else if err := dbi.Build("dbName", "user", "userHost", "password", vhs); err != nil {
		fmt.Println(err)
	}
}
```
