module go.nhat.io/aferocopy/v2

go 1.17

require (
	github.com/spf13/afero v1.9.4
	github.com/stretchr/testify v1.8.1
	go.nhat.io/aferomock v0.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// v2.0.0 has a bug in permission control that does not apply correct permission to the copied files.
retract v2.0.0
