package concourse

var BuiltinResourceTypes = []ResourceType{
	{
		Name: "git",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/git-resource",
			Tag:        "1.15.0",
		},
	},
	{
		Name: "s3",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/s3-resource",
			Tag:        "1.3.0",
		},
	},
	{
		Name: "docker-image",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/docker-image-resource",
			Tag:        "1.8.0",
		},
	},
	{
		Name: "github-release",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/github-release-resource",
			Tag:        "1.9.0",
		},
	},
	{
		Name: "registry-image",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/registry-image-resource",
			Tag:        "1.9.0",
		},
	},
	{
		Name: "pool",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/pool-resource",
			Tag:        "1.4.0",
		},
	},
	{
		Name: "time",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/time-resource",
			Tag:        "1.7.0",
		},
	},
	{
		Name: "semver",
		Source: &ResourceTypeSource{
			Repository: "docker.io/concourse/semver-resource",
			Tag:        "1.7.0",
		},
	},
}
