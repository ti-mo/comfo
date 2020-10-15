package comfo

//go:generate protoc comfo.proto --go_out=paths=source_relative:. --twirp_out=paths=source_relative:.
//go:generate protoc comfo.proto --python_out=../../python/comfo --twirpy_out=../../python/comfo
