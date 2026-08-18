package pkg

import "os"

func osEnvironImpl() []string { return os.Environ() }
