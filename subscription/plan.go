package subscription

var (
	PlanFree    Plan = "free"
	PlanPremium Plan = "premium"
)

type Plan string

func IsPremium(plan string) bool {
	return plan == string(PlanPremium)
}
