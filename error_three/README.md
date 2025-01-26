## github.com/LydinaPavuh/verbose_errors/error_three



Prints wrapped errors in a readable form.

```go
import (
	"fmt"
	"github.com/LydinaPavuh/verbose_errors/error_three"
)

func main() {
	third := error_three.Error("Third", error_three.Error("s3"), error_three.Error("s3.3"))
	second := error_three.Error("Second", error_three.Error("s2"), error_three.Error("s2.1", third))
	first := error_three.Error("Root err", second)

	fmt.Println(error_three.PrintErrorThree(first))
}

-- 1: Root err
---- 2: Second
------ 3: s2.1
-------- 4: Third
---------- 5: s3.3
---------- 5: s3
------ 3: s2
```
