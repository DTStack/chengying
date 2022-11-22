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

package kube

import (
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"testing"
)


func getConn(){
	user := "root"
	password := "dtstack"
	host:= "172.16.10.37"
	port := 3306
	dbname := "dtagent"
	log.ConfigureLogger("/tmp/matrix",100,3,1)
	db,_ := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true", user, password, host, port, dbname))
	model.MYSQLDB = db
	err := Build()
	if err != nil{
		fmt.Printf("err: %v \n",err)
	}
}
func TestGetClusterSts(t *testing.T) {
	getConn()
	info := &ClusterInfo{
		Name:          "ccaqdwa",
		IsDeleted:     0,
	}
	if err := getClusterSts.Get(info,info);err!=nil{
		fmt.Println(err.Error())
	}
	bts,_:=json.Marshal(info)
	fmt.Println(string(bts))
}

func TestInsertClusterSts(t *testing.T) {
	getConn()
	info := &ClusterInfo{
		Name:          "cccccc",
		Type:          "kubernetes",
		Mode:          0,
		Version:       "1.11.1",
		IsDeleted:     0,
		//NetworkPlugin: "flannel",
	}
	i,err := DeployClusterList.InsertK8sCluster(info)
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println("------",i)
}

func TestGetClusterK8sAvailableByMode(t *testing.T) {
	getConn()
	as,err := DeployClusterK8sAvailable.GetClusterK8sAvailableByMode(1)
	if err != nil{
		fmt.Println(err.Error())
	}
	for _,a := range as{
		bts,_:=json.Marshal(a)
		fmt.Println(string(bts))
	}
}

func TestGetRealVersion(t *testing.T) {
	getConn()
	s,err := DeployClusterK8sAvailable.GetRealVersion("")
	if err != nil{
		fmt.Println("err ",err.Error())
	}
	fmt.Println(s)
}

func TestInitmoudle(t *testing.T) {
	getConn()
	s,err := ImportInitMoudle.GetInitMoudle()
	if err != nil{
		fmt.Println(err.Error())
	}

	fmt.Println(s)
}
//func TestInserNamspaceList(t *testing.T) {
//	getConn()
//	ss,err := DeployNamespaceList.GetByCluster(1)
//	if err != nil{
//		fmt.Println("err",err.Error())
//
//	}
//	for _,s := range ss{
//		fmt.Printf("s %+v \n",s)
//	}
//	//err := DeployNamespaceList.UpdateStatus(&s)
//	//if err != nil{
//	//	fmt.Println("err5",err.Error())
//	//}
//}

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
func TestInsertConnectInfo(t *testing.T) {
	getConn()
	config := kubeconfig
	tbsc := DeployNamespaceClientSchema{
		Yaml:        config,
		NamespaceId: 1,
		Filename:    "ffff",
	}
	DeployNamespaceClient.Insert(&tbsc)
}

func TestGetEventCount(t *testing.T) {
	getConn()
	count,err := DeployNamespaceEvent.SelectCount(7)
	if err != nil{
		fmt.Println("err1",err.Error())
		return
	}
	fmt.Println(count)
}
//
//func TestGetnamespacelist(t *testing.T) {
//	getConn()
//	scs,err := DeployNamespaceList.Select("85","valid","desc")
//	if err != nil{
//		fmt.Println("err1",err.Error())
//		return
//	}
//	bts,_ := yaml.Marshal(scs)
//	fmt.Println("bts",string(bts))
//
//}
