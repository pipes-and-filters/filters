test:
	@glide novendor|xargs go test -v
cover:
	@glide novendor|xargs go test -v -covermode=count -coverprofile=cover.profile
