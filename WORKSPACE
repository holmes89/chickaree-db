load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "69de5c704a05ff37862f7e0f5534d4f479418afc21806c887db544a316f3cb6b",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)

http_archive(
    name = "com_google_protobuf",
    sha256 = "9748c0d90e54ea09e5e75fb7fac16edce15d2028d4356f32211cfa3c0e956564",
    strip_prefix = "protobuf-3.11.4",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.11.4.zip"],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

go_repository(
    name = "com_github_boltdb_bolt",
    importpath = "github.com/boltdb/bolt",
    sum = "h1:JQmyP4ZBrce+ZQu0dY660FMfatumYDLun9hBCUVIkF4=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    sum = "h1:ROPKBNFfQgOUMifHyP+KYbvpjbdoFNs+aK7DXlji0Tw=",
    version = "v1.5.2",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    sum = "h1:Khx7svrCpmxxtHBq5j2mp/xVjsi8hQMfNLvJFAlrGgU=",
    version = "v0.5.5",
)

go_repository(
    name = "org_golang_google_grpc_cmd_protoc_gen_go_grpc",
    importpath = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
    sum = "h1:M1YKkFIboKNieVO5DLUEVzQfGwJD30Nv2jfUgzb5UcE=",
    version = "v1.1.0",
)

go_repository(
    name = "org_golang_google_protobuf",
    importpath = "google.golang.org/protobuf",
    sum = "h1:bxAC2xTBsZGibn2RTntX0oH50xLsqy1OxA9tTL3p/lk=",
    version = "v1.26.0",
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:gG67DSER+11cZvqIMb8S8bt0vZtiN6xWYARwirrOSfE=",
    version = "v0.0.0-20210510120138-977fb7262007",
)

go_repository(
    name = "org_golang_x_xerrors",
    importpath = "golang.org/x/xerrors",
    sum = "h1:go1bK/D/BFZV2I8cIQd1NKEZ+0owSTG1fDTci4IqFcE=",
    version = "v0.0.0-20200804184101-5ec99f83aff1",
)

go_repository(
    name = "com_github_armon_circbuf",
    importpath = "github.com/armon/circbuf",
    sum = "h1:QEF07wC0T1rKkctt1RINW/+RMTVmiwxETico2l3gxJA=",
    version = "v0.0.0-20150827004946-bbbad097214e",
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
    name = "com_github_bgentry_speakeasy",
    importpath = "github.com/bgentry/speakeasy",
    sum = "h1:ByYyxL9InA1OWqxJqqp2A5pYHUrCiAL6K3J+LKSsQkY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:ZDRjVQ15GmhC3fiQ8ni8+OwkZQO4DARzQgrnXU1Liz8=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_fatih_color",
    importpath = "github.com/fatih/color",
    sum = "h1:8xPHl4/q1VyqGIPif1F+1V3Y3lSmrq01EabUW3CoW5s=",
    version = "v1.9.0",
)

go_repository(
    name = "com_github_google_btree",
    importpath = "github.com/google/btree",
    sum = "h1:964Od4U6p2jUkFxvCydnIczKteheJEzHRToSGK3Bnlw=",
    version = "v0.0.0-20180813153112-4030bb1f1f0c",
)

go_repository(
    name = "com_github_hashicorp_errwrap",
    importpath = "github.com/hashicorp/errwrap",
    sum = "h1:hLrqtEDnRye3+sgx6z4qVLNuviH3MR5aQ0ykNJa/UYA=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_immutable_radix",
    importpath = "github.com/hashicorp/go-immutable-radix",
    sum = "h1:AKDB1HM5PWEA7i4nhcpwOrO2byshxBjXVn/J/3+z5/0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_msgpack",
    importpath = "github.com/hashicorp/go-msgpack",
    sum = "h1:zKjpN5BK/P5lMYrLmBHdBULWbJ0XpYR+7NGzqkZzoD4=",
    version = "v0.5.3",
)

go_repository(
    name = "com_github_hashicorp_go_multierror",
    importpath = "github.com/hashicorp/go-multierror",
    sum = "h1:B9UzwGQJehnUY1yNrnwREHc3fGbC2xefo8g4TbElacI=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_hashicorp_go_sockaddr",
    importpath = "github.com/hashicorp/go-sockaddr",
    sum = "h1:GeH6tui99pF4NJgfnhp+L6+FfobzVW3Ah46sLo0ICXs=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_syslog",
    importpath = "github.com/hashicorp/go-syslog",
    sum = "h1:KaodqZuhUoZereWVIYmpUgZysurB1kBLX2j0MwMrUAE=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_uuid",
    importpath = "github.com/hashicorp/go-uuid",
    sum = "h1:fv1ep09latC32wFoVwnqcnKJGnMSdBanPczbHAYm1BE=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    importpath = "github.com/hashicorp/golang-lru",
    sum = "h1:CL2msUPvZTLb5O648aiLNJw3hnBxN2+1Jq8rCOH9wdo=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_hashicorp_logutils",
    importpath = "github.com/hashicorp/logutils",
    sum = "h1:dLEQVugN8vlakKOUE3ihGLTZJRB4j+M2cdTm/ORI65Y=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_mdns",
    importpath = "github.com/hashicorp/mdns",
    sum = "h1:XFSOubp8KWB+Jd2PDyaX5xUd5bhSP/+pTDZVDMzZJM8=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_hashicorp_memberlist",
    importpath = "github.com/hashicorp/memberlist",
    sum = "h1:5+RffWKwqJ71YPu9mWsF7ZOscZmwfasdA8kbdC7AO2g=",
    version = "v0.2.2",
)

go_repository(
    name = "com_github_hashicorp_serf",
    importpath = "github.com/hashicorp/serf",
    sum = "h1:EBWvyu9tcRszt3Bxp3KNssBMP1KuHWyO51lz9+786iM=",
    version = "v0.9.5",
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
    name = "com_github_miekg_dns",
    importpath = "github.com/miekg/dns",
    sum = "h1:gPxPSwALAeHJSjarOs00QjVdV9QoBvc1D2ujQUr5BzU=",
    version = "v1.1.26",
)

go_repository(
    name = "com_github_mitchellh_cli",
    importpath = "github.com/mitchellh/cli",
    sum = "h1:tEElEatulEHDeedTxwckzyYMA5c86fbmNIUL1hBIiTg=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    importpath = "github.com/mitchellh/mapstructure",
    sum = "h1:kK7VuFVykgt0LfMSloWYjDOt4TnOcL0AxF0/rDq2VkM=",
    version = "v0.0.0-20160808181253-ca63d7c062ee",
)

go_repository(
    name = "com_github_pascaldekloe_goe",
    importpath = "github.com/pascaldekloe/goe",
    sum = "h1:Lgl0gzECD8GnQ5QCWA8o6BtfL6mDH5rQgM4/fX3avOs=",
    version = "v0.0.0-20180627143212-57f6aae5913c",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_posener_complete",
    importpath = "github.com/posener/complete",
    sum = "h1:NP0eAhjcjImqslEwo/1hq7gpajME0fTLTezBKDqfXqo=",
    version = "v1.2.3",
)

go_repository(
    name = "com_github_ryanuber_columnize",
    importpath = "github.com/ryanuber/columnize",
    sum = "h1:UFr9zpz4xgTnIE5yIMtWAMngCdZ9p/+q6lTbgelo80M=",
    version = "v0.0.0-20160712163229-9b3edd62028f",
)

go_repository(
    name = "com_github_sean_seed",
    importpath = "github.com/sean-/seed",
    sum = "h1:nn5Wsu0esKSJiIVhscUtVbo7ada43DJhG55ua/hjS5I=",
    version = "v0.0.0-20170313163322-e2103e2c3529",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:nOGnQDM7FYENwehXlg/kFVnos3rEvtKTjRvOWSzb6H4=",
    version = "v1.5.1",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
    version = "v0.0.0-20161208181325-20d25e280405",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:ZCJp+EgiOT7lHqUV2J862kp8Qj64Jo6az82+3Td9dZw=",
    version = "v2.2.2",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:psW17arqaxU48Z5kZ0CQnkZWQJsqcURM6tKiBApRjXI=",
    version = "v0.0.0-20200622213623-75b288015ac9",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:SQFwaSi55rU7vdNs9Yr0Z324VNlrF+0wMqRXT4St8ck=",
    version = "v0.0.0-20201020160332-67f06af15bc9",
)

go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    sum = "h1:po9/4sTYwZU9lPhi1tOrb4hCv3qrhiQ77LZfGa2OjwY=",
    version = "v0.1.0",
)

go_repository(
    name = "co_honnef_go_tools",
    importpath = "honnef.co/go/tools",
    sum = "h1:/hemPrYIhOhy8zYrNj+069zDB68us2sMGsfkFJO0iZs=",
    version = "v0.0.0-20190523083050-ea95bdfd59fc",
)

go_repository(
    name = "com_github_burntsushi_toml",
    importpath = "github.com/BurntSushi/toml",
    sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
    version = "v0.3.1",
)

go_repository(
    name = "com_github_census_instrumentation_opencensus_proto",
    importpath = "github.com/census-instrumentation/opencensus-proto",
    sum = "h1:glEXhBS5PSLLv4IXzLA5yPRVX4bilULVyxxbrfOtDAk=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_client9_misspell",
    importpath = "github.com/client9/misspell",
    sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
    version = "v0.3.4",
)

go_repository(
    name = "com_github_cncf_udpa_go",
    importpath = "github.com/cncf/udpa/go",
    sum = "h1:cqQfy1jclcSy/FwLjemeg3SR1yaINm74aQyupQ0Bl8M=",
    version = "v0.0.0-20201120205902-5459f2c99403",
)

go_repository(
    name = "com_github_envoyproxy_go_control_plane",
    importpath = "github.com/envoyproxy/go-control-plane",
    sum = "h1:QyzYnTnPE15SQyUeqU6qLbWxMkwyAyu+vGksa0b7j00=",
    version = "v0.9.9-0.20210217033140-668b12f5399d",
)

go_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    importpath = "github.com/envoyproxy/protoc-gen-validate",
    sum = "h1:EQciDnbrYxy13PgWoY8AqoxGiPrpgBZ1R8UNe3ddc+A=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_golang_glog",
    importpath = "github.com/golang/glog",
    sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
    version = "v0.0.0-20160126235308-23def4e6c14b",
)

go_repository(
    name = "com_github_golang_mock",
    importpath = "github.com/golang/mock",
    sum = "h1:G5FRp8JnTd7RQH5kemVNlMeyXQAztQ3mOWV95KxsXH8=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_google_uuid",
    importpath = "github.com/google/uuid",
    sum = "h1:EVhdT+1Kseyi1/pUmXKaFxYsDNy9RQYkMWRH68J/W7Y=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_prometheus_client_model",
    importpath = "github.com/prometheus/client_model",
    sum = "h1:gQz4mCbXsO+nc9n1hCxHcGA3Zx3Eo+UHZoInFGUIXNM=",
    version = "v0.0.0-20190812154241-14fe0d1b01d4",
)

go_repository(
    name = "com_github_yuin_goldmark",
    importpath = "github.com/yuin/goldmark",
    sum = "h1:ruQGxdhGHe7FWOJPT0mKs5+pD2Xs1Bm/kdGlHO04FmM=",
    version = "v1.2.1",
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    sum = "h1:e0WKqKTd5BnrG8aKH3J3h+QvEIQtSUcf2n5UZ5ZgLtQ=",
    version = "v0.26.0",
)

go_repository(
    name = "org_golang_google_appengine",
    importpath = "google.golang.org/appengine",
    sum = "h1:/wp5JvzpHIxhs/dumFmF7BXTf3Z+dd4uXta4kVyO508=",
    version = "v1.4.0",
)

go_repository(
    name = "org_golang_google_genproto",
    importpath = "google.golang.org/genproto",
    sum = "h1:tzkHckzMzgPr8SC4taTC3AldLr4+oJivSoq1xf/nhsc=",
    version = "v0.0.0-20210510173355-fb37daa5cd7a",
)

go_repository(
    name = "org_golang_x_exp",
    importpath = "golang.org/x/exp",
    sum = "h1:c2HOrn5iMezYjSlGPncknSEr/8x5LELb/ilJbXi9DEA=",
    version = "v0.0.0-20190121172915-509febef88a4",
)

go_repository(
    name = "org_golang_x_lint",
    importpath = "golang.org/x/lint",
    sum = "h1:2M3HP5CCK1Si9FQhwnzYhXdG6DXeebvUHFpre8QvbyI=",
    version = "v0.0.0-20201208152925-83fdc39ff7b5",
)

go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:RM4zey1++hCTbCVQfnWeKs9/IEsaBLA8vTkd0WVtmH4=",
    version = "v0.3.0",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    sum = "h1:vEDujvNQGv4jgYKudGeI/+DAX4Jffq6hpD55MmoEvKs=",
    version = "v0.0.0-20180821212333-d2e6202438be",
)

go_repository(
    name = "org_golang_x_term",
    importpath = "golang.org/x/term",
    sum = "h1:v+OssWQX+hTHEmOBgwxdZxK4zHq3yOs8F9J7mk0PY8E=",
    version = "v0.0.0-20201126162022-7de9c90e9dd1",
)

go_rules_dependencies()

go_register_toolchains(version = "1.16")

gazelle_dependencies()

go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable",
    importpath = "google.golang.org/grpc",
    sum = "h1:uSZWeQJX5j11bIQ4AJoj+McDBo29cY1MCoC1wO3ts+c=",
    version = "v1.37.0",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:p9UgmWI9wKpfYmgaV/IZKGdXc5qEK45tDwwwDyjS26I=",
    version = "v0.0.0-20210510120150-4163338589ed",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:aRYxNxv6iGQlyVaZmk6ZgYEDa+Jg18DxebPSrd6bg1M=",
    version = "v0.3.6",
)

protobuf_deps()
