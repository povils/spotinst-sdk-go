package mcs

import (
	"context"
	"encoding/json"
	"github.com/povils/spotinst-sdk-go/spotinst"
	"github.com/povils/spotinst-sdk-go/spotinst/client"
	"github.com/povils/spotinst-sdk-go/spotinst/util/uritemplates"
	"io/ioutil"
	"net/http"
)

type ClusterCostInput struct {
	ClusterID *string `json:"clusterId,omitempty"`
	ToDate    *string `json:"toDate,omitempty"`
	FromDate  *string `json:"fromDate,omitempty"`
}

type ClusterCostOutput struct {
	ClusterCosts []*ClusterCost `json:"clusterCosts,omitempty"`
}

type ClusterCost struct {
	TotalCost          *float64      `json:"totalCost,omitempty"`
	Namespaces         []*Namespace  `json:"namespaces,omitempty"`
	Deployments        []*Deployment `json:"deployments,omitempty"`
	StandAlonePodsCost *float64      `json:"standAlonePodsCost,omitempty"`
	HeadroomCost       *float64      `json:"headroomCost,omitempty"`
}

type Namespace struct {
	Namespace *string  `json:"namespace,omitempty"`
	Cost      *float64 `json:"cost,omitempty"`
}

type Deployment struct {
	DeploymentName *string  `json:"deploymentName,omitempty"`
	Namespace      *string  `json:"namespace,omitempty"`
	Cost           *float64 `json:"cost,omitempty"`
}

func clusterCostFromJSON(in []byte) (*ClusterCost, error) {
	b := new(ClusterCost)
	if err := json.Unmarshal(in, b); err != nil {
		return nil, err
	}
	return b, nil
}

func clusterCostsFromJSON(in []byte) ([]*ClusterCost, error) {
	var rw client.Response
	if err := json.Unmarshal(in, &rw); err != nil {
		return nil, err
	}
	out := make([]*ClusterCost, len(rw.Response.Items))
	if len(out) == 0 {
		return out, nil
	}
	for i, rb := range rw.Response.Items {
		b, err := clusterCostFromJSON(rb)
		if err != nil {
			return nil, err
		}
		out[i] = b
	}
	return out, nil
}

func clusterCostsFromHttpResponse(resp *http.Response) ([]*ClusterCost, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return clusterCostsFromJSON(body)
}

// GetClusterCosts accepts kubernetes clusterId, fromDate, and toDate and returns a list of costs.
// Dates may be in the format of yyyy-mm-dd or Unix Timestamp (1494751821472)
func (s *ServiceOp) GetClusterCosts(ctx context.Context, input *ClusterCostInput) (*ClusterCostOutput, error) {
	path, err := uritemplates.Expand("/mcs/kubernetes/cluster/{clusterIdentifier}/costs", uritemplates.Values{
		"clusterIdentifier": spotinst.StringValue(input.ClusterID),
	})
	if err != nil {
		return nil, err
	}

	r := client.NewRequest(http.MethodGet, path)

	if input.ToDate != nil {
		r.Params.Set("toDate", *input.ToDate)
	}
	if input.FromDate != nil {
		r.Params.Set("fromDate", *input.FromDate)
	}
	r.Obj = input

	resp, err := client.RequireOK(s.Client.Do(ctx, r))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	costs, err := clusterCostsFromHttpResponse(resp)
	if err != nil {
		return nil, err
	}

	return &ClusterCostOutput{costs}, nil
}
