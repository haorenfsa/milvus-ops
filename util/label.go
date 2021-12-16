package util

func GetComponentLabelByManager(manager string) string {
	switch manager {
	case "helm":
		return "component"
	case "operator":
		return "app.kubernetes.io/component"
	default:
		return ""
	}
}
