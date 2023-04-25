module github.com/konvera/geth-sev

go 1.20

replace (
	github.com/google/go-tpm => github.com/thomasten/go-tpm v0.0.0-20230222180349-bb3cc5560299
	github.com/google/go-tpm-tools => github.com/daniel-weisse/go-tpm-tools v0.0.0-20230105122812-f7474d459dfc
)

require (
	cloud.google.com/go/compute v1.18.0
	cloud.google.com/go/compute/metadata v0.2.3
	cloud.google.com/go/kms v1.10.0
	cloud.google.com/go/logging v1.7.0
	cloud.google.com/go/storage v1.30.1
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.2.2
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2 v2.1.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.6.1
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/aws/aws-sdk-go-v2 v1.17.7
	github.com/aws/aws-sdk-go-v2/config v1.18.19
	github.com/aws/aws-sdk-go-v2/credentials v1.13.18
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.1
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.59
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.26.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.20.7
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.91.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.19.7
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.14.7
	github.com/aws/aws-sdk-go-v2/service/s3 v1.31.0
	github.com/aws/smithy-go v1.13.5
	github.com/edgelesssys/constellation/v2 v2.7.0
	github.com/edgelesssys/go-azguestattestation v0.0.0-20230303085714-62ede861d33f
	github.com/fsnotify/fsnotify v1.6.0
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.11.2
	github.com/google/go-tpm v0.3.3
	github.com/google/go-tpm-tools v0.3.10
	github.com/googleapis/gax-go/v2 v2.7.1
	github.com/gophercloud/gophercloud v1.2.0
	github.com/gophercloud/utils v0.0.0-20230324070755-05e9e7f5ea4d
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/go-kms-wrapping/v2 v2.0.8
	github.com/hashicorp/go-kms-wrapping/wrappers/awskms/v2 v2.0.7
	github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault/v2 v2.0.7
	github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms/v2 v2.0.8
	github.com/microsoft/ApplicationInsights-Go v0.4.4
	github.com/siderolabs/talos/pkg/machinery v1.3.2
	github.com/sigstore/rekor v1.0.1
	github.com/sigstore/sigstore v1.6.0
	github.com/spf13/afero v1.9.5
	github.com/spf13/cobra v1.6.1
	github.com/stretchr/testify v1.8.2
	go.uber.org/goleak v1.2.1
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.6.0
	golang.org/x/mod v0.8.0
	golang.org/x/tools v0.6.0
	google.golang.org/api v0.114.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.29.1
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.26.3
	k8s.io/apiextensions-apiserver v0.26.3
	k8s.io/apimachinery v0.26.3
	k8s.io/client-go v0.26.3
	k8s.io/utils v0.0.0-20230313181309-38a27ef9d749
)

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/iam v0.12.0 // indirect
	cloud.google.com/go/longrunning v0.4.1 // indirect
	code.cloudfoundry.org/clock v0.0.0-20180518195852-02e53af36e6c // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.2.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.21 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.11 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.6 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v0.9.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.44.209 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.7 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cyberphone/json-canonicalization v0.0.0-20210303052042-6bc126869bf4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/go-chi/chi v4.1.2+incompatible // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.1 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.24.1 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/certificate-transparency-go v1.1.4 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-attestation v0.4.4-0.20221011162210-17f9c05652a9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-containerregistry v0.13.0 // indirect
	github.com/google/go-sev-guest v0.4.1 // indirect
	github.com/google/go-tspi v0.3.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/logger v1.1.1 // indirect
	github.com/google/trillian v1.5.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.2 // indirect
	github.com/hashicorp/go-secure-stdlib/awsutil v0.1.6 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jedisct1/go-minisign v0.0.0-20211028175153-1c139d1cc84b // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.2 // indirect
	github.com/letsencrypt/boulder v0.0.0-20221109233200-85aa52084eaf // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sassoftware/relic v0.0.0-20210427151427-dfb082b79b74 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.5.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tent/canonical-json-go v0.0.0-20130607151641-96e4ba3a7613 // indirect
	github.com/theupdateframework/go-tuf v0.5.2 // indirect
	github.com/titanous/rocacheck v0.0.0-20171023193734-afe73141d399 // indirect
	github.com/transparency-dev/merkle v0.0.1 // indirect
	go.mongodb.org/mongo-driver v1.10.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20220823124025-807a23277127 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230320184635-7606e756e683 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.90.0 // indirect
	k8s.io/kube-openapi v0.0.0-20230228151317-19cbebb19cb7 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
