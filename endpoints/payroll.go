package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"

	act "github.com/ebikode/payroll-core/domain/activity_log"
	pyr "github.com/ebikode/payroll-core/domain/payroll"
	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
	ut "github.com/ebikode/payroll-core/utils"
	"github.com/go-chi/chi"
	// "fmt"
)

// GetPayrollEndpoint fetch single payroll
func GetPayrollEndpoint(ps pyr.PayrollService, userType string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		payrollID := chi.URLParam(r, "payrollID")
		var employeeID string
		// if userType == admin then get the employeeId from the request parameter
		if userType == "admin" {
			employeeID = chi.URLParam(r, "employeeID")
		} else {
			// Get Employee Token Data
			tokenData := r.Context().Value("tokenData").(*md.EmployeeTokenData)
			employeeID = string(tokenData.EmployeeID)
		}

		var payroll *md.Payroll
		payroll = ps.GetPayroll(employeeID, payrollID)
		resp := ut.Message(true, "")
		resp["payroll"] = payroll
		ut.Respond(w, r, resp)
	}
}

// GetPayrollReportsEndpoint ...
func GetPayrollReportsEndpoint(ps pyr.PayrollService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		reports := ps.GetPayrollReports()
		resp := ut.Message(true, "")
		resp["payroll_reports"] = reports
		ut.Respond(w, r, resp)
	}
}

// GetPayrollAllMonthAndYearEndpoint Get all Month Year Combo.
func GetPayrollAllMonthAndYearEndpoint(ps pyr.PayrollService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		filters := ps.GetPayrollAllMonthAndYear()
		resp := ut.Message(true, "")
		resp["payroll_filters"] = filters
		ut.Respond(w, r, resp)
	}
}

// GetPayrollsEndpoint Admin Enpoint for fetching all payrolls
func GetPayrollsEndpoint(ps pyr.PayrollService, userType string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		page, limit := ut.PaginationParams(r)

		var payrolls []*md.Payroll
		payrolls = ps.GetPayrolls(page, limit)

		var nextPage int
		if len(payrolls) == limit {
			nextPage = page + 1
		}

		resp := ut.Message(true, "")
		resp["current_page"] = page
		resp["next_page"] = nextPage
		resp["limit"] = limit
		resp["payrolls"] = payrolls
		ut.Respond(w, r, resp)
	}
}

// GetEmployeePayrollsEndpoint Fetch All employee Payrolls Endpoint
func GetEmployeePayrollsEndpoint(ps pyr.PayrollService, userType string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		page, limit := ut.PaginationParams(r)

		var employeeID string
		// if userType == admin then get the employeeId from the request parameter
		if userType == "admin" {
			employeeID = chi.URLParam(r, "employeeID")
		} else {
			// Get Employee Token Data
			tokenData := r.Context().Value("tokenData").(*md.EmployeeTokenData)
			employeeID = string(tokenData.EmployeeID)
		}

		var payrolls []*md.Payroll
		payrolls = ps.GetEmployeePayrolls(employeeID, page, limit)

		var nextPage int
		if len(payrolls) == limit {
			nextPage = page + 1
		}

		resp := ut.Message(true, "")
		resp["current_page"] = page
		resp["next_page"] = nextPage
		resp["limit"] = limit
		resp["payrolls"] = payrolls
		ut.Respond(w, r, resp)
	}
}

// GetPayrollsByMonthYearEndpoint Fetch All employee account Payrolls Endpoint
func GetPayrollsByMonthYearEndpoint(ps pyr.PayrollService, userType string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		month, _ := strconv.ParseUint(chi.URLParam(r, "month"), 10, 64)
		year, _ := strconv.ParseUint(chi.URLParam(r, "year"), 10, 64)

		var payrolls []*md.Payroll
		payrolls = ps.GetPayrollsByMonthYear(uint(month), uint(year))

		resp := ut.Message(true, "")
		resp["payrolls"] = payrolls
		ut.Respond(w, r, resp)
	}
}

// UpdatePayrollStatusEndpoint ...
func UpdatePayrollStatusEndpoint(pys pyr.PayrollService, acs act.ActivityLogService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get user Token Data
		tokenData := r.Context().Value("tokenData").(*md.AdminTokenData)
		adminID := string(tokenData.AdminID)

		// Set payroll translation parameter
		tParam := tr.TParam{
			Key:          "error.request_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		payload := pyr.Payload{}

		err := json.NewDecoder(r.Body).Decode(&payload)

		// If error occurred
		if err != nil {
			// Respond with an error translated
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}
		// update payroll status
		pys.UpdatePayrollStatus(payload.Status, payload.Month, payload.Year)

		tParam = tr.TParam{
			Key:          "general.update_success",
			TemplateData: nil,
			PluralCount:  nil,
		}

		// Log activity
		aLog := md.ActivityLog{
			AdminID:     adminID,
			AppLocation: "Payroll Status Update Function",
			Action:      "Updated Payroll Status to " + payload.Status + " for " + string(payload.Month) + "/" + string(payload.Year),
		}

		defer acs.CreateActivityLog(aLog)

		resp := ut.Message(true, ut.Translate(tParam, r))
		ut.Respond(w, r, resp)
	}

}

// UpdatePayrollPaymentStatusEndpoint ...
func UpdatePayrollPaymentStatusEndpoint(pys pyr.PayrollService, acs act.ActivityLogService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get Admin Token Data
		tokenData := r.Context().Value("tokenData").(*md.AdminTokenData)
		adminID := string(tokenData.AdminID)

		// Set payroll translation parameter
		tParam := tr.TParam{
			Key:          "error.request_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		payload := pyr.Payload{}

		err := json.NewDecoder(r.Body).Decode(&payload)

		// If error occurred
		if err != nil {
			// Respond with an error translated
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusBadRequest, w, r, resp)
			return
		}
		// update payroll status
		pys.UpdatePayrollPaymentStatus(payload.PaymentStatus, payload.Month, payload.Year)

		tParam = tr.TParam{
			Key:          "general.update_success",
			TemplateData: nil,
			PluralCount:  nil,
		}

		// Log activity
		aLog := md.ActivityLog{
			AdminID:     adminID,
			AppLocation: "Payroll Status Update Function",
			Action:      "Updated Payroll Payment Status to " + payload.PaymentStatus + " for " + string(payload.Month) + "/" + string(payload.Year),
		}

		defer acs.CreateActivityLog(aLog)

		resp := ut.Message(true, ut.Translate(tParam, r))
		ut.Respond(w, r, resp)
	}

}
