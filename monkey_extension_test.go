package monkeyExtensions

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"reflect"

	"github.com/pkg/errors"
)

func TestMonkeyExtension(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MonkeyExtension suit")
}

var _ = Describe("MonkeyExtension", func() {

	Describe("PatchInstanceMethodFlexible func", func() {

		var (
			subject dummy
			res     int
			err     error
		)

		JustBeforeEach(func() {
			res, err = subject.Foo(111)
		})

		Context("when isn't patched", func() {
			It("should return original data", func() {
				Expect(err).To(Succeed())
				Expect(res).To(Equal(111))
			})
		})

		Context("when is patched", func() {

			Context("when patch is valid", func() {

				Context("when replacement func doesn't have params", func() {
					BeforeEach(func() {
						replacement := func() (int, error) { return 222, errors.New("err") }
						PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)

					})

					It("should return patched data", func() {
						Expect(err).To(HaveOccurred())
						Expect(res).To(Equal(222))
					})
				})

				Context("when replacement func has 1 params", func() {
					BeforeEach(func() {
						replacement := func(_ interface{}, param int) (int, error) { return param * 3, errors.New("err") }
						PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)

					})

					It("should return patched data", func() {
						Expect(err).To(HaveOccurred())
						Expect(res).To(Equal(333))
					})
				})

				Context("when replacement uses passed params", func() {
					BeforeEach(func() {
						replacement := func(_ interface{}, param int) (int, error) { return param * 3, errors.New("err") }
						PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)

					})

					It("should return patched data", func() {
						Expect(err).To(HaveOccurred())
						Expect(res).To(Equal(333))
					})
				})
			})

			Context("when patch is not valid", func() {
				Context("when replacement defined with too many input params", func() {
					It("should panic", func() {
						replacement := func(_, _, _ interface{}) (int, error) { return 0, nil }
						Expect(func() {
							PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)
						}).To(Panic())
					})
				})
			})
		})
	})
})

type dummy struct{}

func (d dummy) Foo(a int) (int, error) {
	if a >= 0 {
		return a, nil
	}

	return 0, errors.New("negatives are not allowed")
}
