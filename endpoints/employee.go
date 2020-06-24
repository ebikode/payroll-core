package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	act "github.com/ebikode/payroll-core/domain/activity_log"
	emp "github.com/ebikode/payroll-core/domain/employee"
	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
	ut "github.com/ebikode/payroll-core/utils"
	"github.com/go-chi/chi"
)

// GetEmployeeEndpoint fetch Authenticated employee account
func GetEmployeeEndpoint(emps emp.EmployeeService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenData := r.Context().Value("tokenData").(*md.EmployeeTokenData)
		employeeID := string(tokenData.EmployeeID)

		employee := emps.GetEmployee(employeeID)
		resp := ut.Message(true, "")
		resp["employee"] = employee
		ut.Respond(w, r, resp)
	}
}

// GetEmployeesEndpoint Admin Endpoint for getting employees
func GetEmployeesEndpoint(emps emp.EmployeeService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		page, limit := ut.PaginationParams(r)

		employees := emps.GetAllEmployees(page, limit)

		var nextPage int
		if len(employees) == limit {
			nextPage = page + 1
		}

		resp := ut.Message(true, "")
		resp["current_page"] = page
		resp["next_page"] = nextPage
		resp["limit"] = limit
		resp["employees"] = employees
		ut.Respond(w, r, resp)
	}
}

// CreateEmployeeEndpoint An endpoint for creating employees new account
func CreateEmployeeEndpoint(emps emp.EmployeeService, acs act.ActivityLogService, clientURL, sendGridKey string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get Admin Token Data
		tokenData := r.Context().Value("tokenData").(*md.AdminTokenData)
		adminID := string(tokenData.AdminID)

		employee := md.Employee{}

		err := json.NewDecoder(r.Body).Decode(&employee)

		tParam := tr.TParam{
			Key:          "error.request_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		if err != nil {
			// Respond with an error translated

			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		checkEmployee := emps.GetEmployeeByEmail(employee.Email)

		// if employee already exist send pincode to the employee for verification/Authentication
		if checkEmployee != nil {

			tParam = tr.TParam{
				Key:          "error.email_already_exist",
				TemplateData: nil,
				PluralCount:  nil,
			}

			resp := ut.Message(true, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		// username := strings(employee.FirstName).LowerCase() + "_" + strings(employee.LastName).LowerCase()

		// REMOVE THIS AFTER DEMO
		employee.Password = "EMPASSWORD2020"
		employee.IsEmailVerified = true
		employee.Status = ut.Active
		// employee.Username = username

		// Validate employee input
		err = emp.Validate(employee, r)
		if err != nil {
			validationFields := emp.ValidationFields{}
			b, _ := json.Marshal(err)
			// Respond with an errortranslated
			resp := ut.Message(false, ut.Translate(tParam, r))
			json.Unmarshal(b, &validationFields)
			resp["error"] = validationFields
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		// Create a employee
		newEmployee, errParam, err := emps.CreateEmployee(employee)
		if err != nil {
			fmt.Println("Error:: err")
			cErr := ut.CheckUniqueError(r, err)
			if cErr != nil {
				ut.ErrorRespond(http.StatusBadRequest, w, r, ut.Message(false, cErr.Error()))
				return
			}
			ut.ErrorRespond(http.StatusBadRequest, w, r, ut.Message(false, ut.Translate(errParam, r)))
			return
		}
		tParam = tr.TParam{
			Key:          "general.registration_success",
			TemplateData: nil,
			PluralCount:  nil,
		}
		// fmt.Println(sendGridKey)

		// employeeName := newEmployee.FirstName //fmt.Sprintf("%s %s", newEmployee.FirstName, newEmployee.LastName)
		// Set up Email Data
		// emailText := "Thank you for being part of our Team. Please click the link below to confirm your email address and view your account"
		// emailData := ut.EmailData{
		// 	To: []*mail.Email{
		// 		mail.NewEmail(employeeName, newEmployee.Email),
		// 	},
		// 	PageTitle:     "Email Verification",
		// 	Subject:       "Email Verification: Welcome Aboard!",
		// 	Preheader:     "Employee Account Created! ",
		// 	BodyTitle:     fmt.Sprintf("Welcome, %s", employeeName),
		// 	FirstBodyText: emailText,
		// }
		// emailData.Button.Text = "Verify Email"
		// emailData.Button.URL = fmt.Sprintf("%s/verify-email/%s/%s", clientURL, newEmployee.ID, newEmployee.EmailToken)

		// // Send A Welcome/Verification Email to Employee
		// emailBody := ut.ProcessEmail(emailData)
		// go ut.SendEmail(emailBody, sendGridKey)

		// Decode to json so it can be used in the activity log
		decoded, _ := json.Marshal(newEmployee)

		// Log activity
		aLog := md.ActivityLog{
			AdminID:     adminID,
			AppLocation: "Salary Update Function",
			Action:      "Created employee *" + string(decoded) + "*",
		}
		defer acs.CreateActivityLog(aLog)

		// resp := ut.Message(true, ut.Translate(tParam, r))
		resp := ut.Message(true, ut.Translate(tParam, r))
		// resp["verify_url"] = emailData.Button.URL
		ut.Respond(w, r, resp)

	}

}

// UpdateEmployeeEndpoint - Update authenticated employee account
func UpdateEmployeeEndpoint(emps emp.EmployeeService, acs act.ActivityLogService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get Admin Token Data
		tokenData := r.Context().Value("tokenData").(*md.AdminTokenData)
		adminID := string(tokenData.AdminID)

		employeeID := chi.URLParam(r, "employeeID")
		employee := md.Employee{}

		err := json.NewDecoder(r.Body).Decode(&employee)

		tParam := tr.TParam{
			Key:          "error.request_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		if err != nil {
			// Respond with an errortra nslated

			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		// Validate employee input
		err = emp.ValidateUpdates(employee, r)
		fmt.Println("ValidateUpdates ", err)

		if err != nil {
			validationFields := emp.ValidationFields{}
			b, _ := json.Marshal(err)
			// Respond with an errortranslated
			resp := ut.Message(false, ut.Translate(tParam, r))
			json.Unmarshal(b, &validationFields)
			resp["error"] = validationFields
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}
		// Get the employee
		emp := emps.GetEmployee(employeeID)

		formerEmployerData := emp

		emp.FirstName = employee.FirstName
		emp.LastName = employee.LastName
		emp.Email = employee.Email
		emp.Address = employee.Address
		emp.About = employee.About
		emp.Phone = employee.Phone
		emp.Position = employee.Position
		emp.AccountName = employee.AccountName
		emp.AccountNumber = employee.AccountNumber
		emp.BankName = employee.BankName
		emp.Avatar = employee.Avatar
		emp.Thumb = employee.Thumb

		// Update a employee
		updatedEmp, errParam, err := emps.UpdateEmployee(emp)
		if err != nil {
			// Check if the error is duplication error
			cErr := ut.CheckUniqueError(r, err)
			if cErr != nil {
				ut.ErrorRespond(http.StatusBadRequest, w, r, ut.Message(false, cErr.Error()))
				return
			}
			// Respond with an errortranslated
			ut.ErrorRespond(http.StatusBadRequest, w, r, ut.Message(false, ut.Translate(errParam, r)))
			return
		}

		tParam = tr.TParam{
			Key:          "general.update_success",
			TemplateData: nil,
			PluralCount:  nil,
		}

		// Decode to json so it can be used in the activity log
		decoded, _ := json.Marshal(updatedEmp)
		decodedFormer, _ := json.Marshal(formerEmployerData)

		// Log activity
		aLog := md.ActivityLog{
			AdminID:     adminID,
			AppLocation: "Salary Update Function",
			Action:      "Updated employee from *" + string(decodedFormer) + "* to *" + string(decoded) + "*",
		}
		defer acs.CreateActivityLog(aLog)

		resp := ut.Message(true, ut.Translate(tParam, r))
		resp["employee"] = updatedEmp
		ut.Respond(w, r, resp)

	}

}

// VerifyEmployeeEmailEndpoint - Verify employee Email
func VerifyEmployeeEmailEndpoint(emps emp.EmployeeService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		employeeID := chi.URLParam(r, "employeeID")
		emailToken := chi.URLParam(r, "emailToken")

		checkEmployee := emps.GetEmployee(employeeID)

		// if employee already exist send pincode to the employee for verification/Authentication
		if checkEmployee == nil {

			tParam := tr.TParam{
				Key:          "error.invalid_token",
				TemplateData: nil,
				PluralCount:  nil,
			}

			resp := ut.Message(true, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		isTokenValid := ut.ValidatePassword(checkEmployee.EmailToken, emailToken)
		if !isTokenValid {

			tParam := tr.TParam{
				Key:          "error.invalid_token",
				TemplateData: nil,
				PluralCount:  nil,
			}

			resp := ut.Message(true, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}
		checkEmployee.IsEmailVerified = true
		checkEmployee.Status = "active"

		emps.UpdateEmployee(checkEmployee)

		tParam := tr.TParam{
			Key:          "general.verification_success",
			TemplateData: nil,
			PluralCount:  nil,
		}

		resp := ut.Message(true, ut.Translate(tParam, r))
		ut.Respond(w, r, resp)
	}
}
