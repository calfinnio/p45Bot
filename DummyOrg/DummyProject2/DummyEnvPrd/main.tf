module "awesome_sauce1" {
    source = "../../modules"
    environment = "dev"
    projectOwner = ["tim.alexander@calfinn.io"]
    budgetAlerts = ["Andy.Dufresne@tailspin.com", "oscar.wallace@tailspin.com"]
    
}