package utils

import (
	"time"

	md "github.com/ebikode/payroll-core/model"
	"github.com/uniplaces/carbon"
	"golang.org/x/crypto/bcrypt"
)

// ValidatePincode - validates a employee pincode
func ValidatePincode(employee *md.Employee, pincode string) bool {
	if employee == nil || employee.Pincode == "" || employee.IsPincodeUsed {
		return false
	}

	// Check if
	now := carbon.NewCarbon(time.Now().UTC())
	// Add 10 minutes to pincode sent date since the validity is 10 minutes
	pincodeSentDate := carbon.NewCarbon(employee.PincodeSentAt.UTC().Add(10 * time.Minute))
	// Comparing both time . If now is greater than pincodesentDate
	if now.Gt(pincodeSentDate) {

		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(employee.Pincode), []byte(pincode))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return false
	}

	return true
}

// ValidatePassword - validates a employee password
func ValidatePassword(dbPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return false
	}

	return true
}

// CheckAccountExpiration Check Account Expiration Date
// func CheckAccountExpiration(account *md.Account, clientURL, sendGridKey string) bool {
// 	if account == nil || account.Status != "active" {
// 		return true
// 	}

// 	// Check Expiration date
// 	now := carbon.NewCarbon(time.Now().UTC())
// 	expirationDate := carbon.NewCarbon(account.ExpirationDate.UTC())

// 	// Comparing both time . If now is greater than expirationDate
// 	if now.Gt(expirationDate) {

// 		return true
// 	}

// Add 7 days to the current expiration date
// if today's date is greater send an email alert to the employee for expiration
// letting them know their account expires in 7days
// 	expirationDate = carbon.NewCarbon(account.ExpirationDate.UTC())

// 	// Comparing both time . If now is greater than expirationDate
// 	if now.Gt(expirationDate.AddDays(7)) {
// 		employeeName := fmt.Sprintf("%s %s", account.Employee.FirstName, account.Employee.LastName)
// 		// Set up Email Data
// 		emailText := "The above account expires in less than 7 days. Please visit your account page to renew you Setup/Subscription"
// 		emailSubject := fmt.Sprintf("%s Expiraton")
// 		emailData := EmailData{
// 			To: []*mail.Email{
// 				mail.NewEmail(employeeName, account.Employee.Email),
// 				mail.NewEmail(account.Name, account.Email),
// 			},
// 			PageTitle:     emailSubject,
// 			Subject:       emailSubject,
// 			Preheader:     "in less than 7 days",
// 			BodyTitle:     account.Name,
// 			FirstBodyText: emailText,
// 		}
// 		emailData.Button.Text = "Goto Account"
// 		emailData.Button.URL = fmt.Sprintf("%s/account/%s", clientURL, account.ID)

// 		// Send A Welcome/Verification Email to Employee
// 		emailBody := ProcessEmail(emailData)
// 		go SendEmail(emailBody, sendGridKey)

// 	}

// 	return false
// }

// IsAccountMoreThanAYear Check if account is more than a year
// this is neccessary to bill startups at normal rates
// after their first year
// func IsAccountMoreThanAYear(account *md.Account) bool {
// 	if account == nil {
// 		return true
// 	}

// 	// Check Expiration date
// 	now := carbon.NewCarbon(time.Now().UTC())
// 	createdAtDate := carbon.NewCarbon(account.CreatedAt.UTC())

// 	// Comparing both time . If now is greater than expirationDate
// 	if now.Gt(createdAtDate.AddYear()) {

// 		return true
// 	}

// 	return false
// }
