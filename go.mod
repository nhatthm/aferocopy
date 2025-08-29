module go.nhat.io/aferocopy/v2

go 1.23.0

require (
	github.com/spf13/afero v1.14.0
	github.com/stretchr/testify v1.11.1
	go.nhat.io/aferomock v0.8.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/text v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// v2.0.0 has a bug in permission control that does not apply correct permission to the copied files.
retract v2.0.0
