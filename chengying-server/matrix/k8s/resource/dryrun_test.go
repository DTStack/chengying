// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"dtstack.com/dtstack/easymatrix/addons/easykube/pkg/client/base"
	commonlog "dtstack.com/dtstack/easymatrix/go-common/log"
	"dtstack.com/dtstack/easymatrix/matrix/api/k8s/view"
	clustergenerator "dtstack.com/dtstack/easymatrix/matrix/k8s/cluster"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/constant"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/kube"
	"dtstack.com/dtstack/easymatrix/matrix/k8s/resource/deployment"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	modelkube "dtstack.com/dtstack/easymatrix/matrix/model/kube"
	"dtstack.com/dtstack/easymatrix/matrix/model/kube/union"
	"github.com/jmoiron/sqlx"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var kubeconfig = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZWRENDQXp3Q0NRQ3ZTQlZyUVQ1ZGpEQU5CZ2txaGtpRzl3MEJBUXNGQURCc01Rc3dDUVlEVlFRR0V3SkQKVGpFUk1BOEdBMVVFQ0F3SVdtaGxhbWxoYm1jeEVUQVBCZ05WQkFjTUNFaGhibWQ2YUc5MU1SQXdEZ1lEVlFRSwpEQWQ0WkdWMmIzQnpNUkF3RGdZRFZRUUxEQWQ0WkdWMmIzQnpNUk13RVFZRFZRUUREQXBrZEhOMFlXTnJMbU51Ck1CNFhEVEl3TVRFd01USXhNVGd5TmxvWERUTXdNVEF6TURJeE1UZ3lObG93YkRFTE1Ba0dBMVVFQmhNQ1EwNHgKRVRBUEJnTlZCQWdNQ0Zwb1pXcHBZVzVuTVJFd0R3WURWUVFIREFoSVlXNW5lbWh2ZFRFUU1BNEdBMVVFQ2d3SAplR1JsZG05d2N6RVFNQTRHQTFVRUN3d0hlR1JsZG05d2N6RVRNQkVHQTFVRUF3d0taSFJ6ZEdGamF5NWpiakNDCkFpSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnSVBBRENDQWdvQ2dnSUJBTHh3MWgzcXptTXE4YUErY01vWFQvM0UKRVVBdzRQNUhCZGZYa2VVa1J1ZjV3VHUvdTBvRHFyUnVRNlcwMFJWS09BUFFQVDZpRVhCQkdyKzdTRGhISi8yTAo4VzBnNDRnaDB5N2lWREE1UmErV2dFS2Y4YmRvT3grUkMzME4wb3JyRmtZbDRCR3p2UWtDQXlOV1lZU1lheGxDCkFEeDdoQldhMzFQb0FaT3MxN0pBM3Y0a1lUVzhzTVBHQzFtdnhCSHJ4cGswK3NxeFFLQjl2TzlaZ05zbXpQcXAKdUNsT0xQbkxjdkFrSjFMSXNIS3c4bXo3TlFnVnpVV3hzWEFTMVFnaTdXWlNwRUVINGsySmtQRTRLNU5Mc2g2Mgo0d3hPS0lRRS83TFFsaVlCRlpIZnFYTHZpU0F5MW94SklTbnp4UmdmOFZVSmUrMUxKbGE0UHp6UDZuVUI2cGxUCnlKODQ2Z0QxQStxTmM1endlN2RtMytHSFVJNVVRNXhZejNZRDdQclRIWlJJd3p3SzVzeXNTM28rTEFOL1duYjgKSnlVK2haNG5hcEZUUzh0SVlLWGQ5V0pWUldYNWNzdkd2MXZSY1daK0tySGRtVi9Mc1oxbG04RFd1ODFLV0EyZQpKSmdvMlMrZ3R0RFBEenZqMHVnbUNEVFVRZGpsYjVQWnhFdWI4cVYvUXViMTBCQ3JJSkNZenJHMm91dld6TVJLCmJIUEE0OWtIVVYrR3BSdW8vNDZZYUw2T3haVjRQcnJjWFU1MVlQTVNZRnJiWU1iT2g1b1NxTnkyY1lrdjRhTkUKMDNzQUllclJHMzYvaktnelYzNFFkZGJvTEI5Yk1sZThEaGJQZ3JEV2VGcTdJRlpBVlYwOUcwdGx6R1hNaDE5YQpLZm8yVnNMRDZMT1dLTGwvOUNVRkFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUNvM3ZYcHBJTk9TCkxkSEpTRXRYT3haSi9kMEw4aXZQaEUwZzlRb0h0NUNNUUZWMGlKcVpzZzVLQURDckJSUG42ZnBVZDFhaDROZWoKWW5lUWxzbEpUOFFJN1RIQ1ZiTzBuOUtsbEVxWEJ5c2gxeDdkTGVXdFVaczRxLzB2QVlxbTc3aWYyaVpWTEFzNgprMk81NmhpelNwN2svZC90SXd3a1k3dTN4UkF2b2RWWEM3cnA1ZmViSTRWcDZLUFMwbjF4Yk8wSXBqUzd3KzhPClkxMjFXdmtxUUQyckhCMkUySFJSSWdGTEF2eVFqY01hak9KMzRSUnh4VDljUzFIWEUzVUdDVVl2bHc1NnMySmMKMWYvMXNXYTJNaHFiUUdzVDZJS3dVOXZ2SlArQktoWnBsWnpicUlnZHZyZHVHNWFCbEpnNlNCQTZsUThSTDhnNQp6MmpscGFaTE5jUmYySkdTUm5TTzlDZzJTcU5va2wwaDRRZVRCZEFNeGhPVnk0MmxOc1grLzVOWVhuM0RXRWo1CmoybTV0SFpRU3ptY3BZZmNlY0ZpRXdud3BwOVdMbWlsb0N3SHVON2FNYU9ZakpSVnVPcWJoTXZMUjFLTnpWVm0KUFdVblZJa3FBY2p0TXVKVFh6clZGV081VXEvazF0WkIzdGhNNEd4c1JZcElDbEJmS3c3emFIaGxsZmNvNE53eQprbWVDeENWMnVLZ3AxYkJlZ0hReTVEc3E4OFI5Qkh1b0g0aWpEbzlUTGRpU2NXT0o1TkFEOUJ5K2F3VGpsUk1qCmtrVHgyMFRZdm1VdEF6T2N5R0xCOENEWERwS1B4b1ZBdVREMGFKMy9na3p0aHZQQ1dDamxINGI2cXVMS3lwVHgKOU85OVRZWGw3Y2VNSXZIc2I3L3l4L01KUG9RbUZjcUgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==
    server: https://172.16.8.88/k8s/clusters/c-vckqj
  name: dtstack
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQWFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKTFdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBMVVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4NzhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTSG84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLcXdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1dzdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmTllGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqeGZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNNmRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtVDlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQy9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hcDdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRSlVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZEVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://172.16.8.160:6443
  name: dtstack-172-16-8-160
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQWFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKTFdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBMVVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4NzhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTSG84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLcXdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1dzdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmTllGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqeGZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNNmRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtVDlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQy9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hcDdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRSlVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZEVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://172.16.8.166:6443
  name: dtstack-172-16-8-166
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQWFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKTFdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBMVVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4NzhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTSG84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLcXdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1dzdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmTllGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqeGZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNNmRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtVDlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQy9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hcDdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRSlVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZEVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://172.16.8.169:6443
  name: dtstack-172-16-8-169
contexts:
- context:
    cluster: dtstack
    user: dtstack
  name: dtstack
- context:
    cluster: dtstack-172-16-8-160
    user: dtstack
  name: dtstack-172-16-8-160
- context:
    cluster: dtstack-172-16-8-166
    user: dtstack
  name: dtstack-172-16-8-166
- context:
    cluster: dtstack-172-16-8-169
    user: dtstack
  name: dtstack-172-16-8-169
- context:
    cluster: dtstack-172-16-8-160
    user: test
  name: test
current-context: dtstack
kind: Config
preferences: {}
users:
- name: dtstack
  user:
    token: kubeconfig-user-mtrk8.c-vckqj:jk8k6n72dpxtllmtd4p99kwkkxthxsqg6562b9r2t2v7bdkf2jft9n
- name: test
  user:
    token: eyJhbGciOiJSUzI1NiIsImtpZCI6IjNhTjhlRFZNQl9ld2xGRDhpQ1F3czJDamtHUGNMTVlua2JBVW0yYUhDT1EifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZXYiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiZGVmYXVsdC10b2tlbi1kajl0NiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiOGM4ZTVhZGQtOGYwZi00NDdkLTg3YjAtMTc1ZmQxYTJmYWQ1Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRldjpkZWZhdWx0In0.bxqYVW0HOklhEZxNXECMQy1MqmdGzW4D2MLJ6fnXeAAXolkYPVIgf69eOO0P9JDxtX0HALQz9P982CqnuP0U7M3_OkJUptr55ixQXnkYbPn8L18ZxGdQ0R8tqymoqdfZAmUzzNf1lk9BIRaX_DAJwJmwKWUdrSGdJaS2lPcsXILP0GcC0tbwK6PP3GL2ZKKKlQsDy8Hbw2DdgXkMnx6NyU2BlsW2CWfEDc0__XUfY-0TSA4Evy7DnrD-US06992BF63zirhp8Y_kNd-8VQY83oNCo9vrKwOJYU0s0UP_3BgZpaKUSyrn4alKoj69hz4sXnnARwhbgJxKF_qBYqy47Q`
var testkube = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQWFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKTFdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBMVVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4NzhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTSG84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLcXdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1dzdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmTllGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqeGZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNNmRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtVDlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQy9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hcDdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRSlVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZEVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://172.16.8.160:6443
  name: dtstack-172-16-8-160
contexts:
- context:
    cluster: dtstack-172-16-8-160
    user: test
  name: test
current-context: test
kind: Config
preferences: {}
users:
- name: test
  user:
    token: eyJhbGciOiJSUzI1NiIsImtpZCI6IjNhTjhlRFZNQl9ld2xGRDhpQ1F3czJDamtHUGNMTVlua2JBVW0yYUhDT1EifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJ0ZXN0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tajhxczIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjYzYjBiMmY4LTZiNzYtNGVjZi05NjZlLTY2Y2NiMjQwZTFhYyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDp0ZXN0OmRlZmF1bHQifQ.FE6HdoTYuCGzwyVJUKLBZXxAubQ88srgmk5gDwe4EUMxT_bzb9QSP7BQFK1sPpNhDXQSRFJd_ACVWKmhssycnFnbkWj7VmaicqrjJfh8-tcyX1MvykMbVGSFY3Czm-RoTFuB47ybbs7DHckwcLZDtlQYQjBORK1H91uWBSR_Knal4MulO7suDu4G_WgmT6OYgZjOmh_g81lxkV3GMNHElxSvDkLOq9bp4Ze949tKhC4uSTL7XR4o23mV8xXdAltOzKc9FPnE1wR5qXzGPhZezXzqhW-T38-hswWp1B3iweK1nFJnT__4jTZW40piLuv0jeIvZgSjsesQLSPKCAgOeA
`

func Test1(t *testing.T) {
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	ns := "test"
	cache, err := kube.ClusterNsClientCache.GetClusterNsClient("abc").GetClientCache(kube.IMPORT_KUBECONFIG)
	if err != nil {
		fmt.Println("err1", err.Error())
		return
	}
	err = cache.Connect(testkube, ns)
	if err != nil {
		fmt.Println("err2", err.Error())
		return
	}
	err = deployment.Ping(cache.GetClient(ns), ns)
	if err != nil {
		fmt.Println("err3", err.Error())
	}
}

func Test2(t *testing.T) {
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	ctx := context.Background()
	vo := &view.NamespaceSaveReq{
		Type:       "kubeconfig",
		Namespace:  "em-dev",
		RegistryId: 20,
		Yaml:       "apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZWRENDQXp3Q0NRQ3ZTQlZyUVQ1ZGpEQU5CZ2txaGtpRzl3MEJBUXNGQURCc01Rc3dDUVlEVlFRR0V3SkQKVGpFUk1BOEdBMVVFQ0F3SVdtaGxhbWxoYm1jeEVUQVBCZ05WQkFjTUNFaGhibWQ2YUc5MU1SQXdEZ1lEVlFRSwpEQWQ0WkdWMmIzQnpNUkF3RGdZRFZRUUxEQWQ0WkdWMmIzQnpNUk13RVFZRFZRUUREQXBrZEhOMFlXTnJMbU51Ck1CNFhEVEl3TVRFd01USXhNVGd5TmxvWERUTXdNVEF6TURJeE1UZ3lObG93YkRFTE1Ba0dBMVVFQmhNQ1EwNHgKRVRBUEJnTlZCQWdNQ0Zwb1pXcHBZVzVuTVJFd0R3WURWUVFIREFoSVlXNW5lbWh2ZFRFUU1BNEdBMVVFQ2d3SAplR1JsZG05d2N6RVFNQTRHQTFVRUN3d0hlR1JsZG05d2N6RVRNQkVHQTFVRUF3d0taSFJ6ZEdGamF5NWpiakNDCkFpSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnSVBBRENDQWdvQ2dnSUJBTHh3MWgzcXptTXE4YUErY01vWFQvM0UKRVVBdzRQNUhCZGZYa2VVa1J1ZjV3VHUvdTBvRHFyUnVRNlcwMFJWS09BUFFQVDZpRVhCQkdyKzdTRGhISi8yTAo4VzBnNDRnaDB5N2lWREE1UmErV2dFS2Y4YmRvT3grUkMzME4wb3JyRmtZbDRCR3p2UWtDQXlOV1lZU1lheGxDCkFEeDdoQldhMzFQb0FaT3MxN0pBM3Y0a1lUVzhzTVBHQzFtdnhCSHJ4cGswK3NxeFFLQjl2TzlaZ05zbXpQcXAKdUNsT0xQbkxjdkFrSjFMSXNIS3c4bXo3TlFnVnpVV3hzWEFTMVFnaTdXWlNwRUVINGsySmtQRTRLNU5Mc2g2Mgo0d3hPS0lRRS83TFFsaVlCRlpIZnFYTHZpU0F5MW94SklTbnp4UmdmOFZVSmUrMUxKbGE0UHp6UDZuVUI2cGxUCnlKODQ2Z0QxQStxTmM1endlN2RtMytHSFVJNVVRNXhZejNZRDdQclRIWlJJd3p3SzVzeXNTM28rTEFOL1duYjgKSnlVK2haNG5hcEZUUzh0SVlLWGQ5V0pWUldYNWNzdkd2MXZSY1daK0tySGRtVi9Mc1oxbG04RFd1ODFLV0EyZQpKSmdvMlMrZ3R0RFBEenZqMHVnbUNEVFVRZGpsYjVQWnhFdWI4cVYvUXViMTBCQ3JJSkNZenJHMm91dld6TVJLCmJIUEE0OWtIVVYrR3BSdW8vNDZZYUw2T3haVjRQcnJjWFU1MVlQTVNZRnJiWU1iT2g1b1NxTnkyY1lrdjRhTkUKMDNzQUllclJHMzYvaktnelYzNFFkZGJvTEI5Yk1sZThEaGJQZ3JEV2VGcTdJRlpBVlYwOUcwdGx6R1hNaDE5YQpLZm8yVnNMRDZMT1dLTGwvOUNVRkFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUNvM3ZYcHBJTk9TCkxkSEpTRXRYT3haSi9kMEw4aXZQaEUwZzlRb0h0NUNNUUZWMGlKcVpzZzVLQURDckJSUG42ZnBVZDFhaDROZWoKWW5lUWxzbEpUOFFJN1RIQ1ZiTzBuOUtsbEVxWEJ5c2gxeDdkTGVXdFVaczRxLzB2QVlxbTc3aWYyaVpWTEFzNgprMk81NmhpelNwN2svZC90SXd3a1k3dTN4UkF2b2RWWEM3cnA1ZmViSTRWcDZLUFMwbjF4Yk8wSXBqUzd3KzhPClkxMjFXdmtxUUQyckhCMkUySFJSSWdGTEF2eVFqY01hak9KMzRSUnh4VDljUzFIWEUzVUdDVVl2bHc1NnMySmMKMWYvMXNXYTJNaHFiUUdzVDZJS3dVOXZ2SlArQktoWnBsWnpicUlnZHZyZHVHNWFCbEpnNlNCQTZsUThSTDhnNQp6MmpscGFaTE5jUmYySkdTUm5TTzlDZzJTcU5va2wwaDRRZVRCZEFNeGhPVnk0MmxOc1grLzVOWVhuM0RXRWo1CmoybTV0SFpRU3ptY3BZZmNlY0ZpRXdud3BwOVdMbWlsb0N3SHVON2FNYU9ZakpSVnVPcWJoTXZMUjFLTnpWVm0KUFdVblZJa3FBY2p0TXVKVFh6clZGV081VXEvazF0WkIzdGhNNEd4c1JZcElDbEJmS3c3emFIaGxsZmNvNE53eQprbWVDeENWMnVLZ3AxYkJlZ0hReTVEc3E4OFI5Qkh1b0g0aWpEbzlUTGRpU2NXT0o1TkFEOUJ5K2F3VGpsUk1qCmtrVHgyMFRZdm1VdEF6T2N5R0xCOENEWERwS1B4b1ZBdVREMGFKMy9na3p0aHZQQ1dDamxINGI2cXVMS3lwVHgKOU85OVRZWGw3Y2VNSXZIc2I3L3l4L01KUG9RbUZjcUgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==\n    server: https://172.16.8.88/k8s/clusters/c-vckqj\n  name: dtstack\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQWFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKTFdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBMVVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4NzhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTSG84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLcXdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1dzdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmTllGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqeGZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNNmRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtVDlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQy9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hcDdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRSlVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZEVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K\n    server: https://172.16.8.160:6443\n  name: dtstack-172-16-8-160\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQWFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKTFdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBMVVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4NzhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTSG84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLcXdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1dzdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmTllGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqeGZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNNmRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtVDlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQy9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hcDdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRSlVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZEVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K\n    server: https://172.16.8.166:6443\n  name: dtstack-172-16-8-166\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQWFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKTFdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBMVVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4NzhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTSG84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLcXdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1dzdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmTllGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqeGZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNNmRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtVDlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQy9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hcDdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRSlVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZEVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K\n    server: https://172.16.8.169:6443\n  name: dtstack-172-16-8-169\ncontexts:\n- context:\n    cluster: dtstack\n    user: dtstack\n  name: dtstack\n- context:\n    cluster: dtstack-172-16-8-160\n    user: dtstack\n  name: dtstack-172-16-8-160\n- context:\n    cluster: dtstack-172-16-8-166\n    user: dtstack\n  name: dtstack-172-16-8-166\n- context:\n    cluster: dtstack-172-16-8-169\n    user: dtstack\n  name: dtstack-172-16-8-169\n- context:\n    cluster: dtstack-172-16-8-160\n    user: test\n  name: test\ncurrent-context: dtstack\nkind: Config\npreferences: {}\nusers:\n- name: dtstack\n  user:\n    token: kubeconfig-user-mtrk8.c-vckqj:jk8k6n72dpxtllmtd4p99kwkkxthxsqg6562b9r2t2v7bdkf2jft9n\n- name: test\n  user:\n    token: eyJhbGciOiJSUzI1NiIsImtpZCI6IjNhTjhlRFZNQl9ld2xGRDhpQ1F3czJDamtHUGNMTVlua2JBVW0yYUhDT1EifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZXYiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiZGVmYXVsdC10b2tlbi1kajl0NiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiOGM4ZTVhZGQtOGYwZi00NDdkLTg3YjAtMTc1ZmQxYTJmYWQ1Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRldjpkZWZhdWx0In0.bxqYVW0HOklhEZxNXECMQy1MqmdGzW4D2MLJ6fnXeAAXolkYPVIgf69eOO0P9JDxtX0HALQz9P982CqnuP0U7M3_OkJUptr55ixQXnkYbPn8L18ZxGdQ0R8tqymoqdfZAmUzzNf1lk9BIRaX_DAJwJmwKWUdrSGdJaS2lPcsXILP0GcC0tbwK6PP3GL2ZKKKlQsDy8Hbw2DdgXkMnx6NyU2BlsW2CWfEDc0__XUfY-0TSA4Evy7DnrD-US06992BF63zirhp8Y_kNd-8VQY83oNCo9vrKwOJYU0s0UP_3BgZpaKUSyrn4alKoj69hz4sXnnARwhbgJxKF_qBYqy47Q",
	}
	clusterid := "109"
	user := "admin@dtstack.com"
	Save(ctx, clusterid, user, vo)
}

func Test3(t *testing.T) {
	user := "root"
	password := "dtstack"
	host := "172.16.10.37"
	port := 3306
	dbname := "dtagent"
	db, _ := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true", user, password, host, port, dbname))
	model.MYSQLDB = db
	err := modelkube.Build()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	ctx := context.Background()
	err = InitResource()
	if err != nil {
		fmt.Println("err3", err.Error())
	}
	rsp, err := GetNamespaceListStatus(ctx, "85", "", "", "")
	if err != nil {
		fmt.Println("err", err.Error())
	}
	fmt.Println(rsp)
}

func getConn() {
	user := "root"
	password := "dtstack"
	host := "172.16.10.37"
	port := 3306
	dbname := "dtagent"
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	db, _ := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true", user, password, host, port, dbname))
	model.MYSQLDB = db
	err := modelkube.Build()
	union.Build()
	if err != nil {
		fmt.Printf("err: %v \n", err)
	}
}

func Test4(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	ctx := context.Background()
	ns := "dt-test-2"
	cache, _ := kube.ClusterNsClientCache.GetClusterNsClient("abc").GetClientCache(kube.IMPORT_KUBECONFIG)
	cache.Connect(kubeconfig, ns)
	namespaceInit(ctx, cache, ns, 21)
}

//func Test5(t *testing.T) {
//	getConn()
//	log.ConfigureLogger("/tmp/matrix",100,3,1)
//	commonlog.Config("/tmp/easykube")
//	ns := "dt-test-2"
//	ctx := context.Background()
//	cache,_ := kube.ClusterNsClientCache.GetClusterNsClient("85").GetClientCache(kube.IMPORT_KUBECONFIG)
//	cache.Connect(kubeconfig,ns)
//	resp,err := GetNamespaceStatus(ctx,ns,"85")
//	if err != nil{
//		fmt.Println("err: ",err.Error())
//	}
//	bts,_:=yaml.Marshal(resp)
//	fmt.Println(string(bts))
//}

func Test6(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	ns := "em-dev"
	clusterid := "111"
	ctx := context.Background()
	cache, _ := kube.ClusterNsClientCache.GetClusterNsClient(clusterid).GetClientCache(kube.IMPORT_KUBECONFIG)
	cache.Connect(kubeconfig, ns)
	err := NamespaceDelete(ctx, ns, clusterid)
	if err != nil {
		fmt.Println("err", err.Error())
	}
}

func Test7(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	clusterid := "111"
	ctx := context.Background()
	res, err := GetNamespaceListStatus(ctx, clusterid, "", "true", "")
	if err != nil {
		fmt.Println("err", err.Error())
	}
	bts, _ := json.Marshal(res)
	fmt.Println(string(bts))
}

func Test8(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	c := model.ClusterInfo{
		Id:   111,
		Name: "em-dev",
		Type: "kubernets",
		Mode: 1,
	}
	rsp, err := NamespaceList(&c)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(rsp)
}

func Test9(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	req := &view.AgentGenerateReq{
		Namespace:  "dev",
		RegistryId: 22,
	}
	rsp, err := AgentGenerate(req)
	if err != nil {
		fmt.Println("err", err.Error())
	}
	fmt.Println(rsp.Yaml)
}

func Test10(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	//namespace := "em-dev"
	clusterid := "112"
	status := ""
	descStr := "true"
	ctx := context.Background()
	//rsp,err := GetNamespaceStatus(ctx,namespace,clusterid,status,descStr)
	rsp, err := GetNamespaceListStatus(ctx, clusterid, status, descStr, "")
	if err != nil {
		fmt.Println("err", err.Error())
	}
	bts, _ := json.Marshal(rsp)
	fmt.Println(string(bts))
}

func Test11(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	req := &view.AgentGenerateReq{
		Namespace: "",
	}
	rsp, err := AgentGenerate(req)
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	bts, _ := json.Marshal(rsp)
	fmt.Println(string(bts))
}

func Test12(t *testing.T) {
	getConn()
	np := &view.NamespacePingReq{
		Namespace: "em-agent",
		Ip:        "172.16.100.117",
		Port:      "80",
	}
	err := NamespacePing(context.Background(), "113", np)
	if err != nil {
		fmt.Println("err", err.Error())
	}
}

func Test13(t *testing.T) {
	var s runtime.Object = &corev1.ResourceQuota{}
	gvks, _, _ := base.Schema.ObjectKinds(s)
	gvk := gvks[0]
	fmt.Println(gvk.Kind)
	fmt.Println(gvk.Version)
	fmt.Println(gvk.Group)
}

func Test14(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")

	ns := "em-dev"
	cache, _ := kube.ClusterNsClientCache.GetClusterNsClient("112").GetClientCache(kube.IMPORT_KUBECONFIG)
	cache.Connect(kubeconfig, ns)

	rsp, err := GetService(context.Background(), "em-dev", "112", "DTinsight", "all", "sql")
	if err != nil {
		fmt.Println("err", err.Error())
	}
	bts, _ := json.Marshal(rsp)
	fmt.Println(string(bts))
}

func Test15(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	cache, _ := kube.ClusterNsClientCache.GetClusterNsClient("113").GetClientCache(kube.IMPORT_AGENT)
	cache.Connect("http://172.16.100.117:80", "em-agent")
	res, err := GetNamespaceListStatus(context.Background(), "113", "", "true", "")
	if err != nil {
		fmt.Println("err", err.Error())
	}
	bts, _ := json.Marshal(res)
	fmt.Println(string(bts))
}

func Test16(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	cache, _ := kube.ClusterNsClientCache.GetClusterNsClient("113").GetClientCache(kube.IMPORT_AGENT)
	cache.Connect("http://172.16.100.117:80", "em-agent")
	cc, _ := kube.ClusterNsClientCache.GetClusterNsClient("113").GetClientCache(kube.IMPORT_AGENT)

	c := cc.GetClient("em-agent")
	cct := c.(*kube.RestClient)
	monitorAgent("em-agent", 113, cct)
	//cct.Events(&e)
	//fmt.Println(cct == nil)
	//events := []monitorevents.Event{}
	//err := c.Events(&events)
	//if err != nil{}
	//fmt.Println(err.Error())
	//bts,_:=json.Marshal(e)
	//fmt.Println(string(bts))
}

//
func Test17(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	cache, _ := kube.ClusterNsClientCache.GetClusterNsClient("113").GetClientCache(kube.IMPORT_AGENT)
	cache.Connect("http://172.16.100.117:80", "em-agent")
	req := &view.InstanceReplicaReq{
		Namespace:   "em-agent",
		Replica:     1,
		ProductName: "DTUic",
		ServiceName: "Uic",
	}
	err := InstanceReplica(context.Background(), "113", req)
	if err != nil {
		fmt.Println("err", err.Error())
	}
}

func Test19(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	cache, _ := kube.ClusterNsClientCache.GetClusterNsClient("135").GetClientCache(kube.IMPORT_AGENT)
	cache.Connect("http://172.16.10.234:8899", "dtstack-system")
	ns := &corev1.NamespaceList{}
	err := cache.GetClient("dtstack-system").List(context.Background(), ns, "")
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	bts, _ := json.Marshal(ns)
	fmt.Println(string(bts))
}
func Test18(t *testing.T) {
	getConn()
	log.ConfigureLogger("/tmp/matrix", 100, 3, 1)
	commonlog.Config("/tmp/easykube")
	cinfo, err := modelkube.DeployClusterList.GetClusterById(131)
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	gInfo := &clustergenerator.GeneratorInfo{
		Type:        constant.TYPE_SELF_BUILD,
		HostIp:      "172.16.10.37",
		ClusterInfo: cinfo,
	}
	//err = clustergenerator.GenerateTemplate(gInfo)
	//if err != nil{
	//	fmt.Println("err",err.Error())
	//}

	bts, _ := clustergenerator.GetTemplateFile(gInfo, false)
	fmt.Println(string(bts))
}
