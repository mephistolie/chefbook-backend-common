package subscription

var (
	PlanFree    = "free"
	PlanPremium = "premium"
	PlanMaximum = "maximum"
)

func IsPremium(plan string) bool {
	return plan == PlanPremium || plan == string(PlanMaximum)
}

func IsMaximum(plan string) bool {
	return plan == string(PlanMaximum)
}
