package consensus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rgreen312/owlplace/server/common"
)

type MembershipProvider interface {
	// ListMembers provides a map from nodeIDs to addresses and ports.  Note
	// that the results of this method are subject to change between calls.
	GetMembership() (map[uint64]string, error)
}

type StaticMembershipProvider struct {
	members map[uint64]string
}

// StaticMembershipFromFile constructs a StaticMembershipProvider from a given
// file.  This method will error if the provided file does not exist, or has
// invalid syntax.  A well-formed file is JSON, with mappings from uint64
// nodeIDs to cluster URIs:
//
//    {
//        "1": "backend1:63000",
//        "2": "backend2:63000",
//        "3": "backend3:63000"
//    }
func StaticMembershipFromFile(path string) (*StaticMembershipProvider, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "reading membership file: %s", path)
	}

	var members map[uint64]string
	err = json.Unmarshal([]byte(file), &members)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshalling membership file: %s", path)
	}

	return &StaticMembershipProvider{
		members: members,
	}, nil
}

func (p *StaticMembershipProvider) GetMembership() (map[uint64]string, error) {
	return p.members, nil
}

type KubernetesMembershipProvider struct {
	// The namespace with which to search for pods.
	namespace string
	// API Handle with which to perform k8s queries with.
	clientset *kubernetes.Clientset
}

// NewKubernetesMembershipProvider constructs a MembershipProvider which
// queries a k8s namespace for existing pods when listing cluster members.
//
// TODO: We currently take a namespace as a parameter, although it would be
// more desirable to take an interface can act as a pod provider, as it would
// be easier to test.
func NewKubernetesMembershipProvider(namespace string) (*KubernetesMembershipProvider, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving cluster config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "creating API handle")
	}

	return &KubernetesMembershipProvider{
		namespace: namespace,
		clientset: clientset,
	}, nil
}

func (p *KubernetesMembershipProvider) GetMembership() (map[uint64]string, error) {
	pods, err := p.clientset.CoreV1().Pods("dev").List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "listing dev pods")
	}

	members := make(map[uint64]string)
	for _, pod := range pods.Items {
		members[common.IPToNodeId(pod.Status.PodIP)] = fmt.Sprintf("%s:%d", pod.Status.PodIP, common.ConsensusPort)
	}

	return members, nil
}
