package v2_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/v2"
	"code.cloudfoundry.org/cli/command/v2/shared"
	"code.cloudfoundry.org/cli/command/v2/v2fakes"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("target Command", func() {
	var (
		cmd        v2.TargetCommand
		testUI     *ui.UI
		fakeActor  *v2fakes.FakeTargetActor
		fakeConfig *commandfakes.FakeConfig
		executeErr error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeActor = new(v2fakes.FakeTargetActor)
		fakeConfig = new(commandfakes.FakeConfig)

		cmd = v2.TargetCommand{
			UI:     testUI,
			Actor:  fakeActor,
			Config: fakeConfig,
		}

		fakeConfig.BinaryNameReturns("faceman")
		fakeConfig.ExperimentalReturns(true)
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	It("Displays the experimental warning message", func() {
		Expect(testUI.Out).To(Say(command.ExperimentalWarning))
	})

	Describe("Cloud Controller minimum version warning", func() {
		var (
			minCLIVersion string
			binaryVersion string
			apiVersion    string
		)

		BeforeEach(func() {
			apiVersion = "6.0.0"
			minCLIVersion = "1.0.0"
			fakeConfig.APIVersionReturns(apiVersion)
			fakeConfig.MinCLIVersionReturns(minCLIVersion)
		})

		Context("when the CLI version is less than the recommended minimum", func() {
			BeforeEach(func() {
				binaryVersion = "0.0.0"
				fakeConfig.BinaryVersionReturns(binaryVersion)
			})

			It("displays a recommendation to update the CLI version", func() {
				Expect(testUI.Err).To(Say("Cloud Foundry API version %s requires CLI version %s. You are currently on version %s. To upgrade your CLI, please visit: https://github.com/cloudfoundry/cli#downloads", apiVersion, minCLIVersion, binaryVersion))
			})
		})

		Context("when the CLI version is greater or equal to the recommended minimum", func() {
			BeforeEach(func() {
				binaryVersion = "1.0.0"
				fakeConfig.BinaryVersionReturns(binaryVersion)
			})

			It("does not display a recommendation to update the CLI version", func() {
				Expect(testUI.Err).NotTo(Say("Cloud Foundry API version %s requires CLI version %s. You are currently on version %s. To upgrade your CLI, please visit: https://github.com/cloudfoundry/cli#downloads", apiVersion, minCLIVersion, binaryVersion))
			})
		})

		Context("when an error is encountered while parsing the semver versions", func() {
			BeforeEach(func() {
				fakeConfig.BinaryVersionReturns("&#%")
			})

			It("does not recommend to update the CLI version", func() {
				Expect(testUI.Err).NotTo(Say("Cloud Foundry API version %s requires CLI version %s.", apiVersion, minCLIVersion))
			})
		})
	})

	Context("when the user is not logged in", func() {
		It("returns an error", func() {
			Expect(executeErr).To(MatchError(command.NotLoggedInError{
				BinaryName: "faceman",
			}))
		})
	})

	Context("when getting the current user returns an error", func() {
		var someErr error

		BeforeEach(func() {
			someErr = errors.New("some-current-user-error")
			fakeConfig.AccessTokenReturns("some-access-token")
			fakeConfig.CurrentUserReturns(configv3.User{}, someErr)
		})

		It("returns the same error", func() {
			Expect(executeErr).To(MatchError(someErr))
		})
	})

	Context("when the user is logged in", func() {
		BeforeEach(func() {
			fakeConfig.TargetReturns("some-api-target")
			fakeConfig.APIVersionReturns("1.2.3")
			fakeConfig.AccessTokenReturns("some-access-token")
			fakeConfig.RefreshTokenReturns("some-refresh-token")
			fakeConfig.CurrentUserReturns(configv3.User{
				Name: "some-user",
			}, nil)
		})

		Context("when no arguments are given", func() {
			Context("when no org or space are targeted", func() {
				It("displays how to target an org and space", func() {
					Expect(executeErr).ToNot(HaveOccurred())

					Expect(testUI.Out).To(Say("API endpoint:   some-api-target"))
					Expect(testUI.Out).To(Say("API version:    1.2.3"))
					Expect(testUI.Out).To(Say("User:           some-user"))
					Expect(testUI.Out).To(Say("No org or space targeted, use 'faceman target -o ORG -s SPACE'"))
				})
			})

			Context("when an org but no space is targeted", func() {
				BeforeEach(func() {
					fakeConfig.HasTargetedOrganizationReturns(true)
					fakeConfig.TargetedOrganizationReturns(configv3.Organization{
						GUID: "some-org-guid",
						Name: "some-org",
					})
				})

				It("displays the org and tip to target space", func() {
					Expect(executeErr).ToNot(HaveOccurred())

					Expect(testUI.Out).To(Say("API endpoint:   some-api-target"))
					Expect(testUI.Out).To(Say("API version:    1.2.3"))
					Expect(testUI.Out).To(Say("User:           some-user"))
					Expect(testUI.Out).To(Say("Org:            some-org"))
					Expect(testUI.Out).To(Say("No space targeted, use 'faceman target -s SPACE'"))
				})
			})

			Context("when an org and space are targeted", func() {
				BeforeEach(func() {
					fakeConfig.HasTargetedOrganizationReturns(true)
					fakeConfig.TargetedOrganizationReturns(configv3.Organization{
						GUID: "some-org-guid",
						Name: "some-org",
					})
					fakeConfig.HasTargetedSpaceReturns(true)
					fakeConfig.TargetedSpaceReturns(configv3.Space{
						GUID: "some-space-guid",
						Name: "some-space",
					})
				})

				It("displays the org and space targeted ", func() {
					Expect(executeErr).ToNot(HaveOccurred())

					Expect(testUI.Out).To(Say("API endpoint:   some-api-target"))
					Expect(testUI.Out).To(Say("API version:    1.2.3"))
					Expect(testUI.Out).To(Say("User:           some-user"))
					Expect(testUI.Out).To(Say("Org:            some-org"))
					Expect(testUI.Out).To(Say("Space:          some-space"))
				})
			})
		})

		Context("when space is provided", func() {
			BeforeEach(func() {
				cmd.Space = "some-space"
			})

			Context("when an org is already targeted", func() {
				BeforeEach(func() {
					fakeConfig.HasTargetedOrganizationReturns(true)
					fakeConfig.TargetedOrganizationReturns(configv3.Organization{
						GUID: "some-org-guid",
					})
				})
				Context("when the space exists", func() {
					BeforeEach(func() {
						fakeActor.GetSpaceByOrganizationAndNameReturns(v2action.Space{
							GUID:     "some-space-guid",
							Name:     "some-space",
							AllowSSH: true,
						}, v2action.Warnings{}, nil)
					})

					It("targets the space", func() {
						Expect(executeErr).ToNot(HaveOccurred())

						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(1))
						spaceGUID, spaceName, spaceAllowSSH := fakeConfig.SetSpaceInformationArgsForCall(0)
						Expect(spaceGUID).To(Equal("some-space-guid"))
						Expect(spaceName).To(Equal("some-space"))
						Expect(spaceAllowSSH).To(Equal(true))
					})
				})

				Context("when the space does not exist", func() {
					BeforeEach(func() {
						fakeActor.GetSpaceByOrganizationAndNameReturns(v2action.Space{}, v2action.Warnings{}, v2action.SpaceNotFoundError{Name: "some-space"})
					})

					It("returns a SpaceNotFoundError and does not set change the space", func() {
						Expect(executeErr).To(MatchError(shared.SpaceNotFoundError{Name: "some-space"}))

						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(0))
					})
				})
			})

			Context("when no org is targeted", func() {
				It("returns NoOrgTargeted error", func() {
					Expect(executeErr).To(MatchError(shared.NoOrganizationTargetedError{}))
					Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(0))
				})
			})
		})

		Context("when org is provided", func() {
			BeforeEach(func() {
				cmd.Organization = "some-org"
			})

			Context("when the org does not exist", func() {
				BeforeEach(func() {
					fakeActor.GetOrganizationByNameReturns(
						v2action.Organization{},
						nil,
						v2action.OrganizationNotFoundError{Name: "some-org"})
				})

				It("displays all warnings and returns an org target error", func() {
					Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(0))

					Expect(executeErr).To(MatchError(shared.OrganizationNotFoundError{Name: "some-org"}))
				})
			})

			Context("when the org exists", func() {
				BeforeEach(func() {
					fakeConfig.HasTargetedOrganizationReturns(true)
					fakeConfig.TargetedOrganizationReturns(configv3.Organization{
						GUID: "some-org-guid",
						Name: "some-org",
					})
					fakeActor.GetOrganizationByNameReturns(
						v2action.Organization{
							GUID: "some-org-guid",
						},
						v2action.Warnings{
							"warning-1",
							"warning-2",
						},
						nil)
				})

				Context("when getting the organization's spaces returns an error", func() {
					var err error
					BeforeEach(func() {
						err = errors.New("get-org-spaces-error")
						fakeActor.GetOrganizationSpacesReturns(
							[]v2action.Space{},
							v2action.Warnings{
								"warning-3",
							}, err)
					})

					It("displays all warnings and returns a Get org spaces error", func() {
						Expect(fakeActor.GetOrganizationSpacesCallCount()).To(Equal(1))
						orgGUID := fakeActor.GetOrganizationSpacesArgsForCall(0)
						Expect(orgGUID).To(Equal("some-org-guid"))

						Expect(testUI.Err).To(Say("warning-1"))
						Expect(testUI.Err).To(Say("warning-2"))
						Expect(testUI.Err).To(Say("warning-3"))

						Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(1))
						orgGUID, orgName := fakeConfig.SetOrganizationInformationArgsForCall(0)
						Expect(orgGUID).To(Equal("some-org-guid"))
						Expect(orgName).To(Equal("some-org"))
						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(0))

						Expect(executeErr).To(MatchError(err))
					})
				})

				Context("when there are no spaces in the targeted org", func() {
					It("displays all warnings", func() {
						Expect(executeErr).ToNot(HaveOccurred())

						Expect(testUI.Err).To(Say("warning-1"))
						Expect(testUI.Err).To(Say("warning-2"))
					})

					It("sets the org and unsets the space in the config", func() {
						Expect(executeErr).ToNot(HaveOccurred())

						Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(1))
						orgGUID, orgName := fakeConfig.SetOrganizationInformationArgsForCall(0)
						Expect(orgGUID).To(Equal("some-org-guid"))
						Expect(orgName).To(Equal("some-org"))

						Expect(fakeConfig.UnsetSpaceInformationCallCount()).To(Equal(1))
						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(0))
					})
				})

				Context("when there is only 1 space in the targeted org", func() {
					BeforeEach(func() {
						fakeActor.GetOrganizationSpacesReturns([]v2action.Space{{
							GUID:     "some-space-guid",
							Name:     "some-space",
							AllowSSH: true,
						}}, v2action.Warnings{
							"warning-3",
						}, nil)
					})

					It("displays all warnings", func() {
						Expect(executeErr).ToNot(HaveOccurred())

						Expect(testUI.Err).To(Say("warning-1"))
						Expect(testUI.Err).To(Say("warning-2"))
						Expect(testUI.Err).To(Say("warning-3"))
					})

					It("targets the org and space", func() {
						Expect(executeErr).ToNot(HaveOccurred())

						Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(1))
						orgGUID, orgName := fakeConfig.SetOrganizationInformationArgsForCall(0)
						Expect(orgGUID).To(Equal("some-org-guid"))
						Expect(orgName).To(Equal("some-org"))

						Expect(fakeConfig.UnsetSpaceInformationCallCount()).To(Equal(1))

						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(1))
						spaceGUID, spaceName, spaceAllowSSH := fakeConfig.SetSpaceInformationArgsForCall(0)
						Expect(spaceGUID).To(Equal("some-space-guid"))
						Expect(spaceName).To(Equal("some-space"))
						Expect(spaceAllowSSH).To(Equal(true))
					})
				})

				Context("when there are multiple spaces in the targeted org", func() {
					BeforeEach(func() {
						fakeActor.GetOrganizationSpacesReturns([]v2action.Space{{
							GUID:     "some-space-guid",
							Name:     "some-space",
							AllowSSH: true,
						}, {
							GUID:     "another-space-space-guid",
							Name:     "another-space",
							AllowSSH: true,
						}}, v2action.Warnings{
							"warning-3",
						}, nil)
					})

					It("displays all warnings, sets the org, and clears the existing targetted space from the config", func() {
						Expect(executeErr).ToNot(HaveOccurred())

						Expect(testUI.Err).To(Say("warning-1"))
						Expect(testUI.Err).To(Say("warning-2"))

						Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(1))
						orgGUID, orgName := fakeConfig.SetOrganizationInformationArgsForCall(0)
						Expect(orgGUID).To(Equal("some-org-guid"))
						Expect(orgName).To(Equal("some-org"))

						Expect(fakeConfig.UnsetSpaceInformationCallCount()).To(Equal(1))
						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(0))
					})
				})
			})
		})

		Context("when org and space arguments are provided", func() {
			BeforeEach(func() {
				cmd.Space = "some-space"
				cmd.Organization = "some-org"
			})

			Context("when the org exists", func() {
				BeforeEach(func() {
					fakeActor.GetOrganizationByNameReturns(
						v2action.Organization{
							GUID: "some-org-guid",
							Name: "some-org",
						}, v2action.Warnings{
							"warning-1",
						}, nil)
				})

				Context("when the space exists", func() {
					BeforeEach(func() {
						fakeActor.GetSpaceByOrganizationAndNameReturns(v2action.Space{
							GUID: "some-space-guid",
							Name: "some-space",
						},
							v2action.Warnings{
								"warning-2",
							}, nil)
					})

					It("sets the target org and space", func() {
						Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(1))
						orgGUID, orgName := fakeConfig.SetOrganizationInformationArgsForCall(0)
						Expect(orgGUID).To(Equal("some-org-guid"))
						Expect(orgName).To(Equal("some-org"))

						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(1))
						spaceGUID, spaceName, allowSSH := fakeConfig.SetSpaceInformationArgsForCall(0)
						Expect(spaceGUID).To(Equal("some-space-guid"))
						Expect(spaceName).To(Equal("some-space"))
						Expect(allowSSH).To(BeFalse())
					})

					It("displays all warnings", func() {
						Expect(testUI.Err).To(Say("warning-1"))
						Expect(testUI.Err).To(Say("warning-2"))
					})
				})

				Context("when the space does not exist", func() {
					BeforeEach(func() {
						fakeActor.GetSpaceByOrganizationAndNameReturns(v2action.Space{}, nil, v2action.SpaceNotFoundError{Name: "some-space"})
					})

					It("returns an error and keeps old target", func() {
						Expect(executeErr).To(MatchError(shared.SpaceNotFoundError{Name: "some-space"}))

						Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(0))
						Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(0))
					})
				})
			})

			Context("when the org does not exist", func() {
				BeforeEach(func() {
					fakeActor.GetOrganizationByNameReturns(v2action.Organization{}, nil, v2action.OrganizationNotFoundError{Name: "some-org"})
				})

				It("returns an error and keeps old target", func() {
					Expect(executeErr).To(MatchError(shared.OrganizationNotFoundError{Name: "some-org"}))

					Expect(fakeConfig.SetOrganizationInformationCallCount()).To(Equal(0))
					Expect(fakeConfig.SetSpaceInformationCallCount()).To(Equal(0))
				})
			})
		})
	})
})