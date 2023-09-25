package healthlake

import (
	"context"
	"os"
	"testing"
)

func TestHealthLakeCreateResource(t *testing.T) {
	conf := NewConfig(awsRegion, "823359e39e1b25fb5cd7a9ebce86b8f9")
	newHLR := NewHealthLakeRepository(conf)
	content, err := os.ReadFile("bundle.json")
	if err != nil {
		t.Fatal("failed reading bundle")
		return
	}
	newHLR.CreateResource(context.TODO(), content, bundleResourceType)
}
