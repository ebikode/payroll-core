package jobs

import (
	"fmt"
	"strconv"
	"time"

	apset "github.com/ebikode/payroll-core/domain/app_setting"
	emp "github.com/ebikode/payroll-core/domain/employee"
	pyr "github.com/ebikode/payroll-core/domain/payroll"
	tx "github.com/ebikode/payroll-core/domain/tax"
	md "github.com/ebikode/payroll-core/model"
	ut "github.com/ebikode/payroll-core/utils"
)

var employees []*md.Employee

// RunPayrollGenerationJob - Automated Payroll Generation
func RunPayrollGenerationJob(
	pys pyr.PayrollService, aps apset.AppSettingService,
	emps emp.EmployeeService, txs tx.TaxService,
) {

	genDay := aps.GetAppSettingByKey(ut.PayrollGenerationDayKey, "admin")

	generationDate, _ := strconv.ParseInt(genDay.Value, 10, 64)

	now := time.Now()

	todayDate := now.Day()
	todayMonth := int(now.Month())
	todayYear := now.Year()

	if int(generationDate) >= todayDate {
		lastPayroll := pys.GetLastPayroll()

		if lastPayroll == nil || (int(lastPayroll.Month) != todayMonth && int(lastPayroll.Year) != todayYear) {
			employees := emps.GetAllActivePubEmployee()

			fmt.Println("Payroll Generation Automation Started")

			for _, v := range employees {

				netSalary, tax := generateNetSalaryAndTaxDeductions(v.Salary)

				payroll := md.Payroll{
					EmployeeID:  v.ID,
					GrossSalary: v.Salary.Salary,
					NetSalary:   netSalary,
					Month:       uint(todayMonth),
					Year:        uint(todayYear),
				}

				newPayroll, _, err := pys.CreatePayroll(payroll)

				if err == nil {
					tax.PayrollID = newPayroll.ID
					txs.CreateTax(tax)
				}

			}

			fmt.Println("Payroll Generation Automation Ended")

		}
	}

}

func generateNetSalaryAndTaxDeductions(salary *md.Salary) (float64, md.Tax) {

	GrossSalary := salary.Salary

	tax := md.Tax{
		Pension: percentage(salary.Pension) * GrossSalary,
		Paye:    percentage(salary.Paye) * GrossSalary,
		Nsitf:   percentage(salary.Nsitf) * GrossSalary,
		Nhf:     percentage(salary.Nhf) * GrossSalary,
		Itf:     percentage(salary.Itf) * GrossSalary,
	}

	deductions := tax.Pension + tax.Paye + tax.Nsitf + tax.Nhf + tax.Itf

	netSalary := GrossSalary - deductions

	return netSalary, tax

}

func percentage(value float64) float64 {
	return (value / 100)
}
