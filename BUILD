load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "gazelle")

go_prefix("github.com/orvice/sox")

# bazel rule definition
gazelle(
  prefix = "github.com/orvice/sox",
  name = "gazelle",
  command = "fix",
)
