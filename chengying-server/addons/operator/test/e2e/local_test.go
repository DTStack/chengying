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

package e2e

import (
	"dtstack.com/dtstack/easymatrix/addons/operator/cmd/options"
	statcklog "dtstack.com/dtstack/easymatrix/addons/operator/log"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/apis"
	"dtstack.com/dtstack/easymatrix/addons/operator/pkg/controller"
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
	"runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"testing"
)

var log = logf.Log.WithName("test-main")

func startPrint() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info("operator start")
}

func TestLocal(t *testing.T) {
	fs := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	opt := options.GetOptions()
	opt.AddToFlagSet(fs)
	fs.Parse(os.Args)

	//set klog log to stderr
	klogFlagSet := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlagSet)

	klogFlagSet.Set("logtostderr", "true")
	klogFlagSet.Set("alsologtostderr", "false")

	logf.SetLogger(statcklog.NewLogger(opt))

	startPrint()

	//namespace, err := getNamespace(opt.WatchNamespace)
	namespace := wathchNamespace
	//if err != nil {
	//	log.Error(err, "Failed to get watch namespace")
	//	os.Exit(1)
	//}
	// use incluster config to talk with api-server
	cfg, err := getConfig()
	if err != nil {
		log.Error(err, "unable to get rest configuration")
	}

	// would elect a leader to start the informer,
	//if the operator pod is not leader any more, the pod will exist and than restart
	mgr, err := manager.New(cfg, manager.Options{
		LeaderElection:          true,
		LeaderElectionNamespace: namespace,
		LeaderElectionID:        opt.ElectionLockName,
		Namespace:               namespace,
	})
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}

func getConfig() (*rest.Config, error) {
	apiconfig, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		log.Error(err, "[kube_client]: load kubeconfig fail")
		return nil, err
	}
	config, err := clientcmd.NewNonInteractiveClientConfig(*apiconfig, "", &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		log.Error(err, "create rest config from kubeconfig fail")
		return nil, err
	}
	return config, nil
}

func getNamespace(namespace string) (string, error) {
	if len(namespace) != 0 {
		return namespace, nil
	}
	ns, found := os.LookupEnv("WATCH_NAMESPACE")
	if !found {
		return "", fmt.Errorf("%s must be set", "WATCH_NAMESPACE")
	}
	return ns, nil
}

var wathchNamespace = "workload-test"

var kubeconfig = `apiVersion: v1
kind: Config
clusters:
- name: "dtstack"
  cluster:
    server: "https://172.16.8.88/k8s/clusters/c-vckqj"
    certificate-authority-data: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZWRENDQ\
      Xp3Q0NRQ3ZTQlZyUVQ1ZGpEQU5CZ2txaGtpRzl3MEJBUXNGQURCc01Rc3dDUVlEVlFRR0V3SkQKV\
      GpFUk1BOEdBMVVFQ0F3SVdtaGxhbWxoYm1jeEVUQVBCZ05WQkFjTUNFaGhibWQ2YUc5MU1SQXdEZ\
      1lEVlFRSwpEQWQ0WkdWMmIzQnpNUkF3RGdZRFZRUUxEQWQ0WkdWMmIzQnpNUk13RVFZRFZRUUREQ\
      XBrZEhOMFlXTnJMbU51Ck1CNFhEVEl3TVRFd01USXhNVGd5TmxvWERUTXdNVEF6TURJeE1UZ3lOb\
      G93YkRFTE1Ba0dBMVVFQmhNQ1EwNHgKRVRBUEJnTlZCQWdNQ0Zwb1pXcHBZVzVuTVJFd0R3WURWU\
      VFIREFoSVlXNW5lbWh2ZFRFUU1BNEdBMVVFQ2d3SAplR1JsZG05d2N6RVFNQTRHQTFVRUN3d0hlR\
      1JsZG05d2N6RVRNQkVHQTFVRUF3d0taSFJ6ZEdGamF5NWpiakNDCkFpSXdEUVlKS29aSWh2Y05BU\
      UVCQlFBRGdnSVBBRENDQWdvQ2dnSUJBTHh3MWgzcXptTXE4YUErY01vWFQvM0UKRVVBdzRQNUhCZ\
      GZYa2VVa1J1ZjV3VHUvdTBvRHFyUnVRNlcwMFJWS09BUFFQVDZpRVhCQkdyKzdTRGhISi8yTAo4V\
      zBnNDRnaDB5N2lWREE1UmErV2dFS2Y4YmRvT3grUkMzME4wb3JyRmtZbDRCR3p2UWtDQXlOV1lZU\
      1lheGxDCkFEeDdoQldhMzFQb0FaT3MxN0pBM3Y0a1lUVzhzTVBHQzFtdnhCSHJ4cGswK3NxeFFLQ\
      jl2TzlaZ05zbXpQcXAKdUNsT0xQbkxjdkFrSjFMSXNIS3c4bXo3TlFnVnpVV3hzWEFTMVFnaTdXW\
      lNwRUVINGsySmtQRTRLNU5Mc2g2Mgo0d3hPS0lRRS83TFFsaVlCRlpIZnFYTHZpU0F5MW94SklTb\
      np4UmdmOFZVSmUrMUxKbGE0UHp6UDZuVUI2cGxUCnlKODQ2Z0QxQStxTmM1endlN2RtMytHSFVJN\
      VVRNXhZejNZRDdQclRIWlJJd3p3SzVzeXNTM28rTEFOL1duYjgKSnlVK2haNG5hcEZUUzh0SVlLW\
      GQ5V0pWUldYNWNzdkd2MXZSY1daK0tySGRtVi9Mc1oxbG04RFd1ODFLV0EyZQpKSmdvMlMrZ3R0R\
      FBEenZqMHVnbUNEVFVRZGpsYjVQWnhFdWI4cVYvUXViMTBCQ3JJSkNZenJHMm91dld6TVJLCmJIU\
      EE0OWtIVVYrR3BSdW8vNDZZYUw2T3haVjRQcnJjWFU1MVlQTVNZRnJiWU1iT2g1b1NxTnkyY1lrd\
      jRhTkUKMDNzQUllclJHMzYvaktnelYzNFFkZGJvTEI5Yk1sZThEaGJQZ3JEV2VGcTdJRlpBVlYwO\
      UcwdGx6R1hNaDE5YQpLZm8yVnNMRDZMT1dLTGwvOUNVRkFnTUJBQUV3RFFZSktvWklodmNOQVFFT\
      EJRQURnZ0lCQUNvM3ZYcHBJTk9TCkxkSEpTRXRYT3haSi9kMEw4aXZQaEUwZzlRb0h0NUNNUUZWM\
      GlKcVpzZzVLQURDckJSUG42ZnBVZDFhaDROZWoKWW5lUWxzbEpUOFFJN1RIQ1ZiTzBuOUtsbEVxW\
      EJ5c2gxeDdkTGVXdFVaczRxLzB2QVlxbTc3aWYyaVpWTEFzNgprMk81NmhpelNwN2svZC90SXd3a\
      1k3dTN4UkF2b2RWWEM3cnA1ZmViSTRWcDZLUFMwbjF4Yk8wSXBqUzd3KzhPClkxMjFXdmtxUUQyc\
      khCMkUySFJSSWdGTEF2eVFqY01hak9KMzRSUnh4VDljUzFIWEUzVUdDVVl2bHc1NnMySmMKMWYvM\
      XNXYTJNaHFiUUdzVDZJS3dVOXZ2SlArQktoWnBsWnpicUlnZHZyZHVHNWFCbEpnNlNCQTZsUThST\
      DhnNQp6MmpscGFaTE5jUmYySkdTUm5TTzlDZzJTcU5va2wwaDRRZVRCZEFNeGhPVnk0MmxOc1grL\
      zVOWVhuM0RXRWo1CmoybTV0SFpRU3ptY3BZZmNlY0ZpRXdud3BwOVdMbWlsb0N3SHVON2FNYU9Za\
      kpSVnVPcWJoTXZMUjFLTnpWVm0KUFdVblZJa3FBY2p0TXVKVFh6clZGV081VXEvazF0WkIzdGhNN\
      Ed4c1JZcElDbEJmS3c3emFIaGxsZmNvNE53eQprbWVDeENWMnVLZ3AxYkJlZ0hReTVEc3E4OFI5Q\
      kh1b0g0aWpEbzlUTGRpU2NXT0o1TkFEOUJ5K2F3VGpsUk1qCmtrVHgyMFRZdm1VdEF6T2N5R0xCO\
      ENEWERwS1B4b1ZBdVREMGFKMy9na3p0aHZQQ1dDamxINGI2cXVMS3lwVHgKOU85OVRZWGw3Y2VNS\
      XZIc2I3L3l4L01KUG9RbUZjcUgKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ=="
- name: "dtstack-172-16-8-160"
  cluster:
    server: "https://172.16.8.160:6443"
    certificate-authority-data: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQ\
      WFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKT\
      FdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBM\
      VVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ\
      0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4N\
      zhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTS\
      G84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLc\
      XdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1d\
      zdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmT\
      llGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqe\
      GZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93U\
      UZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNN\
      mRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtV\
      DlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQ\
      y9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hc\
      DdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRS\
      lVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZ\
      EVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJR\
      klDQVRFLS0tLS0K"
- name: "dtstack-172-16-8-166"
  cluster:
    server: "https://172.16.8.166:6443"
    certificate-authority-data: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQ\
      WFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKT\
      FdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBM\
      VVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ\
      0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4N\
      zhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTS\
      G84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLc\
      XdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1d\
      zdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmT\
      llGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqe\
      GZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93U\
      UZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNN\
      mRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtV\
      DlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQ\
      y9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hc\
      DdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRS\
      lVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZ\
      EVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJR\
      klDQVRFLS0tLS0K"
- name: "dtstack-172-16-8-169"
  cluster:
    server: "https://172.16.8.169:6443"
    certificate-authority-data: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3akNDQ\
      WFxZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkcmRXSmwKT\
      FdOaE1CNFhEVEl3TVRFd05UQXlNekF3T0ZvWERUTXdNVEV3TXpBeU16QXdPRm93RWpFUU1BNEdBM\
      VVFQXhNSAphM1ZpWlMxallUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ\
      0VCQUsyQmpTMEdYVmo5ClNZQWpjTFhYWUR2bXo2OStHdWkvbHQ4WHp0aC9ZdG4rU1NxR2NXSms4N\
      zhLaHI5RGdRQ1hNR3VtZXZtbUhFemIKYVVsY29oQXU0Q0NaS0VmYmxncFY2ZzFqNGtreVFxS2hTS\
      G84VDBhU2VLOW53dDZZMzQ4MXVxTG9UWHhubFhPKwpZblJOWjVsQnlVTWRXb2JmSGhzRGJQM2lLc\
      XdkU3NYckhxYWJlNmFDdk5tMlpIelNJbnc1eG0wakZHWXJvS0M5Cll3a3FVYmpWMDEzRXV1bzV1d\
      zdKZ2l1eFdFcDhoSjJ1WEFKNkp3OGZEeml5aUxIay9qTGVyTFM2UnN5dWNtYXEKWUJPclIvS2pmT\
      llGV1hlTWlqNzF3NXZqRTJFcG1DUnRLMFdtRjFpU3gyUFo2ZlBtYmVSd2lpRUVScXFFZGhYbQpqe\
      GZNRnRWSzUzY0NBd0VBQWFNak1DRXdEZ1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93U\
      UZNQU1CCkFmOHdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSTRxcUxnVFRJa1VQbGtJUWRGdjlNN\
      mRseHpJTlZmU0xXU2QKdEJuZGV6NUVmelBRaVBOdHQzbDlHNklrZmpYSGJUcU1rN0NsTWkvUnZtV\
      DlDcjMwc093NHY4WVNOQTFyYjBHYwp2VURJdEl0ZGoyVW8wSUtzQ2dsN0psalhMcVNvd01jczMwQ\
      y9IQXJGdU9pSTlsS2VBdlVIeFhFTG5ST1hxaGozCjhoTG5kS0dTak8rYXRFVSs1b0ZFa0pydS9hc\
      DdjTGxWSWVQVmh3eHpuM2pSS2dmRDVXTmI4aHRpRzg1YjMrNVIKMGt4eGFjNlRyaDBKQkh5cVpRS\
      lVWMUg2djk3NXdqNnZJNU5NcDNvbnlVZjJxaGJzSjVIWEJJRDh2NEtUQms3MQoxVnliMzFwY3NjZ\
      EVibGNPTFlveTFyR1VDT1QxMzNmbFlyTldwc2VXcFJ6Y2NPUE1WYVE9Ci0tLS0tRU5EIENFUlRJR\
      klDQVRFLS0tLS0K"

users:
- name: "dtstack"
  user:
    token: "kubeconfig-user-mtrk8.c-vckqj:jk8k6n72dpxtllmtd4p99kwkkxthxsqg6562b9r2t2v7bdkf2jft9n"


contexts:
- name: "dtstack"
  context:
    user: "dtstack"
    cluster: "dtstack"
- name: "dtstack-172-16-8-160"
  context:
    user: "dtstack"
    cluster: "dtstack-172-16-8-160"
- name: "dtstack-172-16-8-166"
  context:
    user: "dtstack"
    cluster: "dtstack-172-16-8-166"
- name: "dtstack-172-16-8-169"
  context:
    user: "dtstack"
    cluster: "dtstack-172-16-8-169"

current-context: "dtstack"
`
