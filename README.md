# Generate DAO CRUD, List, DeleteAll functions for models.

## Example of usage
``` go
import (
	"github.com/qwertypomy/crud_generator"
	"github.com/qwertypomy/printers/models"
)

func main() {
    crud_generator.GenerateFiles(
        "Printer",
        models.PrintSize{},
        models.Brand{},
        models.PrintingTechnology{},
        models.FunctionType{},
        models.PrintResolution{},
        models.Printer{},
    )
}
```
