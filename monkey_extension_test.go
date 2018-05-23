package monkeyExtensions

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"reflect"

	"github.com/bouk/monkey"
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

		AfterEach(func() {
			monkey.UnpatchAll()
		})

		Context("when isn't patched", func() {
			It("should return original data", func() {
				Expect(err).To(Succeed())
				Expect(res).To(Equal(111))
			})
		})

		Context("when is patched", func() {

			Context("when patch is valid", func() {
				var ret *monkey.PatchGuard

				Context("when replacement func doesn't have params", func() {
					BeforeEach(func() {
						replacement := func() (int, error) { return 222, errors.New("err") }
						ret = PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)
					})

					It("should return patched data", func() {
						Expect(err).To(HaveOccurred())
						Expect(res).To(Equal(222))
						Expect(ret).To(Not(BeNil()))
					})
				})

				Context("when replacement func has 1 params", func() {
					BeforeEach(func() {
						replacement := func(_ interface{}, param int) (int, error) { return param * 3, errors.New("err") }
						ret = PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)
					})

					It("should return patched data", func() {
						Expect(err).To(HaveOccurred())
						Expect(res).To(Equal(333))
						Expect(ret).To(Not(BeNil()))
					})
				})

				Context("when replacement uses passed params", func() {
					BeforeEach(func() {
						replacement := func(_ interface{}, param int) (int, error) { return param * 3, errors.New("err") }
						ret = PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)

					})

					It("should return patched data", func() {
						Expect(err).To(HaveOccurred())
						Expect(res).To(Equal(333))
						Expect(ret).To(Not(BeNil()))
					})
				})

				Context("when unpatched all", func() {
					BeforeEach(func() {
						replacement := func() (int, error) { return 222, errors.New("err") }
						ret = PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)
						monkey.UnpatchAll()
					})

					It("should return not patched data", func() {
						Expect(err).To(Succeed())
						Expect(res).To(Equal(111))
						Expect(ret).To(Not(BeNil()))
					})
				})

				// Issue with unpatch instance method in base monkey project:
				// Unmerged fix (23.05.2018):
				// https://github.com/bouk/monkey/pull/16

				//Context("when unpatched one", func() {
				//	BeforeEach(func() {
				//		replacement := func(_ dummy, a int) (int, error) { return 10000 + a, errors.New("err") }
				//		PatchInstanceMethodFlexible(reflect.TypeOf(subject), "Foo", replacement)
				//		//monkey.PatchInstanceMethod(reflect.TypeOf(subject), "Foo", replacement)
				//		monkey.UnpatchInstanceMethod(reflect.TypeOf(subject), "Foo")
				//	})
				//
				//	It("should return not patched data", func() {
				//		Expect(err).To(Succeed())
				//		Expect(res).To(Equal(111))
				//	})
				//})
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
