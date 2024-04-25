module URLchess

go 1.16

require (
	github.com/andrewbackes/chess v0.0.0-20171122002438-368c396b5300
	github.com/gopherjs/gopherjs v0.0.0-20210503212227-fb464eba2686
)

replace github.com/andrewbackes/chess => github.com/jezek/chess v1.2.0 // Temporary, until merge requests in andrewbackes/chess are accepted.
