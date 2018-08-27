# roadrunner
roadrunner compiles your go executable, setting up watches (via fsnotify) for the package and all dependencies.
When a change occurs anywhere it recompiles and restarts the binary.

# install
`go get github.com/tike/roadrunner`

# run
`roadrunner my/import/path -flag forBinary arg for binary`

