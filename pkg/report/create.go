package report

import (
	"strconv"

	"github.com/haardikdharma10/kubearmor-adapter/pkg/api/wgpolicyk8s.io/v1alpha2"
	policyreport "github.com/haardikdharma10/kubearmor-adapter/pkg/api/wgpolicyk8s.io/v1alpha2"
	pb "github.com/kubearmor/KubeArmor/protobuf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//global variable to keep the source-name static
const PolicyReportSource string = "KubeArmor Policy Engine"

func Create(alert *pb.Alert) (*policyreport.PolicyReport, error) {

	report := &policyreport.PolicyReport{
		ObjectMeta: metav1.ObjectMeta{
			Name: "KubeArmor Policy Report",
		},
		Summary: v1alpha2.PolicyReportSummary{
			Fail: 1,
		},
	}

	//Storing the result obtained after mapping in "r".
	r := newResult(alert)

	//appending to policyreport result.
	report.Results = append(report.Results, r)

	//fmt.Printf("Created policy-report %q.\n", result.GetObjectMeta().GetName())

	return report, nil
}

func newResult(Alert *pb.Alert) *policyreport.PolicyReportResult {

	var sev string

	if Alert.Severity == "1" || Alert.Severity == "2" {
		sev = "low"
	} else if Alert.Severity == "3" || Alert.Severity == "4" || Alert.Severity == "5" {
		sev = "medium"
	} else {
		sev = "high"
	}

	//Mapping:-
	return &policyreport.PolicyReportResult{

		Source: PolicyReportSource,
		Policy: Alert.PolicyName,
		Scored: false,
		// Timestamp:   metav1.Timestamp{Seconds: int64(Alert.UpdatedTime), Nanos: int32(Alert.Timestamp.Nanosecond())},
		Severity:    v1alpha2.PolicyResultSeverity(sev),
		Result:      "fail",
		Description: Alert.Message,
		Category:    Alert.Type,
		Properties: map[string]string{
			"cluster_name":   Alert.ClusterName,
			"host_name":      Alert.HostName,
			"namespace_name": Alert.NamespaceName,
			"pod_name":       Alert.PodName,
			"container_id":   Alert.ContainerID,
			"container_name": Alert.ContainerName,
			"host_pid":       strconv.Itoa(int(Alert.HostPID)),
			"ppid":           strconv.Itoa(int(Alert.PPID)),
			"pid":            strconv.Itoa(int(Alert.PID)),
			"tags":           Alert.Tags,
			//"source" : Alert.Source,
			"operation": Alert.Operation,
			"resource":  Alert.Resource,
			"data":      Alert.Data,
			"action":    Alert.Action,
			//"result" : Alert.Result,
		},
	}
}
