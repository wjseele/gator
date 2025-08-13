module github.com/wjseele/gator

go 1.24.5

replace github.com/wjseele/gator/internal/config v0.0.0 => ./internal/config

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)
