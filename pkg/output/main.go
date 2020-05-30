package output

type Output interface {
	Write (filename string, err error, skipped bool)
	Flush ()
}

