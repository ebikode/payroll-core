package endpoints

import (
	"net/http"
	"strconv"

	tx "github.com/ebikode/payroll-core/domain/tax"
	md "github.com/ebikode/payroll-core/model"
	ut "github.com/ebikode/payroll-core/utils"
	"github.com/go-chi/chi"
)

// GetTaxEndpoint fetch a single tax
func GetTaxEndpoint(txs tx.TaxService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		taxID, _ := strconv.ParseUint(chi.URLParam(r, "taxID"), 10, 64)

		var tax *md.Tax
		tax = txs.GetTax(uint(taxID))
		resp := ut.Message(true, "")
		resp["tax"] = tax
		ut.Respond(w, r, resp)
	}
}

// GetAdminTaxesEndpoint fetch a single tax
func GetAdminTaxesEndpoint(txs tx.TaxService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		page, limit := ut.PaginationParams(r)

		taxes := txs.GetTaxes(page, limit)

		var nextPage int
		if len(taxes) == limit {
			nextPage = page + 1
		}

		resp := ut.Message(true, "")
		resp["current_page"] = page
		resp["next_page"] = nextPage
		resp["limit"] = limit
		resp["taxes"] = taxes
		ut.Respond(w, r, resp)
	}

}

// GetEmployeeTaxesEndpoint Get AllTaxesEndpoint
func GetEmployeeTaxesEndpoint(txs tx.TaxService, userType string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var employeeID string
		if userType == "admin" {
			employeeID = chi.URLParam(r, "employeeID")
		}
		if userType == "employee" {
			tokenData := r.Context().Value("tokenData").(*md.EmployeeTokenData)
			employeeID = string(tokenData.EmployeeID)
		}

		var taxes []*md.Tax
		taxes = txs.GetEmployeeTaxes(employeeID)
		resp := ut.Message(true, "")
		resp["taxes"] = taxes
		ut.Respond(w, r, resp)
	}
}
