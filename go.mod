module go.nhat.io/aferocopy/v2

go 1.25.0

require (
	github.com/spf13/afero v1.15.0
	github.com/stretchr/testify v1.12.1
	go.nhat.io/aferomock v0.9.1
)

require (
	github.com/stretchr/objx v0.5.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/text v0.41.0 // indirect
)

// v2.0.0 has a bug in permission control that does not apply correct permission to the copied files.
retract v2.0.0
