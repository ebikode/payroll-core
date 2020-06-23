package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	adm "github.com/ebikode/payroll-core/domain/admin"
	aud "github.com/ebikode/payroll-core/domain/authd_device"
	usr "github.com/ebikode/payroll-core/domain/employee"
	pyr "github.com/ebikode/payroll-core/domain/payroll"
	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
	ut "github.com/ebikode/payroll-core/utils"
)

// AuthenticateEmployeeEndpoint Authenticate a employee
func AuthenticateEmployeeEndpoint(appSecret string, us usr.EmployeeService, pys pyr.PayrollService, au aud.AuthdDeviceService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		employee := &md.Employee{}

		err := json.NewDecoder(r.Body).Decode(employee)

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
		employee, tParam, err = us.AuthenticateEmployee(employee.Email, employee.Password)

		if err != nil {
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusUnauthorized, w, r, resp)
			return
		}

		// Detect and save the device the employee used in loggin in
		deviceInfo := ut.DetectDevice(r)
		ip := deviceInfo.IP.To4().String()

		aDevice := md.AuthdDevice{
			EmployeeID:     employee.ID,
			IP:             ip,
			Browser:        deviceInfo.Browser,
			BrowserVersion: deviceInfo.BrowserVersion,
			Platform:       deviceInfo.Platform,
			DeviceOS:       deviceInfo.DeviceOS,
			OSVersion:      deviceInfo.OSVersion,
			DeviceType:     deviceInfo.Type,
			AccessType:     "mobile_app",
			Status:         "active",
		}

		audID := ut.RandomBase64String(8, "MDdv")

		aDevice.ID = audID

		audd, _, _ := au.CreateAuthdDevice(aDevice)

		payrolls := []*md.Payroll{}

		dashoardData := us.GetEmployeeDashboardData(employee.ID)

		payrolls = pys.GetEmployeePayrolls(employee.ID, 1, 20)

		// Create JWT token for application
		tk := &md.EmployeeTokenData{
			EmployeeID: employee.ID,
			DeviceID:   audd.ID,
			Username:   employee.Username,
			ExpireOn:   time.Now().Add(time.Duration(31536000)).UTC(),
		}

		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(appSecret))

		tParam = tr.TParam{
			Key:          "general.login_success",
			TemplateData: nil,
			PluralCount:  nil,
		}
		resp := ut.Message(true, ut.Translate(tParam, r))
		resp["token"] = tokenString
		resp["employee"] = employee
		resp["dashboard_data"] = dashoardData
		resp["recent_payrolls"] = payrolls
		ut.Respond(w, r, resp)
	}
}

// AuthenticateAdminEndpoint - Authenticate admin
func AuthenticateAdminEndpoint(appSecret string, as adm.AdminService, pys pyr.PayrollService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		admin := &md.Admin{}

		err := json.NewDecoder(r.Body).Decode(admin)

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
		admin, tParam, err = as.AuthenticateAdmin(admin.Email, admin.Password)

		if err != nil {
			resp := ut.Message(false, ut.Translate(tParam, r))
			ut.ErrorRespond(http.StatusUnauthorized, w, r, resp)
			return
		}
		payrolls := []*md.Payroll{}

		dashoardData := as.GetAdminDashboardData()

		lastPayroll := pys.GetLastPayroll()

		if lastPayroll != nil {
			payrolls = pys.GetPayrollsByMonthYear(lastPayroll.Month, lastPayroll.Year)
		}

		//Create JWT token
		tk := &md.AdminTokenData{
			AdminID:  admin.ID,
			Role:     admin.Role,
			ExpireOn: time.Now().Add(time.Duration(86400)).UTC(),
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(appSecret))

		tParam = tr.TParam{
			Key:          "general.admin_login_success",
			TemplateData: nil,
			PluralCount:  nil,
		}
		resp := ut.Message(true, ut.Translate(tParam, r))
		resp["token"] = tokenString
		resp["admin"] = admin
		resp["dashboard_data"] = dashoardData
		resp["recent_payrolls"] = payrolls
		ut.Respond(w, r, resp)
	}

}
