load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
    ],
    sha256 = "142dd33e38b563605f0d20e89d9ef9eda0fc3cb539a14be1bdb1350de2eda659",
)

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

http_archive(
    name = "rules_proto",
    sha256 = "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208",
    strip_prefix = "rules_proto-97d8af4dc474595af3900dd85cb3a29ad28cc313",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
        "https://github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
    ],
)

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

http_archive(
    name = "rules_python",
    sha256 = "aa96a691d3a8177f3215b14b0edc9641787abaaa30363a080165d06ab65e1161",
    url = "https://github.com/bazelbuild/rules_python/releases/download/0.0.1/rules_python-0.0.1.tar.gz",
)

load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "dc97fccceacd4c6be14e800b2a00693d5e8d07f69ee187babfd04a80a9f8e250",
    strip_prefix = "rules_docker-0.14.1",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.1/rules_docker-v0.14.1.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)
load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

http_archive(
    name = "bazel_toolchains",
    sha256 = "e754d6028845423b2cc7a6c375f9657fe0b0bbb196d76c8de6dd129c3aa74023",
    strip_prefix = "bazel-toolchains-2.2.3",
    urls = [
        "https://github.com/bazelbuild/bazel-toolchains/releases/download/2.2.3/bazel-toolchains-2.2.3.tar.gz",
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-toolchains/releases/download/2.2.3/bazel-toolchains-2.2.3.tar.gz",
    ],
)

load("@bazel_toolchains//rules:rbe_repo.bzl", "rbe_autoconfig")

# Creates a default toolchain config for RBE.
# Use this as is if you are using the rbe_ubuntu16_04 container,
# otherwise refer to RBE docs.
rbe_autoconfig(name = "rbe_default")

GOGOBUILD = """
load("@com_google_protobuf//:protobuf.bzl", "cc_proto_library", "py_proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
filegroup(
    name = "gogo_proto_fg",
    srcs = glob([ "gogoproto/gogo.proto" ]),
    visibility = [ "//visibility:public" ],
)
proto_library(
    name = "gogo_proto",
    srcs = [
        "gogoproto/gogo.proto",
    ],
    deps = [
        "@com_google_protobuf//:descriptor_proto",
    ],
    visibility = ["//visibility:public"],
)
go_proto_library(
    name = "descriptor_go_proto",
    importpath = "github.com/golang/protobuf/protoc-gen-go/descriptor",
    proto = "@com_google_protobuf//:descriptor_proto",
    visibility = ["//visibility:public"],
)
cc_proto_library(
    name = "gogo_proto_cc",
    srcs = [
        "gogoproto/gogo.proto",
    ],
    default_runtime = "@com_google_protobuf//:protobuf",
    protoc = "@com_google_protobuf//:protoc",
    deps = ["@com_google_protobuf//:cc_wkt_protos"],
    visibility = ["//visibility:public"],
)
go_proto_library(
    name = "gogo_proto_go",
    importpath = "gogoproto",
    proto = ":gogo_proto",
    visibility = ["//visibility:public"],
    deps = [
        ":descriptor_go_proto",
    ],
)
py_proto_library(
    name = "gogo_proto_py",
    srcs = [
        "gogoproto/gogo.proto",
    ],
    default_runtime = "@com_google_protobuf//:protobuf_python",
    protoc = "@com_google_protobuf//:protoc",
    visibility = ["//visibility:public"],
    deps = ["@com_google_protobuf//:protobuf_python"],
)
"""

new_git_repository(
    name = "com_github_gogo_protobuf_repo",
    remote = "https://github.com/gogo/protobuf.git",
    branch = "master",
    build_file_content = GOGOBUILD,
)

CONTRAIL_BUILD = """
filegroup(
    name = "contrail_schema",
    srcs = glob([ "schemas/contrail/**/*.yml" ]),
    visibility = [ "//visibility:public" ],
)
"""

new_git_repository(
    name = "contrail_repository",
    remote = "https://github.com/Juniper/contrail.git",
    branch = "master",
    build_file_content = CONTRAIL_BUILD,
)

ASF_BUILD = """
filegroup(
    name = "asf_templates",
    srcs = glob([ "**" ]),
    visibility = [ "//visibility:public" ],
)
filegroup(
    name = "asf_proto",
    srcs = glob([ "pkg/services/baseservices/base.proto" ]),
    visibility = [ "//visibility:public" ],
)
"""

new_git_repository(
    name = "asf_repository",
    remote = "https://github.com/Juniper/asf.git",
    #commit = "ac2649e96024ebed11853a01c83ffd7fc6548919",
    branch = "master",
    build_file_content = ASF_BUILD,
)

go_repository(
    name = "com_github_gogo_protobuf",
    importpath = "github.com/gogo/protobuf/protoc-gen-gogofaster",
    tag = "v1.3.1",
)

go_repository(
    name = "com_github_gogo_protobuf_proto",
    importpath = "github.com/gogo/protobuf",
    tag = "v1.3.1",
)

go_repository(
    name = "co_honnef_go_tools",
    importpath = "honnef.co/go/tools",
    sum = "h1:TEgegKbBqByGUb1Coo1pc2qIdf2xw6v0mYyLSYtyopE=",
    version = "v0.0.1-2019.2.2",
)

go_repository(
    name = "com_github_alecthomas_template",
    importpath = "github.com/alecthomas/template",
    sum = "h1:JYp7IbQjafoB+tBA3gMyHYHrpOtNuDiK/uB5uXxq5wM=",
    version = "v0.0.0-20190718012654-fb15b899a751",
)

go_repository(
    name = "org_golang_x_time",
    build_file_proto_mode = "disable",
    importpath = "golang.org/x/time",
    sum = "h1:/5xXl8Y5W96D+TtHSlonuFqGHIWVuyCkGJLwGh9JJFs=",
    version = "v0.0.0-20191024005414-555d28b269f0",
)

go_repository(
    name = "com_github_alessio_shellescape",
    importpath = "github.com/alessio/shellescape",
    sum = "h1:H/GMMKYPkEIC3DF/JWQz8Pdd+Feifov2EIgGfNpeogI=",
    version = "v0.0.0-20190409004728-b115ca0f9053",
)

go_repository(
    name = "io_k8s_sigs_kind",
    importpath = "sigs.k8s.io/kind",
    sum = "h1:7y7a8EYtGHM+auHmsvzuK5o84SrxPYGidlvfql7j/k4=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_alecthomas_units",
    importpath = "github.com/alecthomas/units",
    sum = "h1:Hs82Z41s6SdL1CELW+XaDYmOH4hkBN4/N9og/AsOv7E=",
    version = "v0.0.0-20190717042225-c3de453c63f4",
)

go_repository(
    name = "com_github_alecthomas_units",
    importpath = "github.com/alecthomas/units",
    sum = "h1:Hs82Z41s6SdL1CELW+XaDYmOH4hkBN4/N9og/AsOv7E=",
    version = "v0.0.0-20190717042225-c3de453c63f4",
)

go_repository(
    name = "com_github_ant31_crd_validation",
    importpath = "github.com/ant31/crd-validation",
    sum = "h1:CDDf61yprxfS7bmBPyhH8pxaobD2VbO3d7laAxJbZos=",
    version = "v0.0.0-20180702145049-30f8a35d0ac2",
)

go_repository(
    name = "com_github_antihax_optional",
    importpath = "github.com/antihax/optional",
    sum = "h1:uZuxRZCz65cG1o6K/xUqImNcYKtmk9ylqaH0itMSvzA=",
    version = "v0.0.0-20180407024304-ca021399b1a6",
)

go_repository(
    name = "com_github_apache_thrift",
    importpath = "github.com/apache/thrift",
    sum = "h1:pODnxUFNcjP9UTLZGTdeh+j16A8lJbRvD3rOtrk/7bs=",
    version = "v0.12.0",
)

go_repository(
    name = "com_github_armon_circbuf",
    importpath = "github.com/armon/circbuf",
    sum = "h1:7Ip0wMmLHLRJdrloDxZfhMm0xrLXZS8+COSu2bXmEQs=",
    version = "v0.0.0-20190214190532-5111143e8da2",
)

go_repository(
    name = "com_github_armon_consul_api",
    importpath = "github.com/armon/consul-api",
    sum = "h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=",
    version = "v0.0.0-20180202201655-eb2c6b5be1b6",
)

go_repository(
    name = "com_github_asaskevich_govalidator",
    importpath = "github.com/asaskevich/govalidator",
    sum = "h1:idn718Q4B6AGu/h5Sxe66HYVdqdGu2l9Iebqhi/AEoA=",
    version = "v0.0.0-20190424111038-f61b66f89f4a",
)

go_repository(
    name = "com_github_auth0_go_jwt_middleware",
    importpath = "github.com/auth0/go-jwt-middleware",
    sum = "h1:irR1cO6eek3n5uquIVaRAsQmZnlsfPuHNz31cXo4eyk=",
    version = "v0.0.0-20170425171159-5493cabe49f7",
)

go_repository(
    name = "com_github_aws_aws_sdk_go",
    importpath = "github.com/aws/aws-sdk-go",
    sum = "h1:MXnqY6SlWySaZAqNnXThOvjRFdiiOuKtC6i7baFdNdU=",
    version = "v1.27.1",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go",
    importpath = "github.com/Azure/azure-sdk-for-go",
    sum = "h1:3D2O4g8AwDwyWkM1HpMFVux/ccQJmGJHXsE004Wsu1Q=",
    version = "v38.0.0+incompatible",
)

go_repository(
    name = "com_github_azure_go_ansiterm",
    importpath = "github.com/Azure/go-ansiterm",
    sum = "h1:w+iIsaOQNcT7OZ575w+acHgRric5iCyQh+xv+KJ4HB8=",
    version = "v0.0.0-20170929234023-d6e3b3328b78",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest",
    importpath = "github.com/Azure/go-autorest/autorest",
    sum = "h1:OZEIaBbMdUE/Js+BQKlpO81XlISgipr6yDJ+PSwsgi4=",
    version = "v0.9.3",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_adal",
    importpath = "github.com/Azure/go-autorest/autorest/adal",
    sum = "h1:pZdL8o72rK+avFWl+p9nE8RWi1JInZrWJYlnpfXJwHk=",
    version = "v0.8.1",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_date",
    importpath = "github.com/Azure/go-autorest/autorest/date",
    sum = "h1:yW+Zlqf26583pE43KhfnhFcdmSWlm5Ew6bxipnr/tbM=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_mocks",
    importpath = "github.com/Azure/go-autorest/autorest/mocks",
    sum = "h1:qJumjCaCudz+OcqE9/XtEPfvtOjOmKaui4EOpFI6zZc=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_to",
    importpath = "github.com/Azure/go-autorest/autorest/to",
    sum = "h1:zebkZaadz7+wIQYgC7GXaz3Wb28yKYfVkkBKwc38VF8=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_validation",
    importpath = "github.com/Azure/go-autorest/autorest/validation",
    sum = "h1:15vMO4y76dehZSq7pAaOLQxC6dZYsSrj2GQpflyM/L4=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_azure_go_autorest_logger",
    importpath = "github.com/Azure/go-autorest/logger",
    sum = "h1:ruG4BSDXONFRrZZJ2GUXDiUyVpayPmb1GnWeHDdaNKY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_azure_go_autorest_tracing",
    importpath = "github.com/Azure/go-autorest/tracing",
    sum = "h1:TRn4WjSnkcSy5AEG3pnbtFSwNtwzjr4VYyQflFE619k=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_bazelbuild_bazel_gazelle",
    importpath = "github.com/bazelbuild/bazel-gazelle",
    sum = "h1:kRymV9q+24Mbeg25fJehw+gvrtVIlwZZAefOSUq4MzU=",
    version = "v0.20.0",
)

go_repository(
    name = "com_github_bazelbuild_buildtools",
    importpath = "github.com/bazelbuild/buildtools",
    sum = "h1:OfyUN/Msd8yqJww6deQ9vayJWw+Jrbe6Qp9giv51QQI=",
    version = "v0.0.0-20190731111112-f720930ceb60",
)

go_repository(
    name = "com_github_beorn7_perks",
    importpath = "github.com/beorn7/perks",
    sum = "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_bifurcation_mint",
    importpath = "github.com/bifurcation/mint",
    sum = "h1:fUjoj2bT6dG8LoEe+uNsKk8J+sLkDbQkJnB6Z1F02Bc=",
    version = "v0.0.0-20180715133206-93c51c6ce115",
)

go_repository(
    name = "com_github_bitly_go_hostpool",
    importpath = "github.com/bitly/go-hostpool",
    sum = "h1:mXoPYz/Ul5HYEDvkta6I8/rnYM5gSdSV2tJ6XbZuEtY=",
    version = "v0.0.0-20171023180738-a3a6125de932",
)

go_repository(
    name = "com_github_bitly_go_simplejson",
    importpath = "github.com/bitly/go-simplejson",
    sum = "h1:6IH+V8/tVMab511d5bn4M7EwGXZf9Hj6i2xSwkNEM+Y=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_blang_semver",
    importpath = "github.com/blang/semver",
    sum = "h1:cQNTCjp13qL8KC3Nbxr/y2Bqb63oX6wdnnjpJbkM4JQ=",
    version = "v3.5.1+incompatible",
)

go_repository(
    name = "com_github_bmizerany_assert",
    importpath = "github.com/bmizerany/assert",
    sum = "h1:DDGfHa7BWjL4YnC6+E63dPcxHo2sUxDIu8g3QgEJdRY=",
    version = "v0.0.0-20160611221934-b7ed37b82869",
)

go_repository(
    name = "com_github_boltdb_bolt",
    importpath = "github.com/boltdb/bolt",
    sum = "h1:JQmyP4ZBrce+ZQu0dY660FMfatumYDLun9hBCUVIkF4=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_brancz_gojsontoyaml",
    importpath = "github.com/brancz/gojsontoyaml",
    sum = "h1:DMb8SuAL9+demT8equqMMzD8C/uxqWmj4cgV7ufrpQo=",
    version = "v0.0.0-20190425155809-e8bd32d46b3d",
)

go_repository(
    name = "com_github_bshuster_repo_logrus_logstash_hook",
    importpath = "github.com/bshuster-repo/logrus-logstash-hook",
    sum = "h1:pgAtgj+A31JBVtEHu2uHuEx0n+2ukqUJnS2vVe5pQNA=",
    version = "v0.4.1",
)

go_repository(
    name = "com_github_bugsnag_bugsnag_go",
    importpath = "github.com/bugsnag/bugsnag-go",
    sum = "h1:tP8hiPv1pGGW3LA6LKy5lW6WG+y9J2xWUdPd3WC452k=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_bugsnag_panicwrap",
    importpath = "github.com/bugsnag/panicwrap",
    sum = "h1:OzrKrRvXis8qEvOkfcxNcYbOd2O7xXS2nnKMEMABFQA=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_burntsushi_toml",
    importpath = "github.com/BurntSushi/toml",
    sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
    version = "v0.3.1",
)

go_repository(
    name = "com_github_burntsushi_xgb",
    importpath = "github.com/BurntSushi/xgb",
    sum = "h1:1BDTz0u9nC3//pOCMdNH+CiXJVYJh5UQNCOBG7jbELc=",
    version = "v0.0.0-20160522181843-27f122750802",
)

go_repository(
    name = "com_github_caddyserver_caddy",
    importpath = "github.com/caddyserver/caddy",
    sum = "h1:i9gRhBgvc5ifchwWtSe7pDpsdS9+Q0Rw9oYQmYUTw1w=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_campoy_embedmd",
    importpath = "github.com/campoy/embedmd",
    sum = "h1:V4kI2qTJJLf4J29RzI/MAt2c3Bl4dQSYPuflzwFH2hY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_cenkalti_backoff",
    importpath = "github.com/cenkalti/backoff",
    sum = "h1:tKJnvO2kl0zmb/jA5UKAt4VoEVw1qxKWjE/Bpp46npY=",
    version = "v2.1.1+incompatible",
)

go_repository(
    name = "com_github_cespare_prettybench",
    importpath = "github.com/cespare/prettybench",
    sum = "h1:p8i+qCbr/dNhS2FoQhRpSS7X5+IlxTa94nRNYXu4fyo=",
    version = "v0.0.0-20150116022406-03b8cfe5406c",
)

go_repository(
    name = "com_github_cespare_xxhash",
    importpath = "github.com/cespare/xxhash",
    sum = "h1:a6HrQnmkObjyL+Gs60czilIUGqrzKutQD6XZog3p+ko=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_cespare_xxhash_v2",
    importpath = "github.com/cespare/xxhash/v2",
    sum = "h1:yTUvW7Vhb89inJ+8irsUqiWjh8iT6sQPZiQzI6ReGkA=",
    version = "v2.1.0",
)

go_repository(
    name = "com_github_chai2010_gettext_go",
    importpath = "github.com/chai2010/gettext-go",
    sum = "h1:7aWHqerlJ41y6FOsEUvknqgXnGmJyJSbjhAWq5pO4F8=",
    version = "v0.0.0-20160711120539-c6fed771bfd5",
)

go_repository(
    name = "com_github_checkpoint_restore_go_criu",
    importpath = "github.com/checkpoint-restore/go-criu",
    sum = "h1:T4nWG1TXIxeor8mAu5bFguPJgSIGhZqv/f0z55KCrJM=",
    version = "v0.0.0-20190109184317-bdb7599cd87b",
)

go_repository(
    name = "com_github_cheekybits_genny",
    importpath = "github.com/cheekybits/genny",
    sum = "h1:a1zrFsLFac2xoM6zG1u72DWJwZG3ayttYLfmLbxVETk=",
    version = "v0.0.0-20170328200008-9127e812e1e9",
)

go_repository(
    name = "com_github_client9_misspell",
    importpath = "github.com/client9/misspell",
    sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
    version = "v0.3.4",
)

go_repository(
    name = "com_github_cloudflare_cfssl",
    importpath = "github.com/cloudflare/cfssl",
    sum = "h1:eOyFuj3h/Vj5e4voOM16NNrHsUR3jhD0duh76LHMj6Y=",
    version = "v0.0.0-20180726162950-56268a613adf",
)

go_repository(
    name = "com_github_clusterhq_flocker_go",
    importpath = "github.com/clusterhq/flocker-go",
    sum = "h1:eIHD9GNM3Hp7kcRW5mvcz7WTR3ETeoYYKwpgA04kaXE=",
    version = "v0.0.0-20160920122132-2b8b7259d313",
)

go_repository(
    name = "com_github_cockroachdb_apd",
    importpath = "github.com/cockroachdb/apd",
    sum = "h1:3LFP3629v+1aKXU5Q37mxmRxX/pIu1nijXydLShEq5I=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_cockroachdb_cockroach_go",
    importpath = "github.com/cockroachdb/cockroach-go",
    sum = "h1:2zRrJWIt/f9c9HhNHAgrRgq0San5gRRUJTBXLkchal0=",
    version = "v0.0.0-20181001143604-e0a95dfd547c",
)

go_repository(
    name = "com_github_codegangsta_negroni",
    importpath = "github.com/codegangsta/negroni",
    sum = "h1:+aYywywx4bnKXWvoWtRfJ91vC59NbEhEY03sZjQhbVY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_container_storage_interface_spec",
    importpath = "github.com/container-storage-interface/spec",
    sum = "h1:qPsTqtR1VUPvMPeK0UnCZMtXaKGyyLPG8gj/wG6VqMs=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_containerd_console",
    importpath = "github.com/containerd/console",
    sum = "h1:GnRy2maqb8vcJhYRN5L+5WyYNKfUG4otiz2zxE182ng=",
    version = "v0.0.0-20170925154832-84eeaae905fa",
)

go_repository(
    name = "com_github_containerd_containerd",
    importpath = "github.com/containerd/containerd",
    sum = "h1:xjvXQWABwS2uiv3TWgQt5Uth60Gu86LTGZXMJkjc7rY=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_containerd_continuity",
    importpath = "github.com/containerd/continuity",
    sum = "h1:4BX8f882bXEDKfWIf0wa8HRvpnBoPszJJXL+TVbBw4M=",
    version = "v0.0.0-20181203112020-004b46473808",
)

go_repository(
    name = "com_github_containerd_typeurl",
    importpath = "github.com/containerd/typeurl",
    sum = "h1:14r0i3IeJj6zkNLigAJiv/TWSR8EY+pxIjv5tFiT+n8=",
    version = "v0.0.0-20190228175220-2a93cfde8c20",
)

go_repository(
    name = "com_github_containernetworking_cni",
    importpath = "github.com/containernetworking/cni",
    sum = "h1:fE3r16wpSEyaqY4Z4oFrLMmIGfBYIKpPrHK31EJ9FzE=",
    version = "v0.7.1",
)

go_repository(
    name = "com_github_coredns_corefile_migration",
    importpath = "github.com/coredns/corefile-migration",
    sum = "h1:kQga1ATFIZdkBtU6c/oJdtASLcCRkDh3fW8vVyVdvUc=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_coreos_bbolt",
    importpath = "github.com/coreos/bbolt",
    sum = "h1:wZwiHHUieZCquLkDL0B8UhzreNWsPHooDAG3q34zk0s=",
    version = "v1.3.2",
)

go_repository(
    name = "com_github_coreos_etcd",
    importpath = "github.com/coreos/etcd",
    sum = "h1:+9RjdC18gMxNQVvSiXvObLu29mOFmkgdsB4cRTlV+EE=",
    version = "v3.3.15+incompatible",
)

go_repository(
    name = "com_github_coreos_go_etcd",
    importpath = "github.com/coreos/go-etcd",
    sum = "h1:bXhRBIXoTm9BYHS3gE0TtQuyNZyeEMux2sDi4oo5YOo=",
    version = "v2.0.0+incompatible",
)

go_repository(
    name = "com_github_coreos_go_oidc",
    importpath = "github.com/coreos/go-oidc",
    sum = "h1:sdJrfw8akMnCuUlaZU3tE/uYXFgfqom8DBE9so9EBsM=",
    version = "v2.1.0+incompatible",
)

go_repository(
    name = "com_github_coreos_go_semver",
    importpath = "github.com/coreos/go-semver",
    sum = "h1:wkHLiw0WNATZnSG7epLsujiMCgPAc9xhjJ4tgnAxmfM=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_coreos_go_systemd",
    importpath = "github.com/coreos/go-systemd",
    sum = "h1:W8b4lQ4tFF21aspRGoBuCNV6V2fFJBF+pm1J6OY8Lys=",
    version = "v0.0.0-20190620071333-e64a0ec8b42a",
)

go_repository(
    name = "com_github_coreos_pkg",
    importpath = "github.com/coreos/pkg",
    sum = "h1:lBNOc5arjvs8E5mO2tbpBpLoyyu8B6e44T7hJy6potg=",
    version = "v0.0.0-20180928190104-399ea9e2e55f",
)

go_repository(
    name = "com_github_coreos_prometheus_operator",
    importpath = "github.com/coreos/prometheus-operator",
    sum = "h1:TF9qaydNeUamLKs0hniaapa4FBz8U8TIlRRtJX987A4=",
    version = "v0.34.0",
)

go_repository(
    name = "com_github_coreos_rkt",
    importpath = "github.com/coreos/rkt",
    sum = "h1:Kkt6sYeEGKxA3Y7SCrY+nHoXkWed6Jr2BBY42GqMymM=",
    version = "v1.30.0",
)

go_repository(
    name = "com_github_cpuguy83_go_md2man",
    importpath = "github.com/cpuguy83/go-md2man",
    sum = "h1:BSKMNlYxDvnunlTymqtgONjNnaRV1sTpcovwwjF22jk=",
    version = "v1.0.10",
)

go_repository(
    name = "com_github_cyphar_filepath_securejoin",
    importpath = "github.com/cyphar/filepath-securejoin",
    sum = "h1:jCwT2GTP+PY5nBz3c/YL5PAIbusElVrPujOBSCj8xRg=",
    version = "v0.2.2",
)

go_repository(
    name = "com_github_cznic_b",
    importpath = "github.com/cznic/b",
    sum = "h1:UHFGPvSxX4C4YBApSPvmUfL8tTvWLj2ryqvT9K4Jcuk=",
    version = "v0.0.0-20180115125044-35e9bbe41f07",
)

go_repository(
    name = "com_github_cznic_fileutil",
    importpath = "github.com/cznic/fileutil",
    sum = "h1:7uSNgsgcarNk4oiN/nNkO0J7KAjlsF5Yv5Gf/tFdHas=",
    version = "v0.0.0-20180108211300-6a051e75936f",
)

go_repository(
    name = "com_github_cznic_golex",
    importpath = "github.com/cznic/golex",
    sum = "h1:CVAqftqbj+exlab+8KJQrE+kNIVlQfJt58j4GxCMF1s=",
    version = "v0.0.0-20170803123110-4ab7c5e190e4",
)

go_repository(
    name = "com_github_cznic_internal",
    importpath = "github.com/cznic/internal",
    sum = "h1:FHpbUtp2K8X53/b4aFNj4my5n+i3x+CQCZWNuHWH/+E=",
    version = "v0.0.0-20180608152220-f44710a21d00",
)

go_repository(
    name = "com_github_cznic_lldb",
    importpath = "github.com/cznic/lldb",
    sum = "h1:AIA+ham6TSJ+XkMe8imQ/g8KPzMUVWAwqUQQdtuMsHs=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_cznic_mathutil",
    importpath = "github.com/cznic/mathutil",
    sum = "h1:XNT/Zf5l++1Pyg08/HV04ppB0gKxAqtZQBRYiYrUuYk=",
    version = "v0.0.0-20180504122225-ca4c9f2c1369",
)

go_repository(
    name = "com_github_cznic_ql",
    importpath = "github.com/cznic/ql",
    sum = "h1:lcKp95ZtdF0XkWhGnVIXGF8dVD2X+ClS08tglKtf+ak=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_cznic_sortutil",
    importpath = "github.com/cznic/sortutil",
    sum = "h1:hxuZop6tSoOi0sxFzoGGYdRqNrPubyaIf9KoBG9tPiE=",
    version = "v0.0.0-20150617083342-4c7342852e65",
)

go_repository(
    name = "com_github_cznic_strutil",
    importpath = "github.com/cznic/strutil",
    sum = "h1:0rkFMAbn5KBKNpJyHQ6Prb95vIKanmAe62KxsrN+sqA=",
    version = "v0.0.0-20171016134553-529a34b1c186",
)

go_repository(
    name = "com_github_cznic_zappy",
    importpath = "github.com/cznic/zappy",
    sum = "h1:YKKpTb2BrXN2GYyGaygIdis1vXbE7SSAG9axGWIMClg=",
    version = "v0.0.0-20160723133515-2533cb5b45cc",
)

go_repository(
    name = "com_github_data_dog_go_sqlmock",
    importpath = "github.com/DATA-DOG/go-sqlmock",
    sum = "h1:CWUqKXe0s8A2z6qCgkP4Kru7wC11YoAnoupUKFDnH08=",
    version = "v1.3.3",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_daviddengcn_go_colortext",
    importpath = "github.com/daviddengcn/go-colortext",
    sum = "h1:uVsMphB1eRx7xB1njzL3fuMdWRN8HtVzoUOItHMwv5c=",
    version = "v0.0.0-20160507010035-511bcaf42ccd",
)

go_repository(
    name = "com_github_deislabs_oras",
    importpath = "github.com/deislabs/oras",
    sum = "h1:RnDoFd3tQYODMiUqxgQ8JxlrlWL0/VMKIKRD01MmNYk=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_denisenkom_go_mssqldb",
    importpath = "github.com/denisenkom/go-mssqldb",
    sum = "h1:tkum0XDgfR0jcVVXuTsYv/erY2NnEDqwRojbxR1rBYA=",
    version = "v0.0.0-20190515213511-eb9f6a1743f3",
)

go_repository(
    name = "com_github_dgrijalva_jwt_go",
    importpath = "github.com/dgrijalva/jwt-go",
    sum = "h1:7qlOGliEKZXTDg6OTjfoBKDXWrumCAMpl/TFQ4/5kLM=",
    version = "v3.2.0+incompatible",
)

go_repository(
    name = "com_github_dgryski_go_sip13",
    importpath = "github.com/dgryski/go-sip13",
    sum = "h1:RMLoZVzv4GliuWafOuPuQDKSm1SJph7uCRnnS61JAn4=",
    version = "v0.0.0-20181026042036-e10d5fee7954",
)

go_repository(
    name = "com_github_dhui_dktest",
    importpath = "github.com/dhui/dktest",
    sum = "h1:kwX5a7EkLcjo7VpsPQSYJcKGbXBXdjI9FGjuUj1jn6I=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_dnaeon_go_vcr",
    importpath = "github.com/dnaeon/go-vcr",
    sum = "h1:r8L/HqC0Hje5AXMu1ooW8oyQyOFv4GxqpL0nRP7SLLY=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_docker_cli",
    importpath = "github.com/docker/cli",
    sum = "h1:2HQmlpI3yI9deH18Q6xiSOIjXD4sLI55Y/gfpa8/558=",
    version = "v0.0.0-20191017083524-a8ff7f821017",
)

go_repository(
    name = "com_github_docker_distribution",
    importpath = "github.com/docker/distribution",
    sum = "h1:a5mlkVzth6W5A4fOsS3D2EO5BUmsJpcB+cRlLU7cSug=",
    version = "v2.7.1+incompatible",
)

go_repository(
    name = "com_github_docker_docker",
    importpath = "github.com/docker/docker",
    replace = "github.com/moby/moby",
    sum = "h1:cvy4lBOYN3gKfKj8Lzz5Q9TfviP+L7koMHY7SvkyTKs=",
    version = "v0.7.3-0.20190826074503-38ab9da00309",
)

go_repository(
    name = "com_github_docker_docker_credential_helpers",
    importpath = "github.com/docker/docker-credential-helpers",
    sum = "h1:zI2p9+1NQYdnG6sMU26EX4aVGlqbInSQxQXLvzJ4RPQ=",
    version = "v0.6.3",
)

go_repository(
    name = "com_github_docker_go_connections",
    importpath = "github.com/docker/go-connections",
    sum = "h1:El9xVISelRB7BuFusrZozjnkIM5YnzCViNKohAFqRJQ=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_docker_go_metrics",
    importpath = "github.com/docker/go-metrics",
    sum = "h1:X0fj836zx99zFu83v/M79DuBn84IL/Syx1SY6Y5ZEMA=",
    version = "v0.0.0-20181218153428-b84716841b82",
)

go_repository(
    name = "com_github_docker_go_units",
    importpath = "github.com/docker/go-units",
    sum = "h1:3uh0PgVws3nIA0Q+MwDC8yjEPf9zjRfZZWXZYDct3Tw=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_docker_libnetwork",
    importpath = "github.com/docker/libnetwork",
    sum = "h1:8rOK787QQFFZJcOLXPiKKidY/ie2OQpblM5gEAaenPs=",
    version = "v0.0.0-20180830151422-a9cd636e3789",
)

go_repository(
    name = "com_github_docker_libtrust",
    importpath = "github.com/docker/libtrust",
    sum = "h1:UhxFibDNY/bfvqU5CAUmr9zpesgbU6SWc8/B4mflAE4=",
    version = "v0.0.0-20160708172513-aabc10ec26b7",
)

go_repository(
    name = "com_github_docker_spdystream",
    importpath = "github.com/docker/spdystream",
    sum = "h1:ZfSZ3P3BedhKGUhzj7BQlPSU4OvT6tfOKe3DVHzOA7s=",
    version = "v0.0.0-20181023171402-6480d4af844c",
)

go_repository(
    name = "com_github_docopt_docopt_go",
    importpath = "github.com/docopt/docopt-go",
    sum = "h1:bWDMxwH3px2JBh6AyO7hdCn/PkvCZXii8TGj7sbtEbQ=",
    version = "v0.0.0-20180111231733-ee0de3bc6815",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    importpath = "github.com/dustin/go-humanize",
    sum = "h1:VSnTsYCnlFHaM2/igO1h6X3HA71jcobQuxemgkq4zYo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_eapache_go_resiliency",
    importpath = "github.com/eapache/go-resiliency",
    sum = "h1:1NtRmCAqadE2FN4ZcN6g90TP3uk8cg9rn9eNK2197aU=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_eapache_go_xerial_snappy",
    importpath = "github.com/eapache/go-xerial-snappy",
    sum = "h1:YEetp8/yCZMuEPMUDHG0CW/brkkEp8mzqk2+ODEitlw=",
    version = "v0.0.0-20180814174437-776d5712da21",
)

go_repository(
    name = "com_github_eapache_queue",
    importpath = "github.com/eapache/queue",
    sum = "h1:YOEu7KNc61ntiQlcEeUIoDTJ2o8mQznoNvUhiigpIqc=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_edsrzf_mmap_go",
    importpath = "github.com/edsrzf/mmap-go",
    sum = "h1:aaQcKT9WumO6JEJcRyTqFVq4XUZiUcKR2/GI31TOcz8=",
    version = "v0.0.0-20170320065105-0bce6a688712",
)

go_repository(
    name = "com_github_elazarl_goproxy",
    importpath = "github.com/elazarl/goproxy",
    sum = "h1:p1yVGRW3nmb85p1Sh1ZJSDm4A4iKLS5QNbvUHMgGu/M=",
    version = "v0.0.0-20170405201442-c4fc26588b6e",
)

go_repository(
    name = "com_github_emicklei_go_restful",
    importpath = "github.com/emicklei/go-restful",
    sum = "h1:CjKsv3uWcCMvySPQYKxO8XX3f9zD4FeZRsW4G0B4ffE=",
    version = "v2.11.1+incompatible",
)

go_repository(
    name = "com_github_euank_go_kmsg_parser",
    importpath = "github.com/euank/go-kmsg-parser",
    sum = "h1:cHD53+PLQuuQyLZeriD1V/esuG4MuU0Pjs5y6iknohY=",
    version = "v2.0.0+incompatible",
)

go_repository(
    name = "com_github_evanphx_json_patch",
    importpath = "github.com/evanphx/json-patch",
    sum = "h1:ouOWdg56aJriqS0huScTkVXPC5IcNrDCXZ6OoTAWu7M=",
    version = "v4.5.0+incompatible",
)

go_repository(
    name = "com_github_exponent_io_jsonpath",
    importpath = "github.com/exponent-io/jsonpath",
    sum = "h1:105gxyaGwCFad8crR9dcMQWvV9Hvulu6hwUh4tWPJnM=",
    version = "v0.0.0-20151013193312-d6023ce2651d",
)

go_repository(
    name = "com_github_fatih_camelcase",
    importpath = "github.com/fatih/camelcase",
    sum = "h1:hxNvNX/xYBp0ovncs8WyWZrOrpBNub/JfaMvbURyft8=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_fatih_color",
    importpath = "github.com/fatih/color",
    sum = "h1:DkWD4oS2D8LGGgTQ6IvwJJXSL5Vp2ffcQg58nFV38Ys=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_fatih_structtag",
    importpath = "github.com/fatih/structtag",
    sum = "h1:6j4mUV/ES2duvnAzKMFkN6/A5mCaNYPD3xfbAkLLOF8=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_flynn_go_shlex",
    importpath = "github.com/flynn/go-shlex",
    sum = "h1:BHsljHzVlRcyQhjrss6TZTdY2VfCqZPbv5k3iBFa2ZQ=",
    version = "v0.0.0-20150515145356-3f9db97f8568",
)

go_repository(
    name = "com_github_fortytw2_leaktest",
    importpath = "github.com/fortytw2/leaktest",
    sum = "h1:u8491cBMTQ8ft8aeV+adlcytMZylmA5nnwwkRZjI8vw=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_fsnotify_fsnotify",
    importpath = "github.com/fsnotify/fsnotify",
    sum = "h1:IXs+QLmnXW2CcXuY+8Mzv/fWEsPGWxqefPtCP5CnV9I=",
    version = "v1.4.7",
)

go_repository(
    name = "com_github_fsouza_fake_gcs_server",
    importpath = "github.com/fsouza/fake-gcs-server",
    sum = "h1:Un0BXUXrRWYSmYyC1Rqm2e2WJfTPyDy/HGMz31emTi8=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_garyburd_redigo",
    importpath = "github.com/garyburd/redigo",
    sum = "h1:0VruCpn7yAIIu7pWVClQC8wxCJEcG3nyzpMSHKi1PQc=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    sum = "h1:Mn26/9ZMNWSw9C9ERFA1PUxfmGpolnw2v0bKOREu5ew=",
    version = "v1.0.1-0.20190212211648-25d852aebe32",
)

go_repository(
    name = "com_github_ghodss",
    importpath = "github.com/ghodss/yaml",
    sum = "h1:Mn26/9ZMNWSw9C9ERFA1PUxfmGpolnw2v0bKOREu5ew=",
    version = "v1.0.1-0.20190212211648-25d852aebe32",
)

go_repository(
    name = "com_github_globalsign_mgo",
    importpath = "github.com/globalsign/mgo",
    sum = "h1:DujepqpGd1hyOd7aW59XpK7Qymp8iy83xq74fLr21is=",
    version = "v0.0.0-20181015135952-eeefdecb41b8",
)

go_repository(
    name = "com_github_go_acme_lego",
    importpath = "github.com/go-acme/lego",
    sum = "h1:5fNN9yRQfv8ymH3DSsxla+4aYeQt2IgfZqHKVnK8f0s=",
    version = "v2.5.0+incompatible",
)

go_repository(
    name = "com_github_go_bindata_go_bindata",
    importpath = "github.com/go-bindata/go-bindata",
    sum = "h1:tR4f0e4VTO7LK6B2YWyAoVEzG9ByG1wrXB4TL9+jiYg=",
    version = "v3.1.1+incompatible",
)

go_repository(
    name = "com_github_go_kit_kit",
    importpath = "github.com/go-kit/kit",
    sum = "h1:wDJmvq38kDhkVxi50ni9ykkdUr1PKgqKOoi01fa0Mdk=",
    version = "v0.9.0",
)

go_repository(
    name = "com_github_go_logfmt_logfmt",
    importpath = "github.com/go-logfmt/logfmt",
    sum = "h1:MP4Eh7ZCb31lleYCFuwm0oe4/YGak+5l1vA2NOE80nA=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_go_logr_logr",
    importpath = "github.com/go-logr/logr",
    sum = "h1:M1Tv3VzNlEHg6uyACnRdtrploV2P7wZqH8BoQMtz0cg=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_go_logr_zapr",
    importpath = "github.com/go-logr/zapr",
    sum = "h1:qXBXPDdNncunGs7XeEpsJt8wCjYBygluzfdLO0G5baE=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_go_openapi_analysis",
    importpath = "github.com/go-openapi/analysis",
    sum = "h1:ophLETFestFZHk3ji7niPEL4d466QjW+0Tdg5VyDq7E=",
    version = "v0.19.2",
)

go_repository(
    name = "com_github_go_openapi_errors",
    importpath = "github.com/go-openapi/errors",
    sum = "h1:a2kIyV3w+OS3S97zxUndRVD46+FhGOUBDFY7nmu4CsY=",
    version = "v0.19.2",
)

go_repository(
    name = "com_github_go_openapi_jsonpointer",
    importpath = "github.com/go-openapi/jsonpointer",
    sum = "h1:gihV7YNZK1iK6Tgwwsxo2rJbD1GTbdm72325Bq8FI3w=",
    version = "v0.19.3",
)

go_repository(
    name = "com_github_go_openapi_jsonreference",
    importpath = "github.com/go-openapi/jsonreference",
    sum = "h1:5cxNfTy0UVC3X8JL5ymxzyoUZmo8iZb+jeTWn7tUa8o=",
    version = "v0.19.3",
)

go_repository(
    name = "com_github_go_openapi_loads",
    importpath = "github.com/go-openapi/loads",
    sum = "h1:rf5ArTHmIJxyV5Oiks+Su0mUens1+AjpkPoWr5xFRcI=",
    version = "v0.19.2",
)

go_repository(
    name = "com_github_go_openapi_runtime",
    importpath = "github.com/go-openapi/runtime",
    sum = "h1:sU6pp4dSV2sGlNKKyHxZzi1m1kG4WnYtWcJ+HYbygjE=",
    version = "v0.19.0",
)

go_repository(
    name = "com_github_go_openapi_spec",
    importpath = "github.com/go-openapi/spec",
    sum = "h1:ixzUSnHTd6hCemgtAJgluaTSGYpLNpJY4mA2DIkdOAo=",
    version = "v0.19.4",
)

go_repository(
    name = "com_github_go_openapi",
    importpath = "github.com/go-openapi/spec",
    sum = "h1:ixzUSnHTd6hCemgtAJgluaTSGYpLNpJY4mA2DIkdOAo=",
    version = "v0.19.4",
)

go_repository(
    name = "com_github_go_openapi_strfmt",
    importpath = "github.com/go-openapi/strfmt",
    sum = "h1:0Dn9qy1G9+UJfRU7TR8bmdGxb4uifB7HNrJjOnV0yPk=",
    version = "v0.19.0",
)

go_repository(
    name = "com_github_go_openapi_swag",
    importpath = "github.com/go-openapi/swag",
    sum = "h1:lTz6Ys4CmqqCQmZPBlbQENR1/GucA2bzYTE12Pw4tFY=",
    version = "v0.19.5",
)

go_repository(
    name = "com_github_go_openapi_validate",
    importpath = "github.com/go-openapi/validate",
    sum = "h1:ky5l57HjyVRrsJfd2+Ro5Z9PjGuKbsmftwyMtk8H7js=",
    version = "v0.19.2",
)

go_repository(
    name = "com_github_go_ozzo_ozzo_validation",
    importpath = "github.com/go-ozzo/ozzo-validation",
    sum = "h1:sUy/in/P6askYr16XJgTKq/0SZhiWsdg4WZGaLsGQkM=",
    version = "v3.5.0+incompatible",
)

go_repository(
    name = "com_github_go_sql_driver_mysql",
    importpath = "github.com/go-sql-driver/mysql",
    sum = "h1:g24URVg0OFbNUTx9qqY1IRZ9D9z3iPyi5zKhQZpNwpA=",
    version = "v1.4.1",
)

go_repository(
    name = "com_github_go_stack_stack",
    importpath = "github.com/go-stack/stack",
    sum = "h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_gobuffalo_envy",
    importpath = "github.com/gobuffalo/envy",
    sum = "h1:GlXgaiBkmrYMHco6t4j7SacKO4XUjvh5pwXh0f4uxXU=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_gobuffalo_flect",
    importpath = "github.com/gobuffalo/flect",
    sum = "h1:xpKq9ap8MbYfhuPCF0dBH854Gp9CxZjr/IocxELFflo=",
    version = "v0.1.5",
)

go_repository(
    name = "com_github_gobuffalo_logger",
    importpath = "github.com/gobuffalo/logger",
    sum = "h1:xw9Ko9EcC5iAFprrjJ6oZco9UpzS5MQ4jAwghsLHdy4=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_gobuffalo_packd",
    importpath = "github.com/gobuffalo/packd",
    sum = "h1:eMwymTkA1uXsqxS0Tpoop3Lc0u3kTfiMBE6nKtQU4g4=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_gobuffalo_packr",
    importpath = "github.com/gobuffalo/packr",
    sum = "h1:hu1fuVR3fXEZR7rXNW3h8rqSML8EVAf6KNm0NKO/wKg=",
    version = "v1.30.1",
)

go_repository(
    name = "com_github_gobuffalo_packr_v2",
    importpath = "github.com/gobuffalo/packr/v2",
    sum = "h1:TFOeY2VoGamPjQLiNDT3mn//ytzk236VMO2j7iHxJR4=",
    version = "v2.5.1",
)

go_repository(
    name = "com_github_gobwas_glob",
    importpath = "github.com/gobwas/glob",
    sum = "h1:A4xDbljILXROh+kObIiy5kIaPYD8e96x1tgBhUI5J+Y=",
    version = "v0.2.3",
)

go_repository(
    name = "com_github_gocql_gocql",
    importpath = "github.com/gocql/gocql",
    sum = "h1:vF83LI8tAakwEwvWZtrIEx7pOySacl2TOxx6eXk4ePo=",
    version = "v0.0.0-20190301043612-f6df8288f9b4",
)

go_repository(
    name = "com_github_godbus_dbus",
    importpath = "github.com/godbus/dbus",
    sum = "h1:WqqLRTsQic3apZUK9qC5sGNfXthmPXzUZ7nQPrNITa4=",
    version = "v4.1.0+incompatible",
)

go_repository(
    name = "com_github_gofrs_flock",
    importpath = "github.com/gofrs/flock",
    sum = "h1:DP+LD/t0njgoPBvT5MJLeliUIVQR03hiKR6vezdwHlc=",
    version = "v0.7.1",
)

go_repository(
    name = "com_github_gofrs_uuid",
    importpath = "github.com/gofrs/uuid",
    sum = "h1:y12jRkkFxsd7GpqdSZ+/KCs/fJbqpEXSGd4+jfEaewE=",
    version = "v3.2.0+incompatible",
)

go_repository(
    name = "com_github_gogo_protobuf",
    importpath = "github.com/gogo/protobuf",
    sum = "h1:DqDEcV5aeaTmdFBePNpYsp3FlcVH/2ISVVM9Qf8PSls=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_golang_glog",
    importpath = "github.com/golang/glog",
    sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
    version = "v0.0.0-20160126235308-23def4e6c14b",
)

go_repository(
    name = "com_github_golang_groupcache",
    importpath = "github.com/golang/groupcache",
    sum = "h1:uHTyIjqVhYRhLbJ8nIiOJHkEZZ+5YoOsAbD3sk82NiE=",
    version = "v0.0.0-20191027212112-611e8accdfc9",
)

go_repository(
    name = "com_github_golang_lint",
    importpath = "github.com/golang/lint",
    sum = "h1:2hRPrmiwPrp3fQX967rNJIhQPtiGXdlQWAxKbKw3VHA=",
    version = "v0.0.0-20180702182130-06c8688daad7",
)

go_repository(
    name = "com_github_golang_migrate_migrate_v4",
    importpath = "github.com/golang-migrate/migrate/v4",
    sum = "h1:LDDOHo/q1W5UDj6PbkxdCv7lv9yunyZHXvxuwDkGo3k=",
    version = "v4.6.2",
)

go_repository(
    name = "com_github_golang_mock",
    importpath = "github.com/golang/mock",
    sum = "h1:qGJ6qTW+x6xX/my+8YUVl4WNpX9B7+/l2tRsHGZ7f2s=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    sum = "h1:gyjaxf+svBWX08ZjK86iN9geUJF0H6gp2IRKX6Nf6/I=",
    version = "v1.3.3",
)

go_repository(
    name = "com_github_golang_snappy",
    importpath = "github.com/golang/snappy",
    sum = "h1:Qgr9rKW7uDUkrbSmQeiDsGa8SjGyCOGtuasMWwvp2P4=",
    version = "v0.0.1",
)

go_repository(
    name = "com_github_golangplus_bytes",
    importpath = "github.com/golangplus/bytes",
    sum = "h1:7xqw01UYS+KCI25bMrPxwNYkSns2Db1ziQPpVq99FpE=",
    version = "v0.0.0-20160111154220-45c989fe5450",
)

go_repository(
    name = "com_github_golangplus_fmt",
    importpath = "github.com/golangplus/fmt",
    sum = "h1:f5gsjBiF9tRRVomCvrkGMMWI8W1f2OBFar2c5oakAP0=",
    version = "v0.0.0-20150411045040-2a5d6d7d2995",
)

go_repository(
    name = "com_github_golangplus_testing",
    importpath = "github.com/golangplus/testing",
    sum = "h1:KhcknUwkWHKZPbFy2P7jH5LKJ3La+0ZeknkkmrSgqb0=",
    version = "v0.0.0-20180327235837-af21d9c3145e",
)

go_repository(
    name = "com_github_google_btree",
    importpath = "github.com/google/btree",
    sum = "h1:0udJVsspx3VBr5FwtLhQQtuAsVc79tTq0ocGIPAU6qo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_cadvisor",
    importpath = "github.com/google/cadvisor",
    sum = "h1:No7G6U/TasplR9uNqyc5Jj0Bet5VSYsK5xLygOf4pUw=",
    version = "v0.34.0",
)

go_repository(
    name = "com_github_google_certificate_transparency_go",
    importpath = "github.com/google/certificate-transparency-go",
    sum = "h1:Yf1aXowfZ2nuboBsg7iYGLmwsOARdV86pfH3g95wXmE=",
    version = "v1.0.21",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    sum = "h1:xsAVV57WRhGj6kEIi8ReJzQlHHqcBYCElAvkovg3B/4=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_google_go_github",
    importpath = "github.com/google/go-github",
    sum = "h1:N0LgJ1j65A7kfXrZnUDaYCs/Sf4rEjNlfyDHW9dolSY=",
    version = "v17.0.0+incompatible",
)

go_repository(
    name = "com_github_google_go_querystring",
    importpath = "github.com/google/go-querystring",
    sum = "h1:Xkwi/a1rcvNg1PPYe5vI8GbeBY/jrVuDX5ASuANWTrk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_gofuzz",
    importpath = "github.com/google/gofuzz",
    sum = "h1:A8PeW59pxE9IoFRqBp37U+mSNaQoZ46F1f0f863XSXw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_martian",
    importpath = "github.com/google/martian",
    sum = "h1:/CP5g8u/VJHijgedC/Legn3BAbAaWPgecwXBIDzw5no=",
    version = "v2.1.0+incompatible",
)

go_repository(
    name = "com_github_google_pprof",
    importpath = "github.com/google/pprof",
    sum = "h1:Jnx61latede7zDD3DiiP4gmNz33uK0U5HDUaF0a/HVQ=",
    version = "v0.0.0-20190515194954-54271f7e092f",
)

go_repository(
    name = "com_github_google_renameio",
    importpath = "github.com/google/renameio",
    sum = "h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_google_shlex",
    importpath = "github.com/google/shlex",
    sum = "h1:El6M4kTTCOh6aBiKaUGG7oYTSPP8MxqL4YI3kZKwcP4=",
    version = "v0.0.0-20191202100458-e7afc7fbc510",
)

go_repository(
    name = "com_github_google_uuid",
    importpath = "github.com/google/uuid",
    sum = "h1:Gkbcsh/GbpXz7lPftLA3P6TYMwjCLYm83jiFQZF/3gY=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_googleapis_gax_go_v2",
    importpath = "github.com/googleapis/gax-go/v2",
    sum = "h1:sjZBwGj9Jlw33ImPtvFviGYvseOtDM7hkSKB7+Tv3SM=",
    version = "v2.0.5",
)

go_repository(
    name = "com_github_googleapis_gnostic",
    importpath = "github.com/googleapis/gnostic",
    sum = "h1:WeAefnSUHlBb0iJKwxFDZdbfGwkd7xRNuV+IpXMJhYk=",
    version = "v0.3.1",
)

go_repository(
    name = "com_github_googlecloudplatform_k8s_cloud_provider",
    importpath = "github.com/GoogleCloudPlatform/k8s-cloud-provider",
    sum = "h1:N7lSsF+R7wSulUADi36SInSQA3RvfO/XclHQfedr0qk=",
    version = "v0.0.0-20190822182118-27a4ced34534",
)

go_repository(
    name = "com_github_gophercloud_gophercloud",
    importpath = "github.com/gophercloud/gophercloud",
    sum = "h1:lD2Bce2xBAMNNcFZ0dObTpXkGLlVIb33RPVUNVpw6ic=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_gopherjs_gopherjs",
    importpath = "github.com/gopherjs/gopherjs",
    sum = "h1:EGx4pi6eqNxGaHF6qqu48+N2wcFQ5qg5FXgOdqsJ5d8=",
    version = "v0.0.0-20181017120253-0766667cb4d1",
)

go_repository(
    name = "com_github_gorilla_context",
    importpath = "github.com/gorilla/context",
    sum = "h1:AWwleXJkX/nhcU9bZSnZoi3h/qGYqQAGhq6zZe/aQW8=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_gorilla_handlers",
    importpath = "github.com/gorilla/handlers",
    sum = "h1:XulKRWSQK5uChr4pEgSE4Tc/OcmnU9GJuSwdog/tZsA=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:gnP5JzjVOuiZD07fKKToCAOjS0yOpj/qPETTXCCS6hw=",
    version = "v1.7.3",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    sum = "h1:q7AeDBpnBk8AogcD4DSag/Ukw/KV+YhzLj2bP5HvKCM=",
    version = "v1.4.1",
)

go_repository(
    name = "com_github_gosuri_uitable",
    importpath = "github.com/gosuri/uitable",
    sum = "h1:M9sMNgSZPyAu1FJZJLpJ16ofL8q5ko2EDUkICsynvlY=",
    version = "v0.0.1",
)

go_repository(
    name = "com_github_gregjones_httpcache",
    importpath = "github.com/gregjones/httpcache",
    sum = "h1:prg2TTpTOcJF1jRWL2zSU1FQNgB0STAFNux8GK82y8k=",
    version = "v0.0.0-20190203031600-7a902570cb17",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_middleware",
    importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
    sum = "h1:THDBEeQ9xZ8JEaCLyLQqXMMdRqNr0QAUJTIkQAUtFjg=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_prometheus",
    importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
    sum = "h1:Ovs26xHkKqVztRpIrF/92BcuyuQ/YW4NSIpoGtfXNho=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
    sum = "h1:5xLhQjsk4zqPf9EHCrja2qFZMx+yBqkO3XgJ14bNnU0=",
    version = "v1.9.4",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_health_probe",
    importpath = "github.com/grpc-ecosystem/grpc-health-probe",
    sum = "h1:UxmGBzaBcWDQuQh9E1iT1dWKQFbizZ+SpTd1EL4MSqs=",
    version = "v0.2.1-0.20181220223928-2bf0a5b182db",
)

go_repository(
    name = "com_github_hailocab_go_hostpool",
    importpath = "github.com/hailocab/go-hostpool",
    sum = "h1:5upAirOpQc1Q53c0bnx2ufif5kANL7bfZWcc6VJWJd8=",
    version = "v0.0.0-20160125115350-e80d13ce29ed",
)

go_repository(
    name = "com_github_hashicorp_errwrap",
    importpath = "github.com/hashicorp/errwrap",
    sum = "h1:hLrqtEDnRye3+sgx6z4qVLNuviH3MR5aQ0ykNJa/UYA=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_multierror",
    importpath = "github.com/hashicorp/go-multierror",
    sum = "h1:iVjPR7a6H0tWELX5NxNe7bYopibicUzc7uPribsnS6o=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_syslog",
    importpath = "github.com/hashicorp/go-syslog",
    sum = "h1:KaodqZuhUoZereWVIYmpUgZysurB1kBLX2j0MwMrUAE=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_version",
    importpath = "github.com/hashicorp/go-version",
    sum = "h1:3vNe/fWF5CBgRIguda1meWhsZHy3m8gCJ5wx+dIzX/E=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    importpath = "github.com/hashicorp/golang-lru",
    sum = "h1:YPkqC67at8FYaadspW/6uE0COsBxS2656RLEr8Bppgk=",
    version = "v0.5.3",
)

go_repository(
    name = "com_github_hashicorp_hcl",
    importpath = "github.com/hashicorp/hcl",
    sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_heketi_heketi",
    importpath = "github.com/heketi/heketi",
    sum = "h1:B2ACAbYsCHkJXKozYVV7p2j+eEy/zNlLsicihMWCk30=",
    version = "v9.0.0+incompatible",
)

go_repository(
    name = "com_github_heketi_rest",
    importpath = "github.com/heketi/rest",
    sum = "h1:nGZBOxRgSMbqjm2/FYDtO6BU4a+hfR7Om9VGQ9tbbSc=",
    version = "v0.0.0-20180404230133-aa6a65207413",
)

go_repository(
    name = "com_github_heketi_tests",
    importpath = "github.com/heketi/tests",
    sum = "h1:oJ/NLadJn5HoxvonA6VxG31lg0d6XOURNA09BTtM4fY=",
    version = "v0.0.0-20151005000721-f3775cbcefd6",
)

go_repository(
    name = "com_github_heketi_utils",
    importpath = "github.com/heketi/utils",
    sum = "h1:dk3GEa55HcRVIyCeNQmwwwH3kIXnqJPNseKOkDD+7uQ=",
    version = "v0.0.0-20170317161834-435bc5bdfa64",
)

go_repository(
    name = "com_github_helm_helm_2to3",
    importpath = "github.com/helm/helm-2to3",
    sum = "h1:1wgpvw93CHlO+mC6AMrfc5xboYyk1MkoI3ku7eI5uvE=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_hpcloud_tail",
    importpath = "github.com/hpcloud/tail",
    sum = "h1:nfCOvKYfkgYP8hkirhJocXT2+zOD8yUNjXaWfTlyFKI=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_huandu_xstrings",
    importpath = "github.com/huandu/xstrings",
    sum = "h1:yPeWdRnmynF7p+lLYz0H2tthW9lqhMJrQV/U7yy4wX0=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_iancoleman_strcase",
    importpath = "github.com/iancoleman/strcase",
    sum = "h1:ECW73yc9MY7935nNYXUkK7Dz17YuSUI9yqRqYS8aBww=",
    version = "v0.0.0-20190422225806-e506e3ef7365",
)

go_repository(
    name = "com_github_imdario_mergo",
    importpath = "github.com/imdario/mergo",
    sum = "h1:CGgOkSJeqMRmt0D9XLWExdT4m4F1vd3FV3VPt+0VxkQ=",
    version = "v0.3.8",
)

go_repository(
    name = "com_github_improbable_eng_thanos",
    importpath = "github.com/improbable-eng/thanos",
    sum = "h1:iZfU7exq+RD5Lnb8n3Eh9MNYoRLeyeGO/85AvEkLg+8=",
    version = "v0.3.2",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    importpath = "github.com/inconshreveable/mousetrap",
    sum = "h1:Z8tu5sraLXCXIcARxBp/8cbvlwVa7Z1NHg9XEKhtSvM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jackc_fake",
    importpath = "github.com/jackc/fake",
    sum = "h1:vr3AYkKovP8uR8AvSGGUK1IDqRa5lAAvEkZG1LKaCRc=",
    version = "v0.0.0-20150926172116-812a484cc733",
)

go_repository(
    name = "com_github_jackc_pgx",
    importpath = "github.com/jackc/pgx",
    sum = "h1:0Vihzu20St42/UDsvZGdNE6jak7oi/UOeMzwMPHkgFY=",
    version = "v3.2.0+incompatible",
)

go_repository(
    name = "com_github_jeffashton_win_pdh",
    importpath = "github.com/JeffAshton/win_pdh",
    sum = "h1:UKkYhof1njT1/xq4SEg5z+VpTgjmNeHwPGRQl7takDI=",
    version = "v0.0.0-20161109143554-76bb4ee9f0ab",
)

go_repository(
    name = "com_github_jimstudt_http_authentication",
    importpath = "github.com/jimstudt/http-authentication",
    sum = "h1:BcF8coBl0QFVhe8vAMMlD+CV8EISiu9MGKLoj6ZEyJA=",
    version = "v0.0.0-20140401203705-3eca13d6893a",
)

go_repository(
    name = "com_github_jmespath_go_jmespath",
    importpath = "github.com/jmespath/go-jmespath",
    sum = "h1:pmfjZENx5imkbgOkpRUYLnmbU7UEFbjtDA2hxJ1ichM=",
    version = "v0.0.0-20180206201540-c2b33e8439af",
)

go_repository(
    name = "com_github_jmoiron_sqlx",
    importpath = "github.com/jmoiron/sqlx",
    sum = "h1:41Ip0zITnmWNR/vHV+S4m+VoUivnWY5E4OJfLZjCJMA=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_joefitzgerald_rainbow_reporter",
    importpath = "github.com/joefitzgerald/rainbow-reporter",
    sum = "h1:AuMG652zjdzI0YCCnXAqATtRBpGXMcAnrajcaTrSeuo=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_joho_godotenv",
    importpath = "github.com/joho/godotenv",
    sum = "h1:Zjp+RcGpHhGlrMbJzXTrZZPrWj+1vfm90La1wgB6Bhc=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_jonboulle_clockwork",
    importpath = "github.com/jonboulle/clockwork",
    sum = "h1:VKV+ZcuP6l3yW9doeqz6ziZGgcynBVQO+obU0+0hcPo=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_json_iterator_go",
    importpath = "github.com/json-iterator/go",
    sum = "h1:KfgG9LzI+pYjr4xvmz/5H4FXjokeP+rlHLhv3iH62Fo=",
    version = "v1.1.7",
)

go_repository(
    name = "com_github_jsonnet_bundler_jsonnet_bundler",
    importpath = "github.com/jsonnet-bundler/jsonnet-bundler",
    sum = "h1:T/HtHFr+mYCRULrH1x/RnoB0prIs0rMkolJhFMXNC9A=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_jstemmer_go_junit_report",
    importpath = "github.com/jstemmer/go-junit-report",
    sum = "h1:rBMNdlhTLzJjJSDIjNEXX1Pz3Hmwmz91v+zycvx9PJc=",
    version = "v0.0.0-20190106144839-af01ea7f8024",
)

go_repository(
    name = "com_github_jtolds_gls",
    importpath = "github.com/jtolds/gls",
    sum = "h1:xdiiI2gbIgH/gLH7ADydsJ1uDOEzR8yvV7C0MuV77Wo=",
    version = "v4.20.0+incompatible",
)

go_repository(
    name = "com_github_julienschmidt_httprouter",
    importpath = "github.com/julienschmidt/httprouter",
    sum = "h1:TDTW5Yz1mjftljbcKqRcrYhd4XeOoI98t+9HbQbYf7g=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_juniper_contrail_go_api",
    importpath = "github.com/Juniper/contrail-go-api",
    sum = "h1:GswUxqPbmwhx7zQj+8msn0jKaVrVRPnWsPuNizqpU1o=",
    version = "v1.1.1-0.20200414151206-7bb4264f1da4",
)

go_repository(
    name = "com_github_kardianos_osext",
    importpath = "github.com/kardianos/osext",
    sum = "h1:iQTw/8FWTuc7uiaSepXwyf3o52HaUYcV+Tu66S3F5GA=",
    version = "v0.0.0-20190222173326-2bc1f35cddc0",
)

go_repository(
    name = "com_github_karrick_godirwalk",
    importpath = "github.com/karrick/godirwalk",
    sum = "h1:BqUm+LuJcXjGv1d2mj3gBiQyrQ57a0rYoAmhvJQ7RDU=",
    version = "v1.10.12",
)

go_repository(
    name = "com_github_kisielk_errcheck",
    importpath = "github.com/kisielk/errcheck",
    sum = "h1:reN85Pxc5larApoH1keMBiu2GWtPqXQ1nc9gx+jOU+E=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_kisielk_gotool",
    importpath = "github.com/kisielk/gotool",
    sum = "h1:AV2c/EiW3KqPNT9ZKl07ehoAGi4C5/01Cfbblndcapg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_klauspost_cpuid",
    importpath = "github.com/klauspost/cpuid",
    sum = "h1:NMpwD2G9JSFOE1/TJjGSo5zG7Yb2bTe7eq1jH+irmeE=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_konsorten_go_windows_terminal_sequences",
    importpath = "github.com/konsorten/go-windows-terminal-sequences",
    sum = "h1:DB17ag19krx9CFsz4o3enTrPXyIXCl+2iCXH/aMAp9s=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_kr_logfmt",
    importpath = "github.com/kr/logfmt",
    sum = "h1:T+h1c/A9Gawja4Y9mFVWj2vyii2bbUNDw3kt9VxK2EY=",
    version = "v0.0.0-20140226030751-b84e30acd515",
)

go_repository(
    name = "com_github_kr_pretty",
    importpath = "github.com/kr/pretty",
    sum = "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_kr_pty",
    importpath = "github.com/kr/pty",
    sum = "h1:hyz3dwM5QLc1Rfoz4FuWJQG5BN7tc6K1MndAUnGpQr4=",
    version = "v1.1.5",
)

go_repository(
    name = "com_github_kr_text",
    importpath = "github.com/kr/text",
    sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_kshvakov_clickhouse",
    importpath = "github.com/kshvakov/clickhouse",
    sum = "h1:PDTYk9VYgbjPAWry3AoDREeMgOVUFij6bh6IjlloHL0=",
    version = "v1.3.5",
)

go_repository(
    name = "com_github_kylelemons_godebug",
    importpath = "github.com/kylelemons/godebug",
    sum = "h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_lib_pq",
    importpath = "github.com/lib/pq",
    sum = "h1:LXpIM/LZ5xGFhOpXAQUIMM1HdyqzVYM13zNdjCEEcA0=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_libopenstorage_openstorage",
    importpath = "github.com/libopenstorage/openstorage",
    sum = "h1:GLPam7/0mpdP8ZZtKjbfcXJBTIA/T1O6CBErVEFEyIM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_liggitt_tabwriter",
    importpath = "github.com/liggitt/tabwriter",
    sum = "h1:9TO3cAIGXtEhnIaL+V+BEER86oLrvS+kWobKpbJuye0=",
    version = "v0.0.0-20181228230101-89fcab3d43de",
)

go_repository(
    name = "com_github_lithammer_dedent",
    importpath = "github.com/lithammer/dedent",
    sum = "h1:VNzHMVCBNG1j0fh3OrsFRkVUwStdDArbgBWoPAffktY=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_lpabon_godbc",
    importpath = "github.com/lpabon/godbc",
    sum = "h1:ilqjArN1UOENJJdM34I2YHKmF/B0gGq4VLoSGy9iAao=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_lucas_clemente_aes12",
    importpath = "github.com/lucas-clemente/aes12",
    sum = "h1:sSeNEkJrs+0F9TUau0CgWTTNEwF23HST3Eq0A+QIx+A=",
    version = "v0.0.0-20171027163421-cd47fb39b79f",
)

go_repository(
    name = "com_github_lucas_clemente_quic_clients",
    importpath = "github.com/lucas-clemente/quic-clients",
    sum = "h1:/P9n0nICT/GnQJkZovtBqridjxU0ao34m7DpMts79qY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_lucas_clemente_quic_go",
    importpath = "github.com/lucas-clemente/quic-go",
    sum = "h1:iQtTSZVbd44k94Lu0U16lLBIG3lrnjDvQongjPd4B/s=",
    version = "v0.10.2",
)

go_repository(
    name = "com_github_lucas_clemente_quic_go_certificates",
    importpath = "github.com/lucas-clemente/quic-go-certificates",
    sum = "h1:zqEC1GJZFbGZA0tRyNZqRjep92K5fujFtFsu5ZW7Aug=",
    version = "v0.0.0-20160823095156-d2f86524cced",
)

go_repository(
    name = "com_github_magiconair_properties",
    importpath = "github.com/magiconair/properties",
    sum = "h1:ZC2Vc7/ZFkGmsVC9KvOjumD+G5lXy2RtTKyzRKO2BQ4=",
    version = "v1.8.1",
)

go_repository(
    name = "com_github_mailru_easyjson",
    importpath = "github.com/mailru/easyjson",
    sum = "h1:aizVhC/NAAcKWb+5QsU1iNOZb4Yws5UO2I+aIprQITM=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_makenowjust_heredoc",
    importpath = "github.com/MakeNowJust/heredoc",
    sum = "h1:eb0Pzkt15Bm7f2FFYv7sjY7NPFi3cPkS3tv1CcrFBWA=",
    version = "v0.0.0-20171113091838-e9091a26100e",
)

go_repository(
    name = "com_github_maorfr_helm_plugin_utils",
    importpath = "github.com/maorfr/helm-plugin-utils",
    sum = "h1:8HI/lWAXE2BrL85p5qi7Snl7UPV5XMKQ1HynkGhNHBg=",
    version = "v0.0.0-20181205064038-588190cb5e3b",
)

go_repository(
    name = "com_github_markbates_inflect",
    importpath = "github.com/markbates/inflect",
    sum = "h1:5fh1gzTFhfae06u3hzHYO9xe3l3v3nW5Pwt3naLTP5g=",
    version = "v1.0.4",
)

go_repository(
    name = "com_github_marten_seemann_qtls",
    importpath = "github.com/marten-seemann/qtls",
    sum = "h1:0yWJ43C62LsZt08vuQJDK1uC1czUc3FJeCLPoNAI4vA=",
    version = "v0.2.3",
)

go_repository(
    name = "com_github_martinlindhe_base36",
    importpath = "github.com/martinlindhe/base36",
    sum = "h1:eYsumTah144C0A8P1T/AVSUk5ZoLnhfYFM3OGQxB52A=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_masterminds_goutils",
    importpath = "github.com/Masterminds/goutils",
    sum = "h1:zukEsf/1JZwCMgHiK3GZftabmxiCw4apj3a28RPBiVg=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_masterminds_semver",
    importpath = "github.com/Masterminds/semver",
    sum = "h1:H65muMkzWKEuNDnfl9d70GUjFniHKHRbFPGBuZ3QEww=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_masterminds_semver_v3",
    importpath = "github.com/Masterminds/semver/v3",
    sum = "h1:2kKm5lb7dKVrt5TYUiAavE6oFc1cFT0057UVGT+JqLk=",
    version = "v3.0.1",
)

go_repository(
    name = "com_github_masterminds_sprig_v3",
    importpath = "github.com/Masterminds/sprig/v3",
    sum = "h1:KSQz7Nb08/3VU9E4ns29dDxcczhOD1q7O1UfM4G3t3g=",
    version = "v3.0.0",
)

go_repository(
    name = "com_github_masterminds_vcs",
    importpath = "github.com/Masterminds/vcs",
    sum = "h1:USF5TvZGYgIpcbNAEMLfFhHqP08tFZVlUVrmTSpqnyA=",
    version = "v1.13.0",
)

go_repository(
    name = "com_github_mattbaird_jsonpatch",
    importpath = "github.com/mattbaird/jsonpatch",
    sum = "h1:+J2gw7Bw77w/fbK7wnNJJDKmw1IbWft2Ul5BzrG1Qm8=",
    version = "v0.0.0-20171005235357-81af80346b1a",
)

go_repository(
    name = "com_github_mattn_go_colorable",
    importpath = "github.com/mattn/go-colorable",
    sum = "h1:6Su7aK7lXmJ/U79bYtBjLNaha4Fs1Rg9plHpcH+vvnE=",
    version = "v0.1.6",
)

go_repository(
    name = "com_github_mattn_go_isatty",
    importpath = "github.com/mattn/go-isatty",
    sum = "h1:wuysRhFDzyxgEmMf5xjvJ2M9dZoWAXNNr5LSBS7uHXY=",
    version = "v0.0.12",
)

go_repository(
    name = "com_github_mattn_go_runewidth",
    importpath = "github.com/mattn/go-runewidth",
    sum = "h1:2BvfKmzob6Bmd4YsL0zygOqfdFnK7GR4QL06Do4/p7Y=",
    version = "v0.0.4",
)

go_repository(
    name = "com_github_mattn_go_shellwords",
    importpath = "github.com/mattn/go-shellwords",
    sum = "h1:9Jok5pILi5S1MnDirGVTufYGtksUs/V2BWUP3ZkeUUI=",
    version = "v1.0.6",
)

go_repository(
    name = "com_github_mattn_go_sqlite3",
    importpath = "github.com/mattn/go-sqlite3",
    sum = "h1:LDdKkqtYlom37fkvqs8rMPFKAMe8+SgjbwZ6ex1/A/Q=",
    version = "v1.11.0",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
    sum = "h1:4hp9jkHxhMHkqkrB3Ix0jegS5sx/RkqARlsWZ6pIwiU=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_maxbrunsfeld_counterfeiter_v6",
    importpath = "github.com/maxbrunsfeld/counterfeiter/v6",
    sum = "h1:g+4J5sZg6osfvEfkRZxJ1em0VT95/UOZgi/l7zi1/oE=",
    version = "v6.2.2",
)

go_repository(
    name = "com_github_mesos_mesos_go",
    importpath = "github.com/mesos/mesos-go",
    sum = "h1:w8V5sOEnxzHZ2kAOy273v/HgbolyI6XI+qe5jx5u+Y0=",
    version = "v0.0.9",
)

go_repository(
    name = "com_github_mholt_certmagic",
    importpath = "github.com/mholt/certmagic",
    sum = "h1:xKE9kZ5C8gelJC3+BNM6LJs1x21rivK7yxfTZMAuY2s=",
    version = "v0.6.2-0.20190624175158-6a42ef9fe8c2",
)

go_repository(
    name = "com_github_microsoft_go_winio",
    importpath = "github.com/Microsoft/go-winio",
    sum = "h1:+hMXMk01us9KgxGb7ftKQt2Xpf5hH/yky+TDA+qxleU=",
    version = "v0.4.14",
)

go_repository(
    name = "com_github_microsoft_hcsshim",
    importpath = "github.com/Microsoft/hcsshim",
    sum = "h1:ZfF0+zZeYdzMIVMZHKtDKJvLHj76XCuVae/jNkjj0IA=",
    version = "v0.8.6",
)

go_repository(
    name = "com_github_miekg_dns",
    importpath = "github.com/miekg/dns",
    sum = "h1:rCMZsU2ScVSYcAsOXgmC6+AKOK+6pmQTOcw03nfwYV0=",
    version = "v1.1.4",
)

go_repository(
    name = "com_github_mindprince_gonvml",
    importpath = "github.com/mindprince/gonvml",
    sum = "h1:v3dy+FJr7gS7nLgYG7YjX/pmUWuFdudcpnoRNHt2heo=",
    version = "v0.0.0-20171110221305-fee913ce8fb2",
)

go_repository(
    name = "com_github_mistifyio_go_zfs",
    importpath = "github.com/mistifyio/go-zfs",
    sum = "h1:gAMO1HM9xBRONLHHYnu5iFsOJUiJdNZo6oqSENd4eW8=",
    version = "v2.1.1+incompatible",
)

go_repository(
    name = "com_github_mitchellh_copystructure",
    importpath = "github.com/mitchellh/copystructure",
    sum = "h1:Laisrj+bAB6b/yJwB5Bt3ITZhGJdqmxquMKeZ+mmkFQ=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    importpath = "github.com/mitchellh/go-homedir",
    sum = "h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_mitchellh_go_wordwrap",
    importpath = "github.com/mitchellh/go-wordwrap",
    sum = "h1:6GlHJ/LTGMrIJbwgdqdl2eEH8o+Exx/0m8ir9Gns0u4=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_hashstructure",
    importpath = "github.com/mitchellh/hashstructure",
    sum = "h1:ZkRJX1CyOoTkar7p/mLS5TZU4nJ1Rn/F8u9dGS02Q3Y=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    importpath = "github.com/mitchellh/mapstructure",
    sum = "h1:fmNYVwqnSfB9mZU6OS2O6GsXM+wcskZDuKQzvN1EDeE=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_mitchellh_reflectwalk",
    importpath = "github.com/mitchellh/reflectwalk",
    sum = "h1:9D+8oIskB4VJBN5SFlmc27fSlIBZaov1Wpk/IfikLNY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_modern_go_concurrent",
    importpath = "github.com/modern-go/concurrent",
    sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
    version = "v0.0.0-20180306012644-bacd9c7ef1dd",
)

go_repository(
    name = "com_github_modern_go_reflect2",
    importpath = "github.com/modern-go/reflect2",
    sum = "h1:9f412s+6RmYXLWZSEzVVgPGK7C2PphHj5RJrvfx9AWI=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_mohae_deepcopy",
    importpath = "github.com/mohae/deepcopy",
    sum = "h1:e+l77LJOEqXTIQihQJVkA6ZxPOUmfPM5e4H7rcpgtSk=",
    version = "v0.0.0-20170603005431-491d3605edfb",
)

go_repository(
    name = "com_github_morikuni_aec",
    importpath = "github.com/morikuni/aec",
    sum = "h1:nP9CBfwrvYnBRgY6qfDQkygYDmYwOilePFkwzv4dU8A=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mrunalp_fileutils",
    importpath = "github.com/mrunalp/fileutils",
    sum = "h1:A4y2IxU1GcIzlcmUlQ6yr/mrvYZhqo+HakAPwgwaa6s=",
    version = "v0.0.0-20160930181131-4ee1cc9a8058",
)

go_repository(
    name = "com_github_munnerz_goautoneg",
    importpath = "github.com/munnerz/goautoneg",
    sum = "h1:10VrZWOtDSvWhgViCi2J6VUp4p/B3pOA/efiMH3KjjI=",
    version = "v0.0.0-20190414153302-2ae31c8b6b30",
)

go_repository(
    name = "com_github_mvdan_xurls",
    importpath = "github.com/mvdan/xurls",
    sum = "h1:OpuDelGQ1R1ueQ6sSryzi6P+1RtBpfQHM8fJwlE45ww=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_mwitkow_go_conntrack",
    importpath = "github.com/mwitkow/go-conntrack",
    sum = "h1:F9x/1yl3T2AeKLr2AMdilSD8+f9bvMnNN8VS5iDtovc=",
    version = "v0.0.0-20161129095857-cc309e4a2223",
)

go_repository(
    name = "com_github_mxk_go_flowrate",
    importpath = "github.com/mxk/go-flowrate",
    sum = "h1:y5//uYreIhSUg3J1GEMiLbxo1LJaP8RfCpH6pymGZus=",
    version = "v0.0.0-20140419014527-cca7078d478f",
)

go_repository(
    name = "com_github_nakagami_firebirdsql",
    importpath = "github.com/nakagami/firebirdsql",
    sum = "h1:P48LjvUQpTReR3TQRbxSeSBsMXzfK0uol7eRcr7VBYQ=",
    version = "v0.0.0-20190310045651-3c02a58cfed8",
)

go_repository(
    name = "com_github_naoina_go_stringutil",
    importpath = "github.com/naoina/go-stringutil",
    sum = "h1:rCUeRUHjBjGTSHl0VC00jUPLz8/F9dDzYI70Hzifhks=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_naoina_toml",
    importpath = "github.com/naoina/toml",
    sum = "h1:PT/lllxVVN0gzzSqSlHEmP8MJB4MY2U7STGxiouV4X8=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_nvveen_gotty",
    importpath = "github.com/Nvveen/Gotty",
    sum = "h1:TngWCqHvy9oXAN6lEVMRuU21PR1EtLVZJmdB18Gu3Rw=",
    version = "v0.0.0-20120604004816-cd527374f1e5",
)

go_repository(
    name = "com_github_nytimes_gziphandler",
    importpath = "github.com/NYTimes/gziphandler",
    sum = "h1:ZUDjpQae29j0ryrS0u/B8HZfJBtBQHjqw2rQ2cqUQ3I=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_oklog_run",
    importpath = "github.com/oklog/run",
    sum = "h1:Ru7dDtJNOyC66gQ5dQmaCa0qIsAUFY3sFpK1Xk8igrw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_oklog_ulid",
    importpath = "github.com/oklog/ulid",
    sum = "h1:EGfNDEx6MqHz8B3uNV6QAib1UR2Lm97sHi3ocA6ESJ4=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_oneofone_xxhash",
    importpath = "github.com/OneOfOne/xxhash",
    sum = "h1:KMrpdQIwFcEqXDklaen+P1axHaj9BSKzvpUUfnHldSE=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_onsi_ginkgo",
    importpath = "github.com/onsi/ginkgo",
    sum = "h1:q/mM8GF/n0shIN8SaAZ0V+jnLPzen6WIVZdiwrRlMlo=",
    version = "v1.10.1",
)

go_repository(
    name = "com_github_onsi_gomega",
    importpath = "github.com/onsi/gomega",
    sum = "h1:XPnZz8VVBHjVsy1vzJmRwIcSwiUO+JFfrv/xGiigmME=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_opencontainers_go_digest",
    importpath = "github.com/opencontainers/go-digest",
    sum = "h1:WzifXhOVOEOuFYOJAW6aQqW0TooG2iki3E3Ii+WN7gQ=",
    version = "v1.0.0-rc1",
)

go_repository(
    name = "com_github_opencontainers_image_spec",
    importpath = "github.com/opencontainers/image-spec",
    sum = "h1:JMemWkRwHx4Zj+fVxWoMCFm/8sYGGrUVojFA6h/TRcI=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_opencontainers_runc",
    importpath = "github.com/opencontainers/runc",
    sum = "h1:yvQ/2Pupw60ON8TYEIGGTAI77yZsWYkiOeHFZWkwlCk=",
    version = "v1.0.0-rc2.0.20190611121236-6cc515888830",
)

go_repository(
    name = "com_github_opencontainers_runtime_spec",
    importpath = "github.com/opencontainers/runtime-spec",
    sum = "h1:O6L965K88AilqnxeYPks/75HLpp4IG+FjeSCI3cVdRg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_opencontainers_selinux",
    importpath = "github.com/opencontainers/selinux",
    sum = "h1:Kx9J6eDG5/24A6DtUquGSpJQ+m2MUTahn4FtGEe8bFg=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_openshift_api",
    importpath = "github.com/openshift/api",
    sum = "h1:6il8W875Oq9vycPkRV5TteLP9IfMEX3lyOl5yN+CtdI=",
    version = "v3.9.1-0.20190924102528-32369d4db2ad+incompatible",
)

go_repository(
    name = "com_github_openshift_client_go",
    importpath = "github.com/openshift/client-go",
    sum = "h1:E++qQ7W1/EdvuMo+YGVbMPn4HihEp7YT5Rghh0VmA9A=",
    version = "v0.0.0-20190923180330-3b6373338c9b",
)

go_repository(
    name = "com_github_openshift_origin",
    importpath = "github.com/openshift/origin",
    sum = "h1:KLVRXtjLhZHVtrcdnuefaI2Bf182EEiTfEVDHokoyng=",
    version = "v0.0.0-20160503220234-8f127d736703",
)

go_repository(
    name = "com_github_openshift_prom_label_proxy",
    importpath = "github.com/openshift/prom-label-proxy",
    sum = "h1:GW8OxGwBbI2kCqjb5PQfVXRAuCJbYyX1RYs9R3ISjck=",
    version = "v0.1.1-0.20191016113035-b8153a7f39f1",
)

go_repository(
    name = "com_github_opentracing_opentracing_go",
    importpath = "github.com/opentracing/opentracing-go",
    sum = "h1:pWlfV3Bxv7k65HYwkikxat0+s3pV4bsqf19k25Ur8rU=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_openzipkin_zipkin_go",
    importpath = "github.com/openzipkin/zipkin-go",
    sum = "h1:yXiysv1CSK7Q5yjGy1710zZGnsbMUIjluWBxtLXHPBo=",
    version = "v0.1.6",
)

go_repository(
    name = "com_github_operator_framework_api",
    importpath = "github.com/operator-framework/api",
    sum = "h1:bMO43IWWPM3HCGIiuM/GyXjtSJWsrhOlzUpZMesVUw0=",
    version = "v0.0.0-20191127212340-9066a6e95573",
)

go_repository(
    name = "com_github_operator_framework_operator_lifecycle_manager",
    importpath = "github.com/operator-framework/operator-lifecycle-manager",
    sum = "h1:rjaihxY50c5C+kbQIK4s36R8zxByATYrgRbua4eiG6o=",
    version = "v0.0.0-20191115003340-16619cd27fa5",
)

go_repository(
    name = "com_github_operator_framework_operator_registry",
    importpath = "github.com/operator-framework/operator-registry",
    sum = "h1:az83WDwgB+tHsmVn+tFq72yQBbaUAye8e4+KkDQmzLs=",
    version = "v1.5.3",
)

go_repository(
    name = "com_github_operator_framework_operator_sdk",
    importpath = "github.com/operator-framework/operator-sdk",
    sum = "h1:8Aj97ZQ7mPehFR/8zselZQ1s8WSdKR4rWHhqacdVehQ=",
    version = "v0.14.1",
)

go_repository(
    name = "com_github_otiai10_copy",
    importpath = "github.com/otiai10/copy",
    sum = "h1:gtBjD8aq4nychvRZ2CyJvFWAw0aja+VHazDdruZKGZA=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_otiai10_curr",
    importpath = "github.com/otiai10/curr",
    sum = "h1:o59bHXu8Ejas8Kq6pjoVJQ9/neN66SM8AKh6wI42BBs=",
    version = "v0.0.0-20190513014714-f5a3d24e5776",
)

go_repository(
    name = "com_github_otiai10_mint",
    importpath = "github.com/otiai10/mint",
    sum = "h1:DxYL0itZyPaR5Z9HILdxSoHx+gNs6Yx+neOGS3IVUk0=",
    version = "v1.2.4",
)

go_repository(
    name = "com_github_pborman_uuid",
    importpath = "github.com/pborman/uuid",
    sum = "h1:J7Q5mO4ysT1dv8hyrUGHb9+ooztCXu1D8MY8DZYsu3g=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_pelletier_go_toml",
    importpath = "github.com/pelletier/go-toml",
    sum = "h1:aetoXYr0Tv7xRU/V4B4IZJ2QcbtMUFoNb3ORp7TzIK4=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_peterbourgon_diskv",
    importpath = "github.com/peterbourgon/diskv",
    sum = "h1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=",
    version = "v2.0.1+incompatible",
)

go_repository(
    name = "com_github_phayes_freeport",
    importpath = "github.com/phayes/freeport",
    sum = "h1:rZQtoozkfsiNs36c7Tdv/gyGNzD1X1XWKO8rptVNZuM=",
    version = "v0.0.0-20171002181615-b8543db493a5",
)

go_repository(
    name = "com_github_pierrec_lz4",
    importpath = "github.com/pierrec/lz4",
    sum = "h1:2xWsjqPFWcplujydGg4WmhC/6fZqK42wMM8aXeqhl0I=",
    version = "v2.0.5+incompatible",
)

go_repository(
    name = "com_github_pkg_errors",
    importpath = "github.com/pkg/errors",
    sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
    version = "v0.9.1",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_pquerna_cachecontrol",
    importpath = "github.com/pquerna/cachecontrol",
    sum = "h1:0XM1XL/OFFJjXsYXlG30spTkV/E9+gmd5GD1w2HE8xM=",
    version = "v0.0.0-20171018203845-0dec1b30a021",
)

go_repository(
    name = "com_github_pquerna_ffjson",
    importpath = "github.com/pquerna/ffjson",
    sum = "h1:7sBb9iOkeq+O7AXlVoH/8zpIcRXX523zMkKKspHjjx8=",
    version = "v0.0.0-20180717144149-af8b230fcd20",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    importpath = "github.com/prometheus/client_golang",
    sum = "h1:JnMpQc6ppsNgw9QPAGF6Dod479itz7lvlsMzzNayLOI=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_prometheus_client_model",
    importpath = "github.com/prometheus/client_model",
    sum = "h1:gQz4mCbXsO+nc9n1hCxHcGA3Zx3Eo+UHZoInFGUIXNM=",
    version = "v0.0.0-20190812154241-14fe0d1b01d4",
)

go_repository(
    name = "com_github_prometheus_common",
    importpath = "github.com/prometheus/common",
    sum = "h1:L+1lyG48J1zAQXA3RBX/nG/B3gjlHq0zTt2tlbJLyCY=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_prometheus_procfs",
    importpath = "github.com/prometheus/procfs",
    sum = "h1:3+auTFlqw+ZaQYJARz6ArODtkaIwtvBTx3N2NehQlL8=",
    version = "v0.0.5",
)

go_repository(
    name = "com_github_prometheus_prometheus",
    importpath = "github.com/prometheus/prometheus",
    sum = "h1:EekL1S9WPoPtJL2NZvL+xo38iMpraOnyEHOiyZygMDY=",
    version = "v2.3.2+incompatible",
)

go_repository(
    name = "com_github_prometheus_tsdb",
    importpath = "github.com/prometheus/tsdb",
    sum = "h1:w1tAGxsBMLkuGrFMhqgcCeBkM5d1YI24udArs+aASuQ=",
    version = "v0.8.0",
)

go_repository(
    name = "com_github_puerkitobio_purell",
    importpath = "github.com/PuerkitoBio/purell",
    sum = "h1:WEQqlqaGbrPkxLJWfBwQmfEAE1Z7ONdDLqrN38tNFfI=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_puerkitobio_urlesc",
    importpath = "github.com/PuerkitoBio/urlesc",
    sum = "h1:d+Bc7a5rLufV/sSk/8dngufqelfh6jnri85riMAaF/M=",
    version = "v0.0.0-20170810143723-de5bf2ad4578",
)

go_repository(
    name = "com_github_quobyte_api",
    importpath = "github.com/quobyte/api",
    sum = "h1:lPHLsuvtjFyk8WhC4uHoHRkScijIHcffTWBBP+YpzYo=",
    version = "v0.1.2",
)

go_repository(
    name = "com_github_rcrowley_go_metrics",
    importpath = "github.com/rcrowley/go-metrics",
    sum = "h1:9ZKAASQSHhDYGoxY8uLVpewe1GDZ2vu2Tr/vTdVAkFQ=",
    version = "v0.0.0-20181016184325-3113b8401b8a",
)

go_repository(
    name = "com_github_remyoudompheng_bigfft",
    importpath = "github.com/remyoudompheng/bigfft",
    sum = "h1:/NRJ5vAYoqz+7sG51ubIDHXeWO8DlTSrToPu6q11ziA=",
    version = "v0.0.0-20170806203942-52369c62f446",
)

go_repository(
    name = "com_github_rican7_retry",
    importpath = "github.com/Rican7/retry",
    sum = "h1:FqK94z34ly8Baa6K+G8Mmza9rYWTKOJk+yckIBB5qVk=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_robfig_cron",
    importpath = "github.com/robfig/cron",
    sum = "h1:jk4/Hud3TTdcrJgUOBgsqrZBarcxl6ADIjSC2iniwLY=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_rogpeppe_fastuuid",
    importpath = "github.com/rogpeppe/fastuuid",
    sum = "h1:gu+uRPtBe88sKxUCEXRoeCvVG90TJmwhiqRpvdhQFng=",
    version = "v0.0.0-20150106093220-6724a57986af",
)

go_repository(
    name = "com_github_rogpeppe_go_internal",
    importpath = "github.com/rogpeppe/go-internal",
    sum = "h1:Usqs0/lDK/NqTkvrmKSwA/3XkZAs7ZAW/eLeQ2MVBTw=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_rubenv_sql_migrate",
    importpath = "github.com/rubenv/sql-migrate",
    sum = "h1:lwDYefgiwhjuAuVnMVUYknoF+Yg9CBUykYGvYoPCNnQ=",
    version = "v0.0.0-20191025130928-9355dd04f4b3",
)

go_repository(
    name = "com_github_rubiojr_go_vhd",
    importpath = "github.com/rubiojr/go-vhd",
    sum = "h1:ht7N4d/B7Ezf58nvMNVF3OlvDlz9pp+WHVcRNS0nink=",
    version = "v0.0.0-20160810183302-0bfd3b39853c",
)

go_repository(
    name = "com_github_russross_blackfriday",
    importpath = "github.com/russross/blackfriday",
    sum = "h1:HyvC0ARfnZBqnXwABFeSZHpKvJHJJfPz81GNueLj0oo=",
    version = "v1.5.2",
)

go_repository(
    name = "com_github_satori_go_uuid",
    importpath = "github.com/satori/go.uuid",
    sum = "h1:0uYX9dsZ2yD7q2RtLRtPSdGDWzjeM3TbMJP9utgA0ww=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_sclevine_spec",
    importpath = "github.com/sclevine/spec",
    sum = "h1:1Jwdf9jSfDl9NVmt8ndHqbTZ7XCCPbh1jI3hkDBHVYA=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_seccomp_libseccomp_golang",
    importpath = "github.com/seccomp/libseccomp-golang",
    sum = "h1:NJjM5DNFOs0s3kYE1WUOr6G8V97sdt46rlXTMfXGWBo=",
    version = "v0.9.1",
)

go_repository(
    name = "com_github_sergi_go_diff",
    importpath = "github.com/sergi/go-diff",
    sum = "h1:Kpca3qRNrduNnOQeazBd0ysaKrUJiIuISHxogkT9RPQ=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_shopify_logrus_bugsnag",
    importpath = "github.com/Shopify/logrus-bugsnag",
    sum = "h1:UrqY+r/OJnIp5u0s1SbQ8dVfLCZJsnvazdBP5hS4iRs=",
    version = "v0.0.0-20171204204709-577dee27f20d",
)

go_repository(
    name = "com_github_shopify_sarama",
    importpath = "github.com/Shopify/sarama",
    sum = "h1:9oksLxC6uxVPHPVYUmq6xhr1BOF/hHobWH2UzO67z1s=",
    version = "v1.19.0",
)

go_repository(
    name = "com_github_shopify_toxiproxy",
    importpath = "github.com/Shopify/toxiproxy",
    sum = "h1:TKdv8HiTLgE5wdJuEML90aBgNWsokNbMijUGhmcoBJc=",
    version = "v2.1.4+incompatible",
)

go_repository(
    name = "com_github_shopspring_decimal",
    importpath = "github.com/shopspring/decimal",
    sum = "h1:Gojs/hac/DoYEM7WEICT45+hNWczIeuL5D21e5/HPAw=",
    version = "v0.0.0-20191009025716-f1972eb1d1f5",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:SPIRibHv4MatM3XXNO2BJeFLZwZ2LvZgfQ5+UNI2im4=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_smartystreets_assertions",
    importpath = "github.com/smartystreets/assertions",
    sum = "h1:zE9ykElWQ6/NYmHa3jpm/yHnI4xSofP+UP6SpjHcSeM=",
    version = "v0.0.0-20180927180507-b2de0cb4f26d",
)

go_repository(
    name = "com_github_smartystreets_goconvey",
    importpath = "github.com/smartystreets/goconvey",
    sum = "h1:pa8hGb/2YqsZKovtsgrwcDH1RZhVbTKCjLp47XpqCDs=",
    version = "v0.0.0-20190330032615-68dc04aab96a",
)

go_repository(
    name = "com_github_soheilhy_cmux",
    importpath = "github.com/soheilhy/cmux",
    sum = "h1:0HKaf1o97UwFjHH9o5XsHUOF+tqmdA7KEzXLpiyaw0E=",
    version = "v0.1.4",
)

go_repository(
    name = "com_github_spaolacci_murmur3",
    importpath = "github.com/spaolacci/murmur3",
    sum = "h1:qLC7fQah7D6K1B0ujays3HV9gkFtllcxhzImRR7ArPQ=",
    version = "v0.0.0-20180118202830-f09979ecbc72",
)

go_repository(
    name = "com_github_spf13_afero",
    importpath = "github.com/spf13/afero",
    sum = "h1:5jhuqJyZCZf2JRofRvN/nIFgIWNzPa3/Vz8mYylgbWc=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_spf13_cast",
    importpath = "github.com/spf13/cast",
    sum = "h1:oget//CVOEoFewqQxwr0Ej5yjygnqGkvggSE/gB35Q8=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_spf13_cobra",
    importpath = "github.com/spf13/cobra",
    sum = "h1:f0B+LkLX6DtmRH1isoNA9VTtNUK9K8xYd28JNNfOv/s=",
    version = "v0.0.5",
)

go_repository(
    name = "com_github_spf13_jwalterweatherman",
    importpath = "github.com/spf13/jwalterweatherman",
    sum = "h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_spf13_pflag",
    importpath = "github.com/spf13/pflag",
    sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
    version = "v1.0.5",
)

go_repository(
    name = "com_github_spf13_viper",
    importpath = "github.com/spf13/viper",
    sum = "h1:yXHLWeravcrgGyFSyCgdYpXQ9dR9c/WED3pg1RhxqEU=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_storageos_go_api",
    importpath = "github.com/storageos/go-api",
    sum = "h1:n+WYaU0kQ6WIiuEyWSgbXqkBx16irO69kYCtwVYoO5s=",
    version = "v0.0.0-20180912212459-343b3eff91fc",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:Hbg2NidpLE8veEBkEZTL3CvlkUIVzuU9jDplZO54c48=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:nOGnQDM7FYENwehXlg/kFVnos3rEvtKTjRvOWSzb6H4=",
    version = "v1.5.1",
)

go_repository(
    name = "com_github_syndtr_gocapability",
    importpath = "github.com/syndtr/gocapability",
    sum = "h1:w58e6FAOMd+rUgOfhaBb+ZVOQIOfUkpv5AAQVmf6hsI=",
    version = "v0.0.0-20160928074757-e7cb7fa329f4",
)

go_repository(
    name = "com_github_thecodeteam_goscaleio",
    importpath = "github.com/thecodeteam/goscaleio",
    sum = "h1:SB5tO98lawC+UK8ds/U2jyfOCH7GTcFztcF5x9gbut4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_tidwall_pretty",
    importpath = "github.com/tidwall/pretty",
    sum = "h1:BP2bjP495BBPaBcS5rmqviTfrOkN5rO5ceKAMRZCRFc=",
    version = "v0.0.0-20180105212114-65a9db5fad51",
)

go_repository(
    name = "com_github_tmc_grpc_websocket_proxy",
    importpath = "github.com/tmc/grpc-websocket-proxy",
    sum = "h1:LnC5Kc/wtumK+WB441p7ynQJzVuNRJiqddSIE3IlSEQ=",
    version = "v0.0.0-20190109142713-0ad062ec5ee5",
)

go_repository(
    name = "com_github_ugorji_go",
    importpath = "github.com/ugorji/go",
    sum = "h1:j4s+tAvLfL3bZyefP2SEWmhBzmuIlH/eqNuPdFPgngw=",
    version = "v1.1.4",
)

go_repository(
    name = "com_github_ugorji_go_codec",
    importpath = "github.com/ugorji/go/codec",
    sum = "h1:3SVOIvH7Ae1KRYyQWRjXWJEA9sS/c/pjvH++55Gr648=",
    version = "v0.0.0-20181204163529-d75b2dcb6bc8",
)

go_repository(
    name = "com_github_urfave_negroni",
    importpath = "github.com/urfave/negroni",
    sum = "h1:kIimOitoypq34K7TG7DUaJ9kq/N4Ofuwi1sjz0KipXc=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_vishvananda_netlink",
    importpath = "github.com/vishvananda/netlink",
    sum = "h1:f1yevOHP+Suqk0rVc13fIkzcLULJbyQcXDba2klljD0=",
    version = "v0.0.0-20171020171820-b2de5d10e38e",
)

go_repository(
    name = "com_github_vishvananda_netns",
    importpath = "github.com/vishvananda/netns",
    sum = "h1:J9gO8RJCAFlln1jsvRba/CWVUnMHwObklfxxjErl1uk=",
    version = "v0.0.0-20171111001504-be1fbeda1936",
)

go_repository(
    name = "com_github_vmware_govmomi",
    importpath = "github.com/vmware/govmomi",
    sum = "h1:7b/SeTUB3tER8ZLGLLLH3xcnB2xeuLULXmfPFqPSRZA=",
    version = "v0.20.1",
)

go_repository(
    name = "com_github_xanzy_go_gitlab",
    importpath = "github.com/xanzy/go-gitlab",
    sum = "h1:rWtwKTgEnXyNUGrOArN7yyc3THRkpYcKXIXia9abywQ=",
    version = "v0.15.0",
)

go_repository(
    name = "com_github_xdg_scram",
    importpath = "github.com/xdg/scram",
    sum = "h1:u40Z8hqBAAQyv+vATcGgV0YCnDjqSL7/q/JyPhhJSPk=",
    version = "v0.0.0-20180814205039-7eeb5667e42c",
)

go_repository(
    name = "com_github_xdg_stringprep",
    importpath = "github.com/xdg/stringprep",
    sum = "h1:d9X0esnoa3dFsV0FG35rAT0RIhYFlPq7MiP+DW89La0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_xeipuuv_gojsonpointer",
    importpath = "github.com/xeipuuv/gojsonpointer",
    sum = "h1:J9EGpcZtP0E/raorCMxlFGSTBrsSlaDGf3jU/qvAE2c=",
    version = "v0.0.0-20180127040702-4e3ac2762d5f",
)

go_repository(
    name = "com_github_xeipuuv_gojsonreference",
    importpath = "github.com/xeipuuv/gojsonreference",
    sum = "h1:EzJWgHovont7NscjpAxXsDA8S8BMYve8Y5+7cuRE7R0=",
    version = "v0.0.0-20180127040603-bd5ef7bd5415",
)

go_repository(
    name = "com_github_xeipuuv_gojsonschema",
    importpath = "github.com/xeipuuv/gojsonschema",
    sum = "h1:ngVtJC9TY/lg0AA/1k48FYhBrhRoFlEmWzsehpNAaZg=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_xenolf_lego",
    importpath = "github.com/xenolf/lego",
    sum = "h1:BTvU+npm3/yjuBd53EvgiFLl5+YLikf2WvHsjRQ4KrY=",
    version = "v0.3.2-0.20160613233155-a9d8cec0e656",
)

go_repository(
    name = "com_github_xiang90_probing",
    importpath = "github.com/xiang90/probing",
    sum = "h1:eY9dn8+vbi4tKz5Qo6v2eYzo7kUS51QINcR5jNpbZS8=",
    version = "v0.0.0-20190116061207-43a291ad63a2",
)

go_repository(
    name = "com_github_xlab_handysort",
    importpath = "github.com/xlab/handysort",
    sum = "h1:j2hhcujLRHAg872RWAV5yaUrEjHEObwDv3aImCaNLek=",
    version = "v0.0.0-20150421192137-fb3537ed64a1",
)

go_repository(
    name = "com_github_xordataexchange_crypt",
    importpath = "github.com/xordataexchange/crypt",
    sum = "h1:ESFSdwYZvkeru3RtdrYueztKhOBCSAAzS4Gf+k0tEow=",
    version = "v0.0.3-0.20170626215501-b2862e3d0a77",
)

go_repository(
    name = "com_github_yvasiyarov_go_metrics",
    importpath = "github.com/yvasiyarov/go-metrics",
    sum = "h1:p7OofyZ509h8DmPLh8Hn+EIIZm/xYhdZHJ9GnXHdr6U=",
    version = "v0.0.0-20150112132944-c25f46c4b940",
)

go_repository(
    name = "com_github_yvasiyarov_gorelic",
    importpath = "github.com/yvasiyarov/gorelic",
    sum = "h1:qMJQYPNdtJ7UNYHjX38KXZtltKTqimMuoQjNnSVIuJg=",
    version = "v0.0.6",
)

go_repository(
    name = "com_github_yvasiyarov_newrelic_platform_go",
    importpath = "github.com/yvasiyarov/newrelic_platform_go",
    sum = "h1:ERexzlUfuTvpE74urLSbIQW0Z/6hF9t8U4NsJLaioAY=",
    version = "v0.0.0-20140908184405-b21fdbd4370f",
)

go_repository(
    name = "com_github_ziutek_mymysql",
    importpath = "github.com/ziutek/mymysql",
    sum = "h1:GB0qdRGsTwQSBVYuVShFBKaXSnSnYYC2d9knnE1LHFs=",
    version = "v1.5.4",
)

go_repository(
    name = "com_gitlab_nyarla_go_crypt",
    importpath = "gitlab.com/nyarla/go-crypt",
    sum = "h1:7gd+rd8P3bqcn/96gOZa3F5dpJr/vEiDQYlNb/y2uNs=",
    version = "v0.0.0-20160106005555-d9a5dc2b789b",
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    sum = "h1:lRi0CHyU+ytlvylOlFKKq0af6JncuyoRh1J+QJBqQx0=",
    version = "v0.45.1",
)

go_repository(
    name = "in_gopkg_airbrake_gobrake_v2",
    importpath = "gopkg.in/airbrake/gobrake.v2",
    sum = "h1:7z2uVWwn7oVeeugY1DtlPAy5H+KYgB1KeKTnqjNatLo=",
    version = "v2.0.9",
)

go_repository(
    name = "in_gopkg_alecthomas_kingpin_v2",
    importpath = "gopkg.in/alecthomas/kingpin.v2",
    sum = "h1:jMFz6MfLP0/4fUyZle81rXUoxOBFi19VUFKVDOQfozc=",
    version = "v2.2.6",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=",
    version = "v1.0.0-20190902080502-41f04d3bba15",
)

go_repository(
    name = "in_gopkg_errgo_v2",
    importpath = "gopkg.in/errgo.v2",
    sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
    version = "v2.1.0",
)

go_repository(
    name = "in_gopkg_fsnotify_v1",
    importpath = "gopkg.in/fsnotify.v1",
    sum = "h1:xOHLXZwVvI9hhs+cLKq5+I5onOuwQLhQwiu63xxlHs4=",
    version = "v1.4.7",
)

go_repository(
    name = "in_gopkg_gcfg_v1",
    importpath = "gopkg.in/gcfg.v1",
    sum = "h1:0HIbH907iBTAntm+88IJV2qmJALDAh8sPekI9Vc1fm0=",
    version = "v1.2.0",
)

go_repository(
    name = "in_gopkg_gemnasium_logrus_airbrake_hook_v2",
    importpath = "gopkg.in/gemnasium/logrus-airbrake-hook.v2",
    sum = "h1:OAj3g0cR6Dx/R07QgQe8wkA9RNjB2u4i700xBkIT4e0=",
    version = "v2.1.2",
)

go_repository(
    name = "in_gopkg_gorp_v1",
    importpath = "gopkg.in/gorp.v1",
    sum = "h1:j3DWlAyGVv8whO7AcIWznQ2Yj7yJkn34B8s63GViAAw=",
    version = "v1.7.2",
)

go_repository(
    name = "in_gopkg_inf_v0",
    importpath = "gopkg.in/inf.v0",
    sum = "h1:73M5CoZyi3ZLMOyDlQh031Cx6N9NDJ2Vvfl76EDAgDc=",
    version = "v0.9.1",
)

go_repository(
    name = "in_gopkg_mcuadros_go_syslog_v2",
    importpath = "gopkg.in/mcuadros/go-syslog.v2",
    sum = "h1:60g8zx1BijSVSgLTzLCW9UC4/+i1Ih9jJ1DR5Tgp9vE=",
    version = "v2.2.1",
)

go_repository(
    name = "in_gopkg_natefinch_lumberjack_v2",
    importpath = "gopkg.in/natefinch/lumberjack.v2",
    sum = "h1:1Lc07Kr7qY4U2YPouBjpCLxpiyxIVoxqXgkXLknAOE8=",
    version = "v2.0.0",
)

go_repository(
    name = "in_gopkg_resty_v1",
    importpath = "gopkg.in/resty.v1",
    sum = "h1:CuXP0Pjfw9rOuY6EP+UvtNvt5DSqHpIxILZKT/quCZI=",
    version = "v1.12.0",
)

go_repository(
    name = "in_gopkg_square_go_jose_v1",
    importpath = "gopkg.in/square/go-jose.v1",
    sum = "h1:/5jmADZB+RiKtZGr4HxsEFOEfbfsjTKsVnqpThUpE30=",
    version = "v1.1.2",
)

go_repository(
    name = "in_gopkg_square_go_jose_v2",
    importpath = "gopkg.in/square/go-jose.v2",
    sum = "h1:orlkJ3myw8CN1nVQHBFfloD+L3egixIa4FvUP6RosSA=",
    version = "v2.2.2",
)

go_repository(
    name = "in_gopkg_tomb_v1",
    importpath = "gopkg.in/tomb.v1",
    sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
    version = "v1.0.0-20141024135613-dd632973f1e7",
)

go_repository(
    name = "in_gopkg_warnings_v0",
    importpath = "gopkg.in/warnings.v0",
    sum = "h1:XM28wIgFzaBmeZ5dNHIpWLQpt/9DGKxk+rCg/22nnYE=",
    version = "v0.1.1",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:obN1ZagJSUGI0Ek/LBmuj4SNLPfIny3KsKFopxRdj10=",
    version = "v2.2.8",
)

go_repository(
    name = "in_gopkg_yaml.v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:/eiJrUcujPVeJ3xlSWaiNi3uSVmDGBK1pDHUHAnao1I=",
    version = "v2.2.4",
)

go_repository(
    name = "in_gopkg_yaml_v3",
    importpath = "gopkg.in/yaml.v3",
    sum = "h1:XZx7nhd5GMaZpmDaEHFVafUZC7ya0fuo7cSJ3UCKYmM=",
    version = "v3.0.0-20191120175047-4206685974f2",
)

go_repository(
    name = "io_etcd_go_bbolt",
    importpath = "go.etcd.io/bbolt",
    sum = "h1:Z/90sZLPOeCy2PwprqkFa25PdkusRzaj9P8zm/KNyvk=",
    version = "v1.3.2",
)

go_repository(
    name = "io_k8s_api",
    build_file_proto_mode = "disable_global",
    importpath = "k8s.io/api",
    replace = "k8s.io/api",
    sum = "h1:VVUE9xTCXP6KUPMf92cQmN88orz600ebexcRRaBTepQ=",
    version = "v0.0.0-20191016110408-35e52d86657a",
)

go_repository(
    name = "io_k8s_apiextensions_apiserver",
    build_file_proto_mode = "disable_global",
    importpath = "k8s.io/apiextensions-apiserver",
    replace = "k8s.io/apiextensions-apiserver",
    sum = "h1:kThoiqgMsSwBdMK/lPgjtYTsEjbUU9nXCA9DyU3feok=",
    version = "v0.0.0-20191016113550-5357c4baaf65",
)

go_repository(
    name = "io_k8s_apimachinery",
    build_file_proto_mode = "disable_global",
    importpath = "k8s.io/apimachinery",
    replace = "k8s.io/apimachinery",
    sum = "h1:Iieh/ZEgT3BWwbLD5qEKcY06jKuPEl6zC7gPSehoLw4=",
    version = "v0.0.0-20191004115801-a2eda9f80ab8",
)

go_repository(
    name = "io_k8s_apiserver",
    importpath = "k8s.io/apiserver",
    replace = "k8s.io/apiserver",
    sum = "h1:leksCBKKBrPJmW1jV4dZUvwqmVtXpKdzpHsqXfFS094=",
    version = "v0.0.0-20191016112112-5190913f932d",
)

go_repository(
    name = "io_k8s_autoscaler",
    importpath = "k8s.io/autoscaler",
    sum = "h1:VZ7xYyNfl07Yh3DK10Dck9r0/8cQ3vPSxlrSQt7tohk=",
    version = "v0.0.0-20190607113959-1b4f1855cb8e",
)

go_repository(
    name = "io_k8s_cli_runtime",
    importpath = "k8s.io/cli-runtime",
    replace = "k8s.io/cli-runtime",
    sum = "h1:8ZfMjkMBzcXEawLsYHg9lDM7aLEVso3NiVKfUTnN56A=",
    version = "v0.0.0-20191016114015-74ad18325ed5",
)

go_repository(
    name = "io_k8s_client_go",
    importpath = "k8s.io/client-go",
    replace = "k8s.io/client-go",
    sum = "h1:C2XVy2z0dV94q9hSSoCuTPp1KOG7IegvbdXuz9VGxoU=",
    version = "v0.0.0-20191016111102-bec269661e48",
)

go_repository(
    name = "io_k8s_cloud_provider",
    importpath = "k8s.io/cloud-provider",
    replace = "k8s.io/cloud-provider",
    sum = "h1:rP/89rnWN2l+2b7Jckg4VXi2dhgu7xs3S+1bKWKrqGE=",
    version = "v0.0.0-20191016115326-20453efc2458",
)

go_repository(
    name = "io_k8s_cluster_bootstrap",
    importpath = "k8s.io/cluster-bootstrap",
    replace = "k8s.io/cluster-bootstrap",
    sum = "h1:ZwG8XnuF+Z4Qmc/XfhFXgbhfgr6YPmVqFbCRNwLG+G8=",
    version = "v0.0.0-20191016115129-c07a134afb42",
)

go_repository(
    name = "io_k8s_code_generator",
    importpath = "k8s.io/code-generator",
    replace = "k8s.io/code-generator",
    sum = "h1:NMYlxaF7rYQJk2E2IyrUhaX81zX24+dmoZdkPw0gJqI=",
    version = "v0.0.0-20191004115455-8e001e5d1894",
)

go_repository(
    name = "io_k8s_component_base",
    importpath = "k8s.io/component-base",
    replace = "k8s.io/component-base",
    sum = "h1:2D+G/CCNVdYc0h9D+tX+0SmtcyQmby6uzNityrps1s0=",
    version = "v0.0.0-20191016111319-039242c015a9",
)

go_repository(
    name = "io_k8s_cri_api",
    importpath = "k8s.io/cri-api",
    replace = "k8s.io/cri-api",
    sum = "h1:ikDtGPX1DVIhl4E36+khq6RVyA65ycfiieBHecQiaX0=",
    version = "v0.0.0-20190828162817-608eb1dad4ac",
)

go_repository(
    name = "io_k8s_csi_translation_lib",
    importpath = "k8s.io/csi-translation-lib",
    replace = "k8s.io/csi-translation-lib",
    sum = "h1:8St7hlu0fkur/6TRtIYgTqjNGvxFqcTxKywmlAvMiVo=",
    version = "v0.0.0-20191016115521-756ffa5af0bd",
)

go_repository(
    name = "io_k8s_gengo",
    importpath = "k8s.io/gengo",
    sum = "h1:eW/6wVuHNZgQJmFesyAxu0cvj0WAHHUuGaLbPcmNY3Q=",
    version = "v0.0.0-20191010091904-7fa3014cb28f",
)

go_repository(
    name = "io_k8s_heapster",
    importpath = "k8s.io/heapster",
    sum = "h1:lUsE/AHOMHpi3MLlBEkaU8Esxm5QhdyCrv1o7ot0s84=",
    version = "v1.2.0-beta.1",
)

go_repository(
    name = "io_k8s_helm",
    importpath = "k8s.io/helm",
    sum = "h1:L+k810plJlaGWEw1EszeT4deK8XVaKxac1oGcuB+WDc=",
    version = "v2.16.1+incompatible",
)

go_repository(
    name = "io_k8s_klog",
    importpath = "k8s.io/klog",
    sum = "h1:Pt+yjF5aB1xDSVbau4VsWe+dQNzA0qv1LlXdC2dF6Q8=",
    version = "v1.0.0",
)

go_repository(
    name = "io_k8s_kube_aggregator",
    importpath = "k8s.io/kube-aggregator",
    replace = "k8s.io/kube-aggregator",
    sum = "h1:Tv+DHbQg2ozCJqmuw5poFX7sxs2mJPUm7MEz3sQEULM=",
    version = "v0.0.0-20191016112429-9587704a8ad4",
)

go_repository(
    name = "io_k8s_kube_controller_manager",
    importpath = "k8s.io/kube-controller-manager",
    replace = "k8s.io/kube-controller-manager",
    sum = "h1:7zqJKTBHA7+oFKu6FLB/di2/zrx+2Khx1hBKJ5oOBcc=",
    version = "v0.0.0-20191016114939-2b2b218dc1df",
)

go_repository(
    name = "io_k8s_kube_openapi",
    importpath = "k8s.io/kube-openapi",
    sum = "h1:Xpe6sK+RY4ZgCTyZ3y273UmFmURhjtoJiwOMbQsXitY=",
    version = "v0.0.0-20190918143330-0270cf2f1c1d",
)

go_repository(
    name = "io_k8s_kube_proxy",
    importpath = "k8s.io/kube-proxy",
    replace = "k8s.io/kube-proxy",
    sum = "h1:4aqKRZx9D18lcLYHeETB6BBYK+Yr+oWV0gRGE1X0wM8=",
    version = "v0.0.0-20191016114407-2e83b6f20229",
)

go_repository(
    name = "io_k8s_kube_scheduler",
    importpath = "k8s.io/kube-scheduler",
    replace = "k8s.io/kube-scheduler",
    sum = "h1:dFyxN/1nxwm8+GCeRJRZDhmH5upr7r/zY7BuY5dJ4Co=",
    version = "v0.0.0-20191016114748-65049c67a58b",
)

go_repository(
    name = "io_k8s_kube_state_metrics",
    importpath = "k8s.io/kube-state-metrics",
    sum = "h1:6vdtgXrrRRMSgnyDmgua+qvgCYv954JNfxXAtDkeLVQ=",
    version = "v1.7.2",
)

go_repository(
    name = "io_k8s_kubectl",
    importpath = "k8s.io/kubectl",
    replace = "k8s.io/kubectl",
    sum = "h1:RBkTKVMF+xsNsSOVc0+HdC0B5gD1sr6s6Cu5w9qNbuQ=",
    version = "v0.0.0-20191016120415-2ed914427d51",
)

go_repository(
    name = "io_k8s_kubelet",
    importpath = "k8s.io/kubelet",
    replace = "k8s.io/kubelet",
    sum = "h1:YXArqZfchiY+62+AyWPWE59wICh7xnAEowHGWggxBXs=",
    version = "v0.0.0-20191016114556-7841ed97f1b2",
)

go_repository(
    name = "io_k8s_kubernetes",
    importpath = "k8s.io/kubernetes",
    sum = "h1:k0f/OVp6Yfv+UMTm6VYKhqjRgcvHh4QhN9coanjrito=",
    version = "v1.16.2",
)

go_repository(
    name = "io_k8s_legacy_cloud_providers",
    importpath = "k8s.io/legacy-cloud-providers",
    replace = "k8s.io/legacy-cloud-providers",
    sum = "h1:cgLCVtQnxjALxIUjjEkiMaKlQZW5sGj6P3+3K5Y/d+8=",
    version = "v0.0.0-20191016115753-cf0698c3a16b",
)

go_repository(
    name = "io_k8s_metrics",
    importpath = "k8s.io/metrics",
    replace = "k8s.io/metrics",
    sum = "h1:VbAmCGT95GvCZaGtW3oLhf7d2FXT+lnWiTtPE8hzPVo=",
    version = "v0.0.0-20191016113814-3b1a734dba6e",
)

go_repository(
    name = "io_k8s_repo_infra",
    importpath = "k8s.io/repo-infra",
    sum = "h1:WD6cPA3q7qxZe6Fwir0XjjGwGMaWbHlHUcjCcOzuRG0=",
    version = "v0.0.0-20181204233714-00fe14e3d1a3",
)

go_repository(
    name = "io_k8s_sample_apiserver",
    importpath = "k8s.io/sample-apiserver",
    replace = "k8s.io/sample-apiserver",
    sum = "h1:8hWqyqHVeKxjwKYfuo6gcq6bag5snqWh9MQk7WRrY9g=",
    version = "v0.0.0-20191016112829-06bb3c9d77c9",
)

go_repository(
    name = "io_k8s_sigs_controller_runtime",
    importpath = "sigs.k8s.io/controller-runtime",
    sum = "h1:wATM6/m+3w8lj8FXNaO6Fs/rq/vqoOjO1Q116Z9NPsg=",
    version = "v0.4.0",
)

go_repository(
    name = "io_k8s_sigs_controller_tools",
    importpath = "sigs.k8s.io/controller-tools",
    sum = "h1:la1h46EzElvWefWLqfsXrnsO3lZjpkI0asTpX6h8PLA=",
    version = "v0.2.4",
)

go_repository(
    name = "io_k8s_sigs_kustomize",
    importpath = "sigs.k8s.io/kustomize",
    sum = "h1:JUufWFNlI44MdtnjUqVnvh29rR37PQFzPbLXqhyOyX0=",
    version = "v2.0.3+incompatible",
)

go_repository(
    name = "io_k8s_sigs_structured_merge_diff",
    importpath = "sigs.k8s.io/structured-merge-diff",
    sum = "h1:6dsH6AYQWbyZmtttJNe8Gq1cXOeS1BdV3eW37zHilAQ=",
    version = "v0.0.0-20190817042607-6149e4549fca",
)

go_repository(
    name = "io_k8s_sigs_testing_frameworks",
    importpath = "sigs.k8s.io/testing_frameworks",
    sum = "h1:vK0+tvjF0BZ/RYFeZ1E6BYBwHJJXhjuZ3TdsEKH+UQM=",
    version = "v0.1.2",
)

go_repository(
    name = "io_k8s_sigs_yaml",
    importpath = "sigs.k8s.io/yaml",
    sum = "h1:4A07+ZFc2wgJwo8YNlQpr1rVlgUDlxXHhPJciaPY5gs=",
    version = "v1.1.0",
)

go_repository(
    name = "io_k8s_utils",
    importpath = "k8s.io/utils",
    sum = "h1:Gi+/O1saihwDqnlmC8Vhv1M5Sp4+rbOmK9TbsLn8ZEA=",
    version = "v0.0.0-20191010214722-8d271d903fe4",
)

go_repository(
    name = "io_opencensus_go",
    importpath = "go.opencensus.io",
    sum = "h1:C9hSCOW830chIVkdja34wa6Ky+IzWllkUinR+BtRZd4=",
    version = "v0.22.0",
)

go_repository(
    name = "io_rsc_letsencrypt",
    importpath = "rsc.io/letsencrypt",
    sum = "h1:DV0d09Ne9E7UUa9ZqWktZ9L2VmybgTgfq7xlfFR/bbU=",
    version = "v0.0.1",
)

go_repository(
    name = "ke_bou_monkey",
    importpath = "bou.ke/monkey",
    sum = "h1:zEMLInw9xvNakzUUPjfS4Ds6jYPqCFx3m7bRmG5NH2U=",
    version = "v1.0.1",
)

go_repository(
    name = "ml_vbom_util",
    importpath = "vbom.ml/util",
    sum = "h1:MksmcCZQWAQJCTA5T0jgI/0sJ51AVm4Z41MrmfczEoc=",
    version = "v0.0.0-20160121211510-db5cfe13f5cc",
)

go_repository(
    name = "org_bitbucket_bertimus9_systemstat",
    importpath = "bitbucket.org/bertimus9/systemstat",
    sum = "h1:N9r8OBSXAgEUfho3SQtZLY8zo6E1OdOMvelvP22aVFc=",
    version = "v0.0.0-20180207000608-0eeff89b0690",
)

go_repository(
    name = "org_golang_google_api",
    importpath = "google.golang.org/api",
    sum = "h1:jbyannxz0XFD3zdjgrSUsaJbgpH4eTrkdhRChkHPfO8=",
    version = "v0.9.0",
)

go_repository(
    name = "org_golang_google_appengine",
    importpath = "google.golang.org/appengine",
    sum = "h1:tycE03LOZYQNhDpS27tcQdAzLCVMaj7QT2SXxebnpCM=",
    version = "v1.6.5",
)

go_repository(
    name = "org_golang_google_genproto",
    importpath = "google.golang.org/genproto",
    sum = "h1:jB9+PJSvu5tBfmJHy/OVapFdjDF3WvpkqRhxqrmzoEU=",
    version = "v0.0.0-20200218151345-dad8c97a84f5",
)

go_repository(
    name = "org_golang_google_grpc",
    importpath = "google.golang.org/grpc",
    sum = "h1:zvIju4sqAGvwKspUQOhwnpcqSbzi7/H6QomNNjTL4sk=",
    version = "v1.27.1",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:1ZiEyfaQIg3Qh0EoqpwAakHVhecoE5wlSg5GjnafJGw=",
    version = "v0.0.0-20200221231518-2aa609cf4a9d",
)

go_repository(
    name = "org_golang_x",
    importpath = "golang.org/x/crypto",
    sum = "h1:ZC1Xn5A1nlpSmQCIva4bZ3ob3lmhYIefc+GU+DLg1Ow=",
    version = "v0.0.0-20191028145041-f83a4685e152",
)

go_repository(
    name = "org_golang_x_exp",
    importpath = "golang.org/x/exp",
    sum = "h1:OeRHuibLsmZkFj773W4LcfAGsSxJgfPONhr8cmO+eLA=",
    version = "v0.0.0-20190510132918-efd6b22b2522",
)

go_repository(
    name = "org_golang_x_image",
    importpath = "golang.org/x/image",
    sum = "h1:KYGJGHOQy8oSi1fDlSpcZF0+juKwk/hEMv5SiwHogR0=",
    version = "v0.0.0-20190227222117-0694c2d4d067",
)

go_repository(
    name = "org_golang_x_lint",
    importpath = "golang.org/x/lint",
    sum = "h1:QzoH/1pFpZguR8NrRHLcO6jKqfv2zpuSqZLgdm7ZmjI=",
    version = "v0.0.0-20190409202823-959b441ac422",
)

go_repository(
    name = "org_golang_x_mobile",
    importpath = "golang.org/x/mobile",
    sum = "h1:Tus/Y4w3V77xDsGwKUC8a/QrV7jScpU557J77lFffNs=",
    version = "v0.0.0-20190312151609-d3739f865fa6",
)

go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:KU7oHjnv3XNWfa5COkzUifxZmxp1TyI7ImMXqFxLwvQ=",
    version = "v0.2.0",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:0mm1VjtFUOIlE1SbDlwjYaDxZVDP2S5ou6y0gSgXHu8=",
    version = "v0.0.0-20200226121028-0de0cce0169b",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    sum = "h1:SVwTIAaPC2U/AvvLNZ2a7OVsmBpC8L5BlwK1whH3hm0=",
    version = "v0.0.0-20190604053449-0f29369cfe45",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:vcxGaoTs7kV8m5Np9uUNQin4BrLOthgV7252N8V+FwY=",
    version = "v0.0.0-20190911185100-cd5d95a43a6e",
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:/WDfKMnPU+m5M4xB+6x4kaepxRw6jWvR5iDRdvjHgy8=",
    version = "v0.0.0-20200223170610-d5e6a3e2c0ae",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:tW2bmiBqwgJj/UpqtC8EpXEZVYOwU0yG4iWbprSVAcs=",
    version = "v0.3.2",
)

go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    sum = "h1:NbrfHxef+IfdI86qCgO/1Siq1BuMH2xG0NqgvCguRhQ=",
    version = "v0.0.0-20200309202150-20ab64c0d93f",
)

go_repository(
    name = "org_golang_x_xerrors",
    importpath = "golang.org/x/xerrors",
    sum = "h1:E7g+9GITq07hpfrRu66IVDexMakfv52eLZ2CXBWiKr4=",
    version = "v0.0.0-20191204190536-9bdfabe68543",
)

go_repository(
    name = "org_gonum_v1_gonum",
    importpath = "gonum.org/v1/gonum",
    sum = "h1:2qZ38BsejXrhuetzb8UxucqrWDZKjypFSZA82hLCpZ4=",
    version = "v0.0.0-20190710053202-4340aa3071a0",
)

go_repository(
    name = "org_gonum_v1_netlib",
    importpath = "gonum.org/v1/netlib",
    sum = "h1:jRyg0XfpwWlhEV8mDfdNGBeSJM2fuyh9Yjrnd8kF2Ts=",
    version = "v0.0.0-20190331212654-76723241ea4e",
)

go_repository(
    name = "org_modernc_cc",
    importpath = "modernc.org/cc",
    sum = "h1:nPibNuDEx6tvYrUAtvDTTw98rx5juGsa5zuDnKwEEQQ=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_golex",
    importpath = "modernc.org/golex",
    sum = "h1:wWpDlbK8ejRfSyi0frMyhilD3JBvtcx2AdGDnU+JtsE=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_mathutil",
    importpath = "modernc.org/mathutil",
    sum = "h1:93vKjrJopTPrtTNpZ8XIovER7iCIH1QU7wNbOQXC60I=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_strutil",
    importpath = "modernc.org/strutil",
    sum = "h1:XVFtQwFVwc02Wk+0L/Z/zDDXO81r5Lhe6iMKmGX3KhE=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_xc",
    importpath = "modernc.org/xc",
    sum = "h1:7ccXrupWZIS3twbUGrtKmHS2DXY6xegFua+6O3xgAFU=",
    version = "v1.0.0",
)

go_repository(
    name = "org_mongodb_go_mongo_driver",
    importpath = "go.mongodb.org/mongo-driver",
    sum = "h1:aeOqSrhl9eDRAap/3T5pCfMBEBxZ0vuXBP+RMtp2KX8=",
    version = "v1.1.0",
)

go_repository(
    name = "org_uber_go_atomic",
    importpath = "go.uber.org/atomic",
    sum = "h1:cxzIVoETapQEqDhQu3QfnvXAV4AlzcvUCxkVUFw3+EU=",
    version = "v1.4.0",
)

go_repository(
    name = "org_uber_go_multierr",
    importpath = "go.uber.org/multierr",
    sum = "h1:HoEmRHQPVSqub6w2z2d2EOVs2fjyFRGyofhKuyDq0QI=",
    version = "v1.1.0",
)

go_repository(
    name = "org_uber_go_zap",
    importpath = "go.uber.org/zap",
    sum = "h1:ORx85nbTijNz8ljznvCMR1ZBIPKFn3jQrag10X2AsuM=",
    version = "v1.10.0",
)

go_repository(
    name = "sh_helm_helm_v3",
    importpath = "helm.sh/helm/v3",
    sum = "h1:gEs30kweCOnLFK9Diq2S8b+VHmWQ2oi465GhqTc3ZxI=",
    version = "v3.0.1",
)

go_repository(
    name = "tools_gotest",
    importpath = "gotest.tools",
    sum = "h1:VsBPFP1AI068pPrMxtb/S8Zkgf9xEmTLJjfM+P5UIEo=",
    version = "v2.2.0+incompatible",
)

go_repository(
    name = "tools_gotest_gotestsum",
    importpath = "gotest.tools/gotestsum",
    sum = "h1:VePOWRsuWFYpfp/G8mbmOZKxO5T3501SEGQRUdvq7h0=",
    version = "v0.3.5",
)

go_repository(
    name = "xyz_gomodules_jsonpatch_v2",
    importpath = "gomodules.xyz/jsonpatch/v2",
    sum = "h1:xyiBuvkD2g5n7cYzx6u2sxQvsAy4QJsZFCzGVdzOXZ0=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_juniper_asf",
    importpath = "github.com/Juniper/asf",
    commit = "ac2649e96024ebed11853a01c83ffd7fc6548919",
)

go_repository(
    name = "com_github_yudai_gotty",
    importpath = "github.com/yudai/gotty",
    sum = "h1:eUFSuV4B2g+Rj+PS3HxhvOGEu2klWRzsl/7z7T/NUJQ=",
    version = "v2.0.0-alpha.3+incompatible",
)

go_repository(
    name = "com_github_flosch_pongo2",
    importpath = "github.com/flosch/pongo2",
    sum = "h1:GY1+t5Dr9OKADM64SYnQjw/w99HMYvQ0A8/JoUkxVmc=",
    version = "v0.0.0-20190707114632-bbf5a6c351f4",
)

go_repository(
    name = "com_github_juju_errors",
    importpath = "github.com/juju/errors",
    sum = "h1:rhqTjzJlm7EbkELJDKMTU7udov+Se0xZkWmugr6zGok=",
    version = "v0.0.0-20181118221551-089d3ea4e4d5",
)

go_repository(
    name = "com_github_volatiletech_sqlboiler",
    importpath = "github.com/volatiletech/sqlboiler",
    sum = "h1:n160O7UQLpZVRnJY6VH5eRNkt7sQdQBZGCCZ3CUy1+g=",
    version = "v3.5.0+incompatible",
)

go_repository(
    name = "com_github_volatiletech_inflect",
    importpath = "github.com/volatiletech/inflect",
    sum = "h1:gI4/tqP6lCY5k6Sg+4k9qSoBXmPwG+xXgMpK7jivD4M=",
    version = "v0.0.0-20170731032912-e7201282ae8d",
)

go_repository(
    name = "com_github_hashicorp_go_cleanhttp",
    importpath = "github.com/hashicorp/go-cleanhttp",
    sum = "h1:dH3aiDG9Jvb5r5+bYHsikaOUIpcM0xvgMXVoDkXMzJM=",
    version = "v0.5.1",
)

go_repository(
    name = "com_github_labstack_echo",
    importpath = "github.com/labstack/echo",
    sum = "h1:pGRcYk231ExFAyoAjAfD85kQzRJCRI8bbnE7CX5OEgg=",
    version = "v3.3.10+incompatible",
)

go_repository(
    name = "com_github_databus23_keystone",
    importpath = "github.com/databus23/keystone",
    sum = "h1:OptdAs3t90tBs6w+lAJVVhBQj3/gqHh1tAQQBL5r08M=",
    version = "v0.0.0-20180111110916-350fd0e663cd",
)

go_repository(
    name = "com_github_labstack_gommon",
    importpath = "github.com/labstack/gommon",
    sum = "h1:JEeO0bvc78PKdyHxloTKiF8BD5iGrH8T6MSeGvSgob0=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:SPIRibHv4MatM3XXNO2BJeFLZwZ2LvZgfQ5+UNI2im4=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_valyala_fasttemplate",
    importpath = "github.com/valyala/fasttemplate",
    sum = "h1:RZqt0yGBsps8NGvLSGW804QQqCUYYLsaOjTVHy1Ocw4=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_valyala_bytebufferpool",
    importpath = "github.com/valyala/bytebufferpool",
    sum = "h1:GqA5TC/0021Y/b9FG4Oi9Mr3q7XYx6KllzawFIhcdPw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_abdullin_seq",
    importpath = "github.com/abdullin/seq",
    sum = "h1:DBNMBMuMiWYu0b+8KMJuWmfCkcxl09JwdlqwDZZ6U14=",
    version = "v0.0.0-20160510034733-d5467c17e7af",
)

go_repository(
    name = "com_github_agext_levenshtein",
    importpath = "github.com/agext/levenshtein",
    sum = "h1:0S/Yg6LYmFJ5stwQeRp6EeOcCbj7xiqQSdNelsXvaqE=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_agl_ed25519",
    importpath = "github.com/agl/ed25519",
    sum = "h1:LoeFxdq5zUCBQPhbQKE6zvoGwHMxCBlqwbH9+9kHoHA=",
    version = "v0.0.0-20150830182803-278e1ec8e8a6",
)

go_repository(
    name = "com_github_aliyun_alibaba_cloud_sdk_go",
    importpath = "github.com/aliyun/alibaba-cloud-sdk-go",
    sum = "h1:APorzFpCcv6wtD5vmRWYqNm4N55kbepL7c7kTq9XI6A=",
    version = "v0.0.0-20190329064014-6e358769c32a",
)

go_repository(
    name = "com_github_aliyun_aliyun_oss_go_sdk",
    importpath = "github.com/aliyun/aliyun-oss-go-sdk",
    sum = "h1:FrF4uxA24DF3ARNXVbUin3wa5fDLaB1Cy8mKks/LRz4=",
    version = "v0.0.0-20190103054945-8205d1f41e70",
)

go_repository(
    name = "com_github_aliyun_aliyun_tablestore_go_sdk",
    importpath = "github.com/aliyun/aliyun-tablestore-go-sdk",
    sum = "h1:ABQ7FF+IxSFHDMOTtjCfmMDMHiCq6EsAoCV/9sFinaM=",
    version = "v4.1.2+incompatible",
)

go_repository(
    name = "com_github_antchfx_xpath",
    importpath = "github.com/antchfx/xpath",
    sum = "h1:ptBAamGVd6CfRsUtyHD+goy2JGhv1QC32v3gqM8mYAM=",
    version = "v0.0.0-20190129040759-c8489ed3251e",
)

go_repository(
    name = "com_github_antchfx_xquery",
    importpath = "github.com/antchfx/xquery",
    sum = "h1:JaCC8jz0zdMLk2m+qCCVLLLM/PL93p84w4pK3aJWj60=",
    version = "v0.0.0-20180515051857-ad5b8c7a47b0",
)

go_repository(
    name = "com_github_aokoli_goutils",
    importpath = "github.com/aokoli/goutils",
    sum = "h1:7fpzNGoJ3VA8qcrm++XEE1QUe0mIwNeLa02Nwq7RDkg=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_apmckinlay_gsuneido",
    importpath = "github.com/apmckinlay/gsuneido",
    sum = "h1:WwxMm9boNuaj5YW+qfRoORxLLJrSRiK1zovCfGNddY0=",
    version = "v0.0.0-20190404155041-0b6cd442a18f",
)

go_repository(
    name = "com_github_apparentlymart_go_cidr",
    importpath = "github.com/apparentlymart/go-cidr",
    sum = "h1:NmIwLZ/KdsjIUlhf+/Np40atNXm/+lZ5txfTJ/SpF+U=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_apparentlymart_go_dump",
    importpath = "github.com/apparentlymart/go-dump",
    sum = "h1:MzVXffFUye+ZcSR6opIgz9Co7WcDx6ZcY+RjfFHoA0I=",
    version = "v0.0.0-20190214190832-042adf3cf4a0",
)

go_repository(
    name = "com_github_apparentlymart_go_textseg",
    importpath = "github.com/apparentlymart/go-textseg",
    sum = "h1:rRmlIsPEEhUTIKQb7T++Nz/A5Q6C9IuX2wFoYVvnCs0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_apparentlymart_go_versions",
    importpath = "github.com/apparentlymart/go-versions",
    sum = "h1:19Seu/H5gq3Ugtx+CGenwF89SDG3S1REX5i6PJj3RK4=",
    version = "v0.0.2-0.20180815153302-64b99f7cb171",
)

go_repository(
    name = "com_github_armon_go_metrics",
    importpath = "github.com/armon/go-metrics",
    sum = "h1:8GUt8eRujhVEGZFFEjBj46YV4rDjvGrNxb0KMWYkL2I=",
    version = "v0.0.0-20180917152333-f0300d1749da",
)

go_repository(
    name = "com_github_armon_go_radix",
    importpath = "github.com/armon/go-radix",
    sum = "h1:F4z6KzEeeQIMeLFa97iZU6vupzoecKdU5TX24SNppXI=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_azure_cli",
    importpath = "github.com/Azure/go-autorest/autorest/azure/cli",
    sum = "h1:pSwNMF0qotgehbQNllUWwJ4V3vnrLKOzHrwDLEZK904=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_azure_go_ntlmssp",
    importpath = "github.com/Azure/go-ntlmssp",
    sum = "h1:pSm8mp0T2OH2CPmPDPtwHPr3VAQaOwVF/JbllOPP4xA=",
    version = "v0.0.0-20180810175552-4a21cbd618b4",
)

go_repository(
    name = "com_github_baiyubin_aliyun_sts_go_sdk",
    importpath = "github.com/baiyubin/aliyun-sts-go-sdk",
    sum = "h1:ZNv7On9kyUzm7fvRZumSyy/IUiSC7AzL0I1jKKtwooA=",
    version = "v0.0.0-20180326062324-cfa1a18b161f",
)

go_repository(
    name = "com_github_bazelbuild_rules_docker",
    importpath = "github.com/bazelbuild/rules_docker",
    sum = "h1:m78lpCWpxwQlLbg/SBKAt9M18aM01KoZGkl8ad0wvLw=",
    version = "v0.14.1",
)

go_repository(
    name = "com_github_bazelbuild_rules_go",
    importpath = "github.com/bazelbuild/rules_go",
    sum = "h1:wzbawlkLtl2ze9w/312NHZ84c7kpUCtlkD8HgFY27sw=",
    version = "v0.0.0-20190719190356-6dae44dc5cab",
)

go_repository(
    name = "com_github_bgentry_go_netrc",
    importpath = "github.com/bgentry/go-netrc",
    sum = "h1:xDfNPAt8lFiC1UJrqV3uuy861HCTo708pDMbjHHdCas=",
    version = "v0.0.0-20140422174119-9fd32a8b3d3d",
)

go_repository(
    name = "com_github_bgentry_speakeasy",
    importpath = "github.com/bgentry/speakeasy",
    sum = "h1:ByYyxL9InA1OWqxJqqp2A5pYHUrCiAL6K3J+LKSsQkY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_bmatcuk_doublestar",
    importpath = "github.com/bmatcuk/doublestar",
    sum = "h1:2bNwBOmhyFEFcoB3tGvTD5xanq+4kyOZlB8wFYbMjkk=",
    version = "v1.1.5",
)

go_repository(
    name = "com_github_census_instrumentation_opencensus_proto",
    importpath = "github.com/census-instrumentation/opencensus-proto",
    sum = "h1:glEXhBS5PSLLv4IXzLA5yPRVX4bilULVyxxbrfOtDAk=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_cheggaaa_pb",
    importpath = "github.com/cheggaaa/pb",
    sum = "h1:wIkZHkNfC7R6GI5w7l/PdAdzXzlrbcI3p8OAlnkTsnc=",
    version = "v1.0.27",
)

go_repository(
    name = "com_github_christrenkamp_goxpath",
    importpath = "github.com/ChrisTrenkamp/goxpath",
    sum = "h1:y8Gs8CzNfDF5AZvjr+5UyGQvQEBL7pwo+v+wX6q9JI8=",
    version = "v0.0.0-20170922090931-c385f95c6022",
)

go_repository(
    name = "com_github_chzyer_logex",
    importpath = "github.com/chzyer/logex",
    sum = "h1:Swpa1K6QvQznwJRcfTfQJmTE72DqScAa40E+fbHEXEE=",
    version = "v1.1.10",
)

go_repository(
    name = "com_github_chzyer_readline",
    importpath = "github.com/chzyer/readline",
    sum = "h1:fY5BOSpyZCqRo5OhCuC+XN+r/bBCmeuuJtjz+bCNIf8=",
    version = "v0.0.0-20180603132655-2972be24d48e",
)

go_repository(
    name = "com_github_chzyer_test",
    importpath = "github.com/chzyer/test",
    sum = "h1:q763qf9huN11kDQavWsoZXJNW3xEE4JJyHa5Q25/sd8=",
    version = "v0.0.0-20180213035817-a1ea475d72b1",
)

go_repository(
    name = "com_github_codegangsta_cli",
    importpath = "github.com/codegangsta/cli",
    sum = "h1:iX1FXEgwzd5+XN6wk5cVHOGQj6Q3Dcp20lUeS4lHNTw=",
    version = "v1.20.0",
)

go_repository(
    name = "com_github_containerd_fifo",
    importpath = "github.com/containerd/fifo",
    sum = "h1:lsjC5ENBl+Zgf38+B0ymougXFp0BaubeIVETltYZTQw=",
    version = "v0.0.0-20191213151349-ff969a566b00",
)

go_repository(
    name = "com_github_containerd_ttrpc",
    importpath = "github.com/containerd/ttrpc",
    sum = "h1:NY8Zk2i7TpkLxrkOASo+KTFq9iNCEmMH2/ZG9OuOw6k=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_dimchansky_utfbom",
    importpath = "github.com/dimchansky/utfbom",
    sum = "h1:FcM3g+nofKgUteL8dm/UpdRXNC9KmADgTpLKsu0TRo4=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_docker_go_events",
    importpath = "github.com/docker/go-events",
    sum = "h1:+pKlWGMw7gf6bQ+oDZB4KHQFypsfjYlq/C4rfL7D3g8=",
    version = "v0.0.0-20190806004212-e31b211e4f1c",
)

go_repository(
    name = "com_github_dylanmei_iso8601",
    importpath = "github.com/dylanmei/iso8601",
    sum = "h1:812NGQDBcqquTfH5Yeo7lwR0nzx/cKdsmf3qMjPURUI=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_dylanmei_winrmtest",
    importpath = "github.com/dylanmei/winrmtest",
    sum = "h1:r1oACdS2XYiAWcfF8BJXkoU8l1J71KehGR+d99yWEDA=",
    version = "v0.0.0-20190225150635-99b7fe2fddf1",
)

go_repository(
    name = "com_github_elazarl_go_bindata_assetfs",
    importpath = "github.com/elazarl/go-bindata-assetfs",
    sum = "h1:G/bYguwHIzWq9ZoyUQqrjTmJbbYn3j3CKKpKinvZLFk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_envoyproxy_go_control_plane",
    importpath = "github.com/envoyproxy/go-control-plane",
    sum = "h1:4cmBvAEBNJaGARUEs3/suWRyfyBfhf7I60WBZq+bv2w=",
    version = "v0.9.1-0.20191026205805-5f8ba28d4473",
)

go_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    importpath = "github.com/envoyproxy/protoc-gen-validate",
    sum = "h1:EQciDnbrYxy13PgWoY8AqoxGiPrpgBZ1R8UNe3ddc+A=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_ericlagergren_decimal",
    importpath = "github.com/ericlagergren/decimal",
    sum = "h1:mMVotm9OVwoOS2IFGRRS5AfMTFWhtf8wj34JEYh47/k=",
    version = "v0.0.0-20191206042408-88212e6cfca9",
)

go_repository(
    name = "com_github_expansiveworlds_instrumentedsql",
    importpath = "github.com/ExpansiveWorlds/instrumentedsql",
    sum = "h1:r+whow+VHd9kAd4UTtQ/rtvcvmwkdryKUcGofpYOp+8=",
    version = "v0.0.0-20171218214018-45abb4b1947d",
)

go_repository(
    name = "com_github_fatih_structs",
    importpath = "github.com/fatih/structs",
    sum = "h1:Q7juDM0QtcnhCpeyLGQKyg4TOIghuNXrkL32pHAUMxo=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_friendsofgo_errors",
    importpath = "github.com/friendsofgo/errors",
    sum = "h1:X6NYxef4efCBdwI7BgS820zFaN7Cphrmb+Pljdzjtgk=",
    version = "v0.9.2",
)

go_repository(
    name = "com_github_go_check_check",
    importpath = "github.com/go-check/check",
    sum = "h1:0gkP6mzaMqkmpcJYCFOLkIBwI7xFExG03bbkOkCvUPI=",
    version = "v0.0.0-20180628173108-788fd7840127",
)

go_repository(
    name = "com_github_go_test_deep",
    importpath = "github.com/go-test/deep",
    sum = "h1:ZrJSEWsXzPOxaZnFteGEfooLba+ju3FYIbOrS+rQd68=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_gogo_googleapis",
    importpath = "github.com/gogo/googleapis",
    sum = "h1:kX1es4djPJrsDhY7aZKJy7aZasdcB5oSOEphMjSB53c=",
    version = "v1.3.2",
)

go_repository(
    name = "com_github_golang_dep",
    importpath = "github.com/golang/dep",
    sum = "h1:WfV5qbGwsBNUDhk+pfI6emWm7SdDFsnSWkqCMNG3BRs=",
    version = "v0.5.4",
)

go_repository(
    name = "com_github_google_go_containerregistry",
    importpath = "github.com/google/go-containerregistry",
    sum = "h1:dRNlOGHybvWTNJFJgF/3t2kTTqVF9Ve4EXqtHzivZd8=",
    version = "v0.0.0-20200320200342-35f57d7d4930",
)

go_repository(
    name = "com_github_gophercloud_utils",
    importpath = "github.com/gophercloud/utils",
    sum = "h1:OgCNGSnEalfkRpn//WGJHhpo7fkP+LhTpvEITZ7CkK4=",
    version = "v0.0.0-20190128072930-fbb6ab446f01",
)

go_repository(
    name = "com_github_hashicorp_aws_sdk_go_base",
    importpath = "github.com/hashicorp/aws-sdk-go-base",
    sum = "h1:zH9hNUdsS+2G0zJaU85ul8D59BGnZBaKM+KMNPAHGwk=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_hashicorp_consul",
    importpath = "github.com/hashicorp/consul",
    sum = "h1:1eDpXAxTh0iPv+1kc9/gfSI2pxRERDsTk/lNGolwHn8=",
    version = "v0.0.0-20171026175957-610f3c86a089",
)

go_repository(
    name = "com_github_hashicorp_go_azure_helpers",
    importpath = "github.com/hashicorp/go-azure-helpers",
    sum = "h1:KhjDnQhCqEMKlt4yH00MCevJQPJ6LkHFdSveXINO6vE=",
    version = "v0.10.0",
)

go_repository(
    name = "com_github_hashicorp_go_checkpoint",
    importpath = "github.com/hashicorp/go-checkpoint",
    sum = "h1:MFYpPZCnQqQTE18jFwSII6eUQrD/oxMFp3mlgcqk5mU=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_hashicorp_go_getter",
    importpath = "github.com/hashicorp/go-getter",
    sum = "h1:l1KB3bHVdvegcIf5upQ5mjcHjs2qsWnKh4Yr9xgIuu8=",
    version = "v1.4.2-0.20200106182914-9813cbd4eb02",
)

go_repository(
    name = "com_github_hashicorp_go_hclog",
    importpath = "github.com/hashicorp/go-hclog",
    sum = "h1:Yv9YzBlAETjy6AOX9eLBZ3nshNVRREgerT/3nvxlGho=",
    version = "v0.0.0-20181001195459-61d530d6c27f",
)

go_repository(
    name = "com_github_hashicorp_go_immutable_radix",
    importpath = "github.com/hashicorp/go-immutable-radix",
    sum = "h1:0nA8i+6Rwqaq9xlpmVxxTwk6rxiEhX+E6Wh4vPNHiS8=",
    version = "v0.0.0-20180129170900-7f3cd4390caa",
)

go_repository(
    name = "com_github_hashicorp_go_msgpack",
    importpath = "github.com/hashicorp/go-msgpack",
    sum = "h1:SFT72YqIkOcLdWJUYcriVX7hbrZpwc/f7h8aW2NUqrA=",
    version = "v0.5.4",
)

go_repository(
    name = "com_github_hashicorp_go_plugin",
    importpath = "github.com/hashicorp/go-plugin",
    sum = "h1:hRho44SAoNu1CBtn5r8Q9J3rCs4ZverWZ4R+UeeNuWM=",
    version = "v1.0.1-0.20190610192547-a1bc61569a26",
)

go_repository(
    name = "com_github_hashicorp_go_retryablehttp",
    importpath = "github.com/hashicorp/go-retryablehttp",
    sum = "h1:AoISa4P4IsW0/m4T6St8Yw38gTl5GtBAgfkhYh1xAz4=",
    version = "v0.5.2",
)

go_repository(
    name = "com_github_hashicorp_go_rootcerts",
    importpath = "github.com/hashicorp/go-rootcerts",
    sum = "h1:Rqb66Oo1X/eSV1x66xbDccZjhJigjg0+e82kpwzSwCI=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_safetemp",
    importpath = "github.com/hashicorp/go-safetemp",
    sum = "h1:2HR189eFNrjHQyENnQMMpCiBAsRxzbTMIgBhEyExpmo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_slug",
    importpath = "github.com/hashicorp/go-slug",
    sum = "h1:/jAo8dNuLgSImoLXaX7Od7QB4TfYCVPam+OpAt5bZqc=",
    version = "v0.4.1",
)

go_repository(
    name = "com_github_hashicorp_go_sockaddr",
    importpath = "github.com/hashicorp/go-sockaddr",
    sum = "h1:7YOlAIO2YWnJZkQp7B5eFykaIY7C9JndqAFQyVV5BhM=",
    version = "v0.0.0-20180320115054-6d291a969b86",
)

go_repository(
    name = "com_github_hashicorp_go_tfe",
    importpath = "github.com/hashicorp/go-tfe",
    sum = "h1:7XZ/ZoPyYoeuNXaWWW0mJOq016y0qb7I4Q0P/cagyu8=",
    version = "v0.3.27",
)

go_repository(
    name = "com_github_hashicorp_go_uuid",
    importpath = "github.com/hashicorp/go-uuid",
    sum = "h1:fv1ep09latC32wFoVwnqcnKJGnMSdBanPczbHAYm1BE=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_hashicorp_hcl_v2",
    importpath = "github.com/hashicorp/hcl/v2",
    sum = "h1:iRly8YaMwTBAKhn1Ybk7VSdzbnopghktCD031P8ggUE=",
    version = "v2.3.0",
)

go_repository(
    name = "com_github_hashicorp_hil",
    importpath = "github.com/hashicorp/hil",
    sum = "h1:2yzhWGdgQUWZUCNK+AoO35V+HTsgEmcM4J9IkArh7PI=",
    version = "v0.0.0-20190212112733-ab17b08d6590",
)

go_repository(
    name = "com_github_hashicorp_memberlist",
    importpath = "github.com/hashicorp/memberlist",
    sum = "h1:qSsCiC0WYD39lbSitKNt40e30uorm2Ss/d4JGU1hzH8=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_hashicorp_serf",
    importpath = "github.com/hashicorp/serf",
    sum = "h1:ZbgmOQt8DOg796figP87/EFCVx2v2h9yRvwHF/zceX4=",
    version = "v0.0.0-20160124182025-e4ec8cc423bb",
)

go_repository(
    name = "com_github_hashicorp_terraform",
    importpath = "github.com/hashicorp/terraform",
    sum = "h1:lTTswsCcmTOhTwuUl2NdjtJBCNdGqZmRGQi0cjFHYOM=",
    version = "v0.12.24",
)

go_repository(
    name = "com_github_hashicorp_terraform_config_inspect",
    importpath = "github.com/hashicorp/terraform-config-inspect",
    sum = "h1:Pc5TCv9mbxFN6UVX0LH6CpQrdTM5YjbVI2w15237Pjk=",
    version = "v0.0.0-20191212124732-c6ae6269b9d7",
)

go_repository(
    name = "com_github_hashicorp_terraform_svchost",
    importpath = "github.com/hashicorp/terraform-svchost",
    sum = "h1:hjyO2JsNZUKT1ym+FAdlBEkGPevazYsmVgIMw7dVELg=",
    version = "v0.0.0-20191011084731-65d371908596",
)

go_repository(
    name = "com_github_hashicorp_vault",
    importpath = "github.com/hashicorp/vault",
    sum = "h1:4x0lHxui/ZRp/B3E0Auv1QNBJpzETqHR2kQD3mHSBJU=",
    version = "v0.10.4",
)

go_repository(
    name = "com_github_hashicorp_yamux",
    importpath = "github.com/hashicorp/yamux",
    sum = "h1:b5rjCoWHc7eqmAS4/qyk21ZsHyb6Mxv/jykxvNTkU4M=",
    version = "v0.0.0-20180604194846-3520598351bb",
)

go_repository(
    name = "com_github_jmank88_nuts",
    importpath = "github.com/jmank88/nuts",
    sum = "h1:3rHp+7YcvtkTPohGBA++MwneB9OlX/rpORvleiRivMQ=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_joyent_triton_go",
    importpath = "github.com/joyent/triton-go",
    sum = "h1:kie3qOosvRKqwij2HGzXWffwpXvcqfPPXRUw8I4F/mg=",
    version = "v0.0.0-20180313100802-d8f9c0314926",
)

go_repository(
    name = "com_github_juju_loggo",
    importpath = "github.com/juju/loggo",
    sum = "h1:MK144iBQF9hTSwBW/9eJm034bVoG30IshVm688T2hi8=",
    version = "v0.0.0-20180524022052-584905176618",
)

go_repository(
    name = "com_github_juju_testing",
    importpath = "github.com/juju/testing",
    sum = "h1:WQM1NildKThwdP7qWrNAFGzp4ijNLw8RlgENkaI4MJs=",
    version = "v0.0.0-20180920084828-472a3e8b2073",
)

go_repository(
    name = "com_github_juniper_contrail",
    importpath = "github.com/Juniper/contrail",
    sum = "h1:31S3FvAcU12o5sLmsAvNjG0p2XAF7lDI6ceFxnpj7K0=",
    version = "v0.0.0-20200330181744-e78e7561c8fd",
)

go_repository(
    name = "com_github_kat_co_vala",
    importpath = "github.com/kat-co/vala",
    sum = "h1:DQVOxR9qdYEybJUr/c7ku34r3PfajaMYXZwgDM7KuSk=",
    version = "v0.0.0-20170210184112-42e1d8b61f12",
)

go_repository(
    name = "com_github_keybase_go_crypto",
    importpath = "github.com/keybase/go-crypto",
    sum = "h1:NARVGAAgEXvoMeNPHhPFt1SBt1VMznA3Gnz9d0qj+co=",
    version = "v0.0.0-20161004153544-93f5b35093ba",
)

go_repository(
    name = "com_github_kr_fs",
    importpath = "github.com/kr/fs",
    sum = "h1:Jskdu9ieNAYnjxsi0LbQp1ulIKZV1LAFgK1tWhpZgl8=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_labstack_echo_v4",
    importpath = "github.com/labstack/echo/v4",
    sum = "h1:8swiwjE5Jkai3RPfZoahp8kjVCRNq+y7Q0hPji2Kz0o=",
    version = "v4.1.16",
)

go_repository(
    name = "com_github_likexian_gokit",
    importpath = "github.com/likexian/gokit",
    sum = "h1:DgtIqqTRFqtbiLJFzuRESwVrxWxfs8OlY6hnPYBa3BM=",
    version = "v0.20.15",
)

go_repository(
    name = "com_github_likexian_simplejson_go",
    importpath = "github.com/likexian/simplejson-go",
    sum = "h1:tNa5q0zLsg6EzUI4npXeUYwk/xuTtMcKdVnRsYdOHzk=",
    version = "v0.0.0-20190502021454-d8787b4bfa0b",
)

go_repository(
    name = "com_github_lusis_go_artifactory",
    importpath = "github.com/lusis/go-artifactory",
    sum = "h1:wnfcqULT+N2seWf6y4yHzmi7GD2kNx4Ute0qArktD48=",
    version = "v0.0.0-20160115162124-7e4ce345df82",
)

go_repository(
    name = "com_github_masterminds_sprig",
    importpath = "github.com/Masterminds/sprig",
    sum = "h1:0gSxPGWS9PAr7U2NsQ2YQg6juRDINkUyuvbb4b2Xm8w=",
    version = "v2.15.0+incompatible",
)

go_repository(
    name = "com_github_masterzen_simplexml",
    importpath = "github.com/masterzen/simplexml",
    sum = "h1:SmVbOZFWAlyQshuMfOkiAx1f5oUTsOGG5IXplAEYeeM=",
    version = "v0.0.0-20160608183007-4572e39b1ab9",
)

go_repository(
    name = "com_github_masterzen_winrm",
    importpath = "github.com/masterzen/winrm",
    sum = "h1:/1RFh2SLCJ+tEnT73+Fh5R2AO89sQqs8ba7o+hx1G0Y=",
    version = "v0.0.0-20190223112901-5e5c9a7fe54b",
)

go_repository(
    name = "com_github_mattn_goveralls",
    importpath = "github.com/mattn/goveralls",
    sum = "h1:7eJB6EqsPhRVxvwEXGnqdO2sJI0PTsrWoTMXEk9/OQc=",
    version = "v0.0.2",
)

go_repository(
    name = "com_github_mitchellh_cli",
    importpath = "github.com/mitchellh/cli",
    sum = "h1:iGBIsUe3+HZ/AD/Vd7DErOt5sU9fa8Uj7A2s1aggv1Y=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_colorstring",
    importpath = "github.com/mitchellh/colorstring",
    sum = "h1:62I3jR2EmQ4l5rM/4FEfDWcRD+abF5XlKShorW5LRoQ=",
    version = "v0.0.0-20190213212951-d06e56a500db",
)

go_repository(
    name = "com_github_mitchellh_go_linereader",
    importpath = "github.com/mitchellh/go-linereader",
    sum = "h1:GRiLv4rgyqjqzxbhJke65IYUf4NCOOvrPOJbV/sPxkM=",
    version = "v0.0.0-20190213213312-1b945b3263eb",
)

go_repository(
    name = "com_github_mitchellh_go_testing_interface",
    importpath = "github.com/mitchellh/go-testing-interface",
    sum = "h1:fzU/JVNcaqHQEcVFAKeR41fkiLdIPrefOvVG1VZ96U0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_iochan",
    importpath = "github.com/mitchellh/iochan",
    sum = "h1:C+X3KsSTLFVBr/tK1eYN/vs4rJcvsiLU338UhYPJWeY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_panicwrap",
    importpath = "github.com/mitchellh/panicwrap",
    sum = "h1:67zIyVakCIvcs69A0FGfZjBdPleaonSgGlXRSRlb6fE=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_prefixedio",
    importpath = "github.com/mitchellh/prefixedio",
    sum = "h1:eD92Am0Qf3rqhsOeA1zwBHSfRkoHrt4o6uORamdmJP8=",
    version = "v0.0.0-20190213213902-5733675afd51",
)

go_repository(
    name = "com_github_mozillazg_go_httpheader",
    importpath = "github.com/mozillazg/go-httpheader",
    sum = "h1:geV7TrjbL8KXSyvghnFm+NyTux/hxwueTSrwhe88TQQ=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_mwitkow_go_proto_validators",
    importpath = "github.com/mwitkow/go-proto-validators",
    sum = "h1:28i1IjGcx8AofiB4N3q5Yls55VEaitzuEPkFJEVgGkA=",
    version = "v0.0.0-20180403085117-0950a7990007",
)

go_repository(
    name = "com_github_nightlyone_lockfile",
    importpath = "github.com/nightlyone/lockfile",
    sum = "h1:RHep2cFKK4PonZJDdEl4GmkabuhbsRMgk/k3uAmxBiA=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_nu7hatch_gouuid",
    importpath = "github.com/nu7hatch/gouuid",
    sum = "h1:VhgPp6v9qf9Agr/56bj7Y/xa04UccTW04VP0Qed4vnQ=",
    version = "v0.0.0-20131221200532-179d4d0c4d8d",
)

go_repository(
    name = "com_github_packer_community_winrmcp",
    importpath = "github.com/packer-community/winrmcp",
    sum = "h1:m3CEgv3ah1Rhy82L+c0QG/U3VyY1UsvsIdkh0/rU97Y=",
    version = "v0.0.0-20180102160824-81144009af58",
)

go_repository(
    name = "com_github_pascaldekloe_goe",
    importpath = "github.com/pascaldekloe/goe",
    sum = "h1:Lgl0gzECD8GnQ5QCWA8o6BtfL6mDH5rQgM4/fX3avOs=",
    version = "v0.0.0-20180627143212-57f6aae5913c",
)

go_repository(
    name = "com_github_pkg_browser",
    importpath = "github.com/pkg/browser",
    sum = "h1:49lOXmGaUpV9Fz3gd7TFZY106KVlPVa5jcYD1gaQf98=",
    version = "v0.0.0-20180916011732-0a3d74bf9ce4",
)

go_repository(
    name = "com_github_pkg_sftp",
    importpath = "github.com/pkg/sftp",
    sum = "h1:4Zv0OGbpkg4yNuUtH0s8rvoYxRCNyT29NVUo6pgPmxI=",
    version = "v1.11.0",
)

go_repository(
    name = "com_github_pmylund_go_cache",
    importpath = "github.com/pmylund/go-cache",
    sum = "h1:n+7K51jLz6a3sCvff3BppuCAkixuDHuJ/C57Vw/XjTE=",
    version = "v2.1.0+incompatible",
)

go_repository(
    name = "com_github_posener_complete",
    importpath = "github.com/posener/complete",
    sum = "h1:LrvDIY//XNo65Lq84G/akBuMGlawHvGBABv8f/ZN6DI=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_pseudomuto_protoc_gen_doc",
    importpath = "github.com/pseudomuto/protoc-gen-doc",
    sum = "h1:Segz6bKr2LCo9bZgm+foCYiyfr4s0BurLzH3MDE7wC0=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_pseudomuto_protokit",
    importpath = "github.com/pseudomuto/protokit",
    sum = "h1:hlnBDcy3YEDXH7kc9gV+NLaN0cDzhDvD1s7Y6FZ8RpM=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_qcloudapi_qcloud_sign_golang",
    importpath = "github.com/QcloudApi/qcloud_sign_golang",
    sum = "h1:DTQ/38ao/CfXsrK0cSAL+h4R/u0VVvfWLZEOlLwEROI=",
    version = "v0.0.0-20141224014652-e4130a326409",
)

go_repository(
    name = "com_github_sdboyer_constext",
    importpath = "github.com/sdboyer/constext",
    sum = "h1:tnWWLf0nI2TI62Wd/ZOea4XYqE+y1sf2pdm+VItsc0c=",
    version = "v0.0.0-20170321163424-836a14457353",
)

go_repository(
    name = "com_github_sean_seed",
    importpath = "github.com/sean-/seed",
    sum = "h1:nn5Wsu0esKSJiIVhscUtVbo7ada43DJhG55ua/hjS5I=",
    version = "v0.0.0-20170313163322-e2103e2c3529",
)

go_repository(
    name = "com_github_shurcool_sanitized_anchor_name",
    importpath = "github.com/shurcooL/sanitized_anchor_name",
    sum = "h1:lrOUn8raSZS/V52c7elGaEyuogqSkEo/Qj2Auo2G1ik=",
    version = "v0.0.0-20151028001915-10ef21a441db",
)

go_repository(
    name = "com_github_sigma_go_inotify",
    importpath = "github.com/sigma/go-inotify",
    sum = "h1:G1nNtZVTzcCvVKMwcG0Vispo3bhc15EbjO5uamiLikI=",
    version = "v0.0.0-20181102212354-c87b6cf5033d",
)

go_repository(
    name = "com_github_streadway_amqp",
    importpath = "github.com/streadway/amqp",
    sum = "h1:2MR0pKUzlP3SGgj5NYJe/zRYDwOu9ku6YHy+Iw7l5DM=",
    version = "v0.0.0-20200108173154-1c71cc93ed71",
)

go_repository(
    name = "com_github_svanharmelen_jsonapi",
    importpath = "github.com/svanharmelen/jsonapi",
    sum = "h1:Z4EH+5EffvBEhh37F0C0DnpklTMh00JOkjW5zK3ofBI=",
    version = "v0.0.0-20180618144545-0c0828c3f16d",
)

go_repository(
    name = "com_github_tencentcloud_tencentcloud_sdk_go",
    importpath = "github.com/tencentcloud/tencentcloud-sdk-go",
    sum = "h1:5Td2b0yfaOvw9M9nZ5Oav6Li9bxUNxt4DgxMfIPpsa0=",
    version = "v3.0.82+incompatible",
)

go_repository(
    name = "com_github_tencentyun_cos_go_sdk_v5",
    importpath = "github.com/tencentyun/cos-go-sdk-v5",
    sum = "h1:iRD1CqtWUjgEVEmjwTMbP1DMzz1HRytOsgx/rlw/vNs=",
    version = "v0.0.0-20190808065407-f07404cefc8c",
)

go_repository(
    name = "com_github_terraform_providers_terraform_provider_openstack",
    importpath = "github.com/terraform-providers/terraform-provider-openstack",
    sum = "h1:adpjqej+F8BAX9dHmuPF47sUIkgifeqBu6p7iCsyj0Y=",
    version = "v1.15.0",
)

go_repository(
    name = "com_github_ulikunitz_xz",
    importpath = "github.com/ulikunitz/xz",
    sum = "h1:pFrO0lVpTBXLpYw+pnLj6TbvHuyjXMfjGeCwSqCVwok=",
    version = "v0.5.5",
)

go_repository(
    name = "com_github_unknwon_com",
    importpath = "github.com/Unknwon/com",
    sum = "h1:tuQ7w+my8a8mkwN7x2TSd7OzTjkZ7rAeSyH4xncuAMI=",
    version = "v0.0.0-20151008135407-28b053d5a292",
)

go_repository(
    name = "com_github_vdemeester_k8s_pkg_credentialprovider",
    importpath = "github.com/vdemeester/k8s-pkg-credentialprovider",
    sum = "h1:czKEIG2Q3YRTgs6x/8xhjVMJD5byPo6cZuostkbTM74=",
    version = "v1.17.4",
)

go_repository(
    name = "com_github_vmihailenco_msgpack",
    importpath = "github.com/vmihailenco/msgpack",
    sum = "h1:RMF1enSPeKTlXrXdOcqjFUElywVZjjC6pqse21bKbEU=",
    version = "v4.0.1+incompatible",
)

go_repository(
    name = "com_github_volatiletech_null",
    importpath = "github.com/volatiletech/null",
    sum = "h1:7wP8m5d/gZ6kW/9GnrLtMCRre2dlEnaQ9Km5OXlK4zg=",
    version = "v8.0.0+incompatible",
)

go_repository(
    name = "com_github_xanzy_ssh_agent",
    importpath = "github.com/xanzy/ssh-agent",
    sum = "h1:TCbipTQL2JiiCprBWx9frJ2eJlCYT00NmctrHxVAr70=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_xlab_treeprint",
    importpath = "github.com/xlab/treeprint",
    sum = "h1:Jpn2j6wHkC9wJv5iMfJhKqrZJx3TahFx+7sbZ7zQdxs=",
    version = "v0.0.0-20161029104018-1d6e34225557",
)

go_repository(
    name = "com_github_yudai_hcl",
    importpath = "github.com/yudai/hcl",
    sum = "h1:tjsK9T2IA3d2FFNxzDP7AJf+EXhyuPd7PB4Z2HrtAoc=",
    version = "v0.0.0-20151013225006-5fa2393b3552",
)

go_repository(
    name = "com_github_zclconf_go_cty",
    importpath = "github.com/zclconf/go-cty",
    sum = "h1:vGMsygfmeCl4Xb6OA5U5XVAaQZ69FvoG7X2jUtQujb8=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_zclconf_go_cty_yaml",
    importpath = "github.com/zclconf/go-cty-yaml",
    sum = "h1:up11wlgAaDvlAGENcFDnZgkn0qUJurso7k6EpURKNF8=",
    version = "v1.0.1",
)

go_repository(
    name = "com_google_cloud_go_bigquery",
    importpath = "cloud.google.com/go/bigquery",
    sum = "h1:hL+ycaJpVE9M7nLoiXb/Pn10ENE2u+oddxbD8uu0ZVU=",
    version = "v1.0.1",
)

go_repository(
    name = "com_google_cloud_go_datastore",
    importpath = "cloud.google.com/go/datastore",
    sum = "h1:Kt+gOPPp2LEPWp8CSfxhsM8ik9CcyE/gYu+0r+RnZvM=",
    version = "v1.0.0",
)

go_repository(
    name = "in_gopkg_cheggaaa_pb_v1",
    importpath = "gopkg.in/cheggaaa/pb.v1",
    sum = "h1:kJdccidYzt3CaHD1crCFTS1hxyhSi059NhOFUf03YFo=",
    version = "v1.0.27",
)

go_repository(
    name = "in_gopkg_data_dog_go_sqlmock_v1",
    importpath = "gopkg.in/DATA-DOG/go-sqlmock.v1",
    sum = "h1:FVCohIoYO7IJoDDVpV2pdq7SgrMH6wHnuTyrdrxJNoY=",
    version = "v1.3.0",
)

go_repository(
    name = "in_gopkg_ini_v1",
    importpath = "gopkg.in/ini.v1",
    sum = "h1:7N3gPTt50s8GuLortA00n8AqRTk75qOP98+mTPpgzRk=",
    version = "v1.42.0",
)

go_repository(
    name = "in_gopkg_mgo_v2",
    importpath = "gopkg.in/mgo.v2",
    sum = "h1:xcEWjVhvbDy+nHP67nPDDpbYrY+ILlfndk4bRioVHaU=",
    version = "v2.0.0-20180705113604-9856a29383ce",
)

go_repository(
    name = "in_gopkg_volatiletech_null_v6",
    importpath = "gopkg.in/volatiletech/null.v6",
    sum = "h1:P+3+n9hUbqSDkSdtusWHVPQRrpRpLiLFzlZ02xXskM0=",
    version = "v6.0.0-20170828023728-0bef4e07ae1b",
)

go_repository(
    name = "io_rsc_binaryregexp",
    importpath = "rsc.io/binaryregexp",
    sum = "h1:HfqmD5MEmC0zvwBuF187nq9mdnXjXsSivRiXN7SmRkE=",
    version = "v0.2.0",
)

go_repository(
    name = "net_starlark_go",
    importpath = "go.starlark.net",
    sum = "h1:WaYQLsW/cWywevmveTnnHnGGnumFCJ4E9nrPsCd0N9c=",
    version = "v0.0.0-20200326215636-e8819e807894",
)

go_repository(
    name = "com_github_go_pg_pg_v10",
    importpath = "github.com/go-pg/pg/v10",
    sum = "h1:8tNEJLtOEw5/Df0BLLBOHCiLaYAiu4uhdngjK955MK8=",
    version = "v10.0.0-beta.2",
)