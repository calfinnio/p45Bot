module "awesome_sauce2" {
  source       = "../../modules"
  environment  = "prd"
  projectOwner = "bob.smith@tailspin.com"
  budgetAlerts = [
    "Ben.Wyatt@tailspin.com",
    "Andy.Dufresne@tailspin.com",
    "oscar.wallace@tailspin.com"
  ]

}