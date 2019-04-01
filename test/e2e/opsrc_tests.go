package e2e

import (
	"fmt"
	"testing"
	"strings"

	operator "github.com/operator-framework/operator-marketplace/pkg/apis/operators/v1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func runOpSrcTests(t *testing.T) {
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()

	// Get global framework variables
	f := test.Global
	// Run tests
	if err := createOpSrcWithInvalidURL(t, f, ctx); err != nil {
		t.Fatal(err)
	}
	if err := createOpSrcWithInvalidEndpoint(t, f, ctx); err != nil {
		t.Fatal(err)
	}
	if err := createOpSrcWithNonexistentRegistryNamespace(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}

// Create OperatorSource with invalid URL
// Expected result: OperatorSource reaches failed state
func createOpSrcWithInvalidURL(t *testing.T, f *test.Framework, ctx *test.TestCtx) error {
	opSrcName := "invalid-url-opsrc"
	// invalidURL is an invalid URI
	invalidURL := "not-a-url"

	// Get test namespace
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("Could not get namespace: %v", err)
	}

	invalidURLOperatorSource := &operator.OperatorSource{
		TypeMeta: metav1.TypeMeta{
			Kind: operator.OperatorSourceKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opSrcName,
			Namespace: namespace,
		},
		Spec: operator.OperatorSourceSpec{
			Type:              "appregistry",
			Endpoint:          invalidURL,
			RegistryNamespace: "marketplace_e2e",
		},
	}
	err = createRuntimeObject(f, ctx, invalidURLOperatorSource)
	if err != nil {
		return err
	}

	// Check that OperatorSource reaches "Failed" state eventually
	resultOperatorSource := &operator.OperatorSource{}
	expectedPhase := "Failed"
	err = wait.Poll(retryInterval, timeout, func() (bool, error) {
		err = WaitForResult(t, f, resultOperatorSource, namespace, opSrcName)
		if err != nil {
			return false, err
		}
		if resultOperatorSource.Status.CurrentPhase.Name == expectedPhase &&
			strings.Contains(resultOperatorSource.Status.CurrentPhase.Message, "Invalid operator source endpoint") {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("OperatorSource never reached expected phase/message, expected %v", expectedPhase)
	}

	return nil
}

// Create OperatorSource with invalid endpoint
// Expected result: OperatorSource stuck in downloading state
func createOpSrcWithInvalidEndpoint(t *testing.T, f *test.Framework, ctx *test.TestCtx) error {
	opSrcName := "invalid-endpoint-opsrc"
	// invalidEndpoint is the invalid endpoint for the OperatorSource
	invalidEndpoint := "https://not-quay.io/cnr"

	// Get test namespace
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("Could not get namespace: %v", err)
	}

	invalidURLOperatorSource := &operator.OperatorSource{
		TypeMeta: metav1.TypeMeta{
			Kind: operator.OperatorSourceKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opSrcName,
			Namespace: namespace,
		},
		Spec: operator.OperatorSourceSpec{
			Type:              "appregistry",
			Endpoint:          invalidEndpoint,
			RegistryNamespace: "marketplace_e2e",
		},
	}
	err = createRuntimeObject(f, ctx, invalidURLOperatorSource)
	if err != nil {
		return err
	}

	// Check that OperatorSource is in "Downloading" state with approprate message
	resultOperatorSource := &operator.OperatorSource{}
	expectedPhase := "Downloading"
	err = wait.Poll(retryInterval, timeout, func() (bool, error) {
		err = WaitForResult(t, f, resultOperatorSource, namespace, opSrcName)
		if err != nil {
			return false, err
		}
		if resultOperatorSource.Status.CurrentPhase.Name == expectedPhase &&
			strings.Contains(resultOperatorSource.Status.CurrentPhase.Message, "no such host") {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("OperatorSource never reached expected phase/message, expected %v", expectedPhase)
	}

	return nil
}

// Create OperatorSource with valid URL but non-existent registry namespace
// Expected result: OperatorSource reaches failed state
func createOpSrcWithNonexistentRegistryNamespace(t *testing.T, f *test.Framework, ctx *test.TestCtx) error {
	opSrcName := "nonexistent-namespace-opsrc"
	// validURL is a valid endpoint for the OperatorSource
	validURL := "https://quay.io/cnr"

	// nonexistentRegistryNamespace is a namespace that does not exist
	// on the app registry
	nonexistentRegistryNamespace := "not-existent-namespace"

	// Get test namespace
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("Could not get namespace: %v", err)
	}

	nonexistentRegistryNamespaceOperatorSource := &operator.OperatorSource{
		TypeMeta: metav1.TypeMeta{
			Kind: operator.OperatorSourceKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opSrcName,
			Namespace: namespace,
		},
		Spec: operator.OperatorSourceSpec{
			Type:              "appregistry",
			Endpoint:          validURL,
			RegistryNamespace: nonexistentRegistryNamespace,
		},
	}
	err = createRuntimeObject(f, ctx, nonexistentRegistryNamespaceOperatorSource)
	if err != nil {
		return err
	}

	// Check that OperatorSource reaches "Failed" state eventually
	resultOperatorSource := &operator.OperatorSource{}
	expectedPhase := "Failed"
	err = wait.Poll(retryInterval, timeout, func() (bool, error) {
		err = WaitForResult(t, f, resultOperatorSource, namespace, opSrcName)
		if err != nil {
			return false, err
		}
		if resultOperatorSource.Status.CurrentPhase.Name == expectedPhase &&
			strings.Contains(resultOperatorSource.Status.CurrentPhase.Message, "The operator source endpoint returned an empty manifest list") {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("OperatorSource never reached expected phase/message, expected %v", expectedPhase)
	}

	return nil
}