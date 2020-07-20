package finalize

import (
	"context"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Finalize struct {
	finalizerID string
	writer      client.Writer
}

type object interface {
	runtime.Object
	GetDeletionTimestamp() *meta.Time
	GetFinalizers() []string
	SetFinalizers(finalizers []string)
}

func New(writer client.Writer, finalizerID string) *Finalize {
	return &Finalize{
		finalizerID: finalizerID,
		writer:      writer,
	}
}

func (f Finalize) Handle(o object, finalize func() error) (finalized bool, err error) {

	if o.GetDeletionTimestamp().IsZero() {
		if !f.hasFinalizer(o.GetFinalizers()) {
			o.SetFinalizers(f.addFinalizer(o.GetFinalizers()))
			return false, f.writer.Update(context.Background(), o)
		}
		return false, nil
	}
	if !f.hasFinalizer(o.GetFinalizers()) {
		return true, nil
	}

	if err := finalize(); err != nil {
		return false, err
	}
	o.SetFinalizers(f.removeFinalizer(o.GetFinalizers()))
	return true, f.writer.Update(context.Background(), o)

}

func (f Finalize) hasFinalizer(finalizers []string) bool {
	for _, finalizerID := range finalizers {
		if f.finalizerID == finalizerID {
			return true
		}
	}
	return false
}

func (f Finalize) removeFinalizer(finalizers []string) []string {
	for i, item := range finalizers {
		if item == f.finalizerID {
			finalizers = append(finalizers[:i], finalizers[i+1:]...)
		}
	}
	return finalizers
}

func (f Finalize) addFinalizer(finalizers []string) []string {
	return append(finalizers, f.finalizerID)
}
