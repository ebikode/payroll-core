package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	act "github.com/ebikode/payroll-core/domain/activity_log"
	emp "github.com/ebikode/payroll-core/domain/employee"
	"github.com/ebikode/payroll-core/domain/salary"
	slr "github.com/ebikode/payroll-core/domain/salary"
	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
	ut "github.com/ebikode/payroll-core/utils"
	"github.com/go-chi/chi"
)

// GetSalaryEndpoint fetch a single salary
func GetSalaryEndpoint(sls slr.SalaryService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		salaryID, _ := strconv.ParseUint(chi.URLParam(r, "salaryID"), 10, 64)

		var salary *md.Salary
		salary = sls.GetSalary(uint(salaryID))
		resp := ut.Message(true, "")
		resp["salary"] = salary
		ut.Respond(w, r, resp)
	}
}

// GetAdminSalariesEndpoint fetch a single salary
func GetAdminSalariesEndpoint(sls slr.SalaryService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		page, limit := ut.PaginationParams(r)

		salaries := sls.GetSalaries(page, limit)

		var nextPage int
		if len(salaries) == limit {
			nextPage = page + 1
		}

		resp := ut.Message(true, "")
		resp["current_page"] = page
		resp["next_page"] = nextPage
		resp["limit"] = limit
		resp["salaries"] = salaries
		ut.Respond(w, r, resp)
	}

}

// CreateSalaryEndpoint ...
func CreateSalaryEndpoint(sls slr.SalaryService, emps emp.EmployeeService, acs act.ActivityLogService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get Admin Token Data
		tokenData := r.Context().Value("tokenData").(*md.AdminTokenData)
		adminID := string(tokenData.AdminID)
		payload := slr.Payload{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		fmt.Println("second Error check", err)

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

		checkEmployee := emps.GetEmployee(payload.EmployeeID)

		if checkEmployee == nil {
			tParam = tr.TParam{
				Key:          "error.employee_not_found",
				TemplateData: nil,
				PluralCount:  nil,
			}
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		checkSalary := sls.GetSalaryByEmployeeID(payload.EmployeeID)

		if checkSalary != nil {
			tParam = tr.TParam{
				Key:          "error.salary_already_exist",
				TemplateData: nil,
				PluralCount:  nil,
			}
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		// Validate salary input
		err = slr.Validate(payload, r)
		if err != nil {
			validationFields := slr.ValidationFields{}
			fmt.Println("third Error check", validationFields)
			b, _ := json.Marshal(err)
			// Respond with an errortranslated
			resp := ut.Message(false, ut.Translate(tParam, r))
			json.Unmarshal(b, &validationFields)
			resp["error"] = validationFields
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return

		}

		salary := md.Salary{
			EmployeeID: payload.EmployeeID,
			Salary:     payload.Salary,
			Pension:    payload.Pension,
			Paye:       payload.Paye,
			Nsitf:      payload.Nsitf,
			Nhf:        payload.Nhf,
			Itf:        payload.Itf,
		}

		// Create a salary
		newSalary, errParam, err := sls.CreateSalary(salary)
		if err != nil {
			// Check if the error is duplisalaryion error
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
			Key:          "general.resource_created",
			TemplateData: nil,
			PluralCount:  nil,
		}

		// Decode to json so it can be used in the activity log
		decoded, _ := json.Marshal(newSalary)

		// Log activity
		aLog := md.ActivityLog{
			AdminID:     adminID,
			AppLocation: "Salary Creation Function",
			Action:      "Created " + checkEmployee.FirstName + " " + checkEmployee.LastName + " salary *" + string(decoded) + "*",
		}
		defer acs.CreateActivityLog(aLog)

		resp := ut.Message(true, ut.Translate(tParam, r))
		resp["salary"] = newSalary
		ut.Respond(w, r, resp)

	}

}

// UpdateSalaryEndpoint
func UpdateSalaryEndpoint(sls slr.SalaryService, acs act.ActivityLogService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get Admin Token Data
		tokenData := r.Context().Value("tokenData").(*md.AdminTokenData)
		adminID := string(tokenData.AdminID)
		// Translation Param
		tParam := tr.TParam{
			Key:          "error.request_error",
			TemplateData: nil,
			PluralCount:  nil,
		}
		// Parse the salary id param
		salaryID, pErr := strconv.ParseUint(chi.URLParam(r, "salaryID"), 10, 64)
		if pErr != nil || uint(salaryID) < 1 {
			// Respond with an error translated
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}
		salaryPayload := slr.Payload{}
		// dECODE THE REQUEST BODY
		err := json.NewDecoder(r.Body).Decode(&salaryPayload)

		if err != nil {
			// Respond with an error translated
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}

		// Validate salary input
		err = slr.ValidateUpdates(salaryPayload, r)
		if err != nil {
			validationFields := salary.ValidationFields{}
			b, _ := json.Marshal(err)
			// Respond with an errortranslated
			resp := ut.Message(false, ut.Translate(tParam, r))
			json.Unmarshal(b, &validationFields)
			resp["error"] = validationFields
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return

		}
		fmt.Printf("salaryID:: %v \n", uint(salaryID))
		// Get the salary
		salary := sls.GetSalary(uint(salaryID))
		// Set former salary to be used on activity log creation
		formerSalary := salary
		// Assign new values
		salary.Salary = salaryPayload.Salary
		salary.Pension = salaryPayload.Pension
		salary.Paye = salaryPayload.Paye
		salary.Nsitf = salaryPayload.Nsitf
		salary.Nhf = salaryPayload.Nhf
		salary.Itf = salaryPayload.Itf

		// Update a salary
		updatedSalary, errParam, err := sls.UpdateSalary(salary)
		if err != nil {
			// Check if the error is duplisalaryion error
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
		decoded, _ := json.Marshal(updatedSalary)
		decodedFormer, _ := json.Marshal(formerSalary)

		// Log activity
		aLog := md.ActivityLog{
			AdminID:     adminID,
			AppLocation: "Salary Update Function",
			Action:      "Updated " + salary.Employee.FirstName + " " + salary.Employee.LastName + " salary from *" + string(decodedFormer) + "* to *" + string(decoded) + "*",
		}
		defer acs.CreateActivityLog(aLog)

		resp := ut.Message(true, ut.Translate(tParam, r))
		resp["salary"] = updatedSalary
		ut.Respond(w, r, resp)

	}

}
