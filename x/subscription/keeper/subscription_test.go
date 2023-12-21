package keeper_test

import (
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/lavanet/lava/testutil/common"
	keepertest "github.com/lavanet/lava/testutil/keeper"
	"github.com/lavanet/lava/utils"
	"github.com/lavanet/lava/utils/sigs"
	pairingtypes "github.com/lavanet/lava/x/pairing/types"
	planstypes "github.com/lavanet/lava/x/plans/types"
	projectstypes "github.com/lavanet/lava/x/projects/types"
	"github.com/lavanet/lava/x/subscription/types"
	"github.com/stretchr/testify/require"
)

type tester struct {
	common.Tester
}

func newTester(t *testing.T) *tester {
	ts := &tester{Tester: *common.NewTester(t)}
	ts.AddPlan("free", common.CreateMockPlan())
	return ts
}

func (ts *tester) getSubscription(consumer string) (types.Subscription, bool) {
	sub, err := ts.QuerySubscriptionCurrent(consumer)
	require.Nil(ts.T, err)
	if sub.Sub == nil {
		return types.Subscription{}, false
	}
	return *sub.Sub, true
}

func getSubscriptionAndFailTestIfNotFound(t *testing.T, ts *tester, consumer string) types.Subscription {
	sub, found := ts.getSubscription(consumer)
	require.True(t, found)
	require.NotNil(t, sub)
	return sub
}

func getProjectAndFailTestIfNotFound(t *testing.T, ts *tester, consumer string, block uint64) projectstypes.Project {
	project, err := ts.GetProjectForDeveloper(consumer, block)
	require.NoError(t, err)
	require.NotNil(t, project)
	return project
}

func TestCreateSubscription(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(2, 0, 4) // 2 sub, 0 adm, 4 dev

	_, sub1Addr := ts.Account("sub1")
	_, sub2Addr := ts.Account("sub2")
	_, dev1Addr := ts.Account("dev1")
	_, dev2Addr := ts.Account("dev2")
	_, dev3Addr := ts.Account("dev3")

	consumers := []string{dev1Addr, dev2Addr, dev3Addr, "invalid"}
	creators := []string{sub1Addr, sub2Addr, "invalid"}

	var plans []planstypes.Plan
	for i := 0; i < 3; i++ {
		plan := ts.Plan("free")
		plan.Index += strconv.Itoa(i + 1)
		plan.Block = ts.BlockHeight()
		err := ts.TxProposalAddPlans(plan)
		require.NoError(t, err)
		plans = append(plans, plan)
	}

	// delete one plan, and advance to next epoch to take effect
	err := ts.TxProposalDelPlans(plans[2].Index)
	require.NoError(t, err)

	ts.AdvanceEpoch()

	template := []struct {
		name      string
		index     string
		creator   int
		consumers []int
		duration  int
		success   bool
	}{
		{
			name:      "create subscriptions",
			index:     plans[0].Index,
			creator:   0,
			consumers: []int{0, 1},
			duration:  1,
			success:   true,
		},
		{
			name:      "invalid creator",
			index:     plans[0].Index,
			creator:   2,
			consumers: []int{2},
			duration:  1,
			success:   false,
		},
		{
			name:      "invalid consumer",
			index:     plans[0].Index,
			creator:   0,
			consumers: []int{3},
			duration:  1,
			success:   false,
		},
		{
			name:      "duration too long",
			index:     plans[0].Index,
			creator:   0,
			consumers: []int{2},
			duration:  13,
			success:   false,
		},
		{
			name:      "insufficient funds",
			index:     plans[0].Index,
			creator:   1,
			consumers: []int{2},
			duration:  1,
			success:   false,
		},
		{
			name:      "invalid plan",
			index:     "",
			creator:   0,
			consumers: []int{2},
			duration:  1,
			success:   false,
		},
		{
			name:      "unknown plan",
			index:     "no-such-plan",
			creator:   0,
			consumers: []int{2},
			duration:  1,
			success:   false,
		},
		{
			name:      "deleted plan",
			index:     plans[2].Index,
			creator:   0,
			consumers: []int{2},
			duration:  1,
			success:   false,
		},
		{
			name:      "upgrade subscription",
			index:     plans[1].Index,
			creator:   0,
			consumers: []int{0},
			duration:  1,
			success:   true,
		},
	}

	for _, tt := range template {
		for _, consumer := range tt.consumers {
			t.Run(tt.name, func(t *testing.T) {
				sub := types.Subscription{
					Creator:   creators[tt.creator],
					Consumer:  consumers[consumer],
					PlanIndex: tt.index,
				}

				_, err := ts.TxSubscriptionBuy(sub.Creator, sub.Consumer, sub.PlanIndex, tt.duration, false)
				if tt.success {
					require.Nil(t, err, tt.name)
					_, found := ts.getSubscription(sub.Consumer)
					require.True(t, found, tt.name)
				} else {
					require.NotNil(t, err, tt.name)
				}
			})
		}
	}
}

func TestSubscriptionExpiration(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 2 sub, 0 adm, 0 dev

	_, sub1Addr := ts.Account("sub1")
	plan := ts.Plan("free")

	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 1, false)
	require.NoError(t, err)
	_, found := ts.getSubscription(sub1Addr)
	require.True(t, found)

	// advance 1 month + epoch, subscription should expire
	ts.AdvanceMonths(1)
	ts.AdvanceEpoch()

	_, found = ts.getSubscription(sub1Addr)
	require.False(t, found)
}

func TestRenewSubscription(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	_, sub1Addr := ts.Account("sub1")
	plan := ts.Plan("free")

	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 6, false)
	require.NoError(t, err)
	_, found := ts.getSubscription(sub1Addr)
	require.True(t, found)

	// fast-forward three months
	ts.AdvanceMonths(3).AdvanceEpoch()
	sub, found := ts.getSubscription(sub1Addr)
	require.True(t, found)
	require.Equal(t, uint64(3), sub.DurationLeft)

	// with 3 months duration left, asking for 12 more should fail
	_, err = ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 12, false)
	require.NotNil(t, err)

	// but 9 additional month (even 10, the extra month extension below)
	_, err = ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 9, false)
	require.NoError(t, err)
	sub, found = ts.getSubscription(sub1Addr)
	require.True(t, found)

	require.Equal(t, uint64(12), sub.DurationLeft)
	require.Equal(t, uint64(9), sub.DurationBought)

	// edit the subscription's plan (allow more CU)
	cuPerEpoch := plan.PlanPolicy.EpochCuLimit
	plan.PlanPolicy.EpochCuLimit += 100
	plan.Price.Amount = plan.Price.Amount.MulRaw(2)

	err = keepertest.SimulatePlansAddProposal(ts.Ctx, ts.Keepers.Plans, []planstypes.Plan{plan}, false)
	require.NoError(t, err)

	// try extending the subscription (we could extend with 1 more month,
	// but since the subscription's plan changed and its new price is increased
	// by more than 5% , the extension should fail)
	_, err = ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 1, false)
	require.NotNil(t, err)
	require.Equal(t, uint64(12), sub.DurationLeft)
	require.Equal(t, uint64(9), sub.DurationBought)

	// get the subscription's plan and make sure it uses the old plan
	plan, found = ts.FindPlan(sub.PlanIndex, sub.PlanBlock)
	require.True(t, found)
	require.Equal(t, cuPerEpoch, plan.PlanPolicy.EpochCuLimit)

	// delete the plan, and try to renew the subscription again
	err = ts.TxProposalDelPlans(plan.Index)
	require.NoError(t, err)

	ts.AdvanceEpoch()

	// fast-forward another month, renewal should fail
	ts.AdvanceMonths(1).AdvanceEpoch()
	_, found = ts.getSubscription(sub1Addr)
	require.True(t, found)
	_, err = ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 10, false)
	require.NotNil(t, err)
}

func TestSubscriptionAdminProject(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	_, sub1Addr := ts.Account("sub1")
	plan := ts.Plan("free")

	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 1, false)
	require.NoError(t, err)

	// a newly created subscription is expected to have one default project,
	// with the subscription address as its developer key
	_, err = ts.GetProjectDeveloperData(sub1Addr, ts.BlockHeight())
	require.NoError(t, err)
}

func TestMonthlyRechargeCU(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 1, 1) // 1 sub, 1 adm, 1 dev

	_, sub1Addr := ts.Account("sub1")
	_, adm1Addr := ts.Account("adm1")
	_, dev1Addr := ts.Account("dev1")
	plan := ts.Plan("free")

	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 3, false)
	require.NoError(t, err)

	// add another project under the subscription
	projectData := projectstypes.ProjectData{
		Name:    "another_project",
		Enabled: true,
		ProjectKeys: []projectstypes.ProjectKey{
			projectstypes.ProjectDeveloperKey(dev1Addr),
		},
		Policy: &planstypes.Policy{
			GeolocationProfile: 1,
			TotalCuLimit:       1000,
			EpochCuLimit:       100,
			MaxProvidersToPair: 3,
		},
	}
	err = ts.TxSubscriptionAddProject(sub1Addr, projectData)
	require.NoError(t, err)

	template := []struct {
		name             string
		subscription     string
		developer        string
		usedCuPerProject uint64 // total sub CU is 1000; each project uses 500
	}{
		{"default project", sub1Addr, sub1Addr, 500},
		{"second project (non-default)", sub1Addr, dev1Addr, 500},
	}
	for ti, tt := range template {
		t.Run(tt.name, func(t *testing.T) {
			block1 := ts.BlockHeight()
			ts.AdvanceEpoch()

			// charge the subscription
			_, err = ts.Keepers.Subscription.ChargeComputeUnitsToSubscription(
				ts.Ctx, tt.subscription, block1, tt.usedCuPerProject)
			require.NoError(t, err)

			// verify the CU charge of the subscription is updated correctly
			sub, found := ts.getSubscription(tt.subscription)
			require.True(t, found)
			require.Equal(t, sub.MonthCuLeft, sub.MonthCuTotal-tt.usedCuPerProject)

			// charge the project
			proj, err := ts.GetProjectForDeveloper(tt.developer, block1)
			require.NoError(t, err)
			err = ts.Keepers.Projects.ChargeComputeUnitsToProject(
				ts.Ctx, proj, block1, tt.usedCuPerProject)
			require.NoError(t, err)

			// verify that project used the CU
			proj, err = ts.GetProjectForDeveloper(tt.developer, block1)
			require.NoError(t, err)
			require.Equal(t, tt.usedCuPerProject, proj.UsedCu)

			block2 := ts.BlockHeight()

			// force fixation entry (by adding project key)
			projKey := []projectstypes.ProjectKey{projectstypes.ProjectAdminKey(adm1Addr)}
			ts.Keepers.Projects.AddKeysToProject(ts.Ctx, projectstypes.ADMIN_PROJECT_NAME, tt.developer, projKey)

			// fast-forward one month
			ts.AdvanceMonths(1).AdvanceEpoch()
			sub, found = ts.getSubscription(sub1Addr)
			require.True(t, found)
			require.Equal(t, sub.DurationBought-uint64(ti+1), sub.DurationLeft)

			block3 := ts.BlockHeight()

			// check that subscription and project have renewed CUs, and that
			// the project created a snapshot for last month
			sub, found = ts.getSubscription(tt.subscription)
			require.True(t, found)
			require.Equal(t, sub.MonthCuLeft, sub.MonthCuTotal)

			proj, err = ts.GetProjectForDeveloper(tt.developer, block1)
			require.NoError(t, err)
			require.Equal(t, tt.usedCuPerProject, proj.UsedCu)
			proj, err = ts.GetProjectForDeveloper(tt.developer, block2)
			require.NoError(t, err)
			require.Equal(t, tt.usedCuPerProject, proj.UsedCu)
			proj, err = ts.GetProjectForDeveloper(tt.developer, block3)
			require.NoError(t, err)
			require.Equal(t, uint64(0), proj.UsedCu)
		})
	}
}

func TestExpiryTime(t *testing.T) {
	ts := newTester(t)

	template := []struct {
		now    [3]int // year, month, day
		res    [3]int // year, month, day
		months int
	}{
		// monthly
		{[3]int{2000, 3, 1}, [3]int{2000, 4, 1}, 1},
		{[3]int{2000, 3, 30}, [3]int{2000, 4, 28}, 1},
		{[3]int{2000, 3, 31}, [3]int{2000, 4, 28}, 1},
		{[3]int{2000, 2, 1}, [3]int{2000, 3, 1}, 1},
		{[3]int{2000, 2, 28}, [3]int{2000, 3, 28}, 1},
		{[3]int{2001, 2, 28}, [3]int{2001, 3, 28}, 1},
		{[3]int{2000, 2, 29}, [3]int{2000, 3, 28}, 1},
		{[3]int{2000, 1, 28}, [3]int{2000, 2, 28}, 1},
		{[3]int{2001, 1, 28}, [3]int{2001, 2, 28}, 1},
		{[3]int{2000, 1, 29}, [3]int{2000, 2, 28}, 1},
		{[3]int{2001, 1, 29}, [3]int{2001, 2, 28}, 1},
		{[3]int{2000, 1, 30}, [3]int{2000, 2, 28}, 1},
		{[3]int{2001, 1, 30}, [3]int{2001, 2, 28}, 1},
		{[3]int{2000, 1, 31}, [3]int{2000, 2, 28}, 1},
		{[3]int{2001, 1, 31}, [3]int{2001, 2, 28}, 1},
		{[3]int{2001, 12, 31}, [3]int{2002, 1, 28}, 1},
		// duration > 1
		{[3]int{2000, 3, 1}, [3]int{2000, 4, 1}, 2},
		{[3]int{2000, 3, 1}, [3]int{2000, 4, 1}, 6},
		{[3]int{2000, 3, 1}, [3]int{2000, 4, 1}, 12},
	}

	plan := ts.Plan("free")

	for _, tt := range template {
		now := time.Date(tt.now[0], time.Month(tt.now[1]), tt.now[2], 12, 0, 0, 0, time.UTC)

		t.Run(now.Format("2006-01-02"), func(t *testing.T) {
			// new account per attempt
			_, sub1Addr := ts.AddAccount("tmp", 0, 10000)

			delta := now.Sub(ts.BlockTime())
			ts.AdvanceBlock(delta)

			_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, tt.months, false)
			require.NoError(t, err)

			sub, found := ts.getSubscription(sub1Addr)
			require.True(t, found)
			require.Equal(t, uint64(tt.months), sub.DurationBought)

			// will expire and remove
			ts.AdvanceMonths(tt.months).AdvanceEpoch()
			ts.AdvanceBlockUntilStale()
		})
	}
}

func TestSubscriptionExpire(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	sub1Acct, sub1Addr := ts.Account("sub1")
	plan := ts.Plan("free")

	coins := common.NewCoins(ts.TokenDenom(), 10000)
	ts.Keepers.BankKeeper.SetBalance(ts.Ctx, sub1Acct.Addr, coins)

	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 1, false)
	require.NoError(t, err)

	block := ts.BlockHeight()

	_, found := ts.getSubscription(sub1Addr)
	require.True(t, found)

	_, err = ts.Keepers.Subscription.ChargeComputeUnitsToSubscription(
		ts.Ctx, sub1Addr, block, 10)
	require.NoError(t, err)

	// fast-forward one month
	ts.AdvanceMonths(1).AdvanceEpoch()

	// subscription no longer searchable, but can still charge for previous usage
	_, found = ts.getSubscription(sub1Addr)
	require.False(t, found)

	_, err = ts.Keepers.Subscription.ChargeComputeUnitsToSubscription(
		ts.Ctx, sub1Addr, block, 10)
	require.NoError(t, err)

	ts.AdvanceBlockUntilStale()

	// subscription no longer charge-able for previous usage
	_, err = ts.Keepers.Subscription.ChargeComputeUnitsToSubscription(
		ts.Ctx, sub1Addr, block, 10)
	require.NotNil(t, err)
}

func TestPrice(t *testing.T) {
	ts := newTester(t)

	template := []struct {
		name     string
		duration int
		discount uint64
		price    int64
		cost     int64
	}{
		{"1 month", 1, 0, 100, 100},
		{"2 months", 2, 0, 100, 200},
		{"11 months", 11, 0, 100, 1100},
		{"yearly without discount", 12, 0, 100, 1200},
		{"yearly with discount", 12, 25, 100, 900},
	}

	for _, tt := range template {
		t.Run(tt.name, func(t *testing.T) {
			// new account per attempt
			sub1Acct, sub1Addr := ts.AddAccount("tmp", 0, 10000)

			plan := ts.Plan("free")
			plan.AnnualDiscountPercentage = tt.discount
			plan.Price = common.NewCoin(ts.TokenDenom(), tt.price)
			err := ts.TxProposalAddPlans(plan)
			require.NoError(t, err)

			_, err = ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, tt.duration, false)
			require.NoError(t, err)

			_, found := ts.getSubscription(sub1Addr)
			require.True(t, found)

			balance := ts.GetBalance(sub1Acct.Addr)
			require.Equal(t, balance, 10000-tt.cost)

			// will expire and remove
			ts.AdvanceMonths(tt.duration)
		})
	}
}

func TestAddProjectToSubscription(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 1, 1) // 1 sub, 0 adm, 2 dev

	_, sub1Addr := ts.Account("sub1")
	_, adm1Addr := ts.Account("adm1")
	_, dev1Addr := ts.Account("dev1")
	plan := ts.Plan("free")

	_, err := ts.TxSubscriptionBuy(sub1Addr, dev1Addr, plan.Index, 1, false)
	require.NoError(t, err)

	template := []struct {
		name         string
		subscription string
		anotherAdmin string
		projectName  string
		success      bool
	}{
		{"project admin = regular account", dev1Addr, adm1Addr, "test1", true},
		{"project admin = subscription payer account", dev1Addr, sub1Addr, "test2", true},
		{"bad subscription account (regular account)", adm1Addr, dev1Addr, "test3", false},
		{"bad subscription account (subscription payer account)", sub1Addr, dev1Addr, "test4", false},
		{"bad projectName (duplicate)", dev1Addr, adm1Addr, "invalid:name", false},
	}

	for _, tt := range template {
		t.Run(tt.name, func(t *testing.T) {
			projectData := projectstypes.ProjectData{
				Name:    tt.projectName,
				Enabled: true,
				ProjectKeys: []projectstypes.ProjectKey{
					projectstypes.ProjectAdminKey(tt.anotherAdmin),
				},
			}
			projectID := projectstypes.ProjectIndex(tt.subscription, tt.projectName)
			err = ts.TxSubscriptionAddProject(tt.subscription, projectData)
			if tt.success {
				require.NoError(t, err)
				proj, err := ts.GetProjectForBlock(projectID, ts.BlockHeight())
				require.NoError(t, err)
				require.Equal(t, tt.subscription, proj.Subscription)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestGetProjectsForSubscription(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(2, 0, 0) // 2 sub, 0 adm, 0 dev

	_, sub1Addr := ts.Account("sub1")
	_, sub2Addr := ts.Account("sub2")
	plan := ts.Plan("free")

	// buy two subscriptions
	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 1, false)
	require.NoError(t, err)
	_, err = ts.TxSubscriptionBuy(sub2Addr, sub2Addr, plan.Index, 1, false)
	require.NoError(t, err)

	// add two projects to the first subscription
	projData1 := projectstypes.ProjectData{
		Name:    "proj1",
		Enabled: true,
		Policy:  &plan.PlanPolicy,
	}
	err = ts.TxSubscriptionAddProject(sub1Addr, projData1)
	require.NoError(t, err)

	projData2 := projectstypes.ProjectData{
		Name:    "proj2",
		Enabled: false,
		Policy:  &plan.PlanPolicy,
	}
	err = ts.TxSubscriptionAddProject(sub1Addr, projData2)
	require.NoError(t, err)

	res1, err := ts.QuerySubscriptionListProjects(sub1Addr)
	require.NoError(t, err)

	res2, err := ts.QuerySubscriptionListProjects(sub2Addr)
	require.NoError(t, err)

	// number of projects +1 to account for auto-generated admin project
	require.Equal(t, 3, len(res1.Projects))
	require.Equal(t, 1, len(res2.Projects))

	err = ts.TxSubscriptionDelProject(sub1Addr, projData2.Name)
	require.NoError(t, err)
}

func TestAddDelProjectForSubscription(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	_, sub1Addr := ts.Account("sub1")
	plan := ts.Plan("free")

	// buy subscription and add project
	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 1, false)
	require.NoError(t, err)

	projData := projectstypes.ProjectData{
		Name:    "proj",
		Enabled: true,
		Policy:  &plan.PlanPolicy,
	}
	err = ts.TxSubscriptionAddProject(sub1Addr, projData)
	require.NoError(t, err)

	ts.AdvanceEpoch()

	res, err := ts.QuerySubscriptionListProjects(sub1Addr)
	require.NoError(t, err)
	require.Equal(t, 2, len(res.Projects))

	// del project to the subscription
	err = ts.TxSubscriptionDelProject(sub1Addr, projData.Name)
	require.NoError(t, err)

	ts.AdvanceEpoch()

	res, err = ts.QuerySubscriptionListProjects(sub1Addr)
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Projects))
}

func TestDelProjectEndSubscription(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	_, sub1Addr := ts.Account("sub1")
	plan := ts.Plan("free")

	// buy subscription
	_, err := ts.TxSubscriptionBuy(sub1Addr, sub1Addr, plan.Index, 1, false)
	require.NoError(t, err)

	// time of buy subscription
	start := ts.BlockTime()

	// add project to the subscription
	projData := projectstypes.ProjectData{
		Name:    "proj",
		Enabled: true,
		Policy:  &plan.PlanPolicy,
	}
	err = ts.TxSubscriptionAddProject(sub1Addr, projData)
	require.NoError(t, err)

	ts.AdvanceEpoch()

	res, err := ts.QuerySubscriptionListProjects(sub1Addr)
	require.NoError(t, err)
	require.Equal(t, 2, len(res.Projects))

	// advance time to just before subscription expiry, so project deletion
	// and the subsequent expiry will occur in the same epoch
	ts.AdvanceMonthsFrom(start, 1)

	// del project to the subscription
	err = ts.TxSubscriptionDelProject(sub1Addr, projData.Name)
	require.NoError(t, err)

	// expire subscription (by advancing an epoch, we are close enough to expiry)
	ts.AdvanceEpoch()

	_, err = ts.QuerySubscriptionListProjects(sub1Addr)
	require.NotNil(t, err)

	// should not panic
	ts.AdvanceBlocks(2 * ts.BlocksToSave())
}

// TestDurationTotal tests that the total duration of the subscription is updated correctly
func TestDurationTotal(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev
	months := 12
	plan := ts.Plan("free")

	_, subAddr := ts.Account("sub1")
	_, err := ts.TxSubscriptionBuy(subAddr, subAddr, plan.Index, months, false)
	require.NoError(t, err)

	for i := 0; i < months-1; i++ {
		subRes, err := ts.QuerySubscriptionCurrent(subAddr)
		sub := subRes.Sub
		require.NoError(t, err)
		require.Equal(t, uint64(i), sub.DurationTotal)
		ts.AdvanceMonths(1)
		ts.AdvanceEpoch()
	}

	// buy extra 4 months and check duration total continues from last count
	subRes, err := ts.QuerySubscriptionCurrent(subAddr)
	require.NoError(t, err)
	durationSoFar := subRes.Sub.DurationTotal

	extraMonths := 4
	_, err = ts.TxSubscriptionBuy(subAddr, subAddr, plan.Index, extraMonths, false)
	require.NoError(t, err)

	for i := 0; i < extraMonths; i++ {
		subRes, err := ts.QuerySubscriptionCurrent(subAddr)
		sub := subRes.Sub
		require.NoError(t, err)
		require.Equal(t, uint64(i)+durationSoFar, sub.DurationTotal)
		ts.AdvanceMonths(1)
		ts.AdvanceEpoch()
	}

	// expire subscription and buy a new one. verify duration total starts from scratch
	ts.AdvanceMonths(1)
	ts.AdvanceEpoch()
	subRes, err = ts.QuerySubscriptionCurrent(subAddr)
	require.NoError(t, err)
	require.Nil(t, subRes.Sub)

	_, err = ts.TxSubscriptionBuy(subAddr, subAddr, plan.Index, extraMonths, false)
	require.NoError(t, err)
	subRes, err = ts.QuerySubscriptionCurrent(subAddr)
	require.NoError(t, err)
	require.Equal(t, uint64(0), subRes.Sub.DurationTotal)
}

// TestSubAutoRenewal is a happy flow test for subscription auto-renewal
// checks that the two methods for enabling auto renewal works
// verifies that subs with auto-renewal enabled get renewed automatically
func TestSubAutoRenewal(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(3, 0, 0) // 2 sub, 0 adm, 0 dev

	plan := ts.Plan("free")
	_, subAddr1 := ts.Account("sub1")
	_, subAddr2 := ts.Account("sub2")
	_, subAddr3 := ts.Account("sub3")

	// buy two subscriptions with enabled auto-renewal in two different ways
	// and one with disabled auto-renewal.
	// verify the auto-renewal flag is true in the first two subs
	_, err := ts.TxSubscriptionBuy(subAddr1, subAddr1, plan.Index, 1, true)
	require.NoError(t, err)
	_, err = ts.TxSubscriptionBuy(subAddr2, subAddr2, plan.Index, 1, false)
	require.NoError(t, err)
	err = ts.TxSubscriptionAutoRenewal(subAddr2, true)
	require.NoError(t, err)
	_, err = ts.TxSubscriptionBuy(subAddr3, subAddr3, plan.Index, 1, false)
	require.NoError(t, err)

	sub1, found := ts.getSubscription(subAddr1)
	require.True(t, found)
	require.True(t, sub1.AutoRenewal)
	sub2, found := ts.getSubscription(subAddr2)
	require.True(t, found)
	require.True(t, sub2.AutoRenewal)
	sub3, found := ts.getSubscription(subAddr3)
	require.True(t, found)
	require.False(t, sub3.AutoRenewal)

	// advance a couple of months to expire and automatically
	// extend all subscriptions. verify that sub1 and sub2 can
	// still be found and their duration left is always 1
	for i := 0; i < 5; i++ {
		ts.AdvanceMonths(1).AdvanceEpoch()

		newSub1, found := ts.getSubscription(subAddr1)
		require.True(t, found)
		require.Equal(t, uint64(1), newSub1.DurationLeft)
		newSub2, found := ts.getSubscription(subAddr2)
		require.True(t, found)
		require.Equal(t, uint64(1), newSub2.DurationLeft)
		_, found = ts.getSubscription(subAddr3)
		require.False(t, found)
	}
}

// TestSubRenewalFailHighPlanPrice checks that auto-renewal fails when the
// original subscription's plan price increased by more than 5%
func TestSubRenewalFailHighPlanPrice(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	_, subAddr1 := ts.Account("sub1")
	plan := ts.Plan("free")

	_, err := ts.TxSubscriptionBuy(subAddr1, subAddr1, plan.Index, 1, true)
	require.NoError(t, err)
	_, found := ts.getSubscription(subAddr1)
	require.True(t, found)

	// edit the subscription's plan (increase the price by 6% and change the policy (shouldn't matter))
	plan.PlanPolicy.EpochCuLimit += 100
	plan.Price.Amount = plan.Price.Amount.MulRaw(106).QuoRaw(100)

	ts.AdvanceEpoch() // advance epoch so the new plan will be appended as a new entry
	err = keepertest.SimulatePlansAddProposal(ts.Ctx, ts.Keepers.Plans, []planstypes.Plan{plan}, false)
	require.NoError(t, err)

	// advance month to make the subscription expire
	ts.AdvanceMonths(1).AdvanceEpoch()

	// the auto-renewal should've failed since the plan price is too high
	// so the subscription should not be found
	_, found = ts.getSubscription(subAddr1)
	require.False(t, found)
}

// TestNextToMonthExpiryQuery checks that the NextToMonthExpiry query works as intended
// scenario - buy 3 subs: 2 at the same time, and one a little after. The query should return the two subs
// then, expire those and expect to get the last one from the query
func TestNextToMonthExpiryQuery(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(3, 0, 0) // 1 sub, 0 adm, 0 dev
	months := 1
	plan := ts.Plan("free")

	_, sub1 := ts.Account("sub1")
	_, sub2 := ts.Account("sub2")
	_, sub3 := ts.Account("sub3")

	// buy 3 subs - 2 at the same time and one a second later
	_, err := ts.TxSubscriptionBuy(sub1, sub1, plan.Index, months, false)
	require.NoError(t, err)
	_, err = ts.TxSubscriptionBuy(sub2, sub2, plan.Index, months, false)
	require.NoError(t, err)
	sub1Obj, found := ts.getSubscription(sub1)
	require.True(t, found)

	ts.AdvanceBlock(time.Second)
	_, err = ts.TxSubscriptionBuy(sub3, sub3, plan.Index, months, false)
	require.NoError(t, err)
	sub3Obj, found := ts.getSubscription(sub3)
	require.True(t, found)
	require.Equal(t, sub3Obj.MonthExpiryTime, sub1Obj.MonthExpiryTime+1) // sub3 should expire one second after sub1

	// query - expect subs 1 and 2 in the output
	res, err := ts.QuerySubscriptionNextToMonthExpiry()
	require.NoError(t, err)
	require.Equal(t, 2, len(res.Subscriptions))

	for _, sub := range res.Subscriptions {
		if sub.Consumer != sub1 && sub.Consumer != sub2 {
			require.Fail(t, "resulting subscription are not sub1 or sub2")
		}
		require.Equal(t, sub1Obj.MonthExpiryTime, sub.MonthExpiry)
	}

	// advance month minus 4 seconds
	ts.AdvanceMonths(1).AdvanceBlock(4 * time.Second)
	ts.AdvanceBlock(time.Second) // separate so advanceMonth would trigger

	// query - expect sub 3 in the output
	res, err = ts.QuerySubscriptionNextToMonthExpiry()
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Subscriptions))
	require.Equal(t, sub3, res.Subscriptions[0].Consumer)
	require.Equal(t, sub3Obj.MonthExpiryTime, res.Subscriptions[0].MonthExpiry)

	// advance another second to expire sub3. Expect empty output from the query
	ts.AdvanceBlock(time.Second)
	res, err = ts.QuerySubscriptionNextToMonthExpiry()
	require.NoError(t, err)
	require.Equal(t, 0, len(res.Subscriptions))
}

// TestPlanRemovedWhenSubscriptionExpires checks that if a subscription is expired
// the plan's refcount is decreased.
// in this test, we buy a subscription, update the plan and expire the subscription
// since the old plan version's refcount is decreased, the old plan should be removed
// note that we had to update the plan because latest versions (of every fixation object)
// is never deleted
func TestPlanRemovedWhenSubscriptionExpires(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev
	months := 1
	plan := ts.Plan("free")

	_, sub1 := ts.Account("sub1")

	// buy sub with plan first version
	_, err := ts.TxSubscriptionBuy(sub1, sub1, plan.Index, months, false)
	require.NoError(t, err)
	oldPlanBlock := ts.BlockHeight()

	// update plan
	ts.AdvanceEpoch()
	plan.OveruseRate++
	err = ts.Keepers.Plans.AddPlan(ts.Ctx, plan, false)
	require.NoError(t, err)

	// expire the subscription
	ts.AdvanceMonths(1)
	ts.AdvanceBlock(6)
	res, err := ts.QuerySubscriptionCurrent(sub1)
	require.NoError(t, err)
	require.Nil(t, res.Sub)

	// wait the stale period and check old plan doesn't exist
	ts.AdvanceBlockUntilStale()
	_, found := ts.Keepers.Plans.FindPlan(ts.Ctx, plan.Index, oldPlanBlock)
	require.False(t, found)
}

func TestSubscriptionUpgrade(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	_, consumer := ts.Account("sub1")
	freePlan := ts.Plan("free")

	// Add premium plan
	upgradedPlan := common.CreateMockPlan()
	upgradedPlan.Index = "premium"
	upgradedPlan.Price = freePlan.Price.AddAmount(math.NewInt(100))
	ts.AddPlan(upgradedPlan.Index, upgradedPlan)

	// Buy free plan
	_, err := ts.TxSubscriptionBuy(consumer, consumer, freePlan.Index, 1, false)
	require.NoError(t, err)
	// Verify subscription found inside getSubscription
	getSubscriptionAndFailTestIfNotFound(t, ts, consumer)

	// Charge CU from project so we can differentiate the old project from the new one
	projectCuUsed := uint64(100)
	project := getProjectAndFailTestIfNotFound(t, ts, consumer, ts.BlockHeight())
	err = ts.Keepers.Projects.ChargeComputeUnitsToProject(ts.Ctx, project, ts.BlockHeight(), projectCuUsed)
	require.NoError(t, err)

	// Validate the charge of CU
	project = getProjectAndFailTestIfNotFound(t, ts, consumer, ts.BlockHeight())
	require.Equal(t, projectCuUsed, project.UsedCu)

	ts.AdvanceEpochs(2)
	sub := getSubscriptionAndFailTestIfNotFound(t, ts, consumer)
	currentDurationTotal := sub.DurationTotal

	// Buy premium plan
	_, err = ts.TxSubscriptionBuy(consumer, consumer, upgradedPlan.Index, 1, false)
	require.NoError(t, err)

	nextEpoch := ts.GetNextEpoch()

	// Test that the subscription and project are not changed until next epoch
	for ts.BlockHeight() < nextEpoch {
		sub := getSubscriptionAndFailTestIfNotFound(t, ts, consumer)
		require.Equal(t, freePlan.Index, sub.PlanIndex, "plan should be free until next epoch. Block: %v", ts.BlockHeight())
		require.Equal(t, currentDurationTotal, sub.DurationTotal)

		project = getProjectAndFailTestIfNotFound(t, ts, consumer, ts.BlockHeight())
		require.Equal(t, projectCuUsed, project.UsedCu)

		ts.AdvanceBlock()
	}

	// Test that the subscription is now updated
	sub = getSubscriptionAndFailTestIfNotFound(t, ts, consumer)
	require.Equal(t, upgradedPlan.Index, sub.PlanIndex)

	// Test that the project is now updated
	project = getProjectAndFailTestIfNotFound(t, ts, consumer, ts.BlockHeight())
	require.Equal(t, uint64(0), project.UsedCu)
}

func TestSubscriptionDowngradeFails(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev

	_, consumer := ts.Account("sub1")
	freePlan := ts.Plan("free")

	// Add premium plan
	upgradedPlan := common.CreateMockPlan()
	upgradedPlan.Index = "premium"
	upgradedPlan.Price = freePlan.Price.AddAmount(math.NewInt(100))
	ts.AddPlan(upgradedPlan.Index, upgradedPlan)

	// Buy premium plan
	_, err := ts.TxSubscriptionBuy(consumer, consumer, upgradedPlan.Index, 1, false)
	require.NoError(t, err)
	// Verify subscription found inside getSubscription
	getSubscriptionAndFailTestIfNotFound(t, ts, consumer)

	ts.AdvanceEpochs(2)

	// Buy premium plan
	_, err = ts.TxSubscriptionBuy(consumer, consumer, freePlan.Index, 1, false)
	require.NotNil(t, err)

	ts.AdvanceEpoch()

	sub := getSubscriptionAndFailTestIfNotFound(t, ts, consumer)
	require.Equal(t, upgradedPlan.Index, sub.PlanIndex)
}

func TestSubscriptionCuExhaustAndUpgrade(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev
	consumerAcc, consumerAddr := ts.Account("sub1")

	spec := ts.AddSpec("testSpec", common.CreateMockSpec()).Spec("testSpec")

	// Setup validator and provider
	testBalance := int64(1000000)
	testStake := int64(100000)

	validationAcc, _ := ts.AddAccount(common.VALIDATOR, 0, testBalance)
	ts.TxCreateValidator(validationAcc, math.NewInt(testBalance))

	_, providerAddr := ts.AddAccount(common.PROVIDER, 0, testBalance)
	err := ts.StakeProviderExtra(providerAddr, spec, testStake, nil, 0, "provider")
	require.NoError(t, err)

	// Trigger changes
	ts.AdvanceEpoch()

	freePlan := ts.Plan("free")

	// Add premium plan
	premiumPlan := common.CreateMockPlan()
	premiumPlan.Index = "premium"
	premiumPlan.Price = freePlan.Price.AddAmount(math.NewInt(100))
	ts.AddPlan(premiumPlan.Index, premiumPlan)

	// Add premium-plus plan
	premiumPlusPlan := common.CreateMockPlan()
	premiumPlusPlan.Index = "premium-plus"
	premiumPlusPlan.Price = premiumPlan.Price.AddAmount(math.NewInt(100))
	ts.AddPlan(premiumPlusPlan.Index, premiumPlusPlan)

	// Buy free plan
	_, err = ts.TxSubscriptionBuy(consumerAddr, consumerAddr, freePlan.Index, 3, false)
	require.NoError(t, err)

	// Verify subscription found inside getSubscription
	getSubscriptionAndFailTestIfNotFound(t, ts, consumerAddr)

	// Send relay
	sessionId := uint64(1)
	relayNum := uint64(1)
	sendRelayPayment := func() {
		relaySession := &pairingtypes.RelaySession{
			Provider:    providerAddr,
			ContentHash: []byte(spec.ApiCollections[0].Apis[0].Name),
			SessionId:   sessionId,
			SpecId:      spec.Index,
			CuSum:       1000,
			Epoch:       int64(ts.EpochStart(ts.BlockHeight())),
			RelayNum:    relayNum,
		}

		sig, err := sigs.Sign(consumerAcc.SK, *relaySession)
		require.Nil(ts.T, err)
		relaySession.Sig = sig

		_, err = ts.TxPairingRelayPayment(providerAddr, relaySession)
		require.NoError(t, err)

		sessionId++
		relayNum++
	}

	// Send relay under the free subscription
	sendRelayPayment()

	// Buy premium plan
	_, err = ts.TxSubscriptionBuy(consumerAddr, consumerAddr, premiumPlan.Index, 1, false)
	require.NoError(t, err)

	// Trigger new subscription
	ts.AdvanceEpoch()

	// Test that the subscription is now updated
	sub := getSubscriptionAndFailTestIfNotFound(t, ts, consumerAddr)
	require.Equal(t, premiumPlan.Index, sub.PlanIndex)

	// Test that the project is now updated
	project := getProjectAndFailTestIfNotFound(t, ts, consumerAddr, ts.BlockHeight())
	require.Equal(t, uint64(0), project.UsedCu)

	// Send relay under the premium subscription
	sendRelayPayment()

	// Buy premium-plus plan
	_, err = ts.TxSubscriptionBuy(consumerAddr, consumerAddr, premiumPlusPlan.Index, 1, false)
	require.NoError(t, err)

	// Trigger new subscription
	ts.AdvanceEpoch()

	// Test that the subscription is now updated
	sub = getSubscriptionAndFailTestIfNotFound(t, ts, consumerAddr)
	require.Equal(t, premiumPlusPlan.Index, sub.PlanIndex)

	// Test that the project is now updated
	project = getProjectAndFailTestIfNotFound(t, ts, consumerAddr, ts.BlockHeight())
	require.Equal(t, uint64(0), project.UsedCu)

	// Send relay under the premium-plus subscription
	sendRelayPayment()

	// Advance month + blocksToSave + 1 to trigger the provider monthly payment
	ts.AdvanceMonths(1)
	ts.AdvanceBlocks(ts.BlocksToSave() + 1)

	// Query provider's rewards
	rewards, err := ts.QueryDualstakingDelegatorRewards(providerAddr, providerAddr, spec.Index)
	require.NoError(t, err)
	require.Len(t, rewards.Rewards, 1)
	reward := rewards.Rewards[0]

	// Verify that provider got rewarded for both subscriptions
	expectedPrice := freePlan.Price.AddAmount(premiumPlan.Price.Amount).AddAmount(premiumPlusPlan.Price.Amount)
	require.Equal(t, expectedPrice, reward.Amount)
}

func TestSubscriptionUpgradeAffectsTimer(t *testing.T) {
	ts := newTester(t)
	ts.SetupAccounts(1, 0, 0) // 1 sub, 0 adm, 0 dev
	_, consumerAddr := ts.Account("sub1")

	freePlan := ts.Plan("free")

	// Add premium plan
	premiumPlan := common.CreateMockPlan()
	premiumPlan.Index = "premium"
	premiumPlan.Price = freePlan.Price.AddAmount(math.NewInt(100))
	ts.AddPlan(premiumPlan.Index, premiumPlan)

	// Add premium-plus plan
	premiumPlusPlan := common.CreateMockPlan()
	premiumPlusPlan.Index = "premium-plus"
	premiumPlusPlan.Price = premiumPlan.Price.AddAmount(math.NewInt(100))
	ts.AddPlan(premiumPlusPlan.Index, premiumPlusPlan)

	// Buy free plan
	_, err := ts.TxSubscriptionBuy(consumerAddr, consumerAddr, freePlan.Index, 3, false)
	require.NoError(t, err)

	// Verify timer for free plan expiration
	verifyTimerStore := func() {
		subTimers := ts.Keepers.Subscription.ExportSubscriptionsTimers(ts.Ctx).TimeEntries
		require.NoError(t, err)
		require.NotNil(t, subTimers)
		require.Len(t, subTimers, 1)
		require.Equal(t, consumerAddr, subTimers[0].Key)
		require.Equal(t, uint64(utils.NextMonth(ts.BlockTime()).UTC().Unix()), subTimers[0].Value)
	}

	verifyTimerStore()

	ts.AdvanceBlock()

	// Buy premium plan
	_, err = ts.TxSubscriptionBuy(consumerAddr, consumerAddr, premiumPlan.Index, 1, false)
	require.NoError(t, err)

	verifyTimerStore()

	ts.AdvanceBlock()

	// Buy premium-plus plan
	_, err = ts.TxSubscriptionBuy(consumerAddr, consumerAddr, premiumPlusPlan.Index, 1, false)
	require.NoError(t, err)

	verifyTimerStore()
}
