package file

type Metadata struct {
	Type Type
}

type File struct {
	Metadata Metadata
	File     []byte
}
